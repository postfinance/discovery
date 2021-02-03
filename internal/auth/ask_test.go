package auth

import (
	"bytes"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAsk(t *testing.T) {
	r := require.New(t)

	asker := NewAsker(
		WithDfltUsername("test"),
		WithPrompt("Enter"),
	)

	r.NotNil(asker)

	var b bytes.Buffer
	creds, err := asker.Ask(strings.NewReader("Test\n"), &b)

	r.Contains(err.Error(), "reading password from stdin")
	r.Nil(creds)
	r.Equal("Enter (test): Enter SecurID Password: ", b.String())
}
