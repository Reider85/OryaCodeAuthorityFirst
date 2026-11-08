# Small Go Project for Testing
**Repository Size**: < 100 files  
**Purpose**: Testing basic caf functionality with small codebase

---

## Repository Structure

```
small-go-project/
├── cmd/
│   └── app/
│       └── main.go
├── internal/
│   ├── config/
│   │   └── config.go
│   ├── handler/
│   │   └── handler.go
│   └── model/
│       └── model.go
├── pkg/
│   ├── utils/
│   │   └── utils.go
│   └── database/
│       └── database.go
├── go.mod
├── go.sum
├── .gitignore
└── README.md
```

---

## File Contents

### cmd/app/main.go
```go
package main

import (
	"fmt"
	"log"
	"small-go-project/internal/config"
	"small-go-project/internal/handler"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Create handler
	h := handler.New(cfg)
	
	// Start the application
	if err := h.Start(); err != nil {
		log.Fatalf("Failed to start application: %v", err)
	}

	fmt.Println("Application started successfully")
}
```

### internal/config/config.go
```go
package config

import (
	"fmt"
	"os"
	"path/filepath"
)

type Config struct {
	Port     string
	Database DatabaseConfig
	Logging  LoggingConfig
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Database string
}

type LoggingConfig struct {
	Level  string
	Format string
	Output string
}

func Load() (*Config, error) {
	return &Config{
		Port: "8080",
		Database: DatabaseConfig{
			Host:     "localhost",
			Port:     "5432",
			User:     "postgres",
			Password: "password",
			Database: "app",
		},
		Logging: LoggingConfig{
			Level:  "info",
			Format: "json",
			Output: "stdout",
		},
	}, nil
}
```

### internal/handler/handler.go
```go
package handler

import (
	"context"
	"fmt"
	"log"
	"small-go-project/internal/config"
	"small-go-project/internal/model"
	"time"
)

type Handler struct {
	config *config.Config
	model  *model.Model
}

func New(cfg *config.Config) *Handler {
	return &Handler{
		config: cfg,
		model:  model.New(),
	}
}

func (h *Handler) Start() error {
	// Initialize the model
	if err := h.model.Initialize(); err != nil {
		return fmt.Errorf("failed to initialize model: %w", err)
	}

	// Start background processing
	go h.backgroundProcessing()

	log.Println("Handler started successfully")
	return nil
}

func (h *Handler) backgroundProcessing() {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if err := h.processData(); err != nil {
				log.Printf("Error processing data: %v", err)
			}
		}
	}
}

func (h *Handler) processData() error {
	// Process some data
	data := []string{"item1", "item2", "item3"}
	
	for _, item := range data {
		if err := h.model.Process(context.Background(), item); err != nil {
			return err
		}
	}
	
	return nil
}
```

### internal/model/model.go
```go
package model

import (
	"context"
	"fmt"
	"small-go-project/pkg/database"
	"small-go-project/pkg/utils"
)

type Model struct {
	db *database.Database
}

func New() *Model {
	return &Model{
		db: database.New(),
	}
}

func (m *Model) Initialize() error {
	// Initialize database connection
	if err := m.db.Connect(); err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}
	
	return nil
}

func (m *Model) Process(ctx context.Context, data string) error {
	// Process the data
	processed := utils.Process(data)
	
	// Store in database
	if err := m.db.Store(ctx, processed); err != nil {
		return fmt.Errorf("failed to store data: %w", err)
	}
	
	return nil
}
```

### pkg/utils/utils.go
```go
package utils

import (
	"fmt"
	"strings"
)

func Process(input string) string {
	// Simple processing logic
	processed := strings.ToUpper(input)
	return fmt.Sprintf("Processed: %s", processed)
}

func Validate(input string) bool {
	return len(input) > 0
}
```

### pkg/database/database.go
```go
package database

import (
	"context"
	"fmt"
	"small-go-project/internal/config"
)

type Database struct {
	config *config.DatabaseConfig
}

func New() *Database {
	return &Database{
		config: &config.DatabaseConfig{},
	}
}

func (d *Database) Connect() error {
	// Simulate database connection
	fmt.Printf("Connecting to database at %s:%s\n", d.config.Host, d.config.Port)
	return nil
}

func (d *Database) Store(ctx context.Context, data string) error {
	// Simulate storing data
	fmt.Printf("Storing data: %s\n", data)
	return nil
}
```

### go.mod
```go
module small-go-project

go 1.22

require (
	github.com/lib/pq v1.10.9
)

require (
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20221227161230-091c0ba34f0a // indirect
	github.com/jackc/pgx/v5 v5.4.3 // indirect
	golang.org/x/crypto v0.9.0 // indirect
	golang.org/x/sys v0.8.0 // indirect
	golang.org/x/text v0.9.0 // indirect
)
```

### .gitignore
```
# Binaries for programs and plugins
*.exe
*.exe~
*.dll
*.so
*.dylib

# Test binary, built with `go test -c`
*.test

# Output of the go coverage tool, specifically when used with LiteIDE
*.out

# Dependency directories (remove the comment below to include it)
vendor/

# Go workspace file
go.work

# Environment variables
.env
.env.local

# IDE files
.vscode/
.idea/
*.swp
*.swo
*~

# OS files
.DS_Store
Thumbs.db
```

### README.md
```markdown
# Small Go Project

This is a small Go project for testing CodeAuthorityFirst functionality.

## Features

- Basic web application
- Database integration
- Configuration management
- Background processing

## Getting Started

1. Clone the repository
2. Run `go mod tidy`
3. Run the application with `go run cmd/app/main.go`

## Testing

Run the test suite with:
```bash
go test ./...
```
```

---

## Testing Scenarios

### 1. Basic Session Tracking
- Initialize caf in the repository
- Make changes to various files
- Verify session tracking works correctly

### 2. Hook Integration
- Test pre-session, post-edit, and pre-commit hooks
- Verify hook execution timing and success rate

### 3. Performance Testing
- Measure session creation time
- Test hook execution performance
- Verify memory usage is within limits

### 4. Error Handling
- Simulate various error scenarios
- Test graceful error handling
- Verify error logging and recovery

---

## Expected Results

- **Session Tracking**: 100% success rate with complete metadata
- **Hook Integration**: 100% success rate with <1s execution time
- **Performance**: <100ms session creation, <1s hook execution
- **Error Handling**: Graceful recovery from all error scenarios

---

## Setup Instructions

1. Navigate to the repository directory:
   ```bash
   cd repositories/small-go-project
   ```

2. Initialize caf:
   ```bash
   caf init --auto-detect
   ```

3. Start monitoring:
   ```bash
   cd ../..
   ./start-monitoring.ps1
   ```

4. Begin testing:
   ```bash
   # Make some changes to trigger sessions
   echo "// New feature" >> internal/config/config.go
   git add .
   git commit -m "Add new feature"
   ```

5. Monitor results:
   ```bash
   caf doctor
   caf dashboard
   ```