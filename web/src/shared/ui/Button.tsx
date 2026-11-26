"use client";

import { type ButtonHTMLAttributes, forwardRef, memo } from "react";

import { cn } from "../lib";

type ButtonVariant =
  | "primary"
  | "secondary"
  | "ghost"
  | "primaryContained"
  | "secondaryContained"
  | "outlined"
  | "greenContained";

const variantClasses: Record<ButtonVariant, string> = {
  primary:
    "border-emerald-600 bg-emerald-600 text-white shadow-sm shadow-zinc-800/5 enabled:hover:bg-emerald-700",
  primaryContained:
    "border-emerald-600 bg-emerald-600 text-white shadow-sm shadow-zinc-800/5 enabled:hover:bg-emerald-700",
  greenContained:
    "border-emerald-600 bg-emerald-600 text-white shadow-sm shadow-zinc-800/5 enabled:hover:bg-emerald-700",
  secondary:
    "border-emerald-100 bg-white text-zinc-600 enabled:hover:border-emerald-200 enabled:hover:bg-emerald-50",
  secondaryContained:
    "border-emerald-100 bg-white text-zinc-600 enabled:hover:border-emerald-200 enabled:hover:bg-emerald-50",
  outlined:
    "border-emerald-100 bg-white text-zinc-600 enabled:hover:border-emerald-200 enabled:hover:bg-emerald-50",
  ghost:
    "bg-transparent px-2 py-2 text-zinc-500 enabled:hover:bg-zinc-100 enabled:hover:text-emerald-800"
};

interface Props extends ButtonHTMLAttributes<HTMLButtonElement> {
  variant?: ButtonVariant;
  size?: "md" | "lg";
}

const Button = forwardRef<HTMLButtonElement, Props>(
  (
    {
      variant = "secondary",
      size = "md",
      type = "button",
      className,
      ...props
    },
    ref
  ) => (
    <button
      ref={ref}
      type={type}
      className={cn(
        "inline-flex items-center justify-center gap-2 rounded-md border border-transparent px-3.5 py-2.5 text-xs leading-4.5 font-medium whitespace-nowrap",
        variantClasses[variant],
        size === "lg" && "w-full px-4.5 py-3 text-sm",
        className
      )}
      {...props}
    />
  )
);

Button.displayName = "Button";
export default memo(Button);
