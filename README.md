# IOHK Golang Backend (Pre-production)

## Table of Contents
- [Introduction](#introduction)
- [Prerequisites](#prerequisites)
- [Installation](#installation)
- [Usage](#usage)
- [Configuration](#configuration)
- [Database Setup](#database-setup)
- [Database Schema](#database-schema)
- [API Playground](#api-playground)
- [Testing](#testing)
- [Troubleshooting](#troubleshooting)
- [Contributing](#contributing)
- [Future Work](#future-work)
- [License](#license)
- [Contact Information](#contact-information)

## Introduction

This project is a Golang-based backend application that serves as the API for a Next.js frontend application (located in a separate repository). It utilizes GraphQL for API queries and mutations and connects to a PostgreSQL database running inside a Docker container. The application is designed to be run locally using Docker for local testing and development purposes.

## Prerequisites

Before you begin, ensure you have the following installed:
- [Docker](https://docs.docker.com/get-docker/) (version 20.10 or later)
- [Docker Compose](https://docs.docker.com/compose/install/) (version 1.29 or later)
- [Make](https://www.gnu.org/software/make/) (version 4.3 or later)
  - macOS: Included with Xcode Command Line Tools or alternatively via Homebrew with `brew install make`
  - Linux: Use your distribution's package manager (e.g., `sudo apt install make` for Ubuntu)
  - Windows: Install via [Chocolatey](https://chocolatey.org/install) with `choco install make`
- [Go](https://golang.org/doc/install) (version 1.23.1 or later)

Please follow the links to find installation instructions for your specific operating system.

## Installation

1. Clone the repository and navigate to the project directory:
   ```
   git clone https://github.com/your-organization/iohk-golang-backend-preprod.git
   cd iohk-golang-backend-preprod
   ```

2. Build the Docker images:
   ```
   make docker-build
   ```

This process will download all necessary Go dependencies and build the application within the Docker environment.

## Running Locally

To run the application locally:

1. Start the Docker containers:
   ```
   make docker-up
   ```

   This command will start both the PostgreSQL database and the Go application.


2. The application should now be running. You can access the GraphQL playground at [http://localhost:8080/playground](http://localhost:8080/playground). You can view some example queries and mutations in the [API Playground](#api-playground) section.


3. To view the logs of the running containers (this is automatically run when you run `make docker-up`):
   ```
   make docker-logs
   ```

   This command will display the logs from all running containers. It's useful for debugging and monitoring the application's behavior.


4. To stop the application and all associated containers:
   ```
   make docker-down
   ```

Note: The application uses the `.env.local` file for configuration by default. If you need to modify any settings, you can edit this file before running `make docker-up`.

## Usage

This project uses a Makefile to simplify common operations. For a full list of available commands, run:

```
make help
```

Here are some useful commands:

- Build the application:
  ```
  make build
  ```

- Run the application:
  ```
  make run
  ```

- Start the Docker containers:
  ```
  make docker-up
  ```

- Stop the Docker containers:
  ```
  make docker-down
  ```

- View Docker logs:
  ```
  make docker-logs
  ```

- Run tests:
  ```
  make test
  ```

- Generate GraphQL code:
  ```
  make generate
  ```

This command should be run after making changes to the GraphQL schema in `graph/schema.graphqls`.

## Configuration

The application uses environment variables for configuration. These are stored in the `.env.local` file. Here's an example of the required variables:

```
POSTGRES_USER=your_username
POSTGRES_PASSWORD=your_password
POSTGRES_DB=your_database_name
POSTGRES_HOST=localhost
POSTGRES_PORT=5432
POSTGRES_SSLMODE=disable
DB_MAX_CONNS=25
DB_MIN_CONNS=5
DB_MAX_CONN_LIFETIME=5h
DB_MAX_CONN_IDLE_TIME=15m
DB_HEALTH_CHECK_PERIOD=1m
APP_PORT=8080
```

## Database Setup

The PostgreSQL database is automatically set up when you run `make docker-up`. The initial schema and any seed data are applied through the [init.sql](scripts/init.sql) file.

If you need to reset the database, you can run:

```
make docker-down
make docker-up
```

This will destroy the existing database and create a new one with the initial schema.

## Database Schema

The application uses a PostgreSQL database with a `customer` table. Below are the details of the table structure and constraints:

### Customer Table Columns

![Customer Table Columns](schema-diagrams/sql-customer-columns.png)

### Customer Table Checks

![Customer Table Checks](schema-diagrams/sql-customer-checks.png)

These checks ensure data integrity by enforcing rules such as:
- The birth date cannot be in the future
- The number of dependants cannot be negative
- The gender must be one of the predefined values: 'MALE', 'FEMALE'

### Initial Schema and Seed Data

The following SQL script [init.sql](scripts/init.sql) is used to generate the initial schema and seed data when the PostgreSQL container starts.

## API Playground

The GraphQL API can be explored using GraphQL Playground, which is available when running the application locally. To access it:

1. Start the application using `make docker-up`
2. Open a web browser and navigate to [http://localhost:8080/playground](http://localhost:8080/playground)

Here you can explore the schema, run queries, and test mutations.

Example query to get all customers:

```
query GetAllCustomers {
  customers {
    id
    name
    surname
    number
    gender
    country
    dependants
    birthDate
  }
}
```

Example mutation:

```
mutation CreateCustomer {
  createCustomer(input: {
    name: "John"
    surname: "Doe"
    number: 123
    gender: MALE
    country: "USA"
    dependants: 2
    birthDate: "1990-01-01"
  }) {
    id
    name
    surname
    number
    gender
    country
    dependants
    birthDate
  }
}
```

## Testing

To run the unit test suite:

```
make test
```

This command will ensure all dependencies are downloaded and then run all the tests in the project.

To run tests with coverage:

```
make coverage
```

This will download dependencies if needed, run the tests, and generate a `coverage.html` file that you can open in your browser to view detailed coverage information.

Note: If you're running tests outside of the Docker environment, these commands will automatically download any missing dependencies before running the tests.

## Troubleshooting

- If you encounter issues with Docker, ensure that the Docker daemon is running on your machine.
- If you see database connection errors, check that the PostgreSQL container is running and that your `.env.local` file has the correct database credentials.
- For any Go-related issues, ensure that your Go version matches the one specified in the `go.mod` file.

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

Please adhere to the existing code style and include appropriate tests for new features.

Before submitting your Pull Request, please:
1. Ensure your code follows the project's coding standards
2. Update the documentation as necessary
3. Add or update tests to cover your changes
4. Ensure all tests pass locally

We appreciate your contributions to improve this project!

## Future Work

- Deployment instructions for AWS and Vercel will be added in future updates.
- Additional API endpoints and features are planned for upcoming releases.

## License

This project is licensed under the MIT License - see the [LICENSE](https://opensource.org/licenses/MIT) file for details.

## Contact Information

For support or questions, please feel free to contact me at [joshwillems.cv@gmail.com](mailto:joshwillems.cv@gmail.com).