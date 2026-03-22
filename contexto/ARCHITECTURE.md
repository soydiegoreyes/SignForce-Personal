# ARCHITECTURE.md — Servidor Skynet: Dashboard LuxSpace + OpenClaw Bot
**Última actualización:** 6 de Marzo de 2026

## Descripción General
Servidor **Skynet** aloja dos proyectos principales:

1. **Dashboard LuxSpace** — App web financiera para administrar un negocio de rentas tipo Airbnb. Registra ingresos, gastos, devoluciones, retiros y proveedores. Incluye carga automática de PDFs de Airbnb y contabilidad separada para **Granada (CDMX)** vs todas las demás propiedades (LuxSpace).

2. **OpenClaw Bot** — Asistente IA integrado con WhatsApp (Claude 3.5 Sonnet) para gestión operativa y administrativa del negocio.

---

## Stack Tecnológico — Dashboard LuxSpace
| Capa | Tecnología | Versión |
|---|---|---|
| Frontend | Next.js (App Router) | 14.2.5 |
| UI | React | 18.3.1 |
| Estilos | Tailwind CSS | 3.4.6 |
| Tipografía | Inter (Google Fonts) | latest |
| Gráficas | Recharts | 2.12.7 |
| Backend | Next.js API Routes | 14 |
| Base de datos primaria | PostgreSQL 16 (Docker) | 16-alpine |
| Base de datos legado | Google Sheets API v4 | — |
| Autenticación | JWT (jose 5.6.3) | 5.6.3 |
| PDF parsing | pdf-parse 2.4.5 + pdfjs-dist 4.0.379 | — |
| Process Manager | PM2 | latest |
| Web Server | Nginx (reverse proxy) | latest |
| OS | Ubuntu 24 (VPS srv1349130) | — |
| Node.js | 22.x (v22.22.0) | LTS |

## Stack Tecnológico — OpenClaw Bot
| Capa | Tecnología | Detalle |
|---|---|---|
| Motor IA | Claude 3.5 Sonnet | `anthropic/claude-3-5-sonnet-20241022` |
| Integración | WhatsApp Web | Vínculo vía QR |
| Despliegue | Docker | Container `openclaw-t4m8-openclaw-1` |
| Puerto externo | 59715 | Hostinger HVPS OpenClaw |

---

## Diseño Visual
- **Color primario (brand):** Luxeberry `#460479`
- **Estilo:** Apple-inspired con glassmorphism, transparencias, backdrop-blur
- **Tipografía:** Inter (sans-serif, similar a SF Pro)
- **Modo oscuro:** Toggle en header (icono sol/luna), estado guardado en localStorage, clase `html.dark`
- **Clase `.glass`:** `background: rgba(255,255,255,0.82)` + `backdrop-filter: blur(20px) saturate(180%)`
- **Dark glass:** `background: rgba(30,30,32,0.72)` + misma blur
- **Bottom dock:** Glass pill con 3 iconos de navegación (Retiros, Dashboard, Gastos) + botón circular "+" separado con gradiente Luxeberry
- **Logos:** Dos versiones en `/public/`: `LuxSpace-logo-oscuro.png` (modo claro) / `LuxSpace-logo-claro.png` (modo oscuro)
- **Anti-zoom mobile:** `viewport: userScalable: false, maximumScale: 1` + `touch-action: manipulation`
- **Iconos:** SVG inline (Heroicons style), NO emojis en la UI
- **Scrollbar:** Oculta globalmente (webkit + Firefox)
- **Pull-to-refresh:** Touch gesture nativo con spinner Luxeberry animado

