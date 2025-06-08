📅 EventManager – Go Backend with Auth, Redis Caching, and Async Emails
A secure and performant backend system built in Go, combining RESTful architecture, JWT-based authentication, event-driven email dispatch, and Redis-powered caching.

🚀 Overview
EventManager is a backend service tailored for managing users in an event-driven context. It provides:

🔐 Secure JWT-based user authentication

💡 Redis caching for performance boosts on frequent reads

🔁 Event-driven emailing to decouple and asynchronously handle side-effects like welcome emails

⚡️ Scalable Go architecture using Gin, PostgreSQL, Redis, and Resend API

This project is a showcase of how a modern Go backend can blend synchronous REST operations with asynchronous, event-based responsibilities while maintaining high performance and security.

🧰 Tech Stack
Layer	Technology
Language	Go
Framework	Gin
ORM	GORM
Database	PostgreSQL
Auth	JWT (access & refresh tokens)
Caching & Pub/Sub	Redis
Email	Resend


🧠 Architecture
This is a hybrid system that blends traditional REST with selective event-driven design and stateful JWT auth:

sql
Copy
Edit
[ Client ]
    ↓
 POST /register
    ↓
 Save user → Hash password → Create tokens
    ↓
 Cache user in Redis → Publish "user:created" event
    ↓
[ Redis Subscriber ]
    ↓
 Consume event → Send welcome email via Resend
🔐 JWT Authentication
This project uses access and refresh tokens:

Access tokens are short-lived (e.g., 15 mins) and used for authenticating routes.

Refresh tokens are long-lived (e.g., 7 days) and used to generate new access tokens without requiring a fresh login.

Auth Flow
plaintext
Copy
Edit
POST /register  → Issue access + refresh tokens
POST /login     → Verify credentials → Issue tokens
GET /me         → Auth-required route using access token
POST /refresh   → Use refresh token to issue new access token
Tokens are signed using a secret stored in your environment config.

📦 Redis Responsibilities
Redis is used for two purposes:

Role	Description
Pub/Sub Engine	Handles async side-effect: welcome email dispatch
Caching Layer	Stores user data by email for fast access

This dual-purpose usage helps minimize DB overhead and decouple performance-sensitive operations from non-critical tasks.

⚙️ Endpoints (Highlights)
POST /register
Registers a user and returns access/refresh tokens. Triggers welcome email event.

POST /login
Logs in a user and returns tokens.

POST /refresh
Generates a new access token using a refresh token.

GET /me
Protected route — returns current user info if access token is valid.

Feature	Value Delivered
🔐 JWT auth	Stateless, scalable user session management
⚡️ Redis caching	Rapid read access and reduced DB pressure
📨 Event-driven emails	Faster response time on registration
🧱 Clean modular structure	Easily maintainable and extendable
📬 Resend API integration	Real-world async external service interaction

🌱 Future Enhancements
Admin interface for user analytics

Email verification and password reset

Monitoring with Prometheus and Grafana

Rate-limiting & brute-force protection

Role-based access control
