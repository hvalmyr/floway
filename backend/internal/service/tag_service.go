package service

import (
	"context"
	"errors"
	"regexp"

	"floway-backend/internal/model"
)

// hexColor matches a strict "#rrggbb" hex string — the only shape the
// frontend's <input type="color"> ever produces, so this is a sanity check
// against malformed input rather than a real color-format parser.
var hexColor = regexp.MustCompile(`^#[0-9a-fA-F]{6}$`)

type TagSearchRepository interface {
	Search(ctx context.Context, query string) ([]model.Tag, error)
	Delete(ctx context.Context, id int64) error
	UpdateColor(ctx context.Context, id int64, color string) (model.Tag, error)
}

// TagService fronts both independent tag tables for the autocomplete/filter
// endpoint (GET /tags?type=product|client_type) — picking the right
// repository by TagType rather than having two near-identical services.
type TagService struct {
	productTags TagSearchRepository
	typeTags    TagSearchRepository
}

func NewTagService(productTags, typeTags TagSearchRepository) *TagService {
	return &TagService{productTags: productTags, typeTags: typeTags}
}

func (s *TagService) Search(ctx context.Context, tagType model.TagType, query string) ([]model.Tag, error) {
	switch tagType {
	case model.TagTypeProduct:
		return s.productTags.Search(ctx, query)
	case model.TagTypeClientType:
		return s.typeTags.Search(ctx, query)
	default:
		return nil, errors.Join(ErrValidation, errors.New("invalid tag type"))
	}
}

// Delete removes the tag definition itself, not just one client's
// assignment — see TagRepository.Delete for the cascade behavior.
func (s *TagService) Delete(ctx context.Context, tagType model.TagType, id int64) error {
	if id == 0 {
		return errors.Join(ErrValidation, errors.New("id is required"))
	}
	switch tagType {
	case model.TagTypeProduct:
		return s.productTags.Delete(ctx, id)
	case model.TagTypeClientType:
		return s.typeTags.Delete(ctx, id)
	default:
		return errors.Join(ErrValidation, errors.New("invalid tag type"))
	}
}

// SetColor sets a tag's background color, e.g. "#f3d9c4" — the color lives
// on the tag definition, so it changes everywhere the tag is shown at once.
func (s *TagService) SetColor(ctx context.Context, tagType model.TagType, id int64, color string) (model.Tag, error) {
	if id == 0 {
		return model.Tag{}, errors.Join(ErrValidation, errors.New("id is required"))
	}
	if !hexColor.MatchString(color) {
		return model.Tag{}, errors.Join(ErrValidation, errors.New("color must be a #rrggbb hex string"))
	}
	switch tagType {
	case model.TagTypeProduct:
		return s.productTags.UpdateColor(ctx, id, color)
	case model.TagTypeClientType:
		return s.typeTags.UpdateColor(ctx, id, color)
	default:
		return model.Tag{}, errors.Join(ErrValidation, errors.New("invalid tag type"))
	}
}
