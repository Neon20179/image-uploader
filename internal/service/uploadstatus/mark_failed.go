package uploadstatus

import (
	"context"

	"github.com/google/uuid"
)

func (s *UploadStatusService) MarkFailed(ctx context.Context, id uuid.UUID) error {
	return s.uploadRepository.MarkFailed(ctx, id)
}
