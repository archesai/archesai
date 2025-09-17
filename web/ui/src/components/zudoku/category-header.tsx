import { cx } from "class-variance-authority";
import type { ReactNode } from "react";

export const CategoryHeading = ({
  children,
  className,
}: {
  children: ReactNode;
  className?: string;
}) => {
  return (
    <div
      className={cx("mb-2 font-semibold text-primary text-sm", className)}
      data-pagefind-ignore="all"
    >
      {children}
    </div>
  );
};
