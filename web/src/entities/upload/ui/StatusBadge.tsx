import type { IconType } from "react-icons";
import {
  FiCheck,
  FiCircle,
  FiClock,
  FiLoader,
  FiShield,
  FiX
} from "react-icons/fi";

import { cn } from "@/shared/lib/cn";

import { ProcessStatus, UploadStatus } from "../model/types";

const statuses: {
  [key: string]: {
    label: string;
    color: string;
    icon: IconType;
  };
} = {
  pending: { label: "Pending", color: "neutral", icon: FiClock },
  uploaded: { label: "Uploaded", color: "blue", icon: FiCheck },
  processing: { label: "Processing", color: "amber", icon: FiLoader },
  ready: { label: "Ready", color: "green", icon: FiCheck },
  blocked: { label: "Blocked", color: "red", icon: FiShield },
  failed: { label: "Failed", color: "red", icon: FiX },
  completed: { label: "Completed", color: "green", icon: FiCheck }
};

interface Props {
  status: UploadStatus | ProcessStatus;
  compact?: boolean;
}

const colorClasses: Record<string, string> = {
  neutral: "border-zinc-200 bg-zinc-50 text-zinc-400",
  green: "border-emerald-100 bg-emerald-50 text-emerald-500",
  amber: "border-amber-100 bg-amber-50 text-amber-700",
  blue: "border-sky-100 bg-sky-50 text-sky-700",
  red: "border-red-100 bg-red-50 text-red-600"
};

const StatusBadge = ({ status, compact = false }: Props) => {
  const item = statuses[status] ?? {
    label: status,
    color: "neutral",
    icon: FiCircle
  };

  const Icon = item.icon;

  return (
    <span
      className={cn(
        "inline-flex items-center gap-1 rounded-sm border px-2 py-1 text-3xs",
        colorClasses[item.color],
        compact && "border-transparent bg-transparent pl-0"
      )}
    >
      <Icon size={12} />
      {item.label}
    </span>
  );
};

export default StatusBadge;
