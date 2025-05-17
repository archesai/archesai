import path from 'node:path'

import jestOpenAPI from 'jest-openapi'

// Specify the path to your OpenAPI specification
jestOpenAPI.default(path.resolve(__dirname, '#openapi-spec.yaml'))
