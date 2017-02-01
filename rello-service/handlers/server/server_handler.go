package handler

import (
	"fmt"
	"golang.org/x/net/context"
	"os"
	"time"

	"github.com/pkg/errors"
	"google.golang.org/grpc"

	ambitionSVC "github.com/adamryman/ambition-model/ambition-service"
	pb "github.com/adamryman/ambition-rello/rello-service"
	usersSVC "github.com/adamryman/ambition-users/users-service"

	ambitionClient "github.com/adamryman/ambition-model/ambition-service/generated/client/grpc"
	usersClient "github.com/adamryman/ambition-users/users-service/generated/client/grpc"

	"github.com/davecgh/go-spew/spew"
)

const trelloTime = "2006-01-02T15:04:05.999Z07:00"

// NewService returns a naïve, stateless implementation of Service.
func NewService() pb.RelloServer {
	// TODO: Create ambitionSVC and usersSVC clients
	// TODO switch to grpc clients
	modelConnStr := os.Getenv("MODEL_HOST") + ":" + os.Getenv("MODEL_PORT")
	fmt.Println(modelConnStr)
	modelConn, err := grpc.Dial(modelConnStr, grpc.WithInsecure(), grpc.WithTimeout(time.Second))
	m, err := ambitionClient.New(modelConn)
	if err != nil {
		panic(err)
	}
	usersConnStr := os.Getenv("USERS_HOST") + ":" + os.Getenv("USERS_PORT")
	fmt.Println(usersConnStr)
	usersConn, err := grpc.Dial(usersConnStr, grpc.WithInsecure(), grpc.WithTimeout(time.Second))
	u, err := usersClient.New(usersConn)
	if err != nil {
		panic(err)
	}

	return relloService{
		model: m,
		users: u,
	}
}

type relloService struct {
	users usersSVC.UsersServer
	model ambitionSVC.AmbitionServer
}

// CheckListWebhook implements Service.
// TODO: Auth middleware using user trello webhookcallback
func (s relloService) CheckListWebhook(ctx context.Context, in *pb.ChecklistUpdate) (*pb.Empty, error) {
	fmt.Println(spew.Sdump(in))
	fmt.Println(spew.Sdump(in.GetAction()))
	fmt.Println(spew.Sdump(in.GetAction().GetMemberCreator))

	user, err := s.users.ReadUser(ctx,
		&usersSVC.User{
			Trello: &usersSVC.TrelloInfo{
				ID: in.GetAction().GetIdMemberCreator(),
			},
			Info: &usersSVC.UserInfo{},
		})
	if err != nil {
		//TODO: Wrap error
		return nil, err
	}

	cItem := in.Action.Data.CheckItem
	switch in.Action.Type {
	case "createCheckItem":
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
		action, err := s.model.ReadAction(ctx,
			&ambitionSVC.Action{
				Name:   cItem.GetName(),
				UserID: user.GetID(),
			})
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
				UserID: user.GetID(),
				Occurrence: &ambitionSVC.Occurrence{
					ActionID: action.GetID(),
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
