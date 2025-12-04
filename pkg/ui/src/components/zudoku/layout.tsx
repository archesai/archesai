import type { ReactNode } from "react";
import { Suspense, useEffect } from "react";
import { useZudoku } from "./context/ZudokuContext";
import { Footer } from "./Footer";
import { Header } from "./header";
import { Main } from "./main";
import { Slot } from "./Slot";
import { Spinner } from "./spinner";
import { cn } from "./utils";

const useScrollToAnchor = () => {
  useEffect(() => {
    const hash = window.location.hash;
    if (hash) {
      const id = hash.substring(1);
      const element = document.getElementById(id);
      if (element) {
        element.scrollIntoView({ behavior: "smooth" });
      }
    }
  }, []);
};

const useScrollToTop = () => {
  useEffect(() => {
    window.scrollTo(0, 0);
  }, []);
};

const LoadingFallback = () => (
  <main className="col-span-full row-span-full grid place-items-center">
    <Spinner />
  </main>
);

export const Layout = ({ children }: { children?: ReactNode }) => {
  const context = useZudoku();
  const authentication =
    "authentication" in context ? context.authentication : undefined;

  useScrollToAnchor();
  useScrollToTop();

  useEffect(() => {
    // Initialize the authentication plugin
    if (
      authentication &&
      typeof authentication === "object" &&
      "onPageLoad" in authentication
    ) {
      const onPageLoad = authentication.onPageLoad;
      if (typeof onPageLoad === "function") {
        onPageLoad();
      }
    }
  }, [authentication]);

  return (
    <>
      <Slot.Target name="layout-before-head" />
      <Header />
      <Slot.Target name="layout-after-head" />

      <div
        className={cn(
          "grid w-full max-w-screen-2xl lg:mx-auto",
          "grid-rows-[0_min-content_1fr] lg:grid-rows-[min-content_1fr] [&:has(>:only-child)]:grid-rows-1",
          "grid-cols-1 lg:grid-cols-[var(--side-nav-width)_1fr]",
        )}
      >
        <Suspense fallback={<LoadingFallback />}>
          <Main>{children}</Main>
        </Suspense>
      </div>
      <Footer />
    </>
  );
};
