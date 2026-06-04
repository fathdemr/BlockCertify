import { ShieldCheck } from 'lucide-react';
import { Link } from 'react-router-dom';

const columns = [
  {
    title: 'Platform',
    links: [
      { label: 'Ana Sayfa', href: '/' },
      { label: 'Diploma Yükle', href: '/admin/yukle' },
      { label: 'Doğrulama', href: '/dogrula' },
    ],
  },
  {
    title: 'Kaynaklar',
    links: [
      { label: 'Belgeleme', href: '#' },
      { label: 'API Referansı', href: '#' },
      { label: 'Destek', href: '#' },
    ],
  },
  {
    title: 'Yasal',
    links: [
      { label: 'Gizlilik Politikası', href: '#' },
      { label: 'Kullanım Koşulları', href: '#' },
      { label: 'KVKK', href: '#' },
    ],
  },
];

export default function Footer() {
  return (
    <footer className="bg-primary-container text-white/70">
      <div className="max-w-[1280px] mx-auto px-6 py-16">
        <div className="grid grid-cols-1 md:grid-cols-4 gap-12">
          <div className="md:col-span-1">
            <div className="flex items-center gap-2 mb-4">
              <ShieldCheck className="w-6 h-6 text-secondary-container" />
              <span className="font-bold text-white text-lg">BlockCertify</span>
            </div>
            <p className="text-sm leading-relaxed">
              Blockchain tabanlı diploma doğrulama platformu. Akademik bilgilerinizi güvenle saklayın.
            </p>
          </div>

          {columns.map((col) => (
            <div key={col.title}>
              <h4 className="text-white font-semibold mb-4 text-sm tracking-wide">{col.title}</h4>
              <ul className="space-y-3">
                {col.links.map((link) => (
                  <li key={link.label}>
                    <Link
                      to={link.href}
                      className="text-sm hover:text-white transition-colors"
                    >
                      {link.label}
                    </Link>
                  </li>
                ))}
              </ul>
            </div>
          ))}
        </div>

        <div className="border-t border-white/10 mt-12 pt-8 flex flex-col md:flex-row items-center justify-between gap-4">
          <p className="text-xs">© 2024 BlockCertify Inc. Tüm hakları saklıdır.</p>
          <p className="text-xs">Blockchain Ağı Aktif</p>
        </div>
      </div>
    </footer>
  );
}
