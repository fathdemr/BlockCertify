import React from 'react';
import { Shield, Github, Twitter, Linkedin } from 'lucide-react';

const Footer: React.FC = () => {
    return (
        <footer className="bg-brand-dark border-t border-gray-800 pt-12 pb-8">
            <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
                <div className="grid grid-cols-1 md:grid-cols-4 gap-8">
                    <div className="col-span-1 md:col-span-2">
                        <div className="flex items-center gap-2 mb-4">
                            <Shield className="h-6 w-6 text-brand-secondary" />
                            <span className="text-lg font-display font-bold">BlockCertify</span>
                        </div>
                        <p className="text-gray-400 max-w-sm">
                            Securing academic integrity through decentralized blockchain technologies.
                            Verified by Polygon, stored permanently on Arweave.
                        </p>
                    </div>
                    <div>
                        <h3 className="text-sm font-semibold text-white uppercase tracking-wider mb-4">Resources</h3>
                        <ul className="space-y-2">
                            <li><a href="#" className="text-gray-400 hover:text-brand-secondary transition-colors">How It Works</a></li>
                            <li><a href="#" className="text-gray-400 hover:text-brand-secondary transition-colors">Whitepaper</a></li>
                            <li><a href="#" className="text-gray-400 hover:text-brand-secondary transition-colors">API Keys</a></li>
                        </ul>
                    </div>
                    <div>
                        <h3 className="text-sm font-semibold text-white uppercase tracking-wider mb-4">Connect</h3>
                        <div className="flex space-x-4">
                            <a href="#" className="text-gray-400 hover:text-brand-secondary transition-colors"><Github className="h-5 w-5" /></a>
                            <a href="#" className="text-gray-400 hover:text-brand-secondary transition-colors"><Twitter className="h-5 w-5" /></a>
                            <a href="#" className="text-gray-400 hover:text-brand-secondary transition-colors"><Linkedin className="h-5 w-5" /></a>
                        </div>
                    </div>
                </div>
                <div className="mt-12 pt-8 border-t border-gray-800 flex flex-col md:flex-row justify-between items-center gap-4">
                    <p className="text-gray-500 text-sm">
                        © {new Date().getFullYear()} BlockCertify. Built for Academic Integrity.
                    </p>
                    <div className="flex space-x-6">
                        <a href="#" className="text-gray-500 hover:text-gray-400 text-sm">Privacy Policy</a>
                        <a href="#" className="text-gray-500 hover:text-gray-400 text-sm">Terms of Service</a>
                    </div>
                </div>
            </div>
        </footer>
    );
};

export default Footer;
