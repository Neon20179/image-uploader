import { UploadStatus } from "@/entities/upload/model/types";

export interface CreateUploadResponse {
  uploadId: string;
  status: UploadStatus;
  uploadUrl: string;
  fields: Record<string, string>;
}
