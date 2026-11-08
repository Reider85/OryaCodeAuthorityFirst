# Large Go Project for Testing
**Repository Size**: 1000+ files  
**Purpose**: Testing caf functionality with large enterprise-scale codebase and complex scenarios

---

## Repository Structure

```
large-go-project/
├── cmd/
│   ├── api-gateway/
│   │   └── main.go
│   ├── user-service/
│   │   └── main.go
│   ├── product-service/
│   │   └── main.go
│   ├── order-service/
│   │   └── main.go
│   ├── payment-service/
│   │   └── main.go
│   ├── notification-service/
│   │   └── main.go
│   ├── reporting-service/
│   │   └── main.go
│   ├── migration-service/
│   │   └── main.go
│   └── admin-service/
│       └── main.go
├── internal/
│   ├── api-gateway/
│   │   ├── handlers/
│   │   │   ├── auth_handler.go
│   │   │   ├── rate_limiter.go
│   │   │   ├── circuit_breaker.go
│   │   │   └── proxy_handler.go
│   │   ├── middleware/
│   │   │   ├── authentication.go
│   │   │   ├── authorization.go
│   │   │   ├── logging.go
│   │   │   ├── metrics.go
│   │   │   ├── tracing.go
│   │   │   └── cors.go
│   │   ├── routes/
│   │   │   ├── auth_routes.go
│   │   │   ├── user_routes.go
│   │   │   ├── product_routes.go
│   │   │   ├── order_routes.go
│   │   │   ├── payment_routes.go
│   │   │   └── admin_routes.go
│   │   ├── config/
│   │   │   ├── config.go
│   │   │   ├── routes.go
│   │   │   └── load_balancer.go
│   │   └── discovery/
│   │       ├── service_discovery.go
│   │       └── health_check.go
│   ├── user-service/
│   │   ├── handlers/
│   │   │   ├── user_handler.go
│   │   │   ├── profile_handler.go
│   │   │   ├── permission_handler.go
│   │   │   └── audit_handler.go
│   │   ├── services/
│   │   │   ├── user_service.go
│   │   │   ├── profile_service.go
│   │   │   ├── permission_service.go
│   │   │   └── audit_service.go
│   │   ├── repositories/
│   │   │   ├── user_repository.go
│   │   │   ├── profile_repository.go
│   │   │   ├── permission_repository.go
│   │   │   └── audit_repository.go
│   │   ├── models/
│   │   │   ├── user.go
│   │   │   ├── profile.go
│   │   │   ├── permission.go
│   │   │   └── audit.go
│   │   ├── cache/
│   │   │   ├── user_cache.go
│   │   │   └── profile_cache.go
│   │   ├── events/
│   │   │   ├── user_events.go
│   │   │   └── event_handlers.go
│   │   └── utils/
│   │       ├── validation.go
│   │       ├── crypto.go
│   │       ├── email.go
│   │       └── file.go
│   ├── product-service/
│   │   ├── handlers/
│   │   │   ├── product_handler.go
│   │   │   ├── category_handler.go
│   │   │   ├── inventory_handler.go
│   │   │   └── pricing_handler.go
│   │   ├── services/
│   │   │   ├── product_service.go
│   │   │   ├── category_service.go
│   │   │   ├── inventory_service.go
│   │   │   └── pricing_service.go
│   │   ├── repositories/
│   │   │   ├── product_repository.go
│   │   │   ├── category_repository.go
│   │   │   ├── inventory_repository.go
│   │   │   └── pricing_repository.go
│   │   ├── models/
│   │   │   ├── product.go
│   │   │   ├── category.go
│   │   │   ├── inventory.go
│   │   │   └── pricing.go
│   │   ├── cache/
│   │   │   ├── product_cache.go
│   │   │   └── category_cache.go
│   │   ├── events/
│   │   │   ├── product_events.go
│   │   │   └── event_handlers.go
│   │   └── utils/
│   │       ├── validation.go
│   │       ├── search.go
│   │       ├── image.go
│   │       └── export.go
│   ├── order-service/
│   │   ├── handlers/
│   │   │   ├── order_handler.go
│   │   │   ├── cart_handler.go
│   │   │   ├── fulfillment_handler.go
│   │   │   └── returns_handler.go
│   │   ├── services/
│   │   │   ├── order_service.go
│   │   │   ├── cart_service.go
│   │   │   ├── fulfillment_service.go
│   │   │   └── returns_service.go
│   │   ├── repositories/
│   │   │   ├── order_repository.go
│   │   │   ├── cart_repository.go
│   │   │   ├── fulfillment_repository.go
│   │   │   └── returns_repository.go
│   │   ├── models/
│   │   │   ├── order.go
│   │   │   ├── cart.go
│   │   │   ├── fulfillment.go
│   │   │   └── returns.go
│   │   ├── cache/
│   │   │   ├── order_cache.go
│   │   │   └── cart_cache.go
│   │   ├── events/
│   │   │   ├── order_events.go
│   │   │   └── event_handlers.go
│   │   └── utils/
│   │       ├── validation.go
│   │       ├── calculation.go
│   │       ├── workflow.go
│   │       └── export.go
│   ├── payment-service/
│   │   ├── handlers/
│   │   │   ├── payment_handler.go
│   │   │   ├── refund_handler.go
│   │   │   └── reconciliation_handler.go
│   │   ├── services/
│   │   │   ├── payment_service.go
│   │   │   ├── refund_service.go
│   │   │   └── reconciliation_service.go
│   │   ├── repositories/
│   │   │   ├── payment_repository.go
│   │   │   ├── refund_repository.go
│   │   │   └── reconciliation_repository.go
│   │   ├── models/
│   │   │   ├── payment.go
│   │   │   ├── refund.go
│   │   │   └── reconciliation.go
│   │   ├── integrations/
│   │   │   ├── stripe.go
│   │   │   ├── paypal.go
│   │   │   └── bank_transfer.go
│   │   └── utils/
│   │       ├── validation.go
│   │       ├── encryption.go
│   │       └── reporting.go
│   ├── notification-service/
│   │   ├── handlers/
│   │   │   ├── email_handler.go
│   │   │   ├── sms_handler.go
│   │   │   ├── push_handler.go
│   │   │   └── webhook_handler.go
│   │   ├── services/
│   │   │   ├── email_service.go
│   │   │   ├── sms_service.go
│   │   │   ├── push_service.go
│   │   │   └── webhook_service.go
│   │   ├── repositories/
│   │   │   ├── email_repository.go
│   │   │   ├── sms_repository.go
│   │   │   ├── push_repository.go
│   │   │   └── webhook_repository.go
│   │   ├── models/
│   │   │   ├── email.go
│   │   │   ├── sms.go
│   │   │   ├── push.go
│   │   │   └── webhook.go
│   │   ├── templates/
│   │   │   ├── email/
│   │   │   │   ├── welcome.html
│   │   │   │   ├── order.html
│   │   │   │   └── reset.html
│   │   │   ├── sms/
│   │   │   │   ├── welcome.txt
│   │   │   │   └── order.txt
│   │   │   └── push/
│   │   │       ├── welcome.json
│   │   │       └── order.json
│   │   └── utils/
│   │       ├── validation.go
│   │       ├── templating.go
│   │       └── delivery.go
│   ├── reporting-service/
│   │   ├── handlers/
│   │   │   ├── sales_handler.go
│   │   │   ├── inventory_handler.go
│   │   │   ├── financial_handler.go
│   │   │   └── user_handler.go
│   │   ├── services/
│   │   │   ├── sales_service.go
│   │   │   ├── inventory_service.go
│   │   │   ├── financial_service.go
│   │   │   └── user_service.go
│   │   ├── repositories/
│   │   │   ├── sales_repository.go
│   │   │   ├── inventory_repository.go
│   │   │   ├── financial_repository.go
│   │   │   └── user_repository.go
│   │   ├── models/
│   │   │   ├── sales.go
│   │   │   ├── inventory.go
│   │   │   ├── financial.go
│   │   │   └── user.go
│   │   ├── exporters/
│   │   │   ├── csv.go
│   │   │   ├── excel.go
│   │   │   └── pdf.go
│   │   └── utils/
│   │       ├── validation.go
│   │       ├── aggregation.go
│   │       └── formatting.go
│   ├── migration-service/
│   │   ├── handlers/
│   │   │   ├── migration_handler.go
│   │   │   ├── rollback_handler.go
│   │   │   └── validation_handler.go
│   │   ├── services/
│   │   │   ├── migration_service.go
│   │   │   ├── rollback_service.go
│   │   │   └── validation_service.go
│   │   ├── repositories/
│   │   │   ├── migration_repository.go
│   │   │   ├── rollback_repository.go
│   │   │   └── validation_repository.go
│   │   ├── models/
│   │   │   ├── migration.go
│   │   │   ├── rollback.go
│   │   │   └── validation.go
│   │   └── utils/
│   │       ├── validation.go
│   │       ├── schema.go
│   │       └── data.go
│   ├── admin-service/
│   │   ├── handlers/
│   │   │   ├── user_admin_handler.go
│   │   │   ├── product_admin_handler.go
│   │   │   ├── order_admin_handler.go
│   │   │   └── system_admin_handler.go
│   │   ├── services/
│   │   │   ├── user_admin_service.go
│   │   │   ├── product_admin_service.go
│   │   │   ├── order_admin_service.go
│   │   │   └── system_admin_service.go
│   │   ├── repositories/
│   │   │   ├── user_admin_repository.go
│   │   │   ├── product_admin_repository.go
│   │   │   ├── order_admin_repository.go
│   │   │   └── system_admin_repository.go
│   │   ├── models/
│   │   │   ├── user_admin.go
│   │   │   ├── product_admin.go
│   │   │   ├── order_admin.go
│   │   │   └── system_admin.go
│   │   └── utils/
│   │       ├── validation.go
│   │       ├── audit.go
│   │       └── backup.go
│   ├── shared/
│   │   ├── config/
│   │   │   ├── config.go
│   │   │   ├── database.go
│   │   │   ├── redis.go
│   │   │   ├── logging.go
│   │   │   ├── metrics.go
│   │   │   └── tracing.go
│   │   ├── database/
│   │   │   ├── connection.go
│   │   │   ├── migrations/
│   │   │   │   ├── 001_initial_schema.sql
│   │   │   │   ├── 002_add_indexes.sql
│   │   │   │   ├── 003_add_foreign_keys.sql
│   │   │   │   ├── 004_add_constraints.sql
│   │   │   │   └── 005_add_partitions.sql
│   │   │   ├── queries/
│   │   │   │   ├── user_queries.go
│   │   │   │   ├── product_queries.go
│   │   │   │   ├── order_queries.go
│   │   │   │   └── payment_queries.go
│   │   │   └── repositories/
│   │   │       ├── base_repository.go
│   │   │       ├── user_repository.go
│   │   │       ├── product_repository.go
│   │   │       └── order_repository.go
│   │   ├── cache/
│   │   │   ├── redis.go
│   │   │   ├── cache.go
│   │   │   └── distributed_lock.go
│   │   ├── auth/
│   │   │   ├── jwt.go
│   │   │   ├── rbac.go
│   │   │   └── oauth.go
│   │   ├── events/
│   │   │   ├── event_bus.go
│   │   │   ├── event_handlers.go
│   │   │   └── event_store.go
│   │   ├── utils/
│   │   │   ├── validation.go
│   │   │   ├── crypto.go
│   │   │   ├── email.go
│   │   │   ├── file.go
│   │   │   ├── string.go
│   │   │   ├── time.go
│   │   │   └── math.go
│   │   ├── middleware/
│   │   │   ├── authentication.go
│   │   │   ├── authorization.go
│   │   │   ├── logging.go
│   │   │   ├── metrics.go
│   │   │   ├── tracing.go
│   │   │   ├── cors.go
│   │   │   ├── rate_limiting.go
│   │   │   └── circuit_breaker.go
│   │   └── monitoring/
│   │       ├── metrics.go
│   │       ├── health.go
│   │       ├── tracing.go
│   │       └── alerting.go
│   └── migrations/
│       ├── 001_initial_schema.sql
│       ├── 002_add_indexes.sql
│       ├── 003_add_foreign_keys.sql
│       ├── 004_add_constraints.sql
│       ├── 005_add_partitions.sql
│       ├── 006_add_audit_triggers.sql
│       ├── 007_add_event_sourcing.sql
│       ├── 008_add_cdc.sql
│       ├── 009_add_data_retention.sql
│       └── 010_add_performance_optimization.sql
├── pkg/
│   ├── database/
│   │   ├── connection.go
│   │   ├── migrations/
│   │   │   ├── runner.go
│   │   │   └── version.go
│   │   ├── queries/
│   │   │   ├── user_queries.go
│   │   │   ├── product_queries.go
│   │   │   ├── order_queries.go
│   │   │   └── payment_queries.go
│   │   └── repositories/
│   │       ├── base_repository.go
│   │       ├── user_repository.go
│   │       ├── product_repository.go
│   │       └── order_repository.go
│   ├── cache/
│   │   ├── redis.go
│   │   ├── cache.go
│   │   └── distributed_lock.go
│   ├── auth/
│   │   ├── jwt.go
│   │   ├── rbac.go
│   │   └── oauth.go
│   ├── events/
│   │   ├── event_bus.go
│   │   ├── event_handlers.go
│   │   └── event_store.go
│   ├── monitoring/
│   │   ├── metrics.go
│   │   ├── health.go
│   │   ├── tracing.go
│   │   └── alerting.go
│   ├── utils/
│   │   ├── validation.go
│   │   ├── crypto.go
│   │   ├── email.go
│   │   ├── file.go
│   │   ├── string.go
│   │   ├── time.go
│   │   └── math.go
│   ├── middleware/
│   │   ├── authentication.go
│   │   ├── authorization.go
│   │   ├── logging.go
│   │   ├── metrics.go
│   │   ├── tracing.go
│   │   ├── cors.go
│   │   ├── rate_limiting.go
│   │   └── circuit_breaker.go
│   └── exporters/
│       ├── csv.go
│       ├── excel.go
│       ├── pdf.go
│       └── json.go
├── configs/
│   ├── development.yaml
│   ├── staging.yaml
│   ├── production.yaml
│   ├── test.yaml
│   └── local.yaml
├── scripts/
│   ├── deploy.sh
│   ├── migrate.sh
│   ├── test.sh
│   ├── backup.sh
│   ├── restore.sh
│   ├── health-check.sh
│   └── performance-test.sh
├── docs/
│   ├── architecture/
│   │   ├── microservices.md
│   │   ├── api-gateway.md
│   │   ├── database.md
│   │   ├── cache.md
│   │   ├── auth.md
│   │   └── monitoring.md
│   ├── api/
│   │   ├── user-api.md
│   │   ├── product-api.md
│   │   ├── order-api.md
│   │   └── payment-api.md
│   ├── deployment/
│   │   ├── docker.md
│   │   ├── kubernetes.md
│   │   └── cloud.md
│   ├── operations/
│   │   ├── monitoring.md
│   │   ├── logging.md
│   │   ├── backup.md
│   │   └── disaster-recovery.md
│   └── development/
│       ├── coding-standards.md
│       ├── testing.md
│       ├── debugging.md
│       └── performance.md
├── test/
│   ├── unit/
│   │   ├── user-service/
│   │   │   ├── user_service_test.go
│   │   │   ├── profile_service_test.go
│   │   │   └── permission_service_test.go
│   │   ├── product-service/
│   │   │   ├── product_service_test.go
│   │   │   ├── category_service_test.go
│   │   │   └── inventory_service_test.go
│   │   ├── order-service/
│   │   │   ├── order_service_test.go
│   │   │   ├── cart_service_test.go
│   │   │   └── fulfillment_service_test.go
│   │   └── shared/
│   │       ├── database_test.go
│   │       ├── cache_test.go
│   │       ├── auth_test.go
│   │       └── utils_test.go
│   ├── integration/
│   │   ├── api-gateway/
│   │   │   ├── api_gateway_test.go
│   │   │   └── proxy_test.go
│   │   ├── user-service/
│   │   │   ├── user_integration_test.go
│   │   │   └── profile_integration_test.go
│   │   ├── product-service/
│   │   │   ├── product_integration_test.go
│   │   │   └── category_integration_test.go
│   │   └── order-service/
│   │       ├── order_integration_test.go
│   │       └── cart_integration_test.go
│   ├── e2e/
│   │   ├── user-flow/
│   │   │   ├── user_registration_test.go
│   │   │   ├── user_login_test.go
│   │   │   └── user_profile_test.go
│   │   ├── order-flow/
│   │   │   ├── order_creation_test.go
│   │   │   ├── order_payment_test.go
│   │   │   └── order_fulfillment_test.go
│   │   └── admin-flow/
│   │       ├── admin_dashboard_test.go
│   │       ├── user_management_test.go
│   │       └── reporting_test.go
│   └── performance/
│       ├── load-test/
│       │   ├── user_service_load_test.go
│       │   ├── product_service_load_test.go
│       │   └── order_service_load_test.go
│       ├── stress-test/
│       │   ├── user_service_stress_test.go
│       │   ├── product_service_stress_test.go
│       │   └── order_service_stress_test.go
│       └── benchmark/
│           ├── database_benchmark.go
│           ├── cache_benchmark.go
│           └── api_benchmark.go
├── docker/
│   ├── Dockerfile
│   ├── docker-compose.yml
│   ├── docker-compose.dev.yml
│   ├── docker-compose.prod.yml
│   ├── docker-compose.test.yml
│   └── k8s/
│       ├── namespace.yaml
│       ├── configmap.yaml
│       ├── secret.yaml
│       ├── deployment.yaml
│       ├── service.yaml
│       ├── ingress.yaml
│       └── hpa.yaml
├── terraform/
│   ├── main.tf
│   ├── variables.tf
│   ├── outputs.tf
│   ├── providers.tf
│   └── modules/
│       ├── vpc.tf
│       ├── ec2.tf
│       ├── rds.tf
│       ├── redis.tf
│       └── s3.tf
├── ci-cd/
│   ├── .github/
│   │   ├── workflows/
│   │   │   ├── ci.yml
│   │   │   ├── cd.yml
│   │   │   ├── test.yml
│   │   │   ├── security.yml
│   │   │   └── performance.yml
│   │   └── actions/
│   │       ├── setup-go.yml
│   │       ├── run-tests.yml
│   │       ├── build-docker.yml
│   │       └── deploy.yml
│   └── gitlab-ci/
│       ├── .gitlab-ci.yml
│       ├── build.yml
│       ├── test.yml
│       ├── deploy.yml
│       └── monitor.yml
├── .gitignore
├── go.mod
├── go.sum
├── Makefile
├── README.md
├── LICENSE
├── CHANGELOG.md
└── CONTRIBUTING.md
```

