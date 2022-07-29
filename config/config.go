package config

import (
	"io/ioutil"
	"strings"

	yaml "gopkg.in/yaml.v2"
)

type Connections map[string]string

type Route struct {
	Name       string `yaml:"name"`
	Method     string `yaml:"method"`
	HTTPMethod string `yaml:"http_method"`
	Args       []string
}

type Routes map[string]Route

func matchRoute(a, b string) (bool, []string) {
	routeParts := strings.Split(a, "/")
	pathParts := strings.Split(b, "/")

	var args []string
	matches := 0
	for key, part := range pathParts {
		routePart := routeParts[key]
		if part == routePart {
			matches++
		}

		if strings.HasPrefix(part, ":") {
			matches++
			args = append(args, part[1:])
		}
	}

	return matches == len(pathParts), args
}

func (r Routes) Find(name string) (Route, bool) {
	// Try to find a direct match, first
	s, ok := r[name]

	// If no direct match, look for a partial match
	if !ok {
		for path, route := range r {
			if ok, args := matchRoute(name, path); ok {
				route.Args = args
				return route, true
			}
		}
	}

	return s, ok
}

type Config struct {
	Routes      Routes      `yaml:"routes"`
	Connections Connections `yaml:"connections"`
}

// Parse config
func Parse(f string) (*Config, error) {
	b, err := ioutil.ReadFile(f)
	if err != nil {
		return nil, err
	}

	var c *Config
	if err := yaml.Unmarshal(b, &c); err != nil {
		return nil, err
	}

	return c, nil
}
