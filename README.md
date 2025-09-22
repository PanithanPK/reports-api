# Reports API

![Reports API](https://raw.githubusercontent.com/PanithanPK/reports-api/refs/heads/main/docs/dockergo.png)

<div align="center">
  <img src="https://img.shields.io/badge/Go-00ADD8?style=for-the-badge&logo=go&logoColor=white" alt="Go">
  <img src="https://img.shields.io/badge/Docker-2496ED?style=for-the-badge&logo=docker&logoColor=white" alt="Docker">
  <img src="https://img.shields.io/badge/MySQL-4479A1?style=for-the-badge&logo=mysql&logoColor=white" alt="MySQL">
  <img src="https://img.shields.io/badge/Fiber-00ACD7?style=for-the-badge&logo=go&logoColor=white" alt="Fiber">
  <img src="https://img.shields.io/badge/Swagger-85EA2D?style=for-the-badge&logo=swagger&logoColor=black" alt="Swagger">
  <img src="https://img.shields.io/badge/Node.js-339933?style=for-the-badge&logo=nodedotjs&logoColor=white" alt="Node.js">
</div>

<div align="center">
  <a href="https://github.com/PanithanPK/reports-api">
    <img src="https://img.shields.io/github/stars/PanithanPK/reports-api?style=social" alt="GitHub Stars">
  </a>
  <a href="https://github.com/GearRata/reports-ui">
    <img src="https://img.shields.io/github/stars/GearRata/reports-ui?style=social" alt="UI GitHub Stars">
  </a>
  <br>
  <a href="https://hub.docker.com/r/lovee12345e/reports-api">
    <img src="https://img.shields.io/docker/pulls/lovee12345e/reports-api?style=flat-square&logo=docker" alt="API Docker Pulls">
  </a>
  <a href="https://hub.docker.com/r/gearmcc/reports-ui">
    <img src="https://img.shields.io/docker/pulls/gearmcc/reports-ui?style=flat-square&logo=docker" alt="UI Docker Pulls">
  </a>
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

## System Requirements

### Hardware Requirements

#### Minimum
- **CPU:** 1 Core (x86_64 or ARM64)
- **RAM:** 512 MB
- **Storage:** 500 MB free space
- **Network:** Internet connection (for downloading dependencies)

#### Recommended
- **CPU:** 2+ Cores (x86_64 or ARM64)
- **RAM:** 2 GB or more
- **Storage:** 2 GB free space (including database)
- **Network:** Broadband Internet Connection

#### Production
- **CPU:** 4+ Cores
- **RAM:** 4 GB or more
- **Storage:** 10 GB+ (SSD recommended)
- **Network:** Dedicated Server/VPS

### Software Requirements

#### Required
- **Go:** Version 1.23.0 or later
  - Download: [https://go.dev/dl/](https://go.dev/dl/)
  - Check version: `go version`

- **MySQL:** Version 8.0 or later
  - Download: [https://dev.mysql.com/downloads/mysql/](https://dev.mysql.com/downloads/mysql/)
  - Or use MariaDB 10.5+

- **Git:** Latest version
  - Download: [https://git-scm.com/downloads](https://git-scm.com/downloads)

#### Optional
- **Docker:** Version 20.10+ and Docker Compose 2.0+
  - Download: [https://www.docker.com/get-started/](https://www.docker.com/get-started/)
  - For containerized deployment

- **Node.js:** Version 16+ (for development tools)
  - Download: [https://nodejs.org/](https://nodejs.org/)
  - Used for standard-version and changelog management

- **MinIO:** Latest version (for file storage)
  - Download: [https://min.io/download](https://min.io/download)

### Supported Operating Systems

#### Linux (Recommended)
- **Ubuntu:** 20.04 LTS, 22.04 LTS, 24.04 LTS
- **CentOS/RHEL:** 8, 9
- **Debian:** 11, 12
- **Alpine Linux:** 3.15+

#### Windows
- **Windows 10:** Build 1903 or later
- **Windows 11:** All versions
- **Windows Server:** 2019, 2022

#### macOS
- **macOS:** 10.15 (Catalina) or later
- **Apple Silicon (M1/M2):** Supported

### Dependencies & Versions

#### Go Modules (from go.mod)
```
Go Toolchain: 1.24.4
Main Dependencies:
├── github.com/gofiber/fiber/v2 v2.52.5
├── github.com/go-sql-driver/mysql v1.9.3
├── github.com/joho/godotenv v1.5.1
├── github.com/swaggo/swag v1.16.6
├── github.com/swaggo/fiber-swagger v1.3.0
├── github.com/go-telegram-bot-api/telegram-bot-api/v5 v5.5.1
├── github.com/minio/minio-go/v7 v7.0.95
├── github.com/nfnt/resize v0.0.0-20180221191011-83c6a9932646
└── golang.org/x/crypto v0.40.0
```

#### Development Tools
```
Node.js Dependencies:
├── standard-version v9.5.0 (for changelog management)
└── npm/yarn (package manager)

Go Tools:
├── swag CLI (go install github.com/swaggo/swag/cmd/swag@latest)
└── air (hot reload - optional)
```

### Prerequisites Check

Before starting installation, verify your system is ready:

```bash
# Check Go
go version
# Should show: go version go1.23.0 or later

# Check MySQL
mysql --version
# Should show: mysql Ver 8.0.x

# Check Git
git --version
# Should show: git version 2.x.x

# Check Docker (if using)
docker --version && docker-compose --version
# Should show: Docker version 20.10.x and Docker Compose version 2.x.x

# Check Node.js (if using)
node --version && npm --version
# Should show: v16.x.x or later
```

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

#### Status Values
The system uses the following status values for problem tracking:
- `0` - รอดำเนินการ (Pending)
- `1` - กำลังดำเนินการ (In Progress)
- `2` - เสร็จสิ้น (Resolved)

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