## Estructura del Proyecto
```
~/airbnb-dashboard/
├── src/
│   ├── app/
│   │   ├── api/
│   │   │   ├── auth/route.ts              ← Login/logout/session (POST/GET/DELETE)
│   │   │   ├── blancos/route.ts           ← GET blancos diarios desde Sheet separado
│   │   │   ├── dashboard/route.ts         ← GET KPIs financieros por periodo
│   │   │   ├── health/route.ts            ← Health check (incluye ?db=1 para PostgreSQL)
│   │   │   ├── pending-pdfs/route.ts      ← PDFs del watcher (carpeta servidor)
│   │   │   ├── propiedades/route.ts       ← CRUD propiedades (nuevo)
│   │   │   ├── sheets/route.ts            ← CRUD genérico Google Sheets
│   │   │   └── upload-pdf/route.ts        ← Parse + confirm PDFs Airbnb
│   │   ├── login/page.tsx
│   │   ├── globals.css
│   │   ├── layout.tsx
│   │   └── page.tsx                       ← Página principal
│   ├── components/
│   │   ├── BlancosDiariosModule.tsx        ← Módulo blancos diarios (nuevo)
│   │   ├── FinancialCharts.tsx
│   │   ├── KPICard.tsx
│   │   ├── PeriodSelector.tsx
│   │   ├── PropiedadesModule.tsx           ← CRUD propiedades (nuevo)
│   │   ├── UserBreakdown.tsx
│   │   ├── forms/
│   │   │   ├── CuentasForm.tsx            ← CRUD Cuentas (nuevo)
│   │   │   ├── DevolucionesForm.tsx
│   │   │   ├── FormField.tsx
│   │   │   ├── GastosFijosForm.tsx
│   │   │   ├── GastosForm.tsx
│   │   │   ├── IngresosForm.tsx
│   │   │   ├── PdfUploader.tsx
│   │   │   ├── ProveedoresForm.tsx
│   │   │   └── RetirosForm.tsx
│   │   ├── tables/
│   │   │   └── DataTable.tsx
│   │   └── ui/
│   │       ├── Modal.tsx
│   │       ├── Tabs.tsx
│   │       └── Toast.tsx
│   ├── lib/
│   │   ├── auth.ts                        ← JWT + 13 usuarios con roles
│   │   ├── calculations.ts                ← Lógica financiera (LuxSpace + Granada)
│   │   ├── formatters.ts                  ← Formato moneda, propiedades, habitaciones
│   │   └── PropiedadesContext.tsx         ← Context React para propiedades (nuevo)
│   ├── middleware.ts                       ← Protección rutas JWT
│   └── services/
│       ├── airbnbPdfParser.ts
│       └── googleSheets.ts
├── uploads/
│   ├── pendientes/
│   ├── procesados/
│   └── errores/
├── public/
│   ├── LuxSpace-logo-oscuro.png
│   └── LuxSpace-logo-claro.png
├── watcher.js
├── .env.local
├── middleware.ts
├── next.config.js
├── package.json
├── tailwind.config.ts
└── tsconfig.json

~/openclaw-bot/
└── (Docker Compose, Hostinger HVPS OpenClaw)
```

## Flujo de Datos
```
[Usuario Web] → [Next.js Frontend (React 18)]
                        │
                        ├── GET /api/dashboard?mes=X&anio=Y
                        │       → calculations.ts → googleSheets.ts (5 hojas)
                        │       → retorna DashboardData (KPIs LuxSpace + Granada)
                        │
                        ├── CRUD /api/sheets?sheet=NombreHoja
                        │       → googleSheets.ts → Google Sheets API v4
                        │
                        ├── GET /api/blancos
                        │       → Google Sheets API → Spreadsheet Blancos Diarios
                        │       → Hojas: "Merida HOY", "Merida MAÑANA", "Granada Hoy"
                        │
                        ├── POST /api/upload-pdf (parse / confirm)
                        │       → airbnbPdfParser.ts → googleSheets.ts
                        │
                        └── /api/auth → auth.ts (JWT cookie 7 días)

[SSH/SFTP Upload] → uploads/pendientes/ → watcher.js (PM2) → pending-results.json

[WhatsApp] → OpenClaw Bot (Docker :59715) → Claude 3.5 Sonnet → WhatsApp
```

## Procesos PM2
| id | Proceso | Archivo | Puerto | Estado |
|---|---|---|---|---|
| 1 | `dashboard` | `npm start` (Next.js) | 3000 | online |
| 0 | `pdf-watcher` | `watcher.js` | — | online |

