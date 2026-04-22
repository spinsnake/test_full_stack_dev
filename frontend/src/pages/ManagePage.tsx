import { useCallback, useEffect, useMemo, useState } from 'react';
import { Link } from 'react-router-dom';
import {
  attachTagToImage,
  createImage,
  createTag,
  deleteImage,
  deleteTag,
  detachTagFromImage,
  fetchAllImages,
  fetchTags,
  updateImage,
  updateTag,
} from '../api';
import type {
  CreateImagePayload,
  CreateTagPayload,
  ImageItem,
  Tag,
  UpdateImagePayload,
  UpdateTagPayload,
} from '../types';

type ImageFormState = {
  id: number | null;
  imageURL: string;
  thumbnailURL: string;
  width: string;
  height: string;
  altText: string;
  source: string;
};

type TagFormState = {
  id: number | null;
  name: string;
  slug: string;
};

const emptyImageForm: ImageFormState = {
  id: null,
  imageURL: '',
  thumbnailURL: '',
  width: '',
  height: '',
  altText: '',
  source: '',
};

const emptyTagForm: TagFormState = {
  id: null,
  name: '',
  slug: '',
};

function mapImageToForm(image: ImageItem): ImageFormState {
  return {
    id: image.id,
    imageURL: image.image_url,
    thumbnailURL: image.thumbnail_url || '',
    width: image.width ? String(image.width) : '',
    height: image.height ? String(image.height) : '',
    altText: image.alt_text || '',
    source: image.source || '',
  };
}

function mapTagToForm(tag: Tag): TagFormState {
  return {
    id: tag.id,
    name: tag.name,
    slug: tag.slug,
  };
}

function toImagePayload(form: ImageFormState): CreateImagePayload | UpdateImagePayload {
  const payload: CreateImagePayload = {
    image_url: form.imageURL.trim(),
  };

  if (form.thumbnailURL.trim()) {
    payload.thumbnail_url = form.thumbnailURL.trim();
  }
  if (form.width.trim()) {
    payload.width = Number(form.width);
  }
  if (form.height.trim()) {
    payload.height = Number(form.height);
  }
  if (form.altText.trim()) {
    payload.alt_text = form.altText.trim();
  }
  if (form.source.trim()) {
    payload.source = form.source.trim();
  }

  return payload;
}

function toTagPayload(form: TagFormState): CreateTagPayload | UpdateTagPayload {
  const payload: CreateTagPayload = {
    name: form.name.trim(),
  };

  if (form.slug.trim()) {
    payload.slug = form.slug.trim();
  }

  return payload;
}

