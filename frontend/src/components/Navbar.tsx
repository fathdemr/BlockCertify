import React from 'react';
import { Shield, Menu, X, LayoutDashboard, LogIn } from 'lucide-react';
import { Link, useLocation } from 'react-router-dom';
import { useAuth } from '../context/AuthContext';

const Navbar: React.FC = () => {
    const [isOpen, setIsOpen] = React.useState(false);
    const { isAuthenticated, isAdmin } = useAuth();
    const location = useLocation();

    const navLinks = [
        { name: 'Home', href: '/' },
        { name: 'Verify', href: '/verify' },
        { name: 'About', href: '/about' },
    ];

    // Hide Navbar if on admin pages (except login)
    if (location.pathname.startsWith('/admin')) {
        return null;
    }

    return (
        <nav className="fixed w-full z-50 bg-brand-dark/80 backdrop-blur-md border-b border-gray-800">
            <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
                <div className="flex items-center justify-between h-16">
                    <div className="flex items-center">
                        <Link to="/" className="flex-shrink-0 flex items-center gap-2">
                            <Shield className="h-8 w-8 text-brand-secondary" />
                            <span className="text-xl font-display font-bold bg-gradient-to-r from-white to-gray-400 bg-clip-text text-transparent">
                                BlockCertify
                            </span>
                        </Link>
                    </div>
                    <div className="hidden md:block">
                        <div className="ml-10 flex items-baseline space-x-8">
                            {navLinks.map((link) => (
                                <Link
                                    key={link.name}
                                    to={link.href}
                                    className={`px-3 py-2 rounded-md text-sm font-medium transition-colors ${location.pathname === link.href
                                            ? 'text-brand-secondary'
                                            : 'text-gray-300 hover:text-white'
                                        }`}
                                >
                                    {link.name}
                                </Link>
                            ))}

                            {isAuthenticated && isAdmin ? (
                                <Link
                                    to="/admin"
                                    className="flex items-center gap-2 bg-brand-accent/10 hover:bg-brand-accent/20 text-brand-accent border border-brand-accent/30 px-4 py-2 rounded-lg text-sm font-semibold transition-all"
                                >
                                    <LayoutDashboard className="h-4 w-4" />
                                    Dashboard
                                </Link>
                            ) : (
                                <Link
                                    to="/login"
                                    className="flex items-center gap-2 bg-brand-primary hover:bg-brand-primary/90 text-white px-4 py-2 rounded-lg text-sm font-semibold transition-all shadow-[0_0_15px_rgba(59,130,246,0.3)]"
                                >
                                    <LogIn className="h-4 w-4" />
                                    Institution Login
                                </Link>
                            )}
                        </div>
                    </div>
                    <div className="flex md:hidden">
                        <button
                            onClick={() => setIsOpen(!isOpen)}
                            className="inline-flex items-center justify-center p-2 rounded-md text-gray-400 hover:text-white focus:outline-none"
                        >
                            {isOpen ? <X className="h-6 w-6" /> : <Menu className="h-6 w-6" />}
                        </button>
                    </div>
                </div>
            </div>

            {/* Mobile menu */}
            {isOpen && (
                <div className="md:hidden bg-brand-dark/95 border-b border-gray-800">
                    <div className="px-2 pt-2 pb-3 space-y-1 sm:px-3">
                        {navLinks.map((link) => (
                            <Link
                                key={link.name}
                                to={link.href}
                                className={`block px-3 py-2 rounded-md text-base font-medium ${location.pathname === link.href
                                        ? 'text-brand-secondary'
                                        : 'text-gray-300 hover:text-white'
                                    }`}
                                onClick={() => setIsOpen(false)}
                            >
                                {link.name}
                            </Link>
                        ))}
                        {isAuthenticated && isAdmin ? (
                            <Link
                                to="/admin"
                                className="block px-3 py-2 rounded-md text-base font-medium text-brand-accent"
                                onClick={() => setIsOpen(false)}
                            >
                                Dashboard
                            </Link>
                        ) : (
                            <Link
                                to="/login"
                                className="block px-3 py-2 rounded-md text-base font-medium text-brand-secondary"
                                onClick={() => setIsOpen(false)}
                            >
                                Institution Login
                            </Link>
                        )}
                    </div>
                </div>
            )}
        </nav>
    );
};

export default Navbar;
