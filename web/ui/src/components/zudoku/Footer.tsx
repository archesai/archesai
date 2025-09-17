import { Slot } from "./Slot";

export const Footer = () => {
  return (
    <footer className="border-t bg-background">
      <div className="mx-auto max-w-screen-2xl px-4 py-8 lg:px-8">
        <Slot.Target name="footer-content" />
        <div className="mt-4 text-center text-muted-foreground text-sm">
          <Slot.Target name="footer-bottom" />
        </div>
      </div>
    </footer>
  );
};
