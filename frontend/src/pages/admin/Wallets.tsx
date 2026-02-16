import React, { useState } from 'react';
import {
    Wallet,
    Key,
    FileJson,
    Shield,
    Eye,
    EyeOff,
    CheckCircle2,
    AlertCircle,
    Copy,
    ExternalLink
} from 'lucide-react';
import { motion } from 'framer-motion';
import api from '../../services/api';

const Wallets: React.FC = () => {
    const [arweaveWallet, setArweaveWallet] = useState<{ name: string; address: string; balance: string } | null>(null);
    const [polygonKey, setPolygonKey] = useState('');
    const [showPolygonKey, setShowPolygonKey] = useState(false);
    const [isSaving, setIsSaving] = useState(false);
    const [successMessage, setSuccessMessage] = useState('');
    const [isLoading, setIsLoading] = useState(false);

    const handleArweaveUpload = async (e: React.ChangeEvent<HTMLInputElement>) => {
        const file = e.target.files?.[0];
        if (file) {
            if (file.type !== 'application/json' && !file.name.endsWith('.json')) {
                alert('Please upload a valid JSON keyfile.');
                return;
            }

            setIsLoading(true);
            const formData = new FormData();
            formData.append('wallet', file);

            try {
                const response = await api.post('/v1/wallet/upload-key-file', formData, {
                    headers: {
                        'Content-Type': 'multipart/form-data',
                    },
                });

                if (response.data) {
                    setArweaveWallet({
                        name: file.name,
                        address: response.data.address,
                        balance: response.data.balance || '0'
                    });
                    setSuccessMessage('Arweave wallet connected successfully!');
                    setTimeout(() => setSuccessMessage(''), 3000);
                }
            } catch (error: any) {
                console.error('Failed to upload wallet:', error);
                alert(error.response?.data?.details || 'Failed to connect Arweave wallet.');
            } finally {
                setIsLoading(false);
            }
        }
    };

    const handleSave = () => {
        setIsSaving(true);
        // Simulate API call
        setTimeout(() => {
            setIsSaving(false);
            setSuccessMessage('Wallet configuration updated successfully!');
            setTimeout(() => setSuccessMessage(''), 3000);
        }, 1500);
    };

    return (
        <div className="space-y-8 pb-12">
            <div>
                <h1 className="text-3xl font-display font-bold text-white mb-2">Wallet Configuration</h1>
                <p className="text-gray-400">Configure your institution's cryptographic keys for Arweave and Polygon.</p>
            </div>

            {successMessage && (
                <motion.div
                    initial={{ opacity: 0, y: -20 }}
                    animate={{ opacity: 1, y: 0 }}
                    className="bg-green-500/10 border border-green-500/30 text-green-400 p-4 rounded-xl flex items-center gap-3"
                >
                    <CheckCircle2 className="h-5 w-5" />
                    {successMessage}
                </motion.div>
            )}

            <div className="grid grid-cols-1 md:grid-cols-2 gap-8">
                {/* Arweave Wallet Section */}
                <motion.div
                    initial={{ opacity: 0, x: -20 }}
                    animate={{ opacity: 1, x: 0 }}
                    className="bg-white/5 backdrop-blur-xl border border-white/10 rounded-3xl p-8 hover:border-brand-secondary/30 transition-all group"
                >
                    <div className="flex items-center gap-4 mb-8">
                        <div className="p-3 rounded-2xl bg-brand-secondary/10 text-brand-secondary group-hover:scale-110 transition-transform">
                            <FileJson className="h-6 w-6" />
                        </div>
                        <div>
                            <h2 className="text-xl font-bold text-white">Arweave Storage</h2>
                            <p className="text-sm text-gray-500">Permanent decentralized storage</p>
                        </div>
                    </div>

                    <div className="space-y-6">
                        <div className="p-6 border-2 border-dashed border-white/10 rounded-2xl hover:border-brand-secondary/50 transition-colors bg-white/5">
                            {!arweaveWallet ? (
                                <label className="flex flex-col items-center gap-3 cursor-pointer">
                                    <div className="p-3 rounded-full bg-white/5 text-gray-400">
                                        {isLoading ? (
                                            <motion.div
                                                animate={{ rotate: 360 }}
                                                transition={{ duration: 1, repeat: Infinity, ease: "linear" }}
                                            >
                                                <AlertCircle className="h-6 w-6 text-brand-secondary" />
                                            </motion.div>
                                        ) : (
                                            <Wallet className="h-6 w-6" />
                                        )}
                                    </div>
                                    <div className="text-center">
                                        <p className="text-white font-medium">{isLoading ? 'Connecting...' : 'Upload Keyfile'}</p>
                                        <p className="text-xs text-gray-500 mt-1">Select your .json keyfile</p>
                                    </div>
                                    <input
                                        type="file"
                                        className="hidden"
                                        accept=".json"
                                        onChange={handleArweaveUpload}
                                    />
                                </label>
                            ) : (
                                <div className="flex items-center justify-between">
                                    <div className="flex items-center gap-4">
                                        <div className="p-2 rounded-lg bg-green-500/20 text-green-400">
                                            <CheckCircle2 className="h-5 w-5" />
                                        </div>
                                        <div>
                                            <p className="text-white text-sm font-medium">{arweaveWallet.name}</p>
                                            <p className="text-xs text-brand-secondary font-bold mt-1">{arweaveWallet.balance} AR</p>
                                            <p className="text-xs text-gray-500 font-mono mt-0.5">{arweaveWallet.address}</p>
                                        </div>
                                    </div>
                                    <button
                                        onClick={() => setArweaveWallet(null)}
                                        className="text-xs text-gray-500 hover:text-red-400 transition-colors"
                                    >
                                        Change
                                    </button>
                                </div>
                            )}
                        </div>

                        <div className="bg-brand-secondary/5 rounded-2xl p-4 border border-brand-secondary/10 flex gap-3 text-xs text-brand-secondary">
                            <AlertCircle className="h-4 w-4 shrink-0" />
                            <p>Keyfiles are stored securely and used only for institutional signing of Arweave transactions.</p>
                        </div>
                    </div>
                </motion.div>

                {/* Polygon Wallet Section */}
                <motion.div
                    initial={{ opacity: 0, x: 20 }}
                    animate={{ opacity: 1, x: 0 }}
                    className="bg-white/5 backdrop-blur-xl border border-white/10 rounded-3xl p-8 hover:border-brand-primary/30 transition-all group"
                >
                    <div className="flex items-center gap-4 mb-8">
                        <div className="p-3 rounded-2xl bg-brand-primary/10 text-brand-primary group-hover:scale-110 transition-transform">
                            <Shield className="h-6 w-6" />
                        </div>
                        <div>
                            <h2 className="text-xl font-bold text-white">Polygon Verification</h2>
                            <p className="text-sm text-gray-500">Blockchain verification layer</p>
                        </div>
                    </div>

                    <div className="space-y-6">
                        <div className="space-y-2">
                            <label className="text-xs font-semibold text-gray-400 uppercase tracking-wider block ml-1">
                                Private Key
                            </label>
                            <div className="relative group/input">
                                <div className="absolute inset-y-0 left-0 pl-4 flex items-center pointer-events-none text-gray-500">
                                    <Key className="h-5 w-5" />
                                </div>
                                <input
                                    type={showPolygonKey ? 'text' : 'password'}
                                    value={polygonKey}
                                    onChange={(e) => setPolygonKey(e.target.value)}
                                    placeholder="Enter your Polygon private key"
                                    className="w-full bg-white/5 border border-white/10 rounded-2xl py-4 pl-12 pr-12 text-white placeholder:text-gray-600 focus:outline-none focus:border-brand-primary/50 focus:ring-1 focus:ring-brand-primary/50 transition-all"
                                />
                                <button
                                    type="button"
                                    onClick={() => setShowPolygonKey(!showPolygonKey)}
                                    className="absolute inset-y-0 right-0 pr-4 flex items-center text-gray-500 hover:text-white transition-colors"
                                >
                                    {showPolygonKey ? <EyeOff className="h-5 w-5" /> : <Eye className="h-5 w-5" />}
                                </button>
                            </div>
                        </div>

                        <div className="flex items-center justify-between text-xs p-4 bg-white/5 rounded-2xl border border-white/10 font-mono text-gray-400">
                            <span>0x7...E42</span>
                            <div className="flex gap-4">
                                <button className="hover:text-white flex items-center gap-1 transition-colors">
                                    <Copy className="h-3 w-3" /> Copy Address
                                </button>
                                <button className="hover:text-white flex items-center gap-1 transition-colors">
                                    <ExternalLink className="h-3 w-3" /> Explorer
                                </button>
                            </div>
                        </div>
                    </div>
                </motion.div>
            </div>

            <div className="flex justify-end p-6 bg-white/5 rounded-3xl border border-white/10">
                <button
                    onClick={handleSave}
                    disabled={isSaving}
                    className="bg-brand-primary hover:bg-brand-primary/90 text-white px-8 py-3 rounded-xl font-bold shadow-[0_0_20px_rgba(59,130,246,0.2)] hover:shadow-[0_0_25px_rgba(59,130,246,0.3)] transition-all flex items-center gap-3 disabled:opacity-50"
                >
                    {isSaving ? 'Processing...' : (
                        <>
                            <CheckCircle2 className="h-5 w-5" />
                            Save Configuration
                        </>
                    )}
                </button>
            </div>
        </div>
    );
};

export default Wallets;
