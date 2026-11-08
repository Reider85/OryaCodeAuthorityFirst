# Medium Go Project for Testing
**Repository Size**: 100-500 files  
**Purpose**: Testing caf functionality with medium-sized codebase and more complex scenarios

---

## Repository Structure

```
medium-go-project/
├── cmd/
│   ├── server/
│   │   └── main.go
│   ├── client/
│   │   └── main.go
│   └── migration/
│       └── main.go
├── internal/
│   ├── api/
│   │   ├── handlers/
│   │   │   ├── user_handler.go
│   │   │   ├── product_handler.go
│   │   │   └── order_handler.go
│   │   ├── middleware/
│   │   │   ├── auth.go
│   │   │   ├── logging.go
│   │   │   └── cors.go
│   │   └── routes/
│   │       └── routes.go
│   ├── config/
│   │   ├── config.go
│   │   ├── database.go
│   │   └── logging.go
│   ├── service/
│   │   ├── user_service.go
│   │   ├── product_service.go
│   │   ├── order_service.go
│   │   └── notification_service.go
│   ├── repository/
│   │   ├── user_repository.go
│   │   ├── product_repository.go
│   │   └── order_repository.go
│   ├── model/
│   │   ├── user.go
│   │   ├── product.go
│   │   └── order.go
│   ├── cache/
│   │   ├── cache.go
│   │   └── redis.go
│   ├── worker/
│   │   ├── worker.go
│   │   └── jobs.go
│   └── utils/
│       ├── validator.go
│       ├── crypto.go
│       ├── email.go
│       └── file.go
├── pkg/
│   ├── database/
│   │   ├── connection.go
│   │   ├── queries.go
│   │   └ migrations/
│   │       ├── 001_initial_schema.sql
│   │       └── 002_add_indexes.sql
│   ├── auth/
│   │   ├── jwt.go
│   │   └── bcrypt.go
│   ├── logger/
│   │   └── logger.go
│   ├── metrics/
│   │   └── metrics.go
│   └── errors/
│       └── errors.go
├── scripts/
│   ├── deploy.sh
│   ├── migrate.sh
│   └── test.sh
├── configs/
│   ├── development.yaml
│   ├── production.yaml
│   └── test.yaml
├── docs/
│   ├── api.md
│   ├── deployment.md
│   └── architecture.md
├── go.mod
├── go.sum
├── .gitignore
├── Dockerfile
├── docker-compose.yml
└── README.md
```

---

## Key Features

### 1. Multi-Component Architecture
- **API Layer**: RESTful API with handlers, middleware, and routes
- **Service Layer**: Business logic with user, product, and order services
- **Repository Layer**: Data access with database operations
- **Model Layer**: Data structures and business entities

### 2. Advanced Features
- **Authentication**: JWT-based authentication with bcrypt
- **Caching**: Redis-based caching system
- **Background Processing**: Worker system for background jobs
- **Metrics**: Application metrics and monitoring
- **Configuration**: Multi-environment configuration system

### 3. Testing Infrastructure
- **Unit Tests**: Comprehensive unit test coverage
- **Integration Tests**: API and database integration tests
- **E2E Tests**: End-to-end testing scenarios

---

## File Contents (Key Files)

### cmd/server/main.go
```go
package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"medium-go-project/internal/api/routes"
	"medium-go-project/internal/config"
	"medium-go-project/pkg/logger"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize logger
	logger := logger.New(cfg.Logging)

	// Setup routes
	router := routes.Setup(cfg, logger)

	// Create HTTP server
	server := &http.Server{
		Addr:    ":" + cfg.Server.Port,
		Handler: router,
	}

	// Start server in goroutine
	go func() {
		logger.Info("Starting server", "port", cfg.Server.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Error("Server shutdown error", "error", err)
	}

	logger.Info("Server exited")
}
```

