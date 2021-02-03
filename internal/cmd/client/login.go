package client

import (
	"os"

	"github.com/postfinance/discovery/internal/auth"
)

type loginCmd struct {
}

func (l loginCmd) Run(g *Globals) error {
	user := os.Getenv("USER")
	asker := auth.NewAsker(auth.WithPrompt("Enter username: "), auth.WithDfltUsername(user))

	c, err := asker.Ask(os.Stdin, os.Stdout)
	if err != nil {
		return err
	}

	cli, err := auth.NewClient(g.OIDCEndpoint, g.OIDCClientID)
	if err != nil {
		return err
	}

	t, err := cli.Token(c.Username, c.Password, os.Stdout)
	if err != nil {
		return err
	}

	return g.saveToken(t)
}
