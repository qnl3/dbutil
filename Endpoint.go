package main

import (
	"fmt"
)

/*
Endpoint ... */
type Endpoint struct {
	Host  string
	Port  int
	Proto string
	Path  string
}

func (endpoint *Endpoint) String() string {
	var out string

	if endpoint.Proto == "unix" {
		out = endpoint.Path
	} else {
		out = fmt.Sprintf("%s:%d", endpoint.Host, endpoint.Port)
	}

	return out
}
