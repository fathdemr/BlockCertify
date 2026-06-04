import type { DiplomaStatus } from '../types';
import { CheckCircle, Clock, XCircle } from 'lucide-react';

interface StatusBadgeProps {
  status: DiplomaStatus;
}

const config = {
  approved: {
    label: 'Onaylandı',
    icon: CheckCircle,
    className: 'bg-emerald-50 text-emerald-700 border border-emerald-200',
  },
  pending: {
    label: 'Beklemede',
    icon: Clock,
    className: 'bg-amber-50 text-amber-700 border border-amber-200',
  },
  rejected: {
    label: 'Reddedildi',
    icon: XCircle,
    className: 'bg-red-50 text-red-700 border border-red-200',
  },
};

export default function StatusBadge({ status }: Readonly<StatusBadgeProps>) {
  const { label, icon: Icon, className } = config[status];
  return (
    <span className={`inline-flex items-center gap-1 px-2.5 py-1 rounded-full text-xs font-medium ${className}`}>
      <Icon className="w-3 h-3" />
      {label}
    </span>
  );
}
