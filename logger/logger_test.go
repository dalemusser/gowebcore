package logger

import "testing"

func TestInit(t *testing.T) {
	// initialise the global logger
	Init("debug")

	// make sure Init() populated lg and L() exposes it
	if L() == nil {
		t.Fatal("logger not initialised")
	}

	// optional: ensure a helper call doesn’t panic
	Info("test log entry", "unit", "logger_test")
}
