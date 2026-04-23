import type { Page, Route } from '@playwright/test';
import type {
  AttachTagPayload,
  CreateImagePayload,
  CreateTagPayload,
  ImageItem,
  ImageTag,
  Tag,
  UpdateImagePayload,
  UpdateTagPayload,
} from '../../src/types';

type MockImageRecord = Omit<ImageItem, 'tags'> & {
  tagIDs: number[];
};

type MockState = {
  images: MockImageRecord[];
  tags: Tag[];
  nextImageID: number;
  nextTagID: number;
};

const timestamp = '2026-04-23T00:00:00Z';

function buildDataImage(label: string, width: number, height: number, background: string) {
  const svg = `
    <svg xmlns="http://www.w3.org/2000/svg" width="${width}" height="${height}" viewBox="0 0 ${width} ${height}">
      <rect width="100%" height="100%" fill="${background}" />
      <text x="50%" y="50%" dominant-baseline="middle" text-anchor="middle" fill="#ffffff"
        font-family="Arial, sans-serif" font-size="${Math.max(24, Math.floor(width / 10))}" font-weight="700">
        ${label}
      </text>
    </svg>
  `.trim();

  return `data:image/svg+xml;charset=utf-8,${encodeURIComponent(svg)}`;
}

function slugify(value: string) {
  return value
    .trim()
    .toLowerCase()
    .replace(/[^a-z0-9]+/g, '-')
    .replace(/^-+|-+$/g, '') || 'tag';
}

function createInitialState(): MockState {
  const tags: Tag[] = [
    {
      id: 1,
      name: 'Nature',
      slug: 'nature',
      deleted_at: null,
      created_at: timestamp,
      updated_at: timestamp,
    },
    {
      id: 2,
      name: 'City',
      slug: 'city',
      deleted_at: null,
      created_at: timestamp,
      updated_at: timestamp,
    },
    {
      id: 3,
      name: 'Sports',
      slug: 'sports',
      deleted_at: null,
      created_at: timestamp,
      updated_at: timestamp,
    },
  ];

  const images: MockImageRecord[] = [
    {
      id: 103,
      image_url: buildDataImage('Forest Path', 1200, 1500, '#3b7d58'),
      thumbnail_url: buildDataImage('Forest Path', 400, 500, '#3b7d58'),
      width: 1200,
      height: 1500,
      alt_text: 'Forest Path',
      source: 'mock-seed',
      deleted_at: null,
      created_at: timestamp,
      updated_at: timestamp,
      tagIDs: [1],
    },
    {
      id: 102,
      image_url: buildDataImage('City Glow', 1200, 900, '#1f5fbf'),
      thumbnail_url: buildDataImage('City Glow', 400, 300, '#1f5fbf'),
      width: 1200,
      height: 900,
      alt_text: 'City Glow',
      source: 'mock-seed',
      deleted_at: null,
      created_at: timestamp,
      updated_at: timestamp,
      tagIDs: [2],
    },
    {
      id: 101,
      image_url: buildDataImage('Courtside Moment', 1080, 1350, '#c96b2c'),
      thumbnail_url: buildDataImage('Courtside Moment', 360, 450, '#c96b2c'),
      width: 1080,
      height: 1350,
      alt_text: 'Courtside Moment',
      source: 'mock-seed',
      deleted_at: null,
      created_at: timestamp,
      updated_at: timestamp,
      tagIDs: [3],
    },
  ];

  return {
    images,
    tags,
    nextImageID: 104,
    nextTagID: 4,
  };
}

function listActiveTags(state: MockState) {
  return [...state.tags].filter((tag) => tag.deleted_at == null).sort((left, right) => right.id - left.id);
}

function serializeImage(record: MockImageRecord, state: MockState): ImageItem {
  const imageTags: ImageTag[] = record.tagIDs
    .map((tagID) => state.tags.find((tag) => tag.id === tagID && tag.deleted_at == null))
    .filter((tag): tag is Tag => Boolean(tag))
    .map((tag) => ({
      id: tag.id,
      name: tag.name,
      slug: tag.slug,
    }));

  return {
    ...record,
    tags: imageTags,
  };
}

function listImages(state: MockState, options?: { cursor?: number | null; limit?: number; tag?: string | null }) {
  const { cursor = null, limit = 50, tag = null } = options ?? {};

  const activeTags = new Set(listActiveTags(state).map((item) => item.id));
  const filtered = [...state.images]
    .filter((image) => image.deleted_at == null)
    .filter((image) => image.tagIDs.every((tagID) => activeTags.has(tagID)))
    .sort((left, right) => right.id - left.id)
    .filter((image) => (cursor ? image.id < cursor : true))
    .filter((image) => {
      if (!tag) {
        return true;
      }

      return image.tagIDs.some((tagID) => {
        const tagRecord = state.tags.find((entry) => entry.id === tagID && entry.deleted_at == null);
        return tagRecord?.slug === tag;
      });
    })
    .slice(0, limit);

  return filtered.map((image) => serializeImage(image, state));
}

async function readJSON<T>(route: Route) {
  return (await route.request().postDataJSON()) as T;
}

async function fulfillJSON(route: Route, body: unknown, status = 200) {
  await route.fulfill({
    status,
    contentType: 'application/json',
    body: JSON.stringify(body),
  });
}

