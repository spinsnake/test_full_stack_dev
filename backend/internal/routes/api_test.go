package routes

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sort"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"

	"github.com/example/test-full-stack-developer/backend/internal/adapter/handler"
	"github.com/example/test-full-stack-developer/backend/internal/apperrors"
	"github.com/example/test-full-stack-developer/backend/internal/entities"
	"github.com/example/test-full-stack-developer/backend/internal/port"
	"github.com/example/test-full-stack-developer/backend/internal/service"
)

func TestImageCRUDAndListLimit(t *testing.T) {
	app := newTestApp()

	firstImage := createImageViaAPI(t, app, entities.CreateImageInput{
		ImageURL: "https://placehold.co/1200x900?text=Image+01",
		Width:    intPtr(1200),
		Height:   intPtr(900),
		Source:   stringPtr("placehold.co"),
	})

	secondImage := createImageViaAPI(t, app, entities.CreateImageInput{
		ImageURL: "https://placehold.co/900x1200?text=Image+02",
		Width:    intPtr(900),
		Height:   intPtr(1200),
		Source:   stringPtr("placehold.co"),
	})

	status, body := performJSONRequest(t, app, http.MethodGet, "/api/images?limit=1", nil)
	if status != http.StatusOK {
		t.Fatalf("expected list status 200, got %d: %s", status, string(body))
	}

	var listResponse dataResponse[[]entities.Image]
	decodeJSON(t, body, &listResponse)

	if len(listResponse.Data) != 1 {
		t.Fatalf("expected 1 image in list, got %d", len(listResponse.Data))
	}
	if listResponse.Data[0].ID != secondImage.ID {
		t.Fatalf("expected latest image id %d, got %d", secondImage.ID, listResponse.Data[0].ID)
	}

	status, body = performJSONRequest(t, app, http.MethodPatch, fmt.Sprintf("/api/images/%d", firstImage.ID), entities.UpdateImageInput{
		AltText: stringPtr("updated alt text"),
	})
	if status != http.StatusOK {
		t.Fatalf("expected update status 200, got %d: %s", status, string(body))
	}

	var updatedResponse dataResponse[entities.Image]
	decodeJSON(t, body, &updatedResponse)

	if updatedResponse.Data.AltText == nil || *updatedResponse.Data.AltText != "updated alt text" {
		t.Fatalf("expected updated alt_text, got %#v", updatedResponse.Data.AltText)
	}

	status, body = performJSONRequest(t, app, http.MethodDelete, fmt.Sprintf("/api/images/%d", firstImage.ID), nil)
	if status != http.StatusOK {
		t.Fatalf("expected delete status 200, got %d: %s", status, string(body))
	}

	status, body = performJSONRequest(t, app, http.MethodGet, fmt.Sprintf("/api/images/%d", firstImage.ID), nil)
	if status != http.StatusNotFound {
		t.Fatalf("expected get deleted image status 404, got %d: %s", status, string(body))
	}
}

func TestTagCRUD(t *testing.T) {
	app := newTestApp()

	tag := createTagViaAPI(t, app, entities.CreateTagInput{
		Name: "Travel",
		Slug: stringPtr("travel"),
	})

	status, body := performJSONRequest(t, app, http.MethodPatch, fmt.Sprintf("/api/tags/%d", tag.ID), entities.UpdateTagInput{
		Name: stringPtr("City"),
		Slug: stringPtr("city"),
	})
	if status != http.StatusOK {
		t.Fatalf("expected tag update status 200, got %d: %s", status, string(body))
	}

	var updatedResponse dataResponse[entities.Tag]
	decodeJSON(t, body, &updatedResponse)

	if updatedResponse.Data.Name != "City" || updatedResponse.Data.Slug != "city" {
		t.Fatalf("expected updated tag City/city, got %s/%s", updatedResponse.Data.Name, updatedResponse.Data.Slug)
	}

	status, body = performJSONRequest(t, app, http.MethodDelete, fmt.Sprintf("/api/tags/%d", tag.ID), nil)
	if status != http.StatusOK {
		t.Fatalf("expected tag delete status 200, got %d: %s", status, string(body))
	}

	status, body = performJSONRequest(t, app, http.MethodGet, fmt.Sprintf("/api/tags/%d", tag.ID), nil)
	if status != http.StatusNotFound {
		t.Fatalf("expected deleted tag status 404, got %d: %s", status, string(body))
	}
}

