import { LogoSVG } from "@/components/logo-svg";
import { Button, buttonVariants } from "@/components/ui/button";
import {
  NavigationMenu,
  NavigationMenuItem,
  NavigationMenuList,
} from "@/components/ui/navigation-menu";
import {
  Sheet,
  SheetContent,
  SheetHeader,
  SheetTrigger,
} from "@/components/ui/sheet";
import { GitHubLogoIcon } from "@radix-ui/react-icons";
import { Menu } from "lucide-react";
import Link from "next/link";
import { useTheme } from "next-themes";
import { useEffect, useState } from "react";

interface RouteProps {
  href: string;
  label: string;
}

const routeList: RouteProps[] = [
  {
    href: "#features",
    label: "Features",
  },
  {
    href: "#testimonials",
    label: "Testimonials",
  },
  {
    href: "#pricing",
    label: "Pricing",
  },
  {
    href: "#faq",
    label: "FAQ",
  },
];

export const Navbar = () => {
  const { resolvedTheme } = useTheme();
  const [isOpen, setIsOpen] = useState<boolean>(false);

  const [isTop, setIsTop] = useState(true);

  useEffect(() => {
    const handleScroll = () => {
      setIsTop(window.scrollY === 0);
    };

    window.addEventListener("scroll", handleScroll);

    handleScroll();

    return () => {
      window.removeEventListener("scroll", handleScroll);
    };
  }, []);

  return (
    <header
      className={`sticky top-0 z-40 w-full ${
        isTop
          ? "bg-transparent"
          : "bg-white dark:bg-background border-b shadow-sm"
      }`}
    >
      <NavigationMenu className="mx-auto">
        <NavigationMenuList className="container h-16 px-2.5 w-screen flex justify-between">
          <div className="flex items-center justify-center gap-3">
            <NavigationMenuItem className="font-bold flex">
              <LogoSVG fill={resolvedTheme === "dark" ? "#FFF" : "#000"} />
            </NavigationMenuItem>
            {/* mobile */}
            <span className="flex md:hidden">
              <Sheet onOpenChange={setIsOpen} open={isOpen}>
                <SheetTrigger className="px-2">
                  <Menu
                    className="flex md:hidden h-5 w-5"
                    onClick={() => setIsOpen(true)}
                  ></Menu>
                </SheetTrigger>

                <SheetContent side={"left"}>
                  <SheetHeader>
                    <LogoSVG
                      fill={resolvedTheme === "dark" ? "#FFF" : "#000"}
                    />
                  </SheetHeader>
                  <nav className="flex flex-col justify-center items-center gap-2 mt-4">
                    {routeList.map(({ href, label }: RouteProps) => (
                      <a
                        className={buttonVariants({ variant: "ghost" })}
                        href={href}
                        key={label}
                        onClick={() => setIsOpen(false)}
                        rel="noreferrer noopener"
                      >
                        {label}
                      </a>
                    ))}
                    <a
                      className={`w-[110px] border ${buttonVariants({
                        variant: "secondary",
                      })}`}
                      href="https://github.com/leoMirandaa/shadcn-landing-page.git"
                      rel="noreferrer noopener"
                      target="_blank"
                    >
                      <GitHubLogoIcon className="mr-2 w-5 h-5" />
                      Github
                    </a>
                  </nav>
                </SheetContent>
              </Sheet>
            </span>
            {/* desktop */}
            <nav className="hidden md:flex gap-2">
              {routeList.map((route: RouteProps, i) => (
                <a
                  className={`text-[17px] ${buttonVariants({
                    variant: "ghost",
                  })}`}
                  href={route.href}
                  key={i}
                  rel="noreferrer noopener"
                >
                  {route.label}
                </a>
              ))}
            </nav>
          </div>
          <div className="hidden md:flex gap-2">
            <Button variant={"outline"}>
              <Link href="/auth/login">Log in</Link>
            </Button>
            <Button>
              <Link href="/auth/register">Sign up for free</Link>
            </Button>
          </div>
        </NavigationMenuList>
      </NavigationMenu>
    </header>
  );
};