---

## Key Features

### 1. Microservices Architecture
- **API Gateway**: Central entry point with routing, load balancing, and rate limiting
- **User Service**: User management, authentication, and authorization
- **Product Service**: Product catalog, inventory, and pricing
- **Order Service**: Order processing, cart management, and fulfillment
- **Payment Service**: Payment processing, refunds, and reconciliation
- **Notification Service**: Email, SMS, push notifications, and webhooks
- **Reporting Service**: Sales, inventory, financial, and user analytics
- **Migration Service**: Database migrations, rollbacks, and validation
- **Admin Service**: Administrative functions and system management

### 2. Advanced Infrastructure
- **Service Discovery**: Dynamic service registration and discovery
- **Load Balancing**: Multiple load balancing strategies
- **Circuit Breaker**: Fault tolerance and system resilience
- **Distributed Tracing**: Request tracing across services
- **Event Sourcing**: Event-driven architecture with CQRS
- **Caching Layer**: Multi-level caching with Redis
- **Database Sharding**: Horizontal scaling with partitioning
- **CDC (Change Data Capture)**: Real-time data synchronization

### 3. Enterprise Features
- **Multi-tenancy**: Support for multiple organizations
- **Audit Logging**: Comprehensive audit trails
- **Data Retention**: Automated data lifecycle management
- **Performance Optimization**: Database indexing, query optimization
- **Security**: JWT authentication, RBAC, OAuth2, encryption
- **Monitoring**: Comprehensive metrics, health checks, and alerting
- **Backup/Recovery**: Automated backup and disaster recovery

