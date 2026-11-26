"use client";
import type { ReactNode } from "react";

import Link from "next/link";
import { usePathname } from "next/navigation";

import { FiImage, FiLayers, FiUploadCloud } from "react-icons/fi";

const AppShell = ({ children }: { children: ReactNode }) => {
  const pathname = usePathname();
  return (
    <div className="flex min-h-dvh mobile:flex-col">
      <aside className="fixed inset-y-0 left-0 flex w-[242px] shrink-0 flex-col border-r border-zinc-200 bg-white px-5 pt-9 tablet:w-[210px] tablet:px-4 mobile:static mobile:w-full mobile:border-r-0 mobile:border-b mobile:px-5 mobile:pt-4.5">
        <Link
          href="/images"
          className="flex items-center gap-3 px-1.5 text-xl font-bold tracking-[-0.05rem] tablet:text-lg mobile:px-0 mobile:text-xl"
        >
          <span className="grid h-10 w-9 place-items-center rounded-xl bg-emerald-600 text-white">
            <FiLayers size={22} />
          </span>
          <span>
            image<span className="font-normal">uploader</span>
          </span>
        </Link>
        <div className="mx-3 mt-14 mb-4 text-3xs font-semibold tracking-[0.0875rem] text-zinc-400 mobile:hidden">
          WORKSPACE
        </div>
        <nav className="mobile:flex mobile:gap-2.5 mobile:pt-4 mobile:pb-2.5">
          <Link
            href="/images"
            className={`my-1 flex items-center gap-3 rounded-lg px-3.5 py-3 text-xs text-zinc-500 hover:bg-zinc-50 mobile:m-0 mobile:flex-1 mobile:justify-center mobile:p-2.5`}
          >
            <FiImage size={19} />
            Images
            {pathname === "/images" && (
              <span className="ml-auto size-1.5 rounded-full bg-emerald-500 mobile:ml-0" />
            )}
          </Link>

          <Link
            href="/uploads"
            className={`my-1 flex items-center gap-3 rounded-lg px-3.5 py-3 text-xs text-zinc-500 hover:bg-zinc-50 mobile:m-0 mobile:flex-1 mobile:justify-center mobile:p-2.5`}
          >
            <FiUploadCloud size={19} />
            Uploads
            {pathname === "/uploads" && (
              <span className="ml-auto size-1.5 rounded-full bg-emerald-500 mobile:ml-0" />
            )}
          </Link>
        </nav>
      </aside>

      <div className="ml-[242px] flex min-w-0 flex-1 flex-col tablet:ml-[210px] mobile:ml-0">
        <header className="flex h-[76px] items-center justify-between border-b border-zinc-200 bg-white/70 px-11 tablet:px-6 mobile:h-13 mobile:px-5">
          <div className="flex items-center gap-4 text-xs text-zinc-500">
            Library <span className="text-zinc-300">/</span>{" "}
            <strong className="font-medium text-zinc-600">
              {pathname === "/uploads" ? "Uploads" : "Images"}
            </strong>
          </div>
        </header>
        <main
          id="main-content"
          className="mx-auto w-full max-w-[1600px] flex-1 px-11 py-12 min-[1600px]:pt-[65px] tablet:px-6 tablet:py-9 mobile:px-5 mobile:py-7.5"
        >
          {children}
        </main>
      </div>
    </div>
  );
};

export default AppShell;
