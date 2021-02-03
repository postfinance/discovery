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

	opts := []auth.ClientOption{}

	if g.CACert != "" {
		pool, err := auth.AppendCertsToSystemPool(g.CACert)
		if err != nil {
			return err
		}

		transport := auth.NewTLSTransportFromCertPool(pool)

		opts = append(opts, auth.WithTransport(transport))
	}

	cli, err := auth.NewClient(g.OIDCEndpoint, g.OIDCClientID, opts...)
	if err != nil {
		return err
	}

	t, err := cli.Token(c.Username, c.Password, os.Stdout)
	if err != nil {
		return err
	}

	return g.saveToken(t)
}
