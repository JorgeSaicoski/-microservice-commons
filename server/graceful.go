package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// GracefulShutdownConfig holds configuration for graceful shutdown
type GracefulShutdownConfig struct {
	Timeout       time.Duration // Maximum time to wait for shutdown
	SignalTimeout time.Duration // Time to wait for signal
}

// DefaultGracefulConfig returns default graceful shutdown configuration
func DefaultGracefulConfig() GracefulShutdownConfig {
	return GracefulShutdownConfig{
		Timeout:       30 * time.Second,
		SignalTimeout: 5 * time.Second,
	}
}

// ShutdownManager handles graceful shutdown of the server
type ShutdownManager struct {
	server *http.Server
	config GracefulShutdownConfig
}

// NewShutdownManager creates a new shutdown manager
func NewShutdownManager(server *http.Server, config GracefulShutdownConfig) *ShutdownManager {
	return &ShutdownManager{
		server: server,
		config: config,
	}
}

// WaitForShutdown waits for shutdown signals and gracefully shuts down the server
func (sm *ShutdownManager) WaitForShutdown() error {
	// Create a channel to receive OS signals
	quit := make(chan os.Signal, 1)

	// Register the channel to receive specific signals
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Block until we receive our signal
	sig := <-quit
	fmt.Printf("\nReceived signal: %v. Initiating graceful shutdown...\n", sig)

	// Create a deadline to wait for
	ctx, cancel := context.WithTimeout(context.Background(), sm.config.Timeout)
	defer cancel()

	// Attempt the graceful shutdown by closing the listener
	// and completing all inflight requests
	if err := sm.server.Shutdown(ctx); err != nil {
		return fmt.Errorf("server forced to shutdown: %w", err)
	}

	fmt.Println("Server shutdown complete")
	return nil
}

// StartWithGracefulShutdown starts the server and handles graceful shutdown
func (sm *ShutdownManager) StartWithGracefulShutdown() error {
	// Start server in a goroutine
	go func() {
		fmt.Printf("Starting server on %s\n", sm.server.Addr)
		if err := sm.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("Server failed to start: %v\n", err)
			os.Exit(1)
		}
	}()

	// Wait for interrupt signal to gracefully shut down the server
	return sm.WaitForShutdown()
}

// ForceShutdown forces immediate shutdown of the server
func (sm *ShutdownManager) ForceShutdown() error {
	fmt.Println("Forcing immediate server shutdown...")
	return sm.server.Close()
}

// setupGracefulShutdown is a helper function to setup graceful shutdown for a server
func setupGracefulShutdown(server *http.Server) *ShutdownManager {
	config := DefaultGracefulConfig()
	return NewShutdownManager(server, config)
}
