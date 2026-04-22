# Tech Stack for Full-Stack Developer Test

เอกสารนี้สรุปชุดเทคโนโลยีที่เลือกใช้สำหรับโจทย์ทดสอบ `Single-Page Application` ประเภท `image gallery` ที่มี `infinite scroll`, `hashtag filter` และแนวทาง deploy แบบ production

## Goal

- พัฒนา `SPA` ที่แสดงรูปภาพพร้อม hashtag
- รองรับ `infinite scroll`
- กด hashtag เพื่อกรองรูปภาพได้
- ออกแบบระบบให้สามารถ deploy ใช้งานจริงได้

## Selected Stack

- Frontend: `React`, `Vite`, `TypeScript`, `Tailwind CSS`
- Backend: `Go`, `Fiber`
- Database: `MySQL 8`
- Reverse Proxy / Static Server: `Nginx`
- Containerization: `Docker`, `Docker Compose`
- CI/CD: `GitHub Actions`
- Deployment OS: `Ubuntu Server 22.04 LTS`
- Demo image source: `placehold.co`

## Why This Stack

- `React` เหมาะกับงาน `SPA` และจัดการ state ของ gallery/filter ได้ชัดเจน
- `Vite` ช่วยให้เริ่มโปรเจกต์เร็วและ build เร็ว เหมาะกับงานทดสอบที่มีเวลาจำกัด
- `TypeScript` ลดความผิดพลาดของ data shape ระหว่าง frontend และ backend
- `Go Fiber` ตรงกับทักษะ `Golang` ในตำแหน่งงาน และเหมาะกับการทำ API ที่เรียบง่ายและเร็ว
- `MySQL` ตรงกับ requirement ในเอกสาร และเหมาะกับโครงสร้างข้อมูลแบบ relation ระหว่างรูปกับ tag
- `Nginx` ใช้ serve ไฟล์ frontend และ reverse proxy ไปยัง backend ได้ง่าย
- `Docker` และ `GitHub Actions` ทำให้ workflow ดู production-ready และ deploy ซ้ำได้ง่าย
- `Ubuntu Server` เป็นตัวเลือกมาตรฐานสำหรับ production deployment

## Application Design

### Frontend

- แสดงรูปแบบ masonry/grid ตามขนาดรูปที่ไม่เท่ากัน
- โหลดข้อมูลเพิ่มเมื่อ scroll ถึงด้านล่างด้วย `IntersectionObserver`
- กด hashtag แล้วเรียก API ใหม่ด้วยเงื่อนไข `tag`
- รองรับ responsive สำหรับ desktop และ mobile

### Backend

- พัฒนาเป็น `REST API`
- แยกส่วน `handler`, `service`, `repository` เพื่อให้โค้ดอ่านง่ายและดูแลต่อได้
- รองรับ query เช่น `limit`, `cursor`, `tag`

ตัวอย่าง endpoint:

- `GET /api/images`
- `GET /api/images?limit=12`
- `GET /api/images?cursor=24`
- `GET /api/images?tag=nature&limit=12`

### Database

โครงสร้างหลัก:

- `images`
- `tags`
- `image_tags`

แนวคิด:

- 1 รูปมีได้หลาย tag
- 1 tag ใช้กับหลายรูปได้
- ใช้ตาราง `image_tags` เป็น many-to-many relation

## Suggested Production Setup

- Server: `2 vCPU`, `4 GB RAM`, `40 GB SSD`
- OS: `Ubuntu Server 22.04 LTS`
- Runtime: `Docker Engine`, `Docker Compose`
- Web entrypoint: `Nginx`
- App services:
  - `frontend`
  - `backend`
  - `mysql`
- CI/CD flow:
  - push code ไปที่ `GitHub`
  - `GitHub Actions` run build/test
  - deploy ไป server ด้วย `docker compose`

## Deployment Flow

1. Build frontend ด้วย `Vite`
2. Build backend เป็น `Go binary`
3. สร้าง Docker images สำหรับ frontend และ backend
4. ใช้ `Nginx` serve frontend และ proxy `/api` ไป backend
5. ใช้ `docker compose up -d` สำหรับ run production services

## Final Recommendation

หากต้องการ stack ที่บาลานซ์ระหว่างความเร็วในการพัฒนา, ความตรงกับ JD และภาพลักษณ์ production-ready ชุดที่เหมาะที่สุดคือ:

- `React + Vite + TypeScript`
- `Go Fiber`
- `MySQL`
- `Nginx`
- `Docker / Docker Compose`
- `GitHub Actions`
- `Ubuntu Server`
