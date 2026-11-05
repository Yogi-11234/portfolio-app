export interface ApiResponse<T> {
  data: T;
  meta?: {
    page: number;
    per_page: number;
    total: number;
  };
  error?: {
    code: string;
    message: string;
    fields?: Record<string, string>;
  };
}

export interface Profile {
  id: string;
  user_id: string;
  full_name: string;
  headline: string;
  bio_md: string;
  bio_html: string;
  avatar_path?: string;
  resume_path?: string;
  location?: string;
  public_email?: string;
  github_url?: string;
  linkedin_url?: string;
  social_links: Record<string, string>;
  open_to_work: boolean;
  created_at: string;
  updated_at: string;
}

export interface Category {
  id: string;
  name: string;
  slug: string;
  type: 'project' | 'post';
  created_at: string;
}

export type SkillGroup = 'language' | 'framework' | 'database' | 'tool' | 'cloud' | 'other';

export interface Skill {
  id: string;
  name: string;
  slug: string;
  group_name: SkillGroup;
  icon_name?: string;
  icon_key?: string;
  brand_color?: string;
  level: number;
  is_visible: boolean;
  order_index: number;
  sort_order?: number;
  created_at?: string;
}

export interface SkillBrief {
  id: string;
  name: string;
  slug: string;
  group_name?: SkillGroup;
  icon_name?: string;
  icon_key?: string;
  brand_color?: string;
  level?: number;
}

export interface Project {
  id: string;
  category_id?: string;
  title: string;
  slug: string;
  summary: string;
  description_md: string;
  description_html: string;
  thumbnail_path?: string;
  demo_url?: string;
  repo_url?: string;
  featured: boolean;
  status: 'draft' | 'published' | 'archived';
  started_at?: string;
  finished_at?: string;
  created_at: string;
  updated_at: string;
  category?: {
    id: string;
    name: string;
    slug: string;
  };
  skills: SkillBrief[];
}

export type ExperienceType = 'work' | 'education' | 'organization' | 'volunteer';

export interface Experience {
  id: string;
  type: ExperienceType;
  title: string;
  organization: string;
  location?: string;
  started_at: string;
  ended_at?: string;
  current: boolean;
  description_md?: string;
  description_html?: string;
  order_index: number;
  is_visible: boolean;
  created_at: string;
}

export interface Certificate {
  id: string;
  name: string;
  title?: string;
  issuer: string;
  issued_at: string;
  expires_at?: string;
  credential_url?: string;
  file_path?: string;
  file_mime?: string;
  order_index: number;
  sort_order?: number;
  is_visible: boolean;
  created_at?: string;
}

export interface Post {
  id: string;
  author_id: string;
  author_name?: string;
  category_id?: string;
  title: string;
  slug: string;
  excerpt: string;
  content_md: string;
  content_html: string;
  cover_path?: string;
  status: 'draft' | 'published' | 'archived';
  reading_minutes: number;
  published_at?: string;
  created_at: string;
  updated_at: string;
  category?: {
    id: string;
    name: string;
    slug: string;
  };
}

export interface ContactPayload {
  name: string;
  email: string;
  subject?: string;
  body: string;
  website?: string; // honeypot
}

export interface User {
  id: string;
  name: string;
  email: string;
  role: string;
  created_at: string;
}

export interface Message {
  id: string;
  name: string;
  email: string;
  subject: string;
  body: string;
  status: 'unread' | 'read' | 'archived' | 'spam';
  user_agent?: string;
  created_at: string;
  read_at?: string;
}

export interface UploadResult {
  path: string;
  public_url: string;
  mime: string;
  size: number;
}

