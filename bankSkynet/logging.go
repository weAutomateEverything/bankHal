package bankSkynet

import (
	"context"
	"github.com/go-kit/kit/log"
	"time"
)

type loggingService struct {
	logger log.Logger
	Service
}

func NewLoggingService(logger log.Logger, s Service) Service {
	return &loggingService{logger, s}
}

func (s loggingService) RecreateNode(ctx context.Context,chatid uint32, nodeName, callerName string) error {
	defer func(begin time.Time) {
		s.logger.Log(
			"method", "sendDeliveryAlert",
			"nodeName", nodeName,
			"callerName", callerName,
			"chat", chatid,
			"took", time.Since(begin),
		)
	}(time.Now())
	return s.Service.RecreateNode(ctx,chatid, nodeName, callerName)
}
