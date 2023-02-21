package client

import (
	"fmt"
	"os"

	"github.com/postfinance/discovery/internal/auth"
)

type loginCmd struct{}

func (l loginCmd) Run(g *Globals) error {
	if g.OIDC.ExternalLoginCmd != "" {
		fmt.Fprintln(os.Stderr, "login with the following command:")
		fmt.Fprintln(os.Stderr, "")
		fmt.Fprint(os.Stderr, g.OIDC.ExternalLoginCmd+"\n")

		return nil
	}

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

	cli, err := auth.NewClient(g.OIDC.Endpoint, g.OIDC.ClientID, opts...)
	if err != nil {
		return err
	}

	t, err := cli.Token(c.Username, c.Password, os.Stdout)
	if err != nil {
		return err
	}

	return g.saveToken(t)
}
