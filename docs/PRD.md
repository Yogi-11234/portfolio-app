# PRD — Web Portofolio + CMS Pribadi
Versi 1.0 · 29 Oktober 2026 · Dokumen acuan utama untuk AI Coding Agent (Antigravity)

## 1. Executive Summary & Objectives
### 1.1 Ringkasan
Web Portofolio + CMS Pribadi adalah satu aplikasi dengan dua sisi: situs publik (hero, showcase proyek, tech stack, pengalaman, sertifikat, blog, form kontak) dan dashboard admin /admin (CRUD semua konten dan inbox pesan). Seluruh konten dikelola lewat CMS sehingga pemilik tidak perlu mengubah kode untuk update portofolio.

### 1.2 Tujuan
1. Local-first: seluruh sistem berjalan penuh di localhost tanpa dependensi cloud.
2. Modular: setiap domain (projects, skills, certificates, dst.) berdiri sendiri: handler -> service -> repository.
3. Efisien resource: backend Go single binary, memori rendah, query terindeks, gambar dioptimalkan.
4. Security-first: auth aman, validasi ketat, sanitasi XSS, rate limit, upload aman.
5. Siap Appwrite: data layer dan storage dibungkus interface sehingga migrasi ke Appwrite Cloud hanya mengganti implementasi, bukan logika bisnis.

### 1.3 Non-Goals (v1)
Multi-user/multi-tenant, komentar publik, pembayaran, i18n penuh, light mode.

### 1.4 Keputusan Teknis Final
- **Backend**: Go 1.23+, Fiber v2, pgx/v5, golang-migrate, goldmark + bluemonday, go-playground/validator, argon2id
- **Frontend**: React 18 + Vite + TypeScript, React Router, TanStack Query, Tailwind CSS v3, shadcn/ui, react-hook-form + zod, framer-motion
- **Database**: PostgreSQL 16/17 lokal (DDL utama = PostgreSQL)
- **Storage**: Filesystem lokal `./storage/uploads` lewat interface `StorageProvider`
- **Auth**: JWT di httpOnly Secure SameSite=Strict cookie + CSRF double-submit untuk request admin mutasi
