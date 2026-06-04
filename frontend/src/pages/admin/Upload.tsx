import { useState, useRef, useEffect, type FormEvent, type ChangeEvent } from 'react';
import { Upload as UploadIcon, FileText, Shield, CheckCircle2 } from 'lucide-react';
import { diplomaApi, universityApi, facultyApi, departmentApi } from '../../services/api';

interface Faculty { id: string; name: string; }
interface Department { id: string; name: string; }
interface University { id: string; name: string; }

interface FormState {
  firstName: string;
  lastName: string;
  email: string;
  studentNumber: string;
  university: string;       // university name (string) — what backend expects
  universityId: string;     // uuid — used to fetch faculties
  faculty: string;
  facultyId: string;
  department: string;
  graduationYear: number;
  nationality: string;
}

const EMPTY: FormState = {
  firstName: '', lastName: '', email: '', studentNumber: '',
  university: '', universityId: '', faculty: '', facultyId: '',
  department: '', graduationYear: new Date().getFullYear(), nationality: 'TR',
};

export default function Upload() {
  const [universities, setUniversities] = useState<University[]>([]);
  const [faculties, setFaculties] = useState<Faculty[]>([]);
  const [departments, setDepartments] = useState<Department[]>([]);
  const [form, setForm] = useState<FormState>(EMPTY);
  const [file, setFile] = useState<File | null>(null);
  const [loading, setLoading] = useState(false);
  const [success, setSuccess] = useState(false);
  const [error, setError] = useState('');
  const fileRef = useRef<HTMLInputElement>(null);

  useEffect(() => {
    universityApi.getAll().then((r) => setUniversities(r.data ?? [])).catch(() => {});
  }, []);

  const onUniversityChange = async (id: string, name: string) => {
    setForm((f) => ({ ...f, universityId: id, university: name, faculty: '', facultyId: '', department: '' }));
    setFaculties([]);
    setDepartments([]);
    if (!id) return;
    try {
      const res = await facultyApi.getByUniversity(id);
      setFaculties(res.data ?? []);
    } catch { /* faculties not critical */ }
  };

  const onFacultyChange = async (id: string, name: string) => {
    setForm((f) => ({ ...f, facultyId: id, faculty: name, department: '' }));
    setDepartments([]);
    if (!id) return;
    try {
      const res = await departmentApi.getByFaculty(id);
      setDepartments(res.data ?? []);
    } catch { /* departments not critical */ }
  };

  const handleFile = (e: ChangeEvent<HTMLInputElement>) => {
    const f = e.target.files?.[0];
    if (f) setFile(f);
  };

  const handleDrop = (e: React.DragEvent) => {
    e.preventDefault();
    const f = e.dataTransfer.files?.[0];
    if (f) setFile(f);
  };

  const handleSubmit = async (e: FormEvent) => {
    e.preventDefault();
    if (!file) { setError('Lütfen PDF diploma belgesini yükleyin.'); return; }
    if (!file.name.toLowerCase().endsWith('.pdf')) { setError('Sadece PDF dosyaları kabul edilmektedir.'); return; }
    setError('');
    setLoading(true);
    try {
      const fd = new FormData();
      fd.append('firstName', form.firstName);
      fd.append('lastName', form.lastName);
      fd.append('email', form.email);
      fd.append('studentNumber', form.studentNumber);
      fd.append('university', form.university);
      fd.append('faculty', form.faculty);
      fd.append('department', form.department);
      fd.append('graduationYear', String(form.graduationYear));
      fd.append('nationality', form.nationality);
      fd.append('diploma', file);   // backend expects field name "diploma"
      await diplomaApi.upload(fd);
      setSuccess(true);
    } catch {
      setError('Diploma yüklenirken bir hata oluştu. Lütfen tekrar deneyin.');
    } finally {
      setLoading(false);
    }
  };

  if (success) {
    return (
      <div className="p-8 flex flex-col items-center justify-center min-h-[60vh]">
        <div className="w-16 h-16 rounded-full bg-secondary/10 flex items-center justify-center mb-4">
          <CheckCircle2 className="w-8 h-8 text-secondary" />
        </div>
        <h2 className="text-xl font-bold text-on-surface mb-2">Diploma Başarıyla Yüklendi!</h2>
        <p className="text-on-surface-variant text-sm mb-6">
          Diploma blockchain'e kaydedildi ve doğrulanmayı bekliyor.
        </p>
        <button
          onClick={() => { setSuccess(false); setFile(null); setForm(EMPTY); setFaculties([]); setDepartments([]); }}
          className="px-6 py-2.5 bg-primary text-on-primary rounded-xl text-sm font-semibold hover:bg-primary/90 transition-all"
        >
          Yeni Diploma Yükle
        </button>
      </div>
    );
  }

  return (
    <div className="p-8 max-w-3xl mx-auto w-full">
      <div className="mb-8">
        <h1 className="text-2xl font-bold text-on-surface">Yeni Diploma Yükle</h1>
        <p className="text-sm text-on-surface-variant mt-1">
          Blockchain tabanlı doğrulama için öğrenci ve mezuniyet verilerini güvenli bir şekilde sisteme aktarın.
        </p>
      </div>

      {error && (
        <div className="mb-6 px-4 py-3 rounded-xl bg-error-container text-sm text-error font-medium">
          {error}
        </div>
      )}

      <form onSubmit={handleSubmit}>
        <div className="bg-white rounded-2xl border border-outline-variant/40 card-shadow p-8 mb-6">
          <div className="flex items-center gap-2 mb-6">
            <Shield className="w-5 h-5 text-secondary" />
            <h2 className="font-semibold text-on-surface text-lg">Akademik Bilgiler</h2>
          </div>

          <div className="grid grid-cols-1 md:grid-cols-2 gap-5">
            <Field label="İsim" placeholder="Örn: Ahmet" value={form.firstName}
              onChange={(v) => setForm((f) => ({ ...f, firstName: v }))} />
            <Field label="Soyisim" placeholder="Örn: Yılmaz" value={form.lastName}
              onChange={(v) => setForm((f) => ({ ...f, lastName: v }))} />
            <Field label="E-posta Adresi" type="email" placeholder="ahmet.yilmaz@edu.tr" value={form.email}
              onChange={(v) => setForm((f) => ({ ...f, email: v }))} />
            <Field label="Öğrenci Numarası" placeholder="2024010203" value={form.studentNumber}
              onChange={(v) => setForm((f) => ({ ...f, studentNumber: v }))} />

            {/* University dropdown */}
            <div className="md:col-span-2">
              <label className="block text-xs font-semibold text-on-surface mb-1.5">Üniversite</label>
              <select
                value={form.universityId}
                onChange={(e) => {
                  const opt = universities.find((u) => u.id === e.target.value);
                  onUniversityChange(e.target.value, opt?.name ?? '');
                }}
                required
                className="w-full px-4 py-2.5 rounded-lg border border-outline-variant/60 text-sm bg-white focus:outline-none focus:border-on-tertiary-container focus:ring-2 focus:ring-on-tertiary-container/10 transition-all"
              >
                <option value="">Üniversite Seçiniz</option>
                {universities.map((u) => (
                  <option key={u.id} value={u.id}>{u.name}</option>
                ))}
              </select>
            </div>

            {/* Faculty dropdown */}
            <div>
              <label className="block text-xs font-semibold text-on-surface mb-1.5">Fakülte</label>
              {faculties.length > 0 ? (
                <select
                  value={form.facultyId}
                  onChange={(e) => {
                    const opt = faculties.find((f) => f.id === e.target.value);
                    onFacultyChange(e.target.value, opt?.name ?? '');
                  }}
                  required
                  className="w-full px-4 py-2.5 rounded-lg border border-outline-variant/60 text-sm bg-white focus:outline-none focus:border-on-tertiary-container focus:ring-2 focus:ring-on-tertiary-container/10 transition-all"
                >
                  <option value="">Fakülte Seçiniz</option>
                  {faculties.map((f) => (
                    <option key={f.id} value={f.id}>{f.name}</option>
                  ))}
                </select>
              ) : (
                <input type="text" placeholder="Örn: Mühendislik Fakültesi" value={form.faculty}
                  onChange={(e) => setForm((f) => ({ ...f, faculty: e.target.value }))}
                  required
                  className="w-full px-4 py-2.5 rounded-lg border border-outline-variant/60 text-sm focus:outline-none focus:border-on-tertiary-container focus:ring-2 focus:ring-on-tertiary-container/10 placeholder:text-outline transition-all"
                />
              )}
            </div>

            {/* Department dropdown */}
            <div>
              <label className="block text-xs font-semibold text-on-surface mb-1.5">Bölüm</label>
              {departments.length > 0 ? (
                <select
                  value={form.department}
                  onChange={(e) => setForm((f) => ({ ...f, department: e.target.value }))}
                  required
                  className="w-full px-4 py-2.5 rounded-lg border border-outline-variant/60 text-sm bg-white focus:outline-none focus:border-on-tertiary-container focus:ring-2 focus:ring-on-tertiary-container/10 transition-all"
                >
                  <option value="">Bölüm Seçiniz</option>
                  {departments.map((d) => (
                    <option key={d.id} value={d.name}>{d.name}</option>
                  ))}
                </select>
              ) : (
                <input type="text" placeholder="Örn: Bilgisayar Mühendisliği" value={form.department}
                  onChange={(e) => setForm((f) => ({ ...f, department: e.target.value }))}
                  required
                  className="w-full px-4 py-2.5 rounded-lg border border-outline-variant/60 text-sm focus:outline-none focus:border-on-tertiary-container focus:ring-2 focus:ring-on-tertiary-container/10 placeholder:text-outline transition-all"
                />
              )}
            </div>

            <div>
              <label className="block text-xs font-semibold text-on-surface mb-1.5">Mezuniyet Yılı</label>
              <input type="number" min={2000} max={2100} value={form.graduationYear}
                onChange={(e) => setForm((f) => ({ ...f, graduationYear: Number(e.target.value) }))}
                required
                className="w-full px-4 py-2.5 rounded-lg border border-outline-variant/60 text-sm focus:outline-none focus:border-on-tertiary-container focus:ring-2 focus:ring-on-tertiary-container/10 transition-all"
              />
            </div>

            <Field label="Uyruk" placeholder="TR" value={form.nationality}
              onChange={(v) => setForm((f) => ({ ...f, nationality: v }))} />
          </div>
        </div>

        {/* File Upload — PDF only */}
        <div className="bg-white rounded-2xl border border-outline-variant/40 card-shadow p-8 mb-6">
          <label className="block text-xs font-semibold text-on-surface mb-1">
            Diploma Belgesi
          </label>
          <p className="text-xs text-on-surface-variant mb-3">Yalnızca PDF formatı kabul edilmektedir. Maks. 10MB.</p>
          <div
            onDrop={handleDrop}
            onDragOver={(e) => e.preventDefault()}
            onClick={() => fileRef.current?.click()}
            className={`border-2 border-dashed rounded-xl p-10 text-center cursor-pointer transition-all ${
              file
                ? 'border-secondary bg-secondary/5'
                : 'border-outline-variant/60 hover:border-on-tertiary-container/50 hover:bg-surface'
            }`}
          >
            <input ref={fileRef} type="file" accept=".pdf" onChange={handleFile} className="hidden" />
            {file ? (
              <div className="flex flex-col items-center gap-2">
                <FileText className="w-8 h-8 text-secondary" />
                <p className="text-sm font-medium text-on-surface">{file.name}</p>
                <p className="text-xs text-on-surface-variant">{(file.size / 1024 / 1024).toFixed(2)} MB</p>
              </div>
            ) : (
              <div className="flex flex-col items-center gap-2">
                <UploadIcon className="w-8 h-8 text-outline" />
                <p className="text-sm text-on-surface font-medium">Dosyayı buraya sürükleyin veya seçin</p>
                <p className="text-xs text-on-surface-variant">PDF · Maksimum 10MB</p>
              </div>
            )}
          </div>
        </div>

        <div className="flex items-center justify-between">
          <div className="flex items-center gap-2 text-xs text-on-surface-variant">
            <Shield className="w-3.5 h-3.5 text-secondary" />
            <span>Verileriniz AES-256 ile şifrelenmektedir.</span>
          </div>
          <button
            type="submit"
            disabled={loading}
            className="flex items-center gap-2 px-8 py-3 bg-primary text-on-primary rounded-xl font-semibold text-sm hover:bg-primary/90 transition-all disabled:opacity-60"
          >
            <UploadIcon className="w-4 h-4" />
            {loading ? 'Yükleniyor...' : 'Diploma Yükle'}
          </button>
        </div>
      </form>
    </div>
  );
}

interface FieldProps {
  label: string;
  value: string;
  onChange: (v: string) => void;
  placeholder?: string;
  type?: string;
}

function Field({ label, value, onChange, placeholder, type = 'text' }: Readonly<FieldProps>) {
  return (
    <div>
      <label className="block text-xs font-semibold text-on-surface mb-1.5">{label}</label>
      <input
        type={type}
        value={value}
        onChange={(e) => onChange(e.target.value)}
        placeholder={placeholder}
        required
        className="w-full px-4 py-2.5 rounded-lg border border-outline-variant/60 text-sm focus:outline-none focus:border-on-tertiary-container focus:ring-2 focus:ring-on-tertiary-container/10 placeholder:text-outline transition-all"
      />
    </div>
  );
}
