package service

import (
	"context"
	"errors"

	"floway-backend/internal/model"
)

type TagSearchRepository interface {
	Search(ctx context.Context, query string) ([]model.Tag, error)
	Delete(ctx context.Context, id int64) error
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
