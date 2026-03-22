# DATABASE_SCHEMA.md — Base de Datos LuxSpace
**Última actualización:** 6 de Marzo de 2026

## Arquitectura de Datos Dual

El sistema usa **dos fuentes de datos en paralelo**:

| Fuente | Uso | Acceso |
|---|---|---|
| **PostgreSQL 16** (Docker `luxspace-postgres`) | Base de datos principal. Contiene todos los datos migrados + nuevos | `pg` npm package, `POSTGRES_*` env vars |
| **Google Sheets API v4** | Fuente legacy activa. El dashboard aún lee/escribe aquí para el CRUD | `googleapis` npm, `GOOGLE_SERVICE_ACCOUNT_KEY` |

**Estado actual:** El código de `calculations.ts` y `sheets/route.ts` sigue usando Google Sheets. PostgreSQL está poblado (103 registros migrados) pero aún no es la fuente principal del dashboard.

---

## PostgreSQL — Conexión
```
Host:     luxspace-postgres (Docker network: luxspace-net)
Puerto:   127.0.0.1:5432 (solo localhost)
DB:       luxspace
Red:      luxspace-net
```

### Roles PostgreSQL
| Rol | Uso |
|---|---|
| `luxspace_admin` | Acceso total, owner de tablas |
| `luxspace_readonly` | Solo lectura (n8n, bots) |
| `luxspace_writer` | Escritura sin DDL (automatizaciones) |

---

## Google Sheets — Spreadsheets

| Spreadsheet | ID | Uso |
|---|---|---|
| **Principal (financiero)** | `1XRHQvTm_4Hcy3R8O7lpgDQ0Z3UCLKMSJ7jpmxCY53Rw` | Ingresos, Gastos, Devoluciones, Retiros, Proveedores, Gastos_Fijos, Cuentas |
| **Blancos Diarios** | `11tbthRn5U7_-Uc9zmec-JKmQB_W4pjhtOYXR_zVhnmo` | Merida HOY, Merida MAÑANA, Granada Hoy, Habitaciones (iCal URLs) |

---

## Tablas PostgreSQL (schema real post-limpieza)

### Tabla: ingresos
| Columna | Tipo | Notas |
|---|---|---|
| `id` | integer PK | Auto-increment |
| `ingreso_id` | varchar | ID corto generado con `generateId()` |
| `fecha_ingreso` | date | Formato YYYY-MM-DD |
| `propiedad_id` | varchar | Nombre amigable: "AKNAI 1", "GRANADA 3" |
| `fuente` | varchar | "Airbnb", "Booking", "Externo" |
| `total_depositado` | numeric(12,2) | Monto principal |
| `saldo_pendiente` | numeric(12,2) | Monto pendiente |
| `ingresos_brutos_airbnb` | numeric(12,2) | Del PDF Airbnb |
| `ajustes_airbnb` | numeric(12,2) | Del PDF Airbnb |
| `tarifa_servicio_airbnb` | numeric(12,2) | Comisión Airbnb |
| `impuesto_6pct` | numeric(12,2) | Calculado del PDF |
| `cuenta_airbnb` | varchar | Nombre del host |
| `monto_total` | numeric(12,2) | ⚠️ Columna extra de migración (tiene datos) |
| `es_granada` | boolean | `true` si propiedad es Granada |
| `created_at` | timestamptz | Auto |
| `updated_at` | timestamptz | Auto |
| `source` | varchar | "google_sheets", "pdf_direct" — origen del registro |

**⚠️ IMPORTANTE:** El cálculo de ingreso en el código usa `total_depositado + saldo_pendiente`, NO `monto_total`.

### Tabla: gastos
| Columna | Tipo | Notas |
|---|---|---|
| `id` | integer PK | — |
| `gasto_id` | varchar | ID corto |
| `fecha_gasto` | date | — |
| `propiedad_id` | varchar | "ROSADA", "GRANADA", "LuxSpace" |
| `categoria` | varchar | Ver lista en formatters.ts |
| `descripcion` | text | — |
| `monto` | numeric | — |
| `usuario_id` | varchar | "Leo", "Axl", "Soporte" |
| `es_granada` | boolean | — |
| `created_at` | timestamptz | — |
| `updated_at` | timestamptz | — |
| `source` | varchar | — |

### Tabla: gastos_fijos
| Columna | Tipo | Notas |
|---|---|---|
| `id` | integer PK | — |
| `gasto_fijo_id` | varchar | "GF-000001" |
| `concepto` | varchar | — |
| `categoria` | varchar | — |
| `periodicidad` | varchar | "mensual", "Quincenal" |
| `propiedad_id` | varchar | — |
| `monto` | numeric | Monto mensual |
| `usuario_id` | varchar | — |
| `activo` | boolean | Solo `true` se suma en KPIs |
| `updated_at_sheets` | varchar | Legado de migración |
| `created_at` | timestamptz | — |
| `updated_at` | timestamptz | — |
| `source` | varchar | — |

### Tabla: devoluciones
| Columna | Tipo | Notas |
|---|---|---|
| `id` | integer PK | — |
| `devolucion_id` | varchar | — |
| `fecha_devolucion` | date | ⚠️ Campo es `fecha_devolucion`, NO `fecha` |
| `habitacion_id` | varchar | Código interno: "M2-2", "CDMX1-1" |
| `propiedad_id` | varchar | Nombre propiedad normalizado |
| `monto` | numeric | — |
| `motivo` | text | — |
| `es_granada` | boolean | — |
| `created_at` | timestamptz | — |
| `updated_at` | timestamptz | — |
| `source` | varchar | — |

