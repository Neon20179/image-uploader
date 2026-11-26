import { api } from "@/shared/api/client";

import type { UploadsResponse } from "../model/types";

export const fetchUploads = async (
  offset: number,
  limit: number,
  signal?: AbortSignal
) => {
  const res = await api.get<UploadsResponse>("/uploads", {
    params: { offset, limit },
    signal
  });
  return res.data;
};
