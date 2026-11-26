"use client";

import {
  type ChangeEvent,
  type DragEvent,
  type FormEvent,
  useRef,
  useState
} from "react";

import {
  FiCheck,
  FiFile,
  FiInfo,
  FiLoader,
  FiPlus,
  FiUploadCloud,
  FiX
} from "react-icons/fi";

import { getErrorMessage } from "@/shared/api/client";
import { formatBytes } from "@/shared/lib/format";
import { Button, Modal } from "@/shared/ui";

import { useCreateUpload } from "../model/useCreateUpload";

const minBytes = Number(process.env.NEXT_PUBLIC_UPLOAD_MIN_BYTES ?? 1024);
const maxBytes = Number(
  process.env.NEXT_PUBLIC_UPLOAD_MAX_BYTES ?? 20 * 1024 * 1024
);

interface Props {
  onCreated: () => void;
}

const CreateUpload = ({ onCreated }: Props) => {
  const [isModalOpen, setIsModalOpen] = useState(false);

  const {
    mutate: createUpload,
    isPending,
    isError,
    error,
    progress,
    reset
  } = useCreateUpload();
  const input = useRef<HTMLInputElement>(null);

  const [file, setFile] = useState<File | null>(null);
  const [dragging, setDragging] = useState(false);
  const [validation, setValidation] = useState("");

  const handleOpenModal = () => {
    reset();
    setFile(null);
    setDragging(false);
    setValidation("");
    setIsModalOpen(true);
  };

  const handleCloseModal = () => {
    if (isPending) return;
    setIsModalOpen(false);
  };

  const handleSelectFile = (selected?: File) => {
    reset();
    setValidation("");
    setFile(null);
    if (!selected) return;
    if (!["image/jpeg", "image/png"].includes(selected.type)) {
      setValidation("Choose a JPEG or PNG image.");
    } else if (selected.size < minBytes || selected.size > maxBytes) {
      setValidation(
        `File size must be between ${formatBytes(minBytes)} and ${formatBytes(maxBytes)}.`
      );
    } else {
      setFile(selected);
    }
  };

  const handleSubmit = (event: FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    if (!file || validation || isPending) return;
    createUpload(file, {
      onSuccess: () => {
        setIsModalOpen(false);
        onCreated();
      }
    });
  };

  const handleFileChange = (event: ChangeEvent<HTMLInputElement>) => {
    handleSelectFile(event.target.files?.[0]);
    event.target.value = "";
  };

  const handleDragOver = (event: DragEvent<HTMLButtonElement>) => {
    event.preventDefault();
    if (!isPending) setDragging(true);
  };

  const handleDrop = (event: DragEvent<HTMLButtonElement>) => {
    event.preventDefault();
    setDragging(false);
    if (!isPending) {
      if (event.dataTransfer.files.length > 1) {
        reset();
        setFile(null);
        setValidation("You can upload one image at a time.");
      } else handleSelectFile(event.dataTransfer.files[0]);
    }
  };

  return (
    <>
      <Button
        variant="primary"
        className="mobile:self-stretch"
        onClick={handleOpenModal}
      >
        <FiPlus size={17} className="shrink-0" />
        New upload
      </Button>

      <Modal isOpen={isModalOpen} onClose={handleCloseModal}>
        <form onSubmit={handleSubmit}>
          <div className="flex items-start justify-between">
            <span className="grid size-12 place-items-center rounded-3.5 border border-emerald-100 bg-zinc-100 text-emerald-500">
              <FiUploadCloud size={23} />
            </span>
            <Button
              variant="ghost"

              disabled={isPending}
              onClick={handleCloseModal}
            >
              <FiX size={20} className="shrink-0" />
            </Button>
          </div>
          <h2
            id="upload-title"
            className="mt-5 mb-2.5 text-2xl font-medium tracking-[-0.0375rem]"
          >
            New upload
          </h2>
          <p
            id="upload-description"
            className="mt-3 mb-6 text-xs leading-[1.6] text-zinc-500"
          >
            Add an image to your media library. After uploading, it will be
            processed automatically.
          </p>
          <input
            ref={input}
            className="sr-only"
            type="file"

            tabIndex={-1}
            accept="image/jpeg,image/png"
            disabled={isPending}
            onChange={handleFileChange}
          />
          <button
            type="button"
            className={`flex w-full flex-col items-center rounded-lg border-[0.09375rem] border-dashed border-emerald-200 bg-zinc-50 px-4.5 py-7 text-zinc-500 enabled:hover:border-emerald-500 enabled:hover:bg-emerald-50 ${dragging ? "border-emerald-500! bg-emerald-50!" : ""}`}
            disabled={isPending}
            onClick={() => input.current?.click()}
            onDragOver={handleDragOver}
            onDragLeave={() => setDragging(false)}
            onDrop={handleDrop}
          >
            <span className="mb-4 text-zinc-500">
              {file ? <FiFile size={28} /> : <FiUploadCloud size={28} />}
            </span>
            <strong className="max-w-full text-xs font-medium [overflow-wrap:anywhere]">
              {file ? file.name : "Choose a file or drag it here"}
            </strong>
            <span className="mt-2 text-2xs text-zinc-400">
              {file
                ? `${formatBytes(file.size)} · Click to replace`
                : `JPEG or PNG · up to ${formatBytes(maxBytes)}`}
            </span>
            {!file && (
              <span className="mt-4.5 rounded-sm border border-emerald-100 bg-white px-3.5 py-2 text-2xs text-zinc-600">
                Choose file
              </span>
            )}
            {file && (
              <span className="mt-2 flex items-center gap-1 text-2xs text-emerald-500">
                <FiCheck size={14} />
                Ready to upload
              </span>
            )}
          </button>
          {(validation || isError) && (
            <p
              className="mt-3 rounded-sm bg-red-50 px-3 py-2.5 text-xs leading-[1.6] text-zinc-500"
              role="alert"
            >
              {validation || getErrorMessage(error)}
            </p>
          )}
          {isPending && (
            <div className="mt-4" role="status">
              <div className="mb-2 flex justify-between text-xs text-zinc-500">
                <span>
                  {progress < 99
                    ? "Uploading image…"
                    : "Waiting for storage to respond…"}
                </span>
                <strong>{progress}%</strong>
              </div>
              <progress
                className="block h-1.5 w-full appearance-none overflow-hidden rounded-sm border-0 accent-emerald-600 [&::-webkit-progress-bar]:bg-zinc-100 [&::-webkit-progress-value]:bg-emerald-500 [&::-webkit-progress-value]:transition-[width] [&::-webkit-progress-value]:duration-200 [&::-moz-progress-bar]:bg-emerald-500"
                value={progress}
                max={100}
              />
            </div>
          )}
          <div className="mt-4.5 flex items-start gap-2 text-2xs leading-[1.8] text-zinc-400">
            <FiInfo size={16} className="mt-px shrink-0" />
            <span>
              The image will appear in your library once processing and
              validation are complete.
            </span>
          </div>
          <div className="mt-6 flex justify-end gap-2 border-t border-zinc-100 pt-5">
            <Button
              className="mobile:px-3 mobile:py-2 mobile:text-xs"
              disabled={isPending}
              onClick={handleCloseModal}
            >
              Cancel
            </Button>
            <Button
              type="submit"
              className="mobile:px-3 mobile:py-2 mobile:text-xs"
              variant="primary"
              disabled={!file || !!validation || isPending}
            >
              {isPending ? (
                <FiLoader size={16} className="animate-spin shrink-0" />
              ) : (
                <FiUploadCloud size={16} className="shrink-0" />
              )}
              {isPending
                ? "Uploading…"
                : isError
                  ? "Retry upload"
                  : "Upload image"}
            </Button>
          </div>
        </form>
      </Modal>
    </>
  );
};

export default CreateUpload;
