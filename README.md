# 🔐 Osto Assignment – Secure CLI Authentication System

A production-inspired command-line authentication system built with **Go**, **PostgreSQL**, **GORM**, and **Docker**. The application supports secure user authentication, session management, account lockout protection, and optional Two-Factor Authentication (TOTP).

---

## ✨ Features

* ✅ User Registration
* ✅ Secure Login
* ✅ Password Hashing (bcrypt)
* ✅ Session Management
* ✅ "Who Am I" Command
* ✅ Logout
* ✅ Account Lockout after Multiple Failed Login Attempts
* ✅ Two-Factor Authentication (TOTP)
* ✅ Google Authenticator Compatible
* ✅ PostgreSQL Database
* ✅ GORM ORM
* ✅ Docker & Docker Compose Support
* ✅ Interactive CLI using Readline
* ✅ Command History
* ✅ Auto Migration

---

# 🛠 Tech Stack

| Technology     | Purpose                   |
| -------------- | ------------------------- |
| Go             | Backend                   |
| PostgreSQL     | Database                  |
| GORM           | ORM                       |
| Docker         | Containerization          |
| Docker Compose | Multi-container Setup     |
| bcrypt         | Password Hashing          |
| TOTP           | Two-Factor Authentication |
| Readline       | Interactive CLI           |

---

# 📂 Project Structure

```text
.
├── cmd/
│   └── main.go
│
├── internal/
│   ├── cli/
│   ├── database/
│   ├── models/
│   ├── repository/
│   ├── services/
│   ├── totp/
│   └── utils/
│
├── Dockerfile
├── docker-compose.yml
├── .env
├── go.mod
└── README.md
```

---

# 🚀 Running Locally

## Clone Repository

```bash
git clone <repository-url>
cd OstoAssignment
```

## Install Dependencies

```bash
go mod download
```

## Configure Environment

Create a `.env` file.

```env
DB_HOST=127.0.0.1
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=ostoassignment
DB_SSLMODE=disable
```

Run:

```bash
go run ./cmd
```

---

# 🐳 Running with Docker

## Build & Start

```bash
docker compose up --build
```

## Stop Containers

```bash
docker compose down
```

## Remove Containers & Database Volume

```bash
docker compose down -v
```

---

# 🔑 Available Commands

```
register
login
logout
whoami
enable2fa
confirm2fa
disable2fa
help
exit
```

---

# 🔒 Security Features

### Password Hashing

Passwords are securely hashed using bcrypt before being stored in the database.

---

### Session Management

Each successful login creates a new authenticated session.

---

### Account Lockout

* Failed login attempts are tracked.
* Account is temporarily locked after multiple failed attempts.
* Lock expires automatically after the configured duration.

---

### Two-Factor Authentication

Supports Time-based One-Time Passwords (TOTP).

Compatible with:

* Google Authenticator
* Microsoft Authenticator
* Authy

---

# 🗄 Database

The application automatically creates the required tables using GORM AutoMigrate.

Tables:

* users
* sessions

---

# 🐳 Docker Architecture

```text
               Docker Network
        ┌───────────────────────────┐

        ┌──────────────┐
        │  auth-cli    │
        │ Go CLI App   │
        └──────┬───────┘
               │
               │ GORM
               │
        ┌──────▼───────┐
        │   auth-db    │
        │ PostgreSQL   │
        └──────────────┘
```

---

# 📸 Sample Workflow

```text
register

↓

login

↓

(Optional)

Enter OTP

↓

Session Created

↓

whoami

↓

logout
```

---

# 📌 Future Improvements

* QR Code Generation for 2FA
* Password Reset
* Email Verification
* Refresh Tokens
* Role-Based Access Control
* Audit Logs
* Unit Tests
* CI/CD Pipeline
* Multi-stage Docker Build

# 📸 Screenshots

## Home

Interactive CLI startup.

<img src="screenshot/home.png" width="900">

---

## User Registration

Create a new account.

<img src="screenshot/register.png" width="900">

---

## User Login

Secure authentication with session creation.

<img src="screenshot/login.png" width="900">

---

## Enable Two-Factor Authentication

Generate a secret compatible with Google Authenticator.

<img src="screenshot/enable-2fa.png" width="900">

---

## QR Code

Scan using Google Authenticator.

<img src="screenshot/2fa.png" width="350">

---

## Current User

Display authenticated user information.

<img src="screenshot/whoami.png" width="900">

---

## Help Menu

Available commands after login.

<img src="screenshot/helpafterlogin.png" width="900">

---

## Logout

Terminate the active session.

<img src="screenshot/logout.png" width="900">

---

# 👨‍💻 Author

**Sahil Patel**

Backend Developer (Go)

GitHub: https://github.com/SahilPatel8826

LinkedIn: https://www.linkedin.com/in/sahil-patel-9264b6397

---

## 📄 License

This project is created for educational and assignment purposes.
