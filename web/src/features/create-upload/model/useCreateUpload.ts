import { useRef, useState } from "react";

import { useMutation, useQueryClient } from "@tanstack/react-query";

import { createUpload } from "../api/createUpload";
import { uploadToS3 } from "../api/uploadToS3";
import { CreateUploadResponse } from "./types";

export const useCreateUpload = () => {
  const [progress, setProgress] = useState(0);
  const uploadRef = useRef<{ file: File; upload: CreateUploadResponse } | null>(
    null
  );

  const client = useQueryClient();

  const mutation = useMutation({
    mutationFn: async (file: File) => {
      setProgress(0);

      if (uploadRef.current?.file !== file) {
        uploadRef.current = { file, upload: await createUpload() };
      }

      const { upload } = uploadRef.current;
      await uploadToS3(file, upload, setProgress);

      setProgress(100);
      uploadRef.current = null;

      return upload.uploadId;
    },
    onSettled: () => {
      void client.invalidateQueries({ queryKey: ["uploads"] });
      void client.invalidateQueries({ queryKey: ["images"] });
    }
  });

  return { ...mutation, progress };
};
