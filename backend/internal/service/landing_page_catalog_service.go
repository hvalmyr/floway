package service

import (
	"context"

	"floway-backend/internal/model"
)

// Narrow interfaces scoped to exactly what LandingPageCatalogService needs —
// this service doesn't own CRUD for any of them (that's LandingPageService/
// LandingPageBlockService/LandingPageFAQService's job), it only aggregates
// reads for the two public decor-site endpoints. Mirrors
// CourseCatalogService's own reasoning.
type LandingPageListRepository interface {
	ListVisibleWithObjectType(ctx context.Context) ([]model.LandingPageWithObjectType, error)
}

type LandingPageLookupRepository interface {
	FindBySlug(ctx context.Context, slug string) (model.LandingPage, error)
}

type ObjectTypeLookupRepository interface {
	List(ctx context.Context) ([]model.ObjectType, error)
}

type LandingPageBlockListRepository interface {
	ListByLandingPageID(ctx context.Context, landingPageID int64) ([]model.LandingPageBlock, error)
}

type LandingPageFAQListRepository interface {
	ListByLandingPageID(ctx context.Context, landingPageID int64) ([]model.LandingPageFAQItem, error)
}

// LandingPageCatalogService assembles the decor site's two public "read the
// whole tree in one shot" responses: the homepage/nav listing (visible
// pages + their object type) and a single landing page (page -> object type
// + blocks + FAQ items).
type LandingPageCatalogService struct {
	list       LandingPageListRepository
	lookup     LandingPageLookupRepository
	objectType ObjectTypeLookupRepository
	blocks     LandingPageBlockListRepository
	faqItems   LandingPageFAQListRepository
}

func NewLandingPageCatalogService(
	list LandingPageListRepository,
	lookup LandingPageLookupRepository,
	objectType ObjectTypeLookupRepository,
	blocks LandingPageBlockListRepository,
	faqItems LandingPageFAQListRepository,
) *LandingPageCatalogService {
	return &LandingPageCatalogService{list: list, lookup: lookup, objectType: objectType, blocks: blocks, faqItems: faqItems}
}

// ListVisible is the public homepage/nav aggregation: every visible landing
// page with its object type, one query.
func (s *LandingPageCatalogService) ListVisible(ctx context.Context) ([]model.LandingPageWithObjectType, error) {
	return s.list.ListVisibleWithObjectType(ctx)
}

// GetFullBySlug is the public landing-page aggregation: the page, its
// object type, its blocks, and its FAQ items. A hidden page is reported as
// not found (same convention as CourseCatalogService.GetFullBySlug — "hide"
// means unpublish, not "leave off the homepage but still directly
// reachable").
func (s *LandingPageCatalogService) GetFullBySlug(ctx context.Context, slug string) (model.LandingPageWithDetail, error) {
	page, err := s.lookup.FindBySlug(ctx, slug)
	if err != nil {
		return model.LandingPageWithDetail{}, err
	}
	if !page.Visible {
		return model.LandingPageWithDetail{}, ErrNotFound
	}

	// object_types is a tiny dictionary (a handful of rows) — listing it all
	// and picking the match avoids a second lookup-by-id repository method
	// just for this one call site.
	objectTypes, err := s.objectType.List(ctx)
	if err != nil {
		return model.LandingPageWithDetail{}, err
	}
	var objectType model.ObjectType
	for _, ot := range objectTypes {
		if ot.ID == page.ObjectTypeID {
			objectType = ot
			break
		}
	}

	blocks, err := s.blocks.ListByLandingPageID(ctx, page.ID)
	if err != nil {
		return model.LandingPageWithDetail{}, err
	}

	faqItems, err := s.faqItems.ListByLandingPageID(ctx, page.ID)
	if err != nil {
		return model.LandingPageWithDetail{}, err
	}

	return model.LandingPageWithDetail{
		LandingPage: page,
		ObjectType:  objectType,
		Blocks:      blocks,
		FAQItems:    faqItems,
	}, nil
}
