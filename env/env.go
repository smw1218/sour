package env

import (
	"os"
)

type Env string

const (
	Local Env = "local"
	Prod  Env = "prod"
)

// EnvVar is the environment variable where the deployment environment is stored
var EnvVar = "DEPLOYMENT_ENVIRONMENT"

var env Env

func Init() {
	ReadFromEnv()
}

// ReadFromEnv will set the internal env value to
// what is set in the EnvVar environment variable.
// This is not thread safe and should only be set
// once at boot-time.
// If the environment is not set or empty, then the
// env defaults to [Local].
func ReadFromEnv() {
	varval := Env(os.Getenv(EnvVar))
	if varval == "" {
		varval = Local
	}
	setEnv(varval)
}

// to lock or not to lock, that is the question
func setEnv(e Env) {
	env = e
}

// Get returns the current env
func Get() Env {
	return env
}

func (e Env) IsLocal() bool {
	return e == Local
}

func (e Env) IsProd() bool {
	return e == Prod
}

// TestWith should only be used to test env-specific
// code. It temporarily sets the Env to e, runs testFunc
// then restores the Env to the previous value. Since
// Get() reads from a shared package variable, using this is
// not safe for t.Parallel
func TestWith(e Env, testFunc func()) {
	saved := env
	env = e
	testFunc()
	setEnv(saved)
}
