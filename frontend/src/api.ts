import type {
  ApiMessageResponse,
  ApiResponse,
  AttachTagPayload,
  CreateImagePayload,
  CreateTagPayload,
  ImageItem,
  Tag,
  UpdateImagePayload,
  UpdateTagPayload,
} from './types';

const API_BASE_URL = (import.meta.env.VITE_API_BASE_URL || '/api').replace(/\/$/, '');
const ADMIN_IMAGE_PAGE_SIZE = 50;

type FetchImagesOptions = {
  cursor?: number | null;
  limit?: number;
  tag?: string;
};

function normalizeImage(image: ImageItem): ImageItem {
  return {
    ...image,
    tags: Array.isArray(image.tags) ? image.tags : [],
  };
}

async function request<T>(path: string, init?: RequestInit): Promise<T> {
  const response = await fetch(path, init);
  if (!response.ok) {
    const body = (await response.json().catch(() => null)) as { error?: string } | null;
    throw new Error(body?.error || `Request failed with status ${response.status}`);
  }

  if (response.status === 204) {
    return undefined as T;
  }

  return response.json() as Promise<T>;
}

function buildImagesURL({ cursor, limit, tag }: FetchImagesOptions) {
  const url = new URL(`${API_BASE_URL}/images`, window.location.origin);
  url.searchParams.set('limit', String(limit ?? ADMIN_IMAGE_PAGE_SIZE));

  if (cursor) {
    url.searchParams.set('cursor', String(cursor));
  }

  if (tag) {
    url.searchParams.set('tag', tag);
  }

  return `${url.pathname}${url.search}`;
}

export async function fetchTags(): Promise<Tag[]> {
  const response = await request<ApiResponse<Tag[]>>(`${API_BASE_URL}/tags`);
  return response.data;
}

export async function fetchImages(options: FetchImagesOptions = {}): Promise<ImageItem[]> {
  const response = await request<ApiResponse<ImageItem[]>>(buildImagesURL(options));
  return response.data.map(normalizeImage);
}

export async function fetchAllImages(): Promise<ImageItem[]> {
  const pages: ImageItem[] = [];
  let cursor: number | null = null;

  for (;;) {
    const chunk = await fetchImages({ cursor, limit: ADMIN_IMAGE_PAGE_SIZE });
    pages.push(...chunk);

    if (chunk.length < ADMIN_IMAGE_PAGE_SIZE) {
      break;
    }

    cursor = chunk[chunk.length - 1]?.id ?? null;
    if (!cursor) {
      break;
    }
  }

  return pages;
}

export async function createImage(payload: CreateImagePayload): Promise<ImageItem> {
  const response = await request<ApiResponse<ImageItem>>(`${API_BASE_URL}/images`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(payload),
  });
  return normalizeImage(response.data);
}

export async function updateImage(id: number, payload: UpdateImagePayload): Promise<ImageItem> {
  const response = await request<ApiResponse<ImageItem>>(`${API_BASE_URL}/images/${id}`, {
    method: 'PATCH',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(payload),
  });
  return normalizeImage(response.data);
}

export async function deleteImage(id: number): Promise<ApiMessageResponse> {
  return request<ApiMessageResponse>(`${API_BASE_URL}/images/${id}`, {
    method: 'DELETE',
  });
}

export async function createTag(payload: CreateTagPayload): Promise<Tag> {
  const response = await request<ApiResponse<Tag>>(`${API_BASE_URL}/tags`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(payload),
  });
  return response.data;
}

export async function updateTag(id: number, payload: UpdateTagPayload): Promise<Tag> {
  const response = await request<ApiResponse<Tag>>(`${API_BASE_URL}/tags/${id}`, {
    method: 'PATCH',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(payload),
  });
  return response.data;
}

export async function deleteTag(id: number): Promise<ApiMessageResponse> {
  return request<ApiMessageResponse>(`${API_BASE_URL}/tags/${id}`, {
    method: 'DELETE',
  });
}

export async function attachTagToImage(imageID: number, payload: AttachTagPayload) {
  const response = await request<ApiResponse<{ image_id: number; tag_id: number }>>(
    `${API_BASE_URL}/images/${imageID}/tags`,
    {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(payload),
    },
  );

  return response.data;
}

export async function detachTagFromImage(imageID: number, tagID: number): Promise<ApiMessageResponse> {
  return request<ApiMessageResponse>(`${API_BASE_URL}/images/${imageID}/tags/${tagID}`, {
    method: 'DELETE',
  });
}