func TestImageTagAttachAndDetach(t *testing.T) {
	app := newTestApp()

	createImageViaAPI(t, app, entities.CreateImageInput{
		ImageURL: "https://placehold.co/800x600?text=Image+01",
		Source:   stringPtr("placehold.co"),
	})

	targetImage := createImageViaAPI(t, app, entities.CreateImageInput{
		ImageURL: "https://placehold.co/800x600?text=Image+02",
		Source:   stringPtr("placehold.co"),
	})

	firstTag := createTagViaAPI(t, app, entities.CreateTagInput{
		Name: "Nature",
		Slug: stringPtr("nature"),
	})

	secondTag := createTagViaAPI(t, app, entities.CreateTagInput{
		Name: "Travel",
		Slug: stringPtr("travel"),
	})

	for _, tagID := range []uint64{firstTag.ID, secondTag.ID} {
		status, body := performJSONRequest(t, app, http.MethodPost, fmt.Sprintf("/api/images/%d/tags", targetImage.ID), entities.AttachTagToImageInput{
			TagID: tagID,
		})
		if status != http.StatusCreated {
			t.Fatalf("expected attach status 201, got %d: %s", status, string(body))
		}
	}

	status, body := performJSONRequest(t, app, http.MethodGet, fmt.Sprintf("/api/images/%d", targetImage.ID), nil)
	if status != http.StatusOK {
		t.Fatalf("expected image detail status 200, got %d: %s", status, string(body))
	}

	var imageResponse dataResponse[entities.Image]
	decodeJSON(t, body, &imageResponse)

	if len(imageResponse.Data.Tags) != 2 {
		t.Fatalf("expected 2 tags on image, got %d", len(imageResponse.Data.Tags))
	}

	tagIDs := []uint64{imageResponse.Data.Tags[0].ID, imageResponse.Data.Tags[1].ID}
	sort.Slice(tagIDs, func(i, j int) bool { return tagIDs[i] < tagIDs[j] })
	if tagIDs[0] != firstTag.ID || tagIDs[1] != secondTag.ID {
		t.Fatalf("expected image tag ids [%d %d], got %v", firstTag.ID, secondTag.ID, tagIDs)
	}

	status, body = performJSONRequest(t, app, http.MethodDelete, fmt.Sprintf("/api/images/%d/tags/%d", targetImage.ID, firstTag.ID), nil)
	if status != http.StatusOK {
		t.Fatalf("expected detach status 200, got %d: %s", status, string(body))
	}

	status, body = performJSONRequest(t, app, http.MethodGet, fmt.Sprintf("/api/images/%d", targetImage.ID), nil)
	if status != http.StatusOK {
		t.Fatalf("expected image detail after detach status 200, got %d: %s", status, string(body))
	}

	decodeJSON(t, body, &imageResponse)
	if len(imageResponse.Data.Tags) != 1 || imageResponse.Data.Tags[0].ID != secondTag.ID {
		t.Fatalf("expected remaining tag id %d, got %#v", secondTag.ID, imageResponse.Data.Tags)
	}
}

type dataResponse[T any] struct {
	Data T `json:"data"`
}

type memoryStore struct {
	mu          sync.Mutex
	nextImageID uint64
	nextTagID   uint64
	images      map[uint64]entities.Image
	tags        map[uint64]entities.Tag
	assignments map[uint64]map[uint64]struct{}
}

type memoryImageRepo struct {
	store *memoryStore
}

type memoryTagRepo struct {
	store *memoryStore
}

type memoryImageTagRepo struct {
	store *memoryStore
}

func newTestApp() *fiber.App {
	store := &memoryStore{
		nextImageID: 1,
		nextTagID:   1,
		images:      make(map[uint64]entities.Image),
		tags:        make(map[uint64]entities.Tag),
		assignments: make(map[uint64]map[uint64]struct{}),
	}

	var (
		imageRepo    port.ImageRepo    = &memoryImageRepo{store: store}
		tagRepo      port.TagRepo      = &memoryTagRepo{store: store}
		imageTagRepo port.ImageTagRepo = &memoryImageTagRepo{store: store}
	)

	imageService := service.NewImageService(imageRepo, 12, 60)
	tagService := service.NewTagService(tagRepo)
	imageTagService := service.NewImageTagService(imageTagRepo, imageRepo, tagRepo)

	imageHandler := handler.NewImageHandler(imageService)
	tagHandler := handler.NewTagHandler(tagService)
	imageTagHandler := handler.NewImageTagHandler(imageTagService)

	app := fiber.New()
	api := app.Group("/api")
	BindImageRoutes(api.Group("/images"), imageHandler, imageTagHandler)
	BindTagRoutes(api.Group("/tags"), tagHandler)

	return app
}

