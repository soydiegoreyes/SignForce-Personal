# Development Workflow — SignForce

## Git Workflow

### Branch Strategy
- `main` — production branch, deployed to VPS
- Feature branches recommended but currently all work goes to main

### Commit Convention
```
type: brief description

Types:
  feat     — new feature
  fix      — bug fix
  rewrite  — full page/component rewrite
  security — security improvement
  cleanup  — code cleanup, dead code removal
  polish   — visual/UX improvement
  docs     — documentation only
```

### Deploy After Push
```bash
# On VPS (Skynet)
cd ~/signforce && git pull && docker compose -f docker-compose.frontend.yml up -d --build
```

## Frontend Development

### Adding a New Page
1. Create HTML file in appropriate directory (registro/, dashboards/, documentFlow/, administracion/)
2. Include these in `<head>`:
   ```html
   <link href="https://fonts.googleapis.com/css2?family=DM+Sans:opsz,wght@9..40,300;9..40,400;9..40,500;9..40,600;9..40,700&display=swap" rel="stylesheet">
   <link href="https://fonts.googleapis.com/css2?family=Material+Symbols+Outlined" rel="stylesheet"/>
   <link rel="stylesheet" href="../signforce-global.css">
   <link rel="icon" type="image/svg+xml" href="/favicon.svg">
   ```
3. Include before `</body>`:
   ```html
   <style>html,body{background:#060a14 !important;background-color:#060a14 !important;}</style>
   ```
4. Use the standard sidebar pattern (see CLAUDE.md)
5. Add orbs: `<div class="sf-orb sf-orb-1"></div><div class="sf-orb sf-orb-2"></div><div class="sf-orb sf-orb-3"></div>`

### Modifying Styles
- **Global changes:** Edit `sffront/signforce-global.css`
- **Page-specific:** Use inline `<style>` in the HTML file
- **Never** create new standalone CSS files — use signforce-global.css

### Color Palette
```
Primary:    #B5C413  (brand green)
Background: #060a14  (deep black)  
Card:       rgba(15,21,36,0.45)
Border:     rgba(255,255,255,0.06)
Text:       #ffffff / rgba(255,255,255,0.55)
Success:    #34d399
Info:       #60a5fa
Warning:    #fbbf24
Error:      #f87171
Purple:     #a78bfa
```

## Backend Development

### Running Locally
```bash
# Start MySQL
docker compose up mysql -d

# Run sfback
cd sfback && go run main.go

# Run sfmiddle (separate terminal)
cd sfmiddle && go run main.go
```

### Adding an Endpoint (sfmiddle)
1. Create handler function in appropriate file under `sfmiddle/handlers/`
2. Register route in `sfmiddle/main.go`
3. If it calls sfback, add corresponding endpoint there too

### Database Changes
1. Modify `signforce_db.sql`
2. Apply manually: `mysql -u root -p signforce_db < signforce_db.sql`
3. Update models in `sfback/models/` and `sfmiddle/models/`

## Testing
- No automated tests exist currently
- Manual testing via browser at https://testsignforce.luxspace.org
- Always test with Cmd+Shift+R (hard refresh) after CSS/HTML changes
- Check browser console for JS errors

## Infrastructure

### Server Access
```bash
ssh root@187.77.6.110
```

### Useful Commands
```bash
docker ps                          # Running containers
docker logs signforce-web --tail 50 # Frontend logs
pm2 status                         # Node.js services
nginx -t                           # Test Nginx config
certbot renew                      # Renew SSL certs
```
