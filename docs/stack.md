# Current Production Stack

เอกสารนี้สรุปสถาปัตยกรรมและเทคโนโลยีที่ใช้อยู่ใน repo ปัจจุบันสำหรับระบบ `image gallery` แบบ `SPA` ที่มี `infinite scroll`, `tag filter`, หน้า `manage`, backend API และ MySQL

## Goal

- พัฒนา `SPA` สำหรับแสดงรูปภาพและจัดการรูปภาพในระบบ
- รองรับ `infinite scroll`
- รองรับการกรองข้อมูลด้วย `tag`
- รองรับการจัดการข้อมูลผ่านหน้า `/manage`
- รองรับการ deploy แบบแยก `frontend`, `backend`, `database`

## Current Stack

- Frontend:
  - `React 18.3.1`
  - `Vite 5.4.21`
  - `TypeScript 5.6.3`
  - `Tailwind CSS 3.4.15`
  - `react-router-dom 6.28.x`
  - `masonry-layout 4.2.2`
  - `imagesloaded 5.0.0`
- Backend:
  - `Go 1.26`
  - `Fiber v2.52.9`
  - `go-sql-driver/mysql 1.9.3`
- Database:
  - `MySQL 8.4`
- Frontend runtime image:
  - `nginx:1.27-alpine`
- Backend runtime image:
  - `alpine:3.21`
- Backend build image:
  - `golang:1.26-alpine`
- Containerization:
  - `Docker`
  - `Dockerfile`
- Local orchestration:
  - `Docker Compose`
- Current cloud deployment target:
  - `Railway`
- Demo image source:
  - `placehold.co`

## Why This Stack

- `React` เหมาะกับงาน `SPA` และจัดการ state ของ gallery / filter / form ได้ชัดเจน
- `Vite` ใช้ build frontend ได้เร็วและเหมาะกับ workflow ของ React ปัจจุบัน
- `TypeScript` ช่วยคุม data shape ระหว่าง frontend กับ backend
- `Go Fiber` เหมาะกับการทำ API ที่เรียบง่าย เร็ว และดูแลง่าย
- `MySQL` เหมาะกับ relational model ของ `images`, `tags`, `image_tags`
- `Nginx` เหมาะกับการ serve static frontend ใน production
- `Docker` ทำให้ build และ deploy ซ้ำได้สม่ำเสมอ
- `Railway` เหมาะกับ repo นี้เพราะแยก deploy เป็น `frontend`, `backend`, `mysql` ได้ง่ายและเชื่อมกับ GitHub ได้ตรง

## Application Design

### Frontend

- แสดงรูปแบบ masonry/grid ตามขนาดรูปที่ไม่เท่ากัน
- โหลดข้อมูลเพิ่มเมื่อ scroll ถึงด้านล่างด้วย `IntersectionObserver`
- กรองข้อมูลด้วย `tag`
- มีหน้า `/manage` สำหรับจัดการ image, tag และการผูก tag

### Backend

- พัฒนาเป็น `REST API`
- แยกส่วน `handler`, `service`, `repository`
- รองรับ query เช่น `limit`, `cursor`, `tag`

ตัวอย่าง endpoint:

- `GET /api/images`
- `GET /api/images?limit=12`
- `GET /api/images?cursor=24`
- `GET /api/images?tag=nature&limit=12`
- `POST /api/images`
- `POST /api/tags`

### Database

โครงสร้างหลัก:

- `images`
- `tags`
- `image_tags`

แนวคิด:

- 1 รูปมีได้หลาย tag
- 1 tag ใช้กับหลายรูปได้
- ใช้ `deleted_at` สำหรับ soft delete
- ใช้ `image_tags` เป็น many-to-many relation

## Current Production-Like Setup

โปรเจคปัจจุบัน deploy แบบแยก service บน `Railway` เป็น 3 ส่วน:

- `frontend`
- `backend`
- `mysql`

รายละเอียดเชิง runtime ปัจจุบัน:

- `frontend`
  - build จาก `frontend/Dockerfile`
  - frontend ถูก build ด้วย `Vite`
  - serve ผ่าน `nginx:1.27-alpine`
  - public entrypoint ของระบบ
- `backend`
  - build จาก `backend/Dockerfile`
  - compile ด้วย `golang:1.26-alpine`
  - run บน `alpine:3.21`
  - เปิด API ที่ port `8080`
- `mysql`
  - local compose ใช้ `mysql:8.4`
  - cloud ใช้ Railway MySQL service

## OS / Software Version

ถ้าอธิบายให้ตรงกับระบบปัจจุบันใน repo:

- Host OS:
  - ในกรณี deploy บน `Railway` เป็น platform-managed environment จึงไม่ได้ pin OS ของ host machine เองใน repo
- Container OS / Runtime ที่ระบบใช้จริง:
  - Frontend runtime: `nginx:1.27-alpine`
  - Backend runtime: `alpine:3.21`
  - Backend build image: `golang:1.26-alpine`
  - Database (local/reference): `mysql:8.4`
- Application software versions:
  - `React 18.3.1`
  - `Vite 5.4.21`
  - `TypeScript 5.6.3`
  - `Go 1.26`
  - `Fiber v2.52.9`
  - `MySQL 8.4`
  - `Nginx 1.27`

## Current Deployment Flow

ปัจจุบัน flow ที่ตรงกับ repo นี้คือ:

1. push code ไปที่ `GitHub`
2. Railway ดึง source จาก branch ที่ผูกไว้
3. Railway build `backend` จาก `backend/Dockerfile`
4. Railway build `frontend` จาก `frontend/Dockerfile`
5. backend start พร้อม env เช่น `AUTO_MIGRATE=true` และ `MOCKDATA=true/false`
6. frontend build ด้วย `VITE_API_BASE_URL=https://<backend-domain>/api`
7. เปิดใช้งานผ่าน Railway public domains

## CI/CD Status

ปัจจุบันใน repo นี้:

- ยังไม่มี `.github/workflows/`
- ยังไม่มี `GitHub Actions` workflow ใน repo
- deployment ปัจจุบันอาศัย `Railway + GitHub integration`
- เมื่อ push code ไปยัง branch ที่เชื่อมกับ Railway ระบบจะ build และ deploy service ที่เกี่ยวข้องใหม่

ดังนั้นถ้าจะเขียนในรายงานให้ตรงกับ implementation ปัจจุบัน ควรใช้คำว่า:

- `Deployment trigger: Git push -> Railway auto deploy`
- ไม่ควรระบุว่าใช้ `GitHub Actions` ถ้ายังไม่ได้เพิ่ม workflow จริงใน repo

## Recommended Production Notes

- `frontend` เป็น public service
- `backend` เป็น public API service ที่ frontend เรียกผ่าน `VITE_API_BASE_URL`
- `mysql` ควรเป็น private/internal only
- backend ใช้ `AUTO_MIGRATE=true` ได้ใน environment ที่ต้องการให้ schema ถูก apply ตอน start
- ใน production จริงควรจำกัดการเข้าถึง `/swagger` เพิ่มเติมหากไม่ต้องการเปิดสาธารณะ
