# Forum

A simple web forum written in Go.

## Features

- User registration and login with hashed passwords
- Session management via cookies
- Create posts with categories
- Comment on posts
- Like and dislike posts and comments
- Filter posts by category, by your posts or by liked posts
- SQLite database
- Dockerfile for containerization

## Running locally

```bash
go run ./cmd/web
```

The server listens on `:8080` and stores data in `forum.db`.

## Testing

```bash
go test ./...
```

## Docker

```bash
docker build -t forum .
docker run -p 8080:8080 forum
```
