# IT Portfolio & CMS Monorepo (Go + React)

A production-grade, local-first personal portfolio and content management system (CMS) built with a Go Fiber v2 backend and React 18 + Vite frontend, backed by PostgreSQL 17. Architecture is decoupled via repository & storage provider patterns for future migration to Appwrite Cloud.

---

## 🏗️ Architecture & Monorepo Structure

```text
d:\porto project\
├── apps/
│   ├── api/                     # Backend (Go 1.23+ / Fiber v2)
│   │   ├── cmd/
│   │   │   ├── server/          # Fiber REST API server (port :8080)
│   │   │   ├── migrate/         # Database migration runner (up/down/version)
│   │   │   └── seed/            # Idempotent admin & profile seeder
│   │   ├── internal/
│   │   │   ├── config/          # Environment configuration & validation
│   │   │   ├── database/        # pgxpool PostgreSQL connection pool
│   │   │   ├── domain/          # Domain models, requests & response envelopes
│   │   │   ├── handler/         # HTTP handlers for public, auth, admin & uploads
│   │   │   ├── middleware/      # Auth, CSRF, CORS, RateLimit, Logger, SecurityHeaders
│   │   │   ├── ratelimit/       # In-memory sliding-window rate limiter
│   │   │   ├── repository/      # PostgreSQL data access layer (parameterized queries)
│   │   │   ├── sanitize/        # Goldmark + Bluemonday UGCPolicy + URL validation
│   │   │   ├── security/        # Argon2id password hashing & timing-safe compare
│   │   │   ├── service/         # Business logic & domain validation
│   │   │   ├── storage/         # StorageProvider & LocalStorageProvider (uploads)
│   │   │   └── validator/       # Struct validation with Indonesian error messages
│   │   └── migrations/          # Pure SQL migrations (10 core tables)
│   └── web/                     # Frontend (React 18 / Vite / Tailwind CSS v3)
│       ├── src/
│       │   ├── components/      # UI components (Navbar, Footer, Cards, Modals)
│       │   │   └── admin/       # Admin layout & navigation
│       │   ├── pages/           # Public pages (Home, Projects, Posts, Detail, 404)
│       │   │   └── admin/       # CMS pages (Dashboard, Profile, Projects, Skills, etc.)
│       │   ├── lib/             # API client with CSRF token injection
│       │   ├── styles/          # Dark mode design tokens & font imports
│       │   └── types/           # TypeScript contracts for all models & payloads
├── scripts/                     # PowerShell automation scripts (Windows native)
│   ├── dev-api.ps1              # Start Go backend server
│   ├── dev-web.ps1              # Start React Vite dev server
│   ├── migrate.ps1              # Run DB migrations (up / down / version)
│   ├── seed.ps1                 # Run database seeder
│   ├── build.ps1                # Build both backend binary & frontend dist
│   └── test.ps1                 # Run Go unit tests & TypeScript type check
├── storage/
│   └── uploads/                 # Local uploaded files (avatars, resumes, projects, etc.)
├── docs/                        # Project specification & PRD
├── Makefile                     # Make command targets
└── package.json                 # Monorepo orchestrator scripts
```

---

## 🔒 Security Architecture

1. **Password Hashing**: Argon2id (`memory=64MB`, `time=3`, `threads=2`, `salt=16B`, `key=32B`).
2. **Authentication**: JWT stored exclusively in secure `httpOnly` cookie (`auth_token`).
3. **Double-Submit CSRF**: Non-httpOnly `csrf_token` cookie sent via `X-CSRF-Token` header for all mutating HTTP methods (`POST`, `PUT`, `PATCH`, `DELETE`), verified using `subtle.ConstantTimeCompare`.
4. **Brute Force & Rate Limiting**:
   - Login rate limit: 5 requests per 15 minutes per IP + email.
   - Account lockout: 15-minute lock after 5 consecutive failed attempts.
   - Contact form rate limit: 5 submissions per hour per IP.
   - Global rate limit: 100 requests per minute per IP.
5. **Markdown Sanitization (Defense-in-Depth)**:
   - Backend: Goldmark (raw HTML disabled) → Bluemonday `UGCPolicy` → stored in `*_html` database columns.
   - Frontend: `DOMPurify.sanitize()` executed before rendering in React.
6. **Upload Security**:
   - Whitelist file extensions: `.jpg`, `.jpeg`, `.png`, `.webp`, `.pdf`.
   - Magic bytes validation using `http.DetectContentType`.
   - Complete rejection of SVG files (both extension and sniffing).
   - Filenames generated using 128-bit cryptographically secure UUID hex strings.
   - Anti path-traversal validation via `filepath.Rel`.
   - Strict 2 MB file size limit.
   - Static files served with `X-Content-Type-Options: nosniff`.
