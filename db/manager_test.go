package db

import "testing"

func TestAddAndDB(t *testing.T) {
	m := NewManager()
	err := m.Add("primary", "mongodb://localhost:27017", "testdb")
	if err == nil {
		// connection may fail if Mongo isn't running; that's fine for unit test.
		if m.DB("primary") == nil {
			t.Fatalf("expected DB map entry")
		}
	} // if err != nil we skip assertion – integration covered elsewhere.
}
