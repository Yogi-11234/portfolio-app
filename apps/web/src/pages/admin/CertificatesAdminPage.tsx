import React, { useState } from 'react';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { useForm } from 'react-hook-form';
import { toast } from 'sonner';
import { Plus, Edit2, Trash2, X, Upload, Loader2, CheckCircle } from 'lucide-react';
import { api } from '@/lib/api';
import { Certificate } from '@/types';
import { AnimatedButton } from '@/components/AnimatedButton';
import { getMediaUrl } from '@/lib/utils';

interface CertificateFormData {
  name: string;
  issuer: string;
  issued_at: string;
  expires_at?: string;
  credential_url?: string;
  order_index: number;
  is_visible: boolean;
}

export const CertificatesAdminPage: React.FC = () => {
  const queryClient = useQueryClient();
  const [modalOpen, setModalOpen] = useState(false);
  const [editingCert, setEditingCert] = useState<Certificate | null>(null);
  const [filePath, setFilePath] = useState<string | undefined>();
  const [fileMime, setFileMime] = useState<string | undefined>();
  const [uploadingFile, setUploadingFile] = useState(false);

  const { data: certificates = [], isLoading } = useQuery({
    queryKey: ['admin-certificates'],
    queryFn: api.getCertificatesAdmin,
  });

  const { register, handleSubmit, reset } = useForm<CertificateFormData>();

  const openCreateModal = () => {
    setEditingCert(null);
    setFilePath(undefined);
    setFileMime(undefined);
    reset({
      name: '',
      issuer: '',
      issued_at: '',
      expires_at: '',
      credential_url: '',
      order_index: 0,
      is_visible: true,
    });
    setModalOpen(true);
  };

  const openEditModal = (cert: Certificate) => {
    setEditingCert(cert);
    setFilePath(cert.file_path);
    setFileMime(cert.file_mime);
    reset({
      name: cert.name || (cert as any).title || '',
      issuer: cert.issuer,
      issued_at: cert.issued_at ? cert.issued_at.split('T')[0] : '',
      expires_at: cert.expires_at ? cert.expires_at.split('T')[0] : '',
      credential_url: cert.credential_url || '',
      order_index: cert.order_index ?? (cert as any).sort_order ?? 0,
      is_visible: cert.is_visible,
    });
    setModalOpen(true);
  };

  const saveMutation = useMutation({
    mutationFn: (data: CertificateFormData) => {
      if (!filePath) {
        throw new Error('Silakan unggah berkas dokumen sertifikat terlebih dahulu.');
      }

      const payload = {
        title: data.name,
        name: data.name,
        issuer: data.issuer,
        issued_at: new Date(data.issued_at).toISOString(),
        expires_at: data.expires_at ? new Date(data.expires_at).toISOString() : undefined,
        credential_url: data.credential_url || undefined,
        file_path: filePath,
        file_mime: fileMime || 'application/pdf',
        sort_order: Number(data.order_index) || 0,
        order_index: Number(data.order_index) || 0,
        is_visible: data.is_visible,
      };

      if (editingCert) {
        return api.updateCertificate(editingCert.id, payload);
      } else {
        return api.createCertificate(payload);
      }
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['admin-certificates'] });
      queryClient.invalidateQueries({ queryKey: ['certificates'] });
      toast.success(editingCert ? 'Sertifikat diperbarui' : 'Sertifikat ditambahkan');
      setModalOpen(false);
    },
    onError: (err: any) => {
      toast.error(err.message || 'Gagal menyimpan sertifikat');
    },
  });

  const deleteMutation = useMutation({
    mutationFn: (id: string) => api.deleteCertificate(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['admin-certificates'] });
      queryClient.invalidateQueries({ queryKey: ['certificates'] });
      toast.success('Sertifikat dihapus');
    },
    onError: (err: any) => {
      toast.error(err.message || 'Gagal menghapus sertifikat');
    },
  });

  const handleFileUpload = async (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0];
    if (!file) return;

    setUploadingFile(true);
    try {
      const res = await api.uploadMedia('certificate', file);
      setFilePath(res.path);
      setFileMime(res.mime);
      toast.success('Berkas sertifikat berhasil diunggah!');
    } catch (err: any) {
      toast.error(err.message || 'Gagal mengunggah berkas');
    } finally {
      setUploadingFile(false);
    }
  };

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
          <h1 className="text-2xl font-bold text-text">Manajemen Sertifikat & Lisensi</h1>
          <p className="text-xs text-text-muted mt-1">
            Kelola dokumen bukti kompetensi, sertifikasi keahlian, dan lisensi resmi
          </p>
        </div>

        <AnimatedButton
          onClick={openCreateModal}
          size="sm"
        >
          <Plus className="w-3.5 h-3.5 inline mr-1" />
          <span>Tambah Sertifikat</span>
        </AnimatedButton>
      </div>

      <div className="rounded-xl bg-surface border border-border overflow-hidden">
        {isLoading ? (
          <div className="py-16 text-center text-xs text-text-faint font-mono">
            Memuat daftar sertifikat...
          </div>
        ) : certificates.length > 0 ? (
          <div className="overflow-x-auto">
            <table className="w-full text-left text-xs">
              <thead className="bg-surface-2 border-b border-border text-text-muted font-mono uppercase tracking-wider text-[11px]">
                <tr>
                  <th className="py-3 px-4">Nama Sertifikasi</th>
                  <th className="py-3 px-4">Penerbit (Issuer)</th>
                  <th className="py-3 px-4">Tanggal Terbit</th>
                  <th className="py-3 px-4">Dokumen</th>
                  <th className="py-3 px-4 text-right">Aksi</th>
                </tr>
              </thead>
              <tbody className="divide-y divide-border/60">
                {certificates.map((cert) => (
                  <tr key={cert.id} className="hover:bg-surface-2/40 transition-colors">
                    <td className="py-3 px-4 font-semibold text-text">{cert.name || (cert as any).title}</td>
                    <td className="py-3 px-4 text-text-muted">{cert.issuer}</td>
                    <td className="py-3 px-4 text-text-faint font-mono">
                      {formatDate(cert.issued_at)}
                    </td>
                    <td className="py-3 px-4">
                      {cert.file_path ? (
                        <a
                          href={getMediaUrl(cert.file_path)}
                          target="_blank"
                          rel="noopener noreferrer"
                          className="text-accent underline font-mono text-[11px]"
                        >
                          Lihat File
                        </a>
                      ) : (
                        <span className="text-text-faint">-</span>
                      )}
                    </td>
                    <td className="py-3 px-4 text-right space-x-1">
                      <button
                        onClick={() => openEditModal(cert)}
                        className="p-1.5 text-text-muted hover:text-text rounded hover:bg-surface-2"
                      >
                        <Edit2 className="w-3.5 h-3.5" />
                      </button>
                      <button
                        onClick={() => {
                          if (confirm(`Hapus sertifikat "${cert.name}"?`)) {
                            deleteMutation.mutate(cert.id);
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
            Belum ada sertifikat terdaftar.
          </div>
        )}
      </div>

      {/* Modal Form */}
      {modalOpen && (
        <div className="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/80 backdrop-blur-sm">
          <div className="relative w-full max-w-md bg-surface border border-border rounded-xl shadow-2xl p-6 space-y-4">
            <div className="flex items-center justify-between pb-3 border-b border-border">
              <h2 className="text-sm font-bold text-text">
                {editingCert ? 'Edit Sertifikat' : 'Tambah Sertifikat Baru'}
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
                <label className="block text-xs font-medium text-text-muted mb-1">
                  Nama Sertifikat / Lisensi *
                </label>
                <input
                  type="text"
                  placeholder="cth. AWS Certified Solutions Architect"
                  {...register('name', { required: true })}
                  className="w-full px-3 py-2 text-xs rounded-md bg-surface-2 border border-border text-text focus:outline-none focus:border-accent"
                />
              </div>

              <div>
                <label className="block text-xs font-medium text-text-muted mb-1">
                  Organisasi Penerbit *
                </label>
                <input
                  type="text"
                  placeholder="cth. Amazon Web Services / Dicoding"
                  {...register('issuer', { required: true })}
                  className="w-full px-3 py-2 text-xs rounded-md bg-surface-2 border border-border text-text focus:outline-none focus:border-accent"
                />
              </div>

              <div className="grid grid-cols-2 gap-4">
                <div>
                  <label className="block text-xs font-medium text-text-muted mb-1">
                    Tanggal Terbit *
                  </label>
                  <input
                    type="date"
                    {...register('issued_at', { required: true })}
                    className="w-full px-3 py-2 text-xs rounded-md bg-surface-2 border border-border text-text focus:outline-none focus:border-accent"
                  />
                </div>

                <div>
                  <label className="block text-xs font-medium text-text-muted mb-1">
                    Kadaluarsa (Opsional)
                  </label>
                  <input
                    type="date"
                    {...register('expires_at')}
                    className="w-full px-3 py-2 text-xs rounded-md bg-surface-2 border border-border text-text focus:outline-none focus:border-accent"
                  />
                </div>
              </div>

              <div>
                <label className="block text-xs font-medium text-text-muted mb-1">
                  URL Kredensial Online (Opsional)
                </label>
                <input
                  type="url"
                  placeholder="https://credly.com/..."
                  {...register('credential_url')}
                  className="w-full px-3 py-2 text-xs rounded-md bg-surface-2 border border-border text-text focus:outline-none focus:border-accent"
                />
              </div>

              {/* Document upload */}
              <div>
                <label className="block text-xs font-medium text-text-muted mb-1">
                  Unggah Dokumen Sertifikat (PDF / Gambar)
                </label>
                <div className="flex items-center gap-3">
                  <label className="inline-flex items-center gap-2 px-3 py-1.5 rounded-md bg-surface-2 border border-border text-xs text-text hover:border-border-strong cursor-pointer">
                    {uploadingFile ? (
                      <Loader2 className="w-3.5 h-3.5 animate-spin" />
                    ) : (
                      <Upload className="w-3.5 h-3.5 text-text-muted" />
                    )}
                    <span>{uploadingFile ? 'Mengunggah...' : 'Pilih Berkas'}</span>
                    <input
                      type="file"
                      accept="application/pdf,image/png,image/jpeg,image/webp"
                      onChange={handleFileUpload}
                      disabled={uploadingFile}
                      className="hidden"
                    />
                  </label>
                  {filePath && (
                    <span className="text-[10px] text-emerald-400 font-mono flex items-center gap-1">
                      <CheckCircle className="w-3 h-3" />
                      <span>{filePath}</span>
                    </span>
                  )}
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
                  {saveMutation.isPending ? 'Menyimpan...' : 'Simpan Sertifikat'}
                </AnimatedButton>
              </div>
            </form>
          </div>
        </div>
      )}
    </div>
  );
};
