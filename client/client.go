package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	pb "github.com/Thesohan/weaveGitHubSearchService/gen/go/protos/github/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	// Set up gRPC dial options with insecure credentials
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))

	// Establish a connection to the gRPC server
	conn, err := grpc.NewClient("localhost:8000", opts...)
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	defer conn.Close()

	// Create a new client for the GithubSearchService
	client := pb.NewGithubSearchServiceClient(conn)
	for {
		// Get user input for username and search term
		var user, searchTerm string
		fmt.Print("\nEnter username: ")
		fmt.Scanln(&user)
		var userPtr *string
		if user == "" {
			userPtr = nil // Pass nil if user is an empty string
		} else {
			userPtr = &user
		}
		fmt.Print("Enter search term: ")
		reader := bufio.NewReader(os.Stdin)
		searchTerm, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}
		fmt.Print(searchTerm)
		// Define the search parameters
		// URL decode the search term
		// Remove newline character from search term
		searchTerm = searchTerm[:len(searchTerm)-1]
		searchTerm = strings.ReplaceAll(searchTerm, " ", "+")
		fmt.Println(searchTerm)
		req := &pb.SearchRequest{
			SearchTerm: searchTerm,
			User:       userPtr, // Set to a username if needed
		}

		// Perform the search request
		resp, err := client.Search(context.Background(), req)
		if err != nil {
			fmt.Printf("Search failed: %v", err)
			continue
		}

		// Print the search results
		for _, result := range resp.Results {
			fmt.Printf("File: %s, Repo: %s\n", result.FileUrl, result.Repo)
		}
	}

}
