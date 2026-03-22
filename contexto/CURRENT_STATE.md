# CURRENT_STATE.md — Estado Actual del Proyecto
**Última actualización:** 6 de Marzo de 2026

---

## Resumen Ejecutivo
Dashboard LuxSpace v2.0 en producción en `https://dashboard.luxspace.org`.
PostgreSQL migrado y funcionando. Base de datos limpiada y optimizada.
13 usuarios con 8 roles distintos. Blancos Diarios implementado.

---

## Cambios Aplicados (orden cronológico)

| # | Cambio | Descripción | Método |
|---|---|---|---|
| 1 | Dashboard v2 | Base: Next.js 14, JWT (Leo/Axl/Admin), 6 módulos CRUD, KPIs, Recharts | Script .sh |
| 2 | Módulo PDF | Parser PDFs Airbnb, impuesto → 6% | Script .sh |
| 3 | Watcher PDF | PM2 watcher para PDFs desde carpeta SSH | Script .sh |
| 4 | UI v3 | Selector periodo compacto, formato $1,203.54, dropdowns propiedades | Script .sh |
| 5 | Dashboard v4 | LuxSpace + Granada dual, Admin→Soporte (canEdit:true) | Script .sh |
| 6 | Fix datos | Fix campos Sheet, fecha_devolucion, mapeo nombres Airbnb | Script .sh |
| 7 | Habitaciones v5 | Columna Habitación, filtro multi-select, friendlyRoomName() | Script .sh |
| 8 | Fix isGranada | isGranada() robusta: startsWith cdmx/cmdx/granada | Manual |
| 9 | Periodo módulos | PeriodSelector en todas las pestañas | Manual |
| 10 | Fix nombres | parseAirbnbName() con acentos/typos, toInternalRoomCode() | Manual |
| 11 | UI Apple Design | Inter, glassmorphism, dark mode, logos, SVG icons, anti-zoom | Manual |
| 12 | Fix Proveedores | WhatsApp +52, botón Llamar, formatPhone() | Manual |
| 13 | Fix PDF Parser | pdf-parse reinstalado, 4 estrategias fecha, validación | Manual |
| 14 | Limpieza Sheet | Eliminación 27 filas con fecha "0-00-01" | Script Node.js |
| 15 | Dominio + SSL | dashboard.luxspace.org, Certbot, Nginx SSL | Manual |
| 16 | Seguridad | UFW (22/80/443), Fail2ban, cookies secure | Manual |
| 17 | Monitoreo | Uptime Kuma (Docker) en status.luxspace.org | Manual |
| 18 | Backups | Cron diario → /root/backups | Manual |
| 19 | Fix middleware | Redirects con X-Forwarded headers | Manual |
| 20 | UI Redesign v25+ | KPI colores, bottom dock + botón "+", botón fecha flotante, toggle LS/GR, pull-to-refresh, modales detalle | Claude chat |
| 21 | OpenClaw Bot | Docker, Claude 3.5 Sonnet, WhatsApp, skills | Docker |
| 22 | PostgreSQL | Migración a Docker (luxspace-postgres), 103 registros, roles, variables .env | Claude Code |
| 23 | Nuevos usuarios | 8 roles (atencion, limpieza_merida, limpieza_granada, lavanderia, mantenimiento x2) + 8 usuarios nuevos | Claude Code |
| 24 | Blancos Diarios | `/api/blancos`, `BlancosDiariosModule.tsx`, integración con Spreadsheet Apps Script | Claude Code |
| 25 | Módulos nuevos | CuentasForm, PropiedadesModule, PropiedadesContext, api/propiedades | Claude Code |
| 26 | DB Limpieza | Eliminación columnas vacías (mes, anio, created_at_sheets, raw_pdf_data), 9 índices creados | pgAdmin SQL |
| 27 | pgAdmin | Docker container + bd.luxspace.org con SSL | Manual |

---

## Estado de Funcionalidades

### ✅ Funcionando Correctamente

**Infraestructura**
- [x] `https://dashboard.luxspace.org` — SSL, Nginx, Next.js :3000
- [x] `https://bd.luxspace.org` — pgAdmin 4 con SSL (cert: `/etc/letsencrypt/live/bd.luxspace.org/`)
- [x] `https://status.luxspace.org` — Uptime Kuma
- [x] PostgreSQL 16 en Docker (`luxspace-postgres`, red `luxspace-net`, puerto `127.0.0.1:5432`)
- [x] PM2: `dashboard` (online) + `pdf-watcher` (online)
- [x] UFW, Fail2ban, backups cron diario en `/root/backups`

