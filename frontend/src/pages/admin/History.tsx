import { useEffect, useState } from 'react';
import { diplomaApi } from '../../services/api';
import type { Diploma } from '../../types';
import StatusBadge from '../../components/StatusBadge';
import { ArrowRight } from 'lucide-react';

export default function History() {
  const [diplomas, setDiplomas] = useState<Diploma[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    diplomaApi.getAll()
      .then((r) => setDiplomas(r.data ?? []))
      .catch(() => {})
      .finally(() => setLoading(false));
  }, []);

  return (
    <div className="p-8">
      <div className="mb-8">
        <h1 className="text-2xl font-bold text-on-surface">Diploma Geçmişi</h1>
        <p className="text-sm text-on-surface-variant mt-1">Tüm diploma kayıtlarını görüntüleyin.</p>
      </div>

      <div className="bg-white rounded-2xl border border-outline-variant/40 card-shadow overflow-hidden">
        {loading ? (
          <div className="flex items-center justify-center py-20">
            <div className="w-6 h-6 border-2 border-secondary border-t-transparent rounded-full animate-spin" />
          </div>
        ) : diplomas.length === 0 ? (
          <div className="flex flex-col items-center justify-center py-20 text-on-surface-variant">
            <p className="text-sm">Henüz kayıt yok.</p>
          </div>
        ) : (
          <div className="overflow-x-auto">
            <table className="w-full">
              <thead>
                <tr className="bg-surface text-xs text-on-surface-variant uppercase tracking-wider">
                  {['Diploma No', 'İsim', 'Program', 'Üniversite', 'Yıl', 'Durum', 'TX Hash', ''].map((h) => (
                    <th key={h} className="px-6 py-3 text-left font-medium">{h}</th>
                  ))}
                </tr>
              </thead>
              <tbody className="divide-y divide-outline-variant/20">
                {diplomas.map((d) => (
                  <tr key={d.id} className="hover:bg-surface/50 transition-colors">
                    <td className="px-6 py-4 text-xs font-mono-code text-on-surface-variant">
                      #{d.diploma_no ?? String(d.id).padStart(4, '0')}
                    </td>
                    <td className="px-6 py-4 font-semibold text-sm text-on-surface">
                      {d.first_name} {d.last_name}
                    </td>
                    <td className="px-6 py-4 text-sm text-on-surface-variant">{d.department}</td>
                    <td className="px-6 py-4 text-sm text-on-surface-variant">{d.university}</td>
                    <td className="px-6 py-4 text-sm text-on-surface-variant">{d.graduation_year}</td>
                    <td className="px-6 py-4"><StatusBadge status={d.status} /></td>
                    <td className="px-6 py-4 text-xs font-mono-code text-outline max-w-[160px] truncate">
                      {d.tx_hash ?? '—'}
                    </td>
                    <td className="px-6 py-4">
                      <button className="text-xs text-on-tertiary-container font-medium flex items-center gap-1 hover:gap-2 transition-all">
                        Detay <ArrowRight className="w-3 h-3" />
                      </button>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        )}
      </div>
    </div>
  );
}
