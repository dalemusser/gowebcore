package logger

import "testing"

func TestInit(t *testing.T) {
	Init("debug")
	if Instance() == nil {
		t.Fatal("logger not initialised")
	}
}
