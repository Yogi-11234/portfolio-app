import React, { useState } from 'react';
import { useNavigate, Link } from 'react-router-dom';
import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import { z } from 'zod';
import { toast } from 'sonner';
import { Lock, Mail, Loader2, ArrowLeft, ShieldAlert } from 'lucide-react';
import { api } from '@/lib/api';
import { AnimatedButton } from '@/components/AnimatedButton';

const loginSchema = z.object({
  email: z.string().email('Format email tidak valid'),
  password: z.string().min(1, 'Password wajib diisi'),
});

type LoginFormData = z.infer<typeof loginSchema>;

export const LoginPage: React.FC = () => {
  const [loading, setLoading] = useState(false);
  const [errorMsg, setErrorMsg] = useState<string | null>(null);
  const navigate = useNavigate();

  const {
    register,
    handleSubmit,
    formState: { errors },
  } = useForm<LoginFormData>({
    resolver: zodResolver(loginSchema),
  });

  const onSubmit = async (data: LoginFormData) => {
    setLoading(true);
    setErrorMsg(null);
    try {
      await api.login(data);
      toast.success('Login berhasil! Mengalihkan ke dashboard...');
      navigate('/admin');
    } catch (err: any) {
      if (err.code === 'ACCOUNT_LOCKED') {
        setErrorMsg('Akun terkunci sementara karena 5 kali percobaan gagal berturut-turut. Silakan tunggu 15 menit.');
      } else if (err.code === 'RATE_LIMIT_EXCEEDED') {
        setErrorMsg('Terlalu banyak percobaan login. Silakan tunggu beberapa saat.');
      } else {
        setErrorMsg(err.message || 'Email atau password salah');
      }
      toast.error('Gagal masuk');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="min-h-screen flex flex-col justify-center items-center px-4 bg-bg">
      <div className="w-full max-w-md">
        <Link
          to="/"
          className="inline-flex items-center gap-1.5 text-xs text-text-muted hover:text-text font-mono mb-6 transition-colors"
        >
          <ArrowLeft className="w-3.5 h-3.5" />
          <span>Kembali ke Beranda</span>
        </Link>

        <div className="rounded-xl bg-surface border border-border p-8 shadow-2xl">
          <div className="text-center mb-6">
            <div className="w-12 h-12 mx-auto rounded-lg bg-surface-2 border border-border flex items-center justify-center text-accent mb-3">
              <Lock className="w-6 h-6" />
            </div>
            <h1 className="text-xl font-bold text-text">Admin CMS Dashboard</h1>
            <p className="text-xs text-text-muted mt-1">
              Masuk untuk mengelola portofolio, konten, dan pesan
            </p>
          </div>

          {errorMsg && (
            <div className="mb-6 p-3.5 rounded-lg bg-red-950/50 border border-red-800 text-red-200 text-xs flex items-start gap-2.5">
              <ShieldAlert className="w-4 h-4 flex-shrink-0 mt-0.5 text-red-400" />
              <span>{errorMsg}</span>
            </div>
          )}

          <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
            <div>
              <label className="block text-xs font-medium text-text-muted mb-1.5">
                Email Administrator
              </label>
              <div className="relative">
                <Mail className="w-4 h-4 absolute left-3 top-1/2 -translate-y-1/2 text-text-faint" />
                <input
                  type="email"
                  placeholder="admin@local.test"
                  {...register('email')}
                  className={`w-full pl-9 pr-3.5 py-2.5 rounded-md bg-surface-2 border text-xs text-text placeholder:text-text-faint focus:outline-none focus:border-accent transition-colors ${
                    errors.email ? 'border-danger' : 'border-border'
                  }`}
                />
              </div>
              {errors.email && (
                <p className="mt-1 text-[11px] text-danger">{errors.email.message}</p>
              )}
            </div>

            <div>
              <label className="block text-xs font-medium text-text-muted mb-1.5">
                Password
              </label>
              <div className="relative">
                <Lock className="w-4 h-4 absolute left-3 top-1/2 -translate-y-1/2 text-text-faint" />
                <input
                  type="password"
                  placeholder="••••••••••••"
                  {...register('password')}
                  className={`w-full pl-9 pr-3.5 py-2.5 rounded-md bg-surface-2 border text-xs text-text placeholder:text-text-faint focus:outline-none focus:border-accent transition-colors ${
                    errors.password ? 'border-danger' : 'border-border'
                  }`}
                />
              </div>
              {errors.password && (
                <p className="mt-1 text-[11px] text-danger">{errors.password.message}</p>
              )}
            </div>

            <div className="pt-2">
              <AnimatedButton
                type="submit"
                disabled={loading}
                fullWidth
              >
                {loading ? (
                  <>
                    <Loader2 className="w-4 h-4 animate-spin inline mr-1" />
                    <span>Memverifikasi...</span>
                  </>
                ) : (
                  <span>Masuk ke Dashboard</span>
                )}
              </AnimatedButton>
            </div>
          </form>
        </div>
      </div>
    </div>
  );
};
