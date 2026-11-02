# Implementation Notes — Fanz Portfolio Redesign
> Generated from ui-ux-pro-max skill analysis. Read MASTER.md for full design system, then per-page overrides in pages/.

## Struktur File
```
design-system/
└── fanz-portfolio/
    ├── MASTER.md          ← Design system global (source of truth)
    └── pages/
        ├── index.md       ← Hero / landing page overrides
        ├── projects.md    ← Projects grid overrides
        ├── about.md       ← Skills orbital overrides
        ├── certificates.md← Certificates archive overrides
        └── contact.md     ← About/contact overrides
```

## Cara Pakai
Sebelum mengerjakan halaman apapun:
1. Baca `design-system/fanz-portfolio/MASTER.md`
2. Cek apakah ada `pages/<nama-halaman>.md`
3. Kalau ada, aturan di page file **override** MASTER
4. Kalau tidak ada, pakai MASTER saja

---

## Temuan Analisis (Prioritas)

### 🔴 CRITICAL

#### 1. Konten ketutupan floating nav dock
- **Masalah:** Nav dock `fixed bottom-6` menutupi konten di semua halaman
- **Fix:** Tambah `padding-bottom: 6rem` (atau `pb-24`) ke `<main>` atau `<body>` di setiap halaman
- **File:** semua `templates/*.html`
- **Skill rule:** "Fixed nav should not obscure content — add padding equal to nav height"

#### 2. z-index conflict di halaman Skills
- **Masalah:** Three.js canvas menutup teks HUD dan nav dock
- **Fix:** Three.js container `z-index: 1`, header `z-index: 100`, main `z-index: 50`, nav dock `z-index: 200`
- **File:** `templates/about.html`

---

### 🟠 HIGH

#### 3. Terlalu banyak font (3 font families)
- **Masalah:** Newsreader + Inter + JetBrains Mono = berat, tidak konsisten
- **Skill rekomendasikan:** Pilih salah satu stack:
  - **Option A (Developer):** `JetBrains Mono` (heading/label) + `IBM Plex Sans` (body)
  - **Option B (Ultra minimal):** `Inter` saja dengan weight variations
- **Hapus:** Newsreader dari semua template
- **File:** semua `templates/*.html` dan `static/theme.css`

#### 4. Accent crimson terlalu dominan
- **Masalah:** `#df3048` dipakai di mana-mana — glow, border, dot, badge
- **Fix:** Crimson HANYA untuk:
  - Active nav item dot
  - CTA button primary
  - Error/destructive state
- **Sisanya:** gunakan `rgba(255,255,255,0.1)` atau monochrome
- **File:** `static/theme.css` dan semua templates

#### 5. Terlalu banyak efek glass berlapis
- **Masalah:** Setiap elemen punya `backdrop-filter`, multiple `box-shadow`, glow
- **Fix:** Sederhanakan — hanya hero panel yang pakai full glass effect
  - Cards: solid `#111318` dengan border tipis `rgba(255,255,255,0.08)` saja
  - Hapus aurora blob animasi dari semua halaman kecuali hero
- **File:** `static/theme.css`, `static/aurora.css`

---

### 🟡 MEDIUM

#### 6. Spacing tidak konsisten
- **Masalah:** Padding/gap campur antara `rem`, `px`, Tailwind arbitrary values
- **Fix:** Standardisasi ke **8px rhythm** — semua spacing kelipatan 8
  - `space-xs: 4px`, `space-sm: 8px`, `space-md: 16px`, `space-lg: 24px`, `space-xl: 32px`, `space-2xl: 48px`
- **File:** `static/theme.css`

#### 7. Background terlalu gelap, kehilangan depth
- **Masalah:** `#08090C` hampir hitam total, cards tidak terlihat jelas
- **Skill rekomendasikan OLED dark:** `#000000` atau `#111318` base, card `#1C1C1E`
- **Fix:** Perjelas perbedaan layer dengan tonal step yang lebih jelas

#### 8. Halaman certificates belum ada route
- **Masalah:** `handleCertificates()` ada di handlers.go tapi tidak di-register di main.go
- **Fix:** Tambah route `/certificates` di `main.go`
- **File:** `main.go`

---

### 🟢 LOW

#### 9. Header brand link masih `href="#"`
- **Masalah:** Klik logo "Fanz" di header tidak kemana-mana
- **Fix:** Ganti `href="#"` → `href="/"`
- **File:** semua templates kecuali index.html

#### 10. Foto portrait dari URL eksternal
- **Masalah:** Hero image masih pakai URL Google (`lh3.googleusercontent.com`)
- **Fix:** Ganti ke `/static/portrait.png` yang sudah ada
- **File:** `templates/index.html`

---

## Design Tokens yang Harus Diupdate di theme.css

```css
:root {
  /* Background layers — OLED dark */
  --bg:              #000000;
  --bg-surface:      #111318;
  --bg-card:         #1C1C1E;
  --bg-card-hover:   #242428;

  /* Text */
  --text-primary:    #F5F5F7;
  --text-secondary:  rgba(245,245,247,0.6);
  --text-muted:      rgba(245,245,247,0.35);

  /* Accent — SATU warna saja */
  --accent:          #df3048;
  --accent-dim:      rgba(223,48,72,0.15);

  /* Border */
  --border:          rgba(255,255,255,0.08);
  --border-bright:   rgba(255,255,255,0.15);

  /* Spacing rhythm (8px base) */
  --space-1: 4px;
  --space-2: 8px;
  --space-3: 12px;
  --space-4: 16px;
  --space-6: 24px;
  --space-8: 32px;
  --space-12: 48px;
  --space-16: 64px;
}
```

## Font Stack yang Direkomendasikan

```css
/* Option A: Developer (rekomendasi utama) */
@import url('https://fonts.googleapis.com/css2?family=JetBrains+Mono:wght@400;500;600;700&family=IBM+Plex+Sans:wght@300;400;500;600;700&display=swap');

--font-heading: 'JetBrains Mono', monospace;
--font-body:    'IBM Plex Sans', sans-serif;
--font-mono:    'JetBrains Mono', monospace;

/* Option B: Ultra minimal */
@import url('https://fonts.googleapis.com/css2?family=Inter:wght@300;400;500;600;700&display=swap');

--font-heading: 'Inter', sans-serif;
--font-body:    'Inter', sans-serif;
```
