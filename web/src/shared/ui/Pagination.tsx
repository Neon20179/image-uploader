import { FiChevronLeft, FiChevronRight } from "react-icons/fi";

import Button from "./Button";

interface Props {
  page: number;
  pageSize: number;
  totalCount: number | undefined;
  disabled: boolean;
  onPageChange: (page: number) => void;
  onPageSizeChange: (size: number) => void;
}

const Pagination = ({
  page,
  pageSize,
  totalCount,
  disabled,
  onPageChange,
  onPageSizeChange
}: Props) => {
  const pageCount =
    totalCount === undefined
      ? undefined
      : Math.max(1, Math.ceil(totalCount / pageSize));
  const firstItem = totalCount === 0 ? 0 : page * pageSize + 1;
  const lastItem = Math.min((page + 1) * pageSize, totalCount ?? 0);

  const renderPageSizeOption = (size: number) => (
    <option key={size} value={size}>
      {size}
    </option>
  );

  return (
    <div className="flex min-h-[69px] items-center justify-between gap-4 border-t border-zinc-100 px-5.5 py-4 text-2xs mobile:flex-wrap mobile:p-3.5">
      <span className="text-zinc-400">
        {totalCount === undefined
          ? disabled
            ? "Loading…"
            : "Total unavailable"
          : totalCount > 0
            ? `Showing ${firstItem}–${lastItem} of ${totalCount}`
            : "No records"}
      </span>
      <div className="flex items-center gap-6 mobile:ml-auto mobile:gap-4">
        <label className="flex items-center gap-2.5 text-zinc-400">
          Rows per page{" "}
          <select
            className="rounded border border-zinc-200 bg-white px-2 py-1 text-2xs text-zinc-600"
            value={pageSize}
            disabled={disabled}
            onChange={e => onPageSizeChange(Number(e.target.value))}
          >
            {[10, 20, 50].map(renderPageSizeOption)}
          </select>
        </label>
        <div className="flex items-center gap-1.5">
          <Button
            className="border-zinc-200 p-1"
            variant="secondary"

            disabled={disabled || page === 0}
            onClick={() => onPageChange(page - 1)}
          >
            <FiChevronLeft size={16} className="shrink-0" />
          </Button>
          <span className="grid h-7 min-w-7 place-items-center rounded border border-emerald-100 bg-zinc-100 px-1.5 text-xs whitespace-nowrap text-emerald-600">
            {page + 1}
            {pageCount !== undefined && ` / ${pageCount}`}
          </span>
          <Button
            className="border-zinc-200 p-1"
            variant="secondary"

            disabled={
              disabled || pageCount === undefined || page + 1 >= pageCount
            }
            onClick={() => onPageChange(page + 1)}
          >
            <FiChevronRight size={16} className="shrink-0" />
          </Button>
        </div>
      </div>
    </div>
  );
};

export default Pagination;
