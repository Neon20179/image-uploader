import type { HTMLAttributes } from "react";

import { cn } from "../lib";

const TableRow = ({
  className,
  ...props
}: HTMLAttributes<HTMLTableRowElement>) => (
  <tr
    className={cn(
      "border-b border-zinc-100 last:border-b-0 hover:bg-zinc-50",
      className
    )}
    {...props}
  />
);

export default TableRow;
