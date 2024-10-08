import { Button } from "@/components/ui/button";
import { buttonVariants } from "@/components/ui/button";
import { GitHubLogoIcon } from "@radix-ui/react-icons";
import { useRouter } from "next/navigation";

export const Hero = () => {
  const router = useRouter();
  return (
    <section className="z-10 container grid place-items-center py-20 md:py-32 gap-10 h-screen mt-10 max-w-4xl">
      <div className="text-center space-y-6 z-10 ">
        <main className="text-5xl md:text-5xl font-bold">
          Elevate Your Business with AI-Driven Solutions
        </main>

        <p className="text-xl text-muted-foreground md:w-10/12 mx-auto">
          Create intelligent chatbots, generate vibrant AI visuals, and
          integrate seamlessly using our API or no-code widgets.
        </p>

        <div className="space-y-4 md:space-y-0 md:space-x-4">
          <Button
            className="w-full md:w-1/3"
            onClick={() => router.push("/chatbots")}
          >
            Get Started
          </Button>

          <a
            className={`w-full md:w-1/3 ${buttonVariants({
              variant: "outline",
            })}`}
            href="https://github.com/leoMirandaa/shadcn-landing-page.git"
            rel="noreferrer noopener"
            target="_blank"
          >
            Github Repository
            <GitHubLogoIcon className="ml-2 w-5 h-5" />
          </a>
        </div>
      </div>

      {/* Hero cards sections */}
      <div className="z-10">{/* <HeroCards /> */}</div>

      {/* Shadow effect */}
      <div className="shadow"></div>
    </section>
  );
};
