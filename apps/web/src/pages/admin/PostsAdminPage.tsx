import React, { useState } from 'react';
import { Link } from 'react-router-dom';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { useForm } from 'react-hook-form';
import { toast } from 'sonner';
import { Plus, Edit2, Trash2, X, Upload, Loader2 } from 'lucide-react';
import { api } from '@/lib/api';
import { Post } from '@/types';
import { AnimatedButton } from '@/components/AnimatedButton';
import { getMediaUrl } from '@/lib/utils';

interface PostFormData {
  title: string;
  slug: string;
  category_id: string;
  excerpt: string;
  content_md: string;
  status: 'draft' | 'published' | 'archived';
  reading_minutes?: number;
}

export const PostsAdminPage: React.FC = () => {
  const queryClient = useQueryClient();
  const [modalOpen, setModalOpen] = useState(false);
  const [editingPost, setEditingPost] = useState<Post | null>(null);
  const [coverPath, setCoverPath] = useState<string | undefined>();
  const [uploadingCover, setUploadingCover] = useState(false);

  const { data: postsData, isLoading } = useQuery({
    queryKey: ['admin-posts'],
    queryFn: () => api.getPostsAdmin({ per_page: 100 }),
  });

  const { data: categories = [] } = useQuery({
    queryKey: ['categories-post'],
    queryFn: () => api.getCategories('post'),
  });

  const posts = postsData?.posts || [];

  const { register, handleSubmit, reset, setValue } = useForm<PostFormData>();

  const openCreateModal = () => {
    setEditingPost(null);
    setCoverPath(undefined);
    reset({
      title: '',
      slug: '',
      category_id: '',
      excerpt: '',
      content_md: '',
      status: 'draft',
      reading_minutes: 3,
    });
    setModalOpen(true);
  };

  const openEditModal = (post: Post) => {
    setEditingPost(post);
    setCoverPath(post.cover_path);
    reset({
      title: post.title,
      slug: post.slug,
      category_id: post.category_id || '',
      excerpt: post.excerpt,
      content_md: post.content_md,
      status: post.status,
      reading_minutes: post.reading_minutes,
    });
    setModalOpen(true);
  };

  const saveMutation = useMutation({
    mutationFn: (data: PostFormData) => {
      const payload = {
        title: data.title,
        slug: data.slug || undefined,
        category_id: data.category_id || undefined,
        excerpt: data.excerpt,
        content_md: data.content_md,
        cover_path: coverPath,
        status: data.status,
        reading_minutes: data.reading_minutes ? Number(data.reading_minutes) : undefined,
      };

      if (editingPost) {
        return api.updatePost(editingPost.id, payload);
      } else {
        return api.createPost(payload);
      }
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['admin-posts'] });
      queryClient.invalidateQueries({ queryKey: ['posts'] });
      toast.success(editingPost ? 'Artikel diperbarui' : 'Artikel diterbitkan/disimpan');
      setModalOpen(false);
    },
    onError: (err: any) => {
      toast.error(err.message || 'Gagal menyimpan artikel');
    },
  });

  const deleteMutation = useMutation({
    mutationFn: (id: string) => api.deletePost(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['admin-posts'] });
      queryClient.invalidateQueries({ queryKey: ['posts'] });
      toast.success('Artikel dihapus');
    },
    onError: (err: any) => {
      toast.error(err.message || 'Gagal menghapus artikel');
    },
  });

  const handleCoverUpload = async (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0];
    if (!file) return;

    setUploadingCover(true);
    try {
      const res = await api.uploadMedia('post', file);
      setCoverPath(res.path);
      toast.success('Cover gambar berhasil diunggah!');
    } catch (err: any) {
      toast.error(err.message || 'Gagal mengunggah cover');
    } finally {
      setUploadingCover(false);
    }
  };

  return (
    <div className="space-y-6">
      <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
        <div>
          <h1 className="text-2xl font-bold text-text">Manajemen Artikel & Blog</h1>
          <p className="text-xs text-text-muted mt-1">
            Kelola publikasi artikel teknis, status draft, kategori, dan estimasi waktu baca
          </p>
        </div>

        <AnimatedButton
          onClick={openCreateModal}
          size="sm"
        >
          <Plus className="w-3.5 h-3.5 inline mr-1" />
          <span>Tulis Artikel Baru</span>
        </AnimatedButton>
      </div>

      <div className="rounded-xl bg-surface border border-border overflow-hidden">
        {isLoading ? (
          <div className="py-16 text-center text-xs text-text-faint font-mono">
            Memuat daftar artikel...
          </div>
        ) : posts.length > 0 ? (
          <div className="overflow-x-auto">
            <table className="w-full text-left text-xs">
              <thead className="bg-surface-2 border-b border-border text-text-muted font-mono uppercase tracking-wider text-[11px]">
                <tr>
                  <th className="py-3 px-4">Judul Artikel</th>
                  <th className="py-3 px-4">Kategori</th>
                  <th className="py-3 px-4">Status</th>
                  <th className="py-3 px-4">Waktu Baca</th>
                  <th className="py-3 px-4">Tanggal Terbit</th>
                  <th className="py-3 px-4 text-right">Aksi</th>
                </tr>
              </thead>
              <tbody className="divide-y divide-border/60">
                {posts.map((post) => (
                  <tr key={post.id} className="hover:bg-surface-2/40 transition-colors">
                    <td className="py-3 px-4">
                      <div className="font-semibold text-text">{post.title}</div>
                      <div className="text-[11px] text-text-faint font-mono">/{post.slug}</div>
                    </td>
                    <td className="py-3 px-4 text-text-muted">
                      {post.category?.name || '-'}
                    </td>
                    <td className="py-3 px-4">
                      <span
                        className={`px-2 py-0.5 rounded text-[10px] font-mono font-semibold ${
                          post.status === 'published'
                            ? 'bg-emerald-950/80 text-emerald-400 border border-emerald-800'
                            : post.status === 'draft'
                            ? 'bg-amber-950/80 text-amber-400 border border-amber-800'
                            : 'bg-zinc-800 text-text-faint border border-border'
                        }`}
                      >
                        {post.status.toUpperCase()}
                      </span>
                    </td>
                    <td className="py-3 px-4 text-text-faint font-mono">
                      {post.reading_minutes} menit
                    </td>
                    <td className="py-3 px-4 text-text-faint font-mono">
                      {post.published_at
                        ? new Date(post.published_at).toLocaleDateString('id-ID', {
                            day: 'numeric',
                            month: 'short',
                            year: 'numeric',
                          })
                        : '-'}
                    </td>
                    <td className="py-3 px-4 text-right space-x-1">
                      <button
                        onClick={() => openEditModal(post)}
                        className="p-1.5 text-text-muted hover:text-text rounded hover:bg-surface-2"
                      >
                        <Edit2 className="w-3.5 h-3.5" />
                      </button>
                      <button
                        onClick={() => {
                          if (confirm(`Hapus artikel "${post.title}"?`)) {
                            deleteMutation.mutate(post.id);
                          }
                        }}
                        className="p-1.5 text-text-muted hover:text-danger rounded hover:bg-surface-2"
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
            Belum ada artikel yang ditulis.
          </div>
        )}
      </div>

      {/* Modal Form */}
      {modalOpen && (
        <div className="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/80 backdrop-blur-sm overflow-y-auto">
          <div className="relative w-full max-w-2xl my-8 bg-surface border border-border rounded-xl shadow-2xl p-6 space-y-4">
            <div className="flex items-center justify-between pb-3 border-b border-border">
              <h2 className="text-sm font-bold text-text">
                {editingPost ? 'Edit Artikel' : 'Tulis Artikel Baru'}
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
              <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
                <div>
                  <label className="block text-xs font-medium text-text-muted mb-1">
                    Judul Artikel *
                  </label>
                  <input
                    type="text"
                    {...register('title', { required: true })}
                    onChange={(e) => {
                      setValue('title', e.target.value);
                      if (!editingPost) {
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
                  <div className="flex items-center justify-between mb-1">
                    <label className="block text-xs font-medium text-text-muted">
                      Kategori Artikel
                    </label>
                    <Link
                      to="/admin/categories"
                      target="_blank"
                      className="text-[10px] text-accent hover:underline font-mono"
                    >
                      + Buka Menu Kategori
                    </Link>
                  </div>
                  <select
                    {...register('category_id')}
                    className="w-full px-3 py-2 text-xs rounded-md bg-surface-2 border border-border text-text focus:outline-none focus:border-accent"
                  >
                    <option value="">Tanpa Kategori (Opsional)</option>
                    {categories.map((c) => (
                      <option key={c.id} value={c.id}>
                        {c.name}
                      </option>
                    ))}
                  </select>
                  {categories.length === 0 && (
                    <p className="mt-1.5 text-[11px] text-amber-400/90 leading-tight">
                      ⚠️ Belum ada kategori artikel. Di menu Kategori, pastikan memilih tipe <strong>Kategori Artikel</strong> saat menambah/mengedit kategori agar muncul di sini.
                    </p>
                  )}
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
                  Cuplikan / Excerpt Singkat *
                </label>
                <textarea
                  rows={2}
                  {...register('excerpt', { required: true })}
                  placeholder="Ringkasan isi artikel untuk preview card..."
                  className="w-full px-3 py-2 text-xs rounded-md bg-surface-2 border border-border text-text focus:outline-none focus:border-accent"
                />
              </div>

              <div>
                <label className="block text-xs font-medium text-text-muted mb-1">
                  Konten Lengkap (Markdown) *
                </label>
                <textarea
                  rows={8}
                  {...register('content_md', { required: true })}
                  placeholder="Tuliskan artikel lengkap menggunakan sintaks markdown..."
                  className="w-full p-3 text-xs font-mono rounded-md bg-surface-2 border border-border text-text focus:outline-none focus:border-accent"
                />
              </div>

              <div>
                <label className="block text-xs font-medium text-text-muted mb-1">
                  Cover Gambar (Opsional)
                </label>
                <div className="flex items-center gap-4">
                  {coverPath && (
                    <div className="w-20 aspect-video rounded overflow-hidden border border-border bg-surface-2">
                      <img
                        src={getMediaUrl(coverPath)}
                        alt="Cover"
                        className="w-full h-full object-cover"
                      />
                    </div>
                  )}
                  <label className="inline-flex items-center gap-2 px-3 py-1.5 rounded-md bg-surface-2 border border-border text-xs text-text hover:border-border-strong cursor-pointer">
                    {uploadingCover ? (
                      <Loader2 className="w-3.5 h-3.5 animate-spin" />
                    ) : (
                      <Upload className="w-3.5 h-3.5 text-text-muted" />
                    )}
                    <span>{uploadingCover ? 'Mengunggah...' : 'Pilih Gambar Cover'}</span>
                    <input
                      type="file"
                      accept="image/png,image/jpeg,image/webp"
                      onChange={handleCoverUpload}
                      disabled={uploadingCover}
                      className="hidden"
                    />
                  </label>
                </div>
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
                  {saveMutation.isPending ? 'Menyimpan...' : 'Simpan Artikel'}
                </AnimatedButton>
              </div>
            </form>
          </div>
        </div>
      )}
    </div>
  );
};