### 4. DevOps Integration
- **Containerization**: Docker and Kubernetes deployment
- **Infrastructure as Code**: Terraform for infrastructure management
- **CI/CD**: GitHub Actions and GitLab CI/CD pipelines
- **Environment Management**: Development, staging, production environments
- **Security Scanning**: Automated security vulnerability scanning

---

## Testing Scenarios

### 1. Large-Scale Session Tracking
- Initialize caf in the repository
- Make changes across multiple services and components
- Verify session tracking captures all changes correctly
- Test performance with large number of files

### 2. Multi-Service Hook Integration
- Test pre-session, post-edit, and pre-commit hooks across all services
- Verify hook execution timing and success rate
- Test with complex file dependencies and interactions

### 3. Performance Testing
- Measure session creation time with large codebase
- Test hook execution performance with thousands of files
- Verify memory usage and resource consumption
- Test under load conditions

### 4. Error Handling and Recovery
- Simulate various error scenarios across services
- Test graceful error handling and recovery
- Verify error logging and monitoring
- Test system resilience under failure conditions

### 5. Multi-Environment Testing
- Test development, staging, and production configurations
- Verify environment-specific behavior
- Test configuration changes and their impact
- Validate deployment and migration processes

### 6. Load and Stress Testing
- Simulate high-traffic scenarios
- Test system performance under load
- Verify session tracking accuracy under stress
- Test database performance and caching

