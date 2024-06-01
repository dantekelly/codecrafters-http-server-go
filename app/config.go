package main

import (
	"errors"
	"strings"
)

type Config struct {
	directory string
}

type ConfigBuilder struct {
	directory *string
}

func (b *ConfigBuilder) Directory(d string) *ConfigBuilder {
	if !strings.HasSuffix(d, "/") {
		d += "/"
	}

	b.directory = &d
	return b
}

func (b *ConfigBuilder) Build() (*Config, error) {
	if b.directory == nil {
		return nil, errors.New("directory is required")
	}

	return &Config{directory: *b.directory}, nil
}
