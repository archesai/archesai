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
