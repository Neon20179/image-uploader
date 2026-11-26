export type UploadStatus =
  "pending" | "uploaded" | "processing" | "ready" | "blocked" | "failed";

export type ProcessStatus = "pending" | "completed" | "failed";

export interface UploadRecord {
  id: string;
  url: string;
  status: UploadStatus;
  normalizationStatus: ProcessStatus;
  embeddingStatus: ProcessStatus;
  ocrStatus: ProcessStatus;
  createdAt: string;
  updatedAt: string;
}

export interface UploadsResponse {
  uploads: UploadRecord[];
  count: number;
}
