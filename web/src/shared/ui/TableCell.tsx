import type { TdHTMLAttributes } from "react";

import { cn } from "../lib";

interface Props extends TdHTMLAttributes<HTMLTableCellElement> {
  compact?: boolean;
}

const TableCell = ({ compact = false, className, ...props }: Props) => (
  <td
    className={cn(
      "h-[82px] py-4 text-xs text-zinc-600 first:pl-6",
      compact ? "px-4" : "px-5",
      className
    )}
    {...props}
  />
);

export default TableCell;
