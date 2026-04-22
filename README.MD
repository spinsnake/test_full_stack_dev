# Full-Stack Developer Test

## Database Migration

ส่วนนี้ใช้ `golang-migrate` เพื่อจัดการ schema ของ MySQL โดยอ่านค่าการเชื่อมต่อจากไฟล์ `backend/.env` ผ่านสคริปต์ `backend/scripts/migrate.ps1`

### โครงสร้างไฟล์ที่เกี่ยวข้อง

- `backend/migration/*.up.sql` สำหรับ migrate ขึ้น
- `backend/migration/*.down.sql` สำหรับ rollback
- `backend/scripts/migrate.ps1` ตัวช่วยอ่าน env และเรียก `migrate`
- `backend/.env` ค่าจริงในเครื่อง
- `backend/.env.example` template

### ก่อนเริ่ม

1. ตรวจว่า MySQL รันอยู่และมีฐานข้อมูล `gallery_db`
2. ตั้งค่าใน `backend/.env` อย่างน้อย:
   - `MYSQL_HOST`
   - `MYSQL_PORT`
   - `MYSQL_USER`
   - `MYSQL_PASSWORD`
   - `MYSQL_DATABASE`
3. ติดตั้ง `migrate` ที่รองรับ MySQL driver

ตัวอย่างติดตั้ง:

```powershell
go install -tags "mysql" github.com/golang-migrate/migrate/v4/cmd/migrate@latest
```

### วิธีรัน Migration

รันจาก root โปรเจกต์ (`D:\Test_full_stack_developer`)

```powershell
# ดู version ปัจจุบัน
powershell -ExecutionPolicy Bypass -File backend/scripts/migrate.ps1 version

# migrate ขึ้นทั้งหมด
powershell -ExecutionPolicy Bypass -File backend/scripts/migrate.ps1 up

# rollback ลง 1 step
powershell -ExecutionPolicy Bypass -File backend/scripts/migrate.ps1 down 1

# force version (ใช้ตอนแก้ dirty state)
powershell -ExecutionPolicy Bypass -File backend/scripts/migrate.ps1 force 0
```

ถ้าต้องใช้ env ไฟล์อื่น:

```powershell
powershell -ExecutionPolicy Bypass -File backend/scripts/migrate.ps1 up -EnvFile .env.staging
```


