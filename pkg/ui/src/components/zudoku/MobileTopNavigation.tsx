import { Menu, X } from "lucide-react";
import { useState } from "react";
import { Button } from "#components/shadcn/button";
import { Sheet, SheetContent, SheetTrigger } from "#components/shadcn/sheet";

export const MobileTopNavigation = () => {
  const [isOpen, setIsOpen] = useState(false);

  return (
    <div className="lg:hidden">
      <Sheet
        onOpenChange={setIsOpen}
        open={isOpen}
      >
        <SheetTrigger asChild>
          <Button
            size="icon"
            variant="ghost"
          >
            <Menu size={20} />
            <span className="sr-only">Toggle menu</span>
          </Button>
        </SheetTrigger>
        <SheetContent
          className="w-[250px]"
          side="right"
        >
          <nav className="flex flex-col gap-4">
            <div className="flex items-center justify-between">
              <span className="font-semibold text-lg">Menu</span>
              <Button
                onClick={() => setIsOpen(false)}
                size="icon"
                variant="ghost"
              >
                <X size={20} />
                <span className="sr-only">Close menu</span>
              </Button>
            </div>
          </nav>
        </SheetContent>
      </Sheet>
    </div>
  );
};
