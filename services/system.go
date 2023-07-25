package services

import (
	"context"
	"github.com/d3v-friends/ms-accounts-grpc/models"
	"time"
)

type SystemImpl struct {
}

func (x *SystemImpl) ReadAccountIndex(ctx context.Context, _ *Empty) (*AccountIndex, error) {
	system, err := models.FindSystem(ctx)
	if err != nil {
		return nil, err
	}

	return &AccountIndex{
		Identifier: system.Data.Identifier,
		Property:   system.Data.Property,
		Permission: system.Data.Permission,
		UpdatedAt:  system.Data.CreatedAt.Format(time.RFC3339),
	}, nil
}

func (x *SystemImpl) UpdateKeys(ctx context.Context, i *IUpdateKeys) (*Empty, error) {
	_, err := models.UpdateSystem(ctx, i)
	if err != nil {
		return nil, err
	}
	return &Empty{}, nil
}

func (x *SystemImpl) mustEmbedUnimplementedSystemsServer() {
}
