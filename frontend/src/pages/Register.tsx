import React, { useState } from 'react';
import { useNavigate, Link } from 'react-router-dom';
import { motion } from 'framer-motion';
import { Shield, Lock, User, Mail, School, Loader2, ArrowLeft, Check, Search, ChevronDown } from 'lucide-react';
import { useAuth } from '../context/AuthContext';
import api from '../services/api';

interface University {
    id: string;
    name: string;
}

const Register: React.FC = () => {
    const [formData, setFormData] = useState({
        firstName: '',
        lastName: '',
        email: '',
        institution: '',
        password: '',
        confirmPassword: ''
    });
    const [universities, setUniversities] = useState<University[]>([]);
    const [fetchingUniversities, setFetchingUniversities] = useState(true);
    const [searchTerm, setSearchTerm] = useState('');
    const [isDropdownOpen, setIsDropdownOpen] = useState(false);

    const [loading, setLoading] = useState(false);
    const [error, setError] = useState('');
    const [success, setSuccess] = useState(false);
    const { register } = useAuth();
    const navigate = useNavigate();

    const fetchUniversities = async () => {
        if (universities.length > 0) return;

        setFetchingUniversities(true);
        try {
            const response = await api.get('/v1/universities');
            if (response && response.data) {
                setUniversities(response.data);
            }
        } catch (err) {
            console.error('Failed to fetch universities', err);
            setError('Failed to load universities. Please try again.');
        } finally {
            setFetchingUniversities(false);
        }
    };

    const toggleDropdown = () => {
        if (!isDropdownOpen) {
            fetchUniversities();
        }
        setIsDropdownOpen(!isDropdownOpen);
    };

    const filteredUniversities = universities.filter(u =>
        u.name && u.name.toLowerCase().includes(searchTerm.toLowerCase())
    );

    const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
        setFormData({
            ...formData,
            [e.target.name]: e.target.value
        });
    };

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        setError('');

        if (formData.password !== formData.confirmPassword) {
            setError('Passwords do not match');
            return;
        }

        setLoading(true);

        try {
            const { confirmPassword, ...registerData } = formData;
            await register(registerData);
            setSuccess(true);
            setTimeout(() => {
                navigate('/login');
            }, 3000);
        } catch (err: any) {
            setError(err.message || 'Registration failed');
        } finally {
            setLoading(false);
        }
    };

    if (success) {
        return (
            <div className="min-h-screen pt-32 pb-20 px-4 flex items-center justify-center">
                <motion.div
                    initial={{ opacity: 0, scale: 0.95 }}
                    animate={{ opacity: 1, scale: 1 }}
                    className="w-full max-w-md p-8 rounded-3xl bg-white/5 border border-white/10 backdrop-blur-xl shadow-2xl text-center"
                >
                    <div className="inline-flex items-center justify-center p-3 rounded-2xl bg-green-500/10 border border-green-500/20 mb-4">
                        <Shield className="h-8 w-8 text-green-500" />
                    </div>
                    <h1 className="text-3xl font-display font-bold mb-4">Registration Successful!</h1>
                    <p className="text-gray-400 mb-6">Your institutional account has been created. Redirecting to login page...</p>
                    <div className="flex justify-center">
                        <Loader2 className="h-6 w-6 animate-spin text-brand-primary" />
                    </div>
                </motion.div>
            </div>
        );
    }

    return (
        <div className="min-h-screen pt-32 pb-20 px-4 flex items-center justify-center">
            <motion.div
                initial={{ opacity: 0, scale: 0.95 }}
                animate={{ opacity: 1, scale: 1 }}
                className="w-full max-w-xl p-8 rounded-3xl bg-white/5 border border-white/10 backdrop-blur-xl shadow-2xl"
            >
                <div className="text-center mb-8">
                    <div className="inline-flex items-center justify-center p-3 rounded-2xl bg-brand-primary/10 border border-brand-primary/20 mb-4">
                        <Shield className="h-8 w-8 text-brand-primary" />
                    </div>
                    <h1 className="text-3xl font-display font-bold mb-2">Admin Registration</h1>
                    <p className="text-gray-400">Register your institution on BlockCertify</p>
                </div>

                {error && (
                    <motion.div
                        initial={{ opacity: 0, y: -10 }}
                        animate={{ opacity: 1, y: 0 }}
                        className="mb-6 p-4 bg-red-500/10 border border-red-500/20 rounded-xl text-red-500 text-sm text-center"
                    >
                        {error}
                    </motion.div>
                )}

                <form onSubmit={handleSubmit} className="space-y-4">
                    <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                        <div>
                            <label className="block text-sm font-medium text-gray-400 mb-2">First Name</label>
                            <div className="relative">
                                <User className="absolute left-4 top-1/2 -translate-y-1/2 h-5 w-5 text-gray-500" />
                                <input
                                    type="text"
                                    name="firstName"
                                    value={formData.firstName}
                                    onChange={handleChange}
                                    placeholder="John"
                                    className="w-full pl-12 pr-4 py-3 bg-white/5 border border-white/10 rounded-xl focus:outline-none focus:ring-2 focus:ring-brand-primary/50 focus:border-brand-primary transition-all"
                                    required
                                />
                            </div>
                        </div>
                        <div>
                            <label className="block text-sm font-medium text-gray-400 mb-2">Last Name</label>
                            <div className="relative">
                                <User className="absolute left-4 top-1/2 -translate-y-1/2 h-5 w-5 text-gray-500" />
                                <input
                                    type="text"
                                    name="lastName"
                                    value={formData.lastName}
                                    onChange={handleChange}
                                    placeholder="Doe"
                                    className="w-full pl-12 pr-4 py-3 bg-white/5 border border-white/10 rounded-xl focus:outline-none focus:ring-2 focus:ring-brand-primary/50 focus:border-brand-primary transition-all"
                                    required
                                />
                            </div>
                        </div>
                    </div>

                    <div>
                        <label className="block text-sm font-medium text-gray-400 mb-2">E-mail</label>
                        <div className="relative">
                            <Mail className="absolute left-4 top-1/2 -translate-y-1/2 h-5 w-5 text-gray-500" />
                            <input
                                type="email"
                                name="email"
                                value={formData.email}
                                onChange={handleChange}
                                placeholder="institution@university.edu"
                                className="w-full pl-12 pr-4 py-3 bg-white/5 border border-white/10 rounded-xl focus:outline-none focus:ring-2 focus:ring-brand-primary/50 focus:border-brand-primary transition-all"
                                required
                            />
                        </div>
                    </div>

                    <div>
                        <label className="block text-sm font-medium text-gray-400 mb-2">University / Institution</label>
                        <div className="relative">
                            <div
                                onClick={toggleDropdown}
                                className={`w-full pl-12 pr-10 py-3 bg-white/5 border border-white/10 rounded-xl flex items-center justify-between cursor-pointer hover:border-brand-primary/50 transition-all ${isDropdownOpen ? 'ring-2 ring-brand-primary/50 border-brand-primary' : ''}`}
                            >
                                <School className="absolute left-4 top-1/2 -translate-y-1/2 h-5 w-5 text-gray-500" />
                                <span className={formData.institution ? 'text-white' : 'text-gray-500'}>
                                    {fetchingUniversities ? 'Loading universities...' : (formData.institution || 'Select your university')}
                                </span>
                                <ChevronDown className={`h-5 w-5 text-gray-500 transition-transform ${isDropdownOpen ? 'rotate-180' : ''}`} />
                            </div>

                            {isDropdownOpen && (
                                <motion.div
                                    initial={{ opacity: 0, y: 5 }}
                                    animate={{ opacity: 1, y: 0 }}
                                    className="absolute z-50 w-full mt-2 bg-[#1a1c20] border border-white/10 rounded-xl shadow-2xl overflow-hidden backdrop-blur-xl"
                                >
                                    <div className="p-3 border-b border-white/5">
                                        <div className="relative">
                                            <Search className="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-gray-500" />
                                            <input
                                                type="text"
                                                placeholder="Search university..."
                                                className="w-full pl-10 pr-4 py-2 bg-white/5 border border-white/10 rounded-lg text-sm focus:outline-none focus:border-brand-primary/50"
                                                value={searchTerm}
                                                onChange={(e) => setSearchTerm(e.target.value)}
                                                onClick={(e) => e.stopPropagation()}
                                            />
                                        </div>
                                    </div>
                                    <div className="max-h-60 overflow-y-auto">
                                        {filteredUniversities.length > 0 ? (
                                            filteredUniversities.map((u) => (
                                                <div
                                                    key={u.id}
                                                    onClick={() => {
                                                        setFormData({ ...formData, institution: u.name });
                                                        setIsDropdownOpen(false);
                                                        setSearchTerm('');
                                                    }}
                                                    className="px-4 py-3 hover:bg-white/5 cursor-pointer flex items-center justify-between group transition-colors"
                                                >
                                                    <span className="text-sm text-gray-300 group-hover:text-white">{u.name}</span>
                                                    {formData.institution === u.name && (
                                                        <Check className="h-4 w-4 text-brand-primary" />
                                                    )}
                                                </div>
                                            ))
                                        ) : (
                                            <div className="px-4 py-8 text-center text-sm text-gray-500">
                                                No universities found
                                            </div>
                                        )}
                                    </div>
                                </motion.div>
                            )}
                        </div>
                    </div>

                    <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                        <div>
                            <label className="block text-sm font-medium text-gray-400 mb-2">Password</label>
                            <div className="relative">
                                <Lock className="absolute left-4 top-1/2 -translate-y-1/2 h-5 w-5 text-gray-500" />
                                <input
                                    type="password"
                                    name="password"
                                    value={formData.password}
                                    onChange={handleChange}
                                    placeholder="••••••••"
                                    className="w-full pl-12 pr-4 py-3 bg-white/5 border border-white/10 rounded-xl focus:outline-none focus:ring-2 focus:ring-brand-primary/50 focus:border-brand-primary transition-all"
                                    required
                                    minLength={8}
                                />
                            </div>
                        </div>
                        <div>
                            <label className="block text-sm font-medium text-gray-400 mb-2">Confirm Password</label>
                            <div className="relative">
                                <Lock className="absolute left-4 top-1/2 -translate-y-1/2 h-5 w-5 text-gray-500" />
                                <input
                                    type="password"
                                    name="confirmPassword"
                                    value={formData.confirmPassword}
                                    onChange={handleChange}
                                    placeholder="••••••••"
                                    className="w-full pl-12 pr-4 py-3 bg-white/5 border border-white/10 rounded-xl focus:outline-none focus:ring-2 focus:ring-brand-primary/50 focus:border-brand-primary transition-all"
                                    required
                                />
                            </div>
                        </div>
                    </div>

                    <button
                        type="submit"
                        disabled={loading}
                        className="w-full py-4 mt-4 bg-brand-primary hover:bg-brand-primary/90 text-white rounded-xl font-display font-bold transition-all disabled:opacity-50 flex items-center justify-center gap-2"
                    >
                        {loading ? <Loader2 className="h-5 w-5 animate-spin" /> : 'Create Account'}
                    </button>
                </form>

                <div className="mt-8 text-center">
                    <Link
                        to="/login"
                        className="inline-flex items-center gap-2 text-sm text-gray-400 hover:text-white transition-colors"
                    >
                        <ArrowLeft className="h-4 w-4" />
                        Already have an account? Sign In
                    </Link>
                </div>
            </motion.div>
        </div>
    );
};

export default Register;
