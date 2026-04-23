import Masonry from 'masonry-layout';
import imagesLoaded from 'imagesloaded';
import { useCallback, useEffect, useMemo, useRef, useState } from 'react';
import { Link } from 'react-router-dom';
import { fetchImages, fetchTags } from '../api';
import type { ImageItem, Tag } from '../types';

const CHUNK_OPTIONS = [10, 30, 50] as const;
const DEFAULT_CHUNK_SIZE = 10;

type GalleryState = {
  images: ImageItem[];
  nextCursor: number | null;
  hasMore: boolean;
};

function LoadingCard() {
  return (
    <div className="gallery-item">
      <div className="gallery-card animate-pulse overflow-hidden">
        <div className="aspect-[3/4] bg-stone-200/80" />
        <div className="space-y-3 p-4">
          <div className="h-4 w-2/3 rounded-full bg-stone-200/90" />
          <div className="h-3 w-1/2 rounded-full bg-stone-200/70" />
          <div className="flex gap-2">
            <div className="h-7 w-20 rounded-full bg-stone-200/80" />
            <div className="h-7 w-24 rounded-full bg-stone-200/70" />
          </div>
        </div>
      </div>
    </div>
  );
}

function ImageCard({
  image,
  selectedTag,
  onTagSelect,
}: {
  image: ImageItem;
  selectedTag: string;
  onTagSelect: (slug: string) => void;
}) {
  const ratio = image.width && image.height ? `${image.width} / ${image.height}` : '4 / 5';
  const imageSrc = image.thumbnail_url || image.image_url;

  return (
    <div className="gallery-item">
      <article data-testid={`gallery-card-${image.id}`} className="gallery-card overflow-hidden">
        <div className="relative bg-stone-100" style={{ aspectRatio: ratio }}>
          <img
            src={imageSrc}
            alt={image.alt_text || `Gallery item ${image.id}`}
            className="h-full w-full object-cover"
            loading="lazy"
          />
          <div className="pointer-events-none absolute inset-x-0 bottom-0 h-24 bg-gradient-to-t from-black/50 via-black/15 to-transparent" />
          <div className="absolute bottom-3 left-3 right-3 space-y-2 text-paper">
            <p className="truncate font-display text-lg leading-none">
              {image.alt_text || `Frame ${image.id}`}
            </p>
            <div className="flex flex-wrap items-center gap-2">
              <p className="text-xs uppercase tracking-[0.3em] text-paper/78">
                {image.width && image.height ? `${image.width} x ${image.height}` : 'Open format'}
              </p>
              <span className="max-w-full truncate rounded-full border border-white/15 bg-black/45 px-2.5 py-1 text-[10px] uppercase tracking-[0.2em] text-white backdrop-blur-xl">
                {image.source || 'gallery'}
              </span>
            </div>
          </div>
        </div>
        <div className="space-y-4 p-4">
          <div className="flex flex-wrap gap-2">
            {image.tags.map((tag) => (
              <button
                key={`${image.id}-${tag.id}`}
                type="button"
                onClick={() => onTagSelect(tag.slug)}
                className={`rounded-full border px-3 py-1 text-xs font-semibold uppercase tracking-[0.16em] transition ${selectedTag === tag.slug
                    ? 'border-transparent bg-clay text-white shadow-[0_10px_24px_rgba(0,113,227,0.22)]'
                    : 'border-black/10 bg-white/92 text-stone-700 hover:border-clay/30 hover:text-clay'
                  }`}
              >
                #{tag.slug}
              </button>
            ))}
          </div>
        </div>
      </article>
    </div>
  );
}

async function getGalleryChunk(
  cursor: number | null,
  tag: string,
  limit: number,
): Promise<GalleryState> {
  const images = await fetchImages({ cursor, tag, limit });

  return {
    images,
    nextCursor: images.length === limit ? images[images.length - 1]?.id ?? null : null,
    hasMore: images.length === limit,
  };
}

