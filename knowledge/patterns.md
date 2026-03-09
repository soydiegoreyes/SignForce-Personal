# Code Patterns — SignForce

## Frontend Patterns

### Standard Page Template
```html
<!DOCTYPE html>
<html lang="es">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>SignForce — Page Title</title>
  <link href="https://fonts.googleapis.com/css2?family=DM+Sans:opsz,wght@9..40,300;9..40,400;9..40,500;9..40,600;9..40,700&display=swap" rel="stylesheet">
  <link href="https://fonts.googleapis.com/css2?family=Material+Symbols+Outlined" rel="stylesheet"/>
  <link rel="stylesheet" href="../signforce-global.css">
  <link rel="icon" type="image/svg+xml" href="/favicon.svg">
  <meta name="theme-color" content="#060a14">
  <style>
    body{background:#060a14;font-family:'DM Sans',sans-serif;color:#fff;margin:0;min-height:100vh}
    /* page-specific styles here */
  </style>
</head>
<body>
  <div style="position:fixed;top:0;left:0;right:0;height:3px;background:linear-gradient(90deg,#B5C413,#d4e157,#B5C413);z-index:9999;"></div>
  <div class="sf-orb sf-orb-1"></div><div class="sf-orb sf-orb-2"></div><div class="sf-orb sf-orb-3"></div>

  <!-- HEADER -->
  <header class="sf-header">...</header>

  <!-- LAYOUT -->
  <div class="sf-layout">
    <aside class="sf-sidebar">...</aside>
    <main class="sf-content">...</main>
  </div>

  <style>html,body{background:#060a14 !important;background-color:#060a14 !important;}</style>
</body>
</html>
```

### Sidebar Pattern
```html
<aside class="sf-sidebar">
  <!-- Logo card -->
  <div class="sf-sidebar-card" style="padding:1.1rem 1.25rem;">
    <div style="display:flex;align-items:center;gap:.6rem;">
      <div class="sf-card-icon" style="background:rgba(181,196,19,.1);">
        <span class="material-symbols-outlined" style="font-size:1.1rem;color:#B5C413;">icon_name</span>
      </div>
      <div>
        <div style="font-size:.95rem;font-weight:700;">Title</div>
        <div style="font-size:.68rem;color:rgba(255,255,255,.35);">Subtitle</div>
      </div>
    </div>
  </div>
  <!-- Nav card -->
  <nav class="sf-sidebar-card" style="flex:1;">
    <div class="sf-sidebar-title">Navegación</div>
    <a class="sf-nav-link active" href="#"><span class="material-symbols-outlined">icon</span> Active</a>
    <a class="sf-nav-link" href="#"><span class="material-symbols-outlined">icon</span> Link</a>
  </nav>
</aside>
```

### Glass Card with Gradient Line
```html
<div class="sf-card">
  <div class="sf-card-header">
    <div class="sf-card-icon" style="background:rgba(181,196,19,.1);">
      <span class="material-symbols-outlined" style="font-size:1.1rem;color:#B5C413;">icon</span>
    </div>
    <div class="sf-card-title">Card Title</div>
  </div>
  <!-- content -->
</div>
```

### Hero Gradient Card
```html
<div style="background:rgba(15,21,36,.45);backdrop-filter:blur(16px);border:1px solid rgba(255,255,255,.06);border-radius:20px;padding:1.75rem;position:relative;overflow:hidden;">
  <div style="position:absolute;inset:0;background:linear-gradient(135deg,rgba(181,196,19,.06),rgba(52,211,153,.04),transparent);pointer-events:none;"></div>
  <div style="position:absolute;top:0;left:0;right:0;height:1px;background:linear-gradient(90deg,transparent,rgba(181,196,19,.25),transparent);pointer-events:none;"></div>
  <div style="position:relative;z-index:1;">
    <!-- hero content -->
  </div>
</div>
```

### Button Styles
```html
<!-- Primary (green) -->
<button class="sf-btn sf-btn-primary">Action</button>

<!-- Ghost (transparent) -->
<button class="sf-btn-ghost">Secondary</button>

<!-- Small action buttons -->
<button class="sf-btn-sm sf-btn-edit">Edit</button>
<button class="sf-btn-sm sf-btn-del">Delete</button>
```

## Backend Patterns

### Go Handler Pattern (sfmiddle)
```go
func HandlerName(w http.ResponseWriter, r *http.Request) {
    if r.Method != "POST" {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }
    // Parse request
    // Validate input
    // Call sfback or DB
    // Return JSON response
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response)
}
```

### Database Query Pattern (CURRENT — has SQL injection risk)
```go
// ⚠️ UNSAFE — uses fmt.Sprintf
query := fmt.Sprintf("SELECT %s FROM %s WHERE %s;", attrs, table, wheres)
```

### Database Query Pattern (RECOMMENDED — prepared statements)
```go
// ✅ SAFE — uses placeholders
query := "SELECT name, email FROM users WHERE id = ? AND status = ?"
rows, err := db.Query(query, userID, "active")
```

---
*Add new patterns below this line:*

