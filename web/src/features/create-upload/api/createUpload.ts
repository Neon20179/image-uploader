import { api } from "@/shared/api/client";

import { CreateUploadResponse } from "../model/types";

export const createUpload = async () => {
  const res = await api.post<CreateUploadResponse>("/uploads");
  return res.data;
};
