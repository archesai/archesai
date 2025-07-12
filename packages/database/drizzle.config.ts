import type { Config } from 'drizzle-kit'

import { defineConfig } from 'drizzle-kit'

if (!process.env.DATABASE_URL) {
  throw new Error('Missing DATABASE_URL')
}

const nonPoolingUrl = process.env.DATABASE_URL.replace(':6543', ':5432')

export default defineConfig({
  casing: 'snake_case',
  dbCredentials: { url: nonPoolingUrl },
  dialect: 'postgresql',
  out: './migrations',
  schema: './src/schema/index.ts'
  // verbose: true,
  // strict: true
}) satisfies Config
