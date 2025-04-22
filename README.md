# User Management & Real-Time Chat App

This is a **User Management** system with integrated **Real-Time Chat** built using **Go (Fiber)** on the backend and **Vue 3** with **Pinia** on the frontend. Users can register, log in, see other users, and chat with them in real time.

The app supports:

- Secure **JWT-based authentication**
- Manage user data using **RESTful APIs**
- Live chat using **WebSockets**

---

## Architecture

This app uses the **MVC pattern**, which stands for **Model-View-Controller**. It's a way to organize code so everything is neat, easy to manage, and flexible for future changes.

- **Model (M)**: Defined in the backend using GORM (e.g., `User`, `Message`) and persisted in a **PostgreSQL database**.
- **Controller (C)**: Go backend handlers using Fiber (REST API + WebSocket)
- **View (V)**: Vue 3 frontend UI with Pinia store for state management

---

## Getting Started (Docker Compose)

### Prerequisites

Before you start, you’ll need to have the following tools installed:

- **[Docker](https://www.docker.com/get-started/)**
- **[Docker Compose](https://docs.docker.com/compose/install/)**

Make sure you have both Docker and Docker Compose installed before proceeding.

### Installation

1. **Clone the repository**

   ```bash
   git clone https://github.com/aotsurasak46/user-management.git
   cd user-management
   ```

2. **Environment Configuration**

Copy the example environment file and edit values if needed (do it both .env and ./frontend/.env):

```bash
cp .env.example .env
```

```bash
cd ./frontend
cp .env.example .env
```

3. **Run the project**

   ```bash
   docker-compose up --build
   ```

   This will start:

- Backend API → http://localhost:8080
- Frontend → http://localhost:80
- PostgreSQL → localhost:5432

## API Documentation

Visit Swagger docs after running the app:

```bash
docker-compose up --build
```
