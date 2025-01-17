import type { Config } from 'jest'
import fs from 'fs'

const swcConfig = JSON.parse(fs.readFileSync(`${__dirname}/../.swcrc`, 'utf-8'))

const config: Config = {
  // collectCoverage: true,
  coverageDirectory: './test/coverage-e2e',
  globalSetup: '<rootDir>/test/jest.global-setup.ts',
  moduleDirectories: ['node_modules', '<rootDir>'],
  moduleFileExtensions: ['ts', 'js', 'json'],
  moduleNameMapper: {
    '^@/(.*)$': '<rootDir>/$1'
  },
  rootDir: '..',
  setupFilesAfterEnv: ['<rootDir>/test/jest.setup.ts'],
  testEnvironment: 'node',
  testRegex: '\\.(e2e|integration)-spec\\.ts$',
  testTimeout: 120000,
  transform: {
    '^.+\\.(t|j)s?$': ['@swc/jest', swcConfig]
  },
  forceExit: true
}

export default config