### Tabla: retiros_personales
| Columna | Tipo | Notas |
|---|---|---|
| `id` | integer PK | — |
| `retiro_id` | varchar | — |
| `fecha` | date | Campo es `fecha` (sin prefijo) |
| `usuario_id` | varchar | "Leo", "Axl", "Soporte" |
| `propiedad_id` | varchar | Fallback para origen |
| `monto` | numeric | — |
| `nota` | text | — |
| `origen_fondos` | varchar | "LuxSpace" o "Granada" (campo nuevo ✅) |
| `es_granada` | boolean | — |
| `created_at` | timestamptz | — |
| `updated_at` | timestamptz | — |
| `source` | varchar | — |

**Lógica origen:** `origen_fondos` → fallback a `propiedad_id` → `isGranada()`.

### Tabla: proveedores
| Columna | Tipo |
|---|---|
| `id` | integer PK |
| `proveedor_id` | varchar |
| `nombre` | varchar |
| `telefono` | varchar |
| `servicio` | varchar |
| `created_at_sheets` | varchar (legado) |
| `created_at` | timestamptz |
| `updated_at` | timestamptz |
| `source` | varchar |

### Tabla: pdf_imports (log de PDFs)
| Columna | Tipo | Notas |
|---|---|---|
| `id` | integer PK | — |
| `filename` | varchar | Nombre del PDF |
| `host_name` | varchar | Nombre del host Airbnb |
| `mes` | integer | Mes del reporte |
| `anio` | integer | Año del reporte |
| `rooms_parsed` | integer | Habitaciones extraídas |
| `rows_added` | integer | Filas escritas en Sheets |
| `status` | varchar | "success", "error" |
| `error_msg` | text | Si hubo error |
| `raw_parse_result` | jsonb | Resultado completo del parser |
| `processed_by` | varchar | Usuario que procesó |
| `created_at` | timestamptz | — |

---

## Vistas PostgreSQL

| Vista | Descripción |
|---|---|
| `v_balance_mensual` | Balance por año/mes/es_granada: ingresos, devoluciones, impuestos, gastos, retiros, balance_final |
| `v_gastos_por_categoria` | Gastos agrupados por año/mes/es_granada/categoria/usuario |
| `v_ingresos_por_habitacion` | Ingresos agrupados por año/mes/habitacion/es_granada |
| `v_resumen_mensual` | Resumen con total_registros, ingresos_brutos, impuestos_pdf |

---

## Índices (creados post-limpieza)
```sql
idx_ingresos_fecha           ON ingresos(fecha_ingreso)
idx_ingresos_propiedad       ON ingresos(propiedad_id)
idx_ingresos_es_granada      ON ingresos(es_granada)
idx_ingresos_fecha_granada   ON ingresos(fecha_ingreso, es_granada)
idx_gastos_fecha             ON gastos(fecha_gasto)
idx_gastos_es_granada        ON gastos(es_granada)
idx_devoluciones_fecha       ON devoluciones(fecha_devolucion)
idx_devoluciones_es_granada  ON devoluciones(es_granada)
idx_retiros_fecha            ON retiros_personales(fecha)
```

---

## Propiedades y Habitaciones

### 9 Propiedades activas + entidad general

| Propiedad | Código | Ciudad | Habitaciones | Contabilidad |
|---|---|---|---|---|
| ROSADA | M1 | Mérida | 3 | LuxSpace |
| AKNAI | M2 | Mérida | 3 | LuxSpace |
| INAI | M3 | Mérida | 3 | LuxSpace |
| TINY | M4 | Mérida | 2 | LuxSpace |
| SANTIAGO | M5 | Mérida | 1 | LuxSpace |
| CASA DE PIEDRA | M6 | Mérida | 7 | LuxSpace |
| CARLOTA | M7 | Mérida | 5 | LuxSpace |
| ROMINA | M8 | Mérida | 4 | LuxSpace |
| GRANADA | CDMX1 | CDMX | 4 | **Granada (separada)** |
| LuxSpace | — | — | — | Solo gastos generales |

**Nota:** KUKA & NARANJO fue removida de la lista activa en formatters.ts.

### Habitaciones (39 total)
```
ROSADA (M1):         M1-1, M1-2, M1-3
AKNAI (M2):          M2-1, M2-2, M2-3
INAI (M3):           M3-1, M3-2, M3-3
TINY (M4):           M4-1, M4-2
SANTIAGO (M5):       M5
CASA DE PIEDRA (M6): M6-1, M6-2, M6-3, M6-4, M6-A, M6-B, M6-SUITE
CARLOTA (M7):        M7-1, M7-2, M7-3, M7-4, M7-5
ROMINA (M8):         M8-1, M8-2, M8-3, M8-4
GRANADA (CDMX1):     CDMX1-1, CDMX1-2, CDMX1-3, CDMX1-4
```

### Blancos Diarios — Hoja "Habitaciones"
El spreadsheet de Blancos tiene su propia hoja `Habitaciones` con columnas:
`Habitacion | Tamaño de cama | Zona | Url_iCal`

Zonas: `Zona 1` (AKNAI, INAI, TINY), `Zona 2` (ROSADA, SANTIAGO, CASA DE PIEDRA, CARLOTA, ROMINA), `Granada`
