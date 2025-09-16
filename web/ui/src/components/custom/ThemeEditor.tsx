import { useTheme } from "next-themes";
import type React from "react";
import { useEffect, useMemo, useState } from "react";

import { Callout } from "#components/custom/Callout";
import {
  ClipboardPasteIcon,
  DownloadIcon,
  MoonIcon,
  RotateCcwIcon,
  SunIcon,
} from "#components/custom/icons";
import { Link } from "#components/primitives/link";
import { Alert, AlertDescription, AlertTitle } from "#components/shadcn/alert";
import { Button } from "#components/shadcn/button";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
} from "#components/shadcn/card";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "#components/shadcn/dialog";
import { Progress } from "#components/shadcn/progress";
import { Switch } from "#components/shadcn/switch";
import { Textarea } from "#components/shadcn/textarea";
import { baseColors } from "#lib/base-colors";
import { cn } from "#lib/utils";

const availableRadius = [0, 0.3, 0.6, 1];

const camelToKebabCase = (str: string): string =>
  str.replace(/([A-Z])/g, "-$1").toLowerCase();

type ResolvedTheme = "dark" | "light";

export const ThemeEditor = (): React.ReactElement => {
  const { resolvedTheme, setTheme } = useTheme() as {
    resolvedTheme: ResolvedTheme | undefined;
    setTheme: (theme: string) => void;
  };
  const [color, setColor] = useState<string>();
  const [radius, setRadius] = useState<number>();
  const [customCss, setCustomCss] = useState("");
  const [isPasteDialogOpen, setIsPasteDialogOpen] = useState(false);

  const activeColor = useMemo(() => {
    return baseColors.find((c) => c.name === color);
  }, [color]);

  useEffect(() => {
    if (activeColor && resolvedTheme && resolvedTheme in activeColor.cssVars) {
      const cssVars =
        activeColor.cssVars[resolvedTheme as keyof typeof activeColor.cssVars];
      Object.entries(cssVars).forEach(([key, value]) => {
        document.documentElement.style.setProperty(
          `--${camelToKebabCase(key)}`,
          value as string,
        );
      });
    }

    if (typeof radius === "number") {
      document.documentElement.style.setProperty(
        "--radius",
        `${String(radius)}rem`,
      );
    } else {
      document.documentElement.style.removeProperty("--radius");
    }

    return () => {
      document.documentElement.style.removeProperty("--radius");

      if (
        !activeColor?.cssVars ||
        !resolvedTheme ||
        !(resolvedTheme in activeColor.cssVars)
      )
        return;

      const cssVars =
        activeColor.cssVars[resolvedTheme as keyof typeof activeColor.cssVars];
      Object.entries(cssVars).forEach(([key]) => {
        document.documentElement.style.removeProperty(
          `--${camelToKebabCase(key)}`,
        );
      });
    };
  }, [activeColor, resolvedTheme, radius]);

  const handleReset = () => {
    setColor(undefined);
    setRadius(undefined);
    setCustomCss("");
  };

  const handlePasteTheme = (pastedCss: string) => {
    setCustomCss(pastedCss);
  };

  return (
    <div className="not-prose">
      <style>{customCss}</style>
      <div className="mt-4 flex gap-2">
        <Button
          onClick={handleReset}
          size="sm"
          variant="outline"
        >
          <RotateCcwIcon
            className="me-2"
            size={16}
          />{" "}
          Reset Theme
        </Button>
        <Dialog>
          <DialogTrigger asChild>
            <Button
              size="sm"
              variant="outline"
            >
              <DownloadIcon
                className="me-2"
                size={16}
              />{" "}
              Get Theme Config
            </Button>
          </DialogTrigger>
          <DialogContent className="max-w-[666px]">
            <DialogHeader>
              <DialogTitle>Theme </DialogTitle>
              <DialogDescription>
                Copy and paste the following code into your Zudoku config.
              </DialogDescription>
            </DialogHeader>
            {/* <SyntaxHighlight
              className="max-h-[350px]"
              code={JSON.stringify(themeConfig, null, 2)}
              language="css"
              showLanguageIndicator
            /> */}
          </DialogContent>
        </Dialog>
        <Dialog
          onOpenChange={setIsPasteDialogOpen}
          open={isPasteDialogOpen}
        >
          <DialogTrigger asChild>
            <Button
              size="sm"
              variant="outline"
            >
              <ClipboardPasteIcon
                className="me-2"
                size={16}
              />
              Paste theme
            </Button>
          </DialogTrigger>
          <DialogContent className="max-w-[666px]">
            <DialogHeader>
              <DialogTitle>Paste Custom CSS</DialogTitle>
              <DialogDescription>
                Paste CSS from theme editors like{" "}
                <a
                  className="text-primary underline"
                  href="https://tweakcn.com/"
                  rel="noreferrer"
                  target="_blank"
                >
                  tweakcn.com
                </a>{" "}
                or other shadcn theme generators.
              </DialogDescription>
            </DialogHeader>
            <div className="space-y-4">
              <Textarea
                className="min-h-[200px] font-mono text-sm"
                defaultValue={customCss}
                onChange={(e) => {
                  const css = e.target.value;
                  if (css.trim()) {
                    handlePasteTheme(css);
                  }
                }}
                placeholder="Paste your CSS here..."
              />
              <div className="flex gap-2">
                <Button
                  onClick={() => {
                    setCustomCss("");
                    setIsPasteDialogOpen(false);
                  }}
                  variant="outline"
                >
                  Clear
                </Button>
                <Button
                  onClick={() => {
                    setIsPasteDialogOpen(false);
                  }}
                >
                  Close
                </Button>
              </div>
            </div>
          </DialogContent>
        </Dialog>
        <Button
          asChild
          size="sm"
          variant="link"
        >
          <Link href="/docs/customization/colors-theme">Documentation</Link>
        </Button>
      </div>
      <div className="border-px border-border my-2 border-b border-dashed" />
      <div className="grid grid-cols-[minmax(0,560px)_1fr] gap-2">
        <div className="flex flex-col gap-2">
          <Card>
            {/* <CardHeader className="py-4" /> */}
            <CardContent className="grid grid-cols-1 lg:grid-cols-2">
              <CardHeader className="px-0 py-6">
                <CardDescription>Mode</CardDescription>
              </CardHeader>
              <CardHeader className="px-0 py-6">
                <CardDescription>Radius</CardDescription>
              </CardHeader>
              <div className="flex gap-2">
                <Button
                  className={cn(
                    resolvedTheme === "light" && "border-primary border-2",
                  )}
                  onClick={() => {
                    setTheme("light");
                  }}
                  variant="outline"
                >
                  <SunIcon
                    className="me-2"
                    size={16}
                  />
                  Light
                </Button>
                <Button
                  className={cn(
                    resolvedTheme === "dark" && "border-primary border-2",
                  )}
                  onClick={() => {
                    setTheme("dark");
                  }}
                  variant="outline"
                >
                  <MoonIcon
                    className="me-2"
                    size={16}
                  />
                  Dark
                </Button>
              </div>
              <div className="flex gap-2">
                {availableRadius.map((r) => (
                  <Button
                    className={cn(
                      r === radius && "border-primary border-2",
                      "w-10",
                    )}
                    key={r}
                    onClick={() => {
                      setRadius(r);
                    }}
                    size="sm"
                    variant="outline"
                  >
                    {r}
                  </Button>
                ))}
              </div>
            </CardContent>

            <CardHeader className="flex flex-row items-center gap-2 space-y-0 py-4">
              <CardDescription>Colors</CardDescription>
              <a
                className="text-xs"
                href="https://ui.shadcn.com/themes"
                rel="noreferrer"
                target="_blank"
              >
                by shadcn
              </a>
            </CardHeader>
            <CardContent>
              <div className="grid grid-cols-2 gap-2 md:grid-cols-3 lg:grid-cols-4">
                {baseColors.map((color) => (
                  <Button
                    className={cn(
                      color.name === activeColor?.name &&
                        "border-primary border-2",
                    )}
                    key={color.name}
                    onClick={() => {
                      setColor(color.name);
                    }}
                    size="sm"
                    style={
                      {
                        "--theme-primary":
                          resolvedTheme &&
                          resolvedTheme in (activeColor?.activeColor ?? {})
                            ? activeColor?.activeColor[
                                resolvedTheme as keyof typeof activeColor.activeColor
                              ]
                            : undefined,
                      } as React.CSSProperties
                    }
                    variant="outline"
                  >
                    <div
                      className="me-2 h-4 w-4 rounded-full"
                      style={{
                        backgroundColor:
                          resolvedTheme && resolvedTheme in color.activeColor
                            ? color.activeColor[
                                resolvedTheme as keyof typeof color.activeColor
                              ]
                            : undefined,
                      }}
                    />

                    <div className="flex-1">{color.name}</div>
                  </Button>
                ))}
              </div>
            </CardContent>
          </Card>
          {/* <SyntaxHighlight
            code={`
import { Button } from "zudoku/ui/Button.js";

export const App = () => {
  const [count, setCount] = useState(0);

  return (
    <div>
      <Button onClick={() => setCount(count + 1)}>
        Click me
      </Button>
      <div>Count: {count}</div>
    </div>
  );
};
          `.trim()}
            language="tsx"
            showLanguageIndicator
            showLineNumbers
          /> */}
        </div>
        <div className="overflow-hidden rounded-lg">
          <div className="grid grid-cols-1 gap-2">
            <Card>
              <CardHeader>
                <CardDescription>Button Preview</CardDescription>
              </CardHeader>
              <CardContent className="grid grid-cols-3 gap-2">
                <Button>Button</Button>
                <Button variant="outline">Outline</Button>
                <Button variant="ghost">Ghost</Button>
                <Button variant="link">Link</Button>
                <Button variant="secondary">Secondary</Button>
                <Button variant="destructive">Destructive</Button>
              </CardContent>
            </Card>
            <Card>
              <CardHeader>
                <CardDescription>Controls </CardDescription>
              </CardHeader>
              <CardContent className="grid grid-cols-2 gap-2 text-sm font-medium">
                <div>On</div>
                <Switch defaultChecked={true} />
                <div>Off</div>
                <Switch />
                <div>50%</div>
                <Progress value={50} />
              </CardContent>
            </Card>
            <Card>
              <CardHeader>
                <CardDescription>Alerts</CardDescription>
              </CardHeader>
              <CardContent>
                <Alert>
                  <AlertTitle>Alert</AlertTitle>
                  <AlertDescription>
                    This is an alert. It is used to display important
                    information.
                  </AlertDescription>
                </Alert>
                <Callout type="info">
                  This is a callout. It is used to display important
                  information.
                </Callout>
                <Callout type="caution">
                  This is a callout. It is used to display important
                  information.
                </Callout>
                <Callout type="danger">
                  This is a callout. It is used to display important
                  information.
                </Callout>
                <Callout type="tip">
                  This is a callout. It is used to display important
                  information.
                </Callout>
                <Callout type="note">
                  This is a callout. It is used to display important
                  information.
                </Callout>
              </CardContent>
            </Card>
          </div>
        </div>
      </div>
    </div>
  );
};
export const ThemeEditorPage = (): React.ReactElement => {
  return (
    <div className="flex flex-col gap-3 pt-6">
      <div className="text-4xl font-extrabold">Color in Your App.</div>
      <div>Hand-picked themes that you can copy and paste into your apps.</div>

      <ThemeEditor />
    </div>
  );
};

export default ThemeEditorPage;
