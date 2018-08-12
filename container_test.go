package main

import (
	"testing"
)

func TestEnvironmentContainsHostname(t *testing.T) {
	hostname := "foo.bar"

	fixtures := [][]string{
		[]string{},
		[]string{
			"VIRTUAL_HOST=bar.baz",
			"LETSENCRYPT_HOST=bar.baz",
		},
		[]string{
			"LETSENCRYPT_HOST=" + hostname,
		},
		[]string{
			"VIRTUAL_HOST=" + hostname,
		},
	}

	expected := []bool{
		false,
		false,
		true,
		true,
	}

	for i := range fixtures {
		if received := environmentContainsHostname(fixtures[i], hostname); received != expected[i] {
			t.Errorf("Result did not match for variables %v:\nExpected: %v\nReceived: %v\n",
				fixtures[i],
				expected[i],
				received,
			)
		}
	}
}
