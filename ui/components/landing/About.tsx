import { LogoSVG } from "../logo-svg";
import { Statistics } from "./Statistics";

export const About = () => {
  return (
    <section className="container py-24 sm:py-32" id="about">
      <div className="bg-muted/50 border rounded-lg py-12">
        <div className="px-6 flex flex-col-reverse md:flex-row gap-8 md:gap-12 items-center">
          <LogoSVG size="sm" />
          <div className="bg-green-0 flex flex-col justify-between">
            <div className="pb-6">
              <h2 className="text-3xl md:text-4xl font-bold">
                <span className="bg-gradient-to-b from-primary/60 to-primary text-transparent bg-clip-text">
                  About{" "}
                </span>
                Arches AI
              </h2>
              <p className="text-xl text-muted-foreground mt-4">
                At Arches AI, we are revolutionizing the way businesses
                integrate artificial intelligence into their workflows. Our
                platform empowers organizations with advanced AI capabilities,
                enabling them to automate processes, gain insights, and deliver
                better outcomes at scale. From AI-driven chatbots to robust data
                analysis tools, Arches AI provides tailored solutions to meet
                the needs of modern enterprises.
              </p>
            </div>

            <Statistics />
          </div>
        </div>
      </div>
    </section>
  );
};
