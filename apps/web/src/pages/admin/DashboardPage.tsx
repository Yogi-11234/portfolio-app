import React from 'react';
import { useQuery } from '@tanstack/react-query';
import { Link } from 'react-router-dom';
import {
  Layers,
  Code2,
  Award,
  BookOpen,
  Inbox,
  Plus,
  ArrowRight,
  User,
  GraduationCap,
} from 'lucide-react';
import { api } from '@/lib/api';
import { AnimatedButton } from '@/components/AnimatedButton';

export const DashboardPage: React.FC = () => {
  const { data: projectsData } = useQuery({
    queryKey: ['admin-projects-stats'],
    queryFn: () => api.getProjectsAdmin({ per_page: 100 }),
  });

  const { data: skills = [] } = useQuery({
    queryKey: ['admin-skills'],
    queryFn: api.getSkillsAdmin,
  });

  const { data: experiences = [] } = useQuery({
    queryKey: ['admin-experiences'],
    queryFn: api.getExperiencesAdmin,
  });

  const { data: certificates = [] } = useQuery({
    queryKey: ['admin-certificates'],
    queryFn: api.getCertificatesAdmin,
  });

  const { data: postsData } = useQuery({
    queryKey: ['admin-posts-stats'],
    queryFn: () => api.getPostsAdmin({ per_page: 100 }),
  });

  const { data: messagesData } = useQuery({
    queryKey: ['admin-messages-stats'],
    queryFn: () => api.getMessagesAdmin({ per_page: 100 }),
  });

  const projects = projectsData?.projects || [];
  const publishedProjects = projects.filter((p) => p.status === 'published').length;

  const posts = postsData?.posts || [];
  const publishedPosts = posts.filter((p) => p.status === 'published').length;

  const messages = messagesData?.messages || [];
  const unreadMessages = messages.filter((m) => m.status === 'unread').length;

  const stats = [
    {
      label: 'Total Proyek',
      value: projects.length,
      sub: `${publishedProjects} dipublikasikan`,
      icon: Layers,
      href: '/admin/projects',
      color: 'text-sky-400',
    },
    {
      label: 'Keahlian / Skills',
      value: skills.length,
      sub: `${skills.filter((s) => s.is_visible).length} aktif`,
      icon: Code2,
      href: '/admin/skills',
      color: 'text-emerald-400',
    },
    {
      label: 'Artikel & Catatan',
      value: posts.length,
      sub: `${publishedPosts} tayang`,
      icon: BookOpen,
      href: '/admin/posts',
      color: 'text-amber-400',
    },
    {
      label: 'Sertifikat & Lisensi',
      value: certificates.length,
      sub: 'Kredensial',
      icon: Award,
      href: '/admin/certificates',
      color: 'text-indigo-400',
    },
    {
      label: 'Pengalaman',
      value: experiences.length,
      sub: 'Riwayat karir & edukasi',
      icon: GraduationCap,
      href: '/admin/experiences',
      color: 'text-teal-400',
    },
    {
      label: 'Pesan Masuk',
      value: messages.length,
      sub: `${unreadMessages} belum dibaca`,
      icon: Inbox,
      href: '/admin/messages',
      color: unreadMessages > 0 ? 'text-rose-400 font-bold' : 'text-text-muted',
    },
  ];

  return (
    <div className="space-y-8">
      {/* Header */}
      <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
        <div>
          <h1 className="text-2xl font-bold text-text">Ringkasan Sistem</h1>
          <p className="text-xs text-text-muted mt-1">
            Status seluruh konten portofolio dan aktivitas pesan masuk
          </p>
        </div>

        {/* Quick Actions */}
        <div className="flex flex-wrap items-center gap-3">
          <AnimatedButton as="link" to="/admin/projects" size="sm">
            <Plus className="w-3.5 h-3.5 inline mr-1" />
            <span>Tambah Proyek</span>
          </AnimatedButton>
          <AnimatedButton as="link" to="/admin/posts" size="sm">
            <Plus className="w-3.5 h-3.5 inline mr-1" />
            <span>Tulis Artikel</span>
          </AnimatedButton>
          <AnimatedButton as="link" to="/admin/profile" size="sm">
            <User className="w-3.5 h-3.5 inline mr-1" />
            <span>Edit Profil</span>
          </AnimatedButton>
        </div>
      </div>

      {/* Stats Grid */}
      <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4">
        {stats.map((s) => {
          const Icon = s.icon;
          return (
            <Link
              key={s.label}
              to={s.href}
              className="p-5 rounded-xl bg-surface border border-border hover:border-border-strong transition-all flex flex-col justify-between group"
            >
              <div className="flex items-start justify-between">
                <span className="text-xs font-medium text-text-muted">{s.label}</span>
                <div className={`p-2 rounded-lg bg-surface-2 border border-border ${s.color}`}>
                  <Icon className="w-4 h-4" />
                </div>
              </div>

              <div className="mt-4">
                <span className="text-3xl font-extrabold text-text">{s.value}</span>
                <p className="text-[11px] text-text-faint font-mono mt-1">{s.sub}</p>
              </div>

              <div className="mt-4 pt-3 border-t border-border/50 flex items-center justify-between text-xs text-text-faint group-hover:text-accent transition-colors">
                <span>Kelola</span>
                <ArrowRight className="w-3.5 h-3.5" />
              </div>
            </Link>
          );
        })}
      </div>

      {/* Recent Messages Preview */}
      <div className="rounded-xl bg-surface border border-border p-6">
        <div className="flex items-center justify-between mb-4">
          <div>
            <h2 className="text-base font-bold text-text">Pesan Masuk Terbaru</h2>
            <p className="text-xs text-text-faint">Pesan dari pengunjung formulir kontak</p>
          </div>
          <Link
            to="/admin/messages"
            className="text-xs font-mono text-accent hover:underline inline-flex items-center gap-1"
          >
            <span>Buka Inbox</span>
            <ArrowRight className="w-3 h-3" />
          </Link>
        </div>

        <div className="divide-y divide-border/60">
          {messages.slice(0, 5).map((m) => (
            <div key={m.id} className="py-3 flex items-center justify-between gap-4">
              <div className="min-w-0">
                <div className="flex items-center gap-2">
                  <span className="text-xs font-semibold text-text truncate">{m.name}</span>
                  <span className="text-[10px] font-mono text-text-faint">({m.email})</span>
                  {m.status === 'unread' && (
                    <span className="px-1.5 py-0.5 rounded text-[9px] font-semibold bg-rose-500/20 text-rose-400 border border-rose-500/30">
                      BARU
                    </span>
                  )}
                </div>
                <p className="text-xs text-text-muted truncate mt-0.5">{m.subject || m.body}</p>
              </div>
              <span className="text-[10px] text-text-faint font-mono whitespace-nowrap">
                {new Date(m.created_at).toLocaleDateString('id-ID', {
                  day: 'numeric',
                  month: 'short',
                })}
              </span>
            </div>
          ))}

          {messages.length === 0 && (
            <p className="py-6 text-center text-xs text-text-faint font-mono">
              Belum ada pesan masuk.
            </p>
          )}
        </div>
      </div>
    </div>
  );
};
