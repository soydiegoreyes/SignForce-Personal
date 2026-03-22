# API_SPEC.md — Endpoints del Dashboard
**Última actualización:** 6 de Marzo de 2026

## Base URL
`https://dashboard.luxspace.org` (Nginx SSL 443 → localhost:3000)

## Autenticación
- JWT almacenado en cookie `auth-token`
- Flags: `httpOnly: true`, `secure: true`, `sameSite: "lax"`
- Expira en 7 días
- Library: `jose` v5.6.3
- Middleware: `src/middleware.ts` protege TODAS las rutas excepto `/login`, `/api/auth`, assets estáticos

### Middleware — Redirect Fix
Usa `X-Forwarded-Proto` y `X-Forwarded-Host` (headers de Nginx) para construir redirects.
**⚠️ NUNCA usar `new URL("/login", request.nextUrl)`** — genera `localhost:3000` detrás del proxy.

---

## POST /api/auth — Login
```json
// Request
{ "username": "leo", "password": "xxx" }
// Response 200
{ "success": true, "username": "leo" }
// Set-Cookie: auth-token=<JWT>
// Response 401
{ "error": "Credenciales incorrectas" }
```

## GET /api/auth — Verificar sesión
```json
// Response 200 (autenticado)
{ "authenticated": true, "username": "leo", "role": "leo", "canEdit": true }
// Response 401 (no autenticado)
{ "authenticated": false }
```

## DELETE /api/auth — Logout
Elimina la cookie `auth-token`.

---

## Usuarios del Sistema

| Login | Role | canEdit | Acceso |
|---|---|---|---|
| `leo` | leo | ✅ | Full dashboard |
| `axl` | axl | ✅ | Full dashboard |
| `soporte` | soporte | ✅ | Full dashboard |
| `marlene` | atencion | ✅ | Dock: Devoluciones, Gastos, Proveedores, Blancos |
| `mauricio` | atencion | ✅ | Dock: Devoluciones, Gastos, Proveedores, Blancos |
| `marcelo` | limpieza_merida | ❌ | Solo Blancos Diarios (Mérida) |
| `luz` | limpieza_merida | ❌ | Solo Blancos Diarios (Mérida) |
| `isela` | limpieza_merida | ❌ | Solo Blancos Diarios (Mérida) |
| `yael` | limpieza_granada | ❌ | Solo Blancos Diarios (Granada) |
| `javier` | limpieza_granada | ❌ | Solo Blancos Diarios (Granada) |
| `fabiola` | lavanderia | ❌ | Blancos Diarios (Mérida + Granada) |
| `jorge` | mantenimiento | ❌ | Solo Blancos Diarios (Mérida) |
| `victor` | mantenimiento_granada | ❌ | Solo Blancos Diarios (Granada) |

**Patrón de contraseñas:** `PrimerNombre123` (ej: `Marcelo123`)
**Variable .env para soporte/marlene/mauricio:** `USER_ADMIN_PASS`, `USER_MARLENE_PASS`, `USER_MAURICIO_PASS`

---

## GET /api/dashboard — KPIs Financieros
```
GET /api/dashboard?mes=1&anio=2026
```
Retorna `DashboardData` completo. Lee 5 hojas de Google Sheets en paralelo.
Ver `FINANCIAL_LOGIC.md` para estructura completa.

---

## /api/sheets — CRUD Genérico Google Sheets

Hojas permitidas: `Ingresos`, `Gastos`, `Gastos_Fijos`, `Devoluciones`, `Retiros_personales`, `Proveedores`, `Cuentas`

### GET
```
GET /api/sheets?sheet=Ingresos
```
```json
{ "headers": ["ingreso_id","fecha_ingreso",...], "rows": [...], "count": 46 }
```

### POST — Crear fila
```json
{ "sheet": "Ingresos", "values": { "fecha_ingreso": "2026-01-15", ... } }
// Auto-genera ID si está vacío (usando generateId())
// Auto-asigna created_at si está vacío
```

### PUT — Editar fila
```json
{ "sheet": "Ingresos", "rowIndex": 3, "values": { ... } }
// rowIndex es 0-indexed desde la primera fila de datos
```

### DELETE — Eliminar fila
```json
{ "sheet": "Ingresos", "rowIndex": 3 }
```

**Permisos:** POST, PUT, DELETE requieren `session.canEdit === true`.
**Fix importante:** Usa `values[h] != null ? String(values[h]) : ""` para preservar valores "0".

---

## GET /api/blancos — Blancos Diarios
Lee el Spreadsheet separado de Blancos Diarios (generado por Google Apps Script).

```
GET /api/blancos
```
```json
{
  "merida_hoy": {
    "titulo": "HOY Viernes 6 de Marzo",
    "zonas": [
      { "zona": "Zona 1", "habitaciones": ["Aknai 1 - King - 2 noches", ...] },
      { "zona": "Zona 2", "habitaciones": ["Rosada 2 - Matrimonial - 1 noche", ...] }
    ]
  },
  "merida_manana": { ... },
  "granada_hoy": {
    "titulo": "HOY Viernes 6 de Marzo",
    "zonas": [{ "zona": "", "habitaciones": ["Granada 1 - 1 noche", ...] }]
  },
  "fetched_at": "2026-03-06T18:00:00Z"
}
```

**Spreadsheet ID:** `11tbthRn5U7_-Uc9zmec-JKmQB_W4pjhtOYXR_zVhnmo`  
**Hojas leídas:** `Merida HOY`, `Merida MAÑANA`, `Granada Hoy`  
**Generado por:** Google Apps Script (corre a las 7:30 AM y 3:00 PM diariamente)

---

## /api/propiedades — CRUD Propiedades (nuevo)
Endpoint dedicado para gestión de propiedades.

---

## POST /api/upload-pdf — Carga de PDFs Airbnb

### Acción: parse
```
POST /api/upload-pdf (multipart: action=parse, pdfs=<archivo>)
```
```json
{
  "results": [{
    "filename": "enero-2026.pdf",
    "host": "Diego Leonardo Reyes Chavez",
    "mes": 1, "anio": 2026,
    "rooms": [{
      "alojamiento": "AKNAI 1",
      "ingresos_brutos": 28517.71,
      "ajustes": 0,
      "tarifa_servicio": -842.15,
      "total_neto": 25335.71,
      "impuesto_6pct": 1520.14,
      "ingreso_final": 23815.57
    }]
  }]
}
```

### Acción: confirm
```
POST /api/upload-pdf (multipart: action=confirm, data=<JSON>)
```
```json
{ "success": true, "rowsAdded": 5 }
```

---

## GET /api/health — Health Check
```
GET /api/health        → { "status": "ok" }
GET /api/health?db=1   → { "status": "ok", "db": "connected", "rows": 46 }
```

---

## /api/pending-pdfs — PDFs del Watcher

| Método | Descripción |
|---|---|
| GET | Ver resultados pendientes de confirmar |
| POST `{ reports: [...] }` | Confirmar y guardar en Sheets |
| DELETE | Limpiar pendientes |
