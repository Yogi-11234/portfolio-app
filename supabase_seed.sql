-- ==============================================================================
-- SUPABASE SEED SCRIPT FOR PORTFOLIO
-- Jalankan skrip ini langsung di Supabase SQL Editor (Dashboard -> SQL Editor -> New Query)
-- ==============================================================================

-- 1. Pastikan ekstensi tersedia
CREATE EXTENSION IF NOT EXISTS pgcrypto;
CREATE EXTENSION IF NOT EXISTS citext;

-- 2. USERS (Admin / Author)
INSERT INTO users (id, email, password_hash, name, role)
VALUES (
    '00000000-0000-0000-0000-000000000001',
    'yogi.aditia2004@gmail.com',
    '$argon2id$v=19$m=65536,t=3,p=4$vPq4eC8yT0Y9wZ12345678$mockhashedpasswordforyogiaditia2026',
    'Yogi Aditia',
    'admin'
)
ON CONFLICT (email) DO UPDATE
SET name = EXCLUDED.name;

-- 3. PROFILE
INSERT INTO profile (
    id,
    user_id,
    full_name,
    headline,
    bio_md,
    bio_html,
    location,
    public_email,
    github_url,
    linkedin_url,
    social_links,
    open_to_work
)
VALUES (
    '00000000-0000-0000-0000-000000000001',
    '00000000-0000-0000-0000-000000000001',
    'Yogi Aditia',
    'AI & Backend Software Engineer | Machine Learning Enthusiast',
    'Hi! I am **Yogi Aditia**, an AI and Backend Software Engineer focused on designing scalable distributed systems, high-performance APIs, and computer vision solutions. Experienced with Go, Python, TypeScript, and modern cloud database architectures.',
    '<p>Hi! I am <strong>Yogi Aditia</strong>, an AI and Backend Software Engineer focused on designing scalable distributed systems, high-performance APIs, and computer vision solutions.</p><p>Experienced with Go, Python, TypeScript, Machine Learning models, and modern cloud database architectures.</p>',
    'Jakarta, Indonesia',
    'yogi.aditia2004@gmail.com',
    'https://github.com/Yogi-11234',
    'https://linkedin.com/in/yogi-aditia',
    '{"github": "https://github.com/Yogi-11234", "linkedin": "https://linkedin.com/in/yogi-aditia", "email": "yogi.aditia2004@gmail.com"}'::jsonb,
    true
)
ON CONFLICT (user_id) DO UPDATE
SET full_name = EXCLUDED.full_name,
    headline = EXCLUDED.headline,
    bio_md = EXCLUDED.bio_md,
    bio_html = EXCLUDED.bio_html,
    location = EXCLUDED.location,
    public_email = EXCLUDED.public_email,
    github_url = EXCLUDED.github_url,
    linkedin_url = EXCLUDED.linkedin_url,
    social_links = EXCLUDED.social_links,
    open_to_work = EXCLUDED.open_to_work;

-- 4. CATEGORIES
INSERT INTO categories (id, name, slug, type, sort_order)
VALUES
    ('10000000-0000-0000-0000-000000000001', 'AI & Machine Learning', 'ai-ml', 'project', 1),
    ('10000000-0000-0000-0000-000000000002', 'Web & Backend', 'web-backend', 'project', 2),
    ('10000000-0000-0000-0000-000000000003', 'Computer Vision', 'computer-vision', 'project', 3),
    ('10000000-0000-0000-0000-000000000004', 'Backend Architecture', 'backend-architecture', 'post', 1),
    ('10000000-0000-0000-0000-000000000005', 'Machine Learning', 'machine-learning', 'post', 2)
ON CONFLICT (type, slug) DO UPDATE
SET name = EXCLUDED.name,
    sort_order = EXCLUDED.sort_order;