func (r *memoryImageRepo) Create(_ context.Context, input entities.CreateImageInput) (entities.Image, error) {
	r.store.mu.Lock()
	defer r.store.mu.Unlock()

	now := time.Now().UTC()
	image := entities.Image{
		ID:           r.store.nextImageID,
		ImageURL:     input.ImageURL,
		ThumbnailURL: cloneStringPtr(input.ThumbnailURL),
		Width:        cloneIntPtr(input.Width),
		Height:       cloneIntPtr(input.Height),
		AltText:      cloneStringPtr(input.AltText),
		Source:       cloneStringPtr(input.Source),
		CreatedAt:    now,
		UpdatedAt:    now,
		Tags:         []entities.TagSummary{},
	}

	r.store.images[image.ID] = image
	r.store.nextImageID++

	return cloneImage(image), nil
}

func (r *memoryImageRepo) GetByID(_ context.Context, id uint64) (entities.Image, error) {
	r.store.mu.Lock()
	defer r.store.mu.Unlock()

	image, ok := r.store.images[id]
	if !ok || image.DeletedAt != nil {
		return entities.Image{}, fmt.Errorf("%w: image %d", apperrors.ErrNotFound, id)
	}

	return r.decorateImage(image), nil
}

func (r *memoryImageRepo) List(_ context.Context, filter entities.ImageListFilter) ([]entities.Image, error) {
	r.store.mu.Lock()
	defer r.store.mu.Unlock()

	images := make([]entities.Image, 0, len(r.store.images))
	for _, image := range r.store.images {
		if image.DeletedAt != nil {
			continue
		}
		if filter.Cursor != nil && image.ID >= *filter.Cursor {
			continue
		}
		if filter.TagSlug != "" && !r.imageHasTagSlug(image.ID, filter.TagSlug) {
			continue
		}
		images = append(images, r.decorateImage(image))
	}

	sort.Slice(images, func(i, j int) bool {
		return images[i].ID > images[j].ID
	})

	limit := filter.Limit
	if limit <= 0 || limit > len(images) {
		limit = len(images)
	}

	return cloneImages(images[:limit]), nil
}

func (r *memoryImageRepo) Update(_ context.Context, id uint64, input entities.UpdateImageInput) (entities.Image, error) {
	r.store.mu.Lock()
	defer r.store.mu.Unlock()

	image, ok := r.store.images[id]
	if !ok || image.DeletedAt != nil {
		return entities.Image{}, fmt.Errorf("%w: image %d", apperrors.ErrNotFound, id)
	}

	if input.ImageURL != nil {
		image.ImageURL = *input.ImageURL
	}
	if input.ThumbnailURL != nil {
		image.ThumbnailURL = cloneStringPtr(input.ThumbnailURL)
	}
	if input.Width != nil {
		image.Width = cloneIntPtr(input.Width)
	}
	if input.Height != nil {
		image.Height = cloneIntPtr(input.Height)
	}
	if input.AltText != nil {
		image.AltText = cloneStringPtr(input.AltText)
	}
	if input.Source != nil {
		image.Source = cloneStringPtr(input.Source)
	}
	image.UpdatedAt = time.Now().UTC()

	r.store.images[id] = image
	return r.decorateImage(image), nil
}

func (r *memoryImageRepo) SoftDelete(_ context.Context, id uint64) error {
	r.store.mu.Lock()
	defer r.store.mu.Unlock()

	image, ok := r.store.images[id]
	if !ok || image.DeletedAt != nil {
		return fmt.Errorf("%w: image %d", apperrors.ErrNotFound, id)
	}

	now := time.Now().UTC()
	image.DeletedAt = &now
	image.UpdatedAt = now
	r.store.images[id] = image

	return nil
}

func (r *memoryImageRepo) Exists(_ context.Context, id uint64) (bool, error) {
	r.store.mu.Lock()
	defer r.store.mu.Unlock()

	image, ok := r.store.images[id]
	return ok && image.DeletedAt == nil, nil
}

func (r *memoryImageRepo) imageHasTagSlug(imageID uint64, slug string) bool {
	assignments := r.store.assignments[imageID]
	for tagID := range assignments {
		tag, ok := r.store.tags[tagID]
		if !ok || tag.DeletedAt != nil {
			continue
		}
		if strings.EqualFold(tag.Slug, slug) {
			return true
		}
	}
	return false
}

