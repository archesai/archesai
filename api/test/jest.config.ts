import type { Config } from "jest";

export default async (): Promise<Config> => {
  return {
    globalSetup: "<rootDir>/jest.global-setup.ts",
    moduleDirectories: ["node_modules", "<rootDir>/..", "<rootDir>"],
    moduleFileExtensions: ["js", "json", "ts"],
    moduleNameMapper: {
      "^@/(.*)$": "<rootDir>/../$1",
    },
    preset: "ts-jest",
    rootDir: ".",
    setupFilesAfterEnv: ["<rootDir>/jest.setup.ts"],
    testEnvironment: "node",
    testRegex: ".e2e-spec.ts$",
    testTimeout: 120000,
    transform: {
      "^.+\\.(t|j)s?$": "@swc/jest",
    },
  };
};
