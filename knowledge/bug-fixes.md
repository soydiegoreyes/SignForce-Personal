# Bug Fixes — SignForce

## Fix Log

### 2026-03-08: Navy blue background instead of black
- **Symptom:** Pages showed #0f172a (navy) instead of #060a14 (black)
- **Cause:** Tailwind CDN's `bg-[#0f172a]` on `<body>` class overrode our CSS
- **Fix:** Injected `<style>html,body{background:#060a14 !important;}</style>` before `</body>` on all 35 pages
- **Prevention:** Always add this style block to new pages

### 2026-03-08: Docker port conflict (8080 already allocated)
- **Symptom:** `Bind for 127.0.0.1:8080 failed: port is already allocated`
- **Cause:** docker-compose.frontend.yml was changed from 8090 to 8080
- **Fix:** Changed back to `8090:80`
- **Prevention:** CLAUDE.md documents that port 8090 is mandatory

### 2026-03-08: Sidebar not separated from content
- **Symptom:** Sidebar and main content panels were touching (no gap)
- **Cause:** Flex container missing `gap` property
- **Fix:** Added `style="gap:1.25rem;"` to layout container
- **Files:** add_signers.html, add_signs.html, view_invite_.html

### 2026-03-08: Old lime color (#c8ff00) inconsistent
- **Symptom:** Some pages had neon lime, others had the new brand green
- **Cause:** 35 pages had hardcoded `#c8ff00` in inline styles
- **Fix:** Global find/replace `#c8ff00` → `#B5C413` across all files

---
*Add new bug fixes below this line:*

