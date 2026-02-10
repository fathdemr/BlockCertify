import React, { useState } from 'react';
import { motion, AnimatePresence } from 'framer-motion';
import {
    Upload as UploadIcon,
    FileText,
    User,
    Calendar,
    Hash,
    CheckCircle2,
    Database,
    Cpu,
    ShieldCheck,
    ArrowRight,
    Loader2,
    X,
    Globe
} from 'lucide-react';
import { diplomaService } from '../../services/api';

type UploadStep = 'form' | 'hashing' | 'arweave' | 'polygon' | 'success';

interface DiplomaMetaData {
    firstName: string;
    lastName: string;
    email: string;
    university: string;
    faculty: string;
    department: string;
    graduationYear: number;
    studentNumber: string;
    nationality: string;
}

const UploadDiploma: React.FC = () => {
    const [step, setStep] = useState<UploadStep>('form');
    const [formData, setFormData] = useState<DiplomaMetaData & { file: File | null }>({
        firstName: '',
        lastName: '',
        email: '',
        university: 'Karabuk University',
        faculty: '',
        department: '',
        graduationYear: 2025,
        studentNumber: '',
        nationality: '',
        file: null,
    });

    const [txData, setTxData] = useState({
        arweaveTx: '',
        polygonTx: '',
        fileHash: '',
    });

    const handleFileChange = (e: React.ChangeEvent<HTMLInputElement>) => {
        if (e.target.files && e.target.files[0]) {
            setFormData({ ...formData, file: e.target.files[0] });
        }
    };

    const startUpload = (e: React.FormEvent) => {
        e.preventDefault();
        if (!formData.file) return;

        // Start logical process
        processFlow();
    };

    const processFlow = async () => {
        if (!formData.file) return;

        try {
            setStep('hashing');
            // Logical split: we call the backend once, but update the UI to show progress
            // Mocking the intermediate steps for UI smoothness if the backend is too fast
            await new Promise(r => setTimeout(r, 800));

            setStep('arweave');
            await new Promise(r => setTimeout(r, 800));

            setStep('polygon');

            const uploadFormData = new FormData();
            uploadFormData.append('diploma', formData.file);
            uploadFormData.append('firstName', formData.firstName);
            uploadFormData.append('lastName', formData.lastName);
            uploadFormData.append('email', formData.email);
            uploadFormData.append('university', formData.university);
            uploadFormData.append('faculty', formData.faculty);
            uploadFormData.append('department', formData.department);
            uploadFormData.append('graduationYear', formData.graduationYear.toString());
            uploadFormData.append('studentNumber', formData.studentNumber);
            uploadFormData.append('nationality', formData.nationality);

            const response = await diplomaService.upload(uploadFormData);

            setTxData({
                arweaveTx: response.arweaveTxID,
                polygonTx: response.polygonTxHash,
                fileHash: response.diplomaHash,
            });

            setStep('success');
        } catch (error: any) {
            console.error('Upload failed:', error);
            alert(error.response?.data?.message || 'Upload failed. Please check the backend connection.');
            setStep('form');
        }
    };

    const reset = () => {
        setStep('form');
        setFormData({
            firstName: '',
            lastName: '',
            email: '',
            university: 'Global Tech University',
            faculty: '',
            department: '',
            graduationYear: 2025,
            studentNumber: '',
            nationality: '',
            file: null,
        });
    };

    return (
        <div className="max-w-4xl mx-auto">
            <div className="mb-8">
                <h1 className="text-3xl font-display font-bold mb-2">Issue New Diploma</h1>
                <p className="text-gray-400">Upload a PDF and register its integrity on the blockchain.</p>
            </div>

            <AnimatePresence mode="wait">
                {step === 'form' && (
                    <motion.form
                        key="form"
                        initial={{ opacity: 0, y: 20 }}
                        animate={{ opacity: 1, y: 0 }}
                        exit={{ opacity: 0, y: -20 }}
                        onSubmit={startUpload}
                        className="space-y-6"
                    >
                        <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                            {/* Row 1: First Name & Last Name */}
                            <div>
                                <label className="block text-sm font-medium text-gray-400 mb-2">First Name</label>
                                <div className="relative">
                                    <User className="absolute left-4 top-1/2 -translate-y-1/2 h-5 w-5 text-gray-500" />
                                    <input
                                        type="text"
                                        value={formData.firstName}
                                        onChange={(e) => setFormData({ ...formData, firstName: e.target.value })}
                                        className="w-full pl-12 pr-4 py-3 bg-white/5 border border-white/10 rounded-xl focus:ring-2 focus:ring-brand-primary/50"
                                        placeholder="e.g. John"
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
                                        value={formData.lastName}
                                        onChange={(e) => setFormData({ ...formData, lastName: e.target.value })}
                                        className="w-full pl-12 pr-4 py-3 bg-white/5 border border-white/10 rounded-xl focus:ring-2 focus:ring-brand-primary/50"
                                        placeholder="e.g. Doe"
                                        required
                                    />
                                </div>
                            </div>

                            {/* Row 2: Email & Student Number */}
                            <div>
                                <label className="block text-sm font-medium text-gray-400 mb-2">Email Address</label>
                                <div className="relative">
                                    <FileText className="absolute left-4 top-1/2 -translate-y-1/2 h-5 w-5 text-gray-500" />
                                    <input
                                        type="email"
                                        value={formData.email}
                                        onChange={(e) => setFormData({ ...formData, email: e.target.value })}
                                        className="w-full pl-12 pr-4 py-3 bg-white/5 border border-white/10 rounded-xl focus:ring-2 focus:ring-brand-primary/50"
                                        placeholder="e.g. john.doe@example.com"
                                        required
                                    />
                                </div>
                            </div>
                            <div>
                                <label className="block text-sm font-medium text-gray-400 mb-2">Student Number</label>
                                <div className="relative">
                                    <Hash className="absolute left-4 top-1/2 -translate-y-1/2 h-5 w-5 text-gray-500" />
                                    <input
                                        type="text"
                                        value={formData.studentNumber}
                                        onChange={(e) => setFormData({ ...formData, studentNumber: e.target.value })}
                                        className="w-full pl-12 pr-4 py-3 bg-white/5 border border-white/10 rounded-xl focus:ring-2 focus:ring-brand-primary/50"
                                        placeholder="e.g. 202100456"
                                        required
                                    />
                                </div>
                            </div>

                            {/* Row 3: University (Read-only) & Faculty */}
                            <div>
                                <label className="block text-sm font-medium text-gray-400 mb-2">University</label>
                                <div className="relative">
                                    <FileText className="absolute left-4 top-1/2 -translate-y-1/2 h-5 w-5 text-gray-500" />
                                    <input
                                        type="text"
                                        value={formData.university}
                                        readOnly
                                        className="w-full pl-12 pr-4 py-3 bg-white/5 border border-white/10 rounded-xl opacity-60 cursor-not-allowed"
                                    />
                                </div>
                            </div>
                            <div>
                                <label className="block text-sm font-medium text-gray-400 mb-2">Faculty</label>
                                <div className="relative">
                                    <FileText className="absolute left-4 top-1/2 -translate-y-1/2 h-5 w-5 text-gray-500" />
                                    <input
                                        type="text"
                                        value={formData.faculty}
                                        onChange={(e) => setFormData({ ...formData, faculty: e.target.value })}
                                        className="w-full pl-12 pr-4 py-3 bg-white/5 border border-white/10 rounded-xl focus:ring-2 focus:ring-brand-primary/50"
                                        placeholder="e.g. Faculty of Engineering"
                                        required
                                    />
                                </div>
                            </div>

                            {/* Row 4: Department & Graduation Year */}
                            <div>
                                <label className="block text-sm font-medium text-gray-400 mb-2">Department</label>
                                <div className="relative">
                                    <FileText className="absolute left-4 top-1/2 -translate-y-1/2 h-5 w-5 text-gray-500" />
                                    <input
                                        type="text"
                                        value={formData.department}
                                        onChange={(e) => setFormData({ ...formData, department: e.target.value })}
                                        className="w-full pl-12 pr-4 py-3 bg-white/5 border border-white/10 rounded-xl focus:ring-2 focus:ring-brand-primary/50"
                                        placeholder="e.g. Computer Science"
                                        required
                                    />
                                </div>
                            </div>
                            <div>
                                <label className="block text-sm font-medium text-gray-400 mb-2">Graduation Year</label>
                                <div className="relative">
                                    <Calendar className="absolute left-4 top-1/2 -translate-y-1/2 h-5 w-5 text-gray-500" />
                                    <select
                                        value={formData.graduationYear}
                                        onChange={(e) => setFormData({ ...formData, graduationYear: parseInt(e.target.value) })}
                                        className="w-full pl-12 pr-4 py-3 bg-white/5 border border-white/10 rounded-xl focus:ring-2 focus:ring-brand-primary/50"
                                    >
                                        {[2024, 2025, 2026, 2027].map(year => (
                                            <option key={year} value={year}>{year}</option>
                                        ))}
                                    </select>
                                </div>
                            </div>

                            {/* Row 5: Nationality */}
                            <div className="md:col-span-2">
                                <label className="block text-sm font-medium text-gray-400 mb-2">Nationality</label>
                                <div className="relative">
                                    <Globe className="absolute left-4 top-1/2 -translate-y-1/2 h-5 w-5 text-gray-500" />
                                    <input
                                        type="text"
                                        value={formData.nationality}
                                        onChange={(e) => setFormData({ ...formData, nationality: e.target.value })}
                                        className="w-full pl-12 pr-4 py-3 bg-white/5 border border-white/10 rounded-xl focus:ring-2 focus:ring-brand-primary/50"
                                        placeholder="e.g. Turkish"
                                        required
                                    />
                                </div>
                            </div>
                        </div>

                        <div className="relative">
                            <label className="block text-sm font-medium text-gray-400 mb-4">Diploma PDF File</label>
                            <div
                                className={`border-2 border-dashed rounded-3xl p-12 text-center transition-all ${formData.file ? 'border-brand-success/50 bg-brand-success/5' : 'border-white/10 hover:border-brand-primary/30 bg-white/5'
                                    }`}
                            >
                                <input
                                    type="file"
                                    accept=".pdf"
                                    onChange={handleFileChange}
                                    className="absolute inset-0 w-full h-full opacity-0 cursor-pointer"
                                    required={!formData.file}
                                />
                                {formData.file ? (
                                    <div className="flex flex-col items-center">
                                        <CheckCircle2 className="h-12 w-12 text-brand-success mb-2" />
                                        <p className="text-white font-medium">{formData.file.name}</p>
                                        <button
                                            type="button"
                                            onClick={(e) => {
                                                e.preventDefault();
                                                setFormData({ ...formData, file: null });
                                            }}
                                            className="mt-4 text-xs text-red-400 hover:text-red-300 flex items-center gap-1"
                                        >
                                            <X className="h-3 w-3" /> Remove File
                                        </button>
                                    </div>
                                ) : (
                                    <div className="flex flex-col items-center">
                                        <UploadIcon className="h-12 w-12 text-gray-600 mb-4" />
                                        <p className="text-lg font-medium text-gray-300">Click or drag PDF to upload</p>
                                        <p className="text-sm text-gray-500 mt-1">Maximum file size: 10MB</p>
                                    </div>
                                )}
                            </div>
                        </div>

                        <button
                            type="submit"
                            className="w-full py-5 bg-brand-primary hover:bg-brand-primary/90 text-white rounded-2xl font-display font-bold text-lg flex items-center justify-center gap-3 transition-all shadow-[0_10px_30px_rgba(59,130,246,0.3)]"
                        >
                            Start Blockchain Issuance
                            <ArrowRight className="h-5 w-5" />
                        </button>
                    </motion.form>
                )}

                {(step === 'hashing' || step === 'arweave' || step === 'polygon') && (
                    <motion.div
                        key="processing"
                        initial={{ opacity: 0 }}
                        animate={{ opacity: 1 }}
                        exit={{ opacity: 0 }}
                        className="p-12 rounded-3xl bg-white/5 border border-white/10 text-center"
                    >
                        <div className="relative w-24 h-24 mx-auto mb-8">
                            <div className="absolute inset-0 bg-brand-primary/20 rounded-full blur-2xl animate-pulse" />
                            <Loader2 className="w-24 h-24 text-brand-primary animate-spin" />
                        </div>

                        <h2 className="text-2xl font-display font-bold mb-8 capitalize">{step.replace('hashing', 'Generating Fingerprint').replace('arweave', 'Storing on Permaweb').replace('polygon', 'Registering on Polygon')}...</h2>

                        <div className="max-w-md mx-auto space-y-4">
                            {[
                                { id: 'hashing', label: 'File Hashing', icon: Cpu, active: step === 'hashing', done: step !== 'hashing' },
                                { id: 'arweave', label: 'Arweave Storage', icon: Database, active: step === 'arweave', done: step === 'polygon' },
                                { id: 'polygon', label: 'Polygon Registration', icon: ShieldCheck, active: step === 'polygon', done: false },
                            ].map((item, i) => (
                                <div
                                    key={i}
                                    className={`flex items-center gap-4 p-4 rounded-xl border transition-all ${item.active ? 'bg-brand-primary/10 border-brand-primary/30' :
                                        item.done ? 'bg-brand-success/10 border-brand-success/30 opacity-60' :
                                            'bg-white/5 border-white/5 opacity-40'
                                        }`}
                                >
                                    <item.icon className={`h-6 w-6 ${item.active ? 'text-brand-primary' : item.done ? 'text-brand-success' : 'text-gray-500'}`} />
                                    <span className={`font-medium ${item.active ? 'text-white' : 'text-gray-400'}`}>{item.label}</span>
                                    {item.done && <CheckCircle2 className="h-5 w-5 text-brand-success ml-auto" />}
                                    {item.active && <Loader2 className="h-5 w-5 text-brand-primary animate-spin ml-auto" />}
                                </div>
                            ))}
                        </div>
                    </motion.div>
                )}

                {step === 'success' && (
                    <motion.div
                        key="success"
                        initial={{ opacity: 0, scale: 0.95 }}
                        animate={{ opacity: 1, scale: 1 }}
                        className="p-12 rounded-3xl bg-brand-success/5 border border-brand-success/20 text-center"
                    >
                        <div className="w-20 h-20 bg-brand-success/20 rounded-full flex items-center justify-center mx-auto mb-6">
                            <CheckCircle2 className="h-12 w-12 text-brand-success" />
                        </div>
                        <h2 className="text-3xl font-display font-bold text-white mb-2">Diploma Issued Successfully</h2>
                        <p className="text-gray-400 mb-8">The digital credential is now permanent and verifiable on-chain.</p>

                        <div className="grid grid-cols-1 gap-4 mb-10 text-left">
                            <div className="p-4 bg-white/5 rounded-xl border border-white/5">
                                <p className="text-xs text-gray-500 uppercase mb-1">File Hash (integrity)</p>
                                <code className="text-sm text-brand-secondary break-all">{txData.fileHash}</code>
                            </div>
                            <div className="p-4 bg-white/5 rounded-xl border border-white/5">
                                <p className="text-xs text-gray-500 uppercase mb-1">Arweave Transaction</p>
                                <code className="text-sm text-brand-accent break-all">{txData.arweaveTx}</code>
                            </div>
                            <div className="p-4 bg-white/5 rounded-xl border border-white/5">
                                <p className="text-xs text-gray-500 uppercase mb-1">Polygon Transaction</p>
                                <code className="text-sm text-brand-primary break-all">{txData.polygonTx}</code>
                            </div>
                        </div>

                        <div className="flex flex-col sm:flex-row gap-4">
                            <button
                                onClick={reset}
                                className="flex-grow py-4 bg-white/5 hover:bg-white/10 text-white rounded-xl font-bold transition-all"
                            >
                                Issue Another
                            </button>
                            <button
                                onClick={() => setStep('form')}
                                className="flex-grow py-4 bg-brand-primary hover:bg-brand-primary/90 text-white rounded-xl font-bold transition-all shadow-lg"
                            >
                                Go to Verification
                            </button>
                        </div>
                    </motion.div>
                )}
            </AnimatePresence>
        </div>
    );
};

export default UploadDiploma;
