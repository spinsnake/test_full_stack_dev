## Production Architecture Topology

ใน production ระบบจะถูกแยกออกเป็น 3 ส่วนหลัก คือ frontend, backend และ database โดยกำหนดให้มีเฉพาะ service ที่จำเป็นเท่านั้นที่เปิดให้เข้าถึงจากภายนอก เพื่อเพิ่มความปลอดภัยและลดพื้นที่เสี่ยงของระบบ

โครงสร้างการเข้าถึงจากภายนอก

ผู้ใช้ (User) จะเข้าถึงระบบผ่าน Frontend เป็นหลัก

Frontend ทำหน้าที่เป็น public entrypoint สำหรับการใช้งานผ่าน browser

Backend เปิดให้เข้าถึงผ่าน HTTP/HTTPS เฉพาะสำหรับการเรียก API

Database จะไม่เปิด public access และอนุญาตให้เชื่อมต่อได้เฉพาะจาก Backend ภายใน private network เท่านั้น

บทบาทของแต่ละส่วนใน production

Frontend

ทำหน้าที่แสดงผล UI ของระบบ

build เป็น static asset ด้วย Vite

serve ผ่าน Nginx

ติดต่อกับ backend ผ่าน REST API โดยใช้ VITE_API_BASE_URL

Backend

ทำหน้าที่เป็น API service ของระบบ

รับ request จาก frontend เช่น การดึงรูปภาพ, tag, การเพิ่ม/แก้ไข/ลบข้อมูล

ประมวลผล business logic และเชื่อมต่อกับฐานข้อมูล

สามารถเปิด public endpoint สำหรับ API และ Swagger ได้ตาม environment ที่ใช้งาน

Database

ใช้ MySQL สำหรับเก็บข้อมูลรูปภาพ, tag และความสัมพันธ์ระหว่างข้อมูล

อยู่ใน private/internal network

ไม่เปิดให้ user หรือ browser เข้าถึงโดยตรง

รับการเชื่อมต่อจาก backend เท่านั้น

Traffic Flow
ลำดับการทำงานของระบบใน production เป็นดังนี้

ผู้ใช้เปิดเว็บไซต์ผ่าน Frontend

Nginx ทำหน้าที่ serve static frontend ให้กับ browser

เมื่อ frontend ต้องการข้อมูล เช่น รายการรูปภาพหรือ tag จะส่ง request ไปยัง Backend API

Backend รับ request และประมวลผล business logic

Backend query หรือ update ข้อมูลใน MySQL

MySQL ส่งผลลัพธ์กลับไปยัง Backend

Backend ส่ง response กลับไปยัง Frontend

Frontend นำข้อมูลมาแสดงผลให้ผู้ใช้เห็น

การเปิด Public / Private Access

Frontend : Public

Backend API : Public ตามรูปแบบ deployment

Swagger UI : ควรเปิดเฉพาะ environment สำหรับ development / staging หรือจำกัดสิทธิ์ใน production

Database : Private only

# Database

Database Schema
ประกอบด้วย table ต่อไปนี้
1. images ทำหน้าที่เก็บข้อมูลรูปภาพ

| name | type | note |
| --- | --- | --- |
| id | BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY | Primary Key ใช้สำหรับอ้างอิงข้อมูลแต่ละรูป |
| image_url | VARCHAR(512) NOT NULL | URL ของรูปภาพหลักที่ใช้แสดงหรือเก็บอ้างอิงรูปจริง |
| thumbnail_url | VARCHAR(512) NULL | URL ของรูป thumbnail สำหรับใช้แสดงใน gallery ให้โหลดเร็วขึ้น |
| width | INT UNSIGNED NULL | ความกว้างของรูปภาพ หน่วยเป็น pixel |
| height | INT UNSIGNED NULL | ความสูงของรูปภาพ หน่วยเป็น pixel |
| alt_text | VARCHAR(255) NULL | คำอธิบายรูป |
| source | VARCHAR(100) NULL | แหล่งที่มาของรูป เช่น seed:placehold.co |
| deleted_at | DATETIME(3) NULL | เวลาในการลบแบบ soft delete ถ้าเป็น NULL แปลว่ายังใช้งานอยู่ |
| created_at | DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) | วันเวลาที่สร้างข้อมูลรูปภาพนี้ |
| updated_at | DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3) | วันเวลาที่แก้ไขข้อมูลล่าสุด |

