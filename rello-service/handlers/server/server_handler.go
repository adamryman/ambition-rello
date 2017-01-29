package handler

import (
	"golang.org/x/net/context"

	pb "github.com/adamryman/ambition-rello/rello-service"
)

// NewService returns a naïve, stateless implementation of Service.
func NewService() pb.RelloServer {
	return relloService{}
}

type relloService struct{}

// CheckListWebhook implements Service.
func (s relloService) CheckListWebhook(ctx context.Context, in *pb.ChecklistUpdate) (*pb.Empty, error) {
	var resp pb.Empty
	resp = pb.Empty{}
	return &resp, nil
}
