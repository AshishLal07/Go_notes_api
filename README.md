# Secure Notes API

A secure REST API built with Go (Fiber framework) for user authentication and personal notes management. Only authenticated users can access and manage their own notes.

## Features

- **User Authentication**: Registration and login with JWT tokens
- **Secure Password Handling**: bcrypt hashing for passwords
- **Personal Notes Management**: CRUD operations for notes
- **Authorization**: Users can only access their own notes
- **Pagination & Search**: Notes can be paginated and searched
- **Docker Support**: Complete Docker setup with MySQL
- **Database Seeding**: CLI tool to populate sample data
- **Production Ready**: Proper error handling, validation, and logging

## Tech Stack

- **Backend**: Go 1.21 + Fiber v2
- **Database**: MySQL 8.0 + GORM
- **Authentication**: JWT with custom middleware
- **Password Hashing**: bcrypt
- **Containerization**: Docker + Docker Compose
- **Configuration**: Environment variables with .env support

## Project Structure

\`\`\`
notes-api/
├── cmd/
│ └── seed/ # Database seeding CLI
├── config/ # Database configuration
├── handlers/ # HTTP request handlers
├── middleware/ # Custom middleware (JWT auth)
├── models/ # Data models and validation
├── routes/ # Route definitions
├── utils/ # Utility functions (JWT, validation)
├── docker-compose.yml # Docker services configuration
├── Dockerfile # Application container
├── .env.example # Environment variables template
└── main.go # Application entry point
\`\`\`

## Quick Start

### Prerequisites

- Go 1.21+
- Docker & Docker Compose

### 1. Clone and Setup

```bash
# Clone the repository
git clone <repository-url>
cd notes-api

# Copy environment configuration
cp .env.example .env

# Edit .env with your preferred settings (optional)
nano .env

```

### 2. Docker Commands

```bash
# Start with automatic seeding enabled
docker-compose --profile seed up -d

# Check if services are running
docker-compose ps

# View logs
docker-compose logs -f
```
