package handler

import (
	"fmt"
	"golang.org/x/net/context"
	"time"

	"github.com/pkg/errors"

	ambitionSVC "github.com/adamryman/ambition-model/ambition-service"
	pb "github.com/adamryman/ambition-rello/rello-service"
	usersSVC "github.com/adamryman/ambition-users/users-service"
)

const trelloTime = "2006-01-02T15:04:05.999Z07:00"

// NewService returns a naïve, stateless implementation of Service.
func NewService() pb.RelloServer {
	// TODO: Create ambitionSVC and usersSVC clients
	return relloService{}
}

type relloService struct {
	users usersSVC.UsersServer
	model ambitionSVC.AmbitionServer
}

// CheckListWebhook implements Service.
// TODO: Auth middleware using user trello webhookcallback
func (s relloService) CheckListWebhook(ctx context.Context, in *pb.ChecklistUpdate) (*pb.Empty, error) {
	user, err := s.users.ReadUser(ctx,
		&usersSVC.User{
			Trello: &usersSVC.TrelloInfo{
				ID: in.Action.MemberCreator.Id,
			},
		})
	if err != nil {
		//TODO: Wrap error
		return nil, err
	}

	cItem := in.Action.Data.CheckItem
	switch in.Action.Type {
	case "createCheckItem":
		// TODO: CreateAction
		// TODO: Map CheckItem.Id to Action.Id
		// TODO: Map u.Action.MemberCreator.Id to Action.UserId
		// UserId hardcoded to 1
		_ = in.Action.MemberCreator.Id

		action, err := s.model.CreateAction(ctx,
			&ambitionSVC.Action{
				Name:   cItem.GetName(),
				UserID: user.GetID(),
			})
		if err != nil {
			fmt.Println(errors.Wrap(err, "cannot create action"))
			break
		}
		fmt.Printf("%s action created\n", action.Name)
	case "updateCheckItemStateOnCard":
		if cItem.State == "incomplete" {
			fmt.Printf("%q is being unchecked\n", cItem.Name)
			break
		}
		dateString := &in.Action.Date
		date, err := time.Parse(trelloTime, *dateString)
		_ = date
		if err != nil {
			fmt.Println(errors.Wrap(err, "time parsing failed"))
			break
		}
		// TODO: log occurrence
		_, err = s.model.CreateOccurrence(ctx,
			&ambitionSVC.CreateOccurrenceRequest{
				Occurrence: &ambitionSVC.Occurrence{
					//ActionID:
					// TODO: make sure this time format is fine
					Datetime: date.String(),
				},
			})
		if err != nil {
			fmt.Println(errors.Wrap(err, "error creating occurrence"))
			break
		}

		// TODO: Create Occurrence with the ActionId from CheckItem.Id
		// Maps to action id
	default:
		fmt.Println(in.Action.Type)
	}
	return nil, nil
}
