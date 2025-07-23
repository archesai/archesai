import { drizzle } from 'drizzle-orm/node-postgres'
import { reset as resetDb } from 'drizzle-seed'

import * as schema from '#schema/index'

if (!('DATABASE_URL' in process.env))
  throw new Error('DATABASE_URL not found on .env.development')

async function reset() {
  const url = process.env.DATABASE_URL
  if (!url) {
    throw new Error('DATABASE_URL is required')
  }
  console.log('⏳ Resetting database...')
  const start = Date.now()
  const db = drizzle({
    connection: url,
    schema
  })
  await resetDb(db, schema)

  const end = Date.now()
  console.log(`✅ Reset end & took ${(end - start).toString()}ms`)
  console.log('')
  process.exit(0)
}

reset().catch((err: unknown) => {
  console.error('❌ Reset failed')
  console.error(err)
  process.exit(1)
})
