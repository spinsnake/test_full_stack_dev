# ภาพรวมของระบบ

## Masonry image gallery

ภาพรวมของระบบ จาก requirement ได้ทำการออกแบบระบบโดยแบ่งระบบ และ tech stach ออกเป็นสามส่วนดังนี้1. ส่วน frontent ใช้ Vite, react native, Typescript, Tailwind CSS, Masonry Layout ในการแสดงผล2. ส่วน backend ใช้ Go fiber ในการทำเป็น API service

3. ส่วน Database ใช้ MySql ในการเก็บข้อมูล

4. Containerization ใช้ Docker

5. Reverse Proxy, Nginx

ภาพแสดงส่วน backend frontend และ database

# database

Database Schemaประกอบด้วย table ต่อไปนี้1. images ทำหน้าที่เก็บข้อมูลรูปภาพ

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

วิธีรัน Migrationโปรเจคนี้มี script ช่วยรัน migration อยู่ที่ backend/scripts/migrate.ps1

ตัวอย่างคำสั่ง

```powershell
powershell -ExecutionPolicy Bypass -File backend/scripts/migrate.ps1 up
```

ใช้สำหรับ apply migration ทั้งหมดที่ยังไม่ได้รัน

```powershell
powershell -ExecutionPolicy Bypass -File backend/scripts/migrate.ps1 down 1
```

ใช้สำหรับ rollback ย้อนกลับ 1 step

```powershell
powershell -ExecutionPolicy Bypass -File backend/scripts/migrate.ps1 version
```

ใช้สำหรับตรวจสอบ version ปัจจุบันของฐานข้อมูล

ผลลัพธ์หลังรัน

เมื่อรัน up ระบบจะสร้าง schema หลักของ gallery และเพิ่ม demo data ตาม migration ที่กำหนด

เมื่อรัน down ระบบจะย้อนการเปลี่ยนแปลงตามลำดับ migration ที่ rollback

# frontend

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

Frontend Tech StackFrontend ของระบบพัฒนาด้วย React, Vite และ TypeScript สำหรับสร้าง Single Page Application (SPA) โดยใช้ Tailwind CSS สำหรับจัดการ style, react-router-dom สำหรับ routing ภายในระบบ, masonry-layout และ imagesloaded สำหรับจัดการ image gallery แบบ masonry และใช้ Playwright สำหรับ E2E testing ของฝั่ง frontend

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

วิธีรันรันจากโฟลเดอร์ frontend npm run e2e

| file | สิ่งที่ทดสอบ |
| --- | --- |
| frontend/e2e/gallery.spec.ts (line 1) | ทดสอบหน้า gallery ว่าสามารถโหลดรายการรูปเริ่มต้นได้, เปิด Tags Filter, เลือก tag แล้วจำนวนรูปถูกกรองถูกต้อง, และสามารถกดไปหน้า /manage ได้ |
| frontend/e2e/manage.spec.ts (line 1) | ทดสอบหน้า /manage ว่าสามารถกรอกข้อมูลสร้าง image ใหม่ได้, เลือก tag ตอนสร้าง image ได้, กด Create Image แล้ว image ใหม่ถูกเพิ่มใน catalog พร้อม tag ที่เลือกไว้ |

# backend

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

แนวทางการทำงานของ Backendลำดับการทำงานหลักของ backend คือ

client ส่ง request เข้ามาที่ route

route ส่งต่อไปยัง handler

handler parse request และเรียก service

service จัดการ business logic และเรียก repository

repository query หรือ update ข้อมูลใน MySQL

handler ส่ง response กลับในรูปแบบ JSON

ตัวอย่าง API หลักที่ระบบรองรับ

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

Swagger / OpenAPIBackend มีเอกสาร API ในรูปแบบ OpenAPI และสามารถเปิดดูผ่าน Swagger UI ได้ ใช้สำหรับดูรายการ endpoint, request body, response schema และทดสอบ API ได้จาก browser โดยตรง

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

# deployment

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

Deployment บน Cloudในโปรเจคนี้สามารถ deploy ขึ้น Railway ได้โดยแยก service ออกเป็น

frontend

backend

mysql

ลำดับการ deploy คือ

deploy database (mysql)

deploy backend

deploy frontend

สำหรับ backend สามารถตั้งค่าให้รัน migration อัตโนมัติผ่าน environment variable เช่น

AUTO_MIGRATE=true

MOCKDATA=true

ส่วน frontend สามารถกำหนด backend URL ผ่าน

VITE_API_BASE_URL=https://<backend-domain>/api

การเข้าถึงระบบหลัง Deploy

Frontend ใช้สำหรับเข้าใช้งานระบบผ่าน browser

Backend สามารถเรียกใช้งาน API ได้โดยตรง

Swagger UI สามารถเข้าผ่าน /swagger

OpenAPI spec สามารถเข้าผ่าน /swagger/openapi.yaml
