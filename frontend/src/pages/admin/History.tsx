import React from 'react';
import {
    ExternalLink,
    Search,
    Filter,
    Download
} from 'lucide-react';

const History: React.FC = () => {
    const logs = [
        { id: 'BC-2025-441', name: 'John Doe', type: 'Issuance', status: 'Success', date: '2026-02-02 23:10', hash: '0x742d...44e' },
        { id: 'BC-2025-440', name: 'Jane Smith', type: 'Issuance', status: 'Success', date: '2026-02-02 22:15', hash: '0x815f...12a' },
        { id: 'BC-1234567', name: 'Verification Unit', type: 'Self-Verify', status: 'Valid', date: '2026-02-02 21:05', hash: 'N/A' },
        { id: 'BC-2025-439', name: 'Robert Brown', type: 'Issuance', status: 'Success', date: '2026-02-02 18:42', hash: '0x921c...55d' },
        { id: 'BC-2025-438', name: 'Emily Davis', type: 'Issuance', status: 'Success', date: '2026-02-01 14:20', hash: '0x012b...b3c' },
        { id: 'BC-2025-437', name: 'Michael Wilson', type: 'Issuance', status: 'Success', date: '2026-02-01 10:15', hash: '0x334a...d4e' },
    ];

    return (
        <div className="space-y-8">
            <div className="flex flex-col md:flex-row md:items-center justify-between gap-4">
                <div>
                    <h1 className="text-3xl font-display font-bold mb-2">History & Logs</h1>
                    <p className="text-gray-400">Full audit trail of all blockchain interactions and issued credentials.</p>
                </div>
                <button className="flex items-center gap-2 px-4 py-2 bg-white/5 hover:bg-white/10 text-white border border-white/10 rounded-xl text-sm font-semibold transition-all">
                    <Download className="h-4 w-4" /> Export CSV
                </button>
            </div>

            <div className="flex flex-col md:flex-row gap-4">
                <div className="flex-grow relative">
                    <Search className="absolute left-4 top-1/2 -translate-y-1/2 h-5 w-5 text-gray-500" />
                    <input
                        type="text"
                        placeholder="Search logs..."
                        className="w-full pl-12 pr-4 py-3 bg-white/5 border border-white/10 rounded-xl focus:ring-2 focus:ring-brand-primary/50 outline-none"
                    />
                </div>
                <button className="flex items-center gap-2 px-6 py-3 bg-white/5 border border-white/10 rounded-xl text-gray-400 hover:text-white transition-all">
                    <Filter className="h-4 w-4" /> Filters
                </button>
            </div>

            <div className="bg-white/5 border border-white/10 rounded-3xl overflow-hidden">
                <table className="w-full text-left font-display">
                    <thead>
                        <tr className="bg-white/5 border-b border-white/10">
                            <th className="px-6 py-4 text-xs font-bold text-gray-500 uppercase tracking-widest">Entry ID</th>
                            <th className="px-6 py-4 text-xs font-bold text-gray-500 uppercase tracking-widest">Description</th>
                            <th className="px-6 py-4 text-xs font-bold text-gray-500 uppercase tracking-widest">Type</th>
                            <th className="px-6 py-4 text-xs font-bold text-gray-500 uppercase tracking-widest">Date</th>
                            <th className="px-6 py-4 text-xs font-bold text-gray-500 uppercase tracking-widest">TX Proof</th>
                        </tr>
                    </thead>
                    <tbody className="divide-y divide-white/5 text-sm">
                        {logs.map((log, i) => (
                            <tr key={i} className="hover:bg-white/5 transition-colors group">
                                <td className="px-6 py-4 font-mono text-brand-secondary">{log.id}</td>
                                <td className="px-6 py-4 text-white font-medium">{log.name}</td>
                                <td className="px-6 py-4">
                                    <span className={`px-2 py-1 rounded-md text-[10px] font-bold uppercase ${log.type === 'Issuance' ? 'bg-brand-primary/10 text-brand-primary' : 'bg-brand-accent/10 text-brand-accent'
                                        }`}>
                                        {log.type}
                                    </span>
                                </td>
                                <td className="px-6 py-4 text-gray-400">{log.date}</td>
                                <td className="px-6 py-4">
                                    <button className="p-2 bg-white/5 rounded-lg border border-white/5 opacity-0 group-hover:opacity-100 transition-opacity">
                                        <ExternalLink className="h-4 w-4 text-gray-500 hover:text-white" />
                                    </button>
                                </td>
                            </tr>
                        ))}
                    </tbody>
                </table>
            </div>
        </div>
    );
};

export default History;
