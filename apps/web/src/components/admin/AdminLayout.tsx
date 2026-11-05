import React from 'react';
import { Link, NavLink, Outlet, useNavigate } from 'react-router-dom';
import { useQuery, useQueryClient } from '@tanstack/react-query';
import { toast } from 'sonner';
import {
  LayoutDashboard,
  User,
  Layers,
  Code2,
  GraduationCap,
  Award,
  BookOpen,
  FolderTree,
  Inbox,
  LogOut,
  ExternalLink,
  Shield,
} from 'lucide-react';
import { api } from '@/lib/api';

const navItems = [
  { label: 'Ringkasan', href: '/admin', icon: LayoutDashboard, end: true },
  { label: 'Profil Saya', href: '/admin/profile', icon: User },
  { label: 'Proyek', href: '/admin/projects', icon: Layers },
  { label: 'Keahlian', href: '/admin/skills', icon: Code2 },
  { label: 'Pengalaman', href: '/admin/experiences', icon: GraduationCap },
  { label: 'Sertifikat', href: '/admin/certificates', icon: Award },
  { label: 'Artikel & Blog', href: '/admin/posts', icon: BookOpen },
  { label: 'Kategori', href: '/admin/categories', icon: FolderTree },
  { label: 'Pesan Masuk', href: '/admin/messages', icon: Inbox },
];

export const AdminLayout: React.FC = () => {
  const navigate = useNavigate();
  const queryClient = useQueryClient();

  const { data: user, isLoading, isError } = useQuery({
    queryKey: ['auth-me'],
    queryFn: api.me,
    retry: false,
  });

  const handleLogout = async () => {
    try {
      await api.logout();
      queryClient.clear();
      toast.success('Berhasil logout');
      navigate('/admin/login');
    } catch {
      toast.error('Gagal logout');
    }
  };

  if (isLoading) {
    return (
      <div className="min-h-screen flex items-center justify-center bg-bg font-mono text-xs text-text-faint">
        Memverifikasi sesi admin...
      </div>
    );
  }

  if (isError || !user) {
    navigate('/admin/login', { replace: true });
    return null;
  }

  return (
    <div className="min-h-screen flex bg-bg text-text">
      {/* Sidebar */}
      <aside className="w-64 border-r border-border bg-surface flex flex-col justify-between hidden md:flex fixed top-0 bottom-0 left-0 z-30">
        <div>
          {/* Brand */}
          <div className="p-6 border-b border-border flex items-center justify-between">
            <div className="flex items-center gap-2">
              <div className="w-8 h-8 rounded bg-surface-2 border border-border flex items-center justify-center text-accent">
                <Shield className="w-4 h-4" />
              </div>
              <span className="font-bold text-sm text-text">CMS Admin</span>
            </div>
            <Link
              to="/"
              target="_blank"
              title="Lihat Situs Publik"
              className="p-1.5 rounded text-text-muted hover:text-text hover:bg-surface-2 transition-colors"
            >
              <ExternalLink className="w-4 h-4" />
            </Link>
          </div>

          {/* Navigation Links */}
          <nav className="p-4 space-y-1">
            {navItems.map((item) => {
              const Icon = item.icon;
              return (
                <NavLink
                  key={item.href}
                  to={item.href}
                  end={item.end}
                  className={({ isActive }) =>
                    `flex items-center gap-3 px-3 py-2 text-xs rounded-md font-medium transition-colors ${
                      isActive
                        ? 'bg-accent text-accent-fg font-semibold shadow-sm'
                        : 'text-text-muted hover:text-text hover:bg-surface-2'
                    }`
                  }
                >
                  <Icon className="w-4 h-4" />
                  <span>{item.label}</span>
                </NavLink>
              );
            })}
          </nav>
        </div>

        {/* User Info & Logout */}
        <div className="p-4 border-t border-border bg-surface-2/40">
          <div className="flex items-center justify-between">
            <div className="min-w-0 pr-2">
              <p className="text-xs font-semibold text-text truncate">{user.name}</p>
              <p className="text-[10px] text-text-faint truncate font-mono">{user.email}</p>
            </div>
            <button
              onClick={handleLogout}
              title="Logout"
              className="p-2 rounded-md text-text-muted hover:text-danger hover:bg-surface transition-colors"
            >
              <LogOut className="w-4 h-4" />
            </button>
          </div>
        </div>
      </aside>

      {/* Main Content Area */}
      <div className="flex-1 md:pl-64 flex flex-col min-w-0">
        {/* Mobile Top Header */}
        <header className="md:hidden border-b border-border bg-surface px-4 py-3 flex items-center justify-between sticky top-0 z-30">
          <span className="font-bold text-sm text-text">CMS Admin</span>
          <div className="flex items-center gap-2">
            <Link to="/" className="p-1.5 text-xs text-text-muted hover:text-text">
              Publik
            </Link>
            <button
              onClick={handleLogout}
              className="p-1.5 text-xs text-danger font-medium"
            >
              Keluar
            </button>
          </div>
        </header>

        {/* Mobile Nav strip */}
        <div className="md:hidden flex overflow-x-auto gap-2 p-2 bg-surface-2 border-b border-border">
          {navItems.map((item) => (
            <NavLink
              key={item.href}
              to={item.href}
              end={item.end}
              className={({ isActive }) =>
                `px-3 py-1 text-xs rounded whitespace-nowrap font-medium ${
                  isActive ? 'bg-accent text-accent-fg' : 'text-text-muted bg-surface'
                }`
              }
            >
              {item.label}
            </NavLink>
          ))}
        </div>

        <main className="p-4 sm:p-8 flex-1 max-w-6xl w-full mx-auto">
          <Outlet />
        </main>
      </div>
    </div>
  );
};
