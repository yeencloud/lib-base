package events

import (
	"context"

	"github.com/yeencloud/lib-events/domain"
)

type Publisher interface {
	Publish(ctx context.Context, message domain.PublishableMessage) error
}
