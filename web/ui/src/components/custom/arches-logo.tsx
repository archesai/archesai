/// <reference types="vite/client" />
/// <reference types="vite-plugin-svgr/client" />

import LargeLogoSVG from "@assets/large-logo.svg?react";
import SmallLogoSVG from "@assets/small-logo.svg?react";
import type { JSX } from "react";

export function ArchesLogo({
  size = "lg" as "lg" | "sm",
  scale = 1,
}: {
  scale?: number;
  size?: "lg" | "sm";
}): JSX.Element {
  const content =
    size === "sm" ? (
      <SmallLogoSVG
        className="fill-current text-black transition-all dark:text-white"
        height={scale * 40}
        width={scale * 40}
      />
    ) : (
      <div className="flex items-center gap-1">
        <SmallLogoSVG
          className="fill-current text-primary dark:text-white"
          height={scale * 40}
          width={scale * 40}
        />
        <LargeLogoSVG
          className="fill-current text-primary dark:text-white"
          height={scale * 40}
          width={scale * 80}
        />
      </div>
    );

  return content;
}
