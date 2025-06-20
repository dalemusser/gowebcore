# gowebcore

Reusable Go web-server toolkit: config, logging, HTTP server, middleware, database helpers, background tasks, and authentication adapters.

| Layer | What you get |
|-------|--------------|
| **config** | Flag + file + env loader (`viper`), merged into a single struct. |
| **logger** | Structured `slog` JSON with file:line and automatic OpenTelemetry trace/span IDs. |
| **server** | Chi router, gzip, request-ID, graceful shutdown, and HTTPS via **static PEM _or_ Let’s Encrypt**. |
| **middleware** | CORS (config-driven), CSRF cookie, Prometheus request histogram, `/metrics` route. |
| **render** | Embedded HTML templates + HTMX fragment detection. |
| **asset** | Fingerprinted `/assets/*` with immutable caching, `asset.Path("app.css")` helper. |
| **auth** | JWT (HS256 & RS256) + API-key middleware. |
| **db** | Multi-alias Mongo/DocumentDB manager + generic helpers. |
| **tasks** | Background job runner with graceful stop & Prometheus gauge. |
| **queue** | Redis Streams producer/consumer adapter (SQS & NATS coming). |
| **aws** | S3 presign (PUT/GET) and CloudFront invalidation. |

---
## Local dev stack

```bash
docker compose up -d  # Redis, Postgres, NATS, Prometheus
go run ./examples/clever
```
---

## Quick start

```bash
# create a new service
go mod init github.com/me/myservice
go get github.com/dalemusser/gowebcore@latest
```

```go
// cmd/server/main.go
package main

import (
	"context"

	"github.com/dalemusser/gowebcore/config"
	"github.com/dalemusser/gowebcore/logger"
	"github.com/dalemusser/gowebcore/middleware"
	"github.com/dalemusser/gowebcore/server"
	"github.com/go-chi/chi/v5"
)

type cfg struct{ config.Base }

func main() {
	var c cfg
	_ = config.Load(&c, config.WithEnvPrefix("APP"))

	logger.Init(c.LogLevel)

	r := chi.NewRouter()
	r.Use(middleware.CORSFromConfig(c.Base))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, gowebcore!"))
	})

	srv := server.New(c.Base, r)
	_ = server.Serve(context.Background(), srv, c.CertFile, c.KeyFile)
}
```

```toml
# config.toml
app_name   = "demo"
env        = "dev"
http_port  = 8080
domain     = "demo.local"
enable_tls = false           # autocert if true & domain set
log_level  = "info"

cors_origins = [
  "http://localhost:5173"
]
```

Run it:

```bash
go run ./cmd/server --config=config.toml
```

## Example — minimal service that just uses DefaultCORS()

If you’re prototyping or building a fully public API, you can keep the
wide-open wildcard behavior:

```go
package main

import (
	"context"
	"net/http"

	"github.com/dalemusser/gowebcore/config"
	"github.com/dalemusser/gowebcore/logger"
	"github.com/dalemusser/gowebcore/middleware"
	"github.com/dalemusser/gowebcore/server"
	"github.com/go-chi/chi/v5"
)

type cfg struct{ config.Base }

func main() {
	var c cfg
	_ = config.Load(&c)      // flags / env / file

	logger.Init(c.LogLevel)

	r := chi.NewRouter()
	r.Use(middleware.DefaultCORS())   // ← allows every Origin "*"

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("public API\n"))
	})

	srv := server.New(c.Base, r)
	_ = server.Serve(context.Background(), srv, "", "") // HTTP by default
}
```

DefaultCORS() sets AllowedOrigins to *, plus standard methods and
headers.

- Ideal for internal services, mock APIs, or auth-protected endpoints where origin restrictions aren’t necessary.
- For production front-ends you’ll usually switch to middleware.CORSFromConfig(cfg.Base) so you can whitelist specific domains.

## Background jobs & Redis queue

```go
import (
	"github.com/dalemusser/gowebcore/queue/redis"
	"github.com/dalemusser/gowebcore/tasks"
)

rdb, _ := redis.New("redis://localhost:6379")
producer := redis.NewProducer(rdb)
_ = producer.Publish(ctx, "email", []byte(`{"to":"alice","subj":"hi"}`))

mgr := tasks.New()
redis.NewConsumer(rdb, "email", "workers", emailHandler).Start(mgr)
mgr.Start(ctx)
```

## Metrics & tracing

```go
middleware.RegisterDefaultPrometheus()

shutdown, _ := observability.Init(ctx, observability.Config{
	ServiceName: "demo",
	Endpoint:    "localhost:4318",  // OTLP/HTTP collector
	SampleRatio: 0.25,
})
defer shutdown(ctx)
```

/metrics exposes Go runtime stats and HTTP request histograms.

Every log line includes trace_id & span_id when a span is active.

⸻

