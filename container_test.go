package main

import (
	"strings"
	"testing"
)

func TestEnvironmentContainsHostname(t *testing.T) {
	hostnames := []string{"foo.bar", "baz.bar"}

	fixtures := [][]string{
		[]string{},
		[]string{
			"VIRTUAL_HOST=bar.baz",
		},
		[]string{
			"LETSENCRYPT_HOST=" + strings.Join(hostnames, ","),
		},
		[]string{
			"VIRTUAL_HOST=" + strings.Join(hostnames, ","),
		},
		[]string{
			"VIRTUAL_HOST=baz.foo,baz.bar",
			"LETSENCRYPT_HOST=" + strings.Join(hostnames, ","),
		},
	}

	expected := []bool{
		false,
		false,
		true,
		true,
		true,
	}

	for i := range fixtures {
		if received := environmentContainsHostname(fixtures[i], hostnames); received != expected[i] {
			t.Errorf("Result did not match for variables %v:\nExpected: %v\nReceived: %v\n",
				fixtures[i],
				expected[i],
				received,
			)
		}
	}
}
