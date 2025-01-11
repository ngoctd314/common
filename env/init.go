package env

import (
	"strings"
	"sync"

	"github.com/spf13/viper"
)

var (
	// immutable: global config object
	immutable *config
	// sentinel to make sure that Init func is called only once
	initOnce sync.Once
)

// wrapper Viper as our config
type config struct {
	*viper.Viper
}

func init() {
	immutable = automaticEnv()
}

func automaticEnv(opts ...option) *config {
	cnf := &config{
		viper.New(),
	}
	for _, opt := range opts {
		opt(cnf)
	}

	// inject ENVIRONMENT VARIABLE of current host into *config
	cnf.AutomaticEnv()
	// replace . -> _ because env use _, but in we use . in config file
	replacer := strings.NewReplacer(".", "_")
	cnf.SetEnvKeyReplacer(replacer)

	return cnf
}

// This function is set up to be call only once
// use it only when you want to set up
//   - custom config file path
//   - custom prefix env
func Init(opts ...option) {
	// it is confused when Init config twice or more
	initOnce.Do(func() {
		immutable = automaticEnv(opts...)
	})
}

const (
	DebugMode   = "debug"
	DevMode     = "dev"
	TestingMode = "testing"
	ProdMode    = "prod"
)
