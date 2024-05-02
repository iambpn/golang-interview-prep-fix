package main

import "testing"

func TestRun(t *testing.T) {
	err := run(":3030", "", "", "")

	if err != nil {
		return
	}

	t.Fatal("server error")
}