-- 5. SKILLS
INSERT INTO skills (id, name, slug, group_name, icon_key, brand_color, level, is_visible, sort_order)
VALUES
    ('20000000-0000-0000-0000-000000000001', 'Go (Golang)', 'golang', 'language', 'go', '#00ADD8', 5, true, 1),
    ('20000000-0000-0000-0000-000000000002', 'Python', 'python', 'language', 'python', '#3776AB', 5, true, 2),
    ('20000000-0000-0000-0000-000000000003', 'TypeScript', 'typescript', 'language', 'typescript', '#3178C6', 4, true, 3),
    ('20000000-0000-0000-0000-000000000004', 'JavaScript', 'javascript', 'language', 'javascript', '#F7DF1E', 5, true, 4),
    ('20000000-0000-0000-0000-000000000005', 'FastAPI', 'fastapi', 'framework', 'fastapi', '#009688', 5, true, 5),
    ('20000000-0000-0000-0000-000000000006', 'Fiber', 'fiber', 'framework', 'go', '#00ADD8', 5, true, 6),
    ('20000000-0000-0000-0000-000000000007', 'PyTorch', 'pytorch', 'framework', 'pytorch', '#EE4C2C', 4, true, 7),
    ('20000000-0000-0000-0000-000000000008', 'React', 'react', 'framework', 'react', '#61DAFB', 4, true, 8),
    ('20000000-0000-0000-0000-000000000009', 'PostgreSQL', 'postgresql', 'database', 'postgresql', '#4169E1', 5, true, 9),
    ('20000000-0000-0000-0000-000000000010', 'Redis', 'redis', 'database', 'redis', '#DC382D', 4, true, 10),
    ('20000000-0000-0000-0000-000000000011', 'Supabase', 'supabase', 'database', 'supabase', '#3ECF8E', 5, true, 11),
    ('20000000-0000-0000-0000-000000000012', 'Docker', 'docker', 'cloud', 'docker', '#2496ED', 4, true, 12),
    ('20000000-0000-0000-0000-000000000013', 'Linux', 'linux', 'cloud', 'linux', '#FCC624', 4, true, 13),
    ('20000000-0000-0000-0000-000000000014', 'Vercel', 'vercel', 'cloud', 'vercel', '#FFFFFF', 5, true, 14),
    ('20000000-0000-0000-0000-000000000015', 'Git & GitHub', 'git', 'tool', 'git', '#F05032', 5, true, 15),
    ('20000000-0000-0000-0000-000000000016', 'OpenCV', 'opencv', 'tool', 'opencv', '#5C3EE8', 4, true, 16)
ON CONFLICT (slug) DO UPDATE
SET name = EXCLUDED.name,
    group_name = EXCLUDED.group_name,
    icon_key = EXCLUDED.icon_key,
    brand_color = EXCLUDED.brand_color,
    level = EXCLUDED.level,
    is_visible = EXCLUDED.is_visible,
    sort_order = EXCLUDED.sort_order;

-- 6. PROJECTS
INSERT INTO projects (
    id,
    category_id,
    title,
    slug,
    summary,
    description_md,
    description_html,
    repo_url,
    demo_url,
    status,
    is_featured,
    sort_order,
    started_at,
    ended_at
)
VALUES
    (
        '30000000-0000-0000-0000-000000000001',
        '10000000-0000-0000-0000-000000000001',
        'Autonomous Object Tracking & Detection System',
        'autonomous-object-tracking',
        'Real-time multi-camera object detection and tracking pipeline leveraging YOLOv8, ByteTrack, and OpenCV with GPU acceleration.',
        '### Project Overview
An edge-ready computer vision system capable of real-time object detection and trajectory tracking across multi-camera streams with high FPS and minimal latency.

### Tech Highlights
- **Model**: YOLOv8 & PyTorch
- **Tracking**: ByteTrack & DeepSORT
- **Backend**: Python FastAPI & WebSockets
- **Storage**: PostgreSQL with analytics tracking',
        '<h3>Project Overview</h3><p>An edge-ready computer vision system capable of real-time object detection and trajectory tracking across multi-camera streams with high FPS and minimal latency.</p><h3>Tech Highlights</h3><ul><li><strong>Model</strong>: YOLOv8 & PyTorch</li><li><strong>Tracking</strong>: ByteTrack & DeepSORT</li><li><strong>Backend</strong>: Python FastAPI & WebSockets</li><li><strong>Storage</strong>: PostgreSQL with analytics tracking</li></ul>',
        'https://github.com/Yogi-11234/object-tracking',
        'https://github.com/Yogi-11234',
        'published',
        true,
        1,
        '2024-01-10',
        '2024-05-20'
    ),
    (
        '30000000-0000-0000-0000-000000000002',
        '10000000-0000-0000-0000-000000000002',
        'High-Performance Microservices Engine',
        'high-performance-microservices-engine',
        'Distributed backend microservice framework built in Go with gRPC, PostgreSQL pooling, Redis caching, and JWT auth.',
        '### Architecture
A clean architecture Go backend designed for ultra-low latency, handling thousands of queries per second with pgx connection pooling and Redis caching.

### Key Features
- Sub-millisecond response times
- Double-submit CSRF and argon2id password security
- Rate limiting and structured slog logging',
        '<h3>Architecture</h3><p>A clean architecture Go backend designed for ultra-low latency, handling thousands of queries per second with pgx connection pooling and Redis caching.</p><h3>Key Features</h3><ul><li>Sub-millisecond response times</li><li>Double-submit CSRF and argon2id password security</li><li>Rate limiting and structured slog logging</li></ul>',
        'https://github.com/Yogi-11234/portfolio-app',
        'https://github.com/Yogi-11234',
        'published',
        true,
        2,
        '2024-03-01',
        '2024-07-15'
    ),
    (
        '30000000-0000-0000-0000-000000000003',
        '10000000-0000-0000-0000-000000000001',
        'Enterprise Document RAG Assistant',
        'enterprise-document-rag-assistant',
        'Intelligent Retrieval-Augmented Generation (RAG) assistant leveraging vector embeddings to query internal knowledge bases with source citations.',
        '### Overview
An AI-driven semantic search and Q&A engine that processes PDFs, spreadsheets, and technical docs with zero hallucinations by anchoring responses in retrieved document chunks.',
        '<h3>Overview</h3><p>An AI-driven semantic search and Q&A engine that processes PDFs, spreadsheets, and technical docs with zero hallucinations by anchoring responses in retrieved document chunks.</p>',
        'https://github.com/Yogi-11234',
        'https://github.com/Yogi-11234',
        'published',
        true,
        3,
        '2024-06-01',
        '2024-09-30'
    )
