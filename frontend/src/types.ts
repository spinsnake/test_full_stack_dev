export type ApiResponse<T> = {
  data: T;
};

export type ApiMessageResponse = {
  message: string;
};

export type Tag = {
  id: number;
  name: string;
  slug: string;
  is_active: boolean;
  created_at: string;
  updated_at: string;
};

export type ImageTag = {
  id: number;
  name: string;
  slug: string;
};

export type ImageItem = {
  id: number;
  image_url: string;
  thumbnail_url?: string | null;
  width?: number | null;
  height?: number | null;
  alt_text?: string | null;
  source?: string | null;
  is_active: boolean;
  created_at: string;
  updated_at: string;
  tags: ImageTag[];
};

export type CreateImagePayload = {
  image_url: string;
  thumbnail_url?: string;
  width?: number;
  height?: number;
  alt_text?: string;
  source?: string;
};

export type UpdateImagePayload = Partial<CreateImagePayload>;

export type CreateTagPayload = {
  name: string;
  slug?: string;
};

export type UpdateTagPayload = Partial<CreateTagPayload>;

export type AttachTagPayload = {
  tag_id: number;
};
