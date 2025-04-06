package service

import (
	"github.com/yeencloud/lib-base/domain/errors"
	events "github.com/yeencloud/lib-events"
)

func (bs *BaseService) GetMqSubscriber() (*events.Subscriber, error) {
	if bs.mqSubscriber == nil {
		return nil, &errors.ModuleNotInitializedError{Module: "redis-subscriber"}
	}
	return bs.mqSubscriber, nil
}

func (bs *BaseService) GetMqPublisher() (*events.Publisher, error) {
	if bs.mqPublisher == nil {
		return nil, &errors.ModuleNotInitializedError{Module: "redis-publisher"}
	}
	return bs.mqPublisher, nil
}
