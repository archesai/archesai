import { Breadcrumbs } from "./breadcrumbs";
import { CommandMenu } from "./command-menu";
import { ModeToggle } from "./mode-toggle";
export const PageHeader = ({
  description,
  title,
}: {
  description: string;
  title: string;
}) => {
  return (
    <header className="hidden md:block z-[9] md:-mt-0 bg-background">
      <div className="bg-true-white/80 backdrop-blur-md">
        <div className="hstack justify-between items-center w-full p-4">
          <Breadcrumbs />
          <div className="stack gap-0.5 py-3 pt-0 pr-6 max-w-[64rem]">
            <h2
              aria-level={1}
              className="text-xl text-dark font-semibold origin-left line-clamp-1 transform-none"
              data-testid="page-title"
            >
              {title}
            </h2>
            <span
              aria-level={2}
              className="text-sm text-muted-foreground min-h-[1em] opacity-100 transform-none hidden md:block"
            >
              {description}
            </span>
          </div>
          <div className="hstack items-center gap-2">
            <CommandMenu />
            <ModeToggle />
          </div>
        </div>
      </div>
    </header>
  );
};
