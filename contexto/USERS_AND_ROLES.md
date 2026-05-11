# USERS_AND_ROLES.md — Sistema de Usuarios y Roles
**Última actualización:** 6 de Marzo de 2026

## Tabla Completa de Usuarios

| Login | Nombre | Role | canEdit | Password Env | Acceso |
|---|---|---|---|---|---|
| `leo` | Leo | `leo` | ✅ | `USER_LEO_PASS` | Full dashboard + admin |
| `axl` | Axl | `axl` | ✅ | `USER_AXL_PASS` | Full dashboard + admin |
| `soporte` | Soporte | `soporte` | ✅ | `USER_ADMIN_PASS` | Full dashboard |
| `marlene` | Marlene Aguilera | `atencion` | ✅ | `USER_MARLENE_PASS` | Dock reducido |
| `mauricio` | Mauricio Azcorra | `atencion` | ✅ | `USER_MAURICIO_PASS` | Dock reducido |
| `marcelo` | Marcelo | `limpieza_merida` | ❌ | `USER_MARCELO_PASS` | Solo Blancos Mérida |
| `luz` | Luz | `limpieza_merida` | ❌ | `USER_LUZ_PASS` | Solo Blancos Mérida |
| `cindy` | Cindy | `limpieza_merida` | ❌ | `USER_CINDY_PASS` | Solo Blancos Mérida |
| `yael` | Yael | `limpieza_granada` | ❌ | `USER_YAEL_PASS` | Solo Blancos Granada |
| `javier` | Javier Chavez | `limpieza_granada` | ❌ | `USER_JAVIER_PASS` | Solo Blancos Granada |
| `fabiola` | Fabiola Roque | `lavanderia` | ❌ | `USER_FABIOLA_PASS` | Blancos Mérida + Granada |
| `jorge` | Jorge Chin | `mantenimiento` | ❌ | `USER_JORGE_PASS` | Solo Blancos Mérida |
| `victor` | Victor | `mantenimiento_granada` | ❌ | `USER_VICTOR_PASS` | Solo Blancos Granada |

**Patrón contraseñas:** `PrimerNombre123` (ej: `Marcelo123`, `Fabi123`)

---

## Roles y Permisos

### `leo` / `axl`
- Dashboard financiero completo (LuxSpace + Granada)
- Bottom dock: Retiros, Dashboard, Gastos
- Menú "+": todos los módulos
- canEdit: true

### `soporte`
- Idéntico a leo/axl en permisos
- canEdit: true

### `atencion` (Marlene, Mauricio)
- Dock reducido: Devoluciones, Gastos, Proveedores, Blancos Diarios
- Sin botón "+"
- canEdit: true (pueden crear/editar/eliminar en sus módulos)

### `limpieza_merida` (Marcelo, Luz, Cindy)
- Solo ven: Blancos Diarios — tabs Mérida HOY + Mérida MAÑANA
- No ven Granada HOY
- canEdit: false

### `limpieza_granada` (Yael, Javier)
- Solo ven: Blancos Diarios — tab Granada HOY únicamente
- canEdit: false

### `lavanderia` (Fabiola)
- Blancos Diarios: las 3 tabs (Mérida HOY, MAÑANA, Granada HOY)
- canEdit: false

### `mantenimiento` (Jorge)
- Solo Blancos Diarios — tabs Mérida HOY + MAÑANA
- canEdit: false

### `mantenimiento_granada` (Victor)
- Solo Blancos Diarios — Granada HOY
- canEdit: false

---

## Lógica de Acceso en auth.ts

```typescript
export type UserRole =
  | "admin" | "leo" | "axl" | "soporte" | "atencion"
  | "limpieza_merida" | "limpieza_granada"
  | "lavanderia" | "mantenimiento" | "mantenimiento_granada";

// Cada usuario referencia su ENV key:
const USERS: Record<string, { envKey: string; role: UserRole; canEdit: boolean }> = {
  leo:      { envKey: "USER_LEO_PASS",      role: "leo",      canEdit: true  },
  marlene:  { envKey: "USER_MARLENE_PASS",  role: "atencion", canEdit: true  },
  marcelo:  { envKey: "USER_MARCELO_PASS",  role: "limpieza_merida", canEdit: false },
  // ...etc
};
```

**Función `authenticateUser(username, password)`:**
1. Normaliza username a lowercase
2. Busca en USERS map
3. Compara password con `process.env[user.envKey]`
4. Si match → genera JWT con `{ username, role, canEdit }`, expira 7 días
5. Si no → retorna null → API devuelve 401

---

## Blancos Diarios — Acceso por Rol

```typescript
// ¿Ve tab Mérida HOY/MAÑANA?
canSeeMerida(role) → ["soporte","atencion","limpieza_merida","lavanderia","leo","axl","mantenimiento"]

// ¿Ve tab Granada HOY?
canSeeGranada(role) → ["soporte","atencion","limpieza_granada","lavanderia","leo","axl","mantenimiento_granada"]
```

---

## Notas Importantes

1. **`soporte` y `atencion`** comparten `canEdit: true` pero tienen vistas distintas en el dock.
2. **`admin`** existe como type en TypeScript pero ningún usuario tiene ese role actualmente.
3. El login acepta el nombre en cualquier capitalización (`Marlene`, `MARLENE`, `marlene` → todos funcionan).
4. Los roles de limpieza deberían hacer redirect a `/blancos` post-login, pero esto **aún está pendiente** de implementar en `login/page.tsx`.
