import { Navigate } from 'react-router-dom';
import { useAuth } from '../context/AuthContext';
import type { ReactNode } from 'react';

interface ProtectedRouteProps {
  children: ReactNode;
  requireAdmin?: boolean;
}

export default function ProtectedRoute({ children, requireAdmin }: Readonly<ProtectedRouteProps>) {
  const { user, loading } = useAuth();

  if (loading) {
    return (
      <div className="flex items-center justify-center min-h-screen bg-surface">
        <div className="w-8 h-8 border-2 border-primary border-t-transparent rounded-full animate-spin" />
      </div>
    );
  }

  if (!user) return <Navigate to="/giris" replace />;
  if (requireAdmin && user.role !== 'admin') return <Navigate to="/" replace />;

  return <>{children}</>;
}
