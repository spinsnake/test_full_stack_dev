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

## Run

Run from `backend/`:

```powershell
go run ./cmd/api
```

By default the API listens on `http://localhost:8080`.

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
