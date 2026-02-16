import React from 'react';
import { NavLink, useNavigate } from 'react-router-dom';
import {
    Upload,
    ShieldCheck,
    History,
    LogOut,
    LayoutDashboard,
    Shield,
    Wallet
} from 'lucide-react';
import { useAuth } from '../../context/AuthContext';

const Sidebar: React.FC = () => {
    const { logout, user } = useAuth();
    const navigate = useNavigate();

    const handleLogout = () => {
        logout();
        navigate('/');
    };

    const menuItems = [
        { icon: LayoutDashboard, label: 'Dashboard', href: '/admin' },
        { icon: Upload, label: 'Upload Diploma', href: '/admin/upload' },
        { icon: ShieldCheck, label: 'Verify Diploma', href: '/admin/verify' },
        { icon: History, label: 'History / Logs', href: '/admin/history' },
        { icon: Wallet, label: 'Wallets', href: '/admin/wallets' },
    ];

    return (
        <aside className="w-64 bg-brand-dark border-r border-gray-800 flex flex-col h-screen sticky top-0">
            <div className="p-6">
                <div className="flex items-center gap-2 mb-8">
                    <Shield className="h-8 w-8 text-brand-secondary" />
                    <span className="text-xl font-display font-bold">Admin Panel</span>
                </div>

                <nav className="space-y-2">
                    {menuItems.map((item) => (
                        <NavLink
                            key={item.href}
                            to={item.href}
                            end
                            className={({ isActive }) =>
                                `flex items-center gap-3 px-4 py-3 rounded-xl transition-all ${isActive
                                    ? 'bg-brand-primary text-white shadow-[0_0_15px_rgba(59,130,246,0.3)]'
                                    : 'text-gray-400 hover:text-white hover:bg-white/5'
                                }`
                            }
                        >
                            <item.icon className="h-5 w-5" />
                            <span className="font-medium">{item.label}</span>
                        </NavLink>
                    ))}
                </nav>
            </div>

            <div className="mt-auto p-6 border-t border-gray-800">
                <div className="flex items-center gap-3 mb-6 px-4">
                    <div className="w-10 h-10 rounded-full bg-brand-accent/20 border border-brand-accent/30 flex items-center justify-center">
                        <span className="text-brand-accent font-bold uppercase">{user?.email?.charAt(0) || 'A'}</span>
                    </div>
                    <div className="overflow-hidden">
                        <p className="text-sm font-semibold text-white truncate">{user?.email}</p>
                        <p className="text-xs text-gray-500 uppercase tracking-wider">{user?.role}</p>
                    </div>
                </div>
                <button
                    onClick={handleLogout}
                    className="w-full flex items-center gap-3 px-4 py-3 text-red-400 hover:text-red-300 hover:bg-red-500/5 rounded-xl transition-all"
                >
                    <LogOut className="h-5 w-5" />
                    <span className="font-medium">Sign Out</span>
                </button>
            </div>
        </aside>
    );
};

export default Sidebar;
