package handler

import (
	"github.com/brienze1/crypto-robot-operation-hub/internal/operation-hub/delivery/handler"
	"github.com/golang/mock/gomock"
	"testing"
)

func TestHandler(t *testing.T) {
	ctrl := gomock.NewController(t)

	h := NewMockHandler(ctrl)
	v := NewMockT(ctrl)

	handler.Handler(h, v)
}