export default function GalleryPage() {
  const [tags, setTags] = useState<Tag[]>([]);
  const [selectedTag, setSelectedTag] = useState('');
  const [isTagModalOpen, setIsTagModalOpen] = useState(false);
  const [isMobileMenuOpen, setIsMobileMenuOpen] = useState(false);
  const [pageSize, setPageSize] = useState<number>(DEFAULT_CHUNK_SIZE);
  const [images, setImages] = useState<ImageItem[]>([]);
  const [nextCursor, setNextCursor] = useState<number | null>(null);
  const [hasMore, setHasMore] = useState(true);
  const [isInitialLoading, setIsInitialLoading] = useState(true);
  const [isLoadingMore, setIsLoadingMore] = useState(false);
  const [isTagsLoading, setIsTagsLoading] = useState(true);
  const [errorMessage, setErrorMessage] = useState<string | null>(null);

  const gridRef = useRef<HTMLDivElement | null>(null);
  const sentinelRef = useRef<HTMLDivElement | null>(null);
  const masonryRef = useRef<Masonry | null>(null);
  const requestIdRef = useRef(0);

  const titleCount = useMemo(() => `${images.length}`.padStart(2, '0'), [images.length]);
  const handleSelectTag = useCallback((slug: string) => {
    setSelectedTag(slug);
    setIsTagModalOpen(false);
    setIsMobileMenuOpen(false);
    window.scrollTo({ top: 0, behavior: 'smooth' });
  }, []);

  const handleResetTag = useCallback(() => {
    setSelectedTag('');
    setIsTagModalOpen(false);
    setIsMobileMenuOpen(false);
    window.scrollTo({ top: 0, behavior: 'smooth' });
  }, []);

  const handleOpenTagModal = useCallback(() => {
    setIsMobileMenuOpen(false);
    setIsTagModalOpen(true);
  }, []);

  const loadTags = useCallback(async () => {
    setIsTagsLoading(true);
    try {
      const catalog = await fetchTags();
      setTags(catalog);
    } catch (error) {
      const message = error instanceof Error ? error.message : 'Failed to load tags.';
      setErrorMessage(message);
    } finally {
      setIsTagsLoading(false);
    }
  }, []);

  const loadImages = useCallback(
    async ({ cursor, replace }: { cursor: number | null; replace: boolean }) => {
      const requestId = ++requestIdRef.current;

      if (replace) {
        setIsInitialLoading(true);
      } else {
        setIsLoadingMore(true);
      }

      setErrorMessage(null);

      try {
        const payload = await getGalleryChunk(cursor, selectedTag, pageSize);
        if (requestId !== requestIdRef.current) {
          return;
        }

        setImages((current) => (replace ? payload.images : [...current, ...payload.images]));
        setNextCursor(payload.nextCursor);
        setHasMore(payload.hasMore);
      } catch (error) {
        if (requestId !== requestIdRef.current) {
          return;
        }
        const message = error instanceof Error ? error.message : 'Failed to load images.';
        setErrorMessage(message);
      } finally {
        if (requestId === requestIdRef.current) {
          setIsInitialLoading(false);
          setIsLoadingMore(false);
        }
      }
    },
    [pageSize, selectedTag],
  );

  useEffect(() => {
    void loadTags();
  }, [loadTags]);

  useEffect(() => {
    setImages([]);
    setNextCursor(null);
    setHasMore(true);
    void loadImages({ cursor: null, replace: true });
  }, [loadImages]);

  useEffect(() => {
    const grid = gridRef.current;
    if (!grid) {
      return;
    }

    if (!masonryRef.current) {
      masonryRef.current = new Masonry(grid, {
        itemSelector: '.gallery-item',
        columnWidth: '.gallery-sizer',
        percentPosition: true,
        transitionDuration: '0.18s',
      });
    } else {
      masonryRef.current.reloadItems?.();
      masonryRef.current.layout?.();
    }

    const imageLoader = imagesLoaded(grid);
    const relayout = () => masonryRef.current?.layout?.();

    imageLoader.on?.('progress', relayout);
    imageLoader.on?.('always', relayout);
    relayout();

    return () => {
      imageLoader.off?.('progress', relayout);
      imageLoader.off?.('always', relayout);
    };
  }, [images, isInitialLoading]);

  useEffect(() => {
    return () => {
      masonryRef.current?.destroy?.();
      masonryRef.current = null;
    };
  }, []);

  useEffect(() => {
    if (!isTagModalOpen && !isMobileMenuOpen) {
      return;
    }

    const handleKeyDown = (event: KeyboardEvent) => {
      if (event.key === 'Escape') {
        if (isTagModalOpen) {
          setIsTagModalOpen(false);
          return;
        }

        setIsMobileMenuOpen(false);
      }
    };

    window.addEventListener('keydown', handleKeyDown);
    return () => window.removeEventListener('keydown', handleKeyDown);
  }, [isMobileMenuOpen, isTagModalOpen]);

  useEffect(() => {
    const mediaQuery = window.matchMedia('(min-width: 1024px)');
    const handleChange = (event: MediaQueryListEvent) => {
      if (event.matches) {
        setIsMobileMenuOpen(false);
      }
    };

    if (mediaQuery.matches) {
      setIsMobileMenuOpen(false);
    }

    mediaQuery.addEventListener('change', handleChange);
    return () => mediaQuery.removeEventListener('change', handleChange);
  }, []);

  useEffect(() => {
    const sentinel = sentinelRef.current;
    if (!sentinel || isInitialLoading || isLoadingMore || !hasMore) {
      return;
    }

    const observer = new IntersectionObserver(
      (entries) => {
        const target = entries[0];
        if (!target?.isIntersecting || isLoadingMore || !hasMore) {
          return;
        }

        void loadImages({ cursor: nextCursor, replace: false });
      },
      {
        rootMargin: '320px 0px 320px 0px',
        threshold: 0.01,
      },
    );

    observer.observe(sentinel);

    return () => observer.disconnect();
  }, [hasMore, isInitialLoading, isLoadingMore, loadImages, nextCursor]);

  return (
    <div className="gallery-page min-h-screen overflow-x-hidden bg-white text-ink">

      <header className="fixed inset-x-0 top-0 z-50 border-b border-stone-900/10 bg-white/95 backdrop-blur-xl">
        <div className="mx-auto max-w-7xl px-4 py-4 sm:px-6 lg:hidden">
          <div className="flex items-start justify-between gap-4">
            <h1 className="min-w-0 max-w-[12rem] font-display text-3xl leading-[0.92] text-ink sm:max-w-[18rem] sm:text-4xl">
              Masonry image gallery.
            </h1>

            <button
              type="button"
              onClick={() => setIsMobileMenuOpen((current) => !current)}
              className="gallery-shadowless inline-flex size-14 shrink-0 items-center justify-center rounded-full border border-stone-900/10 bg-white/80 text-ink shadow-card transition hover:border-clay/30 hover:text-clay"
              aria-label="Toggle gallery menu"
              aria-expanded={isMobileMenuOpen}
              aria-controls="mobile-gallery-menu"
            >
              <span className="relative flex h-5 w-6 flex-col justify-between">
                <span
                  className={`block h-0.5 rounded-full bg-current transition ${isMobileMenuOpen ? 'translate-y-[9px] rotate-45' : ''
                    }`}
                />
                <span
                  className={`block h-0.5 rounded-full bg-current transition ${isMobileMenuOpen ? 'opacity-0' : ''
                    }`}
                />
                <span
                  className={`block h-0.5 rounded-full bg-current transition ${isMobileMenuOpen ? '-translate-y-[9px] -rotate-45' : ''
                    }`}
                />
              </span>
            </button>
          </div>

          {isMobileMenuOpen ? (
            <div
              id="mobile-gallery-menu"
              className="gallery-shadowless mt-4 max-h-[calc(100vh-7rem)] space-y-4 overflow-y-auto rounded-[1.75rem] border border-black/5 bg-white/95 p-4 shadow-[0_24px_54px_rgba(0,0,0,0.12)]"
            >
              <Link
                to="/manage"
                className="secondary-button w-full"
                onClick={() => setIsMobileMenuOpen(false)}
              >
                Manage Gallery
              </Link>

              <div className="grid grid-cols-2 gap-3">
                <div className="panel-card min-w-0">
                  <p className="text-[11px] uppercase tracking-[0.25em] text-stone-500">Loaded</p>
                  <p className="mt-3 truncate font-sans text-3xl font-bold text-ink">{titleCount}</p>
                </div>
                <div className="panel-card min-w-0">
                  <p className="text-[11px] uppercase tracking-[0.25em] text-stone-500">Filter</p>
                  <p className="mt-3 truncate font-sans text-2xl font-bold text-ink">
                    {selectedTag ? `#${selectedTag}` : 'All'}
                  </p>
                </div>
                <div className="panel-card col-span-2 min-w-0">
                  <p className="text-[11px] uppercase tracking-[0.25em] text-stone-500">Chunk size</p>
                  <div className="mt-3 flex flex-wrap gap-2 pr-2">
                    {CHUNK_OPTIONS.map((option) => (
                      <button
                        key={option}
                        type="button"
                        onClick={() => setPageSize(option)}
                        className={`rounded-full px-3 py-2 text-sm font-bold transition ${pageSize === option
                            ? 'border-transparent bg-clay text-white shadow-[0_10px_24px_rgba(0,113,227,0.22)]'
                            : 'border border-black/10 bg-white/90 text-stone-700 hover:border-clay/30 hover:text-clay'
                          }`}
                      >
                        {option}
                      </button>
                    ))}
                  </div>
                </div>
              </div>

              <div className="flex flex-wrap gap-3">
                <button
                  type="button"
                  onClick={handleResetTag}
                  className={`tag-chip ${selectedTag === '' ? 'tag-chip-active' : ''}`}
                >
                  All Frames
                </button>

                {!isTagsLoading ? (
                  <button
                    type="button"
                    onClick={handleOpenTagModal}
                    className={`tag-chip ${isTagModalOpen ? 'tag-chip-active' : ''}`}
                  >
                    Tags Filter
                  </button>
                ) : null}
              </div>
            </div>
          ) : null}
        </div>

        <div className="mx-auto hidden max-w-7xl px-4 pb-5 pt-5 sm:px-6 lg:block lg:px-8">
          <div className="grid gap-6 lg:grid-cols-[1.4fr,1fr] lg:items-end">
            <div className="space-y-4">
              <div className="space-y-2">
                <h1 className="max-w-2xl font-display text-3xl leading-[0.92] text-ink sm:text-5xl lg:text-6xl">
                  Masonry image gallery.
                </h1>
              </div>

              <div>
                <Link to="/manage" className="secondary-button w-full sm:w-auto">
                  Manage Gallery
                </Link>
              </div>
            </div>

            <div className="space-y-4">
              <div className="grid grid-cols-1 gap-3 lg:grid-cols-[minmax(0,0.95fr)_minmax(0,1.15fr)_minmax(0,0.95fr)]">
                <div className="panel-card">
                  <p className="text-[11px] uppercase tracking-[0.25em] text-stone-500">Loaded</p>
                  <p className="mt-3 font-sans text-3xl font-bold text-ink">{titleCount}</p>
                </div>
                <div className="panel-card min-w-0 px-5">
                  <p className="text-[11px] uppercase tracking-[0.25em] text-stone-500">Chunk size</p>
                  <div className="mt-3 flex flex-nowrap justify-center gap-2">
                    {CHUNK_OPTIONS.map((option) => (
                      <button
                        key={option}
                        type="button"
                        onClick={() => setPageSize(option)}
                        className={`rounded-full px-3 py-2 text-sm font-bold transition ${pageSize === option
                            ? 'border-transparent bg-clay text-white shadow-[0_10px_24px_rgba(0,113,227,0.22)]'
                            : 'border border-black/10 bg-white/90 text-stone-700 hover:border-clay/30 hover:text-clay'
                          }`}
                      >
                        {option}
                      </button>
                    ))}
                  </div>
                </div>
                <div className="panel-card">
                  <p className="text-[11px] uppercase tracking-[0.25em] text-stone-500">Filter</p>
                  <p className="mt-3 truncate font-sans text-2xl font-bold text-ink">
                    {selectedTag ? `#${selectedTag}` : 'All'}
                  </p>
                </div>
              </div>

              <div className="flex flex-wrap justify-start gap-3 pb-1 lg:justify-end">
                <button
                  type="button"
                  onClick={handleResetTag}
                  className={`tag-chip ${selectedTag === '' ? 'tag-chip-active' : ''}`}
                >
                  All Frames
                </button>

                {!isTagsLoading ? (
                  <button
                    type="button"
                    onClick={handleOpenTagModal}
                    className={`tag-chip ${isTagModalOpen ? 'tag-chip-active' : ''}`}
                  >
                    Tags Filter
                  </button>
                ) : null}
              </div>
            </div>
          </div>
        </div>
      </header>

      {isMobileMenuOpen ? (
        <button
          type="button"
          className="fixed inset-0 z-40 bg-ink/12 lg:hidden"
          onClick={() => setIsMobileMenuOpen(false)}
          aria-label="Close gallery menu overlay"
        />
      ) : null}

      <main className="mx-auto max-w-7xl px-4 pb-20 pt-28 sm:px-6 sm:pt-32 lg:px-8 lg:pt-[14.5rem]">
        {errorMessage ? (
          <div
            data-testid="gallery-error"
            className="mb-8 rounded-[1.5rem] border border-red-200 bg-red-50 px-5 py-4 text-sm text-red-700"
          >
            {errorMessage}
          </div>
        ) : null}

        <div ref={gridRef} className="gallery-grid" data-testid="gallery-grid">
          <div className="gallery-sizer" />

          {isInitialLoading
            ? Array.from({ length: pageSize }).map((_, index) => <LoadingCard key={index} />)
            : images.map((image) => (
              <ImageCard
                key={image.id}
                image={image}
                selectedTag={selectedTag}
                onTagSelect={handleSelectTag}
              />
            ))}
        </div>

        {!isInitialLoading && images.length === 0 && !errorMessage ? (
          <div className="panel-card px-6 py-12 text-center">
            <p className="font-display text-3xl text-ink">No frames found.</p>
            <p className="mt-2 text-sm text-stone-600">
              Try another tag or seed more data in the backend.
            </p>
          </div>
        ) : null}

        <div ref={sentinelRef} className="h-12" />

        {isLoadingMore ? (
          <div className="mt-6 flex items-center justify-center gap-3 rounded-full border border-black/5 bg-white/92 px-5 py-3 text-sm text-stone-700 shadow-card">
            <span className="inline-flex size-3 animate-spin rounded-full border-2 border-clay border-t-transparent" />
            Loading {pageSize} more images...
          </div>
        ) : null}

        {!hasMore && !isInitialLoading && images.length > 0 ? (
          <div className="mt-6 rounded-full border border-black/5 bg-white/92 px-5 py-3 text-center text-sm text-stone-600 shadow-card">
            End of the gallery feed.
          </div>
        ) : null}
      </main>

      {isTagModalOpen ? (
        <div
          data-testid="gallery-tag-modal"
          className="fixed inset-0 z-[70] flex items-center justify-center bg-ink/24 px-4 py-8 backdrop-blur-sm"
          onClick={() => setIsTagModalOpen(false)}
        >
          <div
            className="max-h-[90vh] w-full max-w-3xl overflow-hidden rounded-[2rem] border border-black/5 bg-white/95 shadow-[0_32px_80px_rgba(0,0,0,0.16)]"
            onClick={(event) => event.stopPropagation()}
          >
            <div className="flex items-start justify-between gap-4 border-b border-stone-900/10 px-4 py-4 sm:px-6 sm:py-5">
              <div>
                <p className="text-[11px] uppercase tracking-[0.28em] text-stone-500">Tag Directory</p>
                <h2 className="mt-2 font-display text-3xl leading-none text-ink sm:text-4xl">Tags Filter</h2>
                <p className="mt-2 text-sm text-stone-600">Pick any tag to filter the gallery feed.</p>
              </div>
              <button
                type="button"
                onClick={() => setIsTagModalOpen(false)}
                className="inline-flex size-11 items-center justify-center rounded-full border border-stone-900/10 bg-white/80 text-xl text-stone-700 transition hover:border-clay/30 hover:text-clay"
                aria-label="Close tag modal"
              >
                ×
              </button>
            </div>

            <div className="max-h-[calc(90vh-120px)] overflow-y-auto px-4 py-4 sm:px-6 sm:py-6">
              <div className="mb-4">
                <button
                  type="button"
                  onClick={handleResetTag}
                  className={`tag-chip ${selectedTag === '' ? 'tag-chip-active' : ''}`}
                >
                  All Frames
                </button>
              </div>

              <div className="grid gap-3 sm:grid-cols-2 lg:grid-cols-3">
                {tags.map((tag) => (
                  <button
                    key={tag.id}
                    type="button"
                    onClick={() => handleSelectTag(tag.slug)}
                    className={`rounded-[1.5rem] border px-4 py-4 text-left transition ${selectedTag === tag.slug
                        ? 'border-transparent bg-clay text-white shadow-[0_16px_32px_rgba(0,113,227,0.2)]'
                        : 'border-black/5 bg-white/80 text-stone-800 hover:-translate-y-0.5 hover:border-clay/30 hover:text-clay'
                      }`}
                  >
                    <p className="text-[11px] uppercase tracking-[0.26em] opacity-70">Keyword</p>
                    <p className="mt-3 font-display text-3xl leading-none">#{tag.slug}</p>
                    <p className="mt-2 text-sm font-medium">{tag.name}</p>
                  </button>
                ))}
              </div>
            </div>
          </div>
        </div>
      ) : null}
    </div>
  );
}