Roadmap
- SQS & NATS adapters under queue/
- CLI toolkit (serve, migrate, --version) via cobra
- Deployment recipes: Dockerfile + Kubernetes manifests
- Docs site with copy-paste snippets

____

## How static files are served in gowebcore

gowebcore itself does not auto-register any static route.

It ships the asset package:

```go
import "github.com/dalemusser/gowebcore/asset"

// asset.Handler()  →  http.Handler that serves /assets/* from the embedded FS
// asset.Path("app.css") → "/assets/app.f3c9e2.css"
```

Your service decides if / where to mount it.

Typical usage (as shown in the example app):

```go
r := chi.NewRouter()
r.Mount("/assets", asset.Handler())    // <─ add the static route

r.Get("/", pageHandler)
```

Want additional static folders (e.g., docs, user uploads)?

You add them explicitly:

```go
r.Handle("/docs/*", http.StripPrefix("/docs/", http.FileServer(http.Dir("./docs"))))
```

Why it’s done this way

- Keeps the core opinion-free about URL layout.
- Services that don’t need embedded assets (pure JSON APIs, queue workers) don’t pay any cost.
- You remain free to host large binaries or user uploads on S3/CloudFront instead of inside the Go binary.

So: no automatic static route—you mount asset.Handler() (or any other file server) exactly where you want it.


## CLI worker

### Background worker

```bash
myservice worker --config=prod.toml
```

Auto-registers queue consumers (Redis Streams, SQS, NATS) if configured.

Supports cron tasks via tasks.Cron:

```golang
mgr.Cron("0 */5 * * *", func(ctx context.Context) error {
    return db.Exec(ctx, "DELETE FROM sessions WHERE expires < NOW()")
})
```

## Supported authentication options

| Auth method            | Protocol / Library                                   | Example path\*           | Required env vars†                                                    | Uses `SESSION_HASH_KEY`? |
|------------------------|------------------------------------------------------|--------------------------|-----------------------------------------------------------------------|--------------------------|
| **ClassLink (SAML)**   | SAML 2.0 — `crewjam/saml/samlsp`                     | `examples/classlink_saml` | `CL_SUBDOMAIN`                                                        | ❌ |
| **ClassLink (OIDC)**   | OpenID Connect — `coreos/go-oidc`                    | `examples/classlink_oidc` | `CL_SUBDOMAIN`, `CL_OIDC_CLIENT_ID`, `CL_OIDC_CLIENT_SECRET`          | ✅ |
| **Clever**             | OAuth 2.0 — Clever REST                              | `examples/clever`        | `CLEVER_CLIENT_ID`, `CLEVER_CLIENT_SECRET`, `SESSION_HASH_KEY`        | ✅ |
| **Google Login**       | OpenID Connect — `golang.org/x/oauth2` (Google)      | `examples/google_login`  | `GOOGLE_CLIENT_ID`, `GOOGLE_CLIENT_SECRET`, `SESSION_HASH_KEY`        | ✅ |
| **GitHub OAuth**       | OAuth 2.0 — `golang.org/x/oauth2` (GitHub)           | `examples/github_oauth`  | `GITHUB_CLIENT_ID`, `GITHUB_CLIENT_SECRET`, `SESSION_HASH_KEY`        | ✅ |

\* *Example paths refer to runnable programs in the `examples/` directory.*  
† Rows marked **✅** require the secure session hash key below.

---

### Generating a secure session hash key

`auth.Session` signs its cookie with a 64-byte key (optionally encrypting with
a 32-byte block key). Generate once and load it from your secret store:

```bash
# generate
openssl rand -hex 64 > HASH_KEY

# load for local dev
export SESSION_HASH_KEY=$(cat HASH_KEY)

## Clever Authentication

```
#### Clever

```toml
[clever]
client_id     = "YOUR_ID"
client_secret = "YOUR_SECRET"
redirect      = "https://svc.example.com/auth/callback"
```

```go
provider := oauth.NewClever(cfg.Clever.ClientID, cfg.Clever.ClientSecret, cfg.Clever.Redirect)
middleware.Routes(r, provider, session)
```

## ClassLink (OIDC) Authentication

```bash
export CL_SUBDOMAIN="mydistrict"
export CL_OIDC_CLIENT_ID="..."
export CL_OIDC_CLIENT_SECRET="..."
go run ./examples/classlink_oidc
```

Open http://localhost:8080/ – you’ll be redirected to ClassLink’s OIDC login
and back to the protected home page.

## ClassLink (SAML) Authentication

```bash
export CL_SUBDOMAIN="mydistrict"
go run ./examples/classlink_saml
```

Note: The SAML example auto-generates an RSA key and self-signed
certificate when keyFile / certFile are nil. For production supply real
key/cert files or use Let’s Encrypt and set SAML_KEY_FILE,
SAML_CERT_FILE.



© 2025 Dale Musser & contributors. MIT License.