### internal/api/routes/routes.go
```go
package routes

import (
	"medium-go-project/internal/api/handlers"
	"medium-go-project/internal/api/middleware"
	"medium-go-project/internal/config"
	"medium-go-project/pkg/logger"

	"github.com/gorilla/mux"
)

func Setup(cfg *config.Config, logger *logger.Logger) *mux.Router {
	router := mux.NewRouter()

	// Middleware
	router.Use(middleware.CORS)
	router.Use(middleware.Logging(logger))
	router.Use(middleware.Recovery)

	// Public routes
	public := router.PathPrefix("/api").Subrouter()
	public.HandleFunc("/health", handlers.HealthCheck).Methods("GET")
	public.HandleFunc("/login", handlers.Login).Methods("POST")

	// Protected routes
	protected := router.PathPrefix("/api").Subrouter()
	protected.Use(middleware.Authenticate)
	
	// User routes
	userHandler := handlers.NewUserHandler(cfg, logger)
	protected.HandleFunc("/users", userHandler.CreateUser).Methods("POST")
	protected.HandleFunc("/users/{id}", userHandler.GetUser).Methods("GET")
	protected.HandleFunc("/users/{id}", userHandler.UpdateUser).Methods("PUT")
	protected.HandleFunc("/users/{id}", userHandler.DeleteUser).Methods("DELETE")

	// Product routes
	productHandler := handlers.NewProductHandler(cfg, logger)
	protected.HandleFunc("/products", productHandler.CreateProduct).Methods("POST")
	protected.HandleFunc("/products", productHandler.ListProducts).Methods("GET")
	protected.HandleFunc("/products/{id}", productHandler.GetProduct).Methods("GET")
	protected.HandleFunc("/products/{id}", productHandler.UpdateProduct).Methods("PUT")
	protected.HandleFunc("/products/{id}", productHandler.DeleteProduct).Methods("DELETE")

	// Order routes
	orderHandler := handlers.NewOrderHandler(cfg, logger)
	protected.HandleFunc("/orders", orderHandler.CreateOrder).Methods("POST")
	protected.HandleFunc("/orders", orderHandler.ListOrders).Methods("GET")
	protected.HandleFunc("/orders/{id}", orderHandler.GetOrder).Methods("GET")
	protected.HandleFunc("/orders/{id}", orderHandler.UpdateOrder).Methods("PUT")
	protected.HandleFunc("/orders/{id}", orderHandler.DeleteOrder).Methods("DELETE")

	return router
}
```

### internal/config/config.go
```go
package config

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Server   ServerConfig        `yaml:"server"`
	Database DatabaseConfig      `yaml:"database"`
	Auth     AuthConfig          `yaml:"auth"`
	Cache    CacheConfig         `yaml:"cache"`
	Logging  LoggingConfig       `yaml:"logging"`
	Metrics  MetricsConfig       `yaml:"metrics"`
}

type ServerConfig struct {
	Port         string        `yaml:"port"`
	ReadTimeout  time.Duration `yaml:"read_timeout"`
	WriteTimeout time.Duration `yaml:"write_timeout"`
	IdleTimeout  time.Duration `yaml:"idle_timeout"`
}

type DatabaseConfig struct {
	Host            string `yaml:"host"`
	Port            string `yaml:"port"`
	User            string `yaml:"user"`
	Password        string `yaml:"password"`
	Database        string `yaml:"database"`
	MaxConnections  int    `yaml:"max_connections"`
	MaxIdleTime     string `yaml:"max_idle_time"`
}

type AuthConfig struct {
	JWTSecret       string        `yaml:"jwt_secret"`
	JWTExpiration   time.Duration `yaml:"jwt_expiration"`
	BcryptCost     int           `yaml:"bcrypt_cost"`
}

type CacheConfig struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Password string `yaml:"password"`
	DB       int    `yaml:"db"`
}

type LoggingConfig struct {
	Level  string `yaml:"level"`
	Format string `yaml:"format"`
	Output string `yaml:"output"`
}

type MetricsConfig struct {
	Enabled bool   `yaml:"enabled"`
	Port    string `yaml:"port"`
}

func Load() (*Config, error) {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = "configs/development.yaml"
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	return &config, nil
}
```

