// import {
//   ArrowRightIcon,
//   ArrowUpRightIcon,
//   BracesIcon,
//   CopyCheckIcon,
//   FileJson,
//   LayoutGridIcon,
//   LibraryBigIcon,
//   LockIcon,
//   PlayCircleIcon,
//   SquareCheckIcon,
//   TriangleAlertIcon,
//   UnplugIcon,
// } from "zudoku/icons";

// import {
//   AnimatedLabelIcon,
//   AnimatedPackageIcon,
// } from "./components/AnimatedIcons";
// import { BentoBox, BentoDescription, BentoImage } from "./components/Bento";
// import BentoAddOpenAPI from "./components/BentoAddOpenAPI";
// import { BentoAuthReady } from "./components/BentoAuthReady";
// import { BentoInstall } from "./components/BentoInstall";
// import BentoInternalTools from "./components/BentoInternalTools";
// import { BentoStaticSite } from "./components/BentoStaticSite";
// import { Box } from "./components/Box";
// import { BoxLongshadow } from "./components/BoxLongshadow";
// import PoweredByYou from "./components/PoweredByYou";
// import { Preview } from "./components/Preview";
// import { SparklesText } from "./components/Sparkles";
// import { StartCustomizing } from "./components/StartCustomizing";
// import Zudoku from "./components/Zudoku";
// import DiscordIcon from "./DiscordIcon";
// import GithubIcon from "./GithubIcon";

// import "./LandingPage.css";

// import type { JSX } from "react";

// const Link = ({
//   children,
//   href,
//   target,
// }: {
//   children: React.ReactNode;
//   href: string;
//   target?: string;
// }) => {
//   return (
//     <a
//       className="decoration-2 underline-offset-4 hover:underline"
//       href={href}
//       target={target}
//     >
//       {children}
//     </a>
//   );
// };

// const TechStack = [
//   {
//     alt: "Tailwind CSS",
//     height: 35,
//     href: "https://tailwindcss.com/",
//     src: "/tech/tailwind.svg",
//   },
//   {
//     alt: "React",
//     height: 55,
//     href: "https://react.dev/",
//     src: "/tech/react.svg",
//   },
//   {
//     alt: "TypeScript",
//     height: 45,
//     href: "https://typescriptlang.org/",
//     src: "/tech/typescript.svg",
//   },
//   {
//     alt: "Vite",
//     height: 55,
//     href: "https://vite.dev/",
//     src: "/tech/vite.svg",
//   },
//   {
//     alt: "Radix UI",
//     height: 45,
//     href: "https://radix-ui.com/",
//     src: "/tech/radix.svg",
//   },
// ];

