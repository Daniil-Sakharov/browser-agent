package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Daniil-Sakharov/BrowserAgent/internal/app"
	"github.com/Daniil-Sakharov/BrowserAgent/internal/config"
	"github.com/Daniil-Sakharov/BrowserAgent/pkg/closer"
)

func main() {
	// 1. Load config
	if err := config.Load(); err != nil {
		fmt.Fprintf(os.Stderr, "❌ Failed to load config: %v\n", err)
		os.Exit(1)
	}

	// 2. Setup context with signals
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()
	defer gracefulShutdown()

	// 3. Create app
	a, err := app.New(ctx)
	if err != nil {
		fmt.Fprintf(os.Stderr, "❌ Failed to create app: %v\n", err)
		os.Exit(1)
	}

	// 4. Execute (Cobra CLI)
	if err := a.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "❌ Error: %v\n", err)
		os.Exit(1)
	}
}

func gracefulShutdown() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_ = closer.CloseAll(ctx)
}
