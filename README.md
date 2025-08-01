# Reports API

## Overview
The Reports API is a backend service built with Go that provides endpoints for managing problems, tasks, IP phones, programs, departments, and branches. It also includes authentication routes for user and admin management.

## Features
- **Problem Management**: Create, update, delete, and list problems.
- **Task Management**: Create, update, and track tasks with phone and system details.
- **IP Phone Management**: Manage IP phones.
- **Program Management**: Manage programs.
- **Department Management**: Manage departments.
- **Branch Management**: Manage branches.
- **Authentication**: User and admin registration, login, and logout.

## API Endpoints

### Problem Routes
- `GET /api/v1/problem/list`: List all problems.
- `POST /api/v1/problem/create`: Create a new problem.
- `PUT /api/v1/problem/update/{id}`: Update a problem by ID.
- `DELETE /api/v1/problem/delete/{id}`: Delete a problem by ID.

### Task Routes
- `GET /api/v1/task/list`: List all tasks with detailed information.
- `POST /api/v1/task/create`: Create a new task.
- `PUT /api/v1/task/update/{id}`: Update a task by ID.
- `DELETE /api/v1/task/delete/{id}`: Delete a task by ID.

### IP Phone Routes
- `GET /api/v1/ipphone/list`: List all IP phones.
- `POST /api/v1/ipphone/create`: Create a new IP phone.
- `PUT /api/v1/ipphone/update/{id}`: Update an IP phone by ID.
- `DELETE /api/v1/ipphone/delete/{id}`: Delete an IP phone by ID.

### Program Routes
- `GET /api/v1/program/list`: List all programs.
- `POST /api/v1/program/create`: Create a new program.
- `PUT /api/v1/program/update/{id}`: Update a program by ID.
- `DELETE /api/v1/program/delete/{id}`: Delete a program by ID.

### Department Routes
- `GET /api/v1/department/list`: List all departments.
- `POST /api/v1/department/create`: Create a new department.
- `PUT /api/v1/department/update/{id}`: Update a department by ID.
- `DELETE /api/v1/department/delete/{id}`: Delete a department by ID.

### Branch Routes
- `GET /api/v1/branch/list`: List all branches.
- `POST /api/v1/branch/create`: Create a new branch.
- `PUT /api/v1/branch/update/{id}`: Update a branch by ID.
- `DELETE /api/v1/branch/delete/{id}`: Delete a branch by ID.

### Authentication Routes
- `POST /authEntry/login`: Login.
- `POST /authEntry/registerUser`: Register a new user.
- `POST /authEntry/registerAdmin`: Register a new admin.
- `PUT /authEntry/updateUser`: Update user information.
- `DELETE /authEntry/deleteUser`: Delete a user.
- `POST /authEntry/logout`: Logout.

## Setup

### Prerequisites
- Go 1.18 or later
- MySQL database

### Installation
1. Clone the repository:
   ```bash
   git clone https://github.com/PanithanPK/reports-api.git
   ```
2. Navigate to the project directory:
   ```bash
   cd reports-api
   ```
3. Install dependencies:
   ```bash
   go mod tidy
   ```

### Configuration
Create a `.env` file in the root directory with the following variables:
```env
PORT=5000
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASSWORD=password
DB_NAME=reports
```

### Running the Application
Start the server:
```bash
go run main.go
```

The server will run on `http://localhost:5000`.

### Environment Variables
The application uses the following environment variables:
- `PORT`: Server port (default: 5000)
- `DB_HOST`: Database host
- `DB_PORT`: Database port
- `DB_USER`: Database username
- `DB_PASSWORD`: Database password
- `DB_NAME`: Database name
- `JWT_SECRET`: Secret key for JWT token generation
- `TOKEN_HOUR_LIFESPAN`: JWT token lifespan in hours

## Security

### Authentication
- JWT-based authentication system
- Secure password hashing using bcrypt
- Role-based access control for admin and regular users

### API Security
- HTTPS support for secure communication
- Input validation to prevent injection attacks
- Rate limiting to prevent brute force attacks
- CORS configuration for frontend integration

## Error Handling
The API uses standard HTTP status codes:
- `200 OK`: Request successful
- `201 Created`: Resource created successfully
- `400 Bad Request`: Invalid input
- `401 Unauthorized`: Authentication required
- `403 Forbidden`: Insufficient permissions
- `404 Not Found`: Resource not found
- `500 Internal Server Error`: Server-side error

All error responses include a JSON object with an `error` field containing a descriptive message.

## Deployment

### Docker
A Dockerfile is provided for containerization:
```bash
# Build the Docker image
docker build -t reports-api .

# Run the container
docker run -p 5000:5000 --env-file .env reports-api
```

### Production Considerations
- Use a reverse proxy (Nginx, Traefik) for SSL termination
- Set up proper logging and monitoring
- Configure database connection pooling
- Implement proper backup strategies

## License
This project is licensed under the MIT License.
