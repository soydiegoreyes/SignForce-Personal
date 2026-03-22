# UI_SPEC.md — Especificación de Interfaz de Usuario
**Última actualización:** 6 de Marzo de 2026

## Diseño General
- **Color primario (brand):** Luxeberry `#460479`
- **Estilo:** Apple-inspired con glassmorphism (transparencias + backdrop-blur)
- **Modo Oscuro:** Toggle en header (sol/luna), guarda en localStorage
- **Mobile-first:** modal como bottom sheet en mobile, centrado en desktop
- **Anti-zoom:** viewport `userScalable: false`, inputs `font-size: 16px`
- **Iconos:** SVG inline (Heroicons), NO emojis en ninguna parte de la UI
- **Scrollbar:** Oculta globalmente (webkit + Firefox)
- **Pull-to-refresh:** Touch gesture con spinner Luxeberry en todas las páginas

### Tipografía
Inter (Google Fonts), importado en `layout.tsx`.

### Modo Claro
- Body: `#f5f5f7`
- Glass: `rgba(255,255,255,0.82)` + `backdrop-filter: blur(20px) saturate(180%)`
- Texto: `#1d1d1f`

### Modo Oscuro
- Body: `#000000`
- Glass: `rgba(30,30,32,0.72)` + misma blur
- Texto: `#f5f5f7`

### Logo
```html
<img class="dark:hidden" src="/LuxSpace-logo-oscuro.png">    <!-- Modo claro -->
<img class="hidden dark:block" src="/LuxSpace-logo-claro.png"> <!-- Modo oscuro -->
```

---

## Roles y Vistas

### Roles con acceso al dashboard completo
`leo`, `axl`, `soporte` — Ven: header completo, bottom dock (Retiros/Dashboard/Gastos), botón "+"

### Rol `atencion` (Marlene, Mauricio)
Dock simplificado: **Devoluciones, Gastos, Proveedores, Blancos Diarios**. Sin botón "+".

### Roles de limpieza/operaciones
`limpieza_merida`, `limpieza_granada`, `lavanderia`, `mantenimiento`, `mantenimiento_granada`
Vista exclusiva: **Solo Blancos Diarios** según su zona.

| Rol | Ve Mérida HOY/MAÑANA | Ve Granada HOY |
|---|---|---|
| limpieza_merida | ✅ | ❌ |
| limpieza_granada | ❌ | ✅ |
| lavanderia | ✅ | ✅ |
| mantenimiento | ✅ | ❌ |
| mantenimiento_granada | ❌ | ✅ |

---

## Layout Principal (page.tsx)

### Header (sticky, glass)
- Logo LuxSpace (dos versiones claro/oscuro)
- Toggle **LS / GR** (pills junto al logo, solo visible en tab Dashboard)
  - LS activo: fondo Luxeberry `#460479`, texto blanco
  - GR activo: fondo `red-500`, texto blanco
- Dark mode toggle (sol/luna)
- Avatar del usuario con color por role
- Botón "Salir"

### Bottom Dock (navegación principal — roles maestro/soporte)
Fijo en la parte inferior. Estructura:
- **Pill `.bottom-dock`** (glass): 3 iconos de navegación
- **Botón `.bottom-dock-plus`** (circular, 62×62px): abre menú "+"

| Posición | ID | Label | Icono |
|---|---|---|---|
| Izquierda | `retiros` | Retiros | Flecha descarga |
| Centro | `dashboard` | Dashboard | Barras (w-7 h-7) |
| Derecha | `gastos` | Gastos | Tarjeta crédito |

```css
.bottom-dock {
  background: rgba(255,255,255,0.65);
  backdrop-filter: blur(40px) saturate(200%);
  border: 1.5px solid rgba(255,255,255,0.6);
  border-radius: 99px;
  padding: 8px 20px;
  height: 62px;
}
.bottom-dock-plus {
  width: 62px; height: 62px;
  border-radius: 50%;
  /* mismo glass style */
}
```

### Botón Fecha Flotante
- Posición: `fixed`, `right: 0`
- Estilo: `bg-red-500`, sin esquinas redondeadas, texto blanco bold
- **Single tap:** Expande panel de filtro (flechas ← →, dropdowns mes/año, botón "Listo")
- **Doble tap:** Navega al mes actual sin abrir panel

### Pull-to-Refresh
- Threshold: 40px, máximo 200px, resistencia `dy * 0.5`
- Indicador: círculo glass 40×40px con flecha Luxeberry → spinner durante carga

### Tabs del menú "+" (acceso completo)
| Tab | ID | Componente |
|---|---|---|
| Dashboard | `dashboard` | KPIs + gráficas |
| Retiros | `retiros` | RetirosForm |
| Gastos | `gastos` | GastosForm |
| Blancos Diarios | `blancos_diarios` | BlancosDiariosModule |
| Cargar PDFs | `cargar_pdf` | PdfUploader |
| Ingresos | `ingresos` | IngresosForm |
| Gastos Fijos | `gastos_fijos` | GastosFijosForm |
| Devoluciones | `devoluciones` | DevolucionesForm |
| Proveedores | `proveedores` | ProveedoresForm |
| Cuentas | `cuentas` | CuentasForm (nuevo) |
| Propiedades | `propiedades` | PropiedadesModule (nuevo) |

---

## Tab: Dashboard

