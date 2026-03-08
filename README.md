<div align="center">

# 🟢 SignForce

**Plataforma de Firma Electrónica Avanzada**

[![NOM-151](https://img.shields.io/badge/NOM--151--SCFI--2016-Compliant-B5C413?style=flat-square)](https://www.dof.gob.mx)
[![XAdES](https://img.shields.io/badge/XAdES--BES-ETSI%20EN%20319%20132-60a5fa?style=flat-square)](https://www.etsi.org)
[![ASiC-E](https://img.shields.io/badge/ASiC--E-ETSI%20EN%20319%20162-a78bfa?style=flat-square)](https://www.etsi.org)
[![Go](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat-square&logo=go)](https://go.dev)
[![Python](https://img.shields.io/badge/Python-3.10+-3776AB?style=flat-square&logo=python)](https://python.org)

[Demo en vivo](https://testsignforce.luxspace.org) · [Documentación](#arquitectura) · [Despliegue](#despliegue-rápido)

</div>

---

## 📋 Descripción

SignForce es una plataforma integral para la gestión, firma y protección de documentos electrónicos con **validez legal en México**, conforme a la **NOM-151-SCFI-2016** y el **Código de Comercio (Art. 89 y 97)**.

### Características principales

- **Firma Electrónica Avanzada** — Firma XAdES-BES con certificados X.509 del SAT (e.firma / FIEL)
- **Cifrado AES-256-GCM** — Documentos cifrados en reposo con nonce único por bloque
- **Contenedor ASiC-E** — Empaquetado estándar ETSI con documento original + firmas + manifiesto
- **Verificación biométrica** — Validación facial con prueba de vida para firmantes
- **Flujo multi-firmante** — Invitaciones a múltiples firmantes y observadores con fechas límite
- **Verificación QR** — Código QR único por firma para verificación instantánea de autenticidad
- **Dashboard administrativo** — Gestión de equipos, usuarios, documentos y aprobaciones

---

## 🏗️ Arquitectura

```
┌──────────────────────────────────────────────────────────┐
│                     NGINX (Reverse Proxy + SSL)          │
│                     puerto 80/443                        │
├──────────────┬──────────────┬────────────────────────────┤
│   sffront    │   sfmiddle   │         sfback             │
│   Frontend   │  Middleware   │      API Core             │
│   HTML/CSS   │  Go :5002    │      Go :5001              │
│   Estático   │  Gateway     │  Crypto + DB + Firma       │
├──────────────┴──────────────┴────────────────────────────┤
│                      MySQL 8.0                           │
│                     puerto 3306                          │
├──────────────────────────────────────────────────────────┤
│   emailServ  │     sfia     │                            │
│  Servicio    │   Agente IA  │                            │
│  Email (Go)  │  Python/ML   │                            │
└──────────────┴──────────────┴────────────────────────────┘
```

| Servicio | Tecnología | Puerto | Descripción |
|----------|-----------|--------|-------------|
| `sffront` | HTML/CSS/JS | 80 (Nginx) | Frontend estático con design system unificado |
| `sfmiddle` | Go | 5002 | Middleware/gateway — autenticación, ruteo, sesiones |
| `sfback` | Go | 5001 | API criptográfica — firma XAdES, cifrado AES, certificados |
| `emailServ` | Go | — | Servicio de notificaciones por email |
| `sfia` | Python | — | Agente IA — verificación biométrica, OCR |
| `MySQL` | MySQL 8.0 | 3306 | Base de datos relacional |

---

## ⚡ Despliegue rápido

### Prerequisitos

- Docker y Docker Compose
- Dominio con SSL (Let's Encrypt)

### Solo frontend (demo)

```bash
git clone https://github.com/xpressiceover0/signforce.git
cd signforce
cp .env.example .env  # Editar variables
docker compose -f docker-compose.frontend.yml up -d
```

### Stack completo

```bash
cp .env.example .env  # Configurar MySQL, JWT, etc.
docker compose up -d
```

---

## 🔐 Seguridad

| Capa | Implementación |
|------|---------------|
| Cifrado en reposo | AES-256-GCM con nonce único por bloque |
| Firma digital | RSA + SHA-256 sobre certificados X.509 |
| Formato de firma | XAdES-BES (ETSI EN 319 132) |
| Contenedor | ASiC-E (ETSI EN 319 162) |
| Transporte | TLS 1.3 via Nginx + Let's Encrypt |
| Autenticación | JWT con expiración configurable |
| Verificación | Hash SHA-256 + QR de autenticidad |

---

## 📜 Cumplimiento normativo

- ✅ **NOM-151-SCFI-2016** — Conservación de mensajes de datos y digitalización
- ✅ **Código de Comercio Art. 89 y 97** — Firma electrónica avanzada
- ✅ **ETSI EN 319 132** — Formato XAdES para firmas XML avanzadas
- ✅ **ETSI EN 319 162** — Contenedor ASiC-E para empaquetado de firmas
- ✅ **RFC 3161** — Protocolo de sellado de tiempo (TSP)

---

## 📁 Estructura del proyecto

```
signforce/
├── sfback/                 # API criptográfica (Go)
│   ├── db/                 # Conexión MySQL + queries
│   ├── objects/            # Modelos de datos
│   ├── handlers/           # Endpoints HTTP
│   └── main.go             # Entry point :5001
├── sfmiddle/               # Middleware gateway (Go)
│   └── main.go             # Entry point :5002
├── sffront/                # Frontend estático
│   ├── signforce-global.css # Design system unificado
│   ├── index/              # Landing page
│   ├── registro/           # Login, registro, validación
│   ├── dashboards/         # Admin, root, user, approvals
│   ├── documentFlow/       # Upload, firma, visor, invitaciones
│   ├── administracion/     # Gestión usuarios, llaves, equipos
│   └── 404.html            # Página de error personalizada
├── emailServ/              # Servicio de email (Go)
├── sfia/                   # Agente IA (Python)
├── signforce_db.sql        # Schema MySQL
├── docker-compose.yml      # Stack completo
├── docker-compose.frontend.yml  # Solo frontend
├── nginx.conf              # Configuración Nginx optimizada
├── .env.example            # Variables de entorno requeridas
└── .gitignore              # Exclusiones de git
```

---

## 🎨 Design System

| Token | Valor | Uso |
|-------|-------|-----|
| Background | `#060a14` | Fondo principal (negro profundo) |
| Card glass | `rgba(15,21,36,0.45)` | Paneles con backdrop-blur |
| Brand green | `#B5C413` | Acento principal, botones, badges |
| Text primary | `#ffffff` | Títulos y texto principal |
| Text secondary | `rgba(255,255,255,0.55)` | Texto auxiliar |
| Font | DM Sans | Tipografía principal |

---

## 📄 Licencia

Todos los derechos reservados © 2026 SignForce Technologies.

