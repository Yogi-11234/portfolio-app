-- 000001_init_schema.up.sql
CREATE EXTENSION IF NOT EXISTS pgcrypto;
CREATE EXTENSION IF NOT EXISTS citext;

CREATE OR REPLACE FUNCTION set_updated_at() RETURNS trigger AS $$
BEGIN
    NEW.updated_at = now();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- 1. users (admin CMS)
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email CITEXT NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    name VARCHAR(100) NOT NULL,
    role VARCHAR(20) NOT NULL DEFAULT 'admin' CHECK (role IN ('admin')),
    failed_logins SMALLINT NOT NULL DEFAULT 0,
    locked_until TIMESTAMPTZ,
    last_login_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- 2. profile (satu baris per user)
CREATE TABLE profile (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL UNIQUE REFERENCES users(id) ON DELETE CASCADE,
    full_name VARCHAR(120) NOT NULL,
    headline VARCHAR(160) NOT NULL,
    bio_md TEXT NOT NULL DEFAULT '',
    bio_html TEXT NOT NULL DEFAULT '',
    location VARCHAR(120),
    public_email VARCHAR(160),
    avatar_path TEXT,
    resume_path TEXT,
    github_url TEXT,
    linkedin_url TEXT,
    social_links JSONB NOT NULL DEFAULT '[]'::jsonb,
    open_to_work BOOLEAN NOT NULL DEFAULT false,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- 3. categories (dipakai projects & posts)
CREATE TABLE categories (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(80) NOT NULL,
    slug VARCHAR(100) NOT NULL,
    type VARCHAR(10) NOT NULL CHECK (type IN ('project','post')),
    sort_order INT NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (type, slug)
);

-- 4. projects
CREATE TABLE projects (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    category_id UUID REFERENCES categories(id) ON DELETE SET NULL,
    title VARCHAR(160) NOT NULL,
    slug VARCHAR(180) NOT NULL UNIQUE,
    summary VARCHAR(300) NOT NULL,
    description_md TEXT NOT NULL DEFAULT '',
    description_html TEXT NOT NULL DEFAULT '',
    thumbnail_path TEXT,
    repo_url TEXT,
    demo_url TEXT,
    status VARCHAR(20) NOT NULL DEFAULT 'draft' CHECK (status IN ('draft','published','archived')),
    is_featured BOOLEAN NOT NULL DEFAULT false,
    sort_order INT NOT NULL DEFAULT 0,
    started_at DATE,
    ended_at DATE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_projects_status_sort ON projects (status, is_featured DESC, sort_order, created_at DESC);
CREATE INDEX idx_projects_category ON projects (category_id);

-- 5. skills
CREATE TABLE skills (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(60) NOT NULL UNIQUE,
    slug VARCHAR(80) NOT NULL UNIQUE,
    group_name VARCHAR(20) NOT NULL DEFAULT 'tool' CHECK (group_name IN ('language','framework','database','tool','cloud','other')),
    icon_key VARCHAR(60),
    brand_color CHAR(7) CHECK (brand_color ~ '^#[0-9A-Fa-f]{6}$'),
    level SMALLINT CHECK (level BETWEEN 1 AND 5),
    is_visible BOOLEAN NOT NULL DEFAULT true,
    sort_order INT NOT NULL DEFAULT 0
);

-- 6. project_skills (many-to-many)
CREATE TABLE project_skills (
    project_id UUID NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    skill_id UUID NOT NULL REFERENCES skills(id) ON DELETE CASCADE,
    PRIMARY KEY (project_id, skill_id)
);

CREATE INDEX idx_project_skills_skill ON project_skills (skill_id);

-- 7. experiences
CREATE TABLE experiences (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    type VARCHAR(20) NOT NULL CHECK (type IN ('work','education','organization','volunteer')),
    organization VARCHAR(160) NOT NULL,
    role VARCHAR(160) NOT NULL,
    location VARCHAR(120),
    description_md TEXT NOT NULL DEFAULT '',
    description_html TEXT NOT NULL DEFAULT '',
    start_date DATE NOT NULL,
    end_date DATE,
    is_current BOOLEAN NOT NULL DEFAULT false,
    is_visible BOOLEAN NOT NULL DEFAULT true,
    sort_order INT NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    CHECK (end_date IS NULL OR end_date >= start_date)
);

CREATE INDEX idx_experiences_order ON experiences (is_visible, start_date DESC);

-- 8. certificates
CREATE TABLE certificates (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title VARCHAR(200) NOT NULL,
    issuer VARCHAR(160) NOT NULL,
    credential_id VARCHAR(120),
    credential_url TEXT,
    issued_at DATE NOT NULL,
    expires_at DATE,
    file_path TEXT NOT NULL,
    file_mime VARCHAR(30) NOT NULL CHECK (file_mime IN ('application/pdf','image/png','image/jpeg','image/webp')),
    thumbnail_path TEXT,
    is_visible BOOLEAN NOT NULL DEFAULT true,
    sort_order INT NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_certificates_visible ON certificates (is_visible, issued_at DESC);

-- 9. posts (blog)
CREATE TABLE posts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    author_id UUID NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    category_id UUID REFERENCES categories(id) ON DELETE SET NULL,
    title VARCHAR(200) NOT NULL,
    slug VARCHAR(220) NOT NULL UNIQUE,
    excerpt VARCHAR(320) NOT NULL DEFAULT '',
    content_md TEXT NOT NULL,
    content_html TEXT NOT NULL,
    cover_path TEXT,
    status VARCHAR(20) NOT NULL DEFAULT 'draft' CHECK (status IN ('draft','published','archived')),
    reading_minutes SMALLINT NOT NULL DEFAULT 1,
    published_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_posts_published ON posts (status, published_at DESC);

-- 10. messages (inbox form kontak)
CREATE TABLE messages (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) NOT NULL,
    email VARCHAR(160) NOT NULL,
    subject VARCHAR(160) NOT NULL DEFAULT '',
    body TEXT NOT NULL CHECK (char_length(body) BETWEEN 10 AND 5000),
    status VARCHAR(12) NOT NULL DEFAULT 'unread' CHECK (status IN ('unread','read','archived','spam')),
    ip_hash CHAR(64),
    user_agent VARCHAR(255),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    read_at TIMESTAMPTZ
);

CREATE INDEX idx_messages_status_created ON messages (status, created_at DESC);

-- trigger updated_at
CREATE TRIGGER trg_users_updated BEFORE UPDATE ON users FOR EACH ROW EXECUTE FUNCTION set_updated_at();
CREATE TRIGGER trg_profile_updated BEFORE UPDATE ON profile FOR EACH ROW EXECUTE FUNCTION set_updated_at();
CREATE TRIGGER trg_projects_updated BEFORE UPDATE ON projects FOR EACH ROW EXECUTE FUNCTION set_updated_at();
CREATE TRIGGER trg_experiences_updated BEFORE UPDATE ON experiences FOR EACH ROW EXECUTE FUNCTION set_updated_at();
CREATE TRIGGER trg_certificates_updated BEFORE UPDATE ON certificates FOR EACH ROW EXECUTE FUNCTION set_updated_at();
CREATE TRIGGER trg_posts_updated BEFORE UPDATE ON posts FOR EACH ROW EXECUTE FUNCTION set_updated_at();