### Toggle LuxSpace / Granada
Pills "LS" / "GR" en el header — solo visible en tab dashboard.

### Vista LuxSpace — 6 KPIs
| KPI | Color dark | Color light | Variante |
|---|---|---|---|
| Ingresos Totales | emerald-400 | emerald-600 | positive |
| Devoluciones | blue-400 | gray-900 | refunds |
| Impuestos (6%) | blue-400 | gray-900 | refunds |
| Gastos Totales | red-400 | red-500 | negative |
| Retiros | amber-400 | gray-900 | returns |
| Utilidad (Balance) | dinámico | dinámico | balance |

**Subtitle Gastos Totales:** "F $X + V $Y"

**KPIs clickeables** → Modal detalle:
- **Ingresos:** lista por habitación (verde)
- **Devoluciones:** lista con montos (azul)
- **Gastos:** desglose por categoría (rojo)
- **Retiros:** desglose por usuario con conteo. Formato: "Leo — 3 retiros — $14,748.00"
- **Utilidad:** resumen balance completo

### Vista LuxSpace — Gráficas (Recharts)
- **Barras:** 6 barras resumen (Ingresos, Devoluciones, Impuestos, G.Fijos, G.Var, Retiros)
- **Pie:** Gastos por categoría

---

## Tab: Blancos Diarios (BlancosDiariosModule)

### Fuente de datos
`GET /api/blancos` → Lee hojas del spreadsheet `11tbthRn5U7...` generado por Apps Script.

### Estructura visual
3 tabs: **Mérida HOY** (verde), **Mérida MAÑANA** (ámbar), **Granada HOY** (rojo)

Por cada tab:
- Barra de progreso: `X/Y listas` con color Luxeberry
- Cards por zona (`Zona 1`, `Zona 2`, `Granada`)
- Cada habitación: checkbox tap-to-confirm + nombre + tipo cama + noches badge
- Timestamp "Actualizado: HH:MM"
- Botón "Actualizar" para refetch

**Estado de confirmación:** Solo en memoria (se resetea al salir). No persiste en DB.

### Permisos de tabs por rol
```
limpieza_merida      → Mérida HOY + MAÑANA
limpieza_granada     → Granada HOY
lavanderia           → Las 3 tabs
mantenimiento        → Mérida HOY + MAÑANA
mantenimiento_granada→ Granada HOY
soporte/atencion/leo/axl → Las 3 tabs
```

---

## Tab: Ingresos

### Filtro por Habitación
Botón "Filtrar" con badge contador. Panel checkboxes agrupados por propiedad.
Usa `toInternalRoomCode()` para matching.

### Tabla
| Columna | Render |
|---|---|
| Fecha | texto |
| Monto Total | `formatCurrency(depositado + pendiente)` verde bold |
| Habitación | `friendlyRoomName(propiedad_id)` |
| Fuente | Badge: Airbnb (rosa), Booking (azul), Externo (ámbar) |
| Depositado | formatCurrency |
| Pendiente | formatCurrency o "—" |

### Modal Nuevo Ingreso
Fecha · Propiedad (select 9) · Fuente (3 botones) · Total Depositado · Saldo Pendiente

---

## Tab: Gastos

### Tabla
Fecha · Monto (coral bold) · Usuario · Categoría (badge) · Descripción · Propiedad

### Modal Nuevo Gasto
Fecha · Monto · Usuario (Leo/Axl/Soporte) · Categoría (10 opciones) · Descripción · Propiedad (11 opciones)

---

## Tab: Gastos Fijos
No filtra por periodo. Muestra todos los activos.

---

## Tab: Devoluciones

### Modal
Habitación (select 39 habitaciones → "M1-1 — ROSADA") · Auto-fill propiedad · Fecha · Monto

---

## Tab: Retiros

### Tabla
Fecha · Monto (amber bold) · Usuario · Origen badge (LuxSpace=sky / Granada=purple)

### Modal
Fecha · Monto · Usuario · **Origen de Fondos** (2 botones: LuxSpace azul / Granada rojo)

---

## Tab: Proveedores
Nombre · Teléfono + botón WhatsApp verde + botón Llamar azul · Servicio (badge)
`formatPhone()` agrega +52 automáticamente si número tiene 10 dígitos.

---

## Tab: Cuentas (nuevo)
CRUD para cuentas bancarias del negocio.

---

## Tab: Propiedades (nuevo)
CRUD para gestión de propiedades con `PropiedadesContext.tsx`.

---

## Componentes Reutilizables

### KPICard
Props: `label, value, icon, variant, subtitle, delay, onClick`
Variants: `positive` (verde) · `negative` (rojo) · `refunds` (azul) · `returns` (amber) · `balance` (dinámico) · `default`
Clickeable → modal detalle. Animación `animate-fade-up` con delay escalonado.

### FormField
Types: `text`, `number`, `date`, `select`, `tel`
CSS: min-height 42px, width 100%.

### Modal
Overlay dark · click-to-close · `overflow-x-hidden` · `animate-scale-in`

### DataTable
Props: `columns, data, onEdit, onDelete, canEdit`
Botones Editar/Eliminar solo si `canEdit=true`. Confirmación antes de eliminar.

### Toast
Types: `success` (verde) · `error` (rojo). Auto-dismiss.

### BlancosDiariosModule
Lee `/api/blancos`. Muestra 3 tabs con habitaciones como checkboxes. Ver sección arriba.
