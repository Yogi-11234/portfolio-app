import React, { useState } from 'react';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { useForm } from 'react-hook-form';
import { toast } from 'sonner';
import {
  Plus,
  Edit2,
  Trash2,
  Upload,
  Loader2,
  X,
} from 'lucide-react';
import { api } from '@/lib/api';
import { Project } from '@/types';
import { AnimatedButton } from '@/components/AnimatedButton';
import { getMediaUrl } from '@/lib/utils';

interface ProjectFormData {
  title: string;
  slug: string;
  category_id: string;
  summary: string;
  description_md: string;
  demo_url?: string;
  repo_url?: string;
  featured: boolean;
  status: 'draft' | 'published' | 'archived';
  started_at?: string;
  finished_at?: string;
}

export const ProjectsAdminPage: React.FC = () => {
  const queryClient = useQueryClient();
  const [modalOpen, setModalOpen] = useState(false);
  const [editingProject, setEditingProject] = useState<Project | null>(null);
  const [thumbnailPath, setThumbnailPath] = useState<string | undefined>();
  const [selectedSkills, setSelectedSkills] = useState<string[]>([]);
  const [uploadingThumb, setUploadingThumb] = useState(false);

  // Queries
  const { data: projectsData, isLoading } = useQuery({
    queryKey: ['admin-projects'],
    queryFn: () => api.getProjectsAdmin({ per_page: 100 }),
  });

  const { data: categories = [] } = useQuery({
    queryKey: ['categories-project'],
    queryFn: () => api.getCategories('project'),
  });

  const { data: skills = [] } = useQuery({
    queryKey: ['admin-skills'],
    queryFn: api.getSkillsAdmin,
  });

  const projects = projectsData?.projects || [];

  const {
    register,
    handleSubmit,
    reset,
    setValue,
  } = useForm<ProjectFormData>();

  const openCreateModal = () => {
    setEditingProject(null);
    setThumbnailPath(undefined);
    setSelectedSkills([]);
    reset({
      title: '',
      slug: '',
      category_id: '',
      summary: '',
      description_md: '',
      demo_url: '',
      repo_url: '',
      featured: false,
      status: 'draft',
      started_at: '',
      finished_at: '',
    });
    setModalOpen(true);
  };

  const openEditModal = (proj: Project) => {
    setEditingProject(proj);
    setThumbnailPath(proj.thumbnail_path);
    setSelectedSkills(proj.skills.map((s) => s.id));
    reset({
      title: proj.title,
      slug: proj.slug,
      category_id: proj.category_id || '',
      summary: proj.summary,
      description_md: proj.description_md,
      demo_url: proj.demo_url || '',
      repo_url: proj.repo_url || '',
      featured: proj.featured,
      status: proj.status,
      started_at: proj.started_at ? proj.started_at.split('T')[0] : '',
      finished_at: proj.finished_at ? proj.finished_at.split('T')[0] : '',
    });
    setModalOpen(true);
  };

  // Mutations
  const saveMutation = useMutation({
    mutationFn: async (data: ProjectFormData) => {
      const payload: any = {
        title: data.title,
        slug: data.slug || undefined,
        category_id: data.category_id || undefined,
        summary: data.summary,
        description_md: data.description_md,
        thumbnail_path: thumbnailPath,
        demo_url: data.demo_url || undefined,
        repo_url: data.repo_url || undefined,
        featured: data.featured,
        is_featured: data.featured,
        status: data.status,
        started_at: data.started_at ? new Date(data.started_at).toISOString() : undefined,
        finished_at: data.finished_at ? new Date(data.finished_at).toISOString() : undefined,
        ended_at: data.finished_at ? new Date(data.finished_at).toISOString() : undefined,
        skill_ids: selectedSkills,
      };

      if (editingProject) {
        return api.updateProject(editingProject.id, payload);
      } else {
        return api.createProject(payload);
      }
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['admin-projects'] });
      queryClient.invalidateQueries({ queryKey: ['projects'] });
      toast.success(editingProject ? 'Proyek berhasil diperbarui' : 'Proyek berhasil dibuat');
      setModalOpen(false);
    },
    onError: (err: any) => {
      toast.error(err.message || 'Gagal menyimpan proyek');
    },
  });

  const deleteMutation = useMutation({
    mutationFn: (id: string) => api.deleteProject(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['admin-projects'] });
      toast.success('Proyek berhasil dihapus');
    },
    onError: (err: any) => {
      toast.error(err.message || 'Gagal menghapus proyek');
    },
  });

  const handleThumbnailUpload = async (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0];
    if (!file) return;

    setUploadingThumb(true);
    try {
      const res = await api.uploadMedia('project', file);
      setThumbnailPath(res.path);
      toast.success('Thumbnail berhasil diunggah!');
    } catch (err: any) {
      toast.error(err.message || 'Gagal mengunggah thumbnail');
    } finally {
      setUploadingThumb(false);
    }
  };

  const toggleSkill = (id: string) => {
    setSelectedSkills((prev) =>
      prev.includes(id) ? prev.filter((s) => s !== id) : [...prev, id]
    );
  };

  return (
    <div className="space-y-6">
      <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
        <div>
          <h1 className="text-2xl font-bold text-text">Manajemen Proyek</h1>
          <p className="text-xs text-text-muted mt-1">
            Kelola daftar karya portofolio, status publikasi, dan skill terkait
          </p>
        </div>

        <AnimatedButton
          onClick={openCreateModal}
          size="sm"
        >
          <Plus className="w-3.5 h-3.5 inline mr-1" />
          <span>Tambah Proyek Baru</span>
        </AnimatedButton>
      </div>

      {/* Projects Table */}
      <div className="rounded-xl bg-surface border border-border overflow-hidden">
        {isLoading ? (
          <div className="py-16 text-center text-xs text-text-faint font-mono">
            Memuat daftar proyek...
          </div>
        ) : projects.length > 0 ? (
          <div className="overflow-x-auto">
            <table className="w-full text-left text-xs">
              <thead className="bg-surface-2 border-b border-border text-text-muted font-mono uppercase tracking-wider text-[11px]">
                <tr>
                  <th className="py-3 px-4">Proyek</th>
                  <th className="py-3 px-4">Kategori</th>
                  <th className="py-3 px-4">Status</th>
                  <th className="py-3 px-4">Unggulan</th>
                  <th className="py-3 px-4">Skills</th>
                  <th className="py-3 px-4 text-right">Aksi</th>
                </tr>
              </thead>
              <tbody className="divide-y divide-border/60">
                {projects.map((proj) => (
                  <tr key={proj.id} className="hover:bg-surface-2/40 transition-colors">
                    <td className="py-3 px-4">
                      <div className="font-semibold text-text">{proj.title}</div>
                      <div className="text-[11px] text-text-faint font-mono">/{proj.slug}</div>
                    </td>
                    <td className="py-3 px-4 text-text-muted">
                      {proj.category?.name || '-'}
                    </td>
                    <td className="py-3 px-4">
                      <span
                        className={`px-2 py-0.5 rounded text-[10px] font-mono font-semibold ${
                          proj.status === 'published'
                            ? 'bg-emerald-950/80 text-emerald-400 border border-emerald-800'
                            : proj.status === 'draft'
                            ? 'bg-amber-950/80 text-amber-400 border border-amber-800'
                            : 'bg-zinc-800 text-text-faint border border-border'
                        }`}
                      >
                        {proj.status.toUpperCase()}
                      </span>
                    </td>
                    <td className="py-3 px-4">
                      {proj.featured ? (
                        <span className="text-amber-400 font-bold">★ Ya</span>
                      ) : (
                        <span className="text-text-faint">-</span>
                      )}
                    </td>
                    <td className="py-3 px-4 text-text-faint font-mono">
                      {proj.skills.length} skill
                    </td>
                    <td className="py-3 px-4 text-right space-x-1">
                      <button
                        onClick={() => openEditModal(proj)}
                        className="p-1.5 text-text-muted hover:text-text rounded hover:bg-surface-2"
                        title="Edit"
                      >
                        <Edit2 className="w-3.5 h-3.5" />
                      </button>
                      <button
                        onClick={() => {
                          if (confirm(`Yakin ingin menghapus proyek "${proj.title}"?`)) {
                            deleteMutation.mutate(proj.id);
                          }
                        }}
                        className="p-1.5 text-text-muted hover:text-danger rounded hover:bg-surface-2"
                        title="Hapus"
                      >
                        <Trash2 className="w-3.5 h-3.5" />
                      </button>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        ) : (
          <div className="py-16 text-center text-xs text-text-faint font-mono">
            Belum ada proyek terdaftar. Klik "Tambah Proyek Baru" untuk memulai.
          </div>
        )}
      </div>

      {/* Modal Form */}
      {modalOpen && (
        <div className="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/80 backdrop-blur-sm">
          <div className="relative w-full max-w-2xl max-h-[90vh] flex flex-col bg-surface border border-border rounded-xl shadow-2xl p-6 my-auto">
            {/* Modal Header */}
            <div className="flex items-center justify-between pb-3 border-b border-border flex-shrink-0">
              <h2 className="text-base font-bold text-text">
                {editingProject ? 'Edit Proyek' : 'Tambah Proyek Baru'}
              </h2>
              <button
                onClick={() => setModalOpen(false)}
                className="p-1.5 text-text-muted hover:text-text rounded hover:bg-surface-2"
              >
                <X className="w-4 h-4" />
              </button>
            </div>

            {/* Scrollable Form Area */}
            <form
              onSubmit={handleSubmit((data) => saveMutation.mutate(data))}
              className="flex-1 overflow-y-auto space-y-4 pr-1 mt-4 custom-scrollbar"
            >
              <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
                <div>
                  <label className="block text-xs font-medium text-text-muted mb-1">
                    Judul Proyek *
                  </label>
                  <input
                    type="text"
                    {...register('title', { required: true })}
                    onChange={(e) => {
                      setValue('title', e.target.value);
                      if (!editingProject) {
                        setValue(
                          'slug',
                          e.target.value
                            .toLowerCase()
                            .replace(/[^a-z0-9]+/g, '-')
                            .replace(/(^-|-$)+/g, '')
                        );
                      }
                    }}
                    className="w-full px-3 py-2 text-xs rounded-md bg-surface-2 border border-border text-text focus:outline-none focus:border-accent"
                  />
                </div>

                <div>
                  <label className="block text-xs font-medium text-text-muted mb-1">
                    Slug URL
                  </label>
                  <input
                    type="text"
                    {...register('slug')}
                    className="w-full px-3 py-2 text-xs rounded-md bg-surface-2 border border-border text-text font-mono focus:outline-none focus:border-accent"
                  />
                </div>
              </div>

              <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
                <div>
                  <label className="block text-xs font-medium text-text-muted mb-1">
                    Kategori Proyek
                  </label>
                  <select
                    {...register('category_id')}
                    className="w-full px-3 py-2 text-xs rounded-md bg-surface-2 border border-border text-text focus:outline-none focus:border-accent"
                  >
                    <option value="">Pilih Kategori...</option>
                    {categories.map((c) => (
                      <option key={c.id} value={c.id}>
                        {c.name}
                      </option>
                    ))}
                  </select>
                </div>

                <div>
                  <label className="block text-xs font-medium text-text-muted mb-1">
                    Status Publikasi
                  </label>
                  <select
                    {...register('status')}
                    className="w-full px-3 py-2 text-xs rounded-md bg-surface-2 border border-border text-text focus:outline-none focus:border-accent"
                  >
                    <option value="draft">Draft (Disimpan)</option>
                    <option value="published">Published (Tayang Publik)</option>
                    <option value="archived">Archived (Diarsipkan)</option>
                  </select>
                </div>
              </div>

              <div>
                <label className="block text-xs font-medium text-text-muted mb-1">
                  Ringkasan / Summary Singkat *
                </label>
                <textarea
                  rows={2}
                  {...register('summary', { required: true })}
                  placeholder="Deskripsi singkat proyek untuk kartu pratinjau..."
                  className="w-full px-3 py-2 text-xs rounded-md bg-surface-2 border border-border text-text focus:outline-none focus:border-accent"
                />
              </div>

              <div>
                <label className="block text-xs font-medium text-text-muted mb-1">
                  Deskripsi Lengkap (Markdown) *
                </label>
                <textarea
                  rows={6}
                  {...register('description_md', { required: true })}
                  placeholder="Detail arsitektur, tantangan teknis, dan fitur utama..."
                  className="w-full p-3 text-xs rounded-md bg-surface-2 border border-border text-text font-mono focus:outline-none focus:border-accent"
                />
              </div>

              {/* Thumbnail Upload */}
              <div>
                <label className="block text-xs font-medium text-text-muted mb-1">
                  Gambar Thumbnail
                </label>
                <div className="flex items-center gap-4">
                  {thumbnailPath && (
                    <div className="w-20 aspect-video rounded overflow-hidden border border-border bg-surface-2">
                      <img
                        src={getMediaUrl(thumbnailPath)}
                        alt="Thumb"
                        className="w-full h-full object-cover"
                      />
                    </div>
                  )}
                  <label className="inline-flex items-center gap-2 px-3 py-1.5 rounded-md bg-surface-2 border border-border text-xs text-text hover:border-border-strong cursor-pointer">
                    {uploadingThumb ? (
                      <Loader2 className="w-3.5 h-3.5 animate-spin" />
                    ) : (
                      <Upload className="w-3.5 h-3.5 text-text-muted" />
                    )}
                    <span>{uploadingThumb ? 'Mengunggah...' : 'Unggah Thumbnail'}</span>
                    <input
                      type="file"
                      accept="image/png,image/jpeg,image/webp"
                      onChange={handleThumbnailUpload}
                      disabled={uploadingThumb}
                      className="hidden"
                    />
                  </label>
                </div>
              </div>

              {/* Skills Multi-Select */}
              <div>
                <label className="block text-xs font-medium text-text-muted mb-2">
                  Skills Terkait ({selectedSkills.length} dipilih)
                </label>
                <div className="flex flex-wrap gap-1.5 max-h-32 overflow-y-auto p-2 rounded-md bg-surface-2 border border-border">
                  {skills.map((s) => {
                    const isSelected = selectedSkills.includes(s.id);
                    return (
                      <button
                        type="button"
                        key={s.id}
                        onClick={() => toggleSkill(s.id)}
                        className={`px-2.5 py-1 rounded text-[11px] font-mono transition-colors ${
                          isSelected
                            ? 'bg-accent text-accent-fg font-semibold'
                            : 'bg-surface border border-border text-text-muted hover:text-text'
                        }`}
                      >
                        {s.name}
                      </button>
                    );
                  })}
                </div>
              </div>

              <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
                <div>
                  <label className="block text-xs font-medium text-text-muted mb-1">
                    Demo URL
                  </label>
                  <input
                    type="url"
                    placeholder="https://..."
                    {...register('demo_url')}
                    className="w-full px-3 py-2 text-xs rounded-md bg-surface-2 border border-border text-text focus:outline-none focus:border-accent"
                  />
                </div>

                <div>
                  <label className="block text-xs font-medium text-text-muted mb-1">
                    GitHub Repo URL
                  </label>
                  <input
                    type="url"
                    placeholder="https://github.com/..."
                    {...register('repo_url')}
                    className="w-full px-3 py-2 text-xs rounded-md bg-surface-2 border border-border text-text focus:outline-none focus:border-accent"
                  />
                </div>
              </div>

              <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
                <div>
                  <label className="block text-xs font-medium text-text-muted mb-1">
                    Tanggal Mulai
                  </label>
                  <input
                    type="date"
                    {...register('started_at')}
                    className="w-full px-3 py-2 text-xs rounded-md bg-surface-2 border border-border text-text focus:outline-none focus:border-accent"
                  />
                </div>

                <div>
                  <label className="block text-xs font-medium text-text-muted mb-1">
                    Tanggal Selesai
                  </label>
                  <input
                    type="date"
                    {...register('finished_at')}
                    className="w-full px-3 py-2 text-xs rounded-md bg-surface-2 border border-border text-text focus:outline-none focus:border-accent"
                  />
                </div>
              </div>

              <div>
                <label className="flex items-center gap-2 cursor-pointer">
                  <input
                    type="checkbox"
                    {...register('featured')}
                    className="w-4 h-4 rounded border-border bg-surface-2 text-accent focus:ring-0"
                  />
                  <span className="text-xs font-medium text-text">
                    Tampilkan sebagai Proyek Unggulan di Halaman Depan
                  </span>
                </label>
              </div>

              {/* Modal Footer (Buttons) */}
              <div className="pt-3 border-t border-border flex justify-end items-center gap-3 flex-shrink-0">
                <button
                  type="button"
                  onClick={() => setModalOpen(false)}
                  className="px-4 py-2 rounded-md bg-surface-2 border border-border text-xs text-text hover:border-border-strong"
                >
                  Batal
                </button>
                <AnimatedButton
                  type="submit"
                  size="sm"
                  disabled={saveMutation.isPending}
                >
                  {saveMutation.isPending ? 'Menyimpan...' : 'Simpan Proyek'}
                </AnimatedButton>
              </div>
            </form>
          </div>
        </div>
      )}
    </div>
  );
};