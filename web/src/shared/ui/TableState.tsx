import type { ReactNode } from "react";

import { FiAlertCircle, FiFolder, FiLoader } from "react-icons/fi";

import { getErrorMessage } from "@/shared/api/client";

import Button from "./Button";
import TableCell from "./TableCell";

interface Props {
  columns: number;
  loading?: boolean;
  error?: unknown;
  onRetry?: () => void;
  children?: ReactNode;
  compact?: boolean;
}

const TableState = ({
  columns,
  loading,
  error,
  onRetry,
  children,
  compact
}: Props) => {
  return (
    <tr className="border-b border-zinc-100 last:border-b-0 hover:bg-zinc-50">
      <TableCell colSpan={columns} compact={compact}>
        <div
          className="flex min-h-[330px] flex-col items-center justify-center px-5 py-12.5 text-center whitespace-normal mobile:min-h-[300px]"
          role={error ? "alert" : "status"}
        >
          <span className="mb-5 grid size-14.5 place-items-center rounded-2xl border border-zinc-200 bg-zinc-100 text-zinc-400">
            {loading ? (
              <FiLoader size={24} className="animate-spin" />
            ) : error ? (
              <FiAlertCircle size={24} />
            ) : (
              <FiFolder size={24} />
            )}
          </span>
          <strong className="text-base font-medium text-zinc-600">
            {loading
              ? "Loading data"
              : error
                ? "Could not load the list"
                : "Nothing here yet"}
          </strong>
          <p className="mt-2.5 mb-4 max-w-[390px] text-xs leading-[1.9] text-zinc-400">
            {loading
              ? "This will only take a few seconds"
              : error
                ? getErrorMessage(error)
                : children}
          </p>
          {!!error && (
            <Button variant="secondary" onClick={onRetry}>
              Try again
            </Button>
          )}
        </div>
      </TableCell>
    </tr>
  );
};

export default TableState;
