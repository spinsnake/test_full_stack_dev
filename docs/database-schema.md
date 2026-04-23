# Database Schema

This document explains the current MySQL schema used by the gallery application.

The schema is defined in:

- `backend/migration/000001_create_gallery_tables.up.sql`
- `backend/migration/000002_seed_demo_data.up.sql`

## Design Summary

- `images` stores the core image records shown in the gallery.
- `tags` stores reusable labels used for filtering and grouping images.
- `image_tags` stores the many-to-many relationship between images and tags.
- Soft delete is implemented with `deleted_at`.
  - `deleted_at IS NULL` means the record is active.
  - `deleted_at IS NOT NULL` means the record was soft deleted and should not be returned by normal API queries.

## Table: `images`

Stores each gallery image and its display metadata.

| Field | Type | Purpose |
| --- | --- | --- |
| `id` | `BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY` | Unique identifier for each image. Used by the API, pagination cursor, and joins. |
| `image_url` | `VARCHAR(512) NOT NULL` | Original or main image URL used for the full image asset. |
| `thumbnail_url` | `VARCHAR(512) NULL` | Optional smaller image URL used for faster gallery rendering. |
| `width` | `INT UNSIGNED NULL` | Optional image width in pixels. Used to preserve aspect ratio in the UI. |
| `height` | `INT UNSIGNED NULL` | Optional image height in pixels. Used together with `width` for layout calculations. |
| `alt_text` | `VARCHAR(255) NULL` | Optional human-readable caption or accessibility text shown on the card and used for `alt`. |
| `source` | `VARCHAR(100) NULL` | Optional source label such as `seed:placehold.co` or an external provider/domain. |
| `deleted_at` | `DATETIME(3) NULL` | Soft delete timestamp. `NULL` means the image is visible to the application. |
| `created_at` | `DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3)` | Record creation timestamp. |
| `updated_at` | `DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3)` | Last update timestamp, automatically refreshed by MySQL on update. |

### Image Indexes

| Index | Columns | Purpose |
| --- | --- | --- |
| `idx_images_feed` | `(deleted_at, id)` | Helps gallery queries fetch non-deleted rows ordered by newest first. |

## Table: `tags`

Stores the filterable tag catalog.

| Field | Type | Purpose |
| --- | --- | --- |
| `id` | `BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY` | Unique identifier for each tag. |
| `name` | `VARCHAR(64) NOT NULL` | Human-readable display name such as `Nature` or `Travel`. |
| `slug` | `VARCHAR(64) NOT NULL` | URL/query-safe tag key such as `nature` or `travel`. Used by API filtering. |
| `deleted_at` | `DATETIME(3) NULL` | Soft delete timestamp. `NULL` means the tag is still active in the catalog. |
| `active_name` | `VARCHAR(64) GENERATED ALWAYS AS (...) STORED` | Generated column that keeps `name` only when the row is not deleted. Used to enforce uniqueness among active tags while still allowing re-use after soft delete. |
| `active_slug` | `VARCHAR(64) GENERATED ALWAYS AS (...) STORED` | Generated column that keeps `slug` only when the row is not deleted. Same purpose as `active_name`, but for tag slugs. |
| `created_at` | `DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3)` | Record creation timestamp. |
| `updated_at` | `DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3)` | Last update timestamp, automatically refreshed by MySQL on update. |

### Tag Constraints and Indexes

| Name | Type | Purpose |
| --- | --- | --- |
| `uk_tags_active_name` | `UNIQUE KEY (active_name)` | Prevents duplicate active tag names. Deleted tags do not block reuse. |
| `uk_tags_active_slug` | `UNIQUE KEY (active_slug)` | Prevents duplicate active tag slugs. Deleted tags do not block reuse. |
| `idx_tags_listing` | `INDEX (deleted_at, name)` | Helps list non-deleted tags in name order. |

## Table: `image_tags`

Stores tag assignments for each image.

| Field | Type | Purpose |
| --- | --- | --- |
| `image_id` | `BIGINT UNSIGNED NOT NULL` | Foreign key to `images.id`. |
| `tag_id` | `BIGINT UNSIGNED NOT NULL` | Foreign key to `tags.id`. |
| `created_at` | `DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3)` | Timestamp of when the tag was attached to the image. |

### Image Tag Constraints and Indexes

| Name | Type | Purpose |
| --- | --- | --- |
| `PRIMARY KEY (image_id, tag_id)` | Primary key | Prevents duplicate tag assignment for the same image. |
| `idx_image_tags_tag_image` | `INDEX (tag_id, image_id)` | Helps reverse lookups from tag to image. |
| `fk_image_tags_image` | Foreign key | Links `image_id` to `images.id`; `ON DELETE CASCADE` removes assignments when an image is deleted physically. |
| `fk_image_tags_tag` | Foreign key | Links `tag_id` to `tags.id`; `ON DELETE CASCADE` removes assignments when a tag is deleted physically. |

## Relationship Model

- One image can have many tags.
- One tag can belong to many images.
- The `image_tags` table is the join table for that many-to-many relationship.

## Query Behavior in the App

The backend currently treats a row as active when:

- `images.deleted_at IS NULL`
- `tags.deleted_at IS NULL`

Typical API behavior:

- Gallery feed reads from `images` ordered by `id DESC`.
- Tag filtering joins `images -> image_tags -> tags`.
- Tag assignment checks that both the image and tag are not soft deleted.
- Soft delete only sets `deleted_at`; it does not physically remove the row.

## Seed Data

`000002_seed_demo_data.up.sql` inserts:

- readable tags such as `Nature`, `Travel`, `Food`, and `Technology`
- 100 demo image rows from `placehold.co`
- link rows in `image_tags` so the gallery can be filtered immediately after setup

## Notes for Future Changes

- If image visibility needs more than soft delete, add a dedicated status field with a different meaning from `deleted_at`.
- If tags need hierarchy or grouping, add a separate relation instead of overloading `slug` or `name`.
- If the gallery grows large, consider additional composite indexes for heavier admin or search queries.
