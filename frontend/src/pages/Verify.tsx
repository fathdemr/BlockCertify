import React, { useState } from 'react';
import { motion, AnimatePresence } from 'framer-motion';
import { Search, Loader2, CheckCircle2, XCircle, Shield, FileText, Globe, Clock, Hash, Database } from 'lucide-react';
import { diplomaService } from '../services/api';

interface VerificationResult {
    status: 'success' | 'failed' | 'idle';
    data?: {
        id: string;
        studentName: string;
        university: string;
        degree: string;
        issueDate: string;
        polygonTx: string;
        arweaveUrl: string;
        arweaveTxID: string;
        fileHash: string;
    };
}

const Verify: React.FC = () => {
    const [diplomaId, setDiplomaId] = useState('');
    const [loading, setLoading] = useState(false);
    const [result, setResult] = useState<VerificationResult>({ status: 'idle' });

    const handleVerify = async (e: React.FormEvent) => {
        e.preventDefault();
        if (!diplomaId.trim()) return;

        setLoading(true);
        setResult({ status: 'idle' });

        try {
            const data = await diplomaService.verify(diplomaId);
            if (data.verified) {
                setResult({
                    status: 'success',
                    data: {
                        id: data.diplomaID || diplomaId.toUpperCase(),
                        studentName: data.studentName || 'N/A',
                        university: data.university || 'N/A',
                        degree: data.degree || 'N/A',
                        issueDate: data.issueDate || 'N/A',
                        polygonTx: data.polygonTxHash || 'N/A',
                        arweaveUrl: data.arweaveUrl || 'N/A',
                        arweaveTxID: data.arweaveTxID || 'N/A',
                        fileHash: data.diplomaHash || 'N/A',
                    }
                });
            } else {
                setResult({ status: 'failed' });
            }
        } catch (error) {
            console.error('Verification failed:', error);
            setResult({ status: 'failed' });
        } finally {
            setLoading(false);
        }
    };

    return (
        <div className="min-h-screen pt-32 pb-20 px-4">
            <div className="max-w-4xl mx-auto">
                <div className="text-center mb-12">
                    <motion.h1
                        initial={{ opacity: 0, y: -20 }}
                        animate={{ opacity: 1, y: 0 }}
                        className="text-4xl md:text-5xl font-display font-bold mb-4"
                    >
                        Verify <span className="text-brand-secondary">Diploma</span>
                    </motion.h1>
                    <p className="text-gray-400">Enter the unique Diploma ID to verify its authenticity on the blockchain.</p>
                </div>

                {/* Search Bar */}
                <div className="mb-12">
                    <form onSubmit={handleVerify} className="relative group">
                        <div className="absolute inset-y-0 left-0 pl-6 flex items-center pointer-events-none text-gray-500 group-focus-within:text-brand-secondary transition-colors">
                            <Search className="h-6 w-6" />
                        </div>
                        <input
                            type="text"
                            value={diplomaId}
                            onChange={(e) => setDiplomaId(e.target.value)}
                            placeholder="e.g. BC-2025-XXXX"
                            className="block w-full pl-16 pr-6 py-5 bg-white/5 border border-white/10 rounded-2xl text-xl text-white placeholder-gray-600 focus:outline-none focus:ring-2 focus:ring-brand-secondary/50 focus:border-brand-secondary transition-all backdrop-blur-sm"
                        />
                        <button
                            type="submit"
                            disabled={loading || !diplomaId.trim()}
                            className="absolute inset-y-2 right-2 px-8 bg-brand-primary hover:bg-brand-primary/90 text-white rounded-xl font-display font-bold transition-all disabled:opacity-50 disabled:cursor-not-allowed shadow-lg"
                        >
                            {loading ? <Loader2 className="h-5 w-5 animate-spin" /> : 'Verify'}
                        </button>
                    </form>
                </div>

                <AnimatePresence mode="wait">
                    {loading && (
                        <motion.div
                            key="loading"
                            initial={{ opacity: 0 }}
                            animate={{ opacity: 1 }}
                            exit={{ opacity: 0 }}
                            className="flex flex-col items-center justify-center p-20 text-center"
                        >
                            <div className="relative mb-6">
                                <Shield className="h-16 w-16 text-brand-primary animate-pulse" />
                                <div className="absolute inset-0 bg-brand-primary/20 rounded-full blur-xl animate-pulse" />
                            </div>
                            <h3 className="text-xl font-display font-bold mb-2">Verifying Integrity</h3>
                            <p className="text-gray-400">Querying Polygon and Arweave nodes...</p>
                        </motion.div>
                    )}

                    {result.status === 'success' && result.data && (
                        <motion.div
                            key="success"
                            initial={{ opacity: 0, scale: 0.95 }}
                            animate={{ opacity: 1, scale: 1 }}
                            className="bg-brand-success/5 border border-brand-success/20 rounded-3xl overflow-hidden shadow-[0_0_40px_rgba(16,185,129,0.1)]"
                        >
                            <div className="bg-brand-success/10 px-8 py-6 flex items-center gap-4">
                                <CheckCircle2 className="h-8 w-8 text-brand-success" />
                                <div>
                                    <h3 className="text-xl font-display font-bold text-brand-success">Verification Successful</h3>
                                    <p className="text-brand-success/70 text-sm">Authenticity confirmed by Blockchain.</p>
                                </div>
                            </div>

                            <div className="p-8">
                                <div className="grid grid-cols-1 md:grid-cols-2 gap-8 mb-8">
                                    <div className="space-y-4">
                                        <div className="flex items-center gap-3">
                                            <FileText className="h-5 w-5 text-gray-500" />
                                            <div>
                                                <p className="text-xs text-gray-500 uppercase tracking-wider">Student Name</p>
                                                <p className="text-lg font-semibold">{result.data.studentName}</p>
                                            </div>
                                        </div>
                                        <div className="flex items-center gap-3">
                                            <Globe className="h-5 w-5 text-gray-500" />
                                            <div>
                                                <p className="text-xs text-gray-500 uppercase tracking-wider">University</p>
                                                <p className="text-lg font-semibold">{result.data.university}</p>
                                            </div>
                                        </div>
                                    </div>
                                    <div className="space-y-4">
                                        <div className="flex items-center gap-3">
                                            <Shield className="h-5 w-5 text-gray-500" />
                                            <div>
                                                <p className="text-xs text-gray-500 uppercase tracking-wider">Degree</p>
                                                <p className="text-lg font-semibold">{result.data.degree}</p>
                                            </div>
                                        </div>
                                        <div className="flex items-center gap-3">
                                            <Clock className="h-5 w-5 text-gray-500" />
                                            <div>
                                                <p className="text-xs text-gray-500 uppercase tracking-wider">Issue Date</p>
                                                <p className="text-lg font-semibold">{result.data.issueDate}</p>
                                            </div>
                                        </div>
                                    </div>
                                </div>

                                <div className="border-t border-white/5 pt-8 space-y-4">
                                    <div className="p-6 bg-white/5 rounded-2xl border border-white/10 flex flex-col gap-3">
                                        <div className="flex items-center gap-3 text-gray-400">
                                            <Hash className="h-4 w-4 text-brand-accent" />
                                            <span className="text-sm font-medium uppercase tracking-wider">Polygon Transaction</span>
                                        </div>
                                        <div className="bg-black/20 p-3 rounded-lg border border-white/5 font-mono text-xs text-brand-accent break-all leading-relaxed">
                                            {result.data.polygonTx}
                                        </div>
                                    </div>
                                        <div className="p-6 bg-white/5 rounded-2xl border border-white/10 flex flex-col gap-3">
                                            <div className="flex items-center gap-3 text-gray-400">
                                                <Database className="h-4 w-4 text-brand-secondary" />
                                                <span className="text-sm font-medium uppercase tracking-wider">Arweave Data Link</span>
                                            </div>
                                            <div className="bg-black/20 p-3 rounded-lg border border-white/5 font-mono text-xs text-brand-secondary break-all leading-relaxed">
                                                {result.data.arweaveUrl}
                                            </div>
                                    </div>
                                </div>
                            </div>
                        </motion.div>
                    )}

                    {result.status === 'failed' && (
                        <motion.div
                            key="failed"
                            initial={{ opacity: 0, scale: 0.95 }}
                            animate={{ opacity: 1, scale: 1 }}
                            className="bg-red-500/5 border border-red-500/20 rounded-3xl p-12 text-center"
                        >
                            <XCircle className="h-16 w-16 text-red-500 mx-auto mb-6" />
                            <h3 className="text-2xl font-display font-bold text-red-500 mb-2">Invalid Certificate</h3>
                            <p className="text-gray-400 mb-8 max-w-md mx-auto">
                                The diploma ID provided does not match any blockchain record or is not authentic.
                            </p>
                            <button
                                onClick={() => setResult({ status: 'idle' })}
                                className="text-white/70 hover:text-white underline text-sm transition-colors"
                            >
                                Try Another Search
                            </button>
                        </motion.div>
                    )}

                    {result.status === 'idle' && (
                        <motion.div
                            key="educational"
                            initial={{ opacity: 0 }}
                            animate={{ opacity: 1 }}
                            className="mt-12 grid grid-cols-1 md:grid-cols-2 gap-8"
                        >
                            <div className="p-8 rounded-3xl bg-white/5 border border-white/10">
                                <Shield className="h-8 w-8 text-brand-primary mb-4" />
                                <h3 className="text-xl font-display font-bold mb-3">Why Blockchain?</h3>
                                <p className="text-gray-400 text-sm leading-relaxed">
                                    Blockchain provides an immutable record of the diploma's fingerprint (hash). Even if the paper certificate is lost or the digital file is altered, the blockchain proof remains unchanged and verifiable.
                                </p>
                            </div>
                            <div className="p-8 rounded-3xl bg-white/5 border border-white/10">
                                <Database className="h-8 w-8 text-brand-secondary mb-4" />
                                <h3 className="text-xl font-display font-bold mb-3">Permanent Storage</h3>
                                <p className="text-gray-400 text-sm leading-relaxed">
                                    By using Arweave's permaweb, we ensure that digital diplomas are stored forever. Traditional servers can go down, but decentralized storage guarantees access for decades to come.
                                </p>
                            </div>
                        </motion.div>
                    )}
                </AnimatePresence>
            </div>
        </div>
    );
};

export default Verify;
