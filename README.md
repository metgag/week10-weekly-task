# 🎬 Tickitz Movie API

Backend project for **Tickitz Movie** built with **Go (Gin Gonic)**.  
This project includes **struct validation**, **JWT authentication**, **argon2 password hashing**, uses **PostgreSQL** as the main database, **Redis** for caching, and **Swagger** for API documentation.

---

## 🛠 Tech Stack

- **Go (Golang)** with Gin
- **PostgreSQL** ([pgx - PostgreSQL Driver](https://github.com/jackc/pgx))
- **Redis**
- **JWT** ([jwt-go](https://github.com/golang-jwt/jwt))
- **Argon2** ([argon2](https://pkg.go.dev/golang.org/x/crypto@v0.41.0/argon2))
- **Migrate** ([db migration tools](https://github.com/golang-migrate/migrate))
- **Docker & Docker Compose**
- **Swagger** (via [Swaggo](https://github.com/swaggo/swag))

---

## 🌐 Environment Variables

Copy `.env.example` to `.env` and fill in your configuration:

```env
# PostgreSQL
DB_URL_M=postgres://<PG_USER>:<PG_PASSWORD>@<PG_HOST>:<PG_PORT>/<PG_DATABASE>

# JWT
JWT_SECRET=<YOUR_JWT_SECRET>
JWT_ISSUER=<YOUR_JWT_ISSUER>

# Redis
REDIS_HOST=redis
REDIS_PORT=6379
````

---

## ⚡ Installation

1. **Clone the repository**

```bash
git clone git@github.com:metgag/week10-weekly-task.git
```

2. **Navigate to the project directory**

```bash
cd week10-weekly-task
```

3. **Install dependencies**

```bash
go mod tidy
```

4. **Set up environment variables**
   Create a `.env` file from `.env.example`.

5. **Install migrate for DB migrations**
   Follow the instructions in the migrate documentation.

6. **Run the database migration**

```bash
make migrate-up
```

7. **Start the server**

```bash
go run ./cmd/main.go
```

---

## 🚧 API Documentation

### Admin Routes

| Method | Endpoint      | Body | Description                 |
| ------ | ------------- | ---- | --------------------------- |
| GET    | /admin/orders | —    | Get all orders (Admin only) |

#### Admin Movie Routes

| Method | Endpoint          | Body                  | Description                   |
| ------ | ----------------- | --------------------- | ----------------------------- |
| GET    | /admin/movies     | —                     | Get all movies (Admin only)   |
| POST   | /admin/movies     | title, synopsis, etc. | Create new movie (Admin only) |
| PATCH  | /admin/movies/:id | title, synopsis, etc. | Update movie (Admin only)     |
| DELETE | /admin/movies/:id | —                     | Delete movie (Admin only)     |

---

### Auth Routes

| Method | Endpoint       | Body            | Description             |
| ------ | -------------- | --------------- | ----------------------- |
| POST   | /auth/register | email, password | Register new user       |
| POST   | /auth/login    | email, password | Login                   |
| DELETE | /auth/logout   | —               | Logout (requires token) |

---

### Cinema Routes

| Method | Endpoint                       | Body | Description                             |
| ------ | ------------------------------ | ---- | --------------------------------------- |
| GET    | /cinemas/schedules             | —    | Get cinema schedules                    |
| GET    | /cinemas/:schedule_id/seats    | —    | Get seats for a specific schedule       |
| GET    | /cinemas/:schedule_id/selected | —    | Get cinema name and time for a schedule |

---

### Movie Routes

| Method | Endpoint              | Body | Description                               |
| ------ | --------------------- | ---- | ----------------------------------------- |
| GET    | /movies/upcoming      | —    | Get upcoming movies                       |
| GET    | /movies/popular       | —    | Get popular movies                        |
| GET    | /movies               | —    | Get movies with genre, pagination, search |
| GET    | /movies/:id           | —    | Get movie details by ID                   |
| GET    | /movies/:id/schedules | —    | Get movie schedules                       |
| GET    | /movies/:id/schedule  | —    | Get filtered movie schedule               |
| GET    | /movies/genres        | —    | Get all genres                            |

---

### Orders

| Method | Endpoint | Body                  | Description                  |
| ------ | -------- | --------------------- | ---------------------------- |
| POST   | /orders  | movie_id, seats, etc. | Create new order (User only) |

---

### User Routes

| Method | Endpoint        | Body             | Description                        |
| ------ | --------------- | ---------------- | ---------------------------------- |
| GET    | /users/         | —                | Get user info (User only)          |
| PATCH  | /users/         | user info fields | Update user info (User only)       |
| GET    | /users/orders   | —                | Get user order history (User only) |
| PATCH  | /users/password | password fields  | Update user password (User only)   |

---

### 📡 Swagger Docs

Full API documentation is available via Swagger:
[http://localhost:6011/swagger/index.html](http://localhost:6011/swagger/index.html)

---

## 📄 License

MIT License © 2025 Tickitz

---

## 📧 Contact

**Author:** Slamet Gagah
**Email:** [rayhansjah88@gmail.com](mailto:rayhansjah88@gmail.com)

---

## 🎯 Related Project

**[Tickitz Frontend (React)](https://github.com/metgag/metgag-loket-tickitz)**