**Autenticación**
- [x] JWT 7 días en cookie `auth-token`
- [x] 13 usuarios con 8 roles distintos
- [x] Contraseñas por variable de entorno (patrón `PrimerNombre123` para operaciones)

**Dashboard & KPIs**
- [x] KPIs LuxSpace + Granada por periodo
- [x] Toggle LS/GR como pills en header
- [x] KPI cards clickeables → modales detalle con desglose
- [x] Retiros: desglose por usuario con conteo
- [x] Gráficas Recharts (barras + pie)
- [x] Desglose por usuario (Leo/Axl/Soporte)

**Blancos Diarios**
- [x] `BlancosDiariosModule.tsx` con 3 tabs (Mérida HOY, MAÑANA, Granada HOY)
- [x] `/api/blancos` lee el spreadsheet separado `11tbthRn5U7...`
- [x] Checkboxes de confirmación (estado en memoria)
- [x] Barra de progreso X/Y listas
- [x] Permisos por rol (limpieza_merida solo ve Mérida, etc.)
- [x] Apps Script genera los datos a las 7:30 AM y 3:00 PM

**CRUD Módulos**
- [x] Ingresos, Gastos, Gastos Fijos, Devoluciones, Retiros, Proveedores
- [x] Cuentas (nuevo)
- [x] Propiedades (nuevo, con Context)
- [x] PDF Upload (browser + watcher)
- [x] Auto-generación de IDs, preservación de valores "0" en formularios

**PostgreSQL**
- [x] 7 tablas + 4 vistas
- [x] 9 índices de performance creados
- [x] Columnas vacías eliminadas (mes, anio, created_at_sheets, raw_pdf_data en ingresos/gastos/devoluciones)
- [x] pgAdmin 4 accesible en bd.luxspace.org

### ⚠️ Pendiente / Por Migrar

- [ ] **Migrar cálculos del dashboard de Google Sheets → PostgreSQL** — `calculations.ts` aún lee Sheets; PostgreSQL está poblado pero no es la fuente activa de los KPIs
- [ ] **Re-subir PDFs de Granada (Mayo-Diciembre 2025)** — Las 27 filas con fecha rota fueron eliminadas; hay que re-cargar desde la UI
- [ ] **Verificar `monto_total`** — Columna existe en PostgreSQL con datos pero el código no la usa; evaluar si eliminar
- [ ] **Redirect post-login por rol** — Los roles de limpieza llegan al dashboard completo en vez de ir directo a Blancos Diarios
- [ ] **Vista standalone para limpieza** — Ruta `/blancos` planeada pero no implementada en page.tsx aún

### 🔧 Mejoras Sugeridas
- Paginación en tablas con muchos registros
- Exportar datos a Excel/CSV
- Dashboard acumulado (varios meses)
- Migración completa Google Sheets → PostgreSQL como fuente única
- Integración OpenClaw ↔ Dashboard (consultas financieras por WhatsApp)

---

## Datos en Google Sheets (post-limpieza)

### Ingresos (46 filas)
- Periodo: Mayo 2025 → Enero 2026
- Suma total depositado: $969,097.86
- 23 propiedades distintas, fuente única: "Airbnb"
- **Granada Mayo-Dic 2025 aún por re-subir** (fueron eliminadas por fecha rota)

### Gastos, Devoluciones, Retiros
- Datos de Enero 2026 principalmente
- `origen_fondos` ya está en el header de Retiros_personales ✅

---

## Dependencias Clave (package.json v2.0.0)

| Paquete | Versión |
|---|---|
| next | 14.2.5 |
| react | 18.3.1 |
| recharts | 2.12.7 |
| googleapis | 140.0.1 |
| jose | 5.6.3 |
| pdf-parse | 2.4.5 |
| pdfjs-dist | 4.0.379 |
| pg | 8.20.0 |
| tailwindcss | 3.4.6 |
| typescript | 5.5.3 |

---

## Servidor Skynet — Estado

| Servicio | Estado | URL |
|---|---|---|
| Next.js Dashboard | ✅ Online (PM2 id:1) | dashboard.luxspace.org |
| PDF Watcher | ✅ Online (PM2 id:0) | — |
| PostgreSQL | ✅ Healthy (Docker) | localhost:5432 |
| pgAdmin | ✅ Up (Docker) | bd.luxspace.org |
| OpenClaw Bot | ✅ Up (Docker :59715) | WhatsApp |
| Uptime Kuma | ✅ Healthy (Docker) | status.luxspace.org |
| Nginx | ✅ Activo | 443→3000, 443→5050 |
| UFW | ✅ Activo | Solo 22/80/443 |
| Fail2ban | ✅ Activo | — |
| Backups | ✅ Cron diario | /root/backups |
