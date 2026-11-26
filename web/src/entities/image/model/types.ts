export interface ImageRecord {
  id: string;
  url: string;
  width: number;
  height: number;
  size: number;
  mimeType: string;
  createdAt: string;
}

export interface ImagesResponse {
  images: ImageRecord[];
  count: number;
}
