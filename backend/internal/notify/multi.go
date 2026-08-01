package notify

import (
	"context"
	"errors"

	"floway-backend/internal/model"
)

type Channel interface {
	NotifyNewLead(ctx context.Context, lead model.Lead) error
}

// MultiNotifier fans a lead out to every configured channel independently —
// one channel failing (e.g. SMTP relay down) must not stop the others from
// being tried.
type MultiNotifier struct {
	channels []Channel
}

func NewMultiNotifier(channels ...Channel) *MultiNotifier {
	return &MultiNotifier{channels: channels}
}

func (m *MultiNotifier) NotifyNewLead(ctx context.Context, lead model.Lead) error {
	var errs []error
	for _, ch := range m.channels {
		if err := ch.NotifyNewLead(ctx, lead); err != nil {
			errs = append(errs, err)
		}
	}
	return errors.Join(errs...)
}