ON CONFLICT (slug) DO UPDATE
SET title = EXCLUDED.title,
    category_id = EXCLUDED.category_id,
    summary = EXCLUDED.summary,
    description_md = EXCLUDED.description_md,
    description_html = EXCLUDED.description_html,
    repo_url = EXCLUDED.repo_url,
    demo_url = EXCLUDED.demo_url,
    status = EXCLUDED.status,
    is_featured = EXCLUDED.is_featured,
    sort_order = EXCLUDED.sort_order;

-- 7. PROJECT_SKILLS (Pivot)
INSERT INTO project_skills (project_id, skill_id)
VALUES
    ('30000000-0000-0000-0000-000000000001', '20000000-0000-0000-0000-000000000002'), -- Python
    ('30000000-0000-0000-0000-000000000001', '20000000-0000-0000-0000-000000000007'), -- PyTorch
    ('30000000-0000-0000-0000-000000000001', '20000000-0000-0000-0000-000000000016'), -- OpenCV
    ('30000000-0000-0000-0000-000000000002', '20000000-0000-0000-0000-000000000001'), -- Go
    ('30000000-0000-0000-0000-000000000002', '20000000-0000-0000-0000-000000000009'), -- PostgreSQL
    ('30000000-0000-0000-0000-000000000002', '20000000-0000-0000-0000-000000000010'), -- Redis
    ('30000000-0000-0000-0000-000000000002', '20000000-0000-0000-0000-000000000012'), -- Docker
    ('30000000-0000-0000-0000-000000000003', '20000000-0000-0000-0000-000000000002'), -- Python
    ('30000000-0000-0000-0000-000000000003', '20000000-0000-0000-0000-000000000005'), -- FastAPI
    ('30000000-0000-0000-0000-000000000003', '20000000-0000-0000-0000-000000000011')  -- Supabase
ON CONFLICT (project_id, skill_id) DO NOTHING;

-- 8. EXPERIENCES
INSERT INTO experiences (
    id,
    type,
    organization,
    role,
    location,
    description_md,
    description_html,
    start_date,
    end_date,
    is_current,
    is_visible,
    sort_order
)
VALUES
    (
        '40000000-0000-0000-0000-000000000001',
        'work',
        'Technology Solutions & Innovation',
        'AI & Backend Software Engineer',
        'Jakarta, Indonesia',
        'Architected distributed microservices in Go and Python, deployed computer vision inference models, and managed cloud databases on PostgreSQL.',
        '<p>Architected distributed microservices in Go and Python, deployed computer vision inference models, and managed cloud databases on PostgreSQL.</p>',
        '2024-01-01',
        NULL,
        true,
        true,
        1
    ),
    (
        '40000000-0000-0000-0000-000000000002',
        'work',
        'Digital Software Studio',
        'Software Engineer Intern',
        'Remote',
        'Collaborated on fullstack web applications with React, TypeScript, and REST APIs. Implemented database schema migrations and automated testing.',
        '<p>Collaborated on fullstack web applications with React, TypeScript, and REST APIs. Implemented database schema migrations and automated testing.</p>',
        '2023-06-01',
        '2023-12-31',
        false,
        true,
        2
    ),
    (
        '40000000-0000-0000-0000-000000000003',
        'education',
        'Universitas',
        'Bachelor of Computer Science / Informatics',
        'Indonesia',
        'Specialized in Software Engineering, Algorithms & Data Structures, Artificial Intelligence, and Database Systems.',
        '<p>Specialized in Software Engineering, Algorithms & Data Structures, Artificial Intelligence, and Database Systems.</p>',
        '2021-08-01',
        '2025-07-31',
        true,
        true,
        3
    )
