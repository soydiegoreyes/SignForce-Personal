# PROJECT_INDEX.md — SignForce Repository Map

## Root Files
| File | Purpose |
|------|---------|
| `CLAUDE.md` | AI agent memory — rules, stack, conventions |
| `README.md` | Project documentation — architecture, deployment |
| `PROJECT_INDEX.md` | This file — repository map |
| `.env.example` | Environment variables template (15 vars) |
| `.gitignore` | Git exclusions (binaries, secrets, IDE) |
| `docker-compose.yml` | Full stack (MySQL + Go services + Nginx) |
| `docker-compose.frontend.yml` | Frontend only (Nginx on port 8090) |
| `Dockerfile.frontend` | Multi-stage build for static frontend |
| `nginx.conf` | Nginx config (security headers, gzip, cache) |
| `nginx-full.conf` | Nginx config with API proxy routes |
| `signforce_db.sql` | MySQL schema + seed data |

## sffront/ — Frontend (35 HTML pages)
```
sffront/
├── signforce-global.css        ← Design system (single source of truth)
├── favicon.svg                 ← SVG favicon (green S)
├── 404.html                    ← Custom error page
├── index.html                  ← Landing redirect
├── index/
│   └── index.html              ← Landing page (hero, features, CTA)
├── registro/                   ← Auth & onboarding
│   ├── login.html              ← Split-panel login/register
│   ├── contracts.html          ← Contract list with sign buttons
│   ├── validation.html         ← 4-step registration validation
│   ├── validacion.html         ← Document validation
│   ├── edit_user.html          ← User profile editor
│   ├── invite_user.html        ← Team invitation
│   ├── plan_pay.html           ← Payment/plan selection
│   ├── noauth.html             ← Unauthorized access page
│   └── waitapprove.html        ← Pending approval screen
├── dashboards/                 ← Admin panels
│   ├── dashboard_admin.html    ← Admin dashboard (stats, users, activity)
│   ├── dashboard_root.html     ← Root dashboard (team management)
│   ├── dashboard_user.html     ← User dashboard (team member view)
│   └── dashboard_approvals.html ← Client approval workflow
├── documentFlow/               ← Core document signing workflow
│   ├── upload_doc.html         ← Upload + configure document
│   ├── add_signers.html        ← Add signers with search
│   ├── add_signers_.html       ← Alternate signer flow
│   ├── add_signs.html          ← Place signature boxes on PDF
│   ├── document_viewer.html    ← PDF viewer (DO NOT MODIFY)
│   ├── my_documents.html       ← User's document list
│   ├── documents.html          ← Document browser
│   ├── folders.html            ← Folder management
│   ├── view_invite.html        ← Signing invitation (single doc)
│   ├── view_invite_.html       ← Signing invitation (multi doc)
│   ├── view_signature.html     ← Signature verification
│   ├── template_man.html       ← Template manager
│   ├── upload_zone.html        ← Drag & drop upload
│   └── email_sign_request.html ← Email signing request
├── administracion/             ← Admin tools
│   ├── mykeys.html             ← RSA key management
│   ├── uploadkeys.html         ← Key upload form
│   ├── gestion_usuarios.html   ← User management (CRUD)
│   └── crear_equipo.html       ← Team creation
└── images/                     ← Static images
```

## sfback/ — Cryptographic API (Go :5001)
```
sfback/
├── main.go                     ← Entry point, 6 HTTP endpoints
├── Dockerfile                  ← Container build
├── auth/
│   ├── auth.go                 ← Authentication logic
│   └── jwt.go                  ← JWT token generation/validation
├── configs/confs.go            ← Configuration loading
├── db/conexion.go              ← MySQL connection + generic queries ⚠️ SQL injection
├── docflow/functions.go        ← Document processing functions
├── models/models.go            ← Data models
├── objects/
│   ├── keys.go                 ← RSA key management
│   ├── signtypes.go            ← XAdES signature types
│   └── user.go                 ← User operations
├── utilities/
│   ├── cryptoutils.go          ← AES-256-GCM encryption/decryption
│   ├── maps.go                 ← Map utilities
│   └── misc.go                 ← Helper functions
├── pades.py                    ← PAdES signature (Python)
└── padeshanko.py               ← PAdES with Hanko (Python)
```

## sfmiddle/ — Middleware Gateway (Go :5002)
```
sfmiddle/
├── main.go                     ← Entry point, 30 HTTP endpoints
├── auth/                       ← Auth middleware
├── coms/communications.go      ← Email/notification dispatch
├── configs/confs.go            ← Configuration
├── db/conexion.go              ← MySQL connection ⚠️ SQL injection
├── documentflow/
│   ├── folders.go              ← Folder operations
│   └── invites.go              ← Invitation management
├── handlers/
│   ├── dashboard.go            ← Dashboard data endpoints
│   ├── docflow.go              ← Document workflow handlers
│   ├── documents.go            ← Document CRUD
│   ├── keys.go                 ← Key management endpoints
│   ├── registro.go             ← Registration flow
│   └── users.go                ← User management
├── iapackage/iafunctions.go    ← AI integration bridge
├── models/                     ← Shared data models
├── objects/                    ← Business objects
└── utilities/                  ← Shared utilities
```

## sfia/ — AI Agent (Python)
```
sfia/
├── agent.py                    ← Main agent entry point
├── amiga_llama.py              ← LLM integration
├── crypto_utils.py             ← AES decryption (reads DOCS_KEY env)
├── doc_cache.py                ← Document caching
└── pipeline.py                 ← Processing pipeline
```

## emailServ/ — Email Service (Go)
```
emailServ/
├── mailservice.go              ← SMTP email sender
├── logerstats.go               ← Logging
└── db/conexion.go              ← MySQL connection
```

## docs/ — Documentation
```
docs/
├── architecture.md             ← System architecture
├── system-overview.md          ← What SignForce does
├── development-workflow.md     ← How to develop
├── deployment.md               ← How to deploy
└── decisions.md                ← Technical decisions log
```

## knowledge/ — AI Knowledge Base
```
knowledge/
├── project-knowledge.md        ← Accumulated project insights
├── bug-fixes.md                ← Bug fix history
└── patterns.md                 ← Code patterns and conventions
```
