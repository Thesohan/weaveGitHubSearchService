package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/Thesohan/weaveGitHubSearchService/server/constants"
	github "github.com/Thesohan/weaveGitHubSearchService/server/github"

	pb "github.com/Thesohan/weaveGitHubSearchService/gen/go/protos/github/v1"
	"google.golang.org/grpc"
)

// gitHubSearchService implements the gRPC service.
type gitHubSearchService struct {
	pb.UnimplementedGithubSearchServiceServer
}

// Search handles search requests.
func (s *gitHubSearchService) Search(ctx context.Context, req *pb.SearchRequest) (*pb.SearchResponse, error) {
	// Validate the request payload
	if req.SearchTerm == "" {
		return nil, fmt.Errorf("search term is required")
	}

	query := buildGitHubQuery(req)
	newGitHubClient := github.NewGitHubClient()
	data, err := newGitHubClient.SearchCode(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("error while searching code %v", err)
	}
	var results []*pb.Result
	for _, item := range data.Items {
		results = append(results, &pb.Result{
			FileUrl: item.HTMLURL,
			Repo:    item.Repository.FullName,
		})
	}
	return &pb.SearchResponse{Results: results}, nil
}

// buildGitHubQuery constructs the search query.
func buildGitHubQuery(req *pb.SearchRequest) string {
	query := req.SearchTerm
	if req.User != nil {
		query += "+user:" + *req.User
	}
	return query
}

func main() {
	listener, err := net.Listen(constants.NETWORK_TCP, constants.SERVER_ADDRESS)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterGithubSearchServiceServer(grpcServer, &gitHubSearchService{})

	log.Printf("gRPC server running on port %v...", constants.SERVER_ADDRESS)
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
