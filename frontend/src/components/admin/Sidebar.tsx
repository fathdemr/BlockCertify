import { NavLink, useNavigate } from 'react-router-dom';
import { ShieldCheck, LayoutDashboard, Upload, CheckCircle, Settings, LogOut } from 'lucide-react';
import { useAuth } from '../../context/AuthContext';

const navItems = [
  { icon: LayoutDashboard, label: 'Ana Sayfa', to: '/admin' },
  { icon: Upload, label: 'Diploma Yükle', to: '/admin/yukle' },
  { icon: CheckCircle, label: 'Diploma Doğrula', to: '/admin/dogrula' },
];

export default function Sidebar() {
  const { logout } = useAuth();
  const navigate = useNavigate();

  const handleLogout = async () => {
    await logout();
    navigate('/giris');
  };

  return (
    <aside className="w-[200px] min-h-screen bg-primary-container flex flex-col py-6 shrink-0">
      <div className="px-5 mb-8">
        <div className="flex items-center gap-2">
          <ShieldCheck className="w-6 h-6 text-secondary-container shrink-0" />
          <div>
            <p className="text-white font-bold text-sm leading-tight">BlockCertify</p>
            <p className="text-white/40 text-[10px] tracking-widest uppercase leading-tight">Blockchain Credentialing</p>
          </div>
        </div>
      </div>

      <nav className="flex-1 px-3 space-y-1">
        {navItems.map(({ icon: Icon, label, to }) => (
          <NavLink
            key={to}
            to={to}
            end={to === '/admin'}
            className={({ isActive }) =>
              `flex items-center gap-3 px-3 py-2.5 rounded-lg text-sm font-medium transition-all ${
                isActive
                  ? 'bg-white/15 text-white'
                  : 'text-white/60 hover:bg-white/10 hover:text-white'
              }`
            }
          >
            <Icon className="w-4 h-4 shrink-0" />
            {label}
          </NavLink>
        ))}
      </nav>

      <div className="px-3 mt-4 space-y-1">
        <button className="w-full flex items-center gap-3 px-3 py-2.5 rounded-lg text-sm font-medium text-white/60 hover:bg-white/10 hover:text-white transition-all">
          <Settings className="w-4 h-4 shrink-0" />
          Ayarlar
        </button>
        <button
          onClick={handleLogout}
          className="w-full flex items-center gap-3 px-3 py-2.5 rounded-lg text-sm font-medium text-red-400 hover:bg-red-500/10 hover:text-red-300 transition-all"
        >
          <LogOut className="w-4 h-4 shrink-0" />
          Çıkış Yap
        </button>
      </div>
    </aside>
  );
}
