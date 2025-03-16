# ColdWheels Go Server

ColdWheels is a Go-based server that provides a simple REST API for car listings management. It uses Gin for HTTP routing, GORM for database interactions with PostgreSQL, and JWT for user authentication. It also incorporates image upload functionality with presigned URL generation.

## Table of Contents

- [Features](#features)
- [Prerequisites](#prerequisites)
- [Installation](#installation)
- [Environment Variables](#environment-variables)
- [Usage](#usage)
- [API Endpoints](#api-endpoints)
- [Directory Structure](#directory-structure)
- [License](#license)

## Features

- User signup and login with JWT authentication.
- CRUD operations for cars.
- Pagination support for listing cars.
- Password security using Argon2 hashing.
- Image uploading support that returns presigned URLs for car images.
- CORS handling for frontend integration.

## Prerequisites

- [Go](https://golang.org) (version >= 1.16 recommended)
- [PostgreSQL](https://www.postgresql.org)
- [Docker](https://www.docker.com) (if you prefer containerized deployment)
- [Git](https://git-scm.com) for source control

## Installation

1. **Clone the Repository**

   ```bash
   git clone https://github.com/silverstone-git/coldwheels-backend.git
   cd coldwheels-backend
   ```

2. **Install Dependencies**

   Use Go modules to install dependencies:

   ```bash
   go mod download
   ```

3. **Set Up Database**

   Make sure you have a PostgreSQL instance running. Create a database for the project.

4. **Environment Setup**

   Create a `.env` file in the root directory. See [Environment Variables](#environment-variables) below for details.

5. **Run the Server**

   For local development run:

   ```bash
   go run main.go
   ```

   The server will start on the port specified in your `.env` file.


## Installation using Docker Image

- Download the Docker Compose and create the `.env` file

    ```bash
    curl -L -o ./docker-compose.yaml https://raw.githubusercontent.com/silverstone-git/coldwheels-backend/main/docker-compose.yaml
    touch .env
    ```

- Edit the `.env` to your liking by referring to the [Environment Variables](#environment-variables)

- Pull and Up the [docker image](https://hub.docker.com/r/cyt0/gowheels)

    ```bash
    docker-compose pull
    docker-compose down
    docker-compose up -d
    ```

## Environment Variables

The project requires certain environment variables to:
- Connect to the database.
- Configure JWT authentication.
- Set up CORS.

Create a `.env` file in the project root with, for example:

```dotenv
# Server configuration
PORT="4054"
CORS_ALLOWED_ORIGIN="http://localhost:3000" # if your frontend is on local

# JWT
JWT_SECRET=your_super_secret_key

# Database configuration
DB_HOST=localhost
DB_USER=your_db_username
DB_PASSWORD=your_db_password
DB_NAME=coldwheels_db
DB_SSL=disable

# SPACES
SPACES_BUCKET="your s3 bucket name"
SPACES_KEY="your s3 key"
SPACES_SECRET="your s3 secret"
SPACES_ENDPOINT="your s3 endpoint"
SPACES_REGION="your s3 region"
MY_ORIGIN="http://localhost:4054" # if you are running it on local

```

> **Note:** Ensure that the `.env` file is in the same directory where you run your Go server so that it can be loaded by [godotenv](https://github.com/joho/godotenv).

## Usage

- **Start the Server:**

  ```bash
  go run main.go
  ```

- **Containerized Deployment:**

  You can also build a Docker image using the provided `Dockerfile` and run it via Docker Compose, if available in the repository. For example:

  ```bash
  docker-compose up --build
  ```

## API Endpoints

### Public Routes

- **POST /api/signup**

  Registers a new user.

  **Body:**

  ```json
  {
    "email": "user@example.com",
    "password": "your_password",
    "otherField": "..."
  }
  ```

- **POST /api/login**

  Authenticates a user and returns a JWT token.

  **Body:**

  ```json
  {
    "email": "user@example.com",
    "password": "your_password"
  }
  ```

### Protected Routes

All protected routes require the JWT token to be attached in the `Authorization` header (e.g., `Bearer <token>`).

- **GET /api/cars/:page**

  Retrieves a paginated list of cars belonging to the authenticated user. You can also pass the optional query parameter `pageSize`.

- **POST /api/cars**

  Creates a new car entry.

  **Body Schema (CarRequest):**

  ```json
  {
    "make": "Toyota",
    "modelName": "Corolla",
    "year": 2020,
    "engineSize": 1.8,
    "fuelType": "Petrol",
    "transmission": "Automatic",
    "description": "A reliable car",
    "imageURLs": ["image_key1", "image_key2"]
  }
  ```

- **PUT /api/cars/:id**

  Updates an existing car. Requires the car's ID as a URL parameter.

- **DELETE /api/cars/:id**

  Deletes a specific car.

- **POST /api/cars/upload-images**

  Handles image uploads and returns a list of object keys. The handler uses presigned URLs for secure uploads.

## Directory Structure

Below is a suggested directory structure for clarity:

```
coldwheels/
├── Dockerfile                   # Docker build instructions
├── docker-compose.yml           # Docker Compose file for local dev (if available)
├── .env                         # Environment variable definitions (not checked in)
├── go.mod                       # Go modules file
├── go.sum                       # Go dependencies checksums file
├── main.go                      # Main program file, contains server logic
├── middleware/
│   └── auth.go                  # Authentication middleware
├── lib/
│   └── ...                      # Database model definitions (User, Car, etc.)
├── repository/
│   └── imageupload.go           # Logic for handling image uploads
└── README.md                    # Project README
```

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

