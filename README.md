````markdown
First:v2 – Mini Social Media Backend

A mini social media backend built in Go where users can register, login, post thoughts, follow others, and view a personalized feed. It uses Neo4j to manage relationships and feed generation, Redis for caching, JWT for authentication, and WebSocket for real-time notifications.

---

🔧 Tech Stack

- Language: Go
- Databases: MySQL, Neo4j (Graph DB)
- Cache: Redis
- Auth: JWT (JSON Web Tokens)
- API: REST (net/http or Gin)
- Real-time: WebSocket (Notifications)
- Tools: Postman, Git

---

✨ Features

- 🧑‍💻 User Authentication
  - Secure registration and login with JWT-based session management

- 📝 Post System
  - Users can create and view text-based posts

- 👥 Follow/Unfollow
  - Maintain a social graph using Neo4j for scalable relationship modeling

- 📰 Personalized Feed
  - Feed generation based on following graph using Neo4j traversal
  - Feed cached in Redis to optimize performance

- ❤️ Like System
  - Like or unlike posts with idempotent behavior

- 🔔 Real-time Notifications
  - Implemented using WebSocket (e.g., when followed or post liked)

- 🧠 User Suggestions
  - Suggested users to follow based on mutual connections (graph-based logic)

- 🔍 Search Functionality
  - Search users or posts by keywords (MySQL-based search)

- ⚡ Redis Caching
  - Cached user profiles and feeds to reduce DB hits and improve latency

---

🚀 Getting Started

### 1. Clone the Repository

```bash
git clone https://github.com/RathodViraj/First.git
cd First
````

2. Set Up Environment Variables

Create a `.env` file:

```env
DB_USER=root
DB_PASS=yourpassword
DB_NAME=socialmedia
DB_HOST=localhost

NEO4J_URI=bolt://localhost:7687
NEO4J_USER=neo4j
NEO4J_PASS=your_neo4j_password

JWT_SECRET=your_jwt_secret

REDIS_ADDR=localhost:6379
```

3. Run the Application

```
bash
go run main.go
```

Ensure MySQL, Redis, and Neo4j are running locally before starting.

---

📁 Folder Structure

```
.
├── handler/        # API handlers (HTTP)
├── service/        # Business logic
├── repository/     # DB and graph/cache interactions
├── graph/          # Neo4j graph functions
├── ws/             # WebSocket handling
├── model/          # Data structures (User, Post)
├── utils/          # JWT, hashing, helpers
├── main.go         # Entry point
└── go.mod
```



📬 API Overview (Sample)

* `POST /register` — User registration
* `POST /login` — JWT-based login
* `POST /posts` — Create a new post
* `GET /feed` — Get personalized feed
* `POST /follow/{id}` — Follow a user
* `POST /unfollow/{id}` — Unfollow a user
* `GET /suggestions` — Get user suggestions
* `GET /search/users?query=` — Search users
* `GET /search/posts?query=` — Search posts
* WebSocket: `/ws` — Real-time notification endpoint



🔄 Future Enhancements

* Add comments system
* Notification history and storage
* Post pagination and infinite scrolling
* GraphQL API version
* Role-based access and admin panel

```
👨‍💻 Author

Viraj Rathod
Backend Developer
📧 [virajrathod631@gmail.com](mailto:virajrathod631@gmail.com)
📞 +91-8799242278
[GitHub](https://github.com/RathodViraj) | [LinkedIn](http://www.linkedin.com/in/viraj-rathod-058ba4280)
```


