import { useEffect, useState } from "react";

const useNavigation = () => {
  const [state, setState] = useState<"idle" | "loading">("idle");

  useEffect(() => {
    function onStart() {
      setState("loading");
    }
    function onEnd() {
      setState("idle");
    }

    window.addEventListener("navigation-start", onStart);
    window.addEventListener("navigation-end", onEnd);

    return () => {
      window.removeEventListener("navigation-start", onStart);
      window.removeEventListener("navigation-end", onEnd);
    };
  }, []);

  return { state };
};

const useNProgress = ({ isAnimating }: { isAnimating: boolean }) => {
  const [progress, setProgress] = useState(0);
  const [isFinished, setIsFinished] = useState(!isAnimating);

  useEffect(() => {
    let timer: NodeJS.Timeout;

    if (isAnimating) {
      setIsFinished(false);
      timer = setInterval(() => {
        setProgress((oldProgress) => {
          if (oldProgress >= 0.9) {
            clearInterval(timer);
            return oldProgress;
          }
          const diff = Math.random() * 0.1;
          return Math.min(oldProgress + diff, 0.9);
        });
      }, 200);
    } else {
      setProgress(1);
      timer = setTimeout(() => {
        setIsFinished(true);
        setProgress(0);
      }, 300);
    }

    return () => {
      clearInterval(timer);
      clearTimeout(timer);
    };
  }, [isAnimating]);

  return { isFinished, progress };
};

// EVERYTHING ABOVE IS FAKE

export const PageProgress = () => {
  const navigation = useNavigation();
  const isNavigating = navigation.state === "loading";
  // delay the animation to avoid flickering
  const [isAnimating, setIsAnimating] = useState(false);

  useEffect(() => {
    const timer = setTimeout(() => setIsAnimating(isNavigating), 200);

    return () => clearTimeout(timer);
  }, [isNavigating]);

  const { isFinished, progress } = useNProgress({ isAnimating });

  return (
    <div
      className="absolute right-0 bottom-[-1px] left-0 h-[2px] w-0 bg-primary transition-all duration-300 ease-in-out"
      style={{
        opacity: isFinished ? 0 : 1,
        width: isFinished ? 0 : `${progress * 100}%`,
      }}
    />
  );
};
