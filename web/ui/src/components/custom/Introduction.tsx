import type { JSX } from "react";

import { cn } from "#lib/utils";

export const Box = ({
  children,
  className,
  ...props
}: React.HTMLAttributes<HTMLDivElement> & {
  children: React.ReactNode;
  className?: string;
}): JSX.Element => {
  return (
    <div
      className={cn("rounded-md border border-[black] bg-white", className)}
      {...props}
    >
      {children}
    </div>
  );
};

export const BoxLongshadow = ({
  children,
  className,
  shadowLength = "medium",
  ...props
}: React.HTMLAttributes<HTMLDivElement> & {
  children?: React.ReactNode;
  className?: string;
  shadowLength?: "large" | "medium";
}): JSX.Element => {
  return (
    <Box
      className={cn(
        "overflow-hidden",
        shadowLength === "medium" && "shadow-[3px_3px_0px_0px_rgba(0,0,0,1)]",
        shadowLength === "large" && "shadow-[8px_8px_0px_0px_rgba(0,0,0,1)]",
        className,
      )}
      {...props}
    >
      {children}
    </Box>
  );
};

const CardTitle = ({ children }: { children: React.ReactNode }) => {
  return <div className="text-2xl font-semibold">{children}</div>;
};

const CardContent = ({ children }: { children: React.ReactNode }) => {
  return <div className="flex flex-col gap-2 p-6">{children}</div>;
};

const CardDescription = ({ children }: { children: React.ReactNode }) => {
  return <div className="text-md text-muted-foreground">{children}</div>;
};

const CardHeader = ({ children }: { children: React.ReactNode }) => {
  return (
    <div className="h-50 relative flex w-full items-end border-b border-black bg-[url(/grid.svg)] bg-center bg-repeat p-8">
      {children}
    </div>
  );
};

export const Introduction = (): JSX.Element => {
  return (
    <div className="grid grid-cols-2 grid-rows-2 gap-10 py-10">
      <BoxLongshadow>
        <CardHeader>
          <img
            alt="Zudoku"
            className="h-16 w-16"
            src="/quickstart.svg"
          />
        </CardHeader>

        <CardContent>
          <CardTitle>Quickstart</CardTitle>
          <CardDescription>
            Learn how to install Zudoku, configure your first project, and
            generate your first docs.
          </CardDescription>
        </CardContent>
      </BoxLongshadow>

      <BoxLongshadow>
        <CardHeader>
          <img
            alt="Zudoku"
            className="h-16 w-16"
            src="/themes.svg"
          />
        </CardHeader>
        <CardContent>
          <CardTitle>Themes</CardTitle>
          <CardDescription>
            Learn how to install Zudoku, configure your first project, and
            generate your first docs.
          </CardDescription>
        </CardContent>
      </BoxLongshadow>

      <BoxLongshadow>
        <CardHeader>
          <img
            alt="Zudoku"
            className="z-20 h-16 w-16"
            src="/components.svg"
          />
        </CardHeader>
        <CardContent>
          <CardTitle>Components</CardTitle>
          <CardDescription>
            Learn how to install Zudoku, configure your first project, and
            generate your first docs.
          </CardDescription>
        </CardContent>
      </BoxLongshadow>
      <BoxLongshadow>
        <CardHeader>
          <img
            alt="Zudoku"
            className="z-20 h-16 w-16"
            src="/authentication.svg"
          />
        </CardHeader>
        <CardContent>
          <CardTitle>Authentication</CardTitle>
          <CardDescription>
            Learn how to install Zudoku, configure your first project, and
            generate your first docs.
          </CardDescription>
        </CardContent>
      </BoxLongshadow>
    </div>
  );
};

export default Introduction;
