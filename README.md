# Go Boilerplate
[![Go Report Card](https://goreportcard.com/badge/github.com/BagusAK95/go-boilerplate)](https://goreportcard.com/report/github.com/BagusAK95/go-boilerplate)

This Go RESTful API Boilerplate is engineered to provide a robust, scalable, and production-grade foundation for your next web service. It embraces a clean, Domain-Driven Design (DDD) architecture to ensure maintainability and separation of concerns, empowering you to focus on delivering business value instead of wrestling with infrastructure setup.

## 🌟 Features
- **Clean Architecture**: Separates concerns into distinct layers (domain, application, infrastructure, presentation) for a more organized, testable, and maintainable codebase.
- **RESTful API**: A lightweight and high-performance RESTful API built with Gin, a popular Go web framework.
- **Multiple Database Support**: Supports PostgreSQL, MySQL, and MongoDB. Uses a repository pattern for flexible data management.
- **Database Migration & Seeding**: Manage your database schema and seed data with simple `make` commands.
- **Multiple Cache Support**: Easily connect to Redis or an in-memory cache.
- **Distributed Tracing**: Integrated with Jaeger for distributed tracing, offering insights into request flows across services to simplify debugging and performance monitoring.
- **Request Validation**: Validates incoming HTTP requests using struct tags to ensure data integrity.
- **Context Propagation**: Manages request lifecycles with Go's `context` to handle cancellations and timeouts gracefully.
- **Idempotency Middleware**: Prevents duplicate requests by using a distributed cache, ensuring an operation is processed only once.
- **Rate Limiting**: A distributed rate-limiting middleware to protect your API from excessive traffic and abuse.
- **Circuit Breaker**: Enhances application stability by preventing repeated calls to failing external services.
- **Centralized Error Handling**: A centralized middleware automatically handles errors, converting them into consistent and well-formatted HTTP responses.
- **Email Sending**: Includes a mail sender service with support for HTML templates, allowing for easy and dynamic email generation.
- **Asynchronous Processing**: Offloads long-running tasks to a message bus, ensuring non-blocking API responses.
- **Mock Generation**: Easily generate mocks for interfaces using the `make mock` command, simplifying unit testing.
- **Graceful Shutdown**: Ensures that the server shuts down gracefully, finishing all in-flight requests and cleaning up resources before exiting.
- **Dockerized Environment**: Comes with `Dockerfile` and `docker-compose.yml` for a consistent and easy-to-set-up local development environment.

## 🚀 Quick Start

### Prerequisites
- [Go](https://golang.org/doc/install)
- [Make](https://www.gnu.org/software/make/)
- [Docker](https://docs.docker.com/get-docker/) & [Docker Compose](https://docs.docker.com/compose/install/) (for Docker-based setup)

### 1. Installation
Clone the repository:
```bash
git clone https://github.com/BagusAK95/go-boilerplate.git
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

## 🚧 Roadmap
- [ ] **Alerting**: Integration with Prometheus Alertmanager for handling alerts.
- [ ] **CI/CD Pipeline**: Automated checks for linting, test coverage, and security scanning.
- [ ] **Message Broker Support**: Adding support for Kafka and RabbitMQ.
- [ ] **Authentication**: Implementing OAuth2 with Ory Kratos for identity and user management.
- [ ] **Authorization**: Integration with Ory Keto for permission and access control.
- [ ] **Structured Logging**: Implementing a structured logger (e.g., Logrus).
- [ ] **Request Sanitization**: Middleware to sanitize incoming request data.
- [ ] **Worker Command**: Add worker for processing asynchronous task.
- [ ] **Makefile Dependency Check**: Automatically prompt to install missing tools when running a make command.

## Internal Test
- [ ] Migration MongoDB & MySQL
- [ ] Seeder MongoDB & MySQL
- [ ] Implementation MongoDB & MySQL