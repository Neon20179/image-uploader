package image

import (
	"context"
	"fmt"
)

func (r *ImageRepository) GetImagesCount(ctx context.Context) (int, error) {
	querier := r.querierForContext(ctx)
	ctx, cancel := context.WithTimeout(ctx, querier.OpTimeout())
	defer cancel()

	query := `
		SELECT COUNT(*)
		FROM images;
	`

	row := querier.QueryRow(ctx, query)

	var count int

	err := row.Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("scan error: %w", err)
	}

	return count, nil
}
