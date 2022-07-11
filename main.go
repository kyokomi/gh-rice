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
	var overwrite bool
	flag.StringVar(&comment, "c", "", "comment")
	flag.StringVar(&token, "t", "", "github access token")
	flag.StringVar(&owner, "o", "", "github owner name")
	flag.StringVar(&repo, "r", "", "github repository name")
	flag.StringVar(&branch, "b", "", "github branch name")
	flag.BoolVar(&overwrite, "w", false, "overwrite and update if there are already comments")
	flag.Parse()

	if comment == "" {
		fmt.Println("no comment")
		return // コメントなしは何もしない
	}

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

		// 上書き更新オプションがONの場合
		if overwrite {
			// tokenのuserIDを取得
			u, _, err := client.Users.Get(ctx, "")
			if err != nil {
				log.Fatalln(err)
			}

			comments, _, err := client.Issues.ListComments(ctx, owner, repo, number, nil)
			if err != nil {
				log.Fatalln(err)
			}

			var commentID *int64
			for _, c := range comments {
				if *c.User.ID == *u.ID {
					commentID = c.ID
				}
			}

			if commentID != nil {
				_, _, err := client.Issues.EditComment(ctx, owner, repo, *commentID, &github.IssueComment{
					Body: &comment,
				})
				if err != nil {
					log.Fatalln(err)
				}
				return // 上書き更新で終わり
			}
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
		if *comments[i].Body == comment {
			return true
		}
	}

	return false
}
