package main

import (
	"fmt"

	"golang.org/x/net/context"

	pb "github.com/brotherlogic/display/proto"
	rcpb "github.com/brotherlogic/recordcollection/proto"
)

func (s *Server) Show(ctx context.Context, req *pb.ShowRequest) (*pb.ShowResponse, error) {
	return nil, fmt.Errorf("Not implemented")
}

//ClientUpdate on an updated record
func (s *Server) ClientUpdate(ctx context.Context, req *rcpb.ClientUpdateRequest) (*rcpb.ClientUpdateResponse, error) {
	s.buildPage(ctx)
	return &rcpb.ClientUpdateResponse{}, nil
}
