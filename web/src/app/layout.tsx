import type { Metadata } from "next";

import { AppShell } from "@/widgets/app-shell";

import Providers from "./Providers";
import "./globals.css";

export const metadata: Metadata = {
  title: {
    default: "Images · Image Uploader",
    template: "%s · Image Uploader"
  },
  description: "Image library and upload management"
};

const RootLayout = ({ children }: LayoutProps<"/">) => {
  return (
    <html lang="en">
      <body className="bg-zinc-50 font-sans text-sm text-zinc-800 antialiased">
        <Providers>
          <AppShell>{children}</AppShell>
        </Providers>
      </body>
    </html>
  );
};

export default RootLayout;
