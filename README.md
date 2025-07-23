# Todo App

This is a simple Todo application built with React (TypeScript) for the frontend and Go (Gin) for the backend, using PostgreSQL as the database.

## Features

- User authentication (signup, login, logout)
- User-specific Todo management (add, view, update, delete)
- Task filtering (all, active, completed)
- OpenTelemetry for monitoring (tracing, logging)

## Technologies Used

### Frontend
- React
- TypeScript
- Vite
- React Router DOM
- React Testing Library, Jest

### Backend
- Go
- Gin Web Framework
- PostgreSQL
- go-sqlmock (for testing)
- OpenTelemetry

## Getting Started

### Prerequisites

- Docker and Docker Compose

### Setup and Run

1.  **Clone the repository:**

    ```bash
    git clone https://github.com/your-repo/todo-with-gemini.git
    cd todo-with-gemini
    ```

2.  **Build and run the Docker containers:**

    ```bash
    docker compose up --build
    ```

    This command will:
    -   Build the frontend and backend Docker images.
    -   Start the PostgreSQL database container.
    -   Run database migrations.
    -   Start the backend API server (on port 8080).
    -   Start the frontend development server (on port 5173).
    -   Start OpenTelemetry Collector and Jaeger.

3.  **Access the application:**

    Open your web browser and navigate to `http://localhost:5173`.

4.  **Access Jaeger UI (for tracing):**

    Open your web browser and navigate to `http://localhost:16686`.

## Running Tests

### Backend Tests

To run backend tests, navigate to the `backend` directory and run:

```bash
cd backend
go test ./...
```

### Frontend Tests

To run frontend tests, navigate to the `frontend` directory and run:

```bash
cd frontend
npm test
```

## Project Structure

```
.
├── backend/                # Go backend application
│   ├── internal/
│   │   ├── auth/           # Authentication handlers
│   │   ├── db/             # Database connection
│   │   ├── logging/        # Structured logging with slog
│   │   ├── middleware/     # JWT authentication middleware
│   │   ├── models/         # Data models (User, Task)
│   │   ├── tasks/          # Task management handlers
│   │   └── telemetry/      # OpenTelemetry setup
│   ├── migrations/         # Database migration SQL files
│   └── main.go             # Main application entry point
├── frontend/               # React TypeScript frontend application
│   ├── public/
│   ├── src/
│   │   ├── components/     # Reusable UI components
│   │   ├── context/        # React Context for authentication
│   │   ├── pages/          # Page components
│   │   ├── services/       # API interaction logic
│   │   └── App.tsx         # Main application component
│   ├── index.css
│   ├── main.tsx
│   ├── jest.config.ts
│   └── setupTests.ts
├── docker-compose.yml      # Docker Compose configuration
├── otel-collector-config.yaml # OpenTelemetry Collector configuration
└── README.md               # Project documentation
```
