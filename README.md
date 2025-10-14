# Go Boilerplate
[![Go Report Card](https://goreportcard.com/badge/github.com/goodone-dev/go-boilerplate)](https://goreportcard.com/report/github.com/goodone-dev/go-boilerplate)

This Go RESTful API Boilerplate is engineered to provide a robust, scalable, and production-grade foundation for your next web service. It embraces a clean, Domain-Driven Design (DDD) architecture to ensure maintainability and separation of concerns, empowering you to focus on delivering business value instead of wrestling with infrastructure setup.

<!-- ## TODO: 💡 Motivation -->

## 🌟 Features
- 🏗️ **Clean Architecture**: Separates concerns into distinct layers (domain, application, infrastructure, presentation) for a more organized, testable, and maintainable codebase.
- 🌐 **RESTful API**: A lightweight and high-performance RESTful API built with Gin, a popular Go web framework.
- 🔄 **Live Reload**: Automatically restart the application when file changes are detected.
- 🗃️ **Multiple Database Support**: Supports PostgreSQL, MySQL, and MongoDB. Uses a repository pattern for flexible data management.
- 🌱 **Database Migration & Seeding**: Manage your database schema and seed data with simple `make` commands.
- ⚡ **Multiple Cache Support**: Easily connect to Redis or an in-memory cache.
- 💉 **Dependency Injection**: Switch between database or cache implementations without altering business logic.
- 🔍 **Distributed Tracing**: Integrated with Jaeger for distributed tracing, offering insights into request flows across services to simplify debugging and performance monitoring.
- ✅ **Request Validation**: Validates incoming HTTP requests using struct tags to ensure data integrity.
- 🧹 **Request Sanitization**: Sanitizes incoming request data based on struct tags to prevent XSS and other injection attacks.
- 🔗 **Context Propagation**: Manages request lifecycles with Go's `context` to handle cancellations and timeouts gracefully.
- 🛡️ **Idempotency Middleware**: Prevents duplicate requests by using a distributed cache, ensuring an operation is processed only once.
- 🚦 **Rate Limiting**: A distributed rate-limiting middleware to protect your API from excessive traffic and abuse.
- 🔌 **Circuit Breaker**: Enhances application stability by preventing repeated calls to failing external services.
- ⚠️ **Centralized Error Handling**: A centralized middleware automatically handles errors, converting them into consistent and well-formatted HTTP responses.
- 📧 **Email Sending**: Includes a mail sender service with support for HTML templates, allowing for easy and dynamic email generation.
- 🕒 **Asynchronous Processing**: Offloads long-running tasks to a message bus, ensuring non-blocking API responses.
- 🎭 **Mock Generation**: Easily generate mocks for interfaces using the `make mock` command, simplifying unit testing.
- 🌙 **Graceful Shutdown**: Ensures that the server shuts down gracefully, finishing all in-flight requests and cleaning up resources before exiting.
- 🐳 **Dockerized Environment**: Comes with `Dockerfile` and `docker-compose.yml` for a consistent and easy-to-set-up local development environment.
- 🤖 **CI/CD Pipeline**: Automated checks for building, linting, test coverage, and security scanning.

<details>
<summary><b>## 🚧 Roadmap</b></summary>

- Unit Test
    - Test all usecase
    - CI/CD Pipeline
    - Min 80% Code Coverage (configured with gocover)
- Authentication & Authorization
    - Authentication Middleware
    - Add role-based access control (RBAC)
    - Ory Kratos/Keto Integration
- CORS middleware
- Sanitize Error Message (need to test)
- HTTP Security Middleware
    - `X-Content-Type-Options: nosniff`
    - `X-Frame-Options: DENY`
    - `X-XSS-Protection: 1; mode=block`
    - `Strict-Transport-Security` (HSTS)
    - `Content-Security-Policy`
- Secrets Management
    - Use Vault
- Input Validation
    - Path traversal protection for file operations
    - Request size limits
    - Content-Type validation
- Database Connection Security
    - Add connection encryption for production
    - Implement connection pooling limits (partially done ✅)
    - Add prepared statement caching
- Structured Logging
    - `zerolog` Integration
    - Log levels (DEBUG, INFO, WARN, ERROR)
    - Include request IDs in all logs
    - Log Aggregation (e.g., ELK Stack)
- API Documentation
    - Add swaggo/swag for auto-generated API docs
- Health Check Endpoints
    - `/health` - Basic liveness probe
    - `/health/ready` - Readiness probe (check DB, Redis, etc.)
    - `/metrics` - Prometheus metrics endpoint
- Graceful Degradation
    - Implement fallback mechanisms when external services fail
    - Add retry logic with exponential backoff
    - Circuit breaker is present ✅ but not used in code
- Background Job Processing
    - Implement RabbitMQ/Kafka as mentioned in roadmap
    - Add worker command for processing async jobs
    - Implement job retry and dead-letter queues
- Caching Strategy
    - Query result caching
    - HTTP response caching
    - Cache invalidation strategies
    - Cache warming on startup

</details>

## 🚀 Getting Started

### Prerequisites
- [Go](https://golang.org/doc/install)
- [Make](https://www.gnu.org/software/make/)
- [Docker](https://docs.docker.com/get-docker/) & [Docker Compose](https://docs.docker.com/compose/install/) (for Docker-based setup)

### 1. Installation
Clone the repository:
```bash
git clone https://github.com/goodone-dev/go-boilerplate.git
cd go-boilerplate
```

### 2. Running the Application
You can run the application in two ways:

#### Option 1: With Docker (Recommended)
This is the easiest way to get started, as it handles all services (database, cache, etc.) for you.

1.  **Start the services**:
    ```bash
    make up
    ```
    This command builds and starts the application, database, and other services.

The API will be accessible at `http://localhost:8080`.

To stop all services, run `make down`.

#### Option 2: Locally
This method requires you to run the database and other services on your local machine.

1. **Setup environment variables**:
    ```bash
    cp .env.example .env
    ```
    Update `.env` with your configuration. For local development, ensure it points to your local database and other services.

2.  **Run database migrations**:
    ```bash
    make migrate_up DRIVER=postgres
    ```

3.  **(Optional) Seed the database**:
    ```bash
    make seed DRIVER=postgres
    ```

4.  **Run the application**:
    ```bash
    make run
    ```
    This command will start the Go application. The API will be accessible at `http://localhost:8080`.

## 📂 Project Structure

This project is structured following the principles of **Clean Architecture**. The code is organized into distinct layers, promoting separation of concerns, testability, and maintainability. The dependencies flow inwards, from the outer layers (Infrastructure, Presentation) to the inner layers (Application, Domain).

```
.
├── cmd/
│   └── api/
│       └── main.go             # Entry point of the application. Initializes and starts the server.
├── internal/
│   ├── application/            # Implements use cases by orchestrating domain logic.
│   │   ├── <domain_name>/      # Groups application logic for a specific domain.
│   │   │   ├── delivery/       # Adapters for handling incoming requests (e.g., HTTP, messaging).
│   │   │   │   ├── http/       # HTTP handlers for the domain.
│   │   │   │   └── messaging/  # Message handlers for the domain.
│   │   │   ├── repository/     # Repository implementations for the domain.
│   │   │   └── usecase/        # Business logic and use cases for the domain.
│   │   └── ...
│   ├── config/                 # Configuration loading and management.
│   ├── domain/                 # Contains core entities and interfaces.
│   │   ├── <domain_name>/      # Groups domain logic for a specific business entity.
│   │   │   └── mocks/          # Mocks for domain interfaces.
│   │   └── ...                 
│   ├── infrastructure/         # Provides implementations for external services.
│   │   ├── cache/              # Cache implementations (e.g., Redis).
│   │   ├── database/           # Database implementations (PostgreSQL, MySQL, MongoDB).
│   │   ├── integration/        # Clients for external APIs.
│   │   ├── mail/               # Email sending implementation.
│   │   ├── message/            # Message bus implementation.
│   │   ├── tracer/             # Distributed tracing implementation (e.g., Jaeger).
│   │   └── ...
│   ├── presentation/           # Adapters for incoming requests.
│   │   ├── rest/               # REST API handlers, router, and middleware.
│   │   │   ├── middleware/     # REST API middleware.
│   │   │   └── router/         # REST API router setup.
│   │   ├── messaging/          # Message bus handlers.
│   │   │   ├── middleware/     # Messaging middleware.
│   │   │   └── listener/       # Message bus listener.
│   │   └── ...
│   └── utils/                  # Utility functions shared across the application.
│       ├── error/              # Custom error types and handling.
│       ├── html/               # HTML template rendering utilities.
│       ├── http/               # HTTP client functions.
│       ├── success/            # Standardized success responses.
│       ├── tracer/             # Tracer helper functions.
│       ├── validator/          # Request validation utilities.
│       └── ...
├── migrations/                 # SQL migration files for managing database schema changes.
│   └── <database_name>/        # Migration files for a specific database.
├── seeders/                    # SQL seed files for populating the database with initial data.
│   └── <database_name>/        # Seeder files for a specific database.
├── templates/                  # HTML templates for emails, PDFs, etc.
│   ├── email/                  # Email templates.
│   ├── pdf/                    # PDF templates.
│   └── ...
├── Makefile                    # Makefile with shortcuts for common development commands.
├── docker-compose.yml          # Defines services for the local Docker environment.
└── Dockerfile                  # Dockerfile for building the application image.
```

<!-- ## TODO: 🏗️ Architecture Diagram -->

<!-- ## TODO: 🔧 Development -->

## 🛠️ Tech Stack

| Category              | Technologies                                                                                                          |
| --------------------- | ----------------------------------------------------------------------------------------------------------------------|
| **Framework**         | [gin](https://github.com/gin-gonic/gin)                                                                               |
| **Database**          | [gorm](https://gorm.io/) (PostgreSQL, MySQL), [mongo-driver](https://github.com/mongodb/mongo-go-driver) (MongoDB)    |
| **Cache**             | [go-redis](https://github.com/redis/go-redis)                                                                         |
| **API Client**        | [resty](https://github.com/go-resty/resty)                                                                            |
| **Config**            | [viper](https://github.com/spf13/viper)                                                                               |
| **Validation**        | [validator](https://github.com/go-playground/validator)                                                               |
| **Migration**         | [golang-migrate](https://github.com/golang-migrate/migrate)                                                           |
| **Tracing**           | [opentelemetry](https://opentelemetry.io/)                                                                            |
| **Email**             | [gomail](https://github.com/go-gomail/gomail)                                                                         |
| **Circuit Breaker**   | [gobreaker](https://github.com/sony/gobreaker)                                                                        |
| **Mocking**           | [mockery](https://github.com/vektra/mockery)                                                                          |

<!-- ## 🚧 Roadmap
- [ ] **Alerting**: Integration with Prometheus Alertmanager for handling alerts.
- [ ] **Message Broker Support**: Adding support for Kafka and RabbitMQ.
- [ ] **Authentication**: Implementing OAuth2 with Ory Kratos for identity and user management.
- [ ] **Authorization**: Integration with Ory Keto for permission and access control.
- [ ] **Structured Logging**: Implementing a structured logger (e.g., Logrus).
- [ ] **Worker Command**: Add worker for processing asynchronous task.
- [ ] **Make Generate Command**: Automate the creation of entity, repository, usecase, and handler files.
- [ ] **HTTP Security Middleware**: Add middleware for handling common security headers.
- [ ] **XSS Handling**: Add middleware for Cross-Site Scripting (XSS) protection.
- [ ] **CORS Handling**: Implement middleware for Cross-Origin Resource Sharing (CORS).
- [ ] **Auto Generate Documentation**: Automatically generate API documentation. -->

## 🤝 Contributing
1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📝 License
This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 📧 Contact
**Bagus Abdul Kurniawan**
- Email: hello@goodone.dev
- GitHub: [github.com/goodone-dev](https://github.com/goodone-dev)
- LinkedIn: [linkedin.com/in/bagusak95](https://linkedin.com/in/bagusak95)