7. **SQL Injection**: Parameterized SQL queries throughout all repository operations.

---

## 🚀 Quick Start Guide

### 1. Prerequisites
- **PostgreSQL 16/17** running locally with database `portfolio`.
- **Go 1.23+** (installed at `D:\tools\go`).
- **Node.js 18+** & npm.

### 2. Environment Setup
Copy `.env.example` to `.env` in the root monorepo directory:
```bash
DATABASE_URL=postgres://portfolio:portfolio@localhost:5432/portfolio?sslmode=disable
JWT_SECRET=rahasia-super-aman-minimal-32-karakter-kunci-jwt-portofolio-2026
APP_PORT=8080
STORAGE_LOCAL_DIR=./storage/uploads
```

### 3. Run Migrations & Seed Admin User
Using npm:
```bash
npm run migrate:up
npm run seed
```
Or using PowerShell scripts:
```powershell
.\scripts\migrate.ps1 up
.\scripts\seed.ps1
```

Default seeded administrator:
- **Email**: `admin@local.test`
- **Password**: `Admin123!Porto2026`

---

## 💻 Running the Application

### Option A: npm scripts (Monorepo root)
```bash
npm run dev:api     # Starts Go API server on http://localhost:8080
npm run dev:web     # Starts Vite dev server on http://localhost:5173
npm run test        # Runs backend unit tests & frontend type checks
npm run build       # Builds backend binary & frontend production bundle
```

### Option B: PowerShell scripts (Windows)
```powershell
.\scripts\dev-api.ps1
.\scripts\dev-web.ps1
.\scripts\test.ps1
.\scripts\build.ps1
```

### Option C: Makefile
```bash
make dev-api
make dev-web
make test
make build
```

---

## 🌐 Endpoints & Pages

| Route | Access | Description |
|---|---|---|
| `http://localhost:5173/` | Public | Landing page (Hero, Bio, Skills, Featured Projects, Timeline, Certificates, Blog, Contact) |
| `http://localhost:5173/projects` | Public | All projects archive with category & skill filters |
| `http://localhost:5173/projects/:slug` | Public | Project detail page with sanitized markdown description |
| `http://localhost:5173/posts` | Public | Blog & technical articles archive |
| `http://localhost:5173/posts/:slug` | Public | Article reading page with sanitized content |
| `http://localhost:5173/admin/login` | Public | Admin login portal |
| `http://localhost:5173/admin` | Admin | CMS Dashboard & overview metrics |
| `http://localhost:5173/admin/profile` | Admin | Edit bio, avatar & CV upload, social links |
| `http://localhost:5173/admin/projects` | Admin | CRUD projects & thumbnails |
| `http://localhost:5173/admin/skills` | Admin | CRUD technical skills & color codes |
| `http://localhost:5173/admin/experiences` | Admin | CRUD career & education timeline |
| `http://localhost:5173/admin/certificates` | Admin | CRUD certificates & document uploads |
| `http://localhost:5173/admin/posts` | Admin | CRUD blog articles & cover images |
| `http://localhost:5173/admin/categories` | Admin | CRUD taxonomies |
| `http://localhost:5173/admin/messages` | Admin | Inbox, status updates & contact management |
| `http://localhost:8080/api/v1/health` | Public | Healthcheck endpoint (`{"database": "connected"}`) |
| `http://localhost:8080/uploads/*` | Public | Static uploads with nosniff header |


Panduan lengkap sudah siap. Singkatnya, ada 4 langkah utama:

Alur Singkat
#	File yang diubah	Fungsinya
1️⃣	pages/NamaPage.tsx (file baru)	Tulis konten halamannya
2️⃣	App.tsx	Daftarkan URL route-nya
3️⃣	components/Navbar.tsx	Tambah link di menu navigasi
4️⃣ (opsional)	pages/admin/ + AdminLayout.tsx	Kalau butuh CMS admin juga
Poin Kunci yang Perlu Diingat
App.tsx = pusat kontrol semua routing. Semua halaman harus didaftarkan di sini.
queryKey di useQuery harus unik tiap halaman, bisa berupa string atau array seperti ['projects', categoryId].
api.ts = semua fungsi pemanggil backend sudah tersedia di sini, tinggal pilih yang sesuai.
Route * (Not Found) wajib paling bawah, taruh route baru di atasnya.
Untuk URL dinamis seperti /projects/:slug → ambil parameternya pakai useParams().

Coba mulai dari halaman publik sederhana dulu tanpa data API untuk latihan, lalu tambah useQuery setelah paham strukturnya!