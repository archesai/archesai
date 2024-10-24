import { Sheet, SheetContent, SheetTrigger } from "@/components/ui/sheet";
import { useSidebar } from "@/hooks/useSidebar";
import { Menu } from "lucide-react";

import { Breadcrumbs } from "./breadcrumbs";
import { CommandMenu } from "./command-menu";
import { ModeToggle } from "./mode-toggle";
import { Sidebar } from "./sidebar";
import { Button } from "./ui/button";
import { UserButton } from "./user-button";

export const PageHeader = () => {
  const { toggleSidebar } = useSidebar();
  return (
    <header className="flex w-full items-center justify-between bg-background p-3 py-3">
      <Sheet>
        <div className="flex items-center gap-3">
          <SheetTrigger asChild>
            <Button
              className="mr-3 flex h-8 w-8 md:hidden"
              onClick={toggleSidebar}
              size="icon"
              variant="secondary"
            >
              <Menu className="h-5 w-5" />
            </Button>
          </SheetTrigger>
          <Breadcrumbs />
        </div>

        <SheetContent className="w-64 p-0" side="left">
          <Sidebar />
        </SheetContent>

        <div className="flex flex-grow items-center justify-end gap-3">
          <CommandMenu />
          <ModeToggle />
          <UserButton size="sm" />
        </div>
      </Sheet>
    </header>
  );
};
