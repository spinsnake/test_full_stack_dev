# Frontend

React + Vite + TypeScript frontend for the image gallery SPA.

## Stack

- React 18
- Vite 5
- TypeScript 5
- Tailwind CSS 3
- masonry-layout
- imagesloaded
- react-router-dom

## Routes

- `/` gallery page
- `/manage` management page

## Features

### Gallery

- masonry image feed
- infinite scroll
- chunk size controls: 10, 30, 50
- tag filtering
- tag modal with all available tags
- mobile burger menu for header controls
- responsive layout for desktop and mobile

### Management

- create, update, soft delete images
- create, update, soft delete tags
- attach tags to images
- detach tags from images
- add-tag modal with multi-select
- tag options already attached to an image are hidden in the modal

## Environment

Create `frontend/.env` from `frontend/.env.example`:

```env
VITE_API_BASE_URL=/api
```

Local development uses `/api` because `vite.config.ts` proxies requests to `http://localhost:8080`.

For production:

```env
VITE_API_BASE_URL=https://api.example.com/api
```

## Prerequisites

- Node.js 20+
- npm
- backend API running on `http://localhost:8080`

## Install

```powershell
cd frontend
npm install
```

## Run

```powershell
npm run dev
```

Default URL:

```text
http://localhost:5173
```

## Build

```powershell
npm run build
```

Preview:

```powershell
npm run preview
```

## Backend Endpoints Used

- `GET /api/tags`
- `GET /api/images?limit=<chunkSize>`
- `GET /api/images?limit=<chunkSize>&cursor=<lastImageId>`
- `GET /api/images?limit=<chunkSize>&tag=<slug>`
- `POST /api/images`
- `PATCH /api/images/:imageID`
- `DELETE /api/images/:imageID`
- `POST /api/tags`
- `PATCH /api/tags/:tagID`
- `DELETE /api/tags/:tagID`
- `POST /api/images/:imageID/tags`
- `DELETE /api/images/:imageID/tags/:tagID`

## UI Notes

- gallery cards show image source, dimensions, and tags
- clicking a tag filters the gallery
- the fixed desktop header collapses into a burger menu on mobile
- the management page includes separate image and tag sections

## Project Structure

```text
frontend/
  src/
    App.tsx
    api.ts
    index.css
    main.tsx
    types.ts
    pages/
      GalleryPage.tsx
      ManagePage.tsx
  index.html
  vite.config.ts
  tailwind.config.cjs
  postcss.config.cjs
  package.json
```

## Notes

- The UI expects seeded backend data for the best demo experience.
- The Vite app was scaffolded manually in this repository and builds successfully with `npm run build`.
