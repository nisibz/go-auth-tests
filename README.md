# Go Auth Tests API

This project is a Go-based RESTful API demonstrating user authentication and management functionalities. It utilizes JWT for secure authentication and MongoDB for data persistence.

## Features

- **User Registration**: Allows new users to create an account.
- **User Login**: Authenticates existing users and provides a JWT token.
- **Protected Routes**: Secures certain API endpoints, requiring a valid JWT token for access.
- **User Profile Management**:
  - Get user by ID.
  - List users (with pagination).
  - Update user information (name, email).
  - Delete user by ID.
- **Password Hashing**: Securely stores user passwords using bcrypt.

## Technologies Used

- **Go**: Programming language.
- **Gin**: HTTP web framework.
- **MongoDB**: NoSQL database for storing user data.
- **JWT (JSON Web Tokens)**: For stateless authentication.
- **godotenv**: For managing environment variables in development.
- **slog**: Structured, leveled logging.
- **bcrypt**: For password hashing.

## Getting Started

### Prerequisites

- Go
- MongoDB (running instance)
- Docker (optional, for easier MongoDB setup)

### Setup

1. **Clone the repository:**

   ```bash
   git clone https://github.com/nisibz/go-auth-tests.git
   cd go-auth-tests
   ```

2. **Set up environment variables:**

   Create a `.env` file in the root of the project by copying the example:

   ```bash
   cp .env.example .env
   ```

   Update the `.env` file with your specific configurations. Refer to `.env.example` for required variables.

3. **Install dependencies:**

   ```bash
   go mod tidy
   ```

4. **Run the application:**

   ```bash
   go run cmd/http/main.go
   ```

   The server will start, typically on `0.0.0.0:8080` (or as configured in your `.env` file).

### Docker Setup

You can also run the application using Docker:

```bash
# Copy the example env file
cp .env.example .env

# Start the containers
docker compose up -d
```

## API Endpoints

All endpoints are prefixed with `/api`.

### Auth Routes (`/api/auth`)

- `POST /register`: Register a new user.

  - Request Body: `{ "name": "John Doe", "email": "john.doe@example.com", "password": "securepassword123" }`

  **Example Response:**

  ```json
  {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
  }
  ```

- `POST /login`: Log in an existing user.

  - Request Body: `{ "email": "john.doe@example.com", "password": "securepassword123" }`

  **Example Response:**

  ```json
  {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
  }
  ```

### User Routes (`/api/users`)

_These routes require Bearer Token authentication via the `Authorization` header. The token is obtained from the `/login` or `/register` endpoint._

- `GET /:id`: Get user details by ID.

  **Example Response:**

  ```json
  {
    "id": "682d7fa1c28b28ae7128e452",
    "name": "John Doe",
    "email": "john.doe@example.com",
    "created_at": "2024-01-01T12:00:00Z"
  }
  ```

- `GET /`: List users (supports `limit` and `offset` query parameters).

  - Example: `/api/users?limit=5&offset=10`

  **Example Response:**

  ```json
  [
    {
      "id": "682d7fa1c28b28ae7128e452",
      "name": "John Doe",
      "email": "john.doe@example.com",
      "created_at": "2024-01-01T12:00:00Z"
    },
    {
      "id": "782d7fa1c28b28ae7128e453",
      "name": "Jane Smith",
      "email": "jane.smith@example.com",
      "created_at": "2024-01-02T08:30:00Z"
    }
  ]
  ```

- `PUT /`: Update the authenticated user's details (name, email). The user ID is derived from the JWT token.

  - Request Body: `{ "name": "Johnathan Doe", "email": "johnathan.doe@example.com" }`

  **Example Response:**

  ```json
  {
    "id": "682d7fa1c28b28ae7128e452",
    "name": "Johnathan Doe",
    "email": "johnathan.doe@example.com",
    "created_at": "2024-01-01T12:00:00Z"
  }
  ```

- `DELETE /:id`: Delete a user by ID.

  **Example Response:**

  HTTP Status: 204 No Content

  No response body.
