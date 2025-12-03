import { X } from "lucide-react";
import { useState } from "react";
import { Button } from "#components/shadcn/button";

export const Banner = () => {
  const [isVisible, setIsVisible] = useState(false);
  const [content] = useState<string | null>(null);

  if (!isVisible || !content) {
    return null;
  }

  return (
    <div className="relative flex items-center justify-center bg-primary px-4 py-2 text-primary-foreground">
      <div className="text-center text-sm">{content}</div>
      <Button
        aria-label="Close banner"
        className="absolute right-4 hover:opacity-80"
        onClick={() => setIsVisible(false)}
      >
        <X size={16} />
      </Button>
    </div>
  );
};
