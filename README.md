
# RestWithGo

REST API in Go for managing users and posts (educational example).

Summary: This small API provides endpoints for user sign-up, login, retrieving the authenticated user, and basic CRUD for posts. It is intended as a learning project and a starting point for more complex REST services.

Project structure:
- `main.go`: application entry point and configuration loading.
- `server/`: server initialization and route binding.
- `handlers/`: HTTP handlers (users, posts, home).
- `repository/`: persistence interface and repository wiring.
- `database/`: Dockerfile and database initialization (Postgres).

Quick requirements:
- Go 1.16+ installed.
- Docker (optional, to run the example database).

Contents:
- Prerequisites
- Installation & run
- Environment variables
- Endpoints
- curl examples
- Database schema
- Authentication notes

Prerequisites
- Install `go` (recommended 1.16+).
- (Optional) Install `docker` to use the provided Postgres image.

Installation & run

1. Clone the repository and enter the project folder:

```bash
git clone <repo-url>
cd "Curso 4 - Rest y Mock Services"
```

2. (Optional) Build the database image and run a Postgres container for local testing:

```bash
# from the `database/` folder
docker build . -t platzi-rs-ws-db

# run the DB with example credentials
docker run --rm --name platzi-rs-ws-db -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=postgres -e POSTGRES_DB=rest -p 5432:5432 platzi-rs-ws-db
```

3. Create a `.env` file in the repository root with the required variables (example):

```
PORT=":8080"
JWT_SECRET=change_this_secret
DATABASE_URL=postgres://postgres:postgres@localhost:5432/rest?sslmode=disable
```

4. Run the API:

```bash
# run directly
go run main.go

# or build and run binary
go build -o rest && ./rest
```

Environment variables
- `PORT`: server port to listen on (include leading `:`). Example: `:8080`.
- `JWT_SECRET`: secret key used to sign JWT tokens.
- `DATABASE_URL`: Postgres connection string, e.g. `postgres://user:pass@host:5432/dbname?sslmode=disable`.

Endpoints

All endpoints accept and return JSON. Some routes require a JWT token in the `Authorization` header (the raw token string, without `Bearer` prefix).

- GET /
	- Description: Welcome / health check
	- Authentication: No

- POST /signup
	- Description: Register a new user
	- Body: `{ "email": "user@example.com", "password": "secret" }`
	- Response: `{ "email": "user@example.com", "id": "<user-id>" }`
	- Authentication: No

- POST /login
	- Description: Login and obtain a JWT token
	- Body: `{ "email": "user@example.com", "password": "secret" }`
	- Response: `{ "token": "<jwt-token>" }`
	- Authentication: No

- GET /me
	- Description: Get the authenticated user's data
	- Headers: `Authorization: <token>`
	- Authentication: Yes

- POST /posts
	- Description: Create a new post for the authenticated user
	- Headers: `Authorization: <token>`
	- Body: `{ "post_content": "Hello world" }`
	- Response: `{ "id": "<post-id>", "post_content": "..." }`

- GET /posts
	- Description: List posts (paginated)
	- Query params: `page` (optional, starts at `0`). Page size = 2.
	- Authentication: Yes

- GET /posts/{id}
	- Description: Get a post by id
	- Authentication: Yes

- PUT /posts/{id}
	- Description: Update a post's content (only if it belongs to the authenticated user)
	- Headers: `Authorization: <token>`
	- Body: `{ "post_content": "New content" }`

- DELETE /posts/{id}
	- Description: Delete a post if it belongs to the authenticated user
	- Headers: `Authorization: <token>`

Authentication note
- The application expects the JWT token in the `Authorization` header as the raw token (no `Bearer` prefix). Example:

```
Authorization: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

This applies both to the middleware and to handlers that verify permissions.

curl examples

- Sign up:

```bash
curl -X POST http://localhost:8080/signup \
	-H "Content-Type: application/json" \
	-d '{"email":"dev@example.com","password":"pass1234"}'
