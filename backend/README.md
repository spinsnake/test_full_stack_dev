# Backend

Go + Fiber backend for the image gallery test. It provides CRUD APIs for images and tags, plus tag-to-image assignments, backed by MySQL.

## Stack

- Go `1.26`
- Fiber `v2`
- MySQL `8`
- SQL migrations with `golang-migrate`
- OpenAPI spec + Swagger UI served by the app itself

## Structure

```text
backend/
  cmd/api/                 # application entrypoint
  docs/                    # embedded OpenAPI spec + Swagger UI
  internal/
    adapter/handler/       # HTTP handlers
    adapter/repo/          # database repositories
    entities/              # request/response/domain models
    infra/                 # infrastructure adapters (MySQL)
    port/                  # service/repository interfaces
    routes/                # route binding
    service/               # business logic
  migration/               # SQL migrations
  scripts/migrate.ps1      # migration helper script
```

## Environment

Copy values from `.env.example` into `.env` and set at least:

- `APP_HOST`
- `APP_PORT`
- `MYSQL_HOST`
- `MYSQL_PORT`
- `MYSQL_USER`
- `MYSQL_PASSWORD`
- `MYSQL_DATABASE`

## Migration

Run from the repository root:

```powershell
# apply migrations
powershell -ExecutionPolicy Bypass -File backend/scripts/migrate.ps1 up

# rollback one step
powershell -ExecutionPolicy Bypass -File backend/scripts/migrate.ps1 down 1

# show current version
powershell -ExecutionPolicy Bypass -File backend/scripts/migrate.ps1 version
```

After `up`, the database contains demo seed data from `000002_seed_demo_data`:

- `100` placeholder images from `placehold.co`
- `20` readable tags such as `nature`, `travel`, `food`, `city`
- `1-5` tags attached to each seeded image

Seed tags are inserted with `INSERT IGNORE`, so existing tags like `travel` or `food` are reused instead of causing unique-key conflicts. `down 1` removes seeded images and image-tag assignments, but it does not remove the readable tag catalog because those tags may already belong to other app data.

## Run

Run from `backend/`:

```powershell
go run ./cmd/api
```

By default the API listens on `http://localhost:8080`.

## Docker

The backend container build is defined in [Dockerfile](D:\Test_full_stack_developer\backend\Dockerfile).

From the repository root, you can run the full stack with:

```powershell
docker compose up --build
```

In Docker Compose:

- backend listens internally on `backend:8080`
- MySQL runs as service `mysql`
- migrations are applied by the `migrate` service before the backend starts

## Test

Run all backend tests from `backend/`:

```powershell
go test ./...
```

The current test suite covers:

- image insert, update, soft delete, and list by limit
- tag insert, update, and soft delete
- attach tag to image and remove tag from image

## Swagger

Once the server is running:

- Swagger UI: `http://localhost:8080/swagger`
- OpenAPI YAML: `http://localhost:8080/swagger/openapi.yaml`

## API Summary

Images:

- `POST /api/images`
- `GET /api/images`
- `GET /api/images/:imageID`
- `PATCH /api/images/:imageID`
- `DELETE /api/images/:imageID`

Tags:

- `POST /api/tags`
- `GET /api/tags`
- `GET /api/tags/:tagID`
- `PATCH /api/tags/:tagID`
- `DELETE /api/tags/:tagID`

Image tag assignments:

- `POST /api/images/:imageID/tags`
- `DELETE /api/images/:imageID/tags/:tagID`

## Example Requests

Create a tag:

```powershell
curl -X POST http://localhost:8080/api/tags ^
  -H "Content-Type: application/json" ^
  -d "{\"name\":\"Travel\",\"slug\":\"travel\"}"
```

Create an image:

```powershell
curl -X POST http://localhost:8080/api/images ^
  -H "Content-Type: application/json" ^
  -d "{\"image_url\":\"https://placehold.co/1200x900?text=Travel+01\",\"thumbnail_url\":\"https://placehold.co/400x300?text=Travel+01\",\"width\":1200,\"height\":900,\"source\":\"placehold.co\"}"
```

Attach a tag to an image:

```powershell
curl -X POST http://localhost:8080/api/images/1/tags ^
  -H "Content-Type: application/json" ^
  -d "{\"tag_id\":1}"
```
