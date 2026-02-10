import React from 'react';
import {
    Users,
    FileCheck,
    AlertCircle,
    Clock,
    ArrowUpRight,
    Shield,
    Activity
} from 'lucide-react';
import { motion } from 'framer-motion';

const Dashboard: React.FC = () => {
    const stats = [
        { label: 'Total Diplomas Issued', value: '1,284', icon: FileCheck, color: 'text-brand-primary', trend: '+12%' },
        { label: 'Pending Verifications', value: '0', icon: Clock, color: 'text-brand-secondary', trend: '0%' },
        { label: 'Active Students', value: '850', icon: Users, color: 'text-brand-accent', trend: '+5%' },
        { label: 'System Health', value: '100%', icon: Activity, color: 'text-brand-success', trend: 'Stable' },
    ];

    const recentHistory = [
        { id: 'BC-2025-441', name: 'John Doe', status: 'Success', date: '2 mins ago', hash: '0x742d...44e' },
        { id: 'BC-2025-440', name: 'Jane Smith', status: 'Success', date: '1 hour ago', hash: '0x815f...12a' },
        { id: 'BC-2025-439', name: 'Robert Brown', status: 'Success', date: '4 hours ago', hash: '0x921c...55d' },
        { id: 'BC-2025-438', name: 'Emily Davis', status: 'Success', date: 'Yesterday', hash: '0x012b...b3c' },
    ];

    return (
        <div className="space-y-10">
            <div className="flex flex-col md:flex-row md:items-center justify-between gap-4">
                <div>
                    <h1 className="text-4xl font-display font-bold text-white mb-2">Dashboard Overview</h1>
                    <p className="text-gray-400">Welcome back, Karabuk University Administrator.</p>
                </div>
                <div className="flex items-center gap-3 px-4 py-2 bg-brand-primary/10 border border-brand-primary/20 rounded-xl">
                    <Shield className="h-5 w-5 text-brand-primary" />
                    <span className="text-sm font-semibold text-brand-primary uppercase tracking-wider">Institution Verified</span>
                </div>
            </div>

            {/* Stats Grid */}
            <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-6">
                {stats.map((stat, i) => (
                    <motion.div
                        key={i}
                        initial={{ opacity: 0, y: 20 }}
                        animate={{ opacity: 1, y: 0 }}
                        transition={{ delay: i * 0.1 }}
                        className="p-6 rounded-3xl bg-white/5 border border-white/10 hover:border-white/20 transition-all group"
                    >
                        <div className="flex items-start justify-between mb-4">
                            <div className={`p-3 rounded-2xl bg-white/5 group-hover:scale-110 transition-transform ${stat.color}`}>
                                <stat.icon className="h-6 w-6" />
                            </div>
                            <span className={`text-xs font-bold px-2 py-1 rounded-lg bg-white/5 ${stat.trend.startsWith('+') ? 'text-brand-success' : 'text-gray-400'
                                }`}>
                                {stat.trend}
                            </span>
                        </div>
                        <p className="text-gray-400 text-sm font-medium mb-1">{stat.label}</p>
                        <h3 className="text-2xl font-display font-bold text-white">{stat.value}</h3>
                    </motion.div>
                ))}
            </div>

            <div className="grid grid-cols-1 lg:grid-cols-3 gap-8">
                {/* Recent Activity */}
                <div className="lg:col-span-2 space-y-6">
                    <div className="flex items-center justify-between">
                        <h2 className="text-2xl font-display font-bold">Recent Issuance</h2>
                        <button className="text-sm text-brand-secondary hover:underline flex items-center gap-1">
                            View all <ArrowUpRight className="h-4 w-4" />
                        </button>
                    </div>
                    <div className="bg-white/5 border border-white/10 rounded-3xl overflow-hidden">
                        <table className="w-full text-left">
                            <thead>
                                <tr className="border-b border-white/10">
                                    <th className="px-6 py-4 text-xs font-semibold text-gray-500 uppercase tracking-wider">ID</th>
                                    <th className="px-6 py-4 text-xs font-semibold text-gray-500 uppercase tracking-wider">Student</th>
                                    <th className="px-6 py-4 text-xs font-semibold text-gray-500 uppercase tracking-wider">Date</th>
                                    <th className="px-6 py-4 text-xs font-semibold text-gray-500 uppercase tracking-wider">Status</th>
                                </tr>
                            </thead>
                            <tbody className="divide-y divide-white/5">
                                {recentHistory.map((item, i) => (
                                    <tr key={i} className="hover:bg-white/5 transition-colors cursor-pointer group">
                                        <td className="px-6 py-4">
                                            <span className="font-mono text-xs text-brand-secondary">{item.id}</span>
                                        </td>
                                        <td className="px-6 py-4">
                                            <p className="text-sm font-semibold text-white">{item.name}</p>
                                            <p className="text-xs text-gray-500 font-mono truncate max-w-[100px]">{item.hash}</p>
                                        </td>
                                        <td className="px-6 py-4 text-sm text-gray-400">{item.date}</td>
                                        <td className="px-6 py-4">
                                            <span className="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-brand-success/10 text-brand-success border border-brand-success/20">
                                                {item.status}
                                            </span>
                                        </td>
                                    </tr>
                                ))}
                            </tbody>
                        </table>
                    </div>
                </div>

                {/* Network Status */}
                <div className="space-y-6">
                    <h2 className="text-2xl font-display font-bold">Network Nodes</h2>
                    <div className="space-y-4">
                        {[
                            { name: 'Polygon Mainnet', status: 'Connected', ping: '12ms' },
                            { name: 'Arweave Gateway', status: 'Connected', ping: '45ms' },
                            { name: 'Internal Hashing Engine', status: 'Online', ping: '0.2ms' },
                        ].map((node, i) => (
                            <div key={i} className="p-4 rounded-2xl bg-white/5 border border-white/10 flex items-center justify-between">
                                <div>
                                    <p className="text-sm font-semibold text-white">{node.name}</p>
                                    <div className="flex items-center gap-1.5 mt-1">
                                        <div className="w-2 h-2 rounded-full bg-brand-success" />
                                        <span className="text-xs text-gray-500">{node.status}</span>
                                    </div>
                                </div>
                                <span className="text-xs font-mono text-gray-500">{node.ping}</span>
                            </div>
                        ))}

                        <div className="p-6 rounded-3xl bg-brand-primary/5 border border-brand-primary/20 mt-4">
                            <div className="flex items-center gap-2 text-brand-primary mb-2">
                                <AlertCircle className="h-5 w-5" />
                                <span className="font-bold">Security Tip</span>
                            </div>
                            <p className="text-sm text-gray-400">
                                Regularly verify your institution's private key rotation schedule to maintain high integrity standards.
                            </p>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    );
};

export default Dashboard;
