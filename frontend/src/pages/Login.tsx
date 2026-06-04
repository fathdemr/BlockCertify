import { useState, type FormEvent } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import { ShieldCheck, Eye, EyeOff } from 'lucide-react';
import { useAuth } from '../context/AuthContext';

export default function Login() {
  const navigate = useNavigate();
  const { login } = useAuth();

  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [showPass, setShowPass] = useState(false);
  const [remember, setRemember] = useState(false);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');

  const handleSubmit = async (e: FormEvent) => {
    e.preventDefault();
    setError('');
    setLoading(true);
    try {
      await login(email, password);
      navigate('/admin');
    } catch {
      setError('E-posta veya şifre hatalı.');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="min-h-screen bg-surface flex flex-col">
      {/* Logo */}
      <div className="px-8 py-5">
        <Link to="/" className="flex items-center gap-2 w-fit">
          <ShieldCheck className="w-6 h-6 text-on-tertiary-container" />
          <span className="font-bold text-on-surface text-lg">BlockCertify</span>
        </Link>
      </div>

      {/* Card */}
      <div className="flex-1 flex items-center justify-center px-4 py-12">
        <div className="w-full max-w-md bg-white rounded-2xl card-shadow p-10">
          <h1 className="text-3xl font-bold text-on-surface text-center mb-2">
            Tekrar Hoş Geldiniz
          </h1>
          <p className="text-on-surface-variant text-center text-sm mb-8">
            Sertifikalarınızı yönetmek için giriş yapın.
          </p>

          {/* SSO Buttons */}
          <div className="grid grid-cols-2 gap-3 mb-6">
            <button className="flex items-center justify-center gap-2 px-4 py-2.5 rounded-lg border border-outline-variant/60 text-sm font-medium text-on-surface hover:bg-surface transition-all">
              <svg className="w-4 h-4" viewBox="0 0 24 24">
                <path fill="#4285F4" d="M22.56 12.25c0-.78-.07-1.53-.2-2.25H12v4.26h5.92c-.26 1.37-1.04 2.53-2.21 3.31v2.77h3.57c2.08-1.92 3.28-4.74 3.28-8.09z"/>
                <path fill="#34A853" d="M12 23c2.97 0 5.46-.98 7.28-2.66l-3.57-2.77c-.98.66-2.23 1.06-3.71 1.06-2.86 0-5.29-1.93-6.16-4.53H2.18v2.84C3.99 20.53 7.7 23 12 23z"/>
                <path fill="#FBBC05" d="M5.84 14.09c-.22-.66-.35-1.36-.35-2.09s.13-1.43.35-2.09V7.07H2.18C1.43 8.55 1 10.22 1 12s.43 3.45 1.18 4.93l2.85-2.22.81-.62z"/>
                <path fill="#EA4335" d="M12 5.38c1.62 0 3.06.56 4.21 1.64l3.15-3.15C17.45 2.09 14.97 1 12 1 7.7 1 3.99 3.47 2.18 7.07l3.66 2.84c.87-2.6 3.3-4.53 6.16-4.53z"/>
              </svg>
              Google
            </button>
            <button className="flex items-center justify-center gap-2 px-4 py-2.5 rounded-lg border border-outline-variant/60 text-sm font-medium text-on-surface hover:bg-surface transition-all">
              <svg className="w-4 h-4" viewBox="0 0 24 24" fill="currentColor">
                <path d="M18.71 19.5c-.83 1.24-1.71 2.45-3.05 2.47-1.34.03-1.77-.79-3.29-.79-1.53 0-2 .77-3.27.82-1.31.05-2.3-1.32-3.14-2.53C4.25 17 2.94 12.45 4.7 9.39c.87-1.52 2.43-2.48 4.12-2.51 1.28-.02 2.5.87 3.29.87.78 0 2.26-1.07 3.8-.91.65.03 2.47.26 3.64 1.98-.09.06-2.17 1.28-2.15 3.81.03 3.02 2.65 4.03 2.68 4.04-.03.07-.42 1.44-1.38 2.83M13 3.5c.73-.83 1.94-1.46 2.94-1.5.13 1.17-.34 2.35-1.04 3.19-.69.85-1.83 1.51-2.95 1.42-.15-1.15.41-2.35 1.05-3.11"/>
              </svg>
              Apple
            </button>
          </div>

          <div className="relative mb-6">
            <div className="absolute inset-0 flex items-center">
              <div className="w-full border-t border-outline-variant/40" />
            </div>
            <div className="relative flex justify-center">
              <span className="bg-white px-4 text-xs text-outline">veya e-posta ile</span>
            </div>
          </div>

          {error && (
            <div className="mb-4 px-4 py-3 rounded-lg bg-error-container text-sm text-error font-medium">
              {error}
            </div>
          )}

          <form onSubmit={handleSubmit} className="space-y-4">
            <div>
              <label className="block text-xs font-semibold text-on-surface mb-1.5">
                E-posta Adresi
              </label>
              <input
                type="email"
                value={email}
                onChange={(e) => setEmail(e.target.value)}
                placeholder="ornek@kurum.com"
                required
                className="w-full px-4 py-2.5 rounded-lg border border-outline-variant/60 text-sm focus:outline-none focus:border-on-tertiary-container focus:ring-2 focus:ring-on-tertiary-container/10 placeholder:text-outline transition-all"
              />
            </div>

            <div>
              <div className="flex items-center justify-between mb-1.5">
                <label className="block text-xs font-semibold text-on-surface">Şifre</label>
                <a href="#" className="text-xs text-on-tertiary-container hover:underline">
                  Şifremi Unuttum
                </a>
              </div>
              <div className="relative">
                <input
                  type={showPass ? 'text' : 'password'}
                  value={password}
                  onChange={(e) => setPassword(e.target.value)}
                  placeholder="••••••••"
                  required
                  className="w-full px-4 py-2.5 pr-11 rounded-lg border border-outline-variant/60 text-sm focus:outline-none focus:border-on-tertiary-container focus:ring-2 focus:ring-on-tertiary-container/10 placeholder:text-outline transition-all"
                />
                <button
                  type="button"
                  onClick={() => setShowPass(!showPass)}
                  className="absolute right-3 top-1/2 -translate-y-1/2 text-outline hover:text-on-surface"
                >
                  {showPass ? <EyeOff className="w-4 h-4" /> : <Eye className="w-4 h-4" />}
                </button>
              </div>
            </div>

            <label className="flex items-center gap-2 cursor-pointer">
              <input
                type="checkbox"
                checked={remember}
                onChange={(e) => setRemember(e.target.checked)}
                className="w-4 h-4 rounded border-outline-variant/60"
              />
              <span className="text-sm text-on-surface-variant">Beni hatırla</span>
            </label>

            <button
              type="submit"
              disabled={loading}
              className="w-full py-3 bg-primary text-on-primary rounded-lg font-semibold text-sm hover:bg-primary/90 transition-all disabled:opacity-60"
            >
              {loading ? 'Giriş yapılıyor...' : 'Giriş Yap'}
            </button>
          </form>

          <p className="text-center text-sm text-on-surface-variant mt-6">
            Hesabınız yok mu?{' '}
            <Link to="/kaydol" className="text-on-tertiary-container font-semibold hover:underline">
              Kaydol
            </Link>
          </p>
        </div>
      </div>

      <div className="py-5 text-center">
        <p className="text-xs text-outline tracking-widest uppercase flex items-center justify-center gap-2">
          <ShieldCheck className="w-3.5 h-3.5" />
          Secured by Blockchain Technology
        </p>
      </div>
    </div>
  );
}
