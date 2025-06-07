package db

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Manager holds multiple Mongo/DocumentDB connections identified by alias.
type Manager struct {
	mu        sync.RWMutex
	clients   map[string]*mongo.Client
	databases map[string]*mongo.Database
}

// NewManager returns an empty Manager.
func NewManager() *Manager {
	return &Manager{
		clients:   make(map[string]*mongo.Client),
		databases: make(map[string]*mongo.Database),
	}
}

// Add establishes a new client and stores it under alias.
// Safe to call multiple times; reuses an existing client if already present.
func (m *Manager) Add(alias, uri, dbName string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, ok := m.databases[alias]; ok {
		return nil // already added
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return fmt.Errorf("mongo connect: %w", err)
	}
	m.clients[alias] = client
	m.databases[alias] = client.Database(dbName)
	return nil
}

// DB returns the *mongo.Database registered under alias, or nil if unknown.
func (m *Manager) DB(alias string) *mongo.Database {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.databases[alias]
}

// Ping checks a single alias (or "primary") and returns error if unreachable.
func (m *Manager) Ping(alias string) error {
	db := m.DB(alias)
	if db == nil {
		return errors.New("db alias not found")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	return db.Client().Ping(ctx, nil)
}
