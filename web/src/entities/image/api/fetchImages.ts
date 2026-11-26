import { api } from "@/shared/api/client";

import type { ImagesResponse } from "../model/types";

export const fetchImages = async (
  offset: number,
  limit: number,
  signal?: AbortSignal
) => {
  const res = await api.get<ImagesResponse>("/images", {
    params: { offset, limit },
    signal
  });
  return res.data;
};