// const LandingPage = (): JSX.Element => {
//   return (
//     <div className="z-1 mx-auto flex w-full flex-col items-center gap-25 pt-10 dark:bg-white dark:text-black">
//       <div className="flex w-full max-w-screen-xl flex-col items-center justify-between gap-6 px-10 md:flex-row">
//         <Zudoku />
//         <ul className="flex items-center gap-6">
//           <li>
//             <Link href="/docs">Documentation</Link>
//           </li>
//           <li>
//             <Link href="/docs/components/typography">Components</Link>
//           </li>
//           <li>
//             <Link href="/docs/theme-playground">Themes</Link>
//           </li>
//         </ul>
//         <div className="flex items-center gap-2">
//           <a
//             aria-label="Discord"
//             className="gap-2 rounded-full p-2 transition hover:bg-accent"
//             href="https://discord.zudoku.dev"
//             rel="noreferrer"
//             target="_blank"
//             title="Get help on Discord"
//           >
//             <DiscordIcon className="h-5 w-5" />
//           </a>
//           <a
//             aria-label="GitHub"
//             className="group relative gap-2 rounded-full p-2 transition hover:bg-accent"
//             href="https://github.com/zuplo/zudoku"
//             rel="noreferrer"
//             target="_blank"
//             title="Star us on GitHub"
//           >
//             <SparklesText
//               className="absolute inset-0 bottom-1/2 left-1/2 opacity-0 group-hover:opacity-100"
//               sparklesCount={2}
//             />
//             <GithubIcon className="h-5 w-5" />
//           </a>
//         </div>
//       </div>
//       <Preview />
//       <div className="px-10">
//         <div className="text-center text-3xl font-bold capitalize">
//           Built with a{" "}
//           <span className="bg-gradient-to-br from-[#B6A0FB] via-[#7362EF] to-[#D2C6FF] bg-clip-text text-transparent">
//             modern stack
//           </span>
//         </div>
//         <ul className="mt-14 flex flex-col items-center gap-10 overflow-x-auto md:flex-row md:gap-20">
//           {TechStack.map((tech) => (
//             <li
//               className="shrink-0 scale-95 opacity-100 saturate-0 transition-all ease-in-out hover:scale-100 hover:opacity-100 hover:saturate-100"
//               key={tech.href}
//             >
//               <a
//                 href={tech.href}
//                 rel="noreferrer"
//                 target="_blank"
//               >
//                 <img
//                   alt={tech.alt}
//                   src={tech.src}
//                   style={{ height: tech.height }}
//                 />
//               </a>
//             </li>
//           ))}
//         </ul>
//         <div className="mt-16 flex justify-center">
//           <a
//             className="group group flex w-fit items-center gap-2 self-center rounded-full border border-[#8D83FF] bg-white px-8 py-3 font-medium text-black md:text-xl"
//             href="https://cosmocargo.dev/"
//           >
//             Check Our Live Demo
//             <ArrowRightIcon
//               className="transition-all duration-300 group-hover:translate-x-1"
//               size={20}
//               strokeWidth={2.5}
//             />
//           </a>
//         </div>
//       </div>
//       <div className="w-full">
//         <div className="flex w-full justify-center border-t border-black">
//           <div className="grid w-full max-w-screen-xl grid-cols-1 md:grid-cols-2 lg:grid-cols-3">
//             <div className="border-black p-10 md:border-r lg:col-span-2 xl:border-l">
//               <h2 className="text-3xl font-semibold">Get Started</h2>
//               <p>
//                 Three quick steps will take you from zero to powerful API docs
//                 in minutes.
//               </p>
//             </div>
//             <div className="flex items-end border-black xl:border-r">
//               <div className="mt-auto flex h-full w-full items-end border-t border-black capitalize md:h-1/2 md:justify-end md:border-r-0 md:border-l-0">
//                 <a
//                   className="inline-flex items-center gap-2 p-3 px-10 text-2xl font-semibold decoration-4 underline-offset-4 hover:underline"
//                   href="/docs"
//                 >
//                   Check our Docs{" "}
//                   <ArrowUpRightIcon
//                     size={32}
//                     strokeWidth={1}
//                   />
//                 </a>
//               </div>
//             </div>
//           </div>
//         </div>
//         <div className="mb-10 flex w-full justify-center border-t border-b border-black">
//           <div className="grid w-full max-w-screen-xl grid-cols-1 md:grid-cols-2 lg:grid-cols-3">
//             <div className="border-b border-black md:border-r lg:border-b-0 xl:border-l">
//               <div className="grid grid-rows-[50px_120px_100px] gap-10 p-10">
//                 <img
//                   alt="cli"
//                   className="h-16 w-16"
//                   src="/1.svg"
//                 />
//                 <BentoInstall />
//                 <div className="flex flex-col gap-2">
//                   <h3 className="text-2xl font-semibold">
//                     Install on your CLI
//                   </h3>
//                   <p className="text-muted-foreground">
//                     npm run zudoku install
//                   </p>
//                 </div>
//               </div>
//             </div>
//             <div className="border-b border-black lg:border-r lg:border-b-0">
//               <div className="grid grid-rows-[50px_120px_100px] gap-10 p-10">
//                 <img
//                   alt="cli"
//                   className="h-16 w-16"
//                   src="/2.svg"
//                 />
//                 <BentoAddOpenAPI />
//                 <div>
//                   <h3 className="text-2xl font-semibold">Add your OpenAPI</h3>
//                   <p className="text-muted-foreground">
//                     Bring your OpenAPI schema into the project and add it to the
//                     Zudoku config.
//                   </p>
//                 </div>
//               </div>
//             </div>
//             <StartCustomizing />
//           </div>
//         </div>
//       </div>
//       <div className="flex w-full flex-col items-center px-10">
//         <PoweredByYou />
//         <h3 className="mt-5 mb-20 text-center text-[54px] font-bold capitalize">
//           Packed with powerful
//           <br />
//           features
//         </h3>