2. tags ทำหน้าที่เก็บข้อมูลคำสำคัญ

| name | type | note |
| --- | --- | --- |
| id | BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY | รหัสประจำ tag ใช้เป็น Primary Key สำหรับอ้างอิงแต่ละ tag |
| name | VARCHAR(64) NOT NULL | ชื่อ tag ที่ใช้แสดงผลให้ผู้ใช้เห็น เช่น Nature, Travel |
| slug | VARCHAR(64) NOT NULL | ค่าของ tag ในรูปแบบที่เหมาะกับการใช้ในระบบและ query เช่น nature, travel |
| deleted_at | DATETIME(3) NULL | เวลาในการลบแบบ soft delete ถ้าเป็น NULL แปลว่า tag นี้ยังใช้งานอยู่ |
| active_name | VARCHAR(64) GENERATED ALWAYS AS (...) STORED | generated column ที่เก็บค่า name เฉพาะกรณีที่ยังไม่ถูกลบ ใช้ช่วยบังคับไม่ให้ชื่อ tag ซ้ำกันในข้อมูลที่ยัง active |
| active_slug | VARCHAR(64) GENERATED ALWAYS AS (...) STORED | generated column ที่เก็บค่า slug เฉพาะกรณีที่ยังไม่ถูกลบ ใช้ช่วยบังคับไม่ให้ slug ซ้ำกันในข้อมูลที่ยัง active |
| created_at | DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) | วันเวลาที่สร้างข้อมูล tag นี้ |
| updated_at | DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3) | วันเวลาที่แก้ไขข้อมูลล่าสุด โดยระบบจะอัปเดตให้อัตโนมัติ |

3. image_tags เก็บข้อมูลที่ map ระหว่าง image และ tags โดยความสัมพันธ์จะเป็น many to many

| name | type | note |
| --- | --- | --- |
| image_id | BIGINT UNSIGNED NOT NULL | รหัสของรูปภาพ อ้างอิงไปยัง images.id เพื่อระบุว่า tag นี้ผูกกับรูปไหน |
| tag_id | BIGINT UNSIGNED NOT NULL | รหัสของ tag อ้างอิงไปยัง tags.id เพื่อระบุว่ารูปนี้ถูกผูกกับ tag อะไร |
| created_at | DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) | วันเวลาที่มีการผูก tag นี้เข้ากับรูปภาพ |

4. schema_migrations เก็บสถานะการ migration หากมีการเรียกใช้คำสั่ง migration up / down ว่าอยู่ในสถานะไหน

| name | type | note |
| --- | --- | --- |
| version | BIGINT NOT NULL | หมายเลข version ของ migration ล่าสุดที่ถูก apply กับฐานข้อมูล ใช้บอกว่าตอนนี้ schema อยู่ที่ขั้นไหน |
| dirty | BOOLEAN NOT NULL | ใช้บอกว่าการ migrate ล่าสุดสมบูรณ์หรือไม่ ถ้าเป็น true แปลว่า migration ค้างหรือพังกลางทาง ถ้าเป็น false แปลว่า migration ล่าสุดเสร็จสมบูรณ์ |

ภาพแสดงส่วน database schema

โฟลเดอร์ที่เก็บโครงสร้างแต่ละ tables จะอยู่ในโฟล์เดอร์ backend/migration/

Database Migration

ระบบใช้ migration สำหรับจัดการโครงสร้างฐานข้อมูลและ seed data ของ MySQL โดยไฟล์ migration ถูกเก็บไว้ในโฟลเดอร์ backend/migration/ และแยกเป็นไฟล์ up กับ down

*.up.sql ใช้สำหรับสร้างตารางหรือเพิ่มข้อมูล

*.down.sql ใช้สำหรับ rollback การเปลี่ยนแปลง

ไฟล์ migration หลักของระบบมีดังนี้

000001_create_gallery_tables.up.sql สำหรับสร้างตารางหลักของระบบ

000001_create_gallery_tables.down.sql สำหรับลบตารางหลัก

