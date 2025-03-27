package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"text/tabwriter"
	"time"

	pb "github.com/Thesohan/weaveGitHubSearchService/gen/go/protos/github/v1"
	"github.com/Thesohan/weaveGitHubSearchService/server/github"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	serverAddress string = ":8000"
	tcpNetwork           = "tcp"
)

// gitHubSearchService implements the gRPC service.
type gitHubSearchService struct {
	pb.UnimplementedGithubSearchServiceServer
	ghClient github.IGithubClient
}

// Search handles search requests.
func (s *gitHubSearchService) Search(ctx context.Context, searchReq *pb.SearchRequest) (*pb.SearchResponse, error) {
	// Validate the request payload
	if searchReq.Term == "" {
		return nil, fmt.Errorf("search term is required")
	}

	query := buildGitHubQuery(searchReq)
	data, err := s.ghClient.SearchCode(ctx, query)
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
func buildGitHubQuery(searchReq *pb.SearchRequest) string {
	query := searchReq.Term
	if searchReq.GetUser() != "" {
		query += " user:" + searchReq.GetUser()
	}
	return query
}
func runServer() {
	listener, err := net.Listen(tcpNetwork, serverAddress)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	githubClient := github.NewGitHubClient()
	pb.RegisterGithubSearchServiceServer(grpcServer, &gitHubSearchService{ghClient: githubClient})

	log.Printf("gRPC server running on port %v...", serverAddress)
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}

}
func runClient() {
	// Establish a connection to the gRPC server with insecure credentials
	conn, err := grpc.NewClient("localhost:8000", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	defer conn.Close()

	// Create a new client for the GithubSearchService
	client := pb.NewGithubSearchServiceClient(conn)
	for {
		// Get user input for username and search term
		log.Print("\nEnter username: ")
		var userInput string
		fmt.Scanln(&userInput) // fmt.Scanln` stops after reading the first whitespace, user can't have any whitspace in github
		var user *string
		if userInput != "" {
			user = &userInput
		}
		log.Print("Enter search term: ")
		reader := bufio.NewReader(os.Stdin)
		searchTerm, err := reader.ReadString('\n') // `bufio.NewReader` stops after the given delimiter `\n`. searchTerm can be a complete sentence.
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("Searching for %v", searchTerm)
		// Remove newline character from search term
		searchTerm = strings.TrimSuffix(searchTerm, "\n")
		log.Println(searchTerm)
		req := &pb.SearchRequest{
			Term: searchTerm,
			User: user, // Set to a username if needed
		}

		// Perform the search request
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		resp, err := client.Search(ctx, req)
		if err != nil {
			log.Printf("Search failed: %v", err)
			continue
		}

		// Create a new tabwriter
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', tabwriter.TabIndent)
		// Print header
		fmt.Fprintln(w, "S.N\tFile\tRepo")
		for index, result := range resp.Results {
			fmt.Fprintf(w, "%d\t%s\t%s\n", index, result.FileUrl, result.Repo)
		}
		// Flush the writer to ensure output is displayed
		if err := w.Flush(); err != nil {
			log.Printf("Failed to flush writer: %v", err)
		}
	}

}
func main() {
	go runServer()
	time.Sleep(5 * time.Second) // Waiting for 5 second for server to start
	runClient()
}
