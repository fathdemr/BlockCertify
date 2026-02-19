import React, { createContext, useContext, useState, useEffect } from 'react';
import api from '../services/api';

type Role = 'admin' | 'public';

interface User {
    email: string;
    role: Role;
    token?: string;
}

interface AuthContextType {
    user: User | null;
    login: (email: string, password: string) => Promise<void>;
    register: (data: { firstName: string; lastName: string; email: string; universityID: string; password: string }) => Promise<void>;
    logout: () => void;
    isAuthenticated: boolean;
    isAdmin: boolean;
    loading: boolean;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export const AuthProvider: React.FC<{ children: React.ReactNode }> = ({ children }) => {
    const [user, setUser] = useState<User | null>(null);
    const [loading, setLoading] = useState(true);

    useEffect(() => {
        const savedUser = localStorage.getItem('blockcertify_user');
        if (savedUser) {
            setUser(JSON.parse(savedUser));
        }
        setLoading(false);
    }, []);

    const login = async (email: string, password: string) => {
        try {
            const response = await api.post('/v1/auth/user/login', { email, password });
            const { token, role } = response.data;
            const newUser: User = { email, role, token }; // Note: User interface might need Update
            setUser(newUser);
            localStorage.setItem('blockcertify_user', JSON.stringify(newUser));
        } catch (error: any) {
            const message = error.response?.data?.message || 'Login failed';
            throw new Error(message);
        }
    };

    const register = async (data: { firstName: string; lastName: string; email: string; universityID: string; password: string }) => {
        try {
            await api.post('/v1/auth/user/register/admin', data);
        } catch (error: any) {
            const message = error.response?.data?.error || error.response?.data?.message || 'Registration failed';
            throw new Error(message);
        }
    };

    const logout = async () => {
        try {
            await api.post('/v1/auth/user/logout');
        } catch (error) {
            console.error('Logout error:', error);
        } finally {
            setUser(null);
            localStorage.removeItem('blockcertify_user');
            window.location.href = '/login';
        }
    };

    return (
        <AuthContext.Provider
            value={{
                user,
                login,
                register,
                logout,
                isAuthenticated: !!user,
                isAdmin: user?.role === 'admin',
                loading,
            }}
        >
            {children}
        </AuthContext.Provider>
    );
};

export const useAuth = () => {
    const context = useContext(AuthContext);
    if (context === undefined) {
        throw new Error('useAuth must be used within an AuthProvider');
    }
    return context;
};