async function handleRoute(route: Route, state: MockState) {
  const request = route.request();
  const url = new URL(request.url());
  const path = url.pathname;
  const method = request.method().toUpperCase();

  if (method === 'GET' && path === '/api/tags') {
    return fulfillJSON(route, { data: listActiveTags(state) });
  }

  if (method === 'GET' && path === '/api/images') {
    const cursorParam = url.searchParams.get('cursor');
    const limitParam = url.searchParams.get('limit');
    const tagParam = url.searchParams.get('tag');

    return fulfillJSON(route, {
      data: listImages(state, {
        cursor: cursorParam ? Number(cursorParam) : null,
        limit: limitParam ? Number(limitParam) : 50,
        tag: tagParam,
      }),
    });
  }

  if (method === 'POST' && path === '/api/images') {
    const payload = await readJSON<CreateImagePayload>(route);
    const created: MockImageRecord = {
      id: state.nextImageID++,
      image_url: payload.image_url,
      thumbnail_url: payload.thumbnail_url ?? null,
      width: payload.width ?? null,
      height: payload.height ?? null,
      alt_text: payload.alt_text ?? null,
      source: payload.source ?? null,
      deleted_at: null,
      created_at: timestamp,
      updated_at: timestamp,
      tagIDs: [],
    };

    state.images.push(created);
    return fulfillJSON(route, { data: serializeImage(created, state) }, 201);
  }

  if (method === 'POST' && path === '/api/tags') {
    const payload = await readJSON<CreateTagPayload>(route);
    const slug = slugify(payload.slug || payload.name);
    const created: Tag = {
      id: state.nextTagID++,
      name: payload.name,
      slug,
      deleted_at: null,
      created_at: timestamp,
      updated_at: timestamp,
    };

    state.tags.push(created);
    return fulfillJSON(route, { data: created }, 201);
  }

  const imageMatch = path.match(/^\/api\/images\/(\d+)$/);
  if (imageMatch) {
    const imageID = Number(imageMatch[1]);
    const image = state.images.find((entry) => entry.id === imageID && entry.deleted_at == null);

    if (!image) {
      return fulfillJSON(route, { error: 'Image not found.' }, 404);
    }

    if (method === 'PATCH') {
      const payload = await readJSON<UpdateImagePayload>(route);

      if (typeof payload.image_url === 'string') {
        image.image_url = payload.image_url;
      }
      if (typeof payload.thumbnail_url === 'string') {
        image.thumbnail_url = payload.thumbnail_url;
      }
      if (typeof payload.width === 'number') {
        image.width = payload.width;
      }
      if (typeof payload.height === 'number') {
        image.height = payload.height;
      }
      if (typeof payload.alt_text === 'string') {
        image.alt_text = payload.alt_text;
      }
      if (typeof payload.source === 'string') {
        image.source = payload.source;
      }

      image.updated_at = timestamp;
      return fulfillJSON(route, { data: serializeImage(image, state) });
    }

    if (method === 'DELETE') {
      image.deleted_at = timestamp;
      image.updated_at = timestamp;
      return fulfillJSON(route, { message: `Image #${imageID} deleted.` });
    }
  }

  const tagMatch = path.match(/^\/api\/tags\/(\d+)$/);
  if (tagMatch) {
    const tagID = Number(tagMatch[1]);
    const tag = state.tags.find((entry) => entry.id === tagID && entry.deleted_at == null);

    if (!tag) {
      return fulfillJSON(route, { error: 'Tag not found.' }, 404);
    }

    if (method === 'PATCH') {
      const payload = await readJSON<UpdateTagPayload>(route);

      if (typeof payload.name === 'string') {
        tag.name = payload.name;
      }
      if (typeof payload.slug === 'string') {
        tag.slug = slugify(payload.slug);
      } else if (typeof payload.name === 'string' && !payload.slug) {
        tag.slug = slugify(payload.name);
      }

      tag.updated_at = timestamp;
      return fulfillJSON(route, { data: tag });
    }

    if (method === 'DELETE') {
      tag.deleted_at = timestamp;
      tag.updated_at = timestamp;
      state.images.forEach((image) => {
        image.tagIDs = image.tagIDs.filter((imageTagID) => imageTagID !== tagID);
      });

      return fulfillJSON(route, { message: `Tag #${tag.slug} deleted.` });
    }
  }

  const imageTagMatch = path.match(/^\/api\/images\/(\d+)\/tags(?:\/(\d+))?$/);
  if (imageTagMatch) {
    const imageID = Number(imageTagMatch[1]);
    const tagIDFromPath = imageTagMatch[2] ? Number(imageTagMatch[2]) : null;
    const image = state.images.find((entry) => entry.id === imageID && entry.deleted_at == null);

    if (!image) {
      return fulfillJSON(route, { error: 'Image not found.' }, 404);
    }

    if (method === 'POST') {
      const payload = await readJSON<AttachTagPayload>(route);
      const tag = state.tags.find((entry) => entry.id === payload.tag_id && entry.deleted_at == null);

      if (!tag) {
        return fulfillJSON(route, { error: 'Tag not found.' }, 404);
      }

      if (!image.tagIDs.includes(tag.id)) {
        image.tagIDs.push(tag.id);
      }

      return fulfillJSON(route, { data: { image_id: imageID, tag_id: tag.id } }, 201);
    }

    if (method === 'DELETE' && tagIDFromPath) {
      image.tagIDs = image.tagIDs.filter((entry) => entry !== tagIDFromPath);
      return fulfillJSON(route, { message: 'Tag detached.' });
    }
  }

  return fulfillJSON(route, { error: `Unhandled mock route: ${method} ${path}` }, 500);
}

export async function installMockApi(page: Page) {
  const state = createInitialState();
  await page.route('**/api/**', async (route) => handleRoute(route, state));
  return state;
}
