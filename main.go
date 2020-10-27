package main

import (
	"context"
	"fmt"
	"log"

	"github.com/shurcooL/githubv4"
	"github.com/skanehira/ght/github"
	"golang.org/x/oauth2"
)

func main() {
	src := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: Config.GitHub.Token},
	)
	httpClient := oauth2.NewClient(context.Background(), src)

	client := githubv4.NewClient(httpClient)

	var q struct {
		github.RepositoryOwner `graphql:"repositoryOwner(login: $login)"`
	}

	variables := map[string]interface{}{
		"login":  githubv4.String("skanehira"),
		"first":  githubv4.Int(20),
		"cursor": (*githubv4.String)(nil),
	}
	if err := client.Query(context.Background(), &q, variables); err != nil {
		log.Fatal(err)
	}

	for _, node := range q.RepositoryOwner.Repositories.Nodes {
		fmt.Println(node.NameWithOwner)
	}

	fmt.Println(q.RepositoryOwner.Repositories.PageInfo.EndCursor)
}