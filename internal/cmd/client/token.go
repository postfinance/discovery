package client

import (
	"fmt"
	"time"

	"connectrpc.com/connect"
	"github.com/alecthomas/kong"
	"github.com/postfinance/discovery/internal/server/convert"
	discoveryv1 "github.com/postfinance/discovery/pkg/discoverypb/postfinance/discovery/v1"
	"go.uber.org/zap"
)

type tokenCmd struct {
	Create tokenCreate `cmd:"" help:"Create an access token."`
	Info   tokenInfo   `cmd:"" help:"Get information about created tokens."`
}

type tokenCreate struct {
	ID         string        `arg:"" short:"i" help:"An ID that can identify token (i.e: username)" required:"true"`
	Expiry     time.Duration `short:"e" default:"0" help:"How long (duration) should the token be valid. 0 is forever."`
	Namespaces []string      `short:"n" help:"The namespaces the token has access to." required:"true"`
}

func (t tokenCreate) Run(g *Globals, l *zap.SugaredLogger, c *kong.Context) error {
	cli, err := g.tokenClient()
	if err != nil {
		return err
	}

	ctx, cancel := g.ctx()
	defer cancel()

	r, err := cli.Create(ctx, connect.NewRequest(&discoveryv1.CreateRequest{
		Expires:    t.Expiry.String(),
		Id:         t.ID,
		Namespaces: t.Namespaces,
	}))
	if err != nil {
		return err
	}

	fmt.Println(r.Msg.GetToken())

	return nil
}

type tokenInfo struct {
	Token string `arg:"" short:"i" help:"The jwt token string" required:"true"`
}

func (t tokenInfo) Run(g *Globals, l *zap.SugaredLogger, c *kong.Context) error {
	cli, err := g.tokenClient()
	if err != nil {
		return err
	}

	ctx, cancel := g.ctx()
	defer cancel()

	r, err := cli.Info(ctx, connect.NewRequest(&discoveryv1.InfoRequest{
		Token: t.Token,
	}))
	if err != nil {
		return err
	}

	i := r.Msg.GetTokeninfo()
	expiry := convert.TimeFromPB(i.GetExpiresAt())
	expiryStr := expiry.Format(time.RFC3339)

	if expiry.IsZero() {
		expiryStr = "never"
	}

	fmt.Println("id:", i.GetId())
	fmt.Println("namespaces:", i.GetNamespaces())
	fmt.Println("expiry:", expiryStr)

	return nil
}
