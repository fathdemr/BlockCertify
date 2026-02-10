import React from 'react';
import { motion } from 'framer-motion';
import { ArrowRight, Hexagon, Database, Cpu, ShieldCheck } from 'lucide-react';
import { Link } from 'react-router-dom';

const Home: React.FC = () => {
    return (
        <div className="flex flex-col w-full">
            {/* Hero Section */}
            <section className="relative pt-32 pb-20 md:pt-48 md:pb-32 overflow-hidden px-4 md:px-0">
                <div className="absolute top-0 left-1/2 -translate-x-1/2 w-full h-full -z-10 opacity-20">
                    <div className="absolute top-1/4 left-1/4 w-96 h-96 bg-brand-primary rounded-full blur-[128px]" />
                    <div className="absolute bottom-1/4 right-1/4 w-96 h-96 bg-brand-accent rounded-full blur-[128px]" />
                </div>

                <div className="max-w-7xl mx-auto text-center relative z-10">
                    <motion.div
                        initial={{ opacity: 0, y: 20 }}
                        animate={{ opacity: 1, y: 0 }}
                        transition={{ duration: 0.6 }}
                    >
                        <span className="inline-flex items-center px-4 py-1 rounded-full text-xs font-semibold tracking-wider bg-brand-secondary/10 text-brand-secondary border border-brand-secondary/20 mb-6 font-display">
                            ACADEMIC INTEGRITY ON-CHAIN
                        </span>
                        <h1 className="text-5xl md:text-7xl font-display font-bold leading-tight mb-8">
                            Secure <span className="text-secondary-cyan bg-gradient-to-r from-brand-secondary to-brand-primary bg-clip-text text-transparent">Diploma Verification</span> <br /> on Blockchain
                        </h1>
                        <p className="text-xl text-gray-400 max-w-2xl mx-auto mb-10">
                            BlockCertify uses Arweave for permanent storage and Polygon for tamper-proof verification, ensuring every certificate is authentic and verifiable forever.
                        </p>
                        <div className="flex flex-col sm:flex-row items-center justify-center gap-4">
                            <Link
                                to="/verify"
                                className="group w-full sm:w-auto px-8 py-4 bg-brand-primary hover:bg-brand-primary/90 text-white rounded-xl font-display font-bold text-lg flex items-center justify-center gap-2 transition-all shadow-[0_0_20px_rgba(59,130,246,0.4)]"
                            >
                                Verify Certificate
                                <ArrowRight className="h-5 w-5 group-hover:translate-x-1 transition-transform" />
                            </Link>
                            <a
                                href="#how-it-works"
                                className="w-full sm:w-auto px-8 py-4 bg-white/5 hover:bg-white/10 text-white border border-white/10 rounded-xl font-display font-bold text-lg transition-all"
                            >
                                How It Works
                            </a>
                        </div>
                    </motion.div>
                </div>
            </section>

            {/* Stats / Proof Section */}
            <section className="py-20 border-y border-gray-800 bg-brand-dark/50 backdrop-blur-sm">
                <div className="max-w-7xl mx-auto px-4 grid grid-cols-1 md:grid-cols-3 gap-12">
                    {[
                        { label: 'Network', value: 'Polygon PoS', icon: Cpu, color: 'text-brand-accent' },
                        { label: 'Storage', value: 'Arweave Permaweb', icon: Database, color: 'text-brand-secondary' },
                        { label: 'Security', value: 'Hash Integrity', icon: ShieldCheck, color: 'text-brand-success' },
                    ].map((stat, i) => (
                        <motion.div
                            key={i}
                            whileHover={{ y: -5 }}
                            className="flex flex-col items-center text-center p-6 rounded-2xl bg-white/5 border border-white/10"
                        >
                            <stat.icon className={`h-10 w-10 mb-4 ${stat.color}`} />
                            <h3 className="text-sm font-semibold text-gray-400 uppercase tracking-widest">{stat.label}</h3>
                            <p className="text-2xl font-display font-bold text-white">{stat.value}</p>
                        </motion.div>
                    ))}
                </div>
            </section>

            {/* How It Works Section */}
            <section id="how-it-works" className="py-24 px-4 overflow-hidden">
                <div className="max-w-7xl mx-auto">
                    <div className="text-center mb-16">
                        <h2 className="text-4xl font-display font-bold mb-4">How It Works</h2>
                        <p className="text-gray-400">The architecture behind the decentralized verification system</p>
                    </div>

                    <div className="relative">
                        {/* Connection Line (Desktop) */}
                        <div className="hidden md:block absolute top-[120px] left-0 w-full h-1 bg-gradient-to-r from-brand-primary via-brand-secondary to-brand-accent opacity-20" />

                        <div className="grid grid-cols-1 md:grid-cols-4 gap-8">
                            {[
                                {
                                    title: 'Upload',
                                    desc: 'University uploads diploma PDF to the secure portal.',
                                    icon: <Hexagon className="h-10 w-10 text-brand-primary" />,
                                    step: '01'
                                },
                                {
                                    title: 'Hashing',
                                    desc: 'System generates a unique SHA-256 fingerprint of the file.',
                                    icon: <Cpu className="h-10 w-10 text-brand-secondary" />,
                                    step: '02'
                                },
                                {
                                    title: 'Store',
                                    desc: 'PDF is encrypted and stored permanently on Arweave.',
                                    icon: <Database className="h-10 w-10 text-brand-accent" />,
                                    step: '03'
                                },
                                {
                                    title: 'Register',
                                    desc: 'Hash is recorded on Polygon for immutable verification.',
                                    icon: <ShieldCheck className="h-10 w-10 text-brand-success" />,
                                    step: '04'
                                },
                            ].map((item, i) => (
                                <motion.div
                                    key={i}
                                    initial={{ opacity: 0, y: 20 }}
                                    whileInView={{ opacity: 1, y: 0 }}
                                    viewport={{ once: true }}
                                    transition={{ delay: i * 0.1 }}
                                    className="relative p-8 rounded-2xl bg-white/5 border border-white/10 backdrop-blur-sm z-10"
                                >
                                    <span className="text-4xl font-display font-black text-white/10 absolute top-4 right-6">{item.step}</span>
                                    <div className="mb-6">{item.icon}</div>
                                    <h3 className="text-xl font-display font-bold mb-3">{item.title}</h3>
                                    <p className="text-gray-400 text-sm leading-relaxed">{item.desc}</p>
                                </motion.div>
                            ))}
                        </div>
                    </div>
                </div>
            </section>
        </div>
    );
};

export default Home;
