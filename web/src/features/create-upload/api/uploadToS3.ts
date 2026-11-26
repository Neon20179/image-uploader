import axios from "axios";

import { CreateUploadResponse } from "../model/types";

export const uploadToS3 = async (
  file: File,
  createUploadResponse: CreateUploadResponse,
  onUploadProgress?: (progress: number) => void
) => {
  const form = new FormData();
  Object.entries(createUploadResponse.fields).forEach(([key, value]) =>
    form.append(key, value)
  );
  if (
    !Object.keys(createUploadResponse.fields).some(
      key => key.toLowerCase() === "content-type"
    )
  ) {
    form.append("Content-Type", file.type);
  }
  form.append("file", file);
  await axios.post(createUploadResponse.uploadUrl, form, {
    timeout: 120_000,
    onUploadProgress: onUploadProgress
      ? event =>
          onUploadProgress(
            Math.min(99, Math.round((event.progress ?? 0) * 100))
          )
      : undefined
  });
};
