package server

import (
	"errors"
	"os"
	"regexp"
	"time"

	"github.com/postfinance/profiler"
	"github.com/postfinance/store"
	"github.com/postfinance/store/etcd"
	"github.com/zbindenren/king"
)

// Globals are the global server flags.
type Globals struct {
	Etcd `prefix:"etcd-"`
	// EnvHelp    king.EnvHelpFlag `help:"Show context-sensitive help about environment variables"`
	Debug      bool             `help:"Show debug output"`
	ShowConfig king.ShowConfig  `help:"Show used config files"`
	Version    king.VersionFlag `help:"Show version information"`
	Profiler   profilerFlags    `embed:"true" prefix:"profiler-"`
}

// Etcd contains the etcd configuration.
type Etcd struct {
	Endpoints        []string      `help:"Etcd endpoints" default:"localhost:2379"`
	User             string        `help:"Etcd user"`
	Password         string        `help:"Etcd password"`
	Cert             string        `help:"Etcd certificate"`
	Key              string        `help:"Etcd key"`
	Prefix           string        `help:"Etcd prefix" default:"/discovery"`
	CA               string        `help:"Etcd CA"`
	CertFile         string        `help:"Etcd cert file"`
	KeyFile          string        `help:"Etcd key file"`
	CAFile           string        `help:"Etcd CA file"`
	AutoSyncInterval time.Duration `help:"Etcd autosync interval" default:"10s"`
	DialTimeout      time.Duration `help:"Etcd dial timeout" default:"5s"`
	RequestTimeout   time.Duration `help:"Etcd request timeout" default:"5s"`
}

var (
	isValidPrefix = regexp.MustCompile(`^/[a-z-]+$`)
)

func (e Etcd) backend() (store.Backend, error) {
	if !isValidPrefix.MatchString(e.Prefix) {
		return nil, errors.New("store prefix must start with '/' followed by at least one letter in the range 'a-z'")
	}

	return etcd.New(
		etcd.WithEndpoints(e.Endpoints),
		etcd.WithUsername(e.User),
		etcd.WithPrefix(e.Prefix),
		etcd.WithPassword(e.Password),
		etcd.WithKey(e.Key),
		etcd.WithKeyFile(e.KeyFile),
		etcd.WithCert(e.Cert),
		etcd.WithCertFile(e.CertFile),
		etcd.WithCA(e.CA),
		etcd.WithCAFile(e.CAFile),
		etcd.WithDialTimeout(e.DialTimeout),
		etcd.WithRequestTimeout(e.RequestTimeout),
		etcd.WithAutoSyncInterval(e.AutoSyncInterval),
	)
}

type profilerFlags struct {
	Enabled bool          `help:"Enable the profiler."`
	Listen  string        `help:"The profiles listen address." default:":6666"`
	Timeout time.Duration `help:"The timeout after the pprof handler will be shutdown." default:"5m"`
}

func (p profilerFlags) New(s os.Signal, h ...profiler.Hooker) *profiler.Profiler {
	return profiler.New(
		profiler.WithAddress(p.Listen),
		profiler.WithTimeout(p.Timeout),
		profiler.WithSignal(s),
		profiler.WithHooks(h...),
	)
}
