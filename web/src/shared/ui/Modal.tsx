"use client";

import {
  type HTMLAttributes,
  type MouseEvent,
  type SyntheticEvent,
  memo,
  useEffect,
  useRef
} from "react";

import { IoCloseOutline } from "react-icons/io5";

import { cn } from "../lib";
import Button from "./Button";

interface Props extends Omit<
  HTMLAttributes<HTMLDialogElement>,
  "title" | "onClose"
> {
  isOpen: boolean;
  onClose: () => void;
  backdropClassName?: string;
  closeButtonClassName?: string;
  showCloseButton?: boolean;
}

const Modal = ({
  isOpen,
  onClose,
  children,
  className,
  backdropClassName,
  closeButtonClassName,
  showCloseButton = false,
  ...props
}: Props) => {
  const dialogRef = useRef<HTMLDialogElement>(null);

  useEffect(() => {
    const dialog = dialogRef.current;
    if (!dialog) return;
    if (!isOpen) {
      dialog.close();
      return;
    }
    dialog.showModal();
    const previousOverflow = document.body.style.overflow;
    document.body.style.overflow = "hidden";
    return () => {
      dialog.close();
      document.body.style.overflow = previousOverflow;
    };
  }, [isOpen]);

  const handleCancel = (event: SyntheticEvent<HTMLDialogElement>) => {
    event.preventDefault();
    onClose();
  };

  const handleBackdropClick = (event: MouseEvent<HTMLDialogElement>) => {
    if (event.target !== event.currentTarget) return;
    const rect = event.currentTarget.getBoundingClientRect();
    if (
      event.clientX < rect.left ||
      event.clientX > rect.right ||
      event.clientY < rect.top ||
      event.clientY > rect.bottom
    )
      onClose();
  };

  return (
    <dialog
      ref={dialogRef}
      className={cn(
        "m-auto max-h-[calc(100dvh-2.5rem)] w-[min(500px,calc(100%-2rem))] overflow-y-auto rounded-[0.875rem] border border-emerald-100 bg-white p-7 text-zinc-800 shadow-2xl shadow-zinc-800/15 backdrop:bg-zinc-800/35 backdrop:backdrop-blur-[0.25rem] mobile:p-5.5",
        backdropClassName,
        className
      )}
      onCancel={handleCancel}
      onClick={handleBackdropClick}
      {...props}
    >
      {showCloseButton && (
        <Button
          variant="ghost"
          className={cn(
            "absolute top-7 right-7 mobile:top-5.5 mobile:right-5.5",
            closeButtonClassName
          )}

          onClick={onClose}
        >
          <IoCloseOutline size={20} className="shrink-0" />
        </Button>
      )}
      {children}
    </dialog>
  );
};

export default memo(Modal);
