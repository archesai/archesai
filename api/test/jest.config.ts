import type { Config } from "jest";

export default async (): Promise<Config> => {
  return {
    globalTeardown: "./teardown.ts",
    moduleDirectories: ["node_modules", "<rootDir>/..", "<rootDir>"],
    moduleFileExtensions: ["js", "json", "ts"],
    rootDir: ".",
    testEnvironment: "node",
    testRegex: ".e2e-spec.ts$",
    testTimeout: 120000,
    transform: {
      "^.+\\.(t|j)s?$": "@swc/jest",
    },
  };
};
