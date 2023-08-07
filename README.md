# GoChat
WebSocket multi-room chat built with Go & Next.js

Techs stack:
* Backend: Go, Gin, Gorilla WebSocket
* Frontend: TypeScript, Next.js, Tailwind CSS

# Requirements & Setup
* Go 1.16+
* Node.js 14.6+

1. Install backend dependencies
    > go mod download
2. Install frontend dependencies
    > npm install
3. Setup database, e.g. Postgres from Docker image
    > docker pull postgres
4. Create BD table
    ```sql
    CREATE TABLE "users" (
        "id" bigserial PRIMARY KEY,
        "username" varchar NOT NULL,
        "email" varchar NOT NULL UNIQUE,
        "password" varchar NOT NULL
    )
    ```

# Running
1. Start backend:
    > go run cmd/main.go
2. Start frontend:
    > npm run dev
3. Open <http://localhost:3000/> in browser

# Backend architecture
Code is divided into 4 layers according to the Clean Architecture:
* Handler - serves incoming requests (REST, gRPC, WebSocket, GraphQL)
* Service - business logic
* Repository - abstract storage (DB, in-mem, file, other services)
* Entity - data structure

Chat
* Client - WebSocket connection with some user data
* Room - contains a collection of clients and broadcasts messages
* Hub - collection of rooms