//         <div className="grid w-full max-w-screen-lg grid-cols-12 gap-5">
//           <BentoBox className="col-span-full md:col-span-6 lg:col-span-5">
//             <BentoImage>
//               <BoxLongshadow className="flex h-[120%] w-full flex-col">
//                 <div className="relative z-10 mb-2 flex items-center gap-2 border-b border-black p-5 py-4 text-lg font-medium">
//                   <LibraryBigIcon
//                     size={20}
//                     strokeWidth={1.5}
//                   />{" "}
//                   API Catalog
//                   <div className="absolute -bottom-0.5 left-0 h-1 w-37 bg-black" />
//                 </div>
//                 <div className="flex h-full items-center justify-center gap-10">
//                   <div className="flex flex-col gap-2">
//                     <BoxLongshadow className="flex items-center justify-center bg-[#F2F4FF] p-6 text-black">
//                       <AnimatedPackageIcon delay={0.2} />
//                     </BoxLongshadow>{" "}
//                     Tracking API
//                   </div>
//                   <div className="flex flex-col gap-2">
//                     <BoxLongshadow className="flex items-center justify-center bg-[#F2F4FF] p-6 text-black">
//                       <AnimatedLabelIcon delay={0.5} />
//                     </BoxLongshadow>
//                     Label API
//                   </div>
//                 </div>
//               </BoxLongshadow>
//             </BentoImage>
//             <BentoDescription
//               description="Auto-generate docs from OpenAPI v2/v3 schemas—single or multi-API."
//               title="API Catalog"
//             />
//           </BentoBox>
//           <BentoAuthReady />
//           <BentoBox className="col-span-full md:col-span-6 lg:col-span-4">
//             <BentoImage className="group font-mono">
//               <div className="grid grid-cols-[min-content_1fr_min-content] gap-2">
//                 <Box className="col-span-full grid grid-cols-subgrid items-center gap-4 px-4 py-4">
//                   <LockIcon size={18} />
//                   <div className="flex-1">Authentication</div>
//                   <SquareCheckIcon
//                     fill="#F0F1F4"
//                     size={22}
//                     strokeWidth={1.5}
//                   />
//                 </Box>
//                 <Box className="col-span-full grid grid-cols-subgrid items-center gap-4 px-4 py-4">
//                   <BracesIcon size={18} />
//                   <div className="flex-1">Parameters</div>
//                   <SquareCheckIcon
//                     fill="#F0F1F4"
//                     size={22}
//                     strokeWidth={1.5}
//                   />
//                 </Box>
//                 <Box className="col-span-full grid grid-cols-subgrid items-center gap-4 px-4 py-4">
//                   <FileJson size={18} />
//                   <div className="flex-1">Body</div>
//                   <SquareCheckIcon
//                     fill="#F0F1F4"
//                     size={22}
//                     strokeWidth={1.5}
//                   />
//                 </Box>
//                 <Box className="relative col-span-full grid grid-cols-subgrid items-center justify-center gap-4 px-4 py-4">
//                   <div className="font-bold text-[#FF02BD]">GET</div>
//                   <div className="col-span-2 flex-1 truncate text-[#B4B9C9]">
//                     https://myapi.example.com
//                   </div>
//                   <BoxLongshadow className="absolute -right-2.5 -bottom-2.5 flex items-center gap-2 bg-[#F2F4FF] p-2 px-3 font-bold transition-all duration-300 ease-in-out group-hover:-translate-x-1 group-hover:-translate-y-1 group-hover:scale-105 group-hover:rotate-2 group-hover:bg-black group-hover:text-white">
//                     <span className="top-0 left-0 inline-block group-hover:opacity-0">
//                       Send
//                     </span>
//                     <SparklesText className="absolute left-3 inline-block opacity-0 transition-all duration-300 ease-in-out group-hover:opacity-100">
//                       Send
//                     </SparklesText>
//                     <PlayCircleIcon
//                       className="fill-white transition-all duration-300 ease-in-out group-hover:fill-black"
//                       size={24}
//                       strokeWidth={1.5}
//                     />
//                   </BoxLongshadow>
//                 </Box>
//               </div>
//             </BentoImage>
//             <BentoDescription
//               description="Test endpoints live, with support for API keys and auth."
//               title="Interactive Playground"
//             />
//           </BentoBox>
//           <BentoBox className="col-span-full md:col-span-6 lg:col-span-4">
//             <BentoImage className="group flex items-center justify-center">
//               <div className="flex translate-y-3.5 transform flex-col">
//                 <img
//                   alt="Search"
//                   className="scale-90 transition-all duration-300 ease-in-out group-hover:scale-100"
//                   src="/search/search.svg"
//                 />
//                 <div className="relative -top-5 flex h-20 scale-80 items-center justify-end transition-all duration-300 ease-in-out group-hover:translate-x-1 group-hover:translate-y-1 group-hover:scale-100">
//                   <div className="h-1 w-full flex-shrink-0"></div>
//                   <img
//                     alt="CMD Key"
//                     className="flex-2 transition-all duration-300 ease-in-out group-hover:-rotate-10"
//                     src="/search/cmd.svg"
//                   />
//                   <img
//                     alt="K Key"
//                     className="flex-2 transition-all duration-300 ease-in-out group-hover:-translate-x-3 group-hover:rotate-10"
//                     src="/search/k.svg"
//                   />
//                 </div>
//               </div>
//             </BentoImage>
//             <BentoDescription
//               description="Instant, intelligent search powered by Pagefind, Inkeep, etc."
//               title="Built-in Search"
//             />
//           </BentoBox>
//           <BentoStaticSite />
//         </div>
//       </div>
//       <div className="mt-16 flex w-full max-w-screen-lg flex-col items-center gap-16">
//         <div className="text-center text-5xl font-semibold capitalize">
//           Host it{" "}
//           <span className="bg-gradient-to-br from-[#B6A0FB] via-[#7362EF] to-[#D2C6FF] bg-clip-text text-transparent">
//             Anywhere
//           </span>
//         </div>
//         <div className="grid w-full grid-cols-1 items-center justify-center gap-5 md:grid-cols-2 lg:grid-cols-4">
//           <div className="flex items-center justify-center">
//             <img
//               alt="Vercel"
//               src="/host/vercel.svg"
//             />
//           </div>
//           <div className="flex items-center justify-center">
//             <img
//               alt="Cloudflare"
//               src="/host/cloudflare.svg"
//             />
//           </div>
//           <div className="flex items-center justify-center">
//             <img
//               alt="Netlify"
//               src="/host/netlify.svg"
//             />
//           </div>
//           <div className="flex items-center justify-center">
//             <img
//               alt="Zuplo"
//               src="/host/zuplo.svg"
//             />
//           </div>
//         </div>