```

- Login and get token:

```bash
curl -X POST http://localhost:8080/login \
	-H "Content-Type: application/json" \
	-d '{"email":"dev@example.com","password":"pass1234"}'

# the response contains: {"token": "..."}
```

- Use token to create a post:

```bash
TOKEN="<obtained_token>"
curl -X POST http://localhost:8080/posts \
	-H "Content-Type: application/json" \
	-H "Authorization: $TOKEN" \
	-d '{"post_content":"Hello from curl"}'
```

Database schema

See `database/up.sql` (summary):

- `users` table:
	- `id` VARCHAR(32) PRIMARY KEY
	- `password` VARCHAR(255) NOT NULL
	- `email` VARCHAR(255) UNIQUE NOT NULL
	- `created_at` TIMESTAMP NOT NULL DEFAULT NOW()

- `posts` table:
	- `id` VARCHAR(32) PRIMARY KEY
	- `post_content` VARCHAR(32) NOT NULL
	- `created_at` TIMESTAMP NOT NULL DEFAULT NOW()
	- `user_id` VARCHAR(32) NOT NULL (FK -> users.id)

Developer notes
- Routes are bound in `main.go` (function `BindRoutes`).
- The middleware `CheckAuthMiddleware` allows unauthenticated access to `signup` and `login`, and requires a token for other routes.
- The repository uses `github.com/lib/pq` as the Postgres driver.

Testing & integration
- There are no automated tests included by default. To test locally:
	1. Start Postgres (docker run or local service).
	2. Populate `.env` and run `go run main.go`.
	3. Use the curl examples to validate the flow: signup -> login -> create post -> list.

Contributing
- Open an issue or submit a pull request with improvements (better error handling, configurable pagination, input validation, etc.).

License
- Example repository — check the course author policy for usage and distribution.

**Changelog (recent commits)**

The repository received several updates in the last commits. Quick summary:

- feat: web socket upgrade (commit f47b725)
	- Added real-time WebSocket support via a new `websocket` package (`websocket/hub.go`, `websocket/client.go`).
	- New model `WebSocketMessage` (`models/message.go`) used to wrap broadcast messages.
	- Server integration: the server now creates a `Hub` and exposes a `/ws` endpoint to accept WebSocket clients.
	- When a new post is created (`POST /posts`) the server broadcasts a message of type `"Post_Created"` containing the created post as payload.
	- Added `ListPost` support (repository and DB) and paginated listing handler.

- chore: Dockerfile to production environment (commit c8cdf52)
	- Introduced a multi-stage `Dockerfile` to build a static production binary and a minimal runner image.
	- The Dockerfile sets up a builder stage that compiles the Go binary and a scratch-based runner stage.
	- `go.mod` tidied up (dependency declarations adjusted).

- refactor: tidy and fix on Dockerfile (commit 60feb23)
	- Updated the root `Dockerfile` to use a newer Go version and improved Alpine package commands.
	- `go.mod` Go version updated (now `go 1.24.0`).

Notes & usage for the WebSocket feature

- Endpoint: `GET /ws` — upgrade HTTP connection to WebSocket.
- Message format: messages are JSON objects that follow the `WebSocketMessage` model:

```json
{
	"type": "Post_Created",
	"payload": {
		"id": "<post-id>",
		"post_content": "...",
		"user_id": "...",
		"created_at": "..."
	}
}
```

- Example (browser/Node.js):

```js
const ws = new WebSocket('ws://localhost:8080/ws');
ws.onmessage = (evt) => {
	const msg = JSON.parse(evt.data);
	console.log('Received', msg.type, msg.payload);
};
```

- When creating a post via the REST API (`POST /posts`) the server will broadcast a `Post_Created` message to all connected WebSocket clients.

If you want, I can also add a short example client script (Node.js) or a Postman collection demonstrating the new flows.


