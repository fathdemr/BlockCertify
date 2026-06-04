import { useState, type FormEvent } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import { ShieldCheck, Eye, EyeOff } from 'lucide-react';
import { authApi } from '../services/api';

export default function Register() {
  const navigate = useNavigate();
  const [name, setName] = useState('');
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [showPass, setShowPass] = useState(false);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');

  const handleSubmit = async (e: FormEvent) => {
    e.preventDefault();
    setError('');
    setLoading(true);
    try {
      await authApi.register(name, email, password);
      navigate('/giris');
    } catch {
      setError('Kayıt sırasında bir hata oluştu.');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="min-h-screen bg-surface flex flex-col">
      <div className="px-8 py-5">
        <Link to="/" className="flex items-center gap-2 w-fit">
          <ShieldCheck className="w-6 h-6 text-on-tertiary-container" />
          <span className="font-bold text-on-surface text-lg">BlockCertify</span>
        </Link>
      </div>

      <div className="flex-1 flex items-center justify-center px-4 py-12">
        <div className="w-full max-w-md bg-white rounded-2xl card-shadow p-10">
          <h1 className="text-3xl font-bold text-on-surface text-center mb-2">Hesap Oluştur</h1>
          <p className="text-on-surface-variant text-center text-sm mb-8">
            Blockchain diploma platformuna katılın.
          </p>

          {error && (
            <div className="mb-4 px-4 py-3 rounded-lg bg-error-container text-sm text-error font-medium">
              {error}
            </div>
          )}

          <form onSubmit={handleSubmit} className="space-y-4">
            <div>
              <label className="block text-xs font-semibold text-on-surface mb-1.5">Ad Soyad</label>
              <input
                type="text"
                value={name}
                onChange={(e) => setName(e.target.value)}
                placeholder="Ad Soyad"
                required
                className="w-full px-4 py-2.5 rounded-lg border border-outline-variant/60 text-sm focus:outline-none focus:border-on-tertiary-container focus:ring-2 focus:ring-on-tertiary-container/10 placeholder:text-outline transition-all"
              />
            </div>
            <div>
              <label className="block text-xs font-semibold text-on-surface mb-1.5">E-posta Adresi</label>
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
              <label className="block text-xs font-semibold text-on-surface mb-1.5">Şifre</label>
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

            <button
              type="submit"
              disabled={loading}
              className="w-full py-3 bg-primary text-on-primary rounded-lg font-semibold text-sm hover:bg-primary/90 transition-all disabled:opacity-60"
            >
              {loading ? 'Kaydediliyor...' : 'Kayıt Ol'}
            </button>
          </form>

          <p className="text-center text-sm text-on-surface-variant mt-6">
            Zaten hesabınız var mı?{' '}
            <Link to="/giris" className="text-on-tertiary-container font-semibold hover:underline">
              Giriş Yap
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
