package auth

import (
	"bufio"
	"fmt"
	"io"
	"strings"
	"syscall"

	"golang.org/x/term"
)

// Credentials contains username and password.
type Credentials struct {
	Username string
	Password string
}

// Option is a functional option to configure the Asker.
type Option func(a *Asker)

// NewAsker creates a new Asker.
func NewAsker(opts ...Option) *Asker {
	a := Asker{
		prompt: "Enter Username",
	}

	for _, opt := range opts {
		opt(&a)
	}

	return &a
}

// WithDfltUsername offers a dflt username in prompt.
func WithDfltUsername(username string) Option {
	return func(a *Asker) {
		a.username = username
	}
}

// WithPrompt overrides default `Enter Username` prompt.
func WithPrompt(prompt string) Option {
	return func(a *Asker) {
		stripped := strings.TrimRight(prompt, ":")
		a.prompt = stripped
	}
}

// Asker asks for username and password.
type Asker struct {
	prompt   string
	username string
}

// Ask aks the username and password.
func (a *Asker) Ask(in io.Reader, out io.Writer) (*Credentials, error) {
	var (
		c      Credentials
		prompt string
	)

	prompt = a.prompt

	if a.username != "" {
		prompt = fmt.Sprintf("%s (%s)", prompt, a.username)
	}

	prompt += ": "

	// username prompt, with offered username
	fmt.Fprint(out, prompt)

	reader := bufio.NewReader(in)

	userInput, err := reader.ReadString('\n')
	if err != nil {
		return nil, fmt.Errorf("reading from stdin: %w", err)
	}

	c.Username = strings.TrimSpace(userInput)

	if c.Username == "" {
		// set the default username
		c.Username = a.username
	}

	p, err := readPassword("Enter SecurID Password", out)
	if err != nil {
		return nil, err
	}

	c.Password = strings.TrimSpace(p)

	return &c, nil
}

func readPassword(prompt string, out io.Writer) (string, error) {
	fmt.Fprintf(out, "%s: ", prompt)

	//nolint:unconvert // needs to be done for windows
	// conversion for GOOS=windows
	bytePassword, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return "", fmt.Errorf("reading password from stdin: %w", err)
	}

	fmt.Println("")

	return strings.TrimSpace(string(bytePassword)), nil
}
