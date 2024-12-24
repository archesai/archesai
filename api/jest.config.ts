/**
 * For a detailed explanation regarding each configuration property, visit:
 * https://jestjs.io/docs/configuration
 */

import type { Config } from 'jest'

const config: Config = {
  displayName: 'api',
  // An array of glob patterns indicating a set of files for which coverage information should be collected
  collectCoverageFrom: ['**/*.(t|j)s', '!**/*.module.ts', '!**/main.ts'],

  // The directory where Jest should output its coverage files
  coverageDirectory: 'coverage',
  // Indicates which provider should be used to instrument code for coverage
  coverageProvider: 'v8',

  // An array of directory names to be searched recursively up from the requiring module's location
  moduleDirectories: ['node_modules', '<rootDir>/..', '<rooDir>'],

  // An array of file extensions your modules use
  moduleFileExtensions: ['js', 'json', 'ts'],

  // A map from regular expressions to module names or to arrays of module names that allow to stub out resources with a single module
  moduleNameMapper: {
    '^@/(.*)$': '<rootDir>/../$1'
  },

  // The root directory that Jest should scan for tests and modules within
  rootDir: 'src',

  // The test environment that will be used for testing
  testEnvironment: 'node',
  // The regexp pattern or array of patterns that Jest uses to detect test files
  testRegex: '.*\\.spec\\.ts$',
  // A map from regular expressions to paths to transformers
  transform: {
    '^.+\\.(t|j)s?$': '@swc/jest'
  }
}

export default config