000002_seed_demo_data.up.sql สำหรับเพิ่มข้อมูลตัวอย่าง

000002_seed_demo_data.down.sql สำหรับลบข้อมูลตัวอย่าง

การทำงานของ Migration

up ใช้สำหรับ apply migration จาก version ปัจจุบันไปยัง version ใหม่กว่า

down ใช้สำหรับ rollback migration ย้อนกลับตามจำนวน step ที่กำหนด

วิธีรัน Migration
โปรเจคนี้มี script ช่วยรัน migration อยู่ที่ backend/scripts/migrate.ps1

ตัวอย่างคำสั่ง

```text
powershell -ExecutionPolicy Bypass -File backend/scripts/migrate.ps1 up
```

ใช้สำหรับ apply migration ทั้งหมดที่ยังไม่ได้รัน

```text
powershell -ExecutionPolicy Bypass -File backend/scripts/migrate.ps1 down 1
```

ใช้สำหรับ rollback ย้อนกลับ 1 step

```text
powershell -ExecutionPolicy Bypass -File backend/scripts/migrate.ps1 version
```

ใช้สำหรับตรวจสอบ version ปัจจุบันของฐานข้อมูล

ผลลัพธ์หลังรัน

เมื่อรัน up ระบบจะสร้าง schema หลักของ gallery และเพิ่ม demo data ตาม migration ที่กำหนด

เมื่อรัน down ระบบจะย้อนการเปลี่ยนแปลงตามลำดับ migration ที่ rollback

# Frontend

Frontend

แบ่งออกเป็น 2 หน้าหลัก ได้แก่

หน้า Gallery (/)

เป็นหน้าหลักสำหรับแสดงรูปภาพทั้งหมดในระบบ

แสดงผลรูปในลักษณะ masonry layout เพื่อให้รองรับรูปที่มีขนาดไม่เท่ากัน

รองรับ infinite scroll โดยเมื่อผู้ใช้เลื่อนหน้าจอลงมาด้านล่าง ระบบจะโหลดรูปเพิ่มอัตโนมัติ

รองรับการกรองรูปภาพด้วย tag filter

ผู้ใช้สามารถเลือกจำนวนรูปที่โหลดต่อรอบได้ เช่น 10, 30, 50

ใช้สำหรับการดูภาพรวมของข้อมูลทั้งหมดใน gallery

หน้า Manage (/manage)

ใช้สำหรับจัดการข้อมูลภายในระบบ

รองรับการเพิ่ม แก้ไข และลบข้อมูลรูปภาพ

รองรับการเพิ่ม แก้ไข และลบข้อมูล tag

รองรับการผูก tag เข้ากับ image และการยกเลิกการผูก tag

ใช้เป็นหน้าสำหรับ admin หรือผู้ดูแลระบบในการจัดการข้อมูล gallery

Frontend Tech Stack
Frontend ของระบบพัฒนาด้วย React, Vite และ TypeScript สำหรับสร้าง Single Page Application (SPA) โดยใช้ Tailwind CSS สำหรับจัดการ style, react-router-dom สำหรับ routing ภายในระบบ, masonry-layout และ imagesloaded สำหรับจัดการ image gallery แบบ masonry และใช้ Playwright สำหรับ E2E testing ของฝั่ง frontend

โครงสร้างไฟล์หลักของ frontend

