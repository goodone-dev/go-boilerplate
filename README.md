# Go Boilerplate
[![Go Report Card](https://goreportcard.com/badge/github.com/goodone-dev/go-boilerplate)](https://goreportcard.com/report/github.com/goodone-dev/go-boilerplate)
[![Go Quality Gate](https://github.com/goodone-dev/go-boilerplate/actions/workflows/quality_gate.yml/badge.svg)](https://github.com/goodone-dev/go-boilerplate/actions/workflows/quality_gate.yml)
[![Active Development](https://img.shields.io/badge/Maintenance%20Level-Actively%20Developed-brightgreen.svg)](https://gist.github.com/cheerfulstoic/d107229326a01ff0f333a1d3476e068d)

This Go RESTful API Boilerplate is engineered to provide a robust, scalable, and production-grade foundation for your next web service. It embraces a clean, Domain-Driven Design (DDD) architecture to ensure maintainability and separation of concerns, empowering you to focus on delivering business value instead of wrestling with infrastructure setup.

<!-- ## TODO: 💡 Motivation -->

## 🌟 Features
- 🏗️ **Clean Architecture**: Separates concerns into distinct layers (domain, application, infrastructure, presentation) for a more organized, testable, and maintainable codebase.
- 🌐 **RESTful API**: A lightweight and high-performance RESTful API built with Gin, a popular Go web framework. Includes CORS and HTTP Security middleware.
- 🔄 **Live Reload**: Automatically restart the application when file changes are detected.
- 🗃️ **Multiple Database Support**: Supports PostgreSQL, MySQL, and MongoDB. Uses a repository pattern for flexible data management.
- 🌱 **Database Migration & Seeding**: Manage your database schema and seed data with simple `make` commands.
- ⚡ **Multiple Cache Support**: Easily connect to Redis or an in-memory cache.
- 🧩 **Dependency Injection**: Switch between database or cache implementations without altering business logic.
- 🛠️ **Code Generation**: Automatically generate repository, usecase, and delivery handler with a single `make generate` command.
- 📈 **Observability**: Observability features include distributed tracing, metrics, and logging.
- 🏁 **Health Check**: `/health` endpoint for liveness and readiness probes.
- ✅ **Request Validation**: Validates incoming HTTP requests using struct tags to ensure data integrity.
- 🧹 **Request Sanitization**: Sanitizes incoming request data based on struct tags to prevent XSS and other injection attacks.
- ⏱️ **Context Propagation**: Manages request lifecycles with Go's `context` to handle cancellations and timeouts gracefully.
- 🔄 **Idempotency Handler**: Prevents duplicate requests by using a distributed cache, ensuring an operation is processed only once.
- 🚦 **Rate Limiting**: A distributed rate-limiting middleware to protect your API from excessive traffic and abuse.
- 🔌 **Circuit Breaker**: Enhances application stability by preventing repeated calls to failing external services.
- 📦 **Standardized Response**: Consistent JSON response format across all API endpoints, making it easier for clients to parse and handle responses uniformly.
- ✉️ **Email Sending**: Includes a mail sender service with support for HTML templates, allowing for easy and dynamic email generation.
- 🕒 **Background Job Processing**: Efficiently handle long-running or resource-intensive tasks asynchronously, ensuring responsive API performance and better user experience.
- 🎭 **Mock Generation**: Easily generate mocks for interfaces using the `make mock` command, simplifying unit testing.
- 🌙 **Graceful Shutdown**: Ensures that the server shuts down gracefully, finishing all in-flight requests and cleaning up resources before exiting.
- 🐳 **Dockerized Environment**: Comes with `Dockerfile` and `docker-compose.yml` for a consistent and easy-to-set-up local development environment.
- 🔒 **Pre-Commit Hooks**: Automated git hooks that run code quality checks, including linting, formatting, and security scanning before each commit.
- 🛡️ **Quality Gate CI/CD**: Automated quality checks in the CI/CD pipeline that enforce code quality standards, test coverage requirements, and security scans before deployment.

📌 **[Project Roadmap](https://github.com/users/goodone-dev/projects/3/views/7)** - Track our development progress, upcoming features, and planned improvements on our public roadmap.

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

### 2. Project Setup
Run the following command to prepare your development environment. This will make all necessary shell scripts executable:

```bash
make setup
```

To see all available make commands and their descriptions, run:

```bash
make help
```

### 3. Running the Application
You can run the application in two ways:

#### Option 1: With Docker (Recommended)
This is the easiest way to get started, as it handles all services (database, cache, etc.) for you.

1.  **Start the services**:
    ```bash
    make up
    ```
    This command builds and starts the application, database, and other services. The API by default will be accessible at `http://localhost:8080`.

2. **Stop the services**:
    ```bash
    make down
    ```
    This command stops all services.

#### Option 2: Locally
This method requires you to run the database and other services on your local machine.

1. **Setup environment variables**:
    ```bash
    cp .env.example .env
    ```
    Update `.env` with your configuration. For local development, ensure it points to your local database and other services.

2.  **Run database migrations**:
    ```bash
    make migration_up DRIVER=postgres
    ```

3.  **(Optional) Seed the database**:
    ```bash
    make seeder_up DRIVER=postgres
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
├── .dev/                       # Local development tools, scripts, and configurations.
│   └── script/                 # Local development scripts.
├── .github/                    # GitHub-specific configurations including Actions workflows and issue templates.
│   └── workflow/               # GitHub Actions workflows.
├── cmd/                        # Server commands.
│   ├── api/                    # API server.
│   │   └── main.go             # Entry point of the application. Initializes and starts the server.
│   └── utils/                  # Utility functions shared across the server.
├── internal/                   # Internal packages.
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
│   │   ├── logger/             # Log aggregation implementations.
│   │   ├── mail/               # Email sending implementation.
│   │   ├── message/            # Message bus/broker implementation.
│   │   ├── tracer/             # Distributed tracing implementation.
│   │   └── ...
│   ├── presentation/           # Adapters for incoming requests.
│   │   ├── rest/               # REST API handlers, router, and middleware.
│   │   │   ├── middleware/     # REST API middleware.
│   │   │   └── router/         # REST API router setup.
│   │   └── messaging/          # Message bus/broker handlers.
│   │       ├── middleware/     # Messaging middleware.
│   │       └── listener/       # Message bus/broker listener.
│   └── utils/                  # Utility functions shared across the application.
│       ├── breaker/            # Circuit breaker utilities.
│       ├── html/               # HTML template utilities.
│       ├── http_client/        # HTTP client utilities.
│       ├── http_response/      # HTTP response utilities.
│       ├── sanitizer/          # Request sanitizer utilities.
│       ├── validator/          # Request validation utilities.
│       └── ...
├── migrations/                 # SQL migration files for managing database schema changes.
│   ├── <database_name>/        # Migration files for a specific database.
│   └── ...
├── seeders/                    # SQL seed files for populating the database with initial data.
│   ├── <database_name>/        # Seeder files for a specific database.
│   └── ...
├── templates/                  # HTML templates for emails, PDFs, etc.
│   ├── email/                  # Email templates.
│   ├── pdf/                    # PDF templates.
│   └── ...
├── .env.example                # Example environment variables file.
├── .air.toml                   # Air.toml for local development.
├── .mockery.yml                # Mockery configuration file.
├── .pre-commit-config.yaml     # Pre-commit configuration file.
├── Makefile                    # Makefile with shortcuts for common development commands.
├── Dockerfile                  # Dockerfile for building the application image.
└── docker-compose.yml          # Defines services for the local Docker environment.
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
| **Observability**     | [opentelemetry](https://opentelemetry.io/)                                                                            |
| **Email**             | [gomail](https://github.com/go-gomail/gomail)                                                                         |
| **Circuit Breaker**   | [gobreaker](https://github.com/sony/gobreaker)                                                                        |
| **Mocking**           | [mockery](https://github.com/vektra/mockery)                                                                          |

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
- Web: [goodone.dev](https://www.goodone.dev)
- LinkedIn: [linkedin.com/in/bagusak95](https://linkedin.com/in/bagusak95)
