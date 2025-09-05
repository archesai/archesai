import { base } from "@archesai/eslint/base"
import { react } from "@archesai/eslint/react"

const config = [
  ...react,
  ...base,
  {
    ignores: ["**/generated/**/*.ts"]
  }
]

export default config
