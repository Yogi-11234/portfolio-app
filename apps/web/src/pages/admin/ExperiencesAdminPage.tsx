import React, { useState } from 'react';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { useForm } from 'react-hook-form';
import { toast } from 'sonner';
import { Plus, Edit2, Trash2, X } from 'lucide-react';
import { api } from '@/lib/api';
import { Experience, ExperienceType } from '@/types';
import { AnimatedButton } from '@/components/AnimatedButton';

interface ExperienceFormData {
  type: ExperienceType;
  title: string;
  organization: string;
  location?: string;
  started_at: string;
  ended_at?: string;
  current: boolean;
  description_md?: string;
  order_index: number;
  is_visible: boolean;
}

export const ExperiencesAdminPage: React.FC = () => {
  const queryClient = useQueryClient();
  const [modalOpen, setModalOpen] = useState(false);
  const [editingExp, setEditingExp] = useState<Experience | null>(null);

  const { data: experiences = [], isLoading } = useQuery({
    queryKey: ['admin-experiences'],
    queryFn: api.getExperiencesAdmin,
  });

  const { register, handleSubmit, reset, watch } = useForm<ExperienceFormData>();
  const isCurrent = watch('current');

  const openCreateModal = () => {
    setEditingExp(null);
    reset({
      type: 'work',
      title: '',
      organization: '',
      location: '',
      started_at: '',
      ended_at: '',
      current: false,
      description_md: '',
      order_index: 0,
      is_visible: true,
    });
    setModalOpen(true);
  };

  const openEditModal = (exp: Experience) => {
    setEditingExp(exp);
    const title = exp.title || (exp as any).role || '';
    const startedAt = exp.started_at || (exp as any).start_date || '';
    const endedAt = exp.ended_at || (exp as any).end_date || '';
    const current = exp.current ?? (exp as any).is_current ?? false;
    const orderIndex = exp.order_index ?? (exp as any).sort_order ?? 0;

    reset({
      type: exp.type,
      title: title,
      organization: exp.organization,
      location: exp.location || '',
      started_at: startedAt ? startedAt.split('T')[0] : '',
      ended_at: endedAt ? endedAt.split('T')[0] : '',
      current: current,
      description_md: exp.description_md || '',
      order_index: orderIndex,
      is_visible: exp.is_visible,
    });
    setModalOpen(true);
  };

  const saveMutation = useMutation({
    mutationFn: (data: ExperienceFormData) => {
      let startDateIso: string;
      try {
        startDateIso = new Date(data.started_at).toISOString();
      } catch {
        startDateIso = data.started_at;
      }

      let endDateIso: string | undefined = undefined;
      if (!data.current && data.ended_at) {
        try {
          endDateIso = new Date(data.ended_at).toISOString();
        } catch {
          endDateIso = data.ended_at;
        }
      }

      const payload = {
        type: data.type,
        role: data.title,
        title: data.title,
        organization: data.organization,
        location: data.location || undefined,
        start_date: startDateIso,
        started_at: startDateIso,
        end_date: endDateIso,
        ended_at: endDateIso,
        is_current: data.current,
        current: data.current,
        description_md: data.description_md || undefined,
        sort_order: Number(data.order_index) || 0,
        order_index: Number(data.order_index) || 0,
        is_visible: data.is_visible,
      };

      if (editingExp) {
        return api.updateExperience(editingExp.id, payload);
      } else {
        return api.createExperience(payload);
      }
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['admin-experiences'] });
      queryClient.invalidateQueries({ queryKey: ['experiences'] });
      toast.success(editingExp ? 'Pengalaman diperbarui' : 'Pengalaman ditambahkan');
      setModalOpen(false);
    },
    onError: (err: any) => {
      toast.error(err.message || 'Gagal menyimpan pengalaman');
    },
  });

  const deleteMutation = useMutation({
    mutationFn: (id: string) => api.deleteExperience(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['admin-experiences'] });
      queryClient.invalidateQueries({ queryKey: ['experiences'] });
      toast.success('Pengalaman dihapus');
    },
    onError: (err: any) => {
      toast.error(err.message || 'Gagal menghapus pengalaman');
    },
  });

  const formatDate = (iso: string) => {
    try {
      return new Date(iso).toLocaleDateString('id-ID', { month: 'short', year: 'numeric' });
    } catch {
      return iso;
    }
  };

  return (
    <div className="space-y-6">
      <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
        <div>
          <h1 className="text-2xl font-bold text-text">Manajemen Pengalaman & Edukasi</h1>
          <p className="text-xs text-text-muted mt-1">
            Kelola rekam jejak pekerjaan, pendidikan formal, organisasi, dan kerelawanan
          </p>
        </div>

        <AnimatedButton
          onClick={openCreateModal}
          size="sm"
        >
          <Plus className="w-3.5 h-3.5 inline mr-1" />
          <span>Tambah Pengalaman</span>
        </AnimatedButton>
      </div>

      <div className="rounded-xl bg-surface border border-border overflow-hidden">
        {isLoading ? (
          <div className="py-16 text-center text-xs text-text-faint font-mono">
            Memuat riwayat pengalaman...
          </div>
        ) : experiences.length > 0 ? (
          <div className="overflow-x-auto">
            <table className="w-full text-left text-xs">
              <thead className="bg-surface-2 border-b border-border text-text-muted font-mono uppercase tracking-wider text-[11px]">
                <tr>
                  <th className="py-3 px-4">Tipe</th>
                  <th className="py-3 px-4">Posisi / Gelar</th>
                  <th className="py-3 px-4">Organisasi / Instansi</th>
                  <th className="py-3 px-4">Periode</th>
                  <th className="py-3 px-4">Status</th>
                  <th className="py-3 px-4 text-right">Aksi</th>
                </tr>
              </thead>
              <tbody className="divide-y divide-border/60">
                {experiences.map((exp) => (
                  <tr key={exp.id} className="hover:bg-surface-2/40 transition-colors">
                    <td className="py-3 px-4">
                      <span className="px-2 py-0.5 rounded text-[10px] font-mono bg-surface-2 border border-border uppercase">
                        {exp.type}
                      </span>
                    </td>
                    <td className="py-3 px-4 font-semibold text-text">{exp.title || (exp as any).role}</td>
                    <td className="py-3 px-4 text-text-muted">
                      {exp.organization} {exp.location ? `(${exp.location})` : ''}
                    </td>
                    <td className="py-3 px-4 text-text-faint font-mono">
                      {formatDate(exp.started_at || (exp as any).start_date)} — {
                        (exp.current ?? (exp as any).is_current)
                          ? 'Sekarang'
                          : (exp.ended_at || (exp as any).end_date)
                            ? formatDate(exp.ended_at || (exp as any).end_date)
                            : '-'
                      }
                    </td>
                    <td className="py-3 px-4">
                      {exp.is_visible ? (
                        <span className="text-emerald-400 font-mono text-[11px]">Tayang</span>
                      ) : (
                        <span className="text-text-faint font-mono text-[11px]">Tersembunyi</span>
                      )}
                    </td>
                    <td className="py-3 px-4 text-right space-x-1">
                      <button
                        onClick={() => openEditModal(exp)}
                        className="p-1.5 text-text-muted hover:text-text rounded hover:bg-surface-2"
                      >
                        <Edit2 className="w-3.5 h-3.5" />
                      </button>
                      <button
                        onClick={() => {
                          if (confirm(`Hapus pengalaman "${exp.title}"?`)) {
                            deleteMutation.mutate(exp.id);
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
            Belum ada riwayat pengalaman.
          </div>
        )}
      </div>

      {/* Modal Form */}
      {modalOpen && (
        <div className="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/80 backdrop-blur-sm">
          <div className="relative w-full max-w-lg bg-surface border border-border rounded-xl shadow-2xl p-6 space-y-4">
            <div className="flex items-center justify-between pb-3 border-b border-border">
              <h2 className="text-sm font-bold text-text">
                {editingExp ? 'Edit Pengalaman' : 'Tambah Pengalaman Baru'}
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
              <div className="grid grid-cols-2 gap-4">
                <div>
                  <label className="block text-xs font-medium text-text-muted mb-1">
                    Tipe Kategori *
                  </label>
                  <select
                    {...register('type', { required: true })}
                    className="w-full px-3 py-2 text-xs rounded-md bg-surface-2 border border-border text-text focus:outline-none focus:border-accent"
                  >
                    <option value="work">Pekerjaan (Work)</option>
                    <option value="education">Pendidikan (Education)</option>
                    <option value="organization">Organisasi (Organization)</option>
                    <option value="volunteer">Sukarelawan (Volunteer)</option>
                  </select>
                </div>

                <div>
                  <label className="block text-xs font-medium text-text-muted mb-1">
                    Posisi / Gelar *
                  </label>
                  <input
                    type="text"
                    placeholder="cth. Senior Backend Engineer"
                    {...register('title', { required: true })}
                    className="w-full px-3 py-2 text-xs rounded-md bg-surface-2 border border-border text-text focus:outline-none focus:border-accent"
                  />
                </div>
              </div>

              <div className="grid grid-cols-2 gap-4">
                <div>
                  <label className="block text-xs font-medium text-text-muted mb-1">
                    Organisasi / Perusahaan *
                  </label>
                  <input
                    type="text"
                    placeholder="cth. PT Solusi Teknologi"
                    {...register('organization', { required: true })}
                    className="w-full px-3 py-2 text-xs rounded-md bg-surface-2 border border-border text-text focus:outline-none focus:border-accent"
                  />
                </div>

                <div>
                  <label className="block text-xs font-medium text-text-muted mb-1">
                    Lokasi (Opsional)
                  </label>
                  <input
                    type="text"
                    placeholder="cth. Jakarta / Remote"
                    {...register('location')}
                    className="w-full px-3 py-2 text-xs rounded-md bg-surface-2 border border-border text-text focus:outline-none focus:border-accent"
                  />
                </div>
              </div>

              <div className="grid grid-cols-2 gap-4">
                <div>
                  <label className="block text-xs font-medium text-text-muted mb-1">
                    Tanggal Mulai *
                  </label>
                  <input
                    type="date"
                    {...register('started_at', { required: true })}
                    className="w-full px-3 py-2 text-xs rounded-md bg-surface-2 border border-border text-text focus:outline-none focus:border-accent"
                  />
                </div>

                <div>
                  <label className="block text-xs font-medium text-text-muted mb-1">
                    Tanggal Selesai
                  </label>
                  <input
                    type="date"
                    disabled={isCurrent}
                    {...register('ended_at')}
                    className="w-full px-3 py-2 text-xs rounded-md bg-surface-2 border border-border text-text focus:outline-none focus:border-accent disabled:opacity-50"
                  />
                </div>
              </div>

              <div>
                <label className="flex items-center gap-2 cursor-pointer">
                  <input
                    type="checkbox"
                    {...register('current')}
                    className="w-4 h-4 rounded border-border bg-surface-2 text-accent focus:ring-0"
                  />
                  <span className="text-xs font-medium text-text">
                    Saat ini masih aktif di posisi ini (Masih berjalan)
                  </span>
                </label>
              </div>

              <div>
                <label className="block text-xs font-medium text-text-muted mb-1">
                  Deskripsi / Tanggung Jawab (Markdown)
                </label>
                <textarea
                  rows={4}
                  placeholder="Detail tanggung jawab, teknologi, pencapaian..."
                  {...register('description_md')}
                  className="w-full p-2.5 text-xs font-mono rounded-md bg-surface-2 border border-border text-text focus:outline-none focus:border-accent"
                />
              </div>

              <div>
                <label className="flex items-center gap-2 cursor-pointer">
                  <input
                    type="checkbox"
                    {...register('is_visible')}
                    className="w-4 h-4 rounded border-border bg-surface-2 text-accent focus:ring-0"
                  />
                  <span className="text-xs font-medium text-text">
                    Tampilkan di Halaman Publik
                  </span>
                </label>
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
                  {saveMutation.isPending ? 'Menyimpan...' : 'Simpan Pengalaman'}
                </AnimatedButton>
              </div>
            </form>
          </div>
        </div>
      )}
    </div>
  );
};