ON CONFLICT (id) DO UPDATE
SET organization = EXCLUDED.organization,
    role = EXCLUDED.role,
    location = EXCLUDED.location,
    description_md = EXCLUDED.description_md,
    description_html = EXCLUDED.description_html,
    start_date = EXCLUDED.start_date,
    end_date = EXCLUDED.end_date,
    is_current = EXCLUDED.is_current,
    is_visible = EXCLUDED.is_visible,
    sort_order = EXCLUDED.sort_order;

-- 9. CERTIFICATES
INSERT INTO certificates (
    id,
    title,
    issuer,
    issued_at,
    credential_url,
    file_path,
    file_mime,
    is_visible,
    sort_order
)
VALUES
    (
        '50000000-0000-0000-0000-000000000001',
        'Machine Learning Specialization',
        'DeepLearning.AI & Stanford Online',
        '2024-02-15',
        'https://coursera.org',
        'certificates/sample.pdf',
        'application/pdf',
        true,
        1
    ),
    (
        '50000000-0000-0000-0000-000000000002',
        'Backend Engineering & Cloud Architecture',
        'Google Cloud & Dicoding',
        '2023-11-20',
        'https://dicoding.com',
        'certificates/sample.pdf',
        'application/pdf',
        true,
        2
    ),
    (
        '50000000-0000-0000-0000-000000000003',
        'PostgreSQL Database Administration & Performance',
        'PostgreSQL Certification Community',
        '2023-08-10',
        'https://postgresql.org',
        'certificates/sample.pdf',
        'application/pdf',
        true,
        3
    )
ON CONFLICT (id) DO UPDATE
SET title = EXCLUDED.title,
    issuer = EXCLUDED.issuer,
    issued_at = EXCLUDED.issued_at,
    credential_url = EXCLUDED.credential_url,
    is_visible = EXCLUDED.is_visible,
    sort_order = EXCLUDED.sort_order;

-- 10. POSTS
INSERT INTO posts (
    id,
    author_id,
    category_id,
    title,
    slug,
    excerpt,
    content_md,
    content_html,
    status,
    reading_minutes,
    published_at
)
VALUES
    (
        '60000000-0000-0000-0000-000000000001',
        '00000000-0000-0000-0000-000000000001',
        '10000000-0000-0000-0000-000000000004',
        'Architecting High-Throughput REST APIs with Go and Fiber',
        'architecting-high-throughput-rest-apis-with-go',
        'Learn best practices for structuring Go backend services with clean architecture, pgx connection pooling, and resilient middleware.',
        '## Introduction
When building distributed backend services, performance, simplicity, and low memory usage are paramount. Go provides an ideal foundation.

### Database Connection Pooling
Using pgx/v5 allows zero-overhead connection pooling that handles spikes seamlessly.',
        '<h2>Introduction</h2><p>When building distributed backend services, performance, simplicity, and low memory usage are paramount. Go provides an ideal foundation.</p><h3>Database Connection Pooling</h3><p>Using pgx/v5 allows zero-overhead connection pooling that handles spikes seamlessly.</p>',
        'published',
        5,
        NOW() - INTERVAL '3 days'
    ),
    (
        '60000000-0000-0000-0000-000000000002',
        '00000000-0000-0000-0000-000000000001',
        '10000000-0000-0000-0000-000000000005',
        'Optimizing Computer Vision Pipelines for Low-Power Edge Devices',
        'optimizing-computer-vision-pipelines-for-edge-devices',
        'A practical guide to model quantization, TensorRT acceleration, and frame pipelining for real-time edge AI inference.',
        '## The Edge AI Challenge
Running complex convolutional and transformer models on edge processors requires aggressive optimization without sacrificing detection accuracy.',
        '<h2>The Edge AI Challenge</h2><p>Running complex convolutional and transformer models on edge processors requires aggressive optimization without sacrificing detection accuracy.</p>',
        'published',
        7,
        NOW() - INTERVAL '7 days'
    )
ON CONFLICT (slug) DO UPDATE
SET title = EXCLUDED.title,
    excerpt = EXCLUDED.excerpt,
    content_md = EXCLUDED.content_md,
    content_html = EXCLUDED.content_html,
    status = EXCLUDED.status,
    published_at = EXCLUDED.published_at;

-- SELESAI
