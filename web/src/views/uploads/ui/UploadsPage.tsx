"use client";

import { useState } from "react";

import { FiCheckCircle, FiRefreshCw, FiUploadCloud, FiX } from "react-icons/fi";

import { StatusBadge, useUploads } from "@/entities/upload";
import { CreateUpload } from "@/features/create-upload";
import { formatDate } from "@/shared/lib/format";
import {
  Button,
  Pagination,
  TableCell,
  TableHeadCell,
  TableState
} from "@/shared/ui";

const UploadsPage = () => {
  const [page, setPage] = useState(0);
  const [pageSize, setPageSize] = useState(10);

  const [created, setCreated] = useState(false);

  const query = useUploads({ offset: page * pageSize, limit: pageSize });

  const handleUploadCreated = () => {
    setPage(0);
    setCreated(true);
  };

  const uploads = query.data?.uploads ?? [];
  const totalCount = query.data?.count;
  const lastPage = Math.max(0, Math.ceil((totalCount ?? 0) / pageSize) - 1);

  if (query.isSuccess && !query.isPlaceholderData && page > lastPage) {
    setPage(lastPage);
  }

  const handlePageSizeChange = (size: number) => {
    setPageSize(size);
    setPage(0);
  };

  return (
    <>
      <div className="mb-9 flex items-center justify-between gap-5 mobile:mb-6 mobile:flex-col mobile:items-start">
        <div>
          <h1 className="m-0 text-page-title leading-[1.2] font-medium tracking-[-0.0875rem]">
            Uploads
          </h1>
        </div>
        <CreateUpload onCreated={handleUploadCreated} />
      </div>

      {created && (
        <div
          className="mb-5 flex items-center gap-2.5 rounded-lg border border-emerald-100 bg-emerald-50 px-4 py-2.5 text-xs text-emerald-800"
          role="status"
        >
          <FiCheckCircle size={18} />
          <span>
            File uploaded. Processing status will update automatically.
          </span>
          <Button
            variant="ghost"
            className="ml-auto"

            onClick={() => setCreated(false)}
          >
            <FiX size={16} className="shrink-0" />
          </Button>
        </div>
      )}

      <section className="overflow-hidden rounded-[0.625rem] border border-zinc-200 bg-white shadow-xs shadow-zinc-800/5">
        <div className="flex min-h-[70px] items-center justify-between gap-3 px-5.5 py-4.5 mobile:p-3.5">
          <div className="flex items-center gap-2.5 text-zinc-500">
            <FiUploadCloud size={18} />
            <h2 className="m-0 text-xs font-semibold text-zinc-700 mobile:text-xs">
              All uploads
            </h2>
          </div>
          <div className="flex items-center gap-3.5 mobile:gap-1">
            <span className="flex items-center gap-1.5 text-2xs text-zinc-400 mobile:text-3xs">
              <span className="size-1 rounded-full bg-emerald-300" />
              Auto-refresh
            </span>
            <Button
              variant="ghost"
              disabled={query.isFetching}
              onClick={() => query.refetch()}
            >
              <FiRefreshCw
                size={15}
                className={
                  "shrink-0 " + (query.isFetching ? "animate-spin" : "")
                }
              />
            </Button>
          </div>
        </div>
        <div className="relative overflow-x-auto">
          <table className="w-full border-collapse text-left whitespace-nowrap">
            <thead>
              <tr>
                <TableHeadCell compact scope="col">
                  Upload
                </TableHeadCell>
                <TableHeadCell compact scope="col">
                  Status
                </TableHeadCell>
                <TableHeadCell compact scope="col">
                  Normalization
                </TableHeadCell>
                <TableHeadCell compact scope="col">
                  Embedding
                </TableHeadCell>
                <TableHeadCell compact scope="col">
                  OCR
                </TableHeadCell>
                <TableHeadCell compact scope="col">
                  Created / updated
                </TableHeadCell>
              </tr>
            </thead>

            <tbody>
              {query.isPending ? (
                <TableState compact columns={6} loading />
              ) : query.isError ? (
                <TableState
                  compact
                  columns={6}
                  error={query.error}
                  onRetry={() => void query.refetch()}
                />
              ) : uploads.length === 0 ? (
                <TableState compact columns={6}>
                  Your upload history will appear here.
                  <br />
                  Click “New upload” to add your first image.
                </TableState>
              ) : (
                uploads.map(upload => (
                  <tr
                    key={upload.id}
                    className="border-b border-zinc-100 last:border-b-0 hover:bg-zinc-50"
                  >
                    <TableCell compact>
                      <div className="flex items-center gap-3">
                        <span className="grid h-10 w-9 place-items-center rounded-lg border border-zinc-200 bg-emerald-50 text-zinc-400">
                          <FiUploadCloud size={19} />
                        </span>
                        <div>
                          <span className="text-xs font-semibold text-zinc-700">
                            {upload.id.slice(0, 8)}
                          </span>
                          <span
                            className="mt-1.5 block max-w-[155px] overflow-hidden font-mono text-3xs text-ellipsis text-zinc-400"
                            title={upload.id}
                          >
                            {upload.id}
                          </span>
                        </div>
                      </div>
                    </TableCell>
                    <TableCell compact>
                      <StatusBadge status={upload.status} />
                    </TableCell>
                    <TableCell compact>
                      <StatusBadge
                        status={upload.normalizationStatus}
                        compact
                      />
                    </TableCell>
                    <TableCell compact>
                      <StatusBadge status={upload.embeddingStatus} compact />
                    </TableCell>
                    <TableCell compact>
                      <StatusBadge status={upload.ocrStatus} compact />
                    </TableCell>
                    <TableCell compact className="text-2xs text-zinc-500">
                      {formatDate(upload.createdAt)}
                      <span className="mt-2 block text-3xs text-zinc-400">
                        {formatDate(upload.updatedAt)}
                      </span>
                    </TableCell>
                  </tr>
                ))
              )}
            </tbody>
          </table>
        </div>
        <Pagination
          page={page}
          pageSize={pageSize}
          totalCount={totalCount}
          disabled={query.isFetching}
          onPageChange={setPage}
          onPageSizeChange={handlePageSizeChange}
        />
      </section>
    </>
  );
};

export default UploadsPage;