### internal/service/user_service.go
```go
package service

import (
	"context"
	"errors"
	"medium-go-project/internal/model"
	"medium-go-project/internal/repository"
	"medium-go-project/pkg/auth"
	"time"
)

type UserService struct {
	userRepo    *repository.UserRepository
	productRepo *repository.ProductRepository
	orderRepo   *repository.OrderRepository
}

func NewUserService(userRepo *repository.UserRepository, productRepo *repository.ProductRepository, orderRepo *repository.OrderRepository) *UserService {
	return &UserService{
		userRepo:    userRepo,
		productRepo: productRepo,
		orderRepo:   orderRepo,
	}
}

func (s *UserService) CreateUser(ctx context.Context, user *model.User) error {
	// Validate user
	if err := user.Validate(); err != nil {
		return err
	}

	// Hash password
	hashedPassword, err := auth.HashPassword(user.Password)
	if err != nil {
		return err
	}

	user.Password = hashedPassword
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	return s.userRepo.Create(ctx, user)
}

func (s *UserService) GetUser(ctx context.Context, id string) (*model.User, error) {
	return s.userRepo.GetById(ctx, id)
}

func (s *UserService) UpdateUser(ctx context.Context, user *model.User) error {
	// Validate user
	if err := user.Validate(); err != nil {
		return err
	}

	// Hash password if provided
	if user.Password != "" {
		hashedPassword, err := auth.HashPassword(user.Password)
		if err != nil {
			return err
		}
		user.Password = hashedPassword
	}

	user.UpdatedAt = time.Now()
	return s.userRepo.Update(ctx, user)
}

func (s *UserService) DeleteUser(ctx context.Context, id string) error {
	return s.userRepo.Delete(ctx, id)
}

func (s *UserService) ListUsers(ctx context.Context, limit, offset int) ([]*model.User, error) {
	return s.userRepo.List(ctx, limit, offset)
}
```

### internal/model/user.go
```go
package model

import (
	"errors"
	"regexp"
	"time"
)

type User struct {
	ID        string    `json:"id" db:"id"`
	Email     string    `json:"email" db:"email"`
	Password  string    `json:"-" db:"password"`
	FirstName string    `json:"first_name" db:"first_name"`
	LastName  string    `json:"last_name" db:"last_name"`
	Role      string    `json:"role" db:"role"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

func (u *User) Validate() error {
	if u.Email == "" {
		return errors.New("email is required")
	}

	if u.Password == "" {
		return errors.New("password is required")
	}

	if u.FirstName == "" {
		return errors.New("first name is required")
	}

	if u.LastName == "" {
		return errors.New("last name is required")
	}

	if u.Role == "" {
		return errors.New("role is required")
	}

	// Validate email format
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(u.Email) {
		return errors.New("invalid email format")
	}

	return nil
}
```

### pkg/database/connection.go
```go
package database

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"medium-go-project/internal/config"

	_ "github.com/lib/pq"
)

type Database struct {
	db *sql.DB
}

func New(cfg *config.DatabaseConfig) (*Database, error) {
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Database)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Set connection pool parameters
	db.SetMaxOpenConns(cfg.MaxConnections)
	db.SetMaxIdleConns(cfg.MaxConnections / 2)
	db.SetConnMaxLifetime(time.Hour)

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &Database{db: db}, nil
}

func (d *Database) GetDB() *sql.DB {
	return d.db
}

func (d *Database) Close() error {
	return d.db.Close()
}
```

---

## Testing Scenarios

### 1. Multi-Component Session Tracking
- Initialize caf in the repository
- Make changes across different components (API, service, repository)
- Verify session tracking captures all changes correctly

### 2. Hook Integration Testing
- Test pre-session, post-edit, and pre-commit hooks
- Verify hook execution timing and success rate
- Test with different file types and sizes

### 3. Performance Testing
- Measure session creation time with medium complexity
- Test hook execution performance with multiple files
- Verify memory usage is within limits

### 4. Error Handling and Recovery
- Simulate various error scenarios
- Test graceful error handling
- Verify error logging and recovery

### 5. Multi-Environment Testing
- Test development, production, and test configurations
- Verify environment-specific behavior
- Test configuration changes and their impact

---

## Expected Results

- **Session Tracking**: 100% success rate with complete metadata across all components
- **Hook Integration**: 100% success rate with <1s execution time
- **Performance**: <200ms session creation, <1s hook execution
- **Error Handling**: Graceful recovery from all error scenarios
- **Multi-Environment**: Consistent behavior across all environments

---

## Setup Instructions

1. Navigate to the repository directory:
   ```bash
   cd repositories/medium-go-project
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
   # Make changes to trigger sessions
   echo "// New feature" > internal/api/handlers/user_handler.go
   git add .
   git commit -m "Add new feature"
   ```

5. Monitor results:
   ```bash
   caf doctor
   caf dashboard
   ```

---

## Advanced Testing

### 1. Load Testing
- Simulate multiple concurrent sessions
- Test performance under load
- Verify session tracking accuracy

### 2. Cross-Component Integration
- Test changes spanning multiple components
- Verify session tracking captures dependencies
- Test complex refactoring scenarios

### 3. Configuration Changes
- Test configuration file changes
- Verify session tracking for config updates
- Test environment-specific behavior