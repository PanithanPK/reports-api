# Reports API

<div align="center">
  <img src="/dockergo.png" alt="Reports API Logo" width="700">
</div>

This project is a backend API built with Go for managing a problem reporting system. While a detailed description was not initially provided, this README aims to provide comprehensive information based on the available details and project structure.

## Key Features & Benefits

*   **Backend API:** Provides a robust and efficient backend for handling reports.
*   **Go Implementation:** Leverages the performance and concurrency features of Go.
*   **Dockerized:** Easily deployable and scalable using Docker.
*   **Structured Architecture:** Follows MVC, Repository, and Middleware patterns for maintainability.
*   **API Documentation:** Includes Swagger documentation for easy API usage.
*   **Database Integration:** Designed to connect to a database (potentially MySQL based on other files).
*   **Notification System:** Potentially integrates with Telegram for notifications (based on `PROJECT_OVERVIEW.md`).

## Prerequisites & Dependencies

Before you begin, ensure you have the following installed:

*   **Go (version 1.23 or higher):**  [https://go.dev/dl/](https://go.dev/dl/)
*   **Docker:** [https://www.docker.com/get-started/](https://www.docker.com/get-started/)
*   **Node.js (for some tooling, version may vary):** [https://nodejs.org/](https://nodejs.org/) (check package.json for version requirements if any).  While not strictly required for running the Go backend, it's listed as a technology used in the project.

## Installation & Setup Instructions

1.  **Clone the repository:**

    ```bash
    git clone https://github.com/PanithanPK/reports-api.git
    cd reports-api
    ```

2.  **Build the Docker image (Development):**

    ```bash
    docker build -t <your-dockerhub-username>/reports-api:dev -f Dockerfile.dev .
    ```

    Replace `<your-dockerhub-username>` with your Docker Hub username.

3.  **Run the Docker container:**

    ```bash
    docker run -p 8080:8080 <your-dockerhub-username>/reports-api:dev
    ```

    This assumes the application runs on port 8080.

4.  **Alternatively, use the build script (build2.sh):**

    Execute the `build2.sh` script in the root directory for interactive build options.

    ```bash
    ./build2.sh
    ```

5. **Database Configuration:**

   *   The application likely requires a database. Refer to the code and potentially `.env` files for database connection details. You'll likely need to set up a database server (e.g., MySQL) and configure the connection string.

## API Documentation & Swagger Usage

### Swagger Documentation

This API includes comprehensive Swagger/OpenAPI documentation for all endpoints.

#### Generating Swagger Documentation

1. **Install swag CLI tool:**
   ```bash
   go install github.com/swaggo/swag/cmd/swag@latest
   ```

2. **Generate swagger documentation:**
   ```bash
   swag init
   ```
   This will generate `docs/swagger.json` and `docs/swagger.yaml` files.

#### Accessing Swagger UI

**Option 1: Built-in Swagger UI (if implemented)**
- Start the application
- Navigate to `http://localhost:8080/swagger/index.html`

**Option 2: Online Swagger Editor**
- Go to [https://editor.swagger.io/](https://editor.swagger.io/)
- Upload the generated `docs/swagger.json` file

**Option 3: Local Swagger UI**
- Use Docker to run Swagger UI locally:
  ```bash
  docker run -p 8081:8080 -e SWAGGER_JSON=/swagger.json -v $(pwd)/docs/swagger.json:/swagger.json swaggerapi/swagger-ui
  ```
- Access at `http://localhost:8081`

#### API Usage
- **API_USAGE.md:** Refer to `docs/API_USAGE.md` for detailed instructions and examples
- **Interactive Testing:** Use Swagger UI to test API endpoints directly
- **API Reference:** Complete endpoint documentation with request/response examples

## Configuration Options

The application can be configured using environment variables.  Based on the `Dockerfile` and other files, it looks for `.env` files:

*   `.env.prod`:  Used for production environments.
*   `.env`: Generic environment configuration.

Common configuration options might include:

*   **Database connection details (host, port, username, password, database name)**
*   **Port number:** The port the application listens on.
*   **Telegram Bot configuration (if notifications are enabled)**
*   **Authentication settings (if JWT is used)**

Example `.env` file:

```
DATABASE_HOST=localhost
DATABASE_PORT=3306
DATABASE_USER=your_user
DATABASE_PASSWORD=your_password
DATABASE_NAME=reports_db
PORT=8080
```

## Contributing Guidelines

We welcome contributions to this project!

1.  Fork the repository.
2.  Create a new branch for your feature or bug fix.
3.  Make your changes.
4.  Commit your changes with clear and descriptive commit messages (following [standard-version](https://github.com/conventional-changelog/standard-version) conventions is recommended).
5.  Push your branch to your fork.
6.  Submit a pull request.

## License Information

License information is not specified. Please add a relevant license (e.g., MIT, Apache 2.0) to the repository.

## Acknowledgments

*   This project utilizes the Fiber Framework for Go.
*   Standard Version for changelog management.