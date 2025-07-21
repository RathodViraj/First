# First

A mini social media backend API built with Go and Gin.

## Features
- User registration and authentication (JWT)
- Create, delete, and fetch posts
- Like and unlike posts
- Comment on posts
- Follow and unfollow users
- User feed and profile endpoints
- Redis caching for performance
- MySQL database support

## Project Structure
```
First/
  chachingService/      # Redis caching logic
  db/                   # Database connection logic
  handler/              # HTTP route handlers (controllers)
  middleware/           # Gin middleware (e.g., auth)
  model/                # Data models
  repository/           # Data access layer
  service/              # Business logic
  main.go               # Application entry point
```

## Setup
1. **Clone the repository:**
   ```sh
   git clone github.com/RathodViraj/First/
   cd First
   ```
2. **Configure your database:**
   - Update your MySQL and Redis connection settings in `db/mysql.go` and `db/redis.go` as needed.

3. **Install dependencies:**
   ```sh
   go mod tidy
   ```

4. **Run the application:**
   ```sh
   go run main.go
   ```
   The server will start on `:8080` by default.

## API Endpoints
- `POST   /register` — Register a new user
- `POST   /login` — Login and receive JWT
- `GET    /home` — Recent posts
- `POST   /posts` — Create a post
- `GET    /posts/:id` — Get a post
- `DELETE /posts/:id` — Delete a post
- `POST   /posts/:id/like` — Like a post
- `DELETE /posts/:id/like` — Unlike a post
- `GET    /posts/:id/likes` — Get users who liked a post
- `GET    /posts/:id/comments` — Get comments for a post
- `POST   /posts/:id/comments` — Add a comment
- `GET    /users/:id` — Get user profile
- `DELETE /users/:id` — Delete user
- `GET    /users/:id/home` — Get user feed
- `GET    /users/:id/followers` — Get followers
- `GET    /users/:id/followings` — Get followings
- `GET    /users/:id/mutual` — Get mutual connections
- `POST   /follow/:follower_id/:following_id` — Follow a user
- `DELETE /unfollow/:follower_id/:following_id` — Unfollow a user

## Environment Variables
- `JWT_SECRET` — Secret key for JWT signing (set in `main.go` or as an environment variable)

## Notes
- Make sure MySQL and Redis are running and accessible.
- For production, set `GIN_MODE=release`.

---

**Contributions welcome!** 
