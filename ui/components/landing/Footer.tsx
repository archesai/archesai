import { LogoSVG } from "@/components/logo-svg";

export const Footer = () => {
  return (
    <footer id="footer">
      <hr className="w-11/12 mx-auto" />

      <section className="container py-20 grid grid-cols-2 md:grid-cols-4 xl:grid-cols-6 gap-x-12 gap-y-8">
        <div className="col-span-full xl:col-span-2 -mt-8">
          <LogoSVG scale={0.8} size="sm" />
        </div>

        <div className="flex flex-col gap-2">
          <h3 className="font-bold text-lg">Follow US</h3>
          <div>
            <a
              className="opacity-60 hover:opacity-100"
              href="https://github.com/archesai"
              rel="noreferrer noopener"
            >
              Github
            </a>
          </div>

          <div>
            <a
              className="opacity-60 hover:opacity-100"
              href="https://x.com/archesai"
              rel="noreferrer noopener"
            >
              Twitter
            </a>
          </div>
        </div>

        <div className="flex flex-col gap-2">
          <h3 className="font-bold text-lg">Platforms</h3>
          <div>
            <a
              className="opacity-60 hover:opacity-100"
              href="#"
              rel="noreferrer noopener"
            >
              Web
            </a>
          </div>

          <div>
            <a
              className="opacity-60 hover:opacity-100"
              href="#"
              rel="noreferrer noopener"
            >
              Mobile
            </a>
          </div>

          <div>
            <a
              className="opacity-60 hover:opacity-100"
              href="#"
              rel="noreferrer noopener"
            >
              Desktop
            </a>
          </div>
        </div>

        <div className="flex flex-col gap-2">
          <h3 className="font-bold text-lg">About</h3>
          <div>
            <a
              className="opacity-60 hover:opacity-100"
              href="#"
              rel="noreferrer noopener"
            >
              Features
            </a>
          </div>

          <div>
            <a
              className="opacity-60 hover:opacity-100"
              href="#"
              rel="noreferrer noopener"
            >
              Pricing
            </a>
          </div>

          <div>
            <a
              className="opacity-60 hover:opacity-100"
              href="#"
              rel="noreferrer noopener"
            >
              FAQ
            </a>
          </div>
        </div>

        <div className="flex flex-col gap-2">
          <h3 className="font-bold text-lg">Community</h3>
          <div>
            <a
              className="opacity-60 hover:opacity-100"
              href="#"
              rel="noreferrer noopener"
            >
              Youtube
            </a>
          </div>

          <div>
            <a
              className="opacity-60 hover:opacity-100"
              href="#"
              rel="noreferrer noopener"
            >
              Discord
            </a>
          </div>

          <div>
            <a
              className="opacity-60 hover:opacity-100"
              href="#"
              rel="noreferrer noopener"
            >
              Twitch
            </a>
          </div>
        </div>
      </section>

      <section className="container pb-14 text-center">
        <h3>&copy; 2024 Arches Solutions LLC</h3>
      </section>
    </footer>
  );
};
