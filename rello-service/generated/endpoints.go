package svc

// This file contains methods to make individual endpoints from services,
// request and response types to serve those endpoints, as well as encoders and
// decoders for those types, for all of our supported transport serialization
// formats. It also includes endpoint middlewares.

import (
	"golang.org/x/net/context"

	"github.com/go-kit/kit/endpoint"

	pb "github.com/adamryman/ambition-rello/rello-service"
)

// Endpoints collects all of the endpoints that compose an add service. It's
// meant to be used as a helper struct, to collect all of the endpoints into a
// single parameter.
//
// In a server, it's useful for functions that need to operate on a per-endpoint
// basis. For example, you might pass an Endpoints to a function that produces
// an http.Handler, with each method (endpoint) wired up to a specific path. (It
// is probably a mistake in design to invoke the Service methods on the
// Endpoints struct in a server.)
//
// In a client, it's useful to collect individually constructed endpoints into a
// single type that implements the Service interface. For example, you might
// construct individual endpoints using transport/http.NewClient, combine them
// into an Endpoints, and return it to the caller as a Service.
type Endpoints struct {
	CheckListWebhookEndpoint endpoint.Endpoint
}

// Endpoints

func (e Endpoints) CheckListWebhook(ctx context.Context, in *pb.ChecklistUpdate) (*pb.Empty, error) {
	response, err := e.CheckListWebhookEndpoint(ctx, in)
	if err != nil {
		return nil, err
	}
	return response.(*pb.Empty), nil
}

// Make Endpoints

func MakeCheckListWebhookEndpoint(s pb.RelloServer) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(*pb.ChecklistUpdate)
		v, err := s.CheckListWebhook(ctx, req)
		if err != nil {
			return nil, err
		}
		return v, nil
	}
}

// WrapAll wraps each Endpoint field of struct Endpoints with a
// go-kit/kit/endpoint.Middleware.
// Use this for applying a set of middlewares to every endpoint in the service.
// The middlewares will be applied in the order passed, with the first
// middleware being the outermost middleware.
func (e *Endpoints) WrapAll(middlewares ...endpoint.Middleware) {
	if len(middlewares) == 0 {
		return
	}
	m := endpoint.Chain(middlewares[0], middlewares[1:]...)

	e.CheckListWebhookEndpoint = m(e.CheckListWebhookEndpoint)
}
