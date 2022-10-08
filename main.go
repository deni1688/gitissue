package main

import (
	"deni1688/gitissue/domain"
	"deni1688/gitissue/infra"
	"errors"
	"flag"
	"fmt"
	"net/http"
	"strings"
)

func main() {
	setup := flag.Bool("setup", false, "setup config file")
	path := flag.String("path", "./issues.txt", "please provide file path to parse issues from")
	flag.Parse()

	if *path == "" {
		fmt.Println("Please provide file path to parse issues from")
		return
	}

	c := new(config)
	if *setup {
		if err := c.Setup(); err != nil {
			fmt.Println(err)
		}
		return
	}

	if err := c.Load(); err != nil {
		fmt.Println("Error reading config at $HOME/.config/gitissue.json")
		return
	}

	p, err := getProvider(c)
	if err != nil {
		fmt.Println("Error getting provider: ", err)
		return
	}

	n := infra.NewWebhookNotifier(c.WebHooks)
	s := domain.NewService(p, n, c.Prefix)
	cli := infra.NewCli(*path, s)
	if err = cli.Execute(); err != nil {
		fmt.Println("Error running cli:", err)
		return
	}

	fmt.Println("Done!")
}

func getProvider(c *config) (domain.GitProvider, error) {
	switch {
	case strings.Contains(c.Host, "gitlab"):
		return infra.NewGitlab(c.Token, c.Host, c.Query, http.DefaultClient), nil
	case strings.Contains(c.Host, "github"):
		return infra.NewGithub(c.Token, c.Host, c.Query, http.DefaultClient), nil
	default:
		return nil, errors.New("invalid provider " + c.Host)
	}
}
