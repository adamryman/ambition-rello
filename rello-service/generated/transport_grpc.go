package svc

// This file provides server-side bindings for the gRPC transport.
// It utilizes the transport/grpc.Server.

import (
	"net/http"

	"golang.org/x/net/context"
	"google.golang.org/grpc/metadata"

	grpctransport "github.com/go-kit/kit/transport/grpc"

	// This Service
	pb "github.com/adamryman/ambition-rello/rello-service"
)

// MakeGRPCServer makes a set of endpoints available as a gRPC RelloServer.
func MakeGRPCServer(ctx context.Context, endpoints Endpoints) pb.RelloServer {
	serverOptions := []grpctransport.ServerOption{
		grpctransport.ServerBefore(metadataToContext),
	}
	return &grpcServer{
		// rello

		checklistwebhook: grpctransport.NewServer(
			ctx,
			endpoints.CheckListWebhookEndpoint,
			DecodeGRPCCheckListWebhookRequest,
			EncodeGRPCCheckListWebhookResponse,
			serverOptions...,
		),
	}
}

// grpcServer implements the RelloServer interface
type grpcServer struct {
	checklistwebhook grpctransport.Handler
}

// Methods for grpcServer to implement RelloServer interface

func (s *grpcServer) CheckListWebhook(ctx context.Context, req *pb.ChecklistUpdate) (*pb.Empty, error) {
	_, rep, err := s.checklistwebhook.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return rep.(*pb.Empty), nil
}

// Server Decode

// DecodeGRPCCheckListWebhookRequest is a transport/grpc.DecodeRequestFunc that converts a
// gRPC checklistwebhook request to a user-domain checklistwebhook request. Primarily useful in a server.
func DecodeGRPCCheckListWebhookRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*pb.ChecklistUpdate)
	return req, nil
}

// Client Decode

// DecodeGRPCCheckListWebhookResponse is a transport/grpc.DecodeResponseFunc that converts a
// gRPC checklistwebhook reply to a user-domain checklistwebhook response. Primarily useful in a client.
func DecodeGRPCCheckListWebhookResponse(_ context.Context, grpcReply interface{}) (interface{}, error) {
	reply := grpcReply.(*pb.Empty)
	return reply, nil
}

// Server Encode

// EncodeGRPCCheckListWebhookResponse is a transport/grpc.EncodeResponseFunc that converts a
// user-domain checklistwebhook response to a gRPC checklistwebhook reply. Primarily useful in a server.
func EncodeGRPCCheckListWebhookResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(*pb.Empty)
	return resp, nil
}

// Client Encode

// EncodeGRPCCheckListWebhookRequest is a transport/grpc.EncodeRequestFunc that converts a
// user-domain checklistwebhook request to a gRPC checklistwebhook request. Primarily useful in a client.
func EncodeGRPCCheckListWebhookRequest(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(*pb.ChecklistUpdate)
	return req, nil
}

// Helpers

func metadataToContext(ctx context.Context, md *metadata.MD) context.Context {
	for k, v := range *md {
		if v != nil {
			// The key is added both in metadata format (k) which is all lower
			// and the http.CanonicalHeaderKey of the key so that it can be
			// accessed in either format
			ctx = context.WithValue(ctx, k, v[0])
			ctx = context.WithValue(ctx, http.CanonicalHeaderKey(k), v[0])
		}
	}

	return ctx
}
