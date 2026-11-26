"use client";
import { FC, MouseEvent, ReactNode, useEffect } from "react";

import { AnimatePresence, motion } from "framer-motion";

interface Props {
  isOpen: boolean;
  onClose?: () => void;
  closeOnBackdropClick?: boolean;
  className?: string;
  children?: ReactNode;
  onClick?: (e: MouseEvent<HTMLDivElement>) => void;
}

const Backdrop: FC<Props> = ({
  isOpen,
  onClose,
  closeOnBackdropClick = true,
  className,
  onClick,
  children
}) => {
  const handleBackdropClick = (e: MouseEvent<HTMLDivElement>) => {
    onClick?.(e);
    if (closeOnBackdropClick) onClose?.();
  };

  useEffect(() => {
    if (!isOpen) return;
    const previousOverflow = document.body.style.overflow;
    document.body.style.overflow = "hidden";
    return () => {
      document.body.style.overflow = previousOverflow;
    };
  }, [isOpen]);

  return (
    <AnimatePresence mode="wait" initial={false}>
      {isOpen && (
        <>
          <motion.span
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            exit={{ opacity: 0 }}
            transition={{ duration: 0.35, ease: "easeInOut" }}
            className="fixed inset-0 z-50 bg-zinc-800/35 backdrop-blur-[0.25rem]"
            onClick={handleBackdropClick}
          />
          <motion.div
            initial={{ y: 20, opacity: 0 }}
            animate={{ y: 0, opacity: 1 }}
            exit={{ y: -20, opacity: 0 }}
            transition={{ duration: 0.35, ease: "easeInOut" }}
            className={`fixed inset-0 z-50 flex items-center justify-center px-4 py-5 ${className ? " " + className : ""}`}
          >
            <span className="absolute inset-0" onClick={handleBackdropClick} />
            <div className="relative">{children}</div>
          </motion.div>
        </>
      )}
    </AnimatePresence>
  );
};

export default Backdrop;
