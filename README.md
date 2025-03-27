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
### Mac:
1. Make and Git: `brew install git make`
2. Docker: `brew install docker`
3. Ensure Docker is running:
   ```sh
   open -a Docker
   docker info
   ```
   If Docker is not running, start the Docker daemon manually.

### Linux:
1. Make and Git: `sudo apt-get install git make`
2. Docker: `sudo apt-get install docker`
3. Ensure Docker is running:
   ```sh
   sudo systemctl start docker
   sudo systemctl enable docker  # Optional: Start Docker on boot
   sudo docker info

## Setup Instructions:
1. Clone the repository: `git clone git@github.com:Thesohan/weaveGitHubSearchService.git`
3. Generate Protobuf code (if modified): `make generate`
4. Update environement variables file `.env`: `GITHUB_API_TOKEN=<your_github_token>`

## Running the Service
1. To start the gRPC server and client: `make run`
2. It will ask you to `Enter username`: Enter the username for user level search result (optional)
3. It will ask you to `Enter search term`: Enter the search term to search in the GitHub repository (required)

## Running the Tests
1. To execute a test request, run: `make test`

## Troubleshooting
1. Authentication Errors: Ensure you have set `GITHUB_API_TOKEN` with a valid GitHub token in your env variable.
2. Ensure docker is up and running

### Future Improvements
1. Support for configurable logging
2. Support for pagination
3. Secret management for API tokens
