import { Link, useLocation } from 'react-router-dom';
import { ShieldCheck } from 'lucide-react';

export default function Navbar() {
  const { pathname } = useLocation();
  const isDark = pathname === '/';

  return (
    <header
      className={`fixed top-0 left-0 right-0 z-50 ${
        isDark
          ? 'glass-dark border-b border-white/10'
          : 'glass border-b border-outline-variant/30'
      }`}
    >
      <div className="max-w-[1280px] mx-auto px-6 h-16 flex items-center justify-between">
        <Link to="/" className="flex items-center gap-2">
          <ShieldCheck className={`w-6 h-6 ${isDark ? 'text-on-tertiary-container' : 'text-on-tertiary-container'}`} />
          <span className={`font-bold text-lg ${isDark ? 'text-white' : 'text-on-surface'}`}>
            BlockCertify
          </span>
        </Link>

        <nav className="hidden md:flex items-center gap-8">
          {[
            { label: 'Hakkımızda', href: '#about' },
            { label: 'Nasıl Çalışır', href: '#how' },
            { label: 'Doğrulama', href: '/dogrula' },
          ].map(({ label, href }) => (
            <a
              key={label}
              href={href}
              className={`text-sm font-medium transition-colors ${
                isDark
                  ? 'text-white/70 hover:text-white'
                  : 'text-on-surface-variant hover:text-on-surface'
              }`}
            >
              {label}
            </a>
          ))}
        </nav>

        <Link
          to="/giris"
          className={`px-5 py-2 rounded-lg text-sm font-semibold transition-all ${
            isDark
              ? 'bg-white text-primary hover:bg-white/90'
              : 'bg-primary text-on-primary hover:bg-primary/90'
          }`}
        >
          Giriş Yap
        </Link>
      </div>
    </header>
  );
}
