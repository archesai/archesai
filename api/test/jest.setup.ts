import jestOpenAPI from 'jest-openapi'
import * as path from 'path'

// Specify the path to your OpenAPI specification
jestOpenAPI(path.resolve(__dirname, './openapi-spec.yaml'))
