import React, { useState } from 'react';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { useForm } from 'react-hook-form';
import { toast } from 'sonner';
import { Plus, Edit2, Trash2, X, FolderTree } from 'lucide-react';
import { api } from '@/lib/api';
import { Category } from '@/types';
import { AnimatedButton } from '@/components/AnimatedButton';

interface CategoryFormData {
  name: string;
  slug: string;
  type: 'project' | 'post';
}

export const CategoriesAdminPage: React.FC = () => {
  const queryClient = useQueryClient();
  const [modalOpen, setModalOpen] = useState(false);
  const [editingCat, setEditingCat] = useState<Category | null>(null);

  const { data: projectCategories = [], isLoading: loadingProj } = useQuery({
    queryKey: ['categories-project'],
    queryFn: () => api.getCategories('project'),
  });

  const { data: postCategories = [], isLoading: loadingPost } = useQuery({
    queryKey: ['categories-post'],
    queryFn: () => api.getCategories('post'),
  });

  const { register, handleSubmit, reset, setValue, watch } = useForm<CategoryFormData>();
  const currentType = watch('type');

  const openCreateModal = (defaultType: 'project' | 'post' = 'project') => {
    setEditingCat(null);
    reset({
      name: '',
      slug: '',
      type: defaultType,
    });
    setModalOpen(true);
  };

  const openEditModal = (cat: Category) => {
    setEditingCat(cat);
    reset({
      name: cat.name,
      slug: cat.slug,
      type: cat.type,
    });
    setModalOpen(true);
  };

  const saveMutation = useMutation({
    mutationFn: (data: CategoryFormData) => {
      const payload = {
        name: data.name,
        slug: data.slug || undefined,
        type: data.type,
      };

      if (editingCat) {
        return api.updateCategory(editingCat.id, payload);
      } else {
        return api.createCategory(payload);
      }
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['categories-project'] });
      queryClient.invalidateQueries({ queryKey: ['categories-post'] });
      toast.success(editingCat ? 'Kategori diperbarui' : 'Kategori dibuat');
      setModalOpen(false);
    },
    onError: (err: any) => {
      toast.error(err.message || 'Gagal menyimpan kategori');
    },
  });

  const deleteMutation = useMutation({
    mutationFn: (id: string) => api.deleteCategory(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['categories-project'] });
      queryClient.invalidateQueries({ queryKey: ['categories-post'] });
      toast.success('Kategori dihapus');
    },
    onError: (err: any) => {
      toast.error(err.message || 'Gagal menghapus kategori');
    },
  });

  return (
    <div className="space-y-6">
      <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
        <div>
          <h1 className="text-2xl font-bold text-text">Manajemen Kategori</h1>
          <p className="text-xs text-text-muted mt-1">
            Kelola kategori taksonomi untuk pengelompokan Proyek dan Artikel Blog
          </p>
        </div>

        <AnimatedButton
          onClick={() => openCreateModal('project')}
          size="sm"
        >
          <Plus className="w-3.5 h-3.5 inline mr-1" />
          <span>Tambah Kategori Baru</span>
        </AnimatedButton>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
        {/* Project Categories */}
        <div className="rounded-xl bg-surface border border-border p-5 space-y-4">
          <div className="flex items-center justify-between pb-2 border-b border-border">
            <div className="flex items-center gap-2">
              <FolderTree className="w-4 h-4 text-sky-400" />
              <h2 className="text-sm font-bold text-text">Kategori Proyek ({projectCategories.length})</h2>
            </div>
            <button
              onClick={() => openCreateModal('project')}
              className="px-2.5 py-1 text-xs rounded-lg bg-sky-500/10 text-sky-400 hover:bg-sky-500/20 border border-sky-500/30 font-mono flex items-center gap-1 transition-colors"
            >
              <Plus className="w-3.5 h-3.5" />
              <span>Tambah Proyek</span>
            </button>
          </div>

          <div className="divide-y divide-border/60">
            {projectCategories.map((c) => (
              <div key={c.id} className="py-2.5 flex items-center justify-between">
                <div>
                  <div className="flex items-center gap-2">
                    <p className="text-xs font-semibold text-text">{c.name}</p>
                    <span className="px-1.5 py-0.2 rounded text-[9px] font-mono bg-sky-500/15 text-sky-400 border border-sky-500/30">
                      project
                    </span>
                  </div>
                  <p className="text-[10px] text-text-faint font-mono">/{c.slug}</p>
                </div>
                <div className="flex items-center gap-1">
                  <button
                    onClick={() => openEditModal(c)}
                    className="p-1 text-text-muted hover:text-text rounded hover:bg-surface-2"
                  >
                    <Edit2 className="w-3.5 h-3.5" />
                  </button>
                  <button
                    onClick={() => {
                      if (confirm(`Hapus kategori "${c.name}"?`)) {
                        deleteMutation.mutate(c.id);
                      }
                    }}
                    className="p-1 text-text-muted hover:text-danger rounded hover:bg-surface-2"
                  >
                    <Trash2 className="w-3.5 h-3.5" />
                  </button>
                </div>
              </div>
            ))}
            {projectCategories.length === 0 && !loadingProj && (
              <p className="text-xs text-text-faint py-4 text-center font-mono">
                Belum ada kategori proyek.
              </p>
            )}
          </div>
        </div>

        {/* Post Categories */}
        <div className="rounded-xl bg-surface border border-border p-5 space-y-4">
          <div className="flex items-center justify-between pb-2 border-b border-border">
            <div className="flex items-center gap-2">
              <FolderTree className="w-4 h-4 text-amber-400" />
              <h2 className="text-sm font-bold text-text">Kategori Artikel Blog ({postCategories.length})</h2>
            </div>
            <button
              onClick={() => openCreateModal('post')}
              className="px-2.5 py-1 text-xs rounded-lg bg-amber-500/10 text-amber-400 hover:bg-amber-500/20 border border-amber-500/30 font-mono flex items-center gap-1 transition-colors"
            >
              <Plus className="w-3.5 h-3.5" />
              <span>Tambah Artikel</span>
            </button>
          </div>

          <div className="divide-y divide-border/60">
            {postCategories.map((c) => (
              <div key={c.id} className="py-2.5 flex items-center justify-between">
                <div>
                  <div className="flex items-center gap-2">
                    <p className="text-xs font-semibold text-text">{c.name}</p>
                    <span className="px-1.5 py-0.2 rounded text-[9px] font-mono bg-amber-500/15 text-amber-400 border border-amber-500/30">
                      post
                    </span>
                  </div>
                  <p className="text-[10px] text-text-faint font-mono">/{c.slug}</p>
                </div>
                <div className="flex items-center gap-1">
                  <button
                    onClick={() => openEditModal(c)}
                    className="p-1 text-text-muted hover:text-text rounded hover:bg-surface-2"
                  >
                    <Edit2 className="w-3.5 h-3.5" />
                  </button>
                  <button
                    onClick={() => {
                      if (confirm(`Hapus kategori "${c.name}"?`)) {
                        deleteMutation.mutate(c.id);
                      }
                    }}
                    className="p-1 text-text-muted hover:text-danger rounded hover:bg-surface-2"
                  >
                    <Trash2 className="w-3.5 h-3.5" />
                  </button>
                </div>
              </div>
            ))}
            {postCategories.length === 0 && !loadingPost && (
              <p className="text-xs text-text-faint py-4 text-center font-mono">
                Belum ada kategori artikel.
              </p>
            )}
          </div>
        </div>
      </div>

      {/* Modal Form */}
      {modalOpen && (
        <div className="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/80 backdrop-blur-sm">
          <div className="relative w-full max-w-sm bg-surface border border-border rounded-xl shadow-2xl p-6 space-y-4">
            <div className="flex items-center justify-between pb-3 border-b border-border">
              <h2 className="text-sm font-bold text-text">
                {editingCat ? 'Edit Kategori' : 'Tambah Kategori'}
              </h2>
              <button
                onClick={() => setModalOpen(false)}
                className="p-1.5 text-text-muted hover:text-text rounded hover:bg-surface-2"
              >
                <X className="w-4 h-4" />
              </button>
            </div>

            <form
              onSubmit={handleSubmit((data) => saveMutation.mutate(data))}
              className="space-y-4"
            >
              <div>
                <label className="block text-xs font-medium text-text-muted mb-1.5">
                  Tipe Kategori *
                </label>
                <div className="grid grid-cols-2 gap-2">
                  <button
                    type="button"
                    onClick={() => setValue('type', 'project')}
                    className={`py-2 px-3 rounded-lg border text-xs font-mono flex items-center justify-center gap-1.5 transition-all ${
                      currentType === 'project'
                        ? 'border-sky-500 bg-sky-500/20 text-sky-400 font-bold shadow-[0_0_10px_rgba(56,189,248,0.2)]'
                        : 'border-border bg-surface-2 text-text-muted hover:border-text-faint'
                    }`}
                  >
                    <span>📁 Kategori Proyek</span>
                  </button>
                  <button
                    type="button"
                    onClick={() => setValue('type', 'post')}
                    className={`py-2 px-3 rounded-lg border text-xs font-mono flex items-center justify-center gap-1.5 transition-all ${
                      currentType === 'post'
                        ? 'border-amber-500 bg-amber-500/20 text-amber-400 font-bold shadow-[0_0_10px_rgba(245,158,11,0.2)]'
                        : 'border-border bg-surface-2 text-text-muted hover:border-text-faint'
                    }`}
                  >
                    <span>📝 Kategori Artikel</span>
                  </button>
                </div>
                <input type="hidden" {...register('type', { required: true })} />
              </div>

              <div>
                <label className="block text-xs font-medium text-text-muted mb-1">
                  Nama Kategori *
                </label>
                <input
                  type="text"
                  placeholder="cth. Fullstack Web"
                  {...register('name', { required: true })}
                  onChange={(e) => {
                    setValue('name', e.target.value);
                    if (!editingCat) {
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

              <div className="pt-3 border-t border-border flex justify-end items-center gap-3">
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
                  {saveMutation.isPending ? 'Menyimpan...' : 'Simpan Kategori'}
                </AnimatedButton>
              </div>
            </form>
          </div>
        </div>
      )}
    </div>
  );
};