## Docker Containers
| Container | Imagen | Puerto | Estado |
|---|---|---|---|
| `luxspace-postgres` | postgres:16-alpine | 127.0.0.1:5432→5432 | Up (healthy) |
| `pgadmin` | dpage/pgadmin4 | 127.0.0.1:5050→80 | Up |
| `openclaw-t4m8-openclaw-1` | ghcr.io/hostinger/hvps-openclaw | 0.0.0.0:59715→59715 | Up |
| `uptime-kuma` | louislam/uptime-kuma | 127.0.0.1:3001→3001 | Up (healthy) |

## Mapa de Puertos
| Servicio | Puerto Interno | Puerto Público |
|---|---|---|
| Next.js Dashboard | 3000 | 443 (vía Nginx) |
| PostgreSQL | 5432 | Solo localhost |
| pgAdmin | 5050 | 443 → bd.luxspace.org |
| OpenClaw | 59715 | 59715 |
| Uptime Kuma | 3001 | localhost only |

## Dominios y DNS (Namecheap)
| Subdominio | Destino | Servicio |
|---|---|---|
| `dashboard.luxspace.org` | 187.77.6.110 | Next.js App |
| `status.luxspace.org` | 187.77.6.110 | Uptime Kuma |
| `bd.luxspace.org` | 187.77.6.110 | pgAdmin 4 |

SSL: Let's Encrypt (Certbot, auto-renew)  
Certificados: `/etc/letsencrypt/live/luxspace.org/` y `/etc/letsencrypt/live/bd.luxspace.org/`

## Variables de Entorno (.env.local)
```env
# Google Sheets
GOOGLE_SHEETS_SPREADSHEET_ID=1XRHQvTm_4Hcy3R8O7lpgDQ0Z3UCLKMSJ7jpmxCY53Rw
GOOGLE_SERVICE_ACCOUNT_KEY={"type":"service_account",...}

# JWT
JWT_SECRET=<random hex 64 chars>

# Usuarios maestros
USER_LEO_PASS=<pass>
USER_AXL_PASS=<pass>

# Atención/Soporte
USER_ADMIN_PASS=<pass>       ← usado por: soporte
USER_MARLENE_PASS=<pass>     ← usado por: marlene
USER_MAURICIO_PASS=<pass>    ← usado por: mauricio

# Limpieza Mérida
USER_MARCELO_PASS=Marcelo123
USER_LUZ_PASS=Luz123
USER_ISELA_PASS=Isela123

# Limpieza Granada
USER_YAEL_PASS=Yael123
USER_JAVIER_PASS=Javier123

# Lavandería
USER_FABIOLA_PASS=Fabi123

# Mantenimiento
USER_JORGE_PASS=Jorge123
USER_VICTOR_PASS=Victor123

# Blancos Diarios (spreadsheet separado)
BLANCOS_SPREADSHEET_ID=11tbthRn5U7_-Uc9zmec-JKmQB_W4pjhtOYXR_zVhnmo

# PostgreSQL
POSTGRES_HOST=luxspace-postgres
POSTGRES_PORT=5432
POSTGRES_DB=luxspace
POSTGRES_USER=luxspace_admin
POSTGRES_READONLY_USER=luxspace_readonly
POSTGRES_WRITER_USER=luxspace_writer
POSTGRES_PASSWORD=<pass>
```

## Seguridad
- **Firewall (UFW):** Solo puertos 22, 80, 443
- **Fail2ban:** Activo
- **Cookies:** `auth-token` con `httpOnly: true`, `secure: true`, `sameSite: "lax"`

## Acceso al Servidor
- **Nombre:** Skynet
- **IP pública:** 187.77.6.110
- **Cliente SSH:** Terminus
- **Directorio dashboard:** `~/airbnb-dashboard`
- **Directorio OpenClaw:** `~/openclaw-bot`

## Comandos Útiles
```bash
# Dashboard
cd ~/airbnb-dashboard
pm2 status
pm2 logs dashboard
pm2 restart dashboard
npm run build && pm2 restart dashboard

# Docker
docker ps
docker logs luxspace-postgres
docker exec -it luxspace-postgres psql -U luxspace_admin -d luxspace

# Nginx
nginx -t && systemctl reload nginx
certbot renew --dry-run
```
