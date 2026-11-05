import React, { useState } from 'react';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { useForm } from 'react-hook-form';
import { toast } from 'sonner';
import { Plus, Edit2, Trash2, X, Upload, Loader2 } from 'lucide-react';
import { api } from '@/lib/api';
import { Skill, SkillGroup } from '@/types';
import { AnimatedButton } from '@/components/AnimatedButton';
import { TechIcon, TECH_PRESETS } from '@/components/TechIcon';

interface SkillFormData {
  name: string;
  slug: string;
  group_name: SkillGroup;
  icon_key: string;
  brand_color: string;
  order_index: number;
  is_visible: boolean;
}

const GROUPS: { key: SkillGroup; label: string }[] = [
  { key: 'language', label: 'Bahasa Pemrograman' },
  { key: 'framework', label: 'Framework & Library' },
  { key: 'database', label: 'Database & Cache' },
  { key: 'cloud', label: 'Cloud & Infrastruktur' },
  { key: 'tool', label: 'Alat & DevOps' },
  { key: 'other', label: 'Lainnya' },
];

export const SkillsAdminPage: React.FC = () => {
  const queryClient = useQueryClient();
  const [modalOpen, setModalOpen] = useState(false);
  const [editingSkill, setEditingSkill] = useState<Skill | null>(null);
  const [selectedIcon, setSelectedIcon] = useState<string>('react');
  const [uploadingIcon, setUploadingIcon] = useState(false);

  const { data: skills = [], isLoading } = useQuery({
    queryKey: ['admin-skills'],
    queryFn: api.getSkillsAdmin,
  });

  const { register, handleSubmit, reset, setValue, watch } = useForm<SkillFormData>();

  const watchedName = watch('name') || 'Nama Skill';
  const watchedColor = watch('brand_color') || '#39FF14';

  const openCreateModal = () => {
    setEditingSkill(null);
    setSelectedIcon('react');
    reset({
      name: '',
      slug: '',
      group_name: 'language',
      icon_key: 'react',
      brand_color: '#39FF14',
      order_index: 0,
      is_visible: true,
    });
    setModalOpen(true);
  };

  const openEditModal = (skill: Skill) => {
    setEditingSkill(skill);
    const icon = skill.icon_key || (skill as any).icon_name || 'react';
    setSelectedIcon(icon);
    reset({
      name: skill.name,
      slug: skill.slug,
      group_name: skill.group_name,
      icon_key: icon,
      brand_color: skill.brand_color || '#39FF14',
      order_index: skill.order_index,
      is_visible: skill.is_visible,
    });
    setModalOpen(true);
  };

  const handleIconUpload = async (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0];
    if (!file) return;

    setUploadingIcon(true);
    try {
      const res = await api.uploadMedia('skill', file);
      setSelectedIcon(res.path);
      setValue('icon_key', res.path);
      toast.success('Ikon berhasil diunggah!');
    } catch (err: any) {
      toast.error(err.message || 'Gagal mengunggah ikon');
    } finally {
      setUploadingIcon(false);
    }
  };

  const saveMutation = useMutation({
    mutationFn: (data: SkillFormData) => {
      const payload = {
        name: data.name,
        slug: data.slug || undefined,
        group_name: data.group_name,
        icon_key: selectedIcon || data.icon_key || undefined,
        brand_color: data.brand_color,
        level: 5,
        order_index: Number(data.order_index) || 0,
        sort_order: Number(data.order_index) || 0,
        is_visible: data.is_visible,
      };

      if (editingSkill) {
        return api.updateSkill(editingSkill.id, payload);
      } else {
        return api.createSkill(payload);
      }
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['admin-skills'] });
      queryClient.invalidateQueries({ queryKey: ['skills'] });
      toast.success(editingSkill ? 'Skill diperbarui' : 'Skill ditambahkan');
      setModalOpen(false);
    },
    onError: (err: any) => {
      toast.error(err.message || 'Gagal menyimpan skill');
    },
  });

  const deleteMutation = useMutation({
    mutationFn: (id: string) => api.deleteSkill(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['admin-skills'] });
      queryClient.invalidateQueries({ queryKey: ['skills'] });
      toast.success('Skill dihapus');
    },
    onError: (err: any) => {
      toast.error(err.message || 'Gagal menghapus skill');
    },
  });

  return (
    <div className="space-y-6">
      <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
        <div>
          <h1 className="text-2xl font-bold text-text">Manajemen Keahlian (Skills)</h1>
          <p className="text-xs text-text-muted mt-1">
            Kelola daftar teknologi, logo / ikon, kategori grup, dan tingkat kemahiran
          </p>
        </div>

        <AnimatedButton onClick={openCreateModal} size="sm">
          <Plus className="w-3.5 h-3.5 inline mr-1" />
          <span>Tambah Skill Baru</span>
        </AnimatedButton>
      </div>

      {isLoading ? (
        <div className="py-16 text-center text-xs text-text-faint font-mono">
          Memuat daftar skill...
        </div>
      ) : (
        <div className="space-y-6">
          {GROUPS.map((g) => {
            const groupSkills = skills.filter((s) => s.group_name === g.key);
            if (groupSkills.length === 0) return null;

            return (
              <div key={g.key} className="rounded-xl bg-surface border border-border p-5">
                <div className="flex items-center gap-2 mb-3 pb-2 border-b border-border/60">
                  <h3 className="text-xs font-mono uppercase tracking-wider font-semibold text-text">
                    {g.label} ({groupSkills.length})
                  </h3>
                </div>

                <div className="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 gap-3">
                  {groupSkills.map((s) => {
                    const icon = s.icon_key || (s as any).icon_name;
                    return (
                      <div
                        key={s.id}
                        className="p-3 rounded-xl bg-surface-2 border border-border flex items-center justify-between gap-3 hover:border-accent/40 transition-colors"
                      >
                        <div className="flex items-center gap-3 min-w-0">
                          <div
                            className="w-10 h-10 rounded-lg bg-surface border border-border/80 flex items-center justify-center flex-shrink-0"
                            style={{
                              boxShadow: '3px 3px 6px rgba(0,0,0,0.4), -2px -2px 6px rgba(255,255,255,0.02)',
                            }}
                          >
                            <TechIcon
                              name={s.name}
                              iconKey={icon}
                              brandColor={s.brand_color}
                              className="w-5 h-5"
                            />
                          </div>
                          <div className="min-w-0">
                            <p className="text-xs font-semibold text-text truncate">{s.name}</p>
                            <p className="text-[10px] text-text-faint font-mono">
                              {s.group_name} {!s.is_visible && '• (Tersembunyi)'}
                            </p>
                          </div>
                        </div>

                        <div className="flex items-center gap-1 flex-shrink-0">
                          <button
                            onClick={() => openEditModal(s)}
                            className="p-1.5 text-text-muted hover:text-text rounded hover:bg-surface transition-colors"
                            title="Edit"
                          >
                            <Edit2 className="w-3.5 h-3.5" />
                          </button>
                          <button
                            onClick={() => {
                              if (confirm(`Hapus skill "${s.name}"?`)) {
                                deleteMutation.mutate(s.id);
                              }
                            }}
                            className="p-1.5 text-text-muted hover:text-danger rounded hover:bg-surface transition-colors"
                            title="Hapus"
                          >
                            <Trash2 className="w-3.5 h-3.5" />
                          </button>
                        </div>
                      </div>
                    );
                  })}
                </div>
              </div>
            );
          })}
        </div>
      )}

      {/* Modal Form Tambah / Edit Skill */}
      {modalOpen && (
        <div className="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/80 backdrop-blur-sm">
          <div className="relative w-full max-w-xl max-h-[90vh] flex flex-col bg-surface border border-border rounded-2xl shadow-2xl p-6 my-auto">
            {/* Header Modal (Fixed Top) */}
            <div className="flex items-center justify-between pb-3 border-b border-border flex-shrink-0">
              <h2 className="text-sm font-bold text-text">
                {editingSkill ? 'Edit Keahlian' : 'Tambah Keahlian Baru'}
              </h2>
              <button
                onClick={() => setModalOpen(false)}
                className="p-1.5 text-text-muted hover:text-text rounded hover:bg-surface-2"
              >
                <X className="w-4 h-4" />
              </button>
            </div>

            {/* Form Scrollable Area */}
            <form
              onSubmit={handleSubmit((data) => saveMutation.mutate(data))}
              className="space-y-4 overflow-y-auto pr-1 mt-4 flex-1 custom-scrollbar"
            >
              <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
                <div>
                  <label className="block text-xs font-medium text-text-muted mb-1">
                    Nama Skill *
                  </label>
                  <input
                    type="text"
                    placeholder="cth. Go (Golang)"
                    {...register('name', { required: true })}
                    onChange={(e) => {
                      setValue('name', e.target.value);
                      if (!editingSkill) {
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

              <div>
                <label className="block text-xs font-medium text-text-muted mb-1">
                  Grup Kategori *
                </label>
                <select
                  {...register('group_name', { required: true })}
                  className="w-full px-3 py-2 text-xs rounded-md bg-surface-2 border border-border text-text focus:outline-none focus:border-accent"
                >
                  {GROUPS.map((g) => (
                    <option key={g.key} value={g.key}>
                      {g.label}
                    </option>
                  ))}
                </select>
              </div>

              {/* ----------------- PILIH ICON & UPLOAD ----------------- */}
              <div className="p-3.5 rounded-xl bg-surface-2/60 border border-border space-y-3">
                <div className="flex items-center justify-between">
                  <label className="block text-xs font-bold text-text">
                    Ikon / Logo Teknologi
                  </label>
                  <span className="text-[10px] font-mono text-accent">
                    {selectedIcon ? selectedIcon : 'Belum dipilih'}
                  </span>
                </div>

                {/* Preset Icons Selection */}
                <div>
                  <p className="text-[11px] text-text-muted mb-1.5">
                    Pilih Dari Ikon Populer:
                  </p>
                  <div className="grid grid-cols-4 sm:grid-cols-8 gap-2 max-h-36 overflow-y-auto p-1.5 rounded-lg bg-surface border border-border/80">
                    {TECH_PRESETS.map((p) => {
                      const isSelected = selectedIcon === p.key;
                      return (
                        <button
                          key={p.key}
                          type="button"
                          onClick={() => {
                            setSelectedIcon(p.key);
                            setValue('icon_key', p.key);
                            setValue('brand_color', p.color);
                          }}
                          className={`p-2 rounded-lg flex flex-col items-center justify-center gap-1 transition-all ${
                            isSelected
                              ? 'bg-accent/15 border border-accent shadow-[0_0_10px_rgba(57,255,20,0.3)]'
                              : 'bg-surface-2 border border-border/60 hover:border-border-strong hover:bg-surface-2/80'
                          }`}
                          title={p.label}
                        >
                          <TechIcon name={p.key} iconKey={p.key} brandColor={p.color} className="w-5 h-5" />
                          <span className="text-[9px] font-mono truncate w-full text-center text-text-muted">
                            {p.key}
                          </span>
                        </button>
                      );
                    })}
                  </div>
                </div>

                {/* Upload Custom Icon */}
                <div className="pt-2 border-t border-border/60 flex items-center justify-between">
                  <div className="text-[11px] text-text-muted">
                    Atau Unggah Ikon Sendiri (.png, .webp, .jpg):
                  </div>
                  <label className="inline-flex items-center gap-1.5 px-3 py-1.5 rounded-md bg-surface border border-border text-xs text-text hover:border-accent cursor-pointer transition-colors">
                    {uploadingIcon ? (
                      <Loader2 className="w-3.5 h-3.5 animate-spin" />
                    ) : (
                      <Upload className="w-3.5 h-3.5 text-accent" />
                    )}
                    <span>{uploadingIcon ? 'Mengunggah...' : 'Pilih Berkas'}</span>
                    <input
                      type="file"
                      accept="image/png,image/jpeg,image/webp"
                      onChange={handleIconUpload}
                      disabled={uploadingIcon}
                      className="hidden"
                    />
                  </label>
                </div>
              </div>

              {/* Brand Color */}
              <div>
                <label className="block text-xs font-medium text-text-muted mb-1">
                  Warna Brand (#Hex)
                </label>
                <div className="flex items-center gap-2">
                  <input
                    type="color"
                    {...register('brand_color')}
                    className="w-8 h-8 rounded border border-border bg-surface-2 cursor-pointer"
                  />
                  <input
                    type="text"
                    {...register('brand_color')}
                    placeholder="#39FF14"
                    className="flex-1 px-2.5 py-1.5 text-xs font-mono rounded-md bg-surface-2 border border-border text-text focus:outline-none focus:border-accent"
                  />
                </div>
              </div>

              {/* Pratinjau Card Neumorphic Langsung */}
              <div className="pt-2 border-t border-border/80">
                <span className="text-[11px] font-mono text-text-muted block mb-2">
                  Pratinjau Kartu Portofolio:
                </span>
                <div className="flex justify-center p-3 rounded-xl bg-black/40 border border-border/60">
                  <div className="neumorphic-card flex flex-col items-center justify-center p-3 rounded-2xl w-28 h-28">
                    <div
                      className="w-12 h-12 rounded-xl flex items-center justify-center bg-gradient-to-br from-[#1c1c26] to-[#101016] border border-white/5 shadow-inner"
                      style={{ borderColor: `${watchedColor}40` }}
                    >
                      <TechIcon
                        name={watchedName}
                        iconKey={selectedIcon}
                        brandColor={watchedColor}
                        className="w-6 h-6"
                      />
                    </div>
                    <h4 className="mt-2 text-xs font-semibold text-text truncate max-w-full px-1">
                      {watchedName}
                    </h4>
                  </div>
                </div>
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

              {/* Footer Modal (Buttons) */}
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
                  {saveMutation.isPending ? 'Menyimpan...' : 'Simpan Skill'}
                </AnimatedButton>
              </div>
            </form>
          </div>
        </div>
      )}
    </div>
  );
};

export default SkillsAdminPage;