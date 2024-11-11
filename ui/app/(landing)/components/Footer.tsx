import { ArchesLogo } from "@/components/arches-logo";

export const Footer = () => {
  return (
    <footer id="footer">
      <hr className="mx-auto w-11/12" />

      <section className="container grid grid-cols-2 gap-x-12 gap-y-8 py-20 md:grid-cols-4 xl:grid-cols-6">
        <div className="col-span-full -mt-8 xl:col-span-2">
          <ArchesLogo scale={0.8} size="sm" />
        </div>

        <div className="flex flex-col gap-2">
          <h3 className="text-lg font-bold">Follow US</h3>
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
          <h3 className="text-lg font-bold">Platforms</h3>
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
          <h3 className="text-lg font-bold">About</h3>
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
          <h3 className="text-lg font-bold">Community</h3>
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
