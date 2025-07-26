````markdown
# First:V1 – Mini Social Media Backend

A mini social media backend built in Go where users can register, login, post thoughts, follow others, and view a personalized feed. Redis is used for caching, and JWT ensures secure authentication.

 🔧 Tech Stack

- Language: Go
- Database: MySQL
- Cache: Redis
- Auth: JWT (JSON Web Tokens)
- API: REST (built with net/http or Gin)
- Tools: Postman, Git

 ✨ Features

- 🧑‍💻 User Authentication
  - Registration and secure JWT-based login
  - Role-based access supported (optional)

- 📝 Create Posts
  - Users can post their thoughts in text format

- 👥 Follow/Unfollow
  - Follow other users and build a social graph

- 📰 User Feed
  - Personalized post feed based on followed users
  - Cached using Redis for better performance

- ❤️ Like System
  - Like and unlike any post

- 🧠 User Suggestions
  - Suggested users to follow based on mutual connections

- ⚡ Redis Caching
  - Profile and feed data cached to reduce DB hits

 🚀 Getting Started

# 1. Clone the repo

```bash
git clone https://github.com/RathodViraj/First.git
cd First
````

# 2. Set up environment variables

Create a `.env` file (or use environment export) with the following:

```env
DB_USER=root
DB_PASS=yourpassword
DB_NAME=socialmedia
DB_HOST=localhost
JWT_SECRET=your_jwt_secret
REDIS_ADDR=localhost:6379
```

# 3. Run the app

bash
go run main.go


Make sure Redis and MySQL servers are running locally.

 📁 Folder Structure

```
.
├── handler/        # HTTP handlers (controllers)
├── service/        # Business logic
├── repository/     # DB and cache operations
├── model/          # Data structures (User, Post)
├── utils/          # Utility functions (JWT, hashing)
├── main.go         # Entry point
└── go.mod
```

 📬 API Endpoints (Sample)

 `POST /register` — Register user
 `POST /login` — Login and get JWT
 `POST /posts` — Create post
 `GET /feed` — Get personalized feed
 `POST /follow/{id}` — Follow a user
 `POST /unfollow/{id}` — Unfollow a user
 `GET /suggestions` — Get user suggestions

 📌 Future Improvements

 Add comments and notifications
 Rate-limiting using Redis
 Image upload for posts
 GraphQL version (in progress)

---

 👨‍💻 Author

Viraj Rathod
Backend Developer | [LinkedIn](www.linkedin.com/in/viraj-rathod-058ba4280) | [GitHub](github.com/RathodViraj)

```

```
Let me know if you'd like a shorter version too, or a README template for your next project. Once you push this to GitHub, your project will look much more professional to recruiters and collaborators.
```
