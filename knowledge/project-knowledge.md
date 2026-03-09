# Project Knowledge — SignForce

## Accumulated Insights

### Frontend
- Tailwind CDN loads before signforce-global.css, but our CSS uses !important to win
- `bg-[#0f172a]` in Tailwind body classes was causing navy blue instead of black — solved by injecting `<style>` before `</body>`
- Some pages have both inline `<style>` blocks AND external CSS — the inline override block is necessary for Tailwind conflicts
- `document_viewer.html` has complex PDF canvas rendering — never modify without explicit permission

### Backend
- sfmiddle is the main gateway — all frontend requests go through it
- sfback handles the actual cryptographic operations
- JWT tokens are issued by sfmiddle, validated by sfback
- File paths in MySQL use relative paths from BASE_DIR env variable

### Database
- `signforce_db.sql` contains seed data with real user information (Juan Núñez, María Barrera CVs)
- Face embeddings are stored as Base64 in the `faceembeddings` table
- Signatures table records IP, user agent, timezone for audit trail

### Infrastructure
- Skynet server hosts both SignForce and LuxSpace — they share Nginx but are completely isolated in Docker networks
- Let's Encrypt auto-renews but sometimes needs manual `certbot renew`
- The VPS has limited resources — don't run heavy builds in parallel

---
*Add new knowledge entries below this line:*