//         <a
//           className="group flex w-fit items-center gap-2 self-center rounded-full bg-black px-8 py-3 text-xl text-white transition-all duration-300 hover:drop-shadow"
//           href="https://zudoku.dev/docs"
//         >
//           Learn More{" "}
//           <ArrowRightIcon
//             className="transition-all duration-300 group-hover:translate-x-1"
//             size={20}
//           />
//         </a>
//       </div>
//       <div
//         className="w-full rounded-3xl bg-black p-10 px-10 text-white shadow-[0px_-2px_16px_-4px_rgba(0,_0,_0,_0.5)]"
//         style={{
//           animation: "remove-scale 1s cubic-bezier(0, 0.93, 1, 0.61) forwards",
//           animationRange: "entry 0% cover 15%",
//           animationTimeline: "view()",
//         }}
//       >
//         <div className="mx-auto flex max-w-screen-lg flex-col items-center">
//           <div className="my-10 rounded-full border border-[#8D83FF] p-1 px-3 drop-shadow">
//             Designed to Scale
//           </div>
//           <h3 className="mb-30 text-center text-5xl font-semibold">
//             Supercharge your Docs
//             <br />
//             in every scenario
//           </h3>
//           <div className="grid w-full max-w-screen-lg grid-cols-12 gap-5">
//             <BentoBox className="col-span-full lg:col-span-7">
//               <BentoImage className="flex items-center justify-center">
//                 <div className="flex w-full justify-around">
//                   <img
//                     alt="Puzzle"
//                     src="/puzzle/puzzle-1.svg"
//                   />
//                   <img
//                     alt="Puzzle"
//                     src="/puzzle/puzzle-2.svg"
//                   />
//                   <img
//                     alt="Puzzle"
//                     src="/puzzle/puzzle-3.svg"
//                   />
//                   <img
//                     alt="Puzzle"
//                     src="/puzzle/puzzle-4.svg"
//                   />
//                 </div>
//               </BentoImage>
//               <BentoDescription
//                 description="Easy integration with existing plugins (both community and core) and easy extensibility for creating your own."
//                 title="Supercharged Plugins"
//               />
//             </BentoBox>
//             <BentoBox className="col-span-12 md:col-span-6 lg:col-span-5">
//               <BentoImage>
//                 <BoxLongshadow
//                   className="relative h-[120%] w-11/12 bg-[#F2F4FF] p-6"
//                   shadowLength="large"
//                 >
//                   <code className="font-mono leading-loose whitespace-pre-wrap text-[#9095B4]">
//                     {`# Welcome
// **API** docs rule
// ---
// ## Getting Started
// - Edit the markdowne bar
// - See it live
// ---`}
//                   </code>
//                   <div className="absolute top-5 -right-8 flex flex-col gap-2 rounded-lg border border-black bg-white p-4">
//                     <span className="text-2xl font-bold">Welcome</span>
//                     <span>API docs are what we love.</span>
//                     <hr className="h-[2px] bg-neutral-300" />
//                     <span className="text-xl font-bold">Getting Started</span>
//                     <ul className="list-inside list-disc pl-4">
//                       <li>Edit the markdown</li>
//                       <li>See it live</li>
//                     </ul>
//                     <img
//                       alt="Happy"
//                       className="absolute top-4 -right-5"
//                       src="/happy-lsd.svg"
//                     />
//                   </div>
//                 </BoxLongshadow>
//               </BentoImage>
//               <BentoDescription
//                 description="Generate documentation from markdown files, perfect for SEO and performance."
//                 title="MDX Support"
//               />
//             </BentoBox>

