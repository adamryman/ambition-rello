package middlewares

import (
	pb "github.com/adamryman/ambition-rello/rello-service"
)

func WrapService(in pb.RelloServer) pb.RelloServer {
	return in
}
