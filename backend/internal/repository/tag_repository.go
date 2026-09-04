package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"

	"floway-backend/internal/model"
)

// TagRepository backs both independent tag types (product / client-type) —
// the two never share a table (see migration 00032), so one implementation
// is parameterized by table/join-table names fixed at construction time
// rather than duplicated twice. table/joinTable/joinColumn are always
// compile-time constants (never user input), so building SQL strings from
// them here is safe.
type TagRepository struct {
	db         *pgxpool.Pool
	table      string
	joinTable  string
	joinColumn string
}

func NewProductTagRepository(db *pgxpool.Pool) *TagRepository {
	return &TagRepository{db: db, table: "product_tags", joinTable: "client_product_tags", joinColumn: "product_tag_id"}
}

func NewClientTypeTagRepository(db *pgxpool.Pool) *TagRepository {
	return &TagRepository{db: db, table: "client_type_tags", joinTable: "client_client_type_tags", joinColumn: "client_type_tag_id"}
}

// Search returns tags whose name contains query (case-insensitive),
// capped at 20 — backs both the filter dropdown (blank query) and the
// client-detail autocomplete (as-you-type query).
func (r *TagRepository) Search(ctx context.Context, query string) ([]model.Tag, error) {
	rows, err := r.db.Query(ctx, fmt.Sprintf(`
		SELECT id, name FROM %s WHERE name ILIKE '%%' || $1 || '%%' ORDER BY name LIMIT 20
	`, r.table), query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := []model.Tag{}
	for rows.Next() {
		var item model.Tag
		if err := rows.Scan(&item.ID, &item.Name); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

// FindOrCreateByName matches case-insensitively so "Свадьба" and "свадьба"
// resolve to the same row — the backend half of "autocomplete, create a
// new tag on the fly" (no separate tag-management screen).
func (r *TagRepository) FindOrCreateByName(ctx context.Context, name string) (model.Tag, error) {
	var item model.Tag
	err := r.db.QueryRow(ctx, fmt.Sprintf(`SELECT id, name FROM %s WHERE lower(name) = lower($1)`, r.table), name).Scan(&item.ID, &item.Name)
	if err == nil {
		return item, nil
	}

	// ON CONFLICT DO UPDATE (rather than DO NOTHING) so this still RETURNINGs
	// a row if a concurrent insert won the race between the SELECT above and
	// this INSERT — vanishingly unlikely in a single-admin system, but free
	// to handle correctly here.
	err = r.db.QueryRow(ctx, fmt.Sprintf(`
		INSERT INTO %s (name) VALUES ($1)
		ON CONFLICT (lower(name)) DO UPDATE SET name = EXCLUDED.name
		RETURNING id, name
	`, r.table), name).Scan(&item.ID, &item.Name)
	return item, err
}

// SetForClient replaces the full set of this tag type assigned to a
// client — matches the "combobox holds the current selection" UX the
// frontend uses, so callers never need incremental add/remove endpoints.
func (r *TagRepository) SetForClient(ctx context.Context, clientID int64, tagIDs []int64) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	if _, err := tx.Exec(ctx, fmt.Sprintf(`DELETE FROM %s WHERE client_id = $1`, r.joinTable), clientID); err != nil {
		return err
	}
	for _, tagID := range tagIDs {
		if _, err := tx.Exec(ctx, fmt.Sprintf(`INSERT INTO %s (client_id, %s) VALUES ($1, $2)`, r.joinTable, r.joinColumn), clientID, tagID); err != nil {
			return err
		}
	}
	return tx.Commit(ctx)
}

func (r *TagRepository) ListForClient(ctx context.Context, clientID int64) ([]model.Tag, error) {
	rows, err := r.db.Query(ctx, fmt.Sprintf(`
		SELECT t.id, t.name
		FROM %s j JOIN %s t ON t.id = j.%s
		WHERE j.client_id = $1
		ORDER BY t.name
	`, r.joinTable, r.table, r.joinColumn), clientID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := []model.Tag{}
	for rows.Next() {
		var item model.Tag
		if err := rows.Scan(&item.ID, &item.Name); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}
