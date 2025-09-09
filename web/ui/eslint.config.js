import { base } from "@archesai/eslint/base";
import { react } from "@archesai/eslint/react";

const config = [
  ...react,
  ...base,
  {
    ignores: ["**/shadcn/*"],
  },
];

export default config;

// rollupOptions: {
//   external: [
//     "react",
//     "react-dom",
//     "lucide-react",
//     /@radix-ui/,
//     "@sentry/react",

// declare module "lucide-react/dist/esm/dynamicIconImports.js" {
//   // biome-ignore lint/suspicious/noExplicitAny: Allow any type
//   const icons: Record<string, any>;
//   export default icons;
// }
