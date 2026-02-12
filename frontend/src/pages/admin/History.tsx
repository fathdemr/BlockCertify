import React, { useEffect, useState } from 'react';
import {
    ExternalLink,
    Search,
    Filter,
    Download
} from 'lucide-react';
import { diplomaService } from '../../services/api';


const History: React.FC = () => {
    const [logs, setLogs] = useState<any[]>([]);
    const [loading, setLoading] = useState(true);

    useEffect(() => {
        const fetchDiplomas = async () => {
            try {
                const data = await diplomaService.getAllDiplomas();
                setLogs(data);
            } catch (err) {
                console.error("Failed to fetch diplomas", err);
            } finally {
                setLoading(false);
            }
        };

        fetchDiplomas();
    }, []);

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
                            <th className="px-6 py-4 text-xs font-bold text-gray-500 uppercase tracking-widest">Diploma ID</th>
                            <th className="px-6 py-4 text-xs font-bold text-gray-500 uppercase tracking-widest">Owner</th>
                            <th className="px-6 py-4 text-xs font-bold text-gray-500 uppercase tracking-widest">Department</th>
                            <th className="px-6 py-4 text-xs font-bold text-gray-500 uppercase tracking-widest">Registered Date</th>
                            <th className="px-6 py-4 text-xs font-bold text-gray-500 uppercase tracking-widest">Diploma PDF</th>
                        </tr>
                    </thead>
                    <tbody className="divide-y divide-white/5 text-sm">
                        {loading ? (
                            <tr>
                                <td colSpan={5} className="text-center py-10 text-gray-400">
                                    Loading diplomas...
                                </td>
                            </tr>
                        ) : logs.map((log, i) => (
                            <tr key={i} className="hover:bg-white/5 transition-colors group">
                                <td className="px-6 py-4 font-mono text-brand-secondary">{log.diplomaId}</td>
                                <td className="px-6 py-4 text-white font-medium">{log.userName}</td>
                                <td className="px-6 py-4">
                                    {log.department}
                                </td>
                                <td className="px-6 py-4 text-gray-400">
                                    {new Date(log.createDate).toLocaleString()}
                                </td>
                                <td className="px-6 py-4">
                                    <button
                                        onClick={() => diplomaService.getDiplomaFile(log.diplomaId)}
                                        className={"p-2 bg-white/5 rounded-lg border border-white/5 opacity-80 hover:opacity-100 transition-opacity inline-flex"}
                                    >
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
