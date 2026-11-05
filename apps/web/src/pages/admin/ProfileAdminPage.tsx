import React, { useState, useEffect } from 'react';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { useForm } from 'react-hook-form';
import { toast } from 'sonner';
import {
  Upload,
  FileText,
  Loader2,
  Save,
  CheckCircle,
  Eye,
  Edit3,
} from 'lucide-react';
import { api } from '@/lib/api';
import { AnimatedButton } from '@/components/AnimatedButton';
import { getMediaUrl } from '@/lib/utils';
import { SanitizedHTML } from '@/components/SanitizedHTML';

interface ProfileFormData {
  full_name: string;
  headline: string;
  bio_md: string;
  open_to_work: boolean;
  github: string;
  linkedin: string;
  instagram: string;
  twitter: string;
  email_contact: string;
}

export const ProfileAdminPage: React.FC = () => {
  const queryClient = useQueryClient();
  const [bioTab, setBioTab] = useState<'write' | 'preview'>('write');
  const [uploadingAvatar, setUploadingAvatar] = useState(false);
  const [uploadingResume, setUploadingResume] = useState(false);
  const [avatarPath, setAvatarPath] = useState<string | undefined>();
  const [resumePath, setResumePath] = useState<string | undefined>();

  const { data: profile, isLoading } = useQuery({
    queryKey: ['profile'],
    queryFn: api.getProfile,
  });

  const {
    register,
    handleSubmit,
    reset,
    watch,
  } = useForm<ProfileFormData>();

  const watchedBio = watch('bio_md');

  useEffect(() => {
    if (profile) {
      setAvatarPath(profile.avatar_path);
      setResumePath(profile.resume_path);
      reset({
        full_name: profile.full_name || '',
        headline: profile.headline || '',
        bio_md: profile.bio_md || '',
        open_to_work: profile.open_to_work || false,
        github: profile.social_links?.github || '',
        linkedin: profile.social_links?.linkedin || '',
        instagram: profile.social_links?.instagram || '',
        twitter: profile.social_links?.twitter || '',
        email_contact: profile.social_links?.email || '',
      });
    }
  }, [profile, reset]);

  const updateMutation = useMutation({
    mutationFn: (data: ProfileFormData) => {
      const social_links: Record<string, string> = {};
      if (data.github) social_links.github = data.github;
      if (data.linkedin) social_links.linkedin = data.linkedin;
      if (data.instagram) social_links.instagram = data.instagram;
      if (data.twitter) social_links.twitter = data.twitter;
      if (data.email_contact) social_links.email = data.email_contact;

      return api.updateProfile({
        full_name: data.full_name,
        headline: data.headline,
        bio_md: data.bio_md,
        open_to_work: data.open_to_work,
        avatar_path: avatarPath,
        resume_path: resumePath,
        github_url: data.github || undefined,
        linkedin_url: data.linkedin || undefined,
        public_email: data.email_contact || undefined,
        social_links,
      });
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['profile'] });
      toast.success('Profil berhasil diperbarui!');
    },
    onError: (err: any) => {
      toast.error(err.message || 'Gagal menyimpan profil');
    },
  });

  const handleAvatarUpload = async (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0];
    if (!file) return;

    setUploadingAvatar(true);
    try {
      const res = await api.uploadMedia('avatar', file);
      setAvatarPath(res.path);
      toast.success('Foto profil berhasil diunggah!');
    } catch (err: any) {
      toast.error(err.message || 'Gagal mengunggah foto');
    } finally {
      setUploadingAvatar(false);
    }
  };

  const handleResumeUpload = async (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0];
    if (!file) return;

    setUploadingResume(true);
    try {
      const res = await api.uploadMedia('resume', file);
      setResumePath(res.path);
      toast.success('Dokumen Resume/CV berhasil diunggah!');
    } catch (err: any) {
      toast.error(err.message || 'Gagal mengunggah resume');
    } finally {
      setUploadingResume(false);
    }
  };

  if (isLoading) {
    return (
      <div className="py-20 text-center font-mono text-xs text-text-faint">
        Memuat data profil...
      </div>
    );
  }

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-2xl font-bold text-text">Edit Profil Administrator</h1>
        <p className="text-xs text-text-muted mt-1">
          Kelola informasi publik, biografi markdown, avatar, dan dokumen CV
        </p>
      </div>

      <form onSubmit={handleSubmit((data) => updateMutation.mutate(data))} className="space-y-6">
        {/* Main Info */}
        <div className="rounded-xl bg-surface border border-border p-6 space-y-4">
          <h2 className="text-sm font-bold text-text">Informasi Utama</h2>

          <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
            <div>
              <label className="block text-xs font-medium text-text-muted mb-1.5">
                Nama Lengkap
              </label>
              <input
                type="text"
                {...register('full_name')}
                className="w-full px-3.5 py-2 text-xs rounded-md bg-surface-2 border border-border text-text focus:outline-none focus:border-accent"
              />
            </div>

            <div>
              <label className="block text-xs font-medium text-text-muted mb-1.5">
                Headline / Posisi Profesional
              </label>
              <input
                type="text"
                {...register('headline')}
                className="w-full px-3.5 py-2 text-xs rounded-md bg-surface-2 border border-border text-text focus:outline-none focus:border-accent"
              />
            </div>
          </div>

          {/* Open to work toggle */}
          <div className="pt-2">
            <label className="flex items-center gap-3 cursor-pointer">
              <input
                type="checkbox"
                {...register('open_to_work')}
                className="w-4 h-4 rounded border-border bg-surface-2 text-accent focus:ring-0"
              />
              <span className="text-xs font-medium text-text">
                Status: Open to Work (Tampilkan badge siap menerima proyek/pekerjaan baru)
              </span>
            </label>
          </div>
        </div>

        {/* Media Uploads */}
        <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
          {/* Avatar Upload */}
          <div className="rounded-xl bg-surface border border-border p-6 space-y-4">
            <h2 className="text-sm font-bold text-text">Foto Profil (Avatar)</h2>

            <div className="flex items-center gap-4">
              <div className="w-16 h-16 rounded-full overflow-hidden bg-surface-2 border border-border flex items-center justify-center flex-shrink-0">
                {avatarPath ? (
                  <img
                    src={getMediaUrl(avatarPath)}
                    alt="Avatar"
                    className="w-full h-full object-cover"
                  />
                ) : (
                  <span className="text-[10px] text-text-faint font-mono">Kosong</span>
                )}
              </div>

              <div className="flex-1">
                <label className="inline-flex items-center gap-2 px-3 py-1.5 rounded-md bg-surface-2 border border-border text-xs text-text hover:border-border-strong cursor-pointer">
                  {uploadingAvatar ? (
                    <Loader2 className="w-3.5 h-3.5 animate-spin" />
                  ) : (
                    <Upload className="w-3.5 h-3.5 text-text-muted" />
                  )}
                  <span>{uploadingAvatar ? 'Mengunggah...' : 'Pilih Foto Baru'}</span>
                  <input
                    type="file"
                    accept="image/png,image/jpeg,image/webp"
                    onChange={handleAvatarUpload}
                    disabled={uploadingAvatar}
                    className="hidden"
                  />
                </label>
                <p className="text-[10px] text-text-faint mt-1">Maks. 2 MB (.jpg, .png, .webp)</p>
              </div>
            </div>
          </div>

          {/* Resume Upload */}
          <div className="rounded-xl bg-surface border border-border p-6 space-y-4">
            <h2 className="text-sm font-bold text-text">Dokumen Resume / CV</h2>

            <div className="flex items-center gap-4">
              <div className="w-16 h-16 rounded-lg bg-surface-2 border border-border flex items-center justify-center text-text-faint flex-shrink-0">
                <FileText className="w-6 h-6 text-accent" />
              </div>

              <div className="flex-1">
                <label className="inline-flex items-center gap-2 px-3 py-1.5 rounded-md bg-surface-2 border border-border text-xs text-text hover:border-border-strong cursor-pointer">
                  {uploadingResume ? (
                    <Loader2 className="w-3.5 h-3.5 animate-spin" />
                  ) : (
                    <Upload className="w-3.5 h-3.5 text-text-muted" />
                  )}
                  <span>{uploadingResume ? 'Mengunggah...' : 'Unggah PDF CV'}</span>
                  <input
                    type="file"
                    accept="application/pdf"
                    onChange={handleResumeUpload}
                    disabled={uploadingResume}
                    className="hidden"
                  />
                </label>
                {resumePath && (
                  <p className="text-[10px] text-emerald-400 mt-1 flex items-center gap-1 font-mono">
                    <CheckCircle className="w-3 h-3" />
                    <span>File terpasang: {resumePath}</span>
                  </p>
                )}
                <p className="text-[10px] text-text-faint mt-0.5">Format PDF saja (Maks. 2 MB)</p>
              </div>
            </div>
          </div>
        </div>

        {/* Bio Markdown Editor */}
        <div className="rounded-xl bg-surface border border-border p-6 space-y-4">
          <div className="flex items-center justify-between">
            <div>
              <h2 className="text-sm font-bold text-text">Biografi Diri</h2>
              <p className="text-[11px] text-text-faint">
                Mendukung sintaks Markdown lengkap (otomatis dibersihkan dengan UGCPolicy & DOMPurify)
              </p>
            </div>

            <div className="flex gap-1 p-1 rounded-md bg-surface-2 border border-border text-xs">
              <button
                type="button"
                onClick={() => setBioTab('write')}
                className={`px-3 py-1 rounded flex items-center gap-1.5 transition-colors ${
                  bioTab === 'write' ? 'bg-accent text-accent-fg font-semibold' : 'text-text-muted'
                }`}
              >
                <Edit3 className="w-3.5 h-3.5" />
                <span>Tulis</span>
              </button>
              <button
                type="button"
                onClick={() => setBioTab('preview')}
                className={`px-3 py-1 rounded flex items-center gap-1.5 transition-colors ${
                  bioTab === 'preview' ? 'bg-accent text-accent-fg font-semibold' : 'text-text-muted'
                }`}
              >
                <Eye className="w-3.5 h-3.5" />
                <span>Pratinjau</span>
              </button>
            </div>
          </div>

          {bioTab === 'write' ? (
            <textarea
              rows={8}
              {...register('bio_md')}
              placeholder="Tuliskan biografi, ketertarikan teknologi, atau latar belakang Anda..."
              className="w-full p-3.5 rounded-md bg-surface-2 border border-border text-xs text-text font-mono focus:outline-none focus:border-accent"
            />
          ) : (
            <div className="p-4 rounded-md bg-surface-2 border border-border min-h-[160px]">
              <SanitizedHTML
                html={profile?.bio_html || watchedBio || '<p class="text-text-faint">Tidak ada konten pratinjau</p>'}
                className="prose prose-invert text-xs text-text-muted max-w-none"
              />
            </div>
          )}
        </div>

        {/* Social Links */}
        <div className="rounded-xl bg-surface border border-border p-6 space-y-4">
          <h2 className="text-sm font-bold text-text">Tautan Sosial & Kontak</h2>

          <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
            <div>
              <label className="block text-xs font-medium text-text-muted mb-1.5">
                GitHub URL
              </label>
              <input
                type="url"
                placeholder="https://github.com/username"
                {...register('github')}
                className="w-full px-3.5 py-2 text-xs rounded-md bg-surface-2 border border-border text-text focus:outline-none focus:border-accent"
              />
            </div>

            <div>
              <label className="block text-xs font-medium text-text-muted mb-1.5">
                LinkedIn URL
              </label>
              <input
                type="url"
                placeholder="https://linkedin.com/in/username"
                {...register('linkedin')}
                className="w-full px-3.5 py-2 text-xs rounded-md bg-surface-2 border border-border text-text focus:outline-none focus:border-accent"
              />
            </div>

            <div>
              <label className="block text-xs font-medium text-text-muted mb-1.5">
                Instagram URL
              </label>
              <input
                type="url"
                placeholder="https://instagram.com/username"
                {...register('instagram')}
                className="w-full px-3.5 py-2 text-xs rounded-md bg-surface-2 border border-border text-text focus:outline-none focus:border-accent"
              />
            </div>

            <div>
              <label className="block text-xs font-medium text-text-muted mb-1.5">
                X / Twitter URL
              </label>
              <input
                type="url"
                placeholder="https://x.com/username"
                {...register('twitter')}
                className="w-full px-3.5 py-2 text-xs rounded-md bg-surface-2 border border-border text-text focus:outline-none focus:border-accent"
              />
            </div>

            <div>
              <label className="block text-xs font-medium text-text-muted mb-1.5">
                Email Kontak Publik
              </label>
              <input
                type="email"
                placeholder="kontak@example.com"
                {...register('email_contact')}
                className="w-full px-3.5 py-2 text-xs rounded-md bg-surface-2 border border-border text-text focus:outline-none focus:border-accent"
              />
            </div>
          </div>
        </div>

        {/* Submit */}
        <div className="flex justify-end pt-2">
          <AnimatedButton
            type="submit"
            disabled={updateMutation.isPending}
          >
            {updateMutation.isPending ? (
              <>
                <Loader2 className="w-4 h-4 animate-spin inline mr-1" />
                <span>Menyimpan...</span>
              </>
            ) : (
              <>
                <Save className="w-4 h-4 inline mr-1" />
                <span>Simpan Perubahan Profil</span>
              </>
            )}
          </AnimatedButton>
        </div>
      </form>
    </div>
  );
};