func (r *memoryImageRepo) decorateImage(image entities.Image) entities.Image {
	cloned := cloneImage(image)
	tags := make([]entities.TagSummary, 0)
	for tagID := range r.store.assignments[image.ID] {
		tag, ok := r.store.tags[tagID]
		if !ok || tag.DeletedAt != nil {
			continue
		}
		tags = append(tags, entities.TagSummary{
			ID:   tag.ID,
			Name: tag.Name,
			Slug: tag.Slug,
		})
	}
	sort.Slice(tags, func(i, j int) bool {
		return tags[i].Name < tags[j].Name
	})
	cloned.Tags = tags
	return cloned
}

func (r *memoryTagRepo) Create(_ context.Context, input entities.CreateTagInput) (entities.Tag, error) {
	r.store.mu.Lock()
	defer r.store.mu.Unlock()

	for _, existing := range r.store.tags {
		if existing.DeletedAt != nil {
			continue
		}
		if strings.EqualFold(existing.Name, input.Name) || strings.EqualFold(existing.Slug, valueOrEmpty(input.Slug)) {
			return entities.Tag{}, apperrors.NewConflict("tag name or slug already exists")
		}
	}

	now := time.Now().UTC()
	tag := entities.Tag{
		ID:        r.store.nextTagID,
		Name:      input.Name,
		Slug:      valueOrEmpty(input.Slug),
		CreatedAt: now,
		UpdatedAt: now,
	}

	r.store.tags[tag.ID] = tag
	r.store.nextTagID++

	return cloneTag(tag), nil
}

func (r *memoryTagRepo) GetByID(_ context.Context, id uint64) (entities.Tag, error) {
	r.store.mu.Lock()
	defer r.store.mu.Unlock()

	tag, ok := r.store.tags[id]
	if !ok || tag.DeletedAt != nil {
		return entities.Tag{}, fmt.Errorf("%w: tag %d", apperrors.ErrNotFound, id)
	}

	return cloneTag(tag), nil
}

func (r *memoryTagRepo) List(_ context.Context) ([]entities.Tag, error) {
	r.store.mu.Lock()
	defer r.store.mu.Unlock()

	tags := make([]entities.Tag, 0, len(r.store.tags))
	for _, tag := range r.store.tags {
		if tag.DeletedAt != nil {
			continue
		}
		tags = append(tags, cloneTag(tag))
	}

	sort.Slice(tags, func(i, j int) bool {
		return tags[i].Name < tags[j].Name
	})

	return tags, nil
}

func (r *memoryTagRepo) Update(_ context.Context, id uint64, input entities.UpdateTagInput) (entities.Tag, error) {
	r.store.mu.Lock()
	defer r.store.mu.Unlock()

	tag, ok := r.store.tags[id]
	if !ok || tag.DeletedAt != nil {
		return entities.Tag{}, fmt.Errorf("%w: tag %d", apperrors.ErrNotFound, id)
	}

	nextName := tag.Name
	nextSlug := tag.Slug
	if input.Name != nil {
		nextName = *input.Name
	}
	if input.Slug != nil {
		nextSlug = *input.Slug
	}

	for _, existing := range r.store.tags {
		if existing.ID == id || existing.DeletedAt != nil {
			continue
		}
		if strings.EqualFold(existing.Name, nextName) || strings.EqualFold(existing.Slug, nextSlug) {
			return entities.Tag{}, apperrors.NewConflict("tag name or slug already exists")
		}
	}

	tag.Name = nextName
	tag.Slug = nextSlug
	tag.UpdatedAt = time.Now().UTC()
	r.store.tags[id] = tag

	return cloneTag(tag), nil
}

func (r *memoryTagRepo) SoftDelete(_ context.Context, id uint64) error {
	r.store.mu.Lock()
	defer r.store.mu.Unlock()

	tag, ok := r.store.tags[id]
	if !ok || tag.DeletedAt != nil {
		return fmt.Errorf("%w: tag %d", apperrors.ErrNotFound, id)
	}

	now := time.Now().UTC()
	tag.DeletedAt = &now
	tag.UpdatedAt = now
	r.store.tags[id] = tag

	return nil
}

func (r *memoryTagRepo) Exists(_ context.Context, id uint64) (bool, error) {
	r.store.mu.Lock()
	defer r.store.mu.Unlock()

	tag, ok := r.store.tags[id]
	return ok && tag.DeletedAt == nil, nil
}

