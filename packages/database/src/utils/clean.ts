import { drizzle } from 'drizzle-orm/node-postgres'

import * as schema from '#schema/index'

if (!('DATABASE_URL' in process.env))
  throw new Error('DATABASE_URL not found on .env.development')

const clean = async () => {
  const url = process.env.DATABASE_URL ?? ''
  const db = drizzle({
    connection: url,
    schema
  })

  const tableSchema = db._.schema
  if (!tableSchema) throw new Error('No table schema found')

  console.log('🗑️  Emptying the entire database')

  const queries = Object.values(tableSchema).map((table) => {
    console.log(`🧨 Preparing delete query for table: ${table.dbName}`)
    return table.tsName
  })

  await Promise.all(
    queries.map(async (query) => {
      await db.delete(schema[query])
    })
  )
}
clean().catch((err: unknown) => {
  console.error('❌ Reset failed')
  console.error(err)
  process.exit(1)
})
