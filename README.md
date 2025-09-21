# 🚗 ColdWheels: Fueling Your Next Automotive Masterpiece

**Disrupting the garage-to-gallery pipeline with a hyper-scalable, cloud-native solution for showcasing automotive brilliance.**

## 🚀 Mission: Accelerate Everything

In a world where every second counts, ColdWheels is the ultimate pit stop for developers building next-generation automotive applications. We provide a robust, secure, and ridiculously fast backend, so you can focus on what matters: building stunning user experiences that move people.

## ✨ Why ColdWheels is Your Unfair Advantage:

*   **Blazing-Fast Performance:** Built with Go (Golang), Gin, and a PostgreSQL backend, ColdWheels is engineered for speed. Say goodbye to sluggish APIs and hello to instant gratification.
*   **Rock-Solid Security:** We've integrated JWT-based authentication and Argon2 password hashing to keep your users' data locked down tighter than a lug nut.
*   **Infinite Scalability:** With our S3-compatible image upload and presigned URL generation, you can handle a virtually unlimited number of high-resolution images without breaking a sweat.
*   **Effortless Integration:** A clean, RESTful API and CORS support make it a breeze to connect your frontend, whether it's a web app, mobile app, or a VR showroom.
*   **Developer-First Experience:** We've done the heavy lifting so you don't have to. With a straightforward setup and clear documentation, you'll be up and running in minutes.

## 🛠️ The Tech Stack: A Symphony of Power and Precision

*   **Go (Golang):** The heart of our engine, providing raw speed and concurrency.
*   **Gin:** A high-performance HTTP web framework for crafting elegant APIs.
*   **GORM:** The fantastic ORM library for Go, making database interactions a joy.
*   **PostgreSQL:** The world's most advanced open-source relational database.
*   **JWT (JSON Web Tokens):** The industry standard for secure authentication.
*   **Argon2:** A state-of-the-art password hashing algorithm.
*   **Docker:** For seamless containerization and deployment.

## 🏁 Get Started in 60 Seconds

### Prerequisites
- Docker
- Docker Compose
- Git

### 1. Clone the Repository
```bash
git clone https://github.com/silverstone-git/coldwheels-backend.git
cd coldwheels-backend
```

