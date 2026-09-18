# Thoughts

A distraction-free social media app built around one thing: **reading people's thoughts.**

No endless feeds engineered around engagement. No noise. No pressure to constantly interact.

Just open Thoughts and read what people are thinking.

## ✨ Why Thoughts?

Most social platforms are designed to keep you scrolling.

**Thoughts is designed to let you read.**

The idea is simple:

> Write a thought.
> Share it.
> Read other people's thoughts.
> Move on.

The experience focuses on the content itself rather than likes, notifications, recommendations, and other distractions.

## 🚀 Features

* 📝 Create and share thoughts
* 📖 Read thoughts from other people
* 👤 User profiles
* 👥 Follow and unfollow users
* 🏷️ Tag thoughts
* 🔎 Search thoughts
* 📜 Paginated feed
* 💬 Comments
* ✏️ Edit and delete your thoughts
* 🔄 Optimistic concurrency for updates
* 📚 OpenAPI API specification
* ⚡ Generated Go API types and server interfaces

## 🛠️ Tech Stack

### Backend

* **Go**
* **Chi** — HTTP router
* **PostgreSQL** — database
* **pgx** — PostgreSQL driver
* **golang-migrate** — database migrations
* **oapi-codegen** — OpenAPI → Go code generation

### API Documentation

* **OpenAPI** — API specification
* **Scalar** — interactive API documentation
* **GitHub Pages** — documentation hosting
* **GitHub Actions** — automatic documentation deployment

### Development

* **Podman** — containers
* **gofumpt** — formatting
* **golangci-lint** — linting
* **xh** — HTTP API testing

## 📁 Project Structure

```text
.
├── api/
│   ├── openapi.yaml
│   └── oapi-codegen.yaml
│
├── cmd/
│   ├── api/
│   └── migrate/
│
├── gen/
│   └── api/
│
├── internal/
│   ├── db/
│   ├── env/
│   └── store/
│
├── scripts/
├── web/
│
├── docker-compose.yaml
├── Makefile
└── go.mod
```

## 🔌 API

The API is defined using OpenAPI and Go server code is generated using `oapi-codegen`.

The generated server interface is implemented directly by the application's HTTP layer.

Example:

```go
var _ api.ServerInterface = (*application)(nil)
```

This keeps the OpenAPI contract and the Go implementation in sync.

### API Documentation

📚 **[Interactive API Documentation](https://venkat1abhinav.github.io/thoughts/)**

The documentation is automatically deployed with GitHub Actions whenever the OpenAPI specification changes.

## 🏃 Running Locally

Clone the repository:

```bash
git clone https://github.com/Venkat1abhinav/thoughts.git
cd thoughts
```

Start PostgreSQL:

```bash
podman compose up -d
```

Run migrations:

```bash
make migrate
```

Start the API:

```bash
make run
```

The API will be available at:

```text
http://localhost:3000
```

Health check:

```bash
xh GET :3000/v1/health
```

## 🧬 Generate API Code

After modifying `api/openapi.yaml`:

```bash
make generate
```

This regenerates:

```text
gen/api/api.gen.go
```

## 🧪 Testing

Run tests:

```bash
go test ./...
```

Run formatting:

```bash
gofumpt -w .
```

Run static analysis:

```bash
go vet ./...
golangci-lint run
```

## 🗺️ Philosophy

Thoughts isn't trying to compete with social platforms by adding more features.

It's intentionally **small**.

The goal is to create a place where the primary interaction is simply:

**read.**

Everything else should support that experience rather than compete with it.

## 📌 Status

Thoughts is currently under active development.

The backend API is being built first, with the focus on:

* clean HTTP APIs
* PostgreSQL data modeling
* concurrency
* API contracts
* OpenAPI
* production-oriented Go architecture

## 📄 License

License information will be added as the project develops.
