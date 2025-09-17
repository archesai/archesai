/// <reference types="vite/client" />
/// <reference types="vite-plugin-svgr/client" />

import type { JSX } from "react";
import LargeLogoSVG from "../../../../../assets/large-logo.svg?react";
import SmallLogoSVG from "../../../../../assets/small-logo.svg?react";

export function ArchesLogo({
  size = "lg" as "lg" | "sm",
  scale = 1,
}: {
  scale?: number;
  size?: "lg" | "sm";
}): JSX.Element {
  return (
    <div className="flex items-center justify-center gap-2">
      {size === "sm" ? (
        <SmallLogoSVG
          className="fill-current text-primary dark:text-foreground"
          height={scale * 40}
          width={scale * 40}
        />
      ) : (
        <>
          <SmallLogoSVG
            className="text-primary dark:text-foreground"
            height={scale * 40}
            width={scale * 40}
          />
          <LargeLogoSVG
            className="text-primary dark:text-foreground"
            height={scale * 60}
            width={scale * 80}
          />
        </>
      )}
    </div>
  );
}
