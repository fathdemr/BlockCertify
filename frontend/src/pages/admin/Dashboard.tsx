import { useEffect, useRef, useState } from 'react';
import { Link } from 'react-router-dom';
import { Plus, ArrowRight, Share2, Download, TrendingUp, CheckCircle2 } from 'lucide-react';
import { QRCodeSVG, QRCodeCanvas } from 'qrcode.react';
import { diplomaApi } from '../../services/api';
import type { Diploma } from '../../types';
import StatusBadge from '../../components/StatusBadge';
import { useAuth } from '../../context/AuthContext';

export default function Dashboard() {
  const { user } = useAuth();
  const [diplomas, setDiplomas] = useState<Diploma[]>([]);
  const [loading, setLoading] = useState(true);
  const [selected, setSelected] = useState<Diploma | null>(null);
  const [copied, setCopied] = useState(false);
  const qrCanvasRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    diplomaApi.getAll()
      .then((r) => {
        const data = r.data ?? [];
        setDiplomas(data);
        if (data.length > 0) setSelected(data[0]);
      })
      .catch(() => {})
      .finally(() => setLoading(false));
  }, []);

  const preview = diplomas.slice(0, 4);

  const verifyUrl = selected?.tx_hash
    ? `${window.location.origin}/dogrula?tx=${selected.tx_hash}`
    : '';

  async function handleShare() {
    if (!verifyUrl) return;
    const text = `${selected!.first_name} ${selected!.last_name} diplomasını BlockCertify üzerinde doğrulayabilirsiniz: ${verifyUrl}`;
    if (navigator.share) {
      try {
        await navigator.share({ title: 'BlockCertify Diploma', text, url: verifyUrl });
        return;
      } catch {
        // fall through to clipboard
      }
    }
    await navigator.clipboard.writeText(verifyUrl);
    setCopied(true);
    setTimeout(() => setCopied(false), 2500);
  }

  function handleDownload() {
    const canvas = qrCanvasRef.current?.querySelector('canvas');
    if (!canvas) return;
    const url = canvas.toDataURL('image/png');
    const a = document.createElement('a');
    a.href = url;
    a.download = `diploma-qr-${selected?.diploma_no ?? 'code'}.png`;
    a.click();
  }

  return (
    <div className="p-8">
      <div className="flex items-center justify-between mb-8">
        <div>
          <h1 className="text-2xl font-bold text-on-surface">
            Hoş geldin {user?.name?.split(' ')[0] ?? 'Admin'}
          </h1>
          <p className="text-sm text-on-surface-variant mt-1">
            Blockchain tabanlı sertifika yönetim panelindesiniz.
          </p>
        </div>
        <Link
          to="/admin/yukle"
          className="flex items-center gap-2 px-5 py-2.5 bg-secondary text-on-secondary rounded-xl font-semibold text-sm hover:bg-secondary/90 transition-all"
        >
          <Plus className="w-4 h-4" />
          Yeni Sertifika Oluştur
        </Link>
      </div>

      <div className="grid grid-cols-1 xl:grid-cols-3 gap-6">
        {/* Diploma Table */}
        <div className="xl:col-span-2 bg-white rounded-2xl border border-outline-variant/40 card-shadow overflow-hidden">
          <div className="px-6 py-4 flex items-center justify-between border-b border-outline-variant/30">
            <h2 className="font-semibold text-on-surface">Diplomalarım</h2>
            <Link
              to="/admin/gecmis"
              className="text-xs text-on-tertiary-container font-medium hover:underline flex items-center gap-1"
            >
              Tümünü Gör <ArrowRight className="w-3 h-3" />
            </Link>
          </div>

          {loading ? (
            <div className="flex items-center justify-center py-16">
              <div className="w-6 h-6 border-2 border-secondary border-t-transparent rounded-full animate-spin" />
            </div>
          ) : preview.length === 0 ? (
            <div className="flex flex-col items-center justify-center py-16 text-on-surface-variant">
              <p className="text-sm">Henüz diploma kaydı yok.</p>
              <Link to="/admin/yukle" className="mt-3 text-xs text-on-tertiary-container font-medium hover:underline">
                İlk diplomayı ekle →
              </Link>
            </div>
          ) : (
            <div className="overflow-x-auto">
              <table className="w-full">
                <thead>
                  <tr className="bg-surface text-xs text-on-surface-variant uppercase tracking-wider">
                    <th className="px-6 py-3 text-left font-medium">Diploma No</th>
                    <th className="px-6 py-3 text-left font-medium">İsim</th>
                    <th className="px-6 py-3 text-left font-medium">Program</th>
                    <th className="px-6 py-3 text-left font-medium">Durum</th>
                    <th className="px-6 py-3 text-left font-medium">Eylem</th>
                  </tr>
                </thead>
                <tbody className="divide-y divide-outline-variant/20">
                  {preview.map((d) => (
                    <tr
                      key={d.id}
                      onClick={() => setSelected(d)}
                      className={`cursor-pointer transition-colors ${
                        selected?.id === d.id
                          ? 'bg-secondary/8 border-l-2 border-secondary'
                          : 'hover:bg-surface/50'
                      }`}
                    >
                      <td className="px-6 py-4">
                        <span className="text-xs font-mono-code text-on-surface-variant">
                          #{d.diploma_no ?? String(d.id).slice(0, 8)}
                        </span>
                      </td>
                      <td className="px-6 py-4">
                        <span className="font-semibold text-sm text-on-surface">
                          {d.first_name} {d.last_name}
                        </span>
                      </td>
                      <td className="px-6 py-4 text-sm text-on-surface-variant">{d.department}</td>
                      <td className="px-6 py-4">
                        <StatusBadge status={d.status} />
                      </td>
                      <td className="px-6 py-4">
                        <Link
                          to={`/admin/dogrula`}
                          onClick={(e) => e.stopPropagation()}
                          className="text-xs text-on-tertiary-container font-medium flex items-center gap-1 hover:gap-2 transition-all w-fit"
                        >
                          Git <ArrowRight className="w-3 h-3" />
                        </Link>
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          )}
        </div>

        {/* Right Panel */}
        <div className="space-y-4">
          {/* Share Card */}
          <div className="bg-primary-container rounded-2xl p-6 text-white">
            <div className="flex items-center justify-center mb-4">
              {selected && verifyUrl ? (
                <div className="w-28 h-28 bg-white rounded-xl p-2 flex items-center justify-center">
                  <QRCodeSVG
                    value={verifyUrl}
                    size={96}
                    bgColor="#ffffff"
                    fgColor="#131b2e"
                    level="M"
                  />
                </div>
              ) : (
                <div className="w-28 h-28 bg-white/10 rounded-xl border border-white/20 flex items-center justify-center">
                  <p className="text-xs text-white/40 text-center px-2">Diploma seç</p>
                </div>
              )}
            </div>

            {selected ? (
              <div className="mb-4">
                <h3 className="font-semibold text-center text-sm">{selected.first_name} {selected.last_name}</h3>
                <p className="text-xs text-white/60 text-center mt-0.5">{selected.department}</p>
              </div>
            ) : (
              <div className="mb-4">
                <h3 className="font-semibold text-center">Diplomanı Paylaş</h3>
                <p className="text-xs text-white/60 text-center mt-1 leading-relaxed">
                  Tablodan bir diploma seç, QR kodu oluşturmak ve paylaşmak için.
                </p>
              </div>
            )}

            <div className="space-y-2">
              <button
                onClick={handleShare}
                disabled={!selected}
                className="w-full flex items-center justify-center gap-2 py-2.5 bg-white text-primary rounded-xl text-sm font-semibold hover:bg-white/90 transition-all disabled:opacity-40 disabled:cursor-not-allowed"
              >
                {copied ? (
                  <><CheckCircle2 className="w-4 h-4 text-secondary" /> Kopyalandı!</>
                ) : (
                  <><Share2 className="w-4 h-4" /> Paylaş</>
                )}
              </button>
              <button
                onClick={handleDownload}
                disabled={!selected}
                className="w-full flex items-center justify-center gap-2 py-2.5 border border-white/20 rounded-xl text-sm font-medium hover:bg-white/10 transition-all disabled:opacity-40 disabled:cursor-not-allowed"
              >
                <Download className="w-4 h-4" /> Kodu İndir
              </button>
            </div>
          </div>

          {/* Stats Card */}
          <div className="bg-on-tertiary-container rounded-2xl p-6 text-white">
            <div className="flex items-center gap-2 mb-2">
              <TrendingUp className="w-4 h-4 text-secondary-container" />
              <span className="text-xs font-medium text-secondary-container">Toplam kayıt</span>
            </div>
            <p className="text-3xl font-bold">{diplomas.length}</p>
            <p className="text-sm text-white/60 mt-1">Blockchain'e kaydedildi</p>
          </div>
        </div>
      </div>

      {/* Hidden canvas for QR PNG download */}
      <div ref={qrCanvasRef} className="hidden" aria-hidden="true">
        {selected && verifyUrl && (
          <QRCodeCanvas
            value={verifyUrl}
            size={512}
            bgColor="#ffffff"
            fgColor="#131b2e"
            level="M"
          />
        )}
      </div>
    </div>
  );
}
