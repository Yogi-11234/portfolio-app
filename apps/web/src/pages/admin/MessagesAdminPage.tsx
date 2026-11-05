import React, { useState } from 'react';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { toast } from 'sonner';
import {
  Mail,
  Trash2,
  Archive,
  AlertTriangle,
  Eye,
  X,
  Clock,
  User,
} from 'lucide-react';
import { api } from '@/lib/api';
import { Message } from '@/types';
import { AnimatedButton } from '@/components/AnimatedButton';

type MessageFilterStatus = '' | 'unread' | 'read' | 'archived' | 'spam';

export const MessagesAdminPage: React.FC = () => {
  const queryClient = useQueryClient();
  const [statusFilter, setStatusFilter] = useState<MessageFilterStatus>('');
  const [selectedMessage, setSelectedMessage] = useState<Message | null>(null);

  const { data: messagesData, isLoading } = useQuery({
    queryKey: ['admin-messages', statusFilter],
    queryFn: () =>
      api.getMessagesAdmin({
        status: statusFilter || undefined,
        per_page: 50,
      }),
  });

  const messages = messagesData?.messages || [];

  const updateStatusMutation = useMutation({
    mutationFn: ({ id, status }: { id: string; status: 'unread' | 'read' | 'archived' | 'spam' }) =>
      api.updateMessageStatus(id, status),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['admin-messages'] });
      queryClient.invalidateQueries({ queryKey: ['admin-messages-stats'] });
      toast.success('Status pesan diperbarui');
    },
    onError: (err: any) => {
      toast.error(err.message || 'Gagal mengubah status');
    },
  });

  const deleteMutation = useMutation({
    mutationFn: (id: string) => api.deleteMessage(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['admin-messages'] });
      queryClient.invalidateQueries({ queryKey: ['admin-messages-stats'] });
      toast.success('Pesan dihapus');
      if (selectedMessage) setSelectedMessage(null);
    },
    onError: (err: any) => {
      toast.error(err.message || 'Gagal menghapus pesan');
    },
  });

  const openDetail = (m: Message) => {
    setSelectedMessage(m);
    if (m.status === 'unread') {
      updateStatusMutation.mutate({ id: m.id, status: 'read' });
    }
  };

  const formatDate = (iso: string) => {
    try {
      return new Date(iso).toLocaleString('id-ID', {
        day: 'numeric',
        month: 'short',
        year: 'numeric',
        hour: '2-digit',
        minute: '2-digit',
      });
    } catch {
      return iso;
    }
  };

  return (
    <div className="space-y-6">
      <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
        <div>
          <h1 className="text-2xl font-bold text-text">Kotak Masuk Pesan (Inbox)</h1>
          <p className="text-xs text-text-muted mt-1">
            Daftar pesan dari formulir kontak publik pengunjung portofolio
          </p>
        </div>

        {/* Filter Tabs */}
        <div className="flex gap-1 p-1 rounded-md bg-surface border border-border text-xs">
          {[
            { key: '', label: 'Semua' },
            { key: 'unread', label: 'Belum Dibaca' },
            { key: 'read', label: 'Dibaca' },
            { key: 'archived', label: 'Arsip' },
            { key: 'spam', label: 'Spam' },
          ].map((tab) => (
            <button
              key={tab.key}
              onClick={() => setStatusFilter(tab.key as MessageFilterStatus)}
              className={`px-3 py-1.5 rounded transition-colors ${
                statusFilter === tab.key
                  ? 'bg-accent text-accent-fg font-semibold'
                  : 'text-text-muted hover:text-text'
              }`}
            >
              {tab.label}
            </button>
          ))}
        </div>
      </div>

      {/* Messages Table */}
      <div className="rounded-xl bg-surface border border-border overflow-hidden">
        {isLoading ? (
          <div className="py-16 text-center text-xs text-text-faint font-mono">
            Memuat kotak masuk...
          </div>
        ) : messages.length > 0 ? (
          <div className="overflow-x-auto">
            <table className="w-full text-left text-xs">
              <thead className="bg-surface-2 border-b border-border text-text-muted font-mono uppercase tracking-wider text-[11px]">
                <tr>
                  <th className="py-3 px-4">Pengirim</th>
                  <th className="py-3 px-4">Subjek / Cuplikan</th>
                  <th className="py-3 px-4">Status</th>
                  <th className="py-3 px-4">Waktu Diterima</th>
                  <th className="py-3 px-4 text-right">Aksi</th>
                </tr>
              </thead>
              <tbody className="divide-y divide-border/60">
                {messages.map((m) => (
                  <tr
                    key={m.id}
                    className={`hover:bg-surface-2/40 transition-colors ${
                      m.status === 'unread' ? 'bg-surface-2/20 font-medium' : ''
                    }`}
                  >
                    <td className="py-3 px-4">
                      <div className="text-text font-semibold">{m.name}</div>
                      <div className="text-[11px] text-text-faint font-mono">{m.email}</div>
                    </td>
                    <td className="py-3 px-4 max-w-xs">
                      <div className="text-text truncate">{m.subject || '(Tanpa Subjek)'}</div>
                      <div className="text-[11px] text-text-muted truncate mt-0.5">{m.body}</div>
                    </td>
                    <td className="py-3 px-4">
                      <span
                        className={`px-2 py-0.5 rounded text-[10px] font-mono font-semibold ${
                          m.status === 'unread'
                            ? 'bg-rose-950/80 text-rose-400 border border-rose-800'
                            : m.status === 'read'
                            ? 'bg-emerald-950/80 text-emerald-400 border border-emerald-800'
                            : m.status === 'archived'
                            ? 'bg-zinc-800 text-text-faint border border-border'
                            : 'bg-amber-950/80 text-amber-400 border border-amber-800'
                        }`}
                      >
                        {m.status.toUpperCase()}
                      </span>
                    </td>
                    <td className="py-3 px-4 text-text-faint font-mono whitespace-nowrap">
                      {formatDate(m.created_at)}
                    </td>
                    <td className="py-3 px-4 text-right space-x-1 whitespace-nowrap">
                      <button
                        onClick={() => openDetail(m)}
                        className="p-1.5 text-text-muted hover:text-text rounded hover:bg-surface-2"
                        title="Buka Pesan"
                      >
                        <Eye className="w-3.5 h-3.5" />
                      </button>
                      <button
                        onClick={() => {
                          if (confirm('Hapus pesan ini?')) {
                            deleteMutation.mutate(m.id);
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
            Tidak ada pesan pada filter ini.
          </div>
        )}
      </div>

      {/* Message Detail Modal */}
      {selectedMessage && (
        <div className="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/80 backdrop-blur-sm">
          <div className="relative w-full max-w-lg bg-surface border border-border rounded-xl shadow-2xl p-6 space-y-4">
            <div className="flex items-center justify-between pb-3 border-b border-border">
              <div>
                <h2 className="text-sm font-bold text-text">Detail Pesan Masuk</h2>
                <p className="text-[11px] text-text-faint font-mono">ID: {selectedMessage.id}</p>
              </div>
              <button
                onClick={() => setSelectedMessage(null)}
                className="p-1.5 text-text-muted hover:text-text rounded hover:bg-surface-2"
              >
                <X className="w-4 h-4" />
              </button>
            </div>

            <div className="space-y-3 text-xs">
              <div className="p-3 rounded-lg bg-surface-2 border border-border space-y-1">
                <div className="flex items-center gap-1.5 text-text-faint font-mono">
                  <User className="w-3.5 h-3.5" />
                  <span>Pengirim:</span>
                  <span className="text-text font-bold">{selectedMessage.name}</span>
                  <span>&lt;{selectedMessage.email}&gt;</span>
                </div>
                <div className="flex items-center gap-1.5 text-text-faint font-mono">
                  <Clock className="w-3.5 h-3.5" />
                  <span>Diterima:</span>
                  <span className="text-text">{formatDate(selectedMessage.created_at)}</span>
                </div>
                {selectedMessage.user_agent && (
                  <p className="text-[10px] text-text-faint truncate font-mono">
                    User-Agent: {selectedMessage.user_agent}
                  </p>
                )}
              </div>

              <div>
                <span className="text-[11px] text-text-faint font-mono block mb-1">Subjek:</span>
                <p className="text-xs font-semibold text-text bg-surface-2 p-2.5 rounded border border-border">
                  {selectedMessage.subject || '(Tanpa Subjek)'}
                </p>
              </div>

              <div>
                <span className="text-[11px] text-text-faint font-mono block mb-1">Isi Pesan:</span>
                <div className="text-xs text-text bg-surface-2 p-3.5 rounded border border-border whitespace-pre-wrap leading-relaxed max-h-60 overflow-y-auto">
                  {selectedMessage.body}
                </div>
              </div>
            </div>

            {/* Status change actions */}
            <div className="pt-3 border-t border-border flex items-center justify-between">
              <div className="flex items-center gap-1">
                <button
                  onClick={() =>
                    updateStatusMutation.mutate({ id: selectedMessage.id, status: 'archived' })
                  }
                  className="px-2.5 py-1.5 text-[11px] rounded bg-surface-2 border border-border text-text hover:border-border-strong flex items-center gap-1"
                >
                  <Archive className="w-3 h-3" />
                  <span>Arsipkan</span>
                </button>
                <button
                  onClick={() =>
                    updateStatusMutation.mutate({ id: selectedMessage.id, status: 'spam' })
                  }
                  className="px-2.5 py-1.5 text-[11px] rounded bg-surface-2 border border-border text-amber-400 hover:border-amber-600 flex items-center gap-1"
                >
                  <AlertTriangle className="w-3 h-3" />
                  <span>Tandai Spam</span>
                </button>
              </div>

              <div className="flex items-center gap-2">
                <AnimatedButton
                  as="a"
                  size="sm"
                  href={`mailto:${selectedMessage.email}?subject=Balasan: ${encodeURIComponent(
                    selectedMessage.subject || 'Pesan dari Portofolio'
                  )}`}
                >
                  <Mail className="w-3.5 h-3.5 inline mr-1" />
                  <span>Balas Email</span>
                </AnimatedButton>
              </div>
            </div>
          </div>
        </div>
      )}
    </div>
  );
};
