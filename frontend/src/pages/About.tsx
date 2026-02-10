import React from 'react';
import { motion } from 'framer-motion';
import { Award, BookOpen, Globe, Lock, Workflow } from 'lucide-react';

const About: React.FC = () => {
    return (
        <div className="min-h-screen pt-32 pb-20 px-4">
            <div className="max-w-4xl mx-auto">
                <motion.div
                    initial={{ opacity: 0, y: 20 }}
                    animate={{ opacity: 1, y: 0 }}
                    className="text-center mb-16"
                >
                    <Award className="h-16 w-16 text-brand-secondary mx-auto mb-6" />
                    <h1 className="text-4xl md:text-5xl font-display font-bold mb-6">About the Project</h1>
                    <p className="text-xl text-gray-400">
                        BlockCertify is an academic proof-of-concept designed to solve the global problem of fraudulent diplomas using decentralized ledger technology.
                    </p>
                </motion.div>

                <div className="grid grid-cols-1 md:grid-cols-2 gap-8 mb-16">
                    <div className="p-8 rounded-3xl bg-white/5 border border-white/10">
                        <h3 className="text-xl font-display font-bold mb-4 flex items-center gap-2 text-brand-primary">
                            <BookOpen className="h-5 w-5" /> Academic Purpose
                        </h3>
                        <p className="text-gray-400 leading-relaxed">
                            Traditional verification processes are slow, manual, and prone to forgery. BlockCertify provides an instant, self-service verification layer that universities can integrate directly into their issuing workflow.
                        </p>
                    </div>
                    <div className="p-8 rounded-3xl bg-white/5 border border-white/10">
                        <h3 className="text-xl font-display font-bold mb-4 flex items-center gap-2 text-brand-secondary">
                            <Lock className="h-5 w-5" /> Security First
                        </h3>
                        <p className="text-gray-400 leading-relaxed">
                            By separating the storage of the document (Arweave) from the verification of its integrity (Polygon), we ensure that even if the storage layer is compromised, the proof of authenticity remains immutable.
                        </p>
                    </div>
                </div>

                <div className="p-10 rounded-3xl bg-gradient-to-br from-brand-primary/10 to-brand-accent/10 border border-white/10 relative overflow-hidden mb-16">
                    <Workflow className="h-32 w-32 text-white/5 absolute -bottom-8 -right-8" />
                    <h2 className="text-3xl font-display font-bold mb-6">Technological Choice</h2>
                    <div className="space-y-6 relative z-10">
                        <div>
                            <h4 className="font-bold text-white mb-2 underline decoration-brand-secondary">Polygon (Layer 2)</h4>
                            <p className="text-gray-400">
                                Chosen for its Ethereum-level security with significantly lower transaction costs and faster finality, making it scalable for large university systems.
                            </p>
                        </div>
                        <div>
                            <h4 className="font-bold text-white mb-2 underline decoration-brand-accent">Arweave (Permaweb)</h4>
                            <p className="text-gray-400">
                                A decentralized storage network that offers permanent, one-time payment storage. Unlike IPFS, Arweave guarantees that the diploma PDF will stay accessible for hundreds of years.
                            </p>
                        </div>
                    </div>
                </div>

                <div className="text-center">
                    <p className="text-sm text-gray-500 flex items-center justify-center gap-2">
                        <Globe className="h-4 w-4" /> Built for the Future of Decentralized Education.
                    </p>
                </div>
            </div>
        </div>
    );
};

export default About;
