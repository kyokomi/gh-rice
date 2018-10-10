package main

import (
	"context"
	"flag"
	"fmt"
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
		if isAlreadyComment(client, ctx, owner, repo, number, comment) {
			fmt.Println("already comment")
			return
		}

		_, _, err := client.Issues.CreateComment(ctx, owner, repo, number, &github.IssueComment{
			Body: &comment,
		})
		if err != nil {
			log.Fatalln(err)
		}
	}
}

func isAlreadyComment(client *github.Client, ctx context.Context, owner, repo string, number int, comment string) bool {
	comments, _, err := client.Issues.ListComments(ctx, owner, repo, number, nil)
	if err != nil {
		log.Fatalln(err)
	}

	for i := range comments {
		fmt.Println(*comments[i].Body)
		if *comments[i].Body == comment {
			return true
		}
	}

	return false
}
