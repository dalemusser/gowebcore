package sqs

import "testing"

func TestSkip(t *testing.T) {
	t.Skip("integration test requires real SQS; skipped")
}
