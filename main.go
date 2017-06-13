package main

import (
	"context"
	"flag"
	"log"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

func main() {
	var token, comment, owner, repo, branch string
	flag.StringVar(&comment, "c", "", "comment")
	flag.StringVar(&token, "t", "", "github access token")
	flag.StringVar(&owner, "o", "", "github owner name")
	flag.StringVar(&repo, "r", "", "github repository name")
	flag.StringVar(&branch, "b", "", "github branch name")
	flag.Parse()

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)

	pulls, _, err := client.PullRequests.List(ctx, owner, repo, nil)
	if err != nil {
		log.Fatalln(err)
	}

	for _, r := range pulls {
		if *r.Head.Ref != branch {
			continue
		}

		number := *r.Number
		_, _, err := client.Issues.CreateComment(ctx, owner, repo, number, &github.IssueComment{
			Body: &comment,
		})
		if err != nil {
			log.Fatalln(err)
		}
	}
}
