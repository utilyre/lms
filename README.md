# LMS

The Library Management System (LMS) is a backend API designed to provide the 
core functionalities for managing a library's book collection and user 
interactions.

## Development

### Prerequisites

- [Docker Engine](https://docs.docker.com/engine)

- [Go](https://go.dev)

- [Goose](https://pressly.github.io/goose)

  ```bash
  go install github.com/pressly/goose/v3/cmd/goose@latest
  ```

### Steps

1. Create a .env file at the root of project:

   ```bash
   DB_USER=admin
   DB_PASS=secret
   DB_PORT=5432
   DB_NAME=lms

   BE_PORT=8080
   BE_JWT_SECRET=dontshare
   ```

2. Spin up all services:

   ```bash
   docker compose up
   ```

3. Migrate:

   ```bash
   goose -dir=migrations postgres [DSN] up
   ```
