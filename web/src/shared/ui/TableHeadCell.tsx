import type { ThHTMLAttributes } from "react";

import { cn } from "../lib";

interface Props extends ThHTMLAttributes<HTMLTableCellElement> {
  compact?: boolean;
}

const TableHeadCell = ({ compact = false, className, ...props }: Props) => (
  <th
    scope="col"
    className={cn(
      "border-y border-zinc-100 bg-zinc-50 py-3 text-2xs font-medium text-zinc-400 first:pl-6",
      compact ? "px-4" : "px-5",
      className
    )}
    {...props}
  />
);

export default TableHeadCell;
