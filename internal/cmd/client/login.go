package client

import (
	"fmt"
	"os"
)

type loginCmd struct{}

func (l loginCmd) Run(g *Globals) error {
	fmt.Fprintln(os.Stderr, "login with the following command:")
	fmt.Fprintln(os.Stderr, "")
	fmt.Fprint(os.Stderr, g.OIDC.ExternalLoginCmd+"\n")

	return nil
}
