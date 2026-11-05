import React from 'react';
import { BrowserRouter, Routes, Route } from 'react-router-dom';
import { QueryClient, QueryClientProvider, useQuery } from '@tanstack/react-query';
import { Toaster } from 'sonner';
import { api } from '@/lib/api';
import { Navbar } from '@/components/Navbar';
import { Footer } from '@/components/Footer';
import { HomePage } from '@/pages/HomePage';
import { ProjectsPage } from '@/pages/ProjectsPage';
import { ProjectDetailPage } from '@/pages/ProjectDetailPage';
import { PostsPage } from '@/pages/PostsPage';
import { PostDetailPage } from '@/pages/PostDetailPage';
import { NotFoundPage } from '@/pages/NotFoundPage';

// Admin CMS
import { LoginPage } from '@/pages/admin/LoginPage';
import { AdminLayout } from '@/components/admin/AdminLayout';
import { DashboardPage } from '@/pages/admin/DashboardPage';
import { ProfileAdminPage } from '@/pages/admin/ProfileAdminPage';
import { ProjectsAdminPage } from '@/pages/admin/ProjectsAdminPage';
import { SkillsAdminPage } from '@/pages/admin/SkillsAdminPage';
import { ExperiencesAdminPage } from '@/pages/admin/ExperiencesAdminPage';
import { CertificatesAdminPage } from '@/pages/admin/CertificatesAdminPage';
import { PostsAdminPage } from '@/pages/admin/PostsAdminPage';
import { CategoriesAdminPage } from '@/pages/admin/CategoriesAdminPage';
import { MessagesAdminPage } from '@/pages/admin/MessagesAdminPage';

const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      staleTime: 1000 * 60 * 3, // 3 minutes
      refetchOnWindowFocus: false,
      retry: 1,
    },
  },
});

const PublicLayout: React.FC = () => {
  const { data: profile } = useQuery({
    queryKey: ['profile'],
    queryFn: api.getProfile,
  });

  return (
    <div className="min-h-screen flex flex-col bg-bg text-text antialiased relative">
      <Navbar profile={profile} />
      
      {/* PERBAIKAN: Beri relative z-10 agar konten halaman muncul di atas WaveMesh */}
      <main className="flex-1 relative z-10">
        <Routes>
          <Route path="/" element={<HomePage />} />
          <Route path="/projects" element={<ProjectsPage />} />
          <Route path="/projects/:slug" element={<ProjectDetailPage />} />
          <Route path="/posts" element={<PostsPage />} />
          <Route path="/posts/:slug" element={<PostDetailPage />} />
          <Route path="*" element={<NotFoundPage />} />
        </Routes>
      </main>

      {/* Footer dipanggil di sini dan otomatis ada di semua halaman publik */}
      <Footer profile={profile} />
    </div>
  );
};

export default function App() {
  return (
    <QueryClientProvider client={queryClient}>
      <BrowserRouter>
        <Routes>
          {/* Admin Login Route */}
          <Route path="/admin/login" element={<LoginPage />} />

          {/* Admin CMS Protected Routes */}
          <Route path="/admin" element={<AdminLayout />}>
            <Route index element={<DashboardPage />} />
            <Route path="profile" element={<ProfileAdminPage />} />
            <Route path="projects" element={<ProjectsAdminPage />} />
            <Route path="skills" element={<SkillsAdminPage />} />
            <Route path="experiences" element={<ExperiencesAdminPage />} />
            <Route path="certificates" element={<CertificatesAdminPage />} />
            <Route path="posts" element={<PostsAdminPage />} />
            <Route path="categories" element={<CategoriesAdminPage />} />
            <Route path="messages" element={<MessagesAdminPage />} />
          </Route>

          {/* Public Website Routes */}
          <Route path="/*" element={<PublicLayout />} />
        </Routes>
        <Toaster position="top-right" richColors theme="dark" />
      </BrowserRouter>
    </QueryClientProvider>
  );
}