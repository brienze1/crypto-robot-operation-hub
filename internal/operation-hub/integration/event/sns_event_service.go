package event

import "github.com/brienze1/crypto-robot-operation-hub/internal/operation-hub/domain/model"

type snsEventService struct {
}

func SNSEventService() *snsEventService {
	return &snsEventService{}
}

func (s *snsEventService) Send(model.Client) error {
	return nil
}
