# CLAUDE.md — SignForce Agent Memory

## Project Identity
- **Name:** SignForce
- **Type:** SaaS platform for advanced electronic signatures
- **Legal framework:** NOM-151-SCFI-2016 (Mexico)
- **Owner:** Juan Núñez (majority), Leo (minority stake)
- **Repo:** https://github.com/xpressiceover0/signforce

## Stack
- **Backend API (sfback):** Go 1.24, MySQL 8.0, JWT auth, AES-256-GCM encryption
- **Middleware (sfmiddle):** Go 1.23, request routing, session management, document flow
- **AI Agent (sfia):** Python 3.10+, TensorFlow, biometric verification, OCR
- **Email Service (emailServ):** Go, SMTP
- **Frontend (sffront):** Static HTML/CSS/JS, Tailwind CDN, DM Sans font
- **Infrastructure:** Docker, Nginx, Let's Encrypt SSL
- **Server:** Skynet (187.77.6.110), Ubuntu 24

## Design System
- Brand color: `#B5C413` (R181 G196 B19)
- Background: `#060a14` (deep black)
- Glass panels: `rgba(15,21,36,0.45)` + `backdrop-filter: blur(16px)`
- Font: DM Sans
- Single source of truth: `sffront/signforce-global.css`

## Critical Rules — DO NOT BREAK

### Ports
- Frontend container: **8090:80** — NEVER change to 8080 or any other port
- sfback API: 5001
- sfmiddle gateway: 5002
- Nginx proxy expects 127.0.0.1:8090

### Files NEVER modify without explicit permission
- `sffront/documentFlow/document_viewer.html` — PDF viewer, complex canvas logic
- `sffront/documentFlow/document_viewer.js` — PDF rendering engine
- Any `.env` file — contains production secrets
- `sfback/db/conexion.go` — database layer (has known SQL injection issues, fix requires careful migration)

### Files safe to modify freely
- `sffront/**/*.html` — all frontend pages
- `sffront/**/*.css` — all stylesheets
- `sffront/signforce-global.css` — design system
- `sffront/**/*.js` — frontend JavaScript (preserve element IDs)
- `docs/**` — documentation
- `README.md`, `CLAUDE.md`, `PROJECT_INDEX.md`

## Update Procedure
```bash
cd ~/signforce && git pull && docker compose -f docker-compose.frontend.yml up -d --build
```

## Environment Variables (see .env.example)
DB_USER, DB_PASS, DB_IP, DB_NAME, JWT_KEY, JWT_EXP, DOCS_KEY,
BASE_DIR, GENERIC_FOLDER_PATH, KEYS_PATH, QR_URL_BASE,
EMAIL_USER, EMAIL_PASS, SMTP_SERVER, API_PORT, BACK_URL, MAX_UPLOAD_MEM

## Known Issues
1. SQL Injection in GenericInsert/GenericSelect (sfback/db, sfmiddle/db, emailServ/db) — uses fmt.Sprintf
2. 247 debug print statements in Go code
3. 104 console.log in production JS
4. 76 alert() calls instead of UI notifications
5. 28 fetch() calls without error handling
6. No rate limiting on API
7. Payment table stores card data in plain text
8. HTTP services without TLS (handled by Nginx proxy for now)

## Architecture Flow
```
Browser → Nginx (SSL :443) → signforce-web (:8090)
                            → sfmiddle (:5002) → sfback (:5001) → MySQL
                                               → sfia (Python)
                                               → emailServ
```

## Commit Convention
```
type: description

Types: feat, fix, rewrite, security, cleanup, polish, docs
Example: fix: upload_doc sidebar alignment and gap spacing
```

## Frontend Patterns
- Every page includes `signforce-global.css` as last stylesheet
- Every page has `<style>html,body{background:#060a14 !important;}</style>` before `</body>`
- Sidebar pattern: two cards (logo card + nav card) with green active state
- Header: full-width, fixed, blur backdrop, S logo badge
- Cards: glass panels with gradient top-line
- Buttons: `#B5C413` primary, transparent ghost secondary