### 7. Security and Compliance Testing
- Test authentication and authorization
- Verify audit logging and compliance
- Test data encryption and security
- Validate access controls and permissions

---

## Expected Results

- **Session Tracking**: 100% success rate with complete metadata across all services
- **Hook Integration**: 100% success rate with <2s execution time for large repositories
- **Performance**: <500ms session creation, <2s hook execution for 1000+ files
- **Error Handling**: Graceful recovery from all error scenarios
- **Multi-Environment**: Consistent behavior across all environments
- **Load Testing**: System performance meets SLA requirements
- **Security**: All security requirements met with proper audit trails

---

## Setup Instructions

1. Navigate to the repository directory:
   ```bash
   cd repositories/large-go-project
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
   echo "// New feature" > internal/user-service/handlers/user_handler.go
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

### 1. Microservices Integration Testing
- Test changes spanning multiple microservices
- Verify session tracking captures service interactions
- Test API gateway routing and load balancing
- Validate service discovery and health checks

### 2. Database Performance Testing
- Test database queries and performance
- Verify indexing and optimization strategies
- Test data migration and schema changes
- Validate backup and recovery processes

### 3. Load Testing
- Simulate high-traffic scenarios
- Test system performance under load
- Verify session tracking accuracy under stress
- Test database performance and caching

### 4. Security Testing
- Test authentication and authorization
- Verify audit logging and compliance
- Test data encryption and security
- Validate access controls and permissions

### 5. DevOps Pipeline Testing
- Test CI/CD pipeline functionality
- Verify deployment processes
- Test infrastructure as code
- Validate monitoring and alerting

---

## Monitoring and Observability

### 1. Metrics Collection
- Application metrics (response time, throughput, error rate)
- System metrics (CPU, memory, disk, network)
- Business metrics (user activity, order volume, revenue)
- Custom metrics for specific business logic

### 2. Logging and Tracing
- Structured logging with correlation IDs
- Distributed tracing across services
- Error tracking and alerting
- Performance monitoring and optimization

### 3. Health Checks
- Service health monitoring
- Database connectivity checks
- Cache performance monitoring
- External service dependency checks

### 4. Alerting and Notification
- Real-time alerting for critical issues
- Notification escalation procedures
- Automated incident response
- Performance degradation alerts