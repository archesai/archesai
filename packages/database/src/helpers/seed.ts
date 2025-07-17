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

  // Step 2: Seed main tables
  console.log('🌱 Seeding main tables...')
  const mainTables = {
    accounts: schema.AccountTable,
    apiToken: schema.ApiTokenTable,
    artifacts: schema.ArtifactTable,
    invitations: schema.InvitationTable,
    labels: schema.LabelTable,
    organizations: schema.OrganizationTable,
    pipelines: schema.PipelineTable,
    pipelineSteps: schema.PipelineStepTable,
    runs: schema.RunTable,
    sessions: schema.SessionTable,
    tools: schema.ToolTable,
    users: schema.UserTable
  }
  await seed(db, mainTables).refine((f) => ({
    artifacts: {
      columns: {
        description: f.loremIpsum(),
        mimeType: f.valuesFromArray({
          values: ['image/png', 'image/jpeg', 'application/pdf', 'text/plain']
        }),
        name: f.jobTitle()
      }
    },
    tools: {
      columns: {
        inputMimeType: f.valuesFromArray({
          values: [
            'text/plain',
            'image/png',
            'application/pdf',
            'audio/mpeg',
            'video/mp4'
          ]
        }),
        outputMimeType: f.valuesFromArray({
          values: [
            'text/plain',
            'image/png',
            'application/pdf',
            'audio/mpeg',
            'video/mp4'
          ]
        })
      },
      count: 5
    },
    users: {
      columns: {
        email: f.email(),
        name: f.firstName()
      },
      count: 10
    }
  }))

  // Step 3: Fetch generated IDs
  console.log('🔍 Fetching data for junction tables...')
  const labels = await db.select().from(schema.LabelTable)
  // const pipelineStep = await db.select().from(schema.PipelineStepTable)
  // const runs = await db.select().from(schema.RunTable)
  const artifacts = await db.select().from(schema.ArtifactTable)

  // Step 4: Seed junction tables
  console.log('🌿 Seeding junction tables...')
  const junctionTables = {
    labelsToArtifacts: schema.LabelToArtifactTable,
    parentToChild: schema.ParentToChildTable,
    // parentToChild: schema.,
    pipelineStepToDependency: schema.PipelineStepToDependency
  }
  const labelToArtifactData = labels.map((label) => ({
    artifactId: artifacts[0]!.id, // Assuming you want to link to the first artifact
    labelId: label.id
  }))
  await db.insert(junctionTables.labelsToArtifacts).values(labelToArtifactData)
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

// await db.insert(schema.AccountTable).values({
//   accountId: '1wvFGVyQ7N6mbGQ2msYDlfRaybZBfAmW',
//   password:
//     '84c8af61c0dc266afff1a2b7e3fe2a18:80ffa9d5156e3cdb9e5c340f0554782e6ad3a554a234af70abd2c2b89675ad9f0e0077d2b2bf948f1fbee09857d75dce104df6ab285f7d9554c25f855c53e700',
//   providerId: 'credential',
//   userId: user.id
// })
