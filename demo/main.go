package main

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/google/go-github/v41/github"
	"github.com/pactly/client-go/pactly"
	"log"
	"net/http"
)

func main() {
	pactly.Init("client-go-demo-component", "todo")

	demoGet()
	demoPost()
	demoGithub()
	demoGithubRaw()
	demoRateLimit429()
	demoGithubRateLimit()
}

func demoRateLimit429() {
	for i := 0; i < 13; i++ {
		http.Get("https://www.cloudflare.com/rate-limit-test/")
	}
}

func demoGet() {
	// make some dummy API calls
	http.Get("https://jsonplaceholder.typicode.com/todos/1")
}

func demoPost() {
	post := map[string]interface{}{
		"title":  "title",
		"userId": 2,
		"body":   "body",
	}
	postJSON, err := json.Marshal(post)
	if err != nil {
		log.Fatal(err)
	}

	http.Post("https://jsonplaceholder.typicode.com/posts", "application/json", bytes.NewBuffer(postJSON))
	http.Post("http://jsonplaceholder.typicode.com/posts", "application/json", bytes.NewBuffer(postJSON))
}

func demoGithub() {
	client := github.NewClient(nil)

	// list public repositories for org "github"
	opt := &github.RepositoryListByOrgOptions{Type: "public"}
	client.Repositories.ListByOrg(context.Background(), "github", opt)

	// list public repositories for org "getkin"
	opt = &github.RepositoryListByOrgOptions{Type: "public"}
	client.Repositories.ListByOrg(context.Background(), "getkin", opt)

	demoGithubRepositories(client)
}

func demoGithubRepositories(client *github.Client) {
	owner := "github"
	repo := "github"
	branch := "main"
	ctx := context.Background()
	client.Repositories.Get(context.Background(), owner, repo)
	client.Repositories.GetBranch(context.Background(), owner, repo, branch, true)
	client.Repositories.GetBranchProtection(context.Background(), owner, repo, branch)
	client.Repositories.GetCodeOfConduct(context.Background(), owner, repo)
	client.Repositories.GetPagesInfo(ctx, owner, repo)
}

func demoGithubRateLimit() {
	client := github.NewClient(nil)
	opt := &github.RepositoryListByOrgOptions{Type: "public"}

	for i := 0; i < 62; i++ {
		// list public repositories for org "getkin"
		client.Repositories.ListByOrg(context.Background(), "getkin", opt)
	}
}

func demoGithubRaw() {
	http.Get("https://api.github.com/users/octocat/orgs")
}