### 2. Environment Setup
Create a `.env.production` file in the root directory. You can copy the example if it exists:
```bash
cp .env.example .env.production
```
Populate it with your configuration. See the [Environment Variables](#-environment-variables) section for details.

### 3. Docker Network
Create the required Docker network:
```bash
docker network create --ipv6 --subnet 2001:db8::/64 my_ipv6_network
```

### 4. Launch!

#### Using Docker Compose (Recommended)
This is the easiest way to get up and running.
```bash
docker-compose up -d --build
```

#### Using Docker Run
If you prefer a manual approach:
```bash
# First, build the image
docker build -t wheelsimg .

# Then, run the container
docker run --env-file .env.production --name wheelc --network my_ipv6_network -d -p 4054:4054 wheelsimg
```

## ⚙️ Environment Variables

Your `.env.production` file configures the application. Here are the key variables:

| Variable              | Description                                           | Example                               |
| --------------------- | ----------------------------------------------------- | ------------------------------------- |
| `PORT`                | The port for the server to listen on.                 | `4054`                                |
| `CORS_ALLOWED_ORIGIN` | The origin allowed for CORS requests.                 | `http://localhost:3000`               |
| `JWT_SECRET`          | A strong, secret key for signing JWTs.                | `your_super_secret_key`               |
| `DB_HOST`             | The hostname of your PostgreSQL database.             | `localhost`                           |
| `DB_PORT`             | The port of your PostgreSQL database.                 | `5432`                                |
| `DB_USER`             | The username for the database connection.             | `your_db_username`                    |
| `DB_PASSWORD`         | The password for the database user.                   | `your_db_password`                    |
| `DB_NAME`             | The name of the database to use.                      | `coldwheels_db`                       |
| `DB_SSL`              | The SSL mode for the database connection.             | `disable`                             |
| `S3_BUCKET`           | Your S3-compatible bucket name for image uploads.     | `your-s3-bucket-name`                 |
| `S3_KEY`              | Your S3-compatible access key.                        | `your-s3-key`                         |
| `S3_SECRET`           | Your S3-compatible secret key.                        | `your-s3-secret`                      |
| `S3_ENDPOINT`         | The endpoint URL for your S3-compatible service.      | `your-s3-endpoint`                    |
| `S3_REGION`           | The region of your S3-compatible service.             | `your-s3-region`                      |

##  API Reference

### Public Routes

#### `POST /api/signup`
Registers a new user.

**Request Body:**
```json
{
  "email": "user@example.com",
  "password": "your_secure_password"
}
```

**Success Response (201 Created):**
```json
{
    "ID": "a1b2c3d4-e5f6-g7h8-i9j0-k1l2m3n4o5p6",
    "CreatedAt": "2025-09-21T10:00:00Z",
    "UpdatedAt": "2025-09-21T10:00:00Z",
    "DeletedAt": null,
    "email": "user@example.com",
    "password": "hashed_password_here"
}
```

#### `POST /api/login`
Authenticates a user and returns a JWT.

**Request Body:**
```json
{
  "email": "user@example.com",
  "password": "your_secure_password"
}
```

**Success Response (200 OK):**
```json
{
    "token": "your.jwt.token",
    "success": "success"
}
```

### Protected Routes
All protected routes require an `Authorization` header: `Bearer <your.jwt.token>`.

#### `GET /api/cars/:page`
Retrieves a paginated list of cars for the authenticated user.
- `:page` (URL parameter): The page number to retrieve.
- `pageSize` (query parameter, optional): The number of items per page (default: 6).

**Example Request:**
`GET /api/cars/1?pageSize=10`

**Success Response (200 OK):**
```json
[
    {
        "ID": 1,
        "CreatedAt": "2025-09-21T10:00:00Z",
        "UpdatedAt": "2025-09-21T10:00:00Z",
        "DeletedAt": null,
        "Make": "Toyota",
        "ModelName": "Supra",
        "Year": 2021,
        "EngineSize": 3.0,
        "FuelType": "Petrol",
        "Transmission": "Automatic",
        "Description": "A legendary sports car.",
        "OwnerID": "a1b2c3d4-e5f6-g7h8-i9j0-k1l2m3n4o5p6",
        "ImageURLs": ["https://presigned-url.com/image1.jpg"]
    }
]
```

#### `POST /api/cars`
Creates a new car entry.

**Request Body:**
```json
{
  "make": "Nissan",
  "modelName": "GT-R",
  "year": 2022,
  "engineSize": 3.8,
  "fuelType": "Petrol",
  "transmission": "Automatic",
  "description": "Godzilla in its latest form.",
  "imageURLs": ["s3-object-key-1.jpg", "s3-object-key-2.jpg"]
}
```

**Success Response (201 Created):**
Returns the newly created car object.

#### `PUT /api/cars/:id`
Updates an existing car by its ID.

**Request Body:**
Same as `POST /api/cars`.

**Success Response (200 OK):**
Returns the updated car object.

#### `DELETE /api/cars/:id`
Deletes a car by its ID.

**Success Response (200 OK):**
```json
{
    "message": "Car deleted",
    "success": "success"
}
```

#### `POST /api/cars/upload-images`
Uploads one or more images (up to 10). Send as `multipart/form-data`.

**Success Response (200 OK):**
```json
{
    "urls": [
        "https://your-s3-service.com/bucket/user-id/image1.jpg",
        "https://your-s3-service.com/bucket/user-id/image2.jpg"
    ]
}
```

## 🤝 Join the Revolution

ColdWheels is more than just a backend; it's a community of builders, dreamers, and automotive enthusiasts. We're always looking for contributors to help us push the boundaries of what's possible.

**Have an idea? Found a bug? Want to show off what you've built?**

*   **Open an issue:** [https://github.com/silverstone-git/coldwheels-backend/issues](https://github.com/silverstone-git/coldwheels-backend/issues)
*   **Fork the repo and submit a pull request.**

---

**Licensed under the MIT License. See the [LICENSE](LICENSE) file for details.**
