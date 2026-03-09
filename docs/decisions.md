# Technical Decisions — SignForce

## Decision Log

### 2026-03-08: Design System Unified (#B5C413 + #060a14)
- **Context:** 35 HTML pages had inconsistent styles (mix of Tailwind, inline, external CSS)
- **Decision:** Single CSS file (signforce-global.css) with !important overrides
- **Rationale:** Tailwind CDN generates classes that override custom styles. Using !important at end of body ensures consistency.
- **Trade-off:** Harder to override individual elements, but consistency is critical for the demo/presentation.

### 2026-03-08: AES Key Moved to Environment Variable
- **Context:** AES-256 encryption key was hardcoded in sfia/crypto_utils.py
- **Decision:** Replace with `os.getenv("DOCS_KEY")`
- **Rationale:** CRITICAL security vulnerability. Key was visible in git history.
- **Action needed:** Set DOCS_KEY in .env before deploying full stack. The old key value must be used for backward compatibility with existing encrypted files.

### 2026-03-08: Frontend-Only Docker Deployment
- **Context:** Full stack (MySQL + Go) requires significant configuration
- **Decision:** Deploy only Nginx container serving static files for demo
- **Rationale:** Frontend demo is sufficient for presentation. Backend endpoints return errors but pages render correctly.
- **Future:** Deploy full stack when ready for production testing.

### 2026-03-08: Port 8090 for Frontend Container
- **Context:** Port 8080 was already in use by another service on Skynet
- **Decision:** Use 8090:80 for the signforce-web container
- **Consequence:** NEVER change this port. Nginx host config depends on it.

### 2026-01: Static Frontend (No Framework)
- **Context:** Choice between React/Vue/vanilla HTML
- **Decision:** Static HTML + Tailwind CDN + vanilla JS
- **Rationale:** Faster iteration, no build step, easier deployment. Team familiarity with HTML/CSS.
- **Trade-off:** Duplicated header/sidebar code across pages, harder to maintain.

### 2026-01: Go for Backend Services
- **Context:** Need high-performance cryptographic operations
- **Decision:** Go with standard library HTTP server
- **Rationale:** Go's crypto/x509, crypto/aes packages are production-grade. Good concurrency for handling multiple signing operations.

### 2026-01: MySQL over PostgreSQL
- **Context:** Database choice for document/signature audit trail
- **Decision:** MySQL 8.0
- **Rationale:** Team familiarity, good tooling, sufficient for the use case.

### 2026-01: XAdES-BES Signature Format
- **Context:** Legal requirement for electronic signatures in Mexico
- **Decision:** Implement XAdES-BES (ETSI EN 319 132)
- **Note:** NOM-151 ideally requires XAdES-T (with timestamp). Current implementation is XAdES-BES only. Timestamp server integration is pending.
