package main

import "testing"

func TestPrint(t *testing.T) {
	str := "Roma"
	wanted := "string = Roma"
	if s := Print(str); s != wanted {
		t.Errorf("Print() = %q, want = %q", s, wanted)
	}
}