| file | หน้าที่ | request ที่เกี่ยวข้อง |
| --- | --- | --- |
| frontend/src/main.tsx | เป็น entry point ของ frontend ทำหน้าที่ mount React app ลงใน root, ครอบแอปด้วย BrowserRouter และโหลด stylesheet กลาง | ไม่มีการยิง request |
| frontend/src/App.tsx | กำหนด route หลักของระบบ โดย map / ไปหน้า gallery และ /manage ไปหน้าจัดการข้อมูล | ไม่มีการยิง request |
| frontend/src/api.ts | เป็นศูนย์กลางของการเรียก API ทั้งหมด กำหนด API_BASE_URL, มี wrapper request() สำหรับจัดการ fetch, response และ error ให้เป็นรูปแบบเดียวกัน | GET /api/tags, GET /api/images, POST /api/images, PATCH /api/images/:id, DELETE /api/images/:id, POST /api/tags, PATCH /api/tags/:id, DELETE /api/tags/:id, POST /api/images/:imageID/tags, DELETE /api/images/:imageID/tags/:tagID |
| frontend/src/types.ts | เก็บ type definition ของข้อมูลที่ใช้ใน frontend เช่น ImageItem, Tag, CreateImagePayload, CreateTagPayload เพื่อให้ UI และ API ใช้ data shape เดียวกัน | ไม่มีการยิง request |
| frontend/src/pages/GalleryPage.tsx | เป็นหน้าหลักของระบบ ใช้แสดง gallery, คุม state ของรูป, tag, filter, chunk size, modal และ infinite scroll | เรียก fetchTags() เพื่อโหลด tag และเรียก fetchImages() เพื่อโหลดรูปภาพตาม limit, cursor, tag |
| frontend/src/pages/ManagePage.tsx | เป็นหน้าจัดการข้อมูลสำหรับ admin ใช้ create, update, delete image และ tag รวมถึง attach / detach tag ให้ image | เรียก fetchAllImages(), fetchTags(), createImage(), updateImage(), deleteImage(), createTag(), updateTag(), deleteTag(), attachTagToImage(), detachTagFromImage() |
| frontend/src/index.css | เก็บ style กลางของระบบ เช่น layout, card, button, form field และ style เฉพาะของ gallery/manage page | ไม่มีการยิง request |
| frontend/vite.config.ts | ตั้งค่า Vite dev server และ proxy สำหรับ local development โดยส่ง /api และ /swagger ไป backend ที่ localhost:8080 | ไม่ยิง request เอง แต่กำหนด proxy ให้ request จาก frontend วิ่งไป backend ได้ตอน dev |
| frontend/index.html | เป็น HTML template หลักของ Vite ที่มี div#root สำหรับ mount React app | ไม่มีการยิง request |
| frontend/nginx.conf | ใช้ตอน deploy frontend ผ่าน Nginx เพื่อ serve static SPA และ fallback route กลับไป index.html | ไม่มีการยิง request เอง |
| frontend/playwright.config.ts | ใช้ตั้งค่า Playwright E2E test เช่น test directory, base URL และ web server | ไม่มีการยิง request ของระบบจริง |
| frontend/e2e/support/mockApi.ts | ใช้ mock API สำหรับ E2E test เพื่อจำลองข้อมูล image และ tag โดยไม่ต้องเรียก backend จริง | intercept request ใน test เช่น /api/images, /api/tags, /api/images/:id/tags |
| frontend/e2e/gallery.spec.ts | ทดสอบ flow หน้า gallery เช่น โหลดรูป, filter tag, เปิดหน้า manage | เรียกผ่าน mock API |
| frontend/e2e/manage.spec.ts | ทดสอบ flow หน้า manage เช่น create image, create tag, attach tag | เรียกผ่าน mock API |

Frontend E2E testing

Frontend ใช้ Playwright สำหรับ E2E testing โดย config อยู่ที่ playwright.config.ts (line 1) และคำสั่งรันอยู่ใน frontend/package.json (line 5)

ชุดทดสอบนี้ใช้ mock API จาก frontend/e2e/support/mockApi.ts (line 1) ทำให้สามารถทดสอบ flow ของ UI ได้โดยไม่ต้องพึ่ง backend จริงหรือ database จริง

การทำงานของ Test

Playwright จะเปิด frontend ผ่าน Vite dev server อัตโนมัติที่ http://127.0.0.1:4173 ตาม playwright.config.ts (line 3)

