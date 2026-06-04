import { Outlet } from 'react-router-dom';
import Sidebar from '../components/admin/Sidebar';
import { Bell, HelpCircle, Search } from 'lucide-react';
import { useAuth } from '../context/AuthContext';

export default function AdminLayout() {
  const { user } = useAuth();

  return (
    <div className="flex min-h-screen bg-surface">
      <Sidebar />

      <div className="flex-1 flex flex-col min-w-0">
        {/* Top Bar */}
        <header className="h-14 bg-white border-b border-outline-variant/30 flex items-center px-6 gap-4 shrink-0">
          <div className="flex-1 max-w-xs">
            <div className="relative">
              <Search className="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-outline" />
              <input
                type="text"
                placeholder="Sertifika ara..."
                className="w-full pl-9 pr-4 py-2 text-sm bg-surface rounded-lg border border-outline-variant/50 focus:outline-none focus:border-on-tertiary-container placeholder:text-outline"
              />
            </div>
          </div>

          <div className="flex items-center gap-3 ml-auto">
            <button className="w-8 h-8 flex items-center justify-center rounded-lg hover:bg-surface transition-colors">
              <Bell className="w-4 h-4 text-on-surface-variant" />
            </button>
            <button className="w-8 h-8 flex items-center justify-center rounded-lg hover:bg-surface transition-colors">
              <HelpCircle className="w-4 h-4 text-on-surface-variant" />
            </button>
            <div className="flex items-center gap-2 pl-3 border-l border-outline-variant/30">
              <div className="text-right">
                <p className="text-sm font-semibold text-on-surface leading-tight">
                  {user?.name ?? 'Admin'}
                </p>
                <p className="text-xs text-on-surface-variant capitalize">{user?.role ?? 'Admin'}</p>
              </div>
              <div className="w-8 h-8 rounded-full bg-on-tertiary-container flex items-center justify-center text-white text-xs font-bold shrink-0">
                {(user?.name ?? 'A').charAt(0).toUpperCase()}
              </div>
            </div>
          </div>
        </header>

        {/* Page Content */}
        <main className="flex-1 overflow-auto">
          <Outlet />
        </main>
      </div>
    </div>
  );
}
