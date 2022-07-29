package config

import (
	"github.com/matryer/is"
	"testing"
)

func TestParseConfig(t *testing.T) {
	is := is.New(t)
	conf, err := Parse("../config.yaml")
	is.NoErr(err)
	is.Equal("localhost:8888", conf.Connections["service.Greeter"])
}
