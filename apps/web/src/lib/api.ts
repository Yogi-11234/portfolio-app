import { supabase } from '@/lib/supabase';
import {
  Profile,
  Category,
  Skill,
  Project,
  Experience,
  Certificate,
  Post,
  ContactPayload,
  User,
  Message,
  UploadResult,
} from '@/types';

export class ApiError extends Error {
  code: string;
  fields?: Record<string, string>;

  constructor(message: string, code: string = 'UNKNOWN_ERROR', fields?: Record<string, string>) {
    super(message);
    this.name = 'ApiError';
    this.code = code;
    this.fields = fields;
  }
}

function slugify(text: string): string {
  return text
    .toLowerCase()
    .trim()
    .replace(/[^\w\s-]/g, '')
    .replace(/[\s_-]+/g, '-')
    .replace(/^-+|-+$/g, '');
}

function mdToHtml(md?: string): string {
  if (!md) return '';
  return md
    .split(/\n\n+/)
    .map((p) => `<p>${p.replace(/\n/g, '<br/>')}</p>`)
    .join('');
}

export const api = {
  // ---- AUTH (Supabase Auth SDK) ----
  login: async (credentials: { email: string; password: string }): Promise<User> => {
    const { data, error } = await supabase.auth.signInWithPassword({
      email: credentials.email,
      password: credentials.password,
    });

    if (error) {
      throw new ApiError(error.message, 'AUTH_ERROR');
    }

    const u = data.user;
    return {
      id: u.id,
      email: u.email || '',
      name: u.user_metadata?.name || u.email?.split('@')[0] || 'Admin',
      role: 'admin',
      created_at: u.created_at,
    };
  },

  logout: async () => {
    const { error } = await supabase.auth.signOut();
    if (error) {
      throw new ApiError(error.message, 'LOGOUT_ERROR');
    }
    return { message: 'Berhasil logout' };
  },

  me: async (): Promise<User> => {
    const {
      data: { session },
      error,
    } = await supabase.auth.getSession();

    if (error || !session?.user) {
      throw new ApiError('Sesi login telah berakhir', 'UNAUTHORIZED');
    }

    const u = session.user;
    return {
      id: u.id,
      email: u.email || '',
      name: u.user_metadata?.name || u.email?.split('@')[0] || 'Admin',
      role: 'admin',
      created_at: u.created_at,
    };
  },

  // ---- PUBLIC CONTENT ----
  getProfile: async (): Promise<Profile> => {
    const { data, error } = await supabase
      .from('profile')
      .select('*')
      .limit(1)
      .maybeSingle();

    if (error) {
      throw new ApiError(error.message);
    }
    if (!data) {
      return null as unknown as Profile;
    }

    return {
      ...data,
      bio_html: data.bio_html || mdToHtml(data.bio_md),
      social_links: data.social_links || {},
    } as Profile;
  },

  getCategories: async (type?: 'project' | 'post'): Promise<Category[]> => {
    let query = supabase.from('categories').select('*').order('sort_order', { ascending: true });
    if (type) {
      query = query.eq('type', type);
    }

    const { data, error } = await query;
    if (error) throw new ApiError(error.message);
    return (data || []) as Category[];
  },

  getSkills: async (): Promise<Skill[]> => {
    const { data, error } = await supabase
      .from('skills')
      .select('*')
      .eq('is_visible', true)
      .order('sort_order', { ascending: true });

    if (error) throw new ApiError(error.message);
    return (data || []).map((s: any) => ({
      ...s,
      order_index: s.order_index ?? s.sort_order ?? 0,
    })) as Skill[];
  },

  getProjects: async (params?: {
    category?: string;
    skill?: string;
    featured?: boolean;
    page?: number;
  }) => {
    let query = supabase
      .from('projects')
      .select('*, category:categories(id, name, slug), project_skills(skill:skills(*))', { count: 'exact' })
      .eq('status', 'published')
      .order('is_featured', { ascending: false })
      .order('sort_order', { ascending: true })
      .order('created_at', { ascending: false });

    if (params?.featured !== undefined) {
      query = query.eq('is_featured', params.featured);
    }
    if (params?.category) {
      query = query.eq('category_id', params.category);
    }

    const { data, error, count } = await query;
    if (error) throw new ApiError(error.message);

    let projects = (data || []).map((p: any) => ({
      ...p,
      featured: p.featured ?? p.is_featured ?? false,
      category: p.category || undefined,
      skills: (p.project_skills || []).map((ps: any) => ps.skill || ps.skills).filter(Boolean),
    })) as Project[];

    if (params?.skill) {
      projects = projects.filter((p) =>
        p.skills.some((s) => s.slug === params.skill || s.id === params.skill)
      );
    }

    return {
      projects,
      meta: {
        page: params?.page || 1,
        per_page: projects.length,
        total: count ?? projects.length,
      },
    };
  },

  getProjectBySlug: async (slug: string): Promise<Project> => {
    const { data, error } = await supabase
      .from('projects')
      .select('*, category:categories(id, name, slug), project_skills(skill:skills(*))')
      .eq('slug', slug)
      .single();

    if (error) throw new ApiError(error.message);
    return {
      ...data,
      featured: data.featured ?? data.is_featured ?? false,
      category: data.category || undefined,
      skills: (data.project_skills || []).map((ps: any) => ps.skill || ps.skills).filter(Boolean),
    } as Project;
  },

  getExperiences: async (): Promise<Experience[]> => {
    const { data, error } = await supabase
      .from('experiences')
      .select('*')
      .eq('is_visible', true)
      .order('sort_order', { ascending: true })
      .order('start_date', { ascending: false });

    if (error) throw new ApiError(error.message);
    return (data || []).map((e: any) => ({
      ...e,
      title: e.title || e.role || '',
      role: e.role || e.title || '',
      started_at: e.started_at || e.start_date || '',
      ended_at: e.ended_at || e.end_date || '',
      current: e.current ?? e.is_current ?? false,
      order_index: e.order_index ?? e.sort_order ?? 0,
    })) as Experience[];
  },

  getCertificates: async (): Promise<Certificate[]> => {
    const { data, error } = await supabase
      .from('certificates')
      .select('*')
      .eq('is_visible', true)
      .order('sort_order', { ascending: true })
      .order('issued_at', { ascending: false });

    if (error) throw new ApiError(error.message);
    return (data || []).map((c: any) => ({
      ...c,
      name: c.name || c.title || '',
      title: c.title || c.name || '',
      order_index: c.order_index ?? c.sort_order ?? 0,
    })) as Certificate[];
  },

  getPosts: async (params?: { category?: string; page?: number }) => {
    let query = supabase
      .from('posts')
      .select('*, category:categories(id, name, slug)', { count: 'exact' })
      .eq('status', 'published')
      .order('published_at', { ascending: false });

    if (params?.category) {
      query = query.eq('category_id', params.category);
    }

    const { data, error, count } = await query;
    if (error) throw new ApiError(error.message);

    const posts = (data || []).map((p: any) => ({
      ...p,
      category: p.category || undefined,
    })) as Post[];

    return {
      posts,
      meta: {
        page: params?.page || 1,
        per_page: posts.length,
        total: count ?? posts.length,
      },
    };
  },

  getPostBySlug: async (slug: string): Promise<Post> => {
    const { data, error } = await supabase
      .from('posts')
      .select('*, category:categories(id, name, slug)')
      .eq('slug', slug)
      .single();

    if (error) throw new ApiError(error.message);
    return {
      ...data,
      category: data.category || undefined,
    } as Post;
  },

  sendContactMessage: async (payload: ContactPayload) => {
    if (payload.website) {
      return { message: 'Pesan berhasil dikirim' };
    }

    const { error } = await supabase.from('messages').insert({
      name: payload.name,
      email: payload.email,
      subject: payload.subject || '',
      body: payload.body,
      status: 'unread',
    });

    if (error) throw new ApiError(error.message);
    return { message: 'Pesan berhasil dikirim' };
  },

  // ---- ADMIN MANAGEMENT ----
  updateProfile: async (data: Partial<Profile>): Promise<Profile> => {
    const { id, ...updates } = data;
    if (updates.bio_md && !updates.bio_html) {
      updates.bio_html = mdToHtml(updates.bio_md);
    }

    const { data: existing } = await supabase.from('profile').select('id').limit(1).maybeSingle();
    let query;
    if (existing?.id) {
      query = supabase.from('profile').update(updates).eq('id', existing.id).select().single();
    } else {
      query = supabase.from('profile').insert(updates).select().single();
    }

    const { data: result, error } = await query;
    if (error) throw new ApiError(error.message);
    return result as Profile;
  },

  // Categories (admin)
  createCategory: async (data: { name: string; slug?: string; type: 'project' | 'post' }): Promise<Category> => {
    const payload = {
      name: data.name,
      slug: data.slug || slugify(data.name),
      type: data.type,
    };
    const { data: result, error } = await supabase.from('categories').insert(payload).select().single();
    if (error) throw new ApiError(error.message);
    return result as Category;
  },

  updateCategory: async (
    id: string,
    data: { name: string; slug?: string; type: 'project' | 'post' }
  ): Promise<Category> => {
    const payload = {
      name: data.name,
      slug: data.slug || slugify(data.name),
      type: data.type,
    };
    const { data: result, error } = await supabase.from('categories').update(payload).eq('id', id).select().single();
    if (error) throw new ApiError(error.message);
    return result as Category;
  },

  deleteCategory: async (id: string) => {
    const { error } = await supabase.from('categories').delete().eq('id', id);
    if (error) throw new ApiError(error.message);
    return { message: 'Kategori berhasil dihapus' };
  },

  // Skills (admin)
  getSkillsAdmin: async (): Promise<Skill[]> => {
    const { data, error } = await supabase.from('skills').select('*').order('sort_order', { ascending: true });
    if (error) throw new ApiError(error.message);
    return (data || []).map((s: any) => ({
      ...s,
      order_index: s.order_index ?? s.sort_order ?? 0,
    })) as Skill[];
  },

  createSkill: async (data: any): Promise<Skill> => {
    const payload = {
      name: data.name,
      slug: data.slug || slugify(data.name),
      group_name: data.group_name || 'tool',
      icon_key: data.icon_key,
      brand_color: data.brand_color,
      level: data.level || 5,
      is_visible: data.is_visible ?? true,
      sort_order: Number(data.sort_order ?? data.order_index ?? 0),
    };
    const { data: result, error } = await supabase.from('skills').insert(payload).select().single();
    if (error) throw new ApiError(error.message);
    return result as Skill;
  },

  updateSkill: async (id: string, data: any): Promise<Skill> => {
    const payload: any = {
      name: data.name,
      group_name: data.group_name,
      icon_key: data.icon_key,
      brand_color: data.brand_color,
      level: data.level || 5,
      is_visible: data.is_visible,
      sort_order: Number(data.sort_order ?? data.order_index ?? 0),
    };
    if (data.slug) payload.slug = data.slug;

    const { data: result, error } = await supabase.from('skills').update(payload).eq('id', id).select().single();
    if (error) throw new ApiError(error.message);
    return result as Skill;
  },

  deleteSkill: async (id: string) => {
    const { error } = await supabase.from('skills').delete().eq('id', id);
    if (error) throw new ApiError(error.message);
    return { message: 'Skill berhasil dihapus' };
  },

  // Projects (admin)
  getProjectsAdmin: async (params?: { status?: string; page?: number; per_page?: number }) => {
    let query = supabase
      .from('projects')
      .select('*, category:categories(id, name, slug), project_skills(skill:skills(*))', { count: 'exact' })
      .order('created_at', { ascending: false });

    if (params?.status) {
      query = query.eq('status', params.status);
    }
    const { data, error, count } = await query;
    if (error) throw new ApiError(error.message);

    const projects = (data || []).map((p: any) => ({
      ...p,
      featured: p.featured ?? p.is_featured ?? false,
      category: p.category || undefined,
      skills: (p.project_skills || []).map((ps: any) => ps.skill || ps.skills).filter(Boolean),
    })) as Project[];

    return {
      projects,
      meta: {
        page: params?.page || 1,
        per_page: projects.length,
        total: count ?? projects.length,
      },
    };
  },

  getProjectByID: async (id: string): Promise<Project> => {
    const { data, error } = await supabase
      .from('projects')
      .select('*, category:categories(id, name, slug), project_skills(skill:skills(*))')
      .eq('id', id)
      .single();

    if (error) throw new ApiError(error.message);
    return {
      ...data,
      featured: data.featured ?? data.is_featured ?? false,
      category: data.category || undefined,
      skills: (data.project_skills || []).map((ps: any) => ps.skill || ps.skills).filter(Boolean),
    } as Project;
  },

  createProject: async (data: any): Promise<Project> => {
    const skillIds = data.skill_ids || [];
    const payload = {
      title: data.title,
      slug: data.slug || slugify(data.title),
      category_id: data.category_id || null,
      summary: data.summary,
      description_md: data.description_md || '',
      description_html: data.description_html || mdToHtml(data.description_md),
      thumbnail_path: data.thumbnail_path,
      demo_url: data.demo_url || null,
      repo_url: data.repo_url || null,
      is_featured: data.is_featured ?? data.featured ?? false,
      status: data.status || 'draft',
      started_at: data.started_at ? data.started_at.split('T')[0] : null,
      ended_at: (data.ended_at || data.finished_at) ? (data.ended_at || data.finished_at).split('T')[0] : null,
    };

    const { data: project, error } = await supabase.from('projects').insert(payload).select().single();
    if (error) throw new ApiError(error.message);

    if (skillIds.length > 0) {
      await supabase.from('project_skills').insert(
        skillIds.map((sId: string) => ({ project_id: project.id, skill_id: sId }))
      );
    }
    return project as Project;
  },

  updateProject: async (id: string, data: any): Promise<Project> => {
    const skillIds = data.skill_ids;
    const payload: any = {
      title: data.title,
      summary: data.summary,
      description_md: data.description_md,
      status: data.status,
    };
    if (data.slug) payload.slug = data.slug;
    if (data.category_id !== undefined) payload.category_id = data.category_id || null;
    if (data.description_md) payload.description_html = data.description_html || mdToHtml(data.description_md);
    if (data.thumbnail_path !== undefined) payload.thumbnail_path = data.thumbnail_path;
    if (data.demo_url !== undefined) payload.demo_url = data.demo_url || null;
    if (data.repo_url !== undefined) payload.repo_url = data.repo_url || null;
    if (data.is_featured !== undefined || data.featured !== undefined) {
      payload.is_featured = data.is_featured ?? data.featured;
    }
    if (data.started_at !== undefined) {
      payload.started_at = data.started_at ? data.started_at.split('T')[0] : null;
    }
    if (data.ended_at !== undefined || data.finished_at !== undefined) {
      const eDate = data.ended_at || data.finished_at;
      payload.ended_at = eDate ? eDate.split('T')[0] : null;
    }

    const { data: project, error } = await supabase.from('projects').update(payload).eq('id', id).select().single();
    if (error) throw new ApiError(error.message);

    if (Array.isArray(skillIds)) {
      await supabase.from('project_skills').delete().eq('project_id', id);
      if (skillIds.length > 0) {
        await supabase.from('project_skills').insert(
          skillIds.map((sId: string) => ({ project_id: id, skill_id: sId }))
        );
      }
    }
    return project as Project;
  },

  updateProjectSkills: async (id: string, skillIDs: string[]) => {
    await supabase.from('project_skills').delete().eq('project_id', id);
    if (skillIDs.length > 0) {
      await supabase.from('project_skills').insert(
        skillIDs.map((sId: string) => ({ project_id: id, skill_id: sId }))
      );
    }
    return { message: 'Skills berhasil diperbarui' };
  },

  deleteProject: async (id: string) => {
    const { error } = await supabase.from('projects').delete().eq('id', id);
    if (error) throw new ApiError(error.message);
    return { message: 'Proyek berhasil dihapus' };
  },

  // Experiences (admin)
  getExperiencesAdmin: async (): Promise<Experience[]> => {
    const { data, error } = await supabase
      .from('experiences')
      .select('*')
      .order('sort_order', { ascending: true })
      .order('start_date', { ascending: false });

    if (error) throw new ApiError(error.message);
    return (data || []).map((e: any) => ({
      ...e,
      title: e.title || e.role || '',
      role: e.role || e.title || '',
      started_at: e.started_at || e.start_date || '',
      ended_at: e.ended_at || e.end_date || '',
      current: e.current ?? e.is_current ?? false,
      order_index: e.order_index ?? e.sort_order ?? 0,
    })) as Experience[];
  },

  createExperience: async (data: any): Promise<Experience> => {
    const sDate = data.start_date || data.started_at;
    const eDate = data.end_date || data.ended_at;
    const payload = {
      type: data.type,
      organization: data.organization,
      role: data.role || data.title,
      location: data.location || null,
      description_md: data.description_md || '',
      description_html: data.description_html || mdToHtml(data.description_md),
      start_date: sDate ? sDate.split('T')[0] : new Date().toISOString().split('T')[0],
      end_date: eDate ? eDate.split('T')[0] : null,
      is_current: data.is_current ?? data.current ?? false,
      is_visible: data.is_visible ?? true,
      sort_order: Number(data.sort_order ?? data.order_index ?? 0),
    };

    const { data: result, error } = await supabase.from('experiences').insert(payload).select().single();
    if (error) throw new ApiError(error.message);
    return result as Experience;
  },

  updateExperience: async (id: string, data: any): Promise<Experience> => {
    const sDate = data.start_date || data.started_at;
    const eDate = data.end_date || data.ended_at;
    const payload: any = {
      type: data.type,
      organization: data.organization,
      role: data.role || data.title,
      location: data.location || null,
      description_md: data.description_md,
      description_html: data.description_html || mdToHtml(data.description_md),
      is_current: data.is_current ?? data.current ?? false,
      is_visible: data.is_visible,
      sort_order: Number(data.sort_order ?? data.order_index ?? 0),
    };
    if (sDate) payload.start_date = sDate.split('T')[0];
    if (eDate !== undefined) payload.end_date = eDate ? eDate.split('T')[0] : null;

    const { data: result, error } = await supabase.from('experiences').update(payload).eq('id', id).select().single();
    if (error) throw new ApiError(error.message);
    return result as Experience;
  },

  deleteExperience: async (id: string) => {
    const { error } = await supabase.from('experiences').delete().eq('id', id);
    if (error) throw new ApiError(error.message);
    return { message: 'Pengalaman berhasil dihapus' };
  },

  // Certificates (admin)
  getCertificatesAdmin: async (): Promise<Certificate[]> => {
    const { data, error } = await supabase
      .from('certificates')
      .select('*')
      .order('sort_order', { ascending: true })
      .order('issued_at', { ascending: false });

    if (error) throw new ApiError(error.message);
    return (data || []).map((c: any) => ({
      ...c,
      name: c.name || c.title || '',
      title: c.title || c.name || '',
      order_index: c.order_index ?? c.sort_order ?? 0,
    })) as Certificate[];
  },

  createCertificate: async (data: any): Promise<Certificate> => {
    const iDate = data.issued_at ? data.issued_at.split('T')[0] : new Date().toISOString().split('T')[0];
    const eDate = data.expires_at ? data.expires_at.split('T')[0] : null;
    const payload = {
      title: data.name || data.title,
      issuer: data.issuer,
      issued_at: iDate,
      expires_at: eDate,
      credential_id: data.credential_id || null,
      credential_url: data.credential_url || null,
      file_path: data.file_path,
      file_mime: data.file_mime || 'application/pdf',
      thumbnail_path: data.thumbnail_path || null,
      is_visible: data.is_visible ?? true,
      sort_order: Number(data.sort_order ?? data.order_index ?? 0),
    };

    const { data: result, error } = await supabase.from('certificates').insert(payload).select().single();
    if (error) throw new ApiError(error.message);
    return result as Certificate;
  },

  updateCertificate: async (id: string, data: any): Promise<Certificate> => {
    const iDate = data.issued_at ? data.issued_at.split('T')[0] : undefined;
    const eDate = data.expires_at ? data.expires_at.split('T')[0] : null;
    const payload: any = {
      title: data.name || data.title,
      issuer: data.issuer,
      credential_id: data.credential_id || null,
      credential_url: data.credential_url || null,
      file_path: data.file_path,
      file_mime: data.file_mime || 'application/pdf',
      is_visible: data.is_visible,
      sort_order: Number(data.sort_order ?? data.order_index ?? 0),
    };
    if (iDate) payload.issued_at = iDate;
    if (data.expires_at !== undefined) payload.expires_at = eDate;
    if (data.thumbnail_path !== undefined) payload.thumbnail_path = data.thumbnail_path;

    const { data: result, error } = await supabase.from('certificates').update(payload).eq('id', id).select().single();
    if (error) throw new ApiError(error.message);
    return result as Certificate;
  },

  deleteCertificate: async (id: string) => {
    const { error } = await supabase.from('certificates').delete().eq('id', id);
    if (error) throw new ApiError(error.message);
    return { message: 'Sertifikat berhasil dihapus' };
  },

  // Posts (admin)
  getPostsAdmin: async (params?: { status?: string; page?: number; per_page?: number }) => {
    let query = supabase
      .from('posts')
      .select('*, category:categories(id, name, slug)', { count: 'exact' })
      .order('created_at', { ascending: false });

    if (params?.status) {
      query = query.eq('status', params.status);
    }
    const { data, error, count } = await query;
    if (error) throw new ApiError(error.message);

    const posts = (data || []).map((p: any) => ({
      ...p,
      category: p.category || undefined,
    })) as Post[];

    return {
      posts,
      meta: {
        page: params?.page || 1,
        per_page: posts.length,
        total: count ?? posts.length,
      },
    };
  },

  getPostByID: async (id: string): Promise<Post> => {
    const { data, error } = await supabase
      .from('posts')
      .select('*, category:categories(id, name, slug)')
      .eq('id', id)
      .single();

    if (error) throw new ApiError(error.message);
    return {
      ...data,
      category: data.category || undefined,
    } as Post;
  },

  createPost: async (data: any): Promise<Post> => {
    const {
      data: { user },
    } = await supabase.auth.getUser();

    let authorId = user?.id;
    if (!authorId) {
      const { data: anyUser } = await supabase.from('users').select('id').limit(1).maybeSingle();
      authorId = anyUser?.id;
    }

    const payload = {
      author_id: authorId,
      title: data.title,
      slug: data.slug || slugify(data.title),
      category_id: data.category_id || null,
      excerpt: data.excerpt || '',
      content_md: data.content_md || '',
      content_html: data.content_html || mdToHtml(data.content_md),
      cover_path: data.cover_path || null,
      status: data.status || 'draft',
      reading_minutes: Number(data.reading_minutes) || 1,
      published_at: data.status === 'published' ? new Date().toISOString() : null,
    };

    const { data: result, error } = await supabase.from('posts').insert(payload).select().single();
    if (error) throw new ApiError(error.message);
    return result as Post;
  },

  updatePost: async (id: string, data: any): Promise<Post> => {
    const payload: any = {
      title: data.title,
      category_id: data.category_id || null,
      excerpt: data.excerpt,
      content_md: data.content_md,
      content_html: data.content_html || mdToHtml(data.content_md),
      cover_path: data.cover_path || null,
      status: data.status,
      reading_minutes: Number(data.reading_minutes) || 1,
    };
    if (data.slug) payload.slug = data.slug;
    if (data.status === 'published' && !data.published_at) {
      payload.published_at = new Date().toISOString();
    }

    const { data: result, error } = await supabase.from('posts').update(payload).eq('id', id).select().single();
    if (error) throw new ApiError(error.message);
    return result as Post;
  },

  deletePost: async (id: string) => {
    const { error } = await supabase.from('posts').delete().eq('id', id);
    if (error) throw new ApiError(error.message);
    return { message: 'Postingan berhasil dihapus' };
  },

  // Messages (admin inbox)
  getMessagesAdmin: async (params?: { status?: string; page?: number; per_page?: number }) => {
    let query = supabase.from('messages').select('*', { count: 'exact' }).order('created_at', { ascending: false });
    if (params?.status) {
      query = query.eq('status', params.status);
    }
    const { data, error, count } = await query;
    if (error) throw new ApiError(error.message);

    return {
      messages: (data || []) as Message[],
      meta: {
        page: params?.page || 1,
        per_page: (data || []).length,
        total: count ?? (data || []).length,
      },
    };
  },

  getMessageByID: async (id: string): Promise<Message> => {
    const { data, error } = await supabase.from('messages').select('*').eq('id', id).single();
    if (error) throw new ApiError(error.message);
    return data as Message;
  },

  updateMessageStatus: async (id: string, status: 'unread' | 'read' | 'archived' | 'spam') => {
    const payload: any = { status };
    if (status === 'read') {
      payload.read_at = new Date().toISOString();
    }
    const { error } = await supabase.from('messages').update(payload).eq('id', id);
    if (error) throw new ApiError(error.message);
    return { message: 'Status pesan berhasil diperbarui' };
  },

  deleteMessage: async (id: string) => {
    const { error } = await supabase.from('messages').delete().eq('id', id);
    if (error) throw new ApiError(error.message);
    return { message: 'Pesan berhasil dihapus' };
  },

  // Uploads via Supabase Storage
  uploadMedia: async (
    kind: 'avatar' | 'resume' | 'project' | 'certificate' | 'post' | 'skill',
    file: File
  ): Promise<UploadResult> => {
    const ext = file.name.split('.').pop()?.toLowerCase() || 'png';
    const filename = `${crypto.randomUUID()}.${ext}`;
    const folder = kind.endsWith('s') ? kind : `${kind}s`;
    const relativePath = `${folder}/${filename}`;

    const { error } = await supabase.storage
      .from('portfolio_assets')
      .upload(relativePath, file, {
        cacheControl: '3600',
        upsert: true,
      });

    if (error) {
      throw new ApiError(error.message, 'UPLOAD_ERROR');
    }

    const {
      data: { publicUrl },
    } = supabase.storage.from('portfolio_assets').getPublicUrl(relativePath);

    return {
      path: relativePath,
      public_url: publicUrl,
      mime: file.type,
      size: file.size,
    };
  },

  deleteMedia: async (path: string) => {
    const cleanPath = path.replace(/^\/?portfolio_assets\//, '').replace(/^\/?uploads\//, '');
    const { error } = await supabase.storage.from('portfolio_assets').remove([cleanPath]);
    if (error) {
      throw new ApiError(error.message, 'DELETE_ERROR');
    }
    return { message: 'File berhasil dihapus' };
  },
};
