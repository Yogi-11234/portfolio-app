-- 000001_init_schema.down.sql
DROP TRIGGER IF EXISTS trg_posts_updated ON posts;
DROP TRIGGER IF EXISTS trg_certificates_updated ON certificates;
DROP TRIGGER IF EXISTS trg_experiences_updated ON experiences;
DROP TRIGGER IF EXISTS trg_projects_updated ON projects;
DROP TRIGGER IF EXISTS trg_profile_updated ON profile;
DROP TRIGGER IF EXISTS trg_users_updated ON users;

DROP FUNCTION IF EXISTS set_updated_at();

DROP TABLE IF EXISTS messages;
DROP TABLE IF EXISTS posts;
DROP TABLE IF EXISTS certificates;
DROP TABLE IF EXISTS experiences;
DROP TABLE IF EXISTS project_skills;
DROP TABLE IF EXISTS skills;
DROP TABLE IF EXISTS projects;
DROP TABLE IF EXISTS categories;
DROP TABLE IF EXISTS profile;
DROP TABLE IF EXISTS users;