//             <BentoInternalTools />
//             <BentoBox className="col-span-full lg:col-span-7">
//               <BentoImage className="flex w-full items-center justify-center">
//                 <div className="grid w-full grid-cols-12 gap-4">
//                   <BoxLongshadow className="col-span-3 flex h-20 items-center justify-center text-2xl font-bold sm:col-span-2">
//                     Aa
//                   </BoxLongshadow>
//                   <BoxLongshadow className="col-span-9 flex h-20 items-center justify-center sm:col-span-10">
//                     <span className="font-mono">
//                       &lt;
//                       <span className="text-[#FF02BD]">
//                         OpenPlaygroundButton
//                       </span>
//                       &nbsp;
//                       <span className="hidden sm:inline">{"{...props}"}</span>
//                       <span className="">&nbsp;/&gt;</span>
//                     </span>
//                   </BoxLongshadow>
//                   <BoxLongshadow className="col-span-8 flex h-20 items-center justify-center p-4">
//                     <div className="flex h-10 w-10 flex-shrink-0 items-center justify-center rounded-full border border-black font-bold">
//                       1
//                     </div>
//                     <div className="h-[1px] w-full bg-black" />
//                     <div className="flex h-10 w-10 flex-shrink-0 items-center justify-center rounded-full border border-black font-bold">
//                       2
//                     </div>
//                     <div className="h-[1px] w-full bg-black" />
//                     <div className="flex h-10 w-10 flex-shrink-0 items-center justify-center rounded-full border border-black font-bold">
//                       3
//                     </div>
//                     <div className="h-[1px] w-full bg-black" />
//                     <div className="flex h-10 w-10 flex-shrink-0 items-center justify-center rounded-full border border-black font-bold">
//                       4
//                     </div>
//                   </BoxLongshadow>
//                   <BoxLongshadow className="col-span-4 flex h-20 items-center justify-around px-4">
//                     <TriangleAlertIcon
//                       size={20}
//                       strokeWidth={1.5}
//                     />
//                     <UnplugIcon
//                       className="hidden md:inline-block"
//                       size={20}
//                       strokeWidth={1.5}
//                     />
//                     <CopyCheckIcon
//                       size={20}
//                       strokeWidth={1.5}
//                     />
//                     <LayoutGridIcon
//                       className="hidden sm:inline-block"
//                       size={20}
//                       strokeWidth={1.5}
//                     />
//                   </BoxLongshadow>
//                 </div>
//               </BentoImage>
//               <BentoDescription
//                 description="Create the developer experience you've always dreamed of with a full suite of reusable components (or create your own)."
//                 title="Ready to use Components"
//               />
//             </BentoBox>
//           </div>
//           <div className="my-10 flex flex-col gap-4 md:flex-row">
//             <a
//               className="group flex w-full items-center justify-center gap-2 self-center rounded-full bg-white px-8 py-3 text-lg font-medium text-black md:w-fit"
//               href="https://zudoku.dev/docs/quickstart"
//             >
//               Explore the Docs
//             </a>
//             <div className="group flex w-fit items-center gap-2 self-center rounded-full border border-white bg-black px-8 py-3 font-mono font-medium text-white md:text-lg">
//               npm create zudoku@latest
//             </div>
//           </div>
//         </div>