export default function ManagePage() {
  const [images, setImages] = useState<ImageItem[]>([]);
  const [tags, setTags] = useState<Tag[]>([]);
  const [imageForm, setImageForm] = useState<ImageFormState>(emptyImageForm);
  const [tagForm, setTagForm] = useState<TagFormState>(emptyTagForm);
  const [tagPickerImageID, setTagPickerImageID] = useState<number | null>(null);
  const [tagPickerSelection, setTagPickerSelection] = useState<number[]>([]);
  const [isBooting, setIsBooting] = useState(true);
  const [isRefreshing, setIsRefreshing] = useState(false);
  const [isSavingImage, setIsSavingImage] = useState(false);
  const [isSavingTag, setIsSavingTag] = useState(false);
  const [isApplyingTagPicker, setIsApplyingTagPicker] = useState(false);
  const [statusMessage, setStatusMessage] = useState<string | null>(null);
  const [errorMessage, setErrorMessage] = useState<string | null>(null);

  const selectedImage = useMemo(
    () => images.find((image) => image.id === imageForm.id) ?? null,
    [imageForm.id, images],
  );
  const getAvailableTagsForImage = useCallback(
    (image: ImageItem | null) => {
      if (!image) {
        return [] as Tag[];
      }

      const attachedTagIDs = new Set(image.tags.map((tag) => tag.id));
      return tags.filter((tag) => !attachedTagIDs.has(tag.id));
    },
    [tags],
  );
  const selectedImageAvailableTags = useMemo(
    () => getAvailableTagsForImage(selectedImage),
    [getAvailableTagsForImage, selectedImage],
  );
  const tagPickerImage = useMemo(
    () => images.find((image) => image.id === tagPickerImageID) ?? null,
    [images, tagPickerImageID],
  );
  const tagPickerAvailableTags = useMemo(
    () => getAvailableTagsForImage(tagPickerImage),
    [getAvailableTagsForImage, tagPickerImage],
  );

  const refreshData = useCallback(async (options?: { silent?: boolean }) => {
    if (!options?.silent) {
      setIsRefreshing(true);
    }

    try {
      const [imageCatalog, tagCatalog] = await Promise.all([fetchAllImages(), fetchTags()]);
      setImages(imageCatalog);
      setTags(tagCatalog);
    } finally {
      if (!options?.silent) {
        setIsRefreshing(false);
      }
    }
  }, []);

  useEffect(() => {
    void (async () => {
      try {
        await refreshData({ silent: true });
      } catch (error) {
        const message = error instanceof Error ? error.message : 'Failed to load admin data.';
        setErrorMessage(message);
      } finally {
        setIsBooting(false);
      }
    })();
  }, [refreshData]);

  const resetImageForm = () => {
    setImageForm(emptyImageForm);
  };

  const resetTagForm = () => {
    setTagForm(emptyTagForm);
  };

  const handleImageFieldChange = (key: keyof ImageFormState, value: string | number | null) => {
    setImageForm((current) => ({
      ...current,
      [key]: value ?? '',
    }));
  };

  const handleTagFieldChange = (key: keyof TagFormState, value: string | number | null) => {
    setTagForm((current) => ({
      ...current,
      [key]: value ?? '',
    }));
  };

  const runMutation = async <T,>(task: () => Promise<T>, successMessage: string) => {
    setErrorMessage(null);
    setStatusMessage(null);

    try {
      const result = await task();
      await refreshData();
      setStatusMessage(successMessage);
      return result;
    } catch (error) {
      const message = error instanceof Error ? error.message : 'Operation failed.';
      setErrorMessage(message);
      throw error;
    }
  };

  const handleSaveImage = async (event: React.FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    setIsSavingImage(true);

    try {
      const payload = toImagePayload(imageForm);
      if (imageForm.id) {
        const updated = await runMutation(
          () => updateImage(imageForm.id as number, payload as UpdateImagePayload),
          'Image updated.',
        );
        setImageForm(mapImageToForm(updated));
      } else {
        const created = await runMutation(
          () => createImage(payload as CreateImagePayload),
          'Image created.',
        );
        setImageForm(mapImageToForm(created));
      }
    } finally {
      setIsSavingImage(false);
    }
  };

  const handleDeleteImage = async (image: ImageItem) => {
    if (!window.confirm(`Delete image #${image.id}?`)) {
      return;
    }

    await runMutation(() => deleteImage(image.id), `Image #${image.id} deleted.`);
    if (imageForm.id === image.id) {
      resetImageForm();
    }
  };

  const handleSaveTag = async (event: React.FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    setIsSavingTag(true);

    try {
      const payload = toTagPayload(tagForm);
      if (tagForm.id) {
        const updated = await runMutation(
          () => updateTag(tagForm.id as number, payload as UpdateTagPayload),
          'Tag updated.',
        );
        setTagForm(mapTagToForm(updated));
      } else {
        const created = await runMutation(
          () => createTag(payload as CreateTagPayload),
          'Tag created.',
        );
        setTagForm(mapTagToForm(created));
      }
    } finally {
      setIsSavingTag(false);
    }
  };

  const handleDeleteTag = async (tag: Tag) => {
    if (!window.confirm(`Delete tag #${tag.slug}?`)) {
      return;
    }

    await runMutation(() => deleteTag(tag.id), `Tag #${tag.slug} deleted.`);
    if (tagForm.id === tag.id) {
      resetTagForm();
    }
  };

  const handleDetachTag = async (tagID: number) => {
    if (!selectedImage) {
      return;
    }

    await runMutation(
      () => detachTagFromImage(selectedImage.id, tagID),
      'Tag removed from image.',
    );
  };

  const handleDetachTagFromImageCard = async (imageID: number, tagID: number) => {
    await runMutation(
      () => detachTagFromImage(imageID, tagID),
      `Tag removed from image #${imageID}.`,
    );
  };

  const openTagPicker = (image: ImageItem) => {
    setTagPickerImageID(image.id);
    setTagPickerSelection([]);
  };

  const closeTagPicker = useCallback(() => {
    setTagPickerImageID(null);
    setTagPickerSelection([]);
  }, []);

  const toggleTagPickerSelection = (tagID: number) => {
    setTagPickerSelection((current) =>
      current.includes(tagID) ? current.filter((id) => id !== tagID) : [...current, tagID],
    );
  };

  const handleApplyTagPicker = async () => {
    if (!tagPickerImage || tagPickerSelection.length === 0) {
      return;
    }

    setIsApplyingTagPicker(true);

    try {
      await runMutation(
        () =>
          Promise.all(
            tagPickerSelection.map((tagID) =>
              attachTagToImage(tagPickerImage.id, { tag_id: tagID }),
            ),
          ),
        `${tagPickerSelection.length} tag(s) added to image #${tagPickerImage.id}.`,
      );
      closeTagPicker();
    } finally {
      setIsApplyingTagPicker(false);
    }
  };

  useEffect(() => {
    if (!tagPickerImageID) {
      return;
    }

    const handleKeyDown = (event: KeyboardEvent) => {
      if (event.key === 'Escape') {
        closeTagPicker();
      }
    };

    window.addEventListener('keydown', handleKeyDown);
    return () => window.removeEventListener('keydown', handleKeyDown);
  }, [closeTagPicker, tagPickerImageID]);

  return (
    <div className="min-h-screen overflow-x-hidden bg-paper text-ink">
      <div className="pointer-events-none fixed inset-0 -z-10 bg-grain" />

      <header className="sticky top-0 z-50 border-b border-stone-900/10 bg-paper/90 backdrop-blur-xl">
        <div className="mx-auto flex max-w-7xl flex-wrap items-start justify-between gap-4 px-4 py-5 sm:px-6 lg:items-center lg:px-8">
          <div className="max-w-3xl space-y-2">
            <p className="text-[11px] uppercase tracking-[0.3em] text-clay">Management Console</p>
            <h1 className="font-display text-3xl leading-[0.92] text-ink sm:text-5xl">
              Manage images, tags, and assignments.
            </h1>
          </div>

          <div className="flex w-full flex-col gap-3 sm:w-auto sm:flex-row">
            <button type="button" className="secondary-button w-full sm:w-auto" onClick={() => void refreshData()}>
              {isRefreshing ? 'Refreshing...' : 'Refresh Data'}
            </button>
            <Link to="/" className="secondary-button w-full sm:w-auto">
              Back To Gallery
            </Link>
          </div>
        </div>
      </header>

      <main className="mx-auto max-w-7xl px-4 py-8 sm:px-6 lg:px-8">
        {statusMessage ? (
          <div className="mb-6 rounded-[1.25rem] border border-emerald-200 bg-emerald-50 px-5 py-4 text-sm text-emerald-700">
            {statusMessage}
          </div>
        ) : null}

        {errorMessage ? (
          <div className="mb-6 rounded-[1.25rem] border border-red-200 bg-red-50 px-5 py-4 text-sm text-red-700">
            {errorMessage}
          </div>
        ) : null}

        {isBooting ? (
          <div className="panel-card px-6 py-16 text-center">
            <p className="font-display text-3xl text-ink">Loading management data...</p>
          </div>
        ) : (
          <div className="grid items-start gap-6 xl:grid-cols-[1.15fr,0.85fr]">
            <section className="panel-card min-w-0 p-4 sm:p-6">
              <div className="flex flex-col gap-4 sm:flex-row sm:flex-wrap sm:items-center sm:justify-between">
                <div>
                  <p className="text-[11px] uppercase tracking-[0.28em] text-stone-500">Images</p>
                  <h2 className="mt-2 font-display text-3xl leading-none text-ink sm:text-4xl">Create Or Edit Image</h2>
                </div>
                <button type="button" className="secondary-button w-full sm:w-auto" onClick={resetImageForm}>
                  New Image
                </button>
              </div>

              <form className="mt-6 grid gap-4 md:grid-cols-2" onSubmit={handleSaveImage}>
                <label className="form-field md:col-span-2">
                  <span className="field-label">Image URL</span>
                  <input
                    required
                    className="field-input"
                    value={imageForm.imageURL}
                    onChange={(event) => handleImageFieldChange('imageURL', event.target.value)}
                    placeholder="https://placehold.co/1200x900?text=Frame+01"
                  />
                </label>

                <label className="form-field md:col-span-2">
                  <span className="field-label">Thumbnail URL</span>
                  <input
                    className="field-input"
                    value={imageForm.thumbnailURL}
                    onChange={(event) => handleImageFieldChange('thumbnailURL', event.target.value)}
                    placeholder="https://placehold.co/400x300?text=Thumb+01"
                  />
                </label>

                <label className="form-field">
                  <span className="field-label">Width</span>
                  <input
                    type="number"
                    min="1"
                    className="field-input"
                    value={imageForm.width}
                    onChange={(event) => handleImageFieldChange('width', event.target.value)}
                    placeholder="1200"
                  />
                </label>

                <label className="form-field">
                  <span className="field-label">Height</span>
                  <input
                    type="number"
                    min="1"
                    className="field-input"
                    value={imageForm.height}
                    onChange={(event) => handleImageFieldChange('height', event.target.value)}
                    placeholder="900"
                  />
                </label>

                <label className="form-field">
                  <span className="field-label">Alt Text</span>
                  <input
                    className="field-input"
                    value={imageForm.altText}
                    onChange={(event) => handleImageFieldChange('altText', event.target.value)}
                    placeholder="Gallery frame title"
                  />
                </label>

                <label className="form-field">
                  <span className="field-label">Source</span>
                  <input
                    className="field-input"
                    value={imageForm.source}
                    onChange={(event) => handleImageFieldChange('source', event.target.value)}
                    placeholder="placehold.co"
                  />
                </label>

                <div className="md:col-span-2 flex flex-col gap-3 sm:flex-row sm:flex-wrap sm:justify-end">
                  <button type="submit" className="primary-button w-full sm:w-auto" disabled={isSavingImage}>
                    {isSavingImage ? 'Saving...' : imageForm.id ? 'Update Image' : 'Create Image'}
                  </button>
                  {imageForm.id ? (
                    <button
                      type="button"
                      className="danger-button w-full sm:w-auto"
                      onClick={() => {
                        const image = images.find((entry) => entry.id === imageForm.id);
                        if (image) {
                          void handleDeleteImage(image);
                        }
                      }}
                    >
                      Delete Image
                    </button>
                  ) : null}
                </div>
              </form>

              <div className="mt-8 space-y-6">
                {selectedImage ? (
                  <div className="min-w-0 py-1">
                    <div className="flex flex-col gap-4 lg:flex-row lg:items-start lg:justify-between">
                      <div>
                        <p className="text-[11px] uppercase tracking-[0.26em] text-stone-500">Tag Assignments</p>
                        <h3 className="mt-2 text-lg font-bold text-ink">
                          {selectedImage.alt_text || `Image #${selectedImage.id}`}
                        </h3>
                        <p className="mt-1 text-sm text-stone-500">Attach new tags or remove existing ones.</p>
                      </div>

                      <div className="flex w-full flex-col items-start gap-3 lg:w-auto lg:items-end">
                        <button
                          type="button"
                          className="secondary-button w-full lg:w-auto"
                          onClick={() => openTagPicker(selectedImage)}
                          disabled={selectedImageAvailableTags.length === 0}
                        >
                          {selectedImageAvailableTags.length > 0 ? 'Add Tags' : 'All Tags Added'}
                        </button>
                        <p className="text-sm text-stone-500 lg:text-right">
                          {selectedImageAvailableTags.length > 0
                            ? `${selectedImageAvailableTags.length} tag(s) available to attach.`
                            : 'This image already has every available tag.'}
                        </p>
                      </div>
                    </div>

                    <div className="mt-4 flex flex-wrap gap-2">
                      {selectedImage.tags.length > 0 ? (
                        selectedImage.tags.map((tag) => (
                          <button
                            key={tag.id}
                            type="button"
                            onClick={() => void handleDetachTag(tag.id)}
                            className="inline-flex items-center gap-2 rounded-full border border-stone-900/10 bg-stone-100 px-3 py-2 text-xs font-medium uppercase tracking-[0.18em] text-stone-700 transition hover:border-red-300 hover:bg-red-50 hover:text-red-600"
                          >
                            <span>#{tag.slug}</span>
                            <span aria-hidden="true" className="text-sm leading-none">
                              ×
                            </span>
                          </button>
                        ))
                      ) : (
                        <p className="text-sm text-stone-500">No tags attached yet.</p>
                      )}
                    </div>
                  </div>
                ) : null}

                <div className="min-w-0 rounded-[1.5rem] border border-stone-900/10 bg-white/60 p-4">
                  <div className="flex items-center justify-between gap-3">
                    <div>
                      <p className="text-[11px] uppercase tracking-[0.26em] text-stone-500">Catalog</p>
                      <h3 className="mt-2 text-lg font-bold text-ink">{images.length} active images</h3>
                    </div>
                  </div>

                  <div className="mt-4 max-h-[36rem] space-y-3 overflow-y-auto pr-1">
                    {images.map((image) => (
                      <article
                        key={image.id}
                        className={`rounded-[1.5rem] border p-4 transition ${
                          imageForm.id === image.id
                            ? 'border-transparent bg-ink text-paper shadow-card'
                            : 'border-stone-900/10 bg-white/80 text-ink hover:-translate-y-0.5 hover:border-clay/30'
                        }`}
                      >
                        <div className="flex flex-col gap-4">
                          <div className="flex items-start gap-4">
                            <img
                              src={image.thumbnail_url || image.image_url}
                              alt={image.alt_text || `Image ${image.id}`}
                              className="h-20 w-20 shrink-0 rounded-2xl object-cover"
                            />
                            <div className="min-w-0 flex-1">
                              <div className="flex items-start justify-between gap-3">
                                <div className="min-w-0 flex-1">
                                  <p className="truncate text-lg font-semibold">
                                    {image.alt_text || `Image #${image.id}`}
                                  </p>
                                  <p className="mt-1 truncate text-sm opacity-75">{image.image_url}</p>
                                </div>
                                <div className="flex shrink-0 flex-nowrap items-center gap-2">
                                  <button
                                    type="button"
                                    className={imageForm.id === image.id ? 'add-button-active' : 'add-button'}
                                    onClick={() => openTagPicker(image)}
                                  >
                                    Add Tag
                                  </button>
                                  <button
                                    type="button"
                                    className={imageForm.id === image.id ? 'edit-button-active' : 'edit-button'}
                                    onClick={() => {
                                      setImageForm(mapImageToForm(image));
                                      window.scrollTo({ top: 0, behavior: 'smooth' });
                                    }}
                                  >
                                    Edit
                                  </button>
                                  <button
                                    type="button"
                                    className={`rounded-full px-3 py-2 text-xs font-bold uppercase tracking-[0.18em] ${
                                      imageForm.id === image.id
                                        ? 'border border-red-300/40 bg-red-500/10 text-red-100'
                                        : 'border border-red-200 bg-red-50 text-red-600'
                                    }`}
                                    onClick={() => void handleDeleteImage(image)}
                                  >
                                    Delete
                                  </button>
                                </div>
                              </div>
                            </div>
                          </div>

                          <div className="flex flex-wrap gap-2">
                            {image.tags.map((tag) => (
                              <button
                                key={`${image.id}-${tag.id}`}
                                type="button"
                                onClick={() => void handleDetachTagFromImageCard(image.id, tag.id)}
                                className={`inline-flex items-center gap-2 rounded-full px-3 py-1 text-[11px] font-medium uppercase tracking-[0.18em] transition ${
                                  imageForm.id === image.id
                                    ? 'border border-white/20 bg-white/10 text-paper hover:border-red-300/40 hover:bg-red-500/10 hover:text-red-100'
                                    : 'border border-stone-200 bg-stone-100 text-stone-700 hover:border-red-200 hover:bg-red-50 hover:text-red-600'
                                }`}
                                title={`Remove #${tag.slug} from image`}
                              >
                                <span>#{tag.slug}</span>
                                <span aria-hidden="true" className="text-sm leading-none">
                                  ×
                                </span>
                              </button>
                            ))}
                          </div>
                        </div>
                      </article>
                    ))}
                  </div>
                </div>
              </div>
            </section>

            <section className="panel-card min-w-0 p-4 sm:p-6">
              <div className="flex flex-col gap-4 sm:flex-row sm:flex-wrap sm:items-center sm:justify-between">
                <div>
                  <p className="text-[11px] uppercase tracking-[0.28em] text-stone-500">Tags</p>
                  <h2 className="mt-2 font-display text-3xl leading-none text-ink sm:text-4xl">Create Or Edit Tag</h2>
                </div>
                <button type="button" className="secondary-button w-full sm:w-auto" onClick={resetTagForm}>
                  New Tag
                </button>
              </div>

              <form className="mt-6 grid gap-4" onSubmit={handleSaveTag}>
                <label className="form-field">
                  <span className="field-label">Tag Name</span>
                  <input
                    required
                    className="field-input"
                    value={tagForm.name}
                    onChange={(event) => handleTagFieldChange('name', event.target.value)}
                    placeholder="Nature"
                  />
                </label>

                <label className="form-field">
                  <span className="field-label">Slug</span>
                  <input
                    className="field-input"
                    value={tagForm.slug}
                    onChange={(event) => handleTagFieldChange('slug', event.target.value)}
                    placeholder="nature"
                  />
                </label>

                <div className="flex flex-col gap-3 sm:flex-row sm:flex-wrap">
                  <button type="submit" className="primary-button w-full sm:w-auto" disabled={isSavingTag}>
                    {isSavingTag ? 'Saving...' : tagForm.id ? 'Update Tag' : 'Create Tag'}
                  </button>
                  {tagForm.id ? (
                    <button
                      type="button"
                      className="danger-button w-full sm:w-auto"
                      onClick={() => {
                        const tag = tags.find((entry) => entry.id === tagForm.id);
                        if (tag) {
                          void handleDeleteTag(tag);
                        }
                      }}
                    >
                      Delete Tag
                    </button>
                  ) : null}
                </div>
              </form>

              <div className="mt-8 rounded-[1.5rem] border border-stone-900/10 bg-white/60 p-4">
                <p className="text-[11px] uppercase tracking-[0.26em] text-stone-500">Catalog</p>
                <h3 className="mt-2 text-lg font-bold text-ink">{tags.length} active tags</h3>

                <div className="mt-4 max-h-[52rem] space-y-3 overflow-y-auto pr-1">
                  {tags.map((tag) => (
                    <article
                      key={tag.id}
                      className={`rounded-[1.5rem] border p-4 transition ${
                        tagForm.id === tag.id
                          ? 'border-transparent bg-ink text-paper shadow-card'
                          : 'border-stone-900/10 bg-white/80 text-ink hover:-translate-y-0.5 hover:border-clay/30'
                      }`}
                    >
                      <div className="flex flex-col gap-3 sm:flex-row sm:items-start sm:justify-between">
                        <div>
                          <p className="text-lg font-semibold">{tag.name}</p>
                          <p className="mt-1 text-sm opacity-75">#{tag.slug}</p>
                        </div>
                        <div className="flex shrink-0 flex-nowrap items-center gap-2">
                          <button
                            type="button"
                            className={tagForm.id === tag.id ? 'edit-button-active' : 'edit-button'}
                            onClick={() => {
                              setTagForm(mapTagToForm(tag));
                              window.scrollTo({ top: 0, behavior: 'smooth' });
                            }}
                          >
                            Edit
                          </button>
                          <button
                            type="button"
                            className={`rounded-full px-3 py-2 text-xs font-bold uppercase tracking-[0.18em] ${
                              tagForm.id === tag.id
                                ? 'border border-red-300/40 bg-red-500/10 text-red-100'
                                : 'border border-red-200 bg-red-50 text-red-600'
                            }`}
                            onClick={() => void handleDeleteTag(tag)}
                          >
                            Delete
                          </button>
                        </div>
                      </div>
                    </article>
                  ))}
                </div>
              </div>
            </section>
          </div>
        )}
      </main>

      {tagPickerImage ? (
        <div
          className="fixed inset-0 z-[80] flex items-center justify-center bg-ink/45 px-4 py-8 backdrop-blur-sm"
          onClick={closeTagPicker}
        >
          <div
            className="w-full max-w-3xl overflow-hidden rounded-[2rem] border border-stone-900/10 bg-paper shadow-[0_32px_80px_rgba(21,19,18,0.24)]"
            onClick={(event) => event.stopPropagation()}
          >
            <div className="flex items-start justify-between gap-4 border-b border-stone-900/10 px-5 py-5 sm:px-6">
              <div className="min-w-0">
                <p className="text-[11px] uppercase tracking-[0.28em] text-clay">Add Tags</p>
                <h2 className="mt-2 truncate font-display text-3xl leading-none text-ink sm:text-4xl">
                  {tagPickerImage.alt_text || `Image #${tagPickerImage.id}`}
                </h2>
                <p className="mt-2 text-sm text-stone-600">
                  Select one or more tags to attach. Existing tags are hidden from this list.
                </p>
              </div>
              <button
                type="button"
                onClick={closeTagPicker}
                className="inline-flex size-11 items-center justify-center rounded-full border border-stone-900/10 bg-white/80 text-xl text-stone-700 transition hover:border-clay/30 hover:text-clay"
                aria-label="Close add tag modal"
              >
                ×
              </button>
            </div>

            <div className="max-h-[60vh] overflow-y-auto px-5 py-5 sm:px-6">
              {tagPickerAvailableTags.length > 0 ? (
                <div className="flex flex-wrap gap-3">
                  {tagPickerAvailableTags.map((tag) => {
                    const isSelected = tagPickerSelection.includes(tag.id);

                    return (
                      <button
                        key={tag.id}
                        type="button"
                        onClick={() => toggleTagPickerSelection(tag.id)}
                        className={`rounded-full border px-4 py-3 text-sm font-bold uppercase tracking-[0.18em] transition ${
                          isSelected
                            ? 'border-transparent bg-ink text-paper shadow-card'
                            : 'border-stone-900/10 bg-white/85 text-stone-700 hover:-translate-y-0.5 hover:border-clay/30 hover:text-clay'
                        }`}
                      >
                        #{tag.slug}
                      </button>
                    );
                  })}
                </div>
              ) : (
                <div className="rounded-[1.5rem] border border-stone-900/10 bg-white/70 px-5 py-8 text-center">
                  <p className="font-display text-3xl text-ink">No tags left to add.</p>
                  <p className="mt-2 text-sm text-stone-600">
                    This image already has every active tag in the catalog.
                  </p>
                </div>
              )}
            </div>

            <div className="flex flex-col gap-3 border-t border-stone-900/10 px-5 py-5 sm:flex-row sm:items-center sm:justify-between sm:px-6">
              <p className="text-sm text-stone-600">
                {tagPickerSelection.length > 0
                  ? `${tagPickerSelection.length} tag(s) selected.`
                  : 'Choose tags to attach to this image.'}
              </p>
              <div className="flex flex-col gap-3 sm:flex-row">
                <button type="button" className="secondary-button w-full sm:w-auto" onClick={closeTagPicker}>
                  Cancel
                </button>
                <button
                  type="button"
                  className="primary-button w-full sm:w-auto"
                  disabled={tagPickerSelection.length === 0 || isApplyingTagPicker}
                  onClick={() => void handleApplyTagPicker()}
                >
                  {isApplyingTagPicker ? 'Adding...' : 'Add Selected Tags'}
                </button>
              </div>
            </div>
          </div>
        </div>
      ) : null}
    </div>
  );
}
