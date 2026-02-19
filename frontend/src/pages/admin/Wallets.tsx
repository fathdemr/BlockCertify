import React, { useState, useEffect } from 'react';
import {
    Wallet,
    FileJson,
    Shield,
    CheckCircle2,
    AlertCircle,
    Copy,
    ExternalLink
} from 'lucide-react';
import { ethers } from 'ethers';
import { motion } from 'framer-motion';
import api from '../../services/api';

declare global {
    interface Window {
        ethereum?: any;
    }
}

const Wallets: React.FC = () => {
    const [arweaveWallet, setArweaveWallet] = useState<{ name: string; address: string; balance: string } | null>(null);
    const [polygonAddress, setPolygonAddress] = useState<string | null>(null);
    const [isMetaMaskConnecting, setIsMetaMaskConnecting] = useState(false);
    const [successMessage, setSuccessMessage] = useState('');
    const [isLoading, setIsLoading] = useState(false);

    // Load persisted wallets on mount
    useEffect(() => {
        const savedArweave = localStorage.getItem('blockcertify_arweave_wallet');
        const savedPolygon = localStorage.getItem('blockcertify_polygon_address');

        if (savedArweave) {
            try {
                setArweaveWallet(JSON.parse(savedArweave));
            } catch (e) {
                console.error('Failed to parse saved Arweave wallet', e);
                localStorage.removeItem('blockcertify_arweave_wallet');
            }
        }

        if (savedPolygon) {
            setPolygonAddress(savedPolygon);
        }
    }, []);

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
                    const walletData = {
                        name: file.name,
                        address: response.data.address,
                        balance: response.data.balance || '0'
                    };
                    setArweaveWallet(walletData);
                    localStorage.setItem('blockcertify_arweave_wallet', JSON.stringify(walletData));
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

    const connectMetaMask = async () => {
        if (!window.ethereum) {
            alert('MetaMask is not installed.');
            return;
        }

        setIsMetaMaskConnecting(true);
        try {
            const provider = new ethers.BrowserProvider(window.ethereum);

            // Request accounts
            await provider.send('eth_requestAccounts', []);
            const signer = await provider.getSigner();
            const address = await signer.getAddress();

            // Ensure user is on Polygon Amoy (chainId 80002)
            const network = await provider.getNetwork();
            if (network.chainId !== 80002n) {
                try {
                    await window.ethereum.request({
                        method: 'wallet_switchEthereumChain',
                        params: [{ chainId: '0x13882' }], // 80002 in hex
                    });
                } catch (switchError: any) {
                    // Chain not added yet — add it
                    if (switchError.code === 4902) {
                        await window.ethereum.request({
                            method: 'wallet_addEthereumChain',
                            params: [{
                                chainId: '0x13882',
                                chainName: 'Polygon Amoy Testnet',
                                nativeCurrency: { name: 'MATIC', symbol: 'MATIC', decimals: 18 },
                                rpcUrls: ['https://rpc-amoy.polygon.technology/'],
                                blockExplorerUrls: ['https://amoy.polygonscan.com/'],
                            }],
                        });
                    } else {
                        throw switchError;
                    }
                }
            }

            setPolygonAddress(address);
            localStorage.setItem('blockcertify_polygon_address', address);
            setSuccessMessage('Polygon wallet connected!');
            setTimeout(() => setSuccessMessage(''), 3000);
        } catch (error: any) {
            if (error.code === 4001) {
                alert('Connection rejected.');
            } else {
                alert('Failed to connect MetaMask.');
            }
        } finally {
            setIsMetaMaskConnecting(false);
        }
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
                                        onClick={() => {
                                            setArweaveWallet(null);
                                            localStorage.removeItem('blockcertify_arweave_wallet');
                                        }}
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
                        {!polygonAddress ? (
                            <div className="space-y-4">
                                <div className="p-6 border-2 border-dashed border-white/10 rounded-2xl bg-white/5 text-center">
                                    <div className="p-3 rounded-full bg-white/5 text-gray-400 w-fit mx-auto mb-3">
                                        <Wallet className="h-6 w-6" />
                                    </div>
                                    <p className="text-white font-medium mb-1">MetaMask Connection</p>
                                    <p className="text-xs text-gray-500">Connect your wallet for on-chain verification</p>
                                </div>
                                <button
                                    onClick={connectMetaMask}
                                    disabled={isMetaMaskConnecting}
                                    className="w-full bg-brand-primary hover:bg-brand-primary/90 text-white py-4 rounded-2xl font-bold transition-all flex items-center justify-center gap-3 shadow-[0_0_20px_rgba(59,130,246,0.1)]"
                                >
                                    {isMetaMaskConnecting ? (
                                        <>
                                            <motion.div
                                                animate={{ rotate: 360 }}
                                                transition={{ duration: 1, repeat: Infinity, ease: "linear" }}
                                            >
                                                <AlertCircle className="h-5 w-5" />
                                            </motion.div>
                                            Connecting...
                                        </>
                                    ) : (
                                        <>
                                            <Wallet className="h-5 w-5" />
                                            Connect MetaMask
                                        </>
                                    )}
                                </button>
                            </div>
                        ) : (
                            <div className="space-y-4">
                                <div className="p-6 border border-brand-primary/30 rounded-2xl bg-brand-primary/5 flex items-center justify-between">
                                    <div className="flex items-center gap-4">
                                        <div className="p-2 rounded-lg bg-brand-primary/20 text-brand-primary">
                                            <CheckCircle2 className="h-5 w-5" />
                                        </div>
                                        <div>
                                            <p className="text-white text-sm font-medium">Verified Wallet</p>
                                            <p className="text-xs text-gray-500 font-mono mt-1">
                                                {polygonAddress.slice(0, 6)}...{polygonAddress.slice(-4)}
                                            </p>
                                        </div>
                                    </div>
                                    <button
                                        onClick={() => {
                                            setPolygonAddress(null);
                                            localStorage.removeItem('blockcertify_polygon_address');
                                        }}
                                        className="text-xs text-gray-500 hover:text-red-400 transition-colors"
                                    >
                                        Disconnect
                                    </button>
                                </div>
                                <div className="flex items-center justify-between text-xs p-4 bg-white/5 rounded-2xl border border-white/10 font-mono text-gray-400">
                                    <div className="flex gap-4 w-full justify-between">
                                        <button
                                            onClick={() => {
                                                navigator.clipboard.writeText(polygonAddress);
                                                alert('Address copied to clipboard!');
                                            }}
                                            className="hover:text-white flex items-center gap-1 transition-colors"
                                        >
                                            <Copy className="h-3 w-3" /> Copy Address
                                        </button>
                                        <a
                                            href={`https://amoy.polygonscan.com/address/${polygonAddress}`}
                                            target="_blank"
                                            rel="noopener noreferrer"
                                            className="hover:text-white flex items-center gap-1 transition-colors"
                                        >
                                            <ExternalLink className="h-3 w-3" /> Explorer
                                        </a>
                                    </div>
                                </div>
                            </div>
                        )}

                        <div className="bg-brand-primary/5 rounded-2xl p-4 border border-brand-primary/10 flex gap-3 text-xs text-brand-primary">
                            <Shield className="h-4 w-4 shrink-0" />
                            <p>Connecting with MetaMask ensures that only authorized administrators can sign verification transactions.</p>
                        </div>
                    </div>
                </motion.div>
            </div>

        </div>
    );
};

export default Wallets;
