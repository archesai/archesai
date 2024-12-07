import type { Config } from 'jest'

export default async (): Promise<Config> => {
  return {
    collectCoverage: true,
    coverageDirectory: './test/coverage-e2e',
    globalSetup: '<rootDir>/test/jest.global-setup.ts',
    moduleDirectories: ['node_modules', '<rootDir>'],
    moduleFileExtensions: ['js', 'json', 'ts'],
    moduleNameMapper: {
      '^@/(.*)$': '<rootDir>/$1'
    },
    preset: 'ts-jest',
    rootDir: '..',
    setupFilesAfterEnv: ['<rootDir>/test/jest.setup.ts'],
    testEnvironment: 'node',
    testRegex: '.e2e-spec.ts$',
    testTimeout: 120000,
    transform: {
      '^.+\\.(t|j)s?$': '@swc/jest'
    }
  }
}
