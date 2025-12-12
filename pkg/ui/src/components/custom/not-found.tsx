import { Link } from "@tanstack/react-router";
import type { JSX } from "react";
import { Button } from "../shadcn/button";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "../shadcn/card";

export default function NotFound({
  children,
}: {
  children?: React.ReactNode;
}): JSX.Element {
  return (
    <div className="flex min-h-screen items-center justify-center p-4">
      <Card className="w-full max-w-md text-center">
        <CardHeader>
          <div className="font-bold text-6xl text-muted-foreground">404</div>
          <CardTitle className="text-2xl">Page Not Found</CardTitle>
          <CardDescription>
            The page you&apos;re looking for doesn&apos;t exist or has been
            moved.
            {children}
          </CardDescription>
        </CardHeader>
        <CardContent>
          <div className="flex flex-col gap-2 sm:flex-row">
            <Button
              asChild
              className="flex-1"
            >
              <Link to="/">Go Home</Link>
            </Button>
            <Button
              className="flex-1"
              onClick={() => {
                window.history.back();
              }}
              type="button"
            >
              Go back
            </Button>
          </div>
        </CardContent>
      </Card>
    </div>
  );
}
