"use client";

import { useState } from "react";

import Link from "next/link";

import { FiArrowUpRight, FiImage, FiRefreshCw } from "react-icons/fi";

import { useImages } from "@/entities/image";
import { formatBytes, formatDate } from "@/shared/lib/format";
import {
  Button,
  Pagination,
  TableCell,
  TableHeadCell,
  TableState,
  Thumbnail
} from "@/shared/ui";

const ImagesPage = () => {
  const [page, setPage] = useState(0);
  const [pageSize, setPageSize] = useState(10);

  const query = useImages({ offset: page * pageSize, limit: pageSize });

  const images = query.data?.images ?? [];
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
            Images
          </h1>
          <div className="mt-5 flex items-center gap-2 text-2xs leading-[1.7] text-zinc-400">
            <span className="size-1 shrink-0 rounded-full bg-zinc-400" />
            The library displays images that have been processed successfully.
          </div>
        </div>
      </div>

      <section className="overflow-hidden rounded-[0.625rem] border border-zinc-200 bg-white shadow-xs shadow-zinc-800/5">
        <div className="flex min-h-[70px] items-center justify-between gap-3 px-5.5 py-4.5 mobile:p-3.5">
          <div className="flex items-center gap-2.5 text-zinc-500">
            <FiImage size={18} />
            <h2 className="m-0 text-xs font-semibold text-zinc-700 mobile:text-xs">
              All images
            </h2>
          </div>
          <Button
            variant="ghost"
            disabled={query.isFetching}
            onClick={() => query.refetch()}
          >
            <FiRefreshCw
              size={15}
              className={"shrink-0 " + (query.isFetching ? "animate-spin" : "")}
            />
            Refresh
          </Button>
        </div>
        <div className="relative overflow-x-auto">
          <table className="w-full border-collapse text-left whitespace-nowrap">
            <thead>
              <tr>
                <TableHeadCell scope="col">Image</TableHeadCell>
                <TableHeadCell scope="col">Format</TableHeadCell>
                <TableHeadCell scope="col">Dimensions</TableHeadCell>
                <TableHeadCell scope="col">Size</TableHeadCell>
                <TableHeadCell scope="col">Date added</TableHeadCell>
                <TableHeadCell scope="col">
                  <span className="sr-only">Actions</span>
                </TableHeadCell>
              </tr>
            </thead>

            <tbody>
              {query.isPending ? (
                <TableState columns={6} loading />
              ) : query.isError ? (
                <TableState
                  columns={6}
                  error={query.error}
                  onRetry={() => void query.refetch()}
                />
              ) : images.length === 0 ? (
                <TableState columns={6}>
                  Create an upload to see your processed images here.
                  <br />
                  <Link
                    href="/uploads"
                    className="mt-2 inline-flex items-center gap-1 text-emerald-500"
                  >
                    Go to uploads <FiArrowUpRight size={14} />
                  </Link>
                </TableState>
              ) : (
                images.map(image => (
                  <tr
                    key={image.id}
                    className="border-b border-zinc-100 last:border-b-0 hover:bg-zinc-50"
                  >
                    <TableCell>
                      <div className="flex items-center gap-3">
                        <Thumbnail url={image.url} alt={`Image ${image.id}`} />
                        <div>
                          <a
                            className="text-xs font-semibold text-zinc-700 hover:text-emerald-600 hover:underline"
                            href={image.url}
                            target="_blank"
                            rel="noreferrer"
                          >
                            {image.id.slice(0, 8)}
                          </a>
                          <span
                            className="mt-1.5 block max-w-[155px] overflow-hidden font-mono text-3xs text-ellipsis text-zinc-400"
                            title={image.id}
                          >
                            {image.id}
                          </span>
                        </div>
                      </div>
                    </TableCell>
                    <TableCell>
                      <span className="rounded border border-zinc-200 bg-zinc-50 px-2 py-1 text-3xs font-semibold tracking-[0.01875rem]">
                        {image.mimeType.replace("image/", "").toUpperCase()}
                      </span>
                    </TableCell>
                    <TableCell className="tabular-nums">
                      {image.width} <span className="text-zinc-400">×</span>{" "}
                      {image.height}
                      <span className="ml-1 text-2xs text-zinc-400">px</span>
                    </TableCell>
                    <TableCell className="tabular-nums">
                      {formatBytes(image.size)}
                    </TableCell>
                    <TableCell className="text-2xs text-zinc-500">
                      {formatDate(image.createdAt)}
                    </TableCell>
                    <TableCell>
                      <a
                        className="flex size-7 items-center justify-center rounded text-zinc-400 hover:bg-zinc-100 hover:text-emerald-600"

                        href={image.url}
                        target="_blank"
                        rel="noreferrer"
                      >
                        <FiArrowUpRight size={17} />
                      </a>
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

export default ImagesPage;
