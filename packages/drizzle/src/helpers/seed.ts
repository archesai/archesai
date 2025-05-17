import { reset, seed } from 'drizzle-seed'

import { createClient } from '#helpers/clients'
import * as schema from '#schema/index'

async function main() {
  const url = process.env.DATABASE_URL
  if (!url) {
    throw new Error('DATABASE_URL is required')
  }
  const db = createClient(url)
  await reset(db, schema)
  await seed(db, schema)
  console.log('Seeding completed successfully')
}

main()
  .then(() => {
    console.log('done')
    process.exit(0)
  })
  .catch((e: unknown) => {
    console.error(e)
    process.exit(1)
  })