func (r *memoryImageTagRepo) Attach(_ context.Context, imageID, tagID uint64) error {
	r.store.mu.Lock()
	defer r.store.mu.Unlock()

	if r.store.assignments[imageID] == nil {
		r.store.assignments[imageID] = make(map[uint64]struct{})
	}
	if _, exists := r.store.assignments[imageID][tagID]; exists {
		return apperrors.NewConflict("tag is already attached to image")
	}

	r.store.assignments[imageID][tagID] = struct{}{}
	return nil
}

func (r *memoryImageTagRepo) Detach(_ context.Context, imageID, tagID uint64) error {
	r.store.mu.Lock()
	defer r.store.mu.Unlock()

	assignments := r.store.assignments[imageID]
	if assignments == nil {
		return fmt.Errorf("%w: image tag assignment", apperrors.ErrNotFound)
	}
	if _, exists := assignments[tagID]; !exists {
		return fmt.Errorf("%w: image tag assignment", apperrors.ErrNotFound)
	}

	delete(assignments, tagID)
	return nil
}

func performJSONRequest(t *testing.T, app *fiber.App, method, target string, payload any) (int, []byte) {
	t.Helper()

	var bodyReader *bytes.Reader
	if payload == nil {
		bodyReader = bytes.NewReader(nil)
	} else {
		body, err := json.Marshal(payload)
		if err != nil {
			t.Fatalf("marshal request body: %v", err)
		}
		bodyReader = bytes.NewReader(body)
	}

	req := httptest.NewRequest(method, target, bodyReader)
	if payload != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := app.Test(req, -1)
	if err != nil {
		t.Fatalf("perform request: %v", err)
	}
	defer resp.Body.Close()

	responseBody := new(bytes.Buffer)
	if _, err := responseBody.ReadFrom(resp.Body); err != nil {
		t.Fatalf("read response body: %v", err)
	}

	return resp.StatusCode, responseBody.Bytes()
}

func createImageViaAPI(t *testing.T, app *fiber.App, payload entities.CreateImageInput) entities.Image {
	t.Helper()

	status, body := performJSONRequest(t, app, http.MethodPost, "/api/images", payload)
	if status != http.StatusCreated {
		t.Fatalf("expected create image status 201, got %d: %s", status, string(body))
	}

	var response dataResponse[entities.Image]
	decodeJSON(t, body, &response)
	return response.Data
}

func createTagViaAPI(t *testing.T, app *fiber.App, payload entities.CreateTagInput) entities.Tag {
	t.Helper()

	status, body := performJSONRequest(t, app, http.MethodPost, "/api/tags", payload)
	if status != http.StatusCreated {
		t.Fatalf("expected create tag status 201, got %d: %s", status, string(body))
	}

	var response dataResponse[entities.Tag]
	decodeJSON(t, body, &response)
	return response.Data
}

func decodeJSON(t *testing.T, body []byte, target any) {
	t.Helper()

	if err := json.Unmarshal(body, target); err != nil {
		t.Fatalf("decode json: %v\nbody: %s", err, string(body))
	}
}

func cloneImages(images []entities.Image) []entities.Image {
	cloned := make([]entities.Image, 0, len(images))
	for _, image := range images {
		cloned = append(cloned, cloneImage(image))
	}
	return cloned
}

func cloneImage(image entities.Image) entities.Image {
	cloned := image
	cloned.ThumbnailURL = cloneStringPtr(image.ThumbnailURL)
	cloned.Width = cloneIntPtr(image.Width)
	cloned.Height = cloneIntPtr(image.Height)
	cloned.AltText = cloneStringPtr(image.AltText)
	cloned.Source = cloneStringPtr(image.Source)
	if image.DeletedAt != nil {
		value := *image.DeletedAt
		cloned.DeletedAt = &value
	}
	if image.Tags != nil {
		cloned.Tags = append([]entities.TagSummary(nil), image.Tags...)
	}
	return cloned
}

func cloneTag(tag entities.Tag) entities.Tag {
	cloned := tag
	if tag.DeletedAt != nil {
		value := *tag.DeletedAt
		cloned.DeletedAt = &value
	}
	return cloned
}

func cloneStringPtr(value *string) *string {
	if value == nil {
		return nil
	}
	cloned := *value
	return &cloned
}

func cloneIntPtr(value *int) *int {
	if value == nil {
		return nil
	}
	cloned := *value
	return &cloned
}

func intPtr(value int) *int {
	return &value
}

func stringPtr(value string) *string {
	return &value
}

func valueOrEmpty(value *string) string {
	if value == nil {
		return ""
	}
	return *value
}
