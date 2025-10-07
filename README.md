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

## 🚧 Roadmap
- [ ] **Alerting**: Integration with Prometheus Alertmanager for handling alerts.
- [ ] **CI/CD Pipeline**: Automated checks for linting, test coverage, and security scanning.
- [ ] **Message Broker Support**: Adding support for Kafka and RabbitMQ.
- [ ] **Authentication**: Implementing OAuth2 with Ory Kratos for identity and user management.
- [ ] **Authorization**: Integration with Ory Keto for permission and access control.
- [ ] **Structured Logging**: Implementing a structured logger (e.g., Logrus).
- [ ] **Request Sanitization**: Middleware to sanitize incoming request data.
- [ ] **Worker Command**: Add worker for processing asynchronous task.
Internal Test
- [ ] Migration MongoDB & MySQL
- [ ] Seeder MongoDB & MySQL
- [ ] Implementation MongoDB & MySQL