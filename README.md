# Library service

A Go HTTP service that manages a book inventory backed by Elasticsearch, with Redis-based activity tracking.

Built with [gin-gonic](https://github.com/gin-gonic/gin).

---

## Running

```bash
export ES_URL="http://your-es-host:9200"
go run ./cmd/service
```

Server listens on port **8080**.

---

## API

All routes accept `?username=<name>` as a query param for activity tracking.

### Books

| Method | Route | Description |
|---|---|---|
| `POST` | `/books` | Add a new book |
| `GET` | `/books/:id` | Get a book by id |
| `PUT` | `/books/:id` | Update a book's title |
| `DELETE` | `/books/:id` | Delete a book |

**POST /books** — request body:
```json
{
  "title": "The Master and Margarita",
  "author_name": "Mikhail Bulgakov",
  "price": 9.99,
  "ebook_available": true,
  "publish_date": "1967-01-01"
}
```

Response: `201 Created`
```json
{"id": "<generated_id>"}
```

---

### Search & Stats

| Method | Route | Params | Description |
|---|---|---|---|
| `GET` | `/search` | `title`, `author_name`, `price_range` | Search books by criteria |
| `GET` | `/store` | — | Total books and distinct authors |

**Search examples:**
```bash
# by author only
curl 'localhost:8080/search?author_name=Tolkien&username=Gilad'

# by price range only
curl 'localhost:8080/search?price_range=5,15&username=Gilad'

# by author + price range
curl 'localhost:8080/search?author_name=Tolkien&price_range=5,20&username=Gilad'

# by title + price range
curl 'localhost:8080/search?title=hobbit&price_range=,15&username=Gilad'
```

`price_range` format: `min,max` — either side is optional (`,20` = up to 20, `5,` = 5 and up).

---

### Activity

| Method | Route | Description |
|---|---|---|
| `GET` | `/activity` | Last 3 actions for a user |

```bash
curl 'localhost:8080/activity?username=Gilad'
```

Response:
```json
{
  "actions": [
    {"method": "GET", "route": "/store"},
    {"method": "GET", "route": "/search"},
    {"method": "POST", "route": "/books"}
  ]
}
```

Activity is stored in Redis with a 24-hour TTL. The `/activity` route itself is never recorded.

---

## Architecture

Follows the **Repository pattern** — all storage logic is isolated behind interfaces, keeping HTTP handlers independent of the underlying database.

```
cmd/service/main.go
pkg/
  config/                   ← env-based configuration
  book/
    model.go                ← Book struct
    repository.go           ← Repository interface
    es_repository.go        ← Elasticsearch implementation
  activity/
    repository.go           ← Repository interface
    redis_repository.go     ← Redis implementation
  service/                  ← HTTP layer (routes, handlers, middleware)
    routes.go
    handlers.go
    middleware.go
```
