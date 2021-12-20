package main

import (
	"fmt"

	"golang.org/x/net/context"

	pb "github.com/brotherlogic/display/proto"
)

func (s *Server) Show(ctx context.Context, req *pb.ShowRequest) (*pb.ShowResponse, error) {
	return nil, fmt.Errorf("Not implemented")
}
