import { useSidebar } from "@/components/ui/sidebar";
import { Menu } from "lucide-react";

import { CommandMenu } from "./command-menu";
import { ModeToggle } from "./mode-toggle";
import { Button } from "./ui/button";
import { UserButton } from "./user-button";

export const PageHeader = () => {
  const { toggleSidebar } = useSidebar();
  return (
    <header className="flex w-full items-center justify-between bg-background p-3 py-3">
      <div className="flex items-center gap-3">
        <Button
          className="mr-3 flex h-8 w-8"
          onClick={toggleSidebar}
          size="icon"
          variant="secondary"
        >
          <Menu className="h-5 w-5" />
        </Button>
      </div>

      <div className="flex flex-grow items-center justify-end gap-3">
        <CommandMenu />
        <ModeToggle />
        <UserButton size="sm" />
      </div>
    </header>
  );
};