//         <div className="mx-auto my-25 grid max-w-screen-lg grid-cols-1 items-end gap-4 md:grid-cols-2">
//           <h3 className="text-5xl font-bold capitalize">
//             Join our open
//             <br />
//             community
//             <br />
//             of developers
//           </h3>
//           <div className="my-10 flex flex-1 flex-col items-start gap-4 md:flex-row md:items-end lg:my-0 lg:justify-end">
//             <a
//               className="w-full rounded-full bg-[#7362EF] px-6 py-2 text-lg font-medium text-nowrap text-white md:w-fit"
//               href="https://discord.zudoku.dev"
//             >
//               Join our Discord
//             </a>
//             <a
//               className="relative w-full rounded-full border border-[#7362EF] px-6 py-2 text-lg font-medium text-nowrap md:w-fit"
//               href="https://github.com/zuplo/zudoku"
//             >
//               <SparklesText sparklesCount={4}>Star on GitHub</SparklesText>
//             </a>
//           </div>
//         </div>
//       </div>
//       <div className="mt-10 mb-30 grid w-full max-w-screen-lg grid-cols-1 px-10 md:grid-cols-2">
//         <div className="flex flex-col gap-10">
//           <Zudoku />
//           <h2 className="text-3xl font-medium capitalize">
//             Zudoku is free, open-source,
//             <br />
//             and ready to power your docs.
//           </h2>
//         </div>
//         <div className="mt-20 flex flex-col gap-1 md:mt-0 md:items-end md:justify-end">
//           <a
//             className="decoration-2 underline-offset-4 hover:underline"
//             href="https://github.com/zuplo/zudoku"
//           >
//             View on GitHub
//           </a>
//           <a
//             className="decoration-2 underline-offset-4 hover:underline"
//             href="https://cosmocargo.dev"
//           >
//             See Live Example
//           </a>
//           <a
//             className="decoration-2 underline-offset-4 hover:underline"
//             href="https://zudoku.dev/docs"
//           >
//             Documentation
//           </a>
//         </div>
//       </div>
//     </div>
//   );
// };

// export default LandingPage;
