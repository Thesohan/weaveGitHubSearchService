### Github Search Service
Build a service using gRPC that builds on top of the GitHub API to perform queries for the
provided search phrase and allows for optional filtering down to the user level. You’ll return the
file URL and the repo it was found in.
This is the API you’ll use to perform the GitHub search:
        ```https://docs.github.com/en/rest/reference/search```

This is the API spec that should be implemented:

    service GithubSearchService {
        rpc Search(SearchRequest) returns (SearchResponse);
    }
    
    message SearchRequest {
            string search_term = 1;
            string user = 2;
        }
    
    message SearchResponse {
            repeated Result results = 1;
        }
    
    message Result {
            string file_url = 1;
            string repo = 2;
        }
        
Instructions
1. Create a new, public GitHub repository with only a README
2. Create a new branch and do all your work in that branch
3. Create a PR back into the main branch and send in the URL to the Pull Request


## Prerequisites

1. Ensure you have the following installed on your system:
Go (1.23.3) - Download & Install (https://go.dev/dl/)

2. Protocol Buffers Compiler (protoc) - `brew install protobuf`

3. Buf (for Protobuf management) - Install using Homebrew: `brew install bufbuild/buf/buf`

4. Git - Install using Homebrew: `brew install git`

5. Make - Install using Homebrew: `brew install make`

## Setup Instructions:
1. Clone the repository: `git clone git@github.com:Thesohan/weaveGitHubSearchService.git`
2. Install Go dependencies: `make deps`
3. Generate Protobuf code (if modified): `make generate`
4. Set up environment variables: `export GITHUB_API_TOKEN=<your_github_token>`

## Running the Service
1. To start the gRPC server, run: `make server`
2. Running the Client, run: `make client`
3. To execute a test request, run: `make test`

## Troubleshooting
1. Missing Dependencies: Run go mod tidy to install missing Go dependencies.
2. Protobuf Compilation Issues: Ensure protoc and buf are correctly installed.
3. Authentication Errors: Ensure you have set GITHUB_API_TOKEN with a valid GitHub token in your env variable.
