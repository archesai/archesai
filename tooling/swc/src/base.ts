// import type { Config } from '@swc/core'

export default {
  jsc: {
    externalHelpers: true,
    keepClassNames: true,
    parser: {
      decorators: true,
      dynamicImport: true,
      syntax: 'typescript'
    },
    target: 'esnext',
    transform: {
      decoratorMetadata: true,
      legacyDecorator: true,
      useDefineForClassFields: true
    }
  },
  minify: false,
  module: {
    strict: true,
    type: 'nodenext'
  },
  sourceMaps: false
}
