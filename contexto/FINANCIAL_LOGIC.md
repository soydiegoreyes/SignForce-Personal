# FINANCIAL_LOGIC.md — Lógica Financiera del Dashboard
**Última actualización:** 6 de Marzo de 2026

## Archivo Principal
`src/lib/calculations.ts` → función `calculateDashboard(mes, anio)`

Lee **5 hojas de Google Sheets** en paralelo y calcula todos los KPIs.
**Fuente activa:** Google Sheets (PostgreSQL está poblado pero aún no es la fuente del dashboard).

---

## Contabilidad Dual: LuxSpace vs Granada

### LuxSpace (principal)
Todas las propiedades EXCEPTO Granada: ROSADA, AKNAI, INAI, TINY, SANTIAGO, CASA DE PIEDRA, CARLOTA, ROMINA.

### Granada (separada)
Solo la propiedad **GRANADA** (CDMX). Balance completamente independiente.

### Función isGranada()
```typescript
function isGranada(raw: string): boolean {
  if (!raw) return false;
  const lower = raw.toLowerCase().trim();
  if (lower === "granada") return true;
  if (lower.startsWith("cdmx")) return true;
  if (lower.startsWith("cmdx")) return true;  // typo support
  if (lower.startsWith("granada")) return true;
  return false;
}
```

---

## Fórmulas LuxSpace

| KPI | Fórmula | Fuente |
|---|---|---|
| **Ingresos Brutos** | `SUM(total_depositado + saldo_pendiente)` donde `fecha_ingreso` matchea mes/año Y `isGranada(propiedad_id) = false` | Hoja: Ingresos |
| **Devoluciones** | `SUM(monto)` donde `fecha_devolucion` matchea Y no es Granada | Hoja: Devoluciones |
| **Impuestos** | `Ingresos Brutos × 0.06` | Calculado |
| **Gastos Variables** | `SUM(monto)` donde `fecha_gasto` matchea Y no es Granada | Hoja: Gastos |
| **Gastos Fijos** | `SUM(monto)` donde `activo = TRUE/VERDADERO/1` Y no es Granada | Hoja: Gastos_Fijos |
| **Gastos Totales** | `Gastos Variables + Gastos Fijos` | Calculado |
| **Retiros** | `SUM(monto)` donde `fecha` matchea Y `isGranada(origen_fondos ?? propiedad_id) = false` | Hoja: Retiros_personales |
| **Balance Final** | `Ingresos − Devoluciones − Impuestos − Gastos Totales − Retiros` | Calculado |

## Fórmulas Granada (estructura idéntica)

| KPI | Fuente de filtro |
|---|---|
| Ingresos | `isGranada(propiedad_id) = true` |
| Devoluciones | `isGranada(propiedad_id ?? habitacion_id) = true` |
| Gastos Variables | `isGranada(propiedad_id) = true` |
| Gastos Fijos | `isGranada(propiedad_id) = true` |
| Retiros | `isGranada(origen_fondos ?? propiedad_id) = true` |

---

## Tasa de Impuesto
**6%** sobre ingresos brutos (`ingresos * 0.06`).

---

## Desglose por Usuario
Gastos y retiros rastreados por: **Leo**, **Axl**, **Soporte** (match por `usuario_id.toLowerCase()`).
Componente: `UserBreakdown.tsx`.

---

## Lógica de Retiros — Origen de Fondos
```
1. Si row.origen_fondos existe → usar ese valor
2. Si no → usar row.propiedad_id como fallback
3. Aplicar isGranada() al valor resultante
```
En el formulario UI: dos botones "LuxSpace" / "Granada" que guardan en `origen_fondos`.
**✅ La columna `origen_fondos` ya existe en Google Sheets y en PostgreSQL.**

---

## Filtrado por Periodo

| Módulo | Campo fecha | Tipo | Filtro |
|---|---|---|---|
| Dashboard (KPIs) | Todos | Server-side | `calculateDashboard(mes, anio)` |
| Ingresos | `fecha_ingreso` | date | Client-side `useMemo` |
| Gastos | `fecha_gasto` | date | Client-side `useMemo` |
| Devoluciones | `fecha_devolucion` | date | Client-side `useMemo` |
| Retiros | `fecha` | date | Client-side `useMemo` |
| Gastos Fijos | N/A | — | ❌ Muestra todos (recurrentes) |
| Proveedores | N/A | — | ❌ Sin fecha |

### Función parseDate() — soporta múltiples formatos
```typescript
function parseDate(dateStr: string): { month: number; year: number } | null {
  // 1. YYYY-MM-DD (principal)
  // 2. MM/DD/YYYY
  // 3. new Date() fallback
}
```

---

## Gastos Fijos — Valores activo válidos
```typescript
const a = (r.activo || "").toUpperCase();
if (a === "TRUE" || a === "VERDADERO" || a === "1") { /* suma */ }
```
Esto maneja tanto valores booleanos de PostgreSQL como strings de Google Sheets.

---

## Ingresos por Habitación
El desglose usa `friendlyRoomName(propiedad_id)` para convertir el valor del Sheet al nombre amigable.
El objeto `ingresos_by_property` usa nombres amigables como keys:
```json
{ "AKNAI 1": 26999.21, "GRANADA 3": 19819.34 }
```

---

## Interfaz DashboardData Completa
```typescript
export interface DashboardData {
  periodo_visible: string;  // "Enero 2026"
  mes: number;
  anio: number;

  // ═══ LuxSpace ═══
  total_ingresos_brutos: number;
  total_devoluciones: number;
  impuestos: number;
  gastos_variables: number;
  gastos_fijos: number;
  gastos_totales: number;
  gastos_leo: number;
  gastos_axl: number;
  gastos_soporte: number;
  retiros_totales: number;
  retiros_leo: number;
  retiros_axl: number;
  retiros_soporte: number;
  balance_final: number;
  ingresos_by_property: Record<string, number>;
  gastos_by_category: Record<string, number>;
  gastos_fijos_detail: { concepto: string; monto: number }[];

  // ═══ Granada ═══
  granada_ingresos: number;
  granada_devoluciones: number;
  granada_impuestos: number;
  granada_gastos_variables: number;
  granada_gastos_fijos: number;
  granada_gastos_totales: number;
  granada_retiros: number;
  granada_retiros_leo: number;
  granada_retiros_axl: number;
  granada_retiros_soporte: number;
  granada_balance: number;
  granada_gastos_by_category: Record<string, number>;
  granada_ingresos_by_room: Record<string, number>;

  // ═══ Detalles para modales ═══
  devoluciones_detail: { habitacion: string; monto: number }[];
  retiros_detail: { usuario: string; monto: number; origen: string }[];
  granada_devoluciones_detail: { habitacion: string; monto: number }[];
  granada_retiros_detail: { usuario: string; monto: number }[];
}
```

---

## Formato de Moneda
```typescript
formatCurrency(amount) → "$1,203.54"
// Intl.NumberFormat("es-MX", { style: "currency", currency: "MXN", 2 decimals })
```
