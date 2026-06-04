import { createContext, useContext, useState, useEffect, type ReactNode } from 'react';
import type { User } from '../types';
import { authApi } from '../services/api';

interface AuthContextValue {
  user: User | null;
  loading: boolean;
  login: (email: string, password: string) => Promise<void>;
  logout: () => Promise<void>;
  isAdmin: boolean;
}

const AuthContext = createContext<AuthContextValue | null>(null);

export function AuthProvider({ children }: Readonly<{ children: ReactNode }>) {
  const [user, setUser] = useState<User | null>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const stored = localStorage.getItem('bc_user');
    if (stored) {
      try {
        setUser(JSON.parse(stored));
      } catch {
        localStorage.removeItem('bc_user');
      }
    }
    setLoading(false);
  }, []);

  const login = async (email: string, password: string) => {
    const res = await authApi.login(email, password);
    const userData: User = res.data.user ?? res.data;
    setUser(userData);
    localStorage.setItem('bc_user', JSON.stringify(userData));
  };

  const logout = async () => {
    await authApi.logout().catch(() => {});
    setUser(null);
    localStorage.removeItem('bc_user');
  };

  return (
    <AuthContext.Provider value={{ user, loading, login, logout, isAdmin: user?.role === 'admin' }}>
      {children}
    </AuthContext.Provider>
  );
}

export function useAuth() {
  const ctx = useContext(AuthContext);
  if (!ctx) throw new Error('useAuth must be used inside AuthProvider');
  return ctx;
}