ทุก request ที่ไป **/api/** จะถูก intercept โดย mock API ใน mockApi.ts (line 1)

ทำให้ test เน้นตรวจ behavior ของ frontend โดยตรง เช่น rendering, filtering, form submission และ state update

วิธีรัน
รันจากโฟลเดอร์ frontend npm run e2e

| file | สิ่งที่ทดสอบ |
| --- | --- |
| frontend/e2e/gallery.spec.ts (line 1) | ทดสอบหน้า gallery ว่าสามารถโหลดรายการรูปเริ่มต้นได้, เปิด Tags Filter, เลือก tag แล้วจำนวนรูปถูกกรองถูกต้อง, และสามารถกดไปหน้า /manage ได้ |
| frontend/e2e/manage.spec.ts (line 1) | ทดสอบหน้า /manage ว่าสามารถกรอกข้อมูลสร้าง image ใหม่ได้, เลือก tag ตอนสร้าง image ได้, กด Create Image แล้ว image ใหม่ถูกเพิ่มใน catalog พร้อม tag ที่เลือกไว้ |

# Backend

Backend API Services

ทำหน้าที่สร้าง REST API เพื่อรองรับการจัดการข้อมูลรูปภาพ, tag และความสัมพันธ์ระหว่างรูปภาพกับ tag โดยฐานข้อมูลที่ใช้คือ MySQL 8 และมีการจัดการ schema ผ่านไฟล์ migration

เทคโนโลยีที่ใช้

Go

Fiber

OpenAPI / Swagger UI

โครงสร้างของ Backend

Backend ถูกออกแบบโดยแยกหน้าที่ของแต่ละส่วนออกจากกันเพื่อให้โค้ดอ่านง่ายและดูแลต่อได้ง่าย โดยแบ่งเป็น

| path | หน้าที่ |
| --- | --- |
| backend/cmd/api | เป็น entry point ของ application ใช้สำหรับ start server และคำสั่งเกี่ยวกับ migration |
| backend/internal/adapter/handler | รับ HTTP request และแปลง request/response ระหว่าง client กับ service |
| backend/internal/service | จัดการ business logic ของระบบ เช่น validation, pagination, filter และ flow ของการผูก tag |
| backend/internal/adapter/repo | ติดต่อฐานข้อมูลโดยตรง เช่น query, insert, update, delete |
| backend/internal/entities | เก็บ struct ของ domain model และ request/response payload |
| backend/internal/port | กำหนด interface ของ service และ repository เพื่อแยก dependency ให้ชัดเจน |
| backend/internal/routes | กำหนดเส้นทาง API และ bind route เข้ากับ handler |
| backend/internal/config | อ่านค่า configuration จาก environment variables |
| backend/internal/infra | จัดการ infrastructure เช่น การเชื่อมต่อ MySQL และการรัน migration |
| backend/migration | เก็บไฟล์ SQL migration และ seed data |
| backend/docs | เก็บ OpenAPI spec และไฟล์ Swagger UI |

แนวทางการทำงานของ Backend
ลำดับการทำงานหลักของ backend คือ

client ส่ง request เข้ามาที่ route

route ส่งต่อไปยัง handler

handler parse request และเรียก service

service จัดการ business logic และเรียก repository

repository query หรือ update ข้อมูลใน MySQL

handler ส่ง response กลับในรูปแบบ JSON

ตัวอย่าง API หลักที่ระบบรองรับ

```text
GET /api/images
POST /api/images
PATCH /api/images/:imageID
DELETE /api/images/:imageID
GET /api/tags
POST /api/tags
PATCH /api/tags/:tagID
DELETE /api/tags/:tagID
POST /api/images/:imageID/tags
DELETE /api/images/:imageID/tags/:tagID
```

Swagger / OpenAPI
Backend มีเอกสาร API ในรูปแบบ OpenAPI และสามารถเปิดดูผ่าน Swagger UI ได้ ใช้สำหรับดูรายการ endpoint, request body, response schema และทดสอบ API ได้จาก browser โดยตรง

ถ้ารัน local:

Swagger UI: http://localhost:8080/swagger

OpenAPI YAML: http://localhost:8080/swagger/openapi.yaml

ถ้า deploy ขึ้น server:

Swagger UI: https://<backend-domain>/swagger

OpenAPI YAML: https://<backend-domain>/swagger/openapi.yaml

ถ้ามีการเพิ่ม service/endpoint ใหม่ ขั้นตอนการเพิ่ม service ใน swagger คือ:

เพิ่ม handler / route / service ตามปกติ

แก้ openapi.yaml (line 1) ให้มี path, request, response ใหม่

restart หรือ rebuild app

# Deployment

Deployment

ระบบถูกออกแบบให้สามารถ deploy ได้ทั้งในรูปแบบ local development และบน cloud environment โดยแยก frontend, backend และ database ออกจากกันอย่างชัดเจน

ในระดับ local โปรเจคใช้ Docker Compose สำหรับรัน service หลักทั้งหมดร่วมกัน ได้แก่

mysql สำหรับฐานข้อมูล

backend สำหรับ REST API

frontend สำหรับส่วนแสดงผลของระบบ

ไฟล์ที่ใช้สำหรับ deployment หลัก ได้แก่

docker-compose.yml

backend/Dockerfile

frontend/Dockerfile

frontend/nginx.conf

แนวทางการ Deploy

Frontend ถูก build เป็น static asset ด้วย Vite

จากนั้นนำไป serve ผ่าน Nginx

Backend ถูก build เป็น Go application และเปิดให้บริการผ่าน REST API

Database ใช้ MySQL

การเชื่อมต่อระหว่าง frontend และ backend ใช้ค่า VITE_API_BASE_URL

การจัดการ schema และ seed data ใช้ SQL migration

วิธีการ Deploy บน Production

การ deploy ระบบใน production จะแยก service ออกเป็น 3 ส่วนหลัก ได้แก่ frontend, backend และ database เพื่อให้สามารถจัดการแต่ละส่วนได้อย่างอิสระและลดผลกระทบเมื่อมีการอัปเดตระบบ

องค์ประกอบที่ใช้ในการ Deploy

Frontend

build ด้วย Vite

ได้ผลลัพธ์เป็น static files

serve ผ่าน Nginx

Backend

build เป็น Go binary

ทำงานเป็น REST API service

Database

ใช้ MySQL 8

เก็บข้อมูลของ image, tag และ relationship ต่าง ๆ

Containerization

ใช้ Docker สำหรับ build image ของ frontend และ backend

ใช้ Docker Compose หรือ platform service orchestration สำหรับรันระบบ

ลำดับการ Deploy
ลำดับการ deploy ที่เหมาะสมคือ

Deploy Database

สร้าง MySQL service หรือ database instance

กำหนด volume / persistent storage สำหรับข้อมูล

ตั้งค่า username, password, database name และ network access

ปิด public access ของ database และเปิดให้เชื่อมต่อได้เฉพาะ backend

Run Database Migration

รัน migration เพื่อสร้าง schema ของระบบ

ตัวอย่าง env ที่ใช้

```text
AUTO_MIGRATE=true
MOCKDATA=true
```

Deploy Backend

build backend image จาก source code

inject environment variables เช่น

database connection

app host / port

API prefix

CORS settings

start backend service

ตรวจสอบ health check เช่น /healthz

ตรวจสอบว่า Swagger เปิดได้ที่ /swagger

Deploy Frontend

build frontend โดยกำหนด VITE_API_BASE_URL ให้ชี้ไปยัง backend production

นำ build output ไป serve ผ่าน Nginx

เปิด public domain สำหรับ frontend

ทดสอบว่า frontend สามารถเรียก backend API ได้ถูกต้อง

ตัวอย่าง Flow ของการ Deploy

1. Push code ไปยัง Git repository

2. CI/CD pipeline หรือ deploy process ทำการ build image

3. Deploy database และตรวจสอบ connectivity

4. Run migration

5. Deploy backend และตรวจสอบ health check

6. Deploy frontend

7. Verify end-to-end flow ของระบบ

Environment Variables ที่สำคัญ

ฝั่ง backend

```text
DATABASE_URL หรือ DATABASE_DSN
APP_HOST
APP_PORT
```

AUTO_MIGRATE

MOCKDATA

ฝั่ง frontend

VITE_API_BASE_URL

การตรวจสอบหลัง Deploy
หลัง deploy ควรตรวจสอบอย่างน้อยดังนี้

frontend เปิดหน้าแรกได้

หน้า /manage ใช้งานได้

backend API ตอบกลับปกติ

Swagger ใช้งานได้

migration ถูก apply ครบ

frontend เรียกรูปและ tag จาก backend ได้จริง

# Server Specifications

Server Specifications สำหรับ Production

สำหรับ production environment ของระบบนี้ แนะนำให้เริ่มต้นด้วยสถาปัตยกรรมขนาดเล็กถึงกลาง โดยแยก frontend, backend และ database ออกจากกันอย่างชัดเจน เพื่อให้ดูแลง่ายและสามารถขยายระบบในอนาคตได้

สเปกเริ่มต้นที่เหมาะสมมีดังนี้

Frontend Server

1 vCPU

1-2 GB RAM

ใช้สำหรับ serve static asset ผ่าน Nginx

ภาระโหลดต่ำเมื่อเทียบกับ backend เพราะไม่มี business logic หนัก

Backend Server

2 vCPU

2-4 GB RAM

ใช้รัน Go Fiber API

รองรับการประมวลผล request จาก frontend, query database, filter รูปภาพ และจัดการ CRUD operation

Database Server

2 vCPU

4 GB RAM

40-80 GB SSD

ใช้ MySQL 8

ควรใช้ SSD เพื่อให้ query และ index access มีประสิทธิภาพดีขึ้น

ถ้า deploy บน platform เดียว เช่น Railway หรือ PaaS อื่น อาจ map เป็น 3 services คือ

frontend service

backend service

mysql service

OS / Software Version
สำหรับระบบปัจจุบันใน repo นี้ หาก deploy บน Railway จะเป็นลักษณะ platform-managed environment โดย service หลักยังคงอ้างอิง version จาก Docker image และ application runtime ที่ใช้อยู่จริงในโปรเจค

Host OS

กรณี deploy บน Railway จะไม่ได้ pin host OS เองภายใน repo เพราะเป็น platform-managed environment

Frontend Runtime

ใช้ `nginx:1.27-alpine`

Backend Build Runtime

ใช้ `golang:1.26-alpine`

Backend Runtime

ใช้ `alpine:3.21`

Database

local compose ใช้ `mysql:8.4`

cloud environment ใช้ Railway MySQL service

Application Versions

Frontend ใช้ `React 18.3.1`, `Vite 5.4.21`, `TypeScript 5.6.3`

Backend ใช้ `Go 1.26`, `Fiber v2.52.9`

Database model อ้างอิง `MySQL 8.4`

Container Runtime

ใช้ `Docker` และ `Dockerfile` สำหรับ build image

Local Orchestration

ใช้ `Docker Compose`

Cloud Deployment

ใช้ `Railway` โดยแยก service เป็น `frontend`, `backend`, `mysql`

CI/CD

ปัจจุบันใน repo ยังไม่มี `GitHub Actions` workflow และ deployment ใช้วิธี `Git push -> Railway auto deploy`

TLS / HTTPS

ใช้ HTTPS จาก public domain ที่ platform จัดการให้ หรือ reverse proxy / managed TLS ของ environment ที่ deploy

แนวทางการรองรับ Load
สำหรับระบบ image gallery ลักษณะนี้ ภาระโหลดหลักจะอยู่ที่

การโหลดหน้า gallery

การเรียก GET /api/images

การ filter ด้วย tag

การ scroll เพื่อโหลดรูปเพิ่ม

การอ่านข้อมูลมากกว่าการเขียนข้อมูล

ลักษณะโหลดของระบบจึงเป็น read-heavy workload มากกว่า write-heavy workload
ดังนั้นการออกแบบในระยะแรกควรเน้น

query ที่มี index เหมาะสม

จำกัดจำนวนข้อมูลต่อ request ด้วย limit

ใช้ cursor pagination แทน offset pagination

แยก frontend เป็น static service เพื่อให้ backend รับเฉพาะ API traffic

# Current Stack

Stack

Frontend:

- `React 18.3.1`

- `Vite 5.4.21`

- `TypeScript 5.6.3`

- `Tailwind CSS 3.4.15`

- `react-router-dom 6.28.x`

- `masonry-layout 4.2.2`

- `imagesloaded 5.0.0`

Backend:

- `Go 1.26`

- `Fiber v2.52.9`

- `go-sql-driver/mysql 1.9.3`

Database:

- `MySQL 8.4`

Frontend runtime image:

- `nginx:1.27-alpine`

Backend runtime image:

- `alpine:3.21`

Backend build image:

- `golang:1.26-alpine`

Containerization:

- `Docker`

- `Dockerfile`

Local orchestration:

- `Docker Compose`

Current cloud deployment target:

- `Railway`

Demo image source:

- `placehold.co`
