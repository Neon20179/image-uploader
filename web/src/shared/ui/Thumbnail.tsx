"use client";
import { useState } from "react";

import { FiImage } from "react-icons/fi";

const Thumbnail = ({ url, alt }: { url: string; alt: string }) => {
  const [failedUrl, setFailedUrl] = useState<string>();

  return (
    <span className="grid h-11 w-11.5 shrink-0 place-items-center overflow-hidden rounded-md border border-zinc-200 bg-zinc-100 text-zinc-400">
      {url && failedUrl !== url ? (
        // eslint-disable-next-line @next/next/no-img-element
        <img
          className="size-full object-cover"
          src={url}
          alt={alt}
          loading="lazy"
          onError={() => setFailedUrl(url)}
        />
      ) : (
        <FiImage size={20} />
      )}
    </span>
  );
};

export default Thumbnail;
