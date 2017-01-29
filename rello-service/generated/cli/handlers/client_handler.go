package clienthandler

import (
	pb "github.com/adamryman/ambition-rello/rello-service"
)

// CheckListWebhook implements Service.
func CheckListWebhook(ModelCheckListWebhook pb.Model, ActionCheckListWebhook pb.Action) (*pb.ChecklistUpdate, error) {
	request := pb.ChecklistUpdate{
		Model:  &ModelCheckListWebhook,
		Action: &ActionCheckListWebhook,
	}
	return &request, nil
}
