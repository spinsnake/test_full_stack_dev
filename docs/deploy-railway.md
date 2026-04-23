# Deploy To Railway

This project deploys cleanly to Railway as 3 services in one project:

- `mysql`
- `backend`
- `frontend`

The recommended order is:

1. Deploy `mysql`
2. Deploy `backend`
3. Deploy `frontend`

## 1. Create The MySQL Service

In Railway:

1. Open your project canvas.
2. Click `New`.
3. Choose `Database` -> `MySQL`.
4. Wait until the service is healthy.

Keep the service name as `mysql` or any readable name you prefer.

## 2. Deploy The Backend Service

Create a new service from the same GitHub repo and point it to the backend folder.

Recommended backend settings:

- Source repo: this repository
- Root directory: `backend`
- Builder: `Dockerfile`
- Watch path: `backend/**`

### Backend Variables

Set these in the backend service:

- `APP_NAME=image-gallery-api`
- `APP_HOST=0.0.0.0`
- `API_PREFIX=/api`
- `CORS_ALLOW_ORIGINS=*`
- `DEFAULT_PAGE_LIMIT=12`
- `MAX_PAGE_LIMIT=60`
- `MOCKDATA=true`
- `AUTO_MIGRATE=true`
- `DB_MAX_OPEN_CONNS=25`
- `DB_MAX_IDLE_CONNS=25`
- `DB_CONN_MAX_LIFETIME_MIN=30`
- `DB_CONN_MAX_IDLE_TIME_MIN=10`
- `DATABASE_URL=${{mysql.MYSQL_URL}}`

Notes:

- The app now accepts `DATABASE_URL` and Railway-style MySQL variables.
- The app also accepts Railway's `PORT` automatically, so you do not need to hardcode a deploy port.
- The MySQL connection now enables `multiStatements=true`, which is required for the schema migration file in this repo.
- `MOCKDATA=true` runs `000002_seed_demo_data.up.sql` once after schema migrations.
- `MOCKDATA=false` applies schema migrations only.
- `AUTO_MIGRATE=true` runs migrations automatically when the API container starts.

### Backend Deploy Settings

Under `Deploy`:

- Start command: leave default from Docker image
- Pre-deploy command: optional

Recommended Railway setup:

- `AUTO_MIGRATE=true`
- leave `Pre-deploy command` empty unless you specifically want a separate migration step

With `AUTO_MIGRATE=true`, the backend applies schema migrations on startup and optionally seeds demo data when `MOCKDATA=true`.

### Backend Networking

Generate a public domain for the backend service.

You will get a URL like:

```text
https://backend-production-xxxx.up.railway.app
```

Use that URL in the frontend step below.

## 3. Deploy The Frontend Service

Create another service from the same GitHub repo and point it to the frontend folder.

Recommended frontend settings:

- Source repo: this repository
- Root directory: `frontend`
- Builder: `Dockerfile`
- Watch path: `frontend/**`

### Frontend Variables

Set this build variable in the frontend service:

- `VITE_API_BASE_URL=https://YOUR_BACKEND_DOMAIN/api`

Example:

```text
VITE_API_BASE_URL=https://backend-production-xxxx.up.railway.app/api
```

Because the frontend is compiled by Vite, set this before deploying or redeploy after changing it.

### Frontend Networking

Generate a public domain for the frontend service.

This will be the main app URL you share with users.

## 4. Redeploy Order

If you update both services:

1. Deploy backend first
2. Wait for backend to become healthy
3. Deploy frontend

That avoids frontend builds pointing at an old backend URL or schema.

## 5. Quick Smoke Check

After deploy:

- Open frontend `/`
- Open frontend `/manage`
- Create a new image
- Create a new tag
- Attach a tag to an image
- Open backend `/healthz`
- Open backend `/swagger`

## 6. Railway Service Summary

For this repo, the practical production-like setup is:

- `mysql`: Railway MySQL service
- `backend`: Go API from `backend/`
- `frontend`: Vite app served by Nginx from `frontend/`

The frontend talks to the backend through `VITE_API_BASE_URL`, not through Railway internal proxying.
