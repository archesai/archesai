import type { DatabaseService } from '@archesai/core'
import type {
  PipelineInsertModel,
  PipelineSelectModel
} from '@archesai/database'
import type { PipelineEntity } from '@archesai/schemas'

import { BaseRepository } from '@archesai/core'
import { PipelineTable } from '@archesai/database'
import { PipelineEntitySchema } from '@archesai/schemas'

/**
 * Repository for pipelines.
 */
export class PipelineRepository extends BaseRepository<
  PipelineEntity,
  PipelineInsertModel,
  PipelineSelectModel
> {
  constructor(
    databaseService: DatabaseService<
      PipelineEntity,
      PipelineInsertModel,
      PipelineSelectModel
    >
  ) {
    super(databaseService, PipelineTable, PipelineEntitySchema)
    // this.databaseService = databaseService
  }

  // override async create(data: PipelineInsert): Promise<PipelineEntity> {
  //   const [pipeline] = await this.databaseService.insert(PIPELINE_ENTITY_KEY, [data])
  //   if (!pipeline) {
  //     throw new Error('Failed to create entity')
  //   }

  //   for (const step of data.steps) {
  //     await this.databaseService.insert(PIPELINE_ENTITY_KEY, [
  //       {
  //         name: step.name,
  //         organizationId: data.organizationId,
  //         pipelineId: pipeline.id,
  //         toolId: step.toolId
  //       }
  //     ])

  //     await this.databaseService.db.insert(_PipelineStepDependencies).values(
  //       step.prerequisites.map((prerequisite) => ({
  //         pipelineStepId: step.id,
  //         prerequisiteStepId: prerequisite.pipelineStepId
  //       }))
  //     )
  //   }

  //   return this.findOne(pipeline.id)
  // }

  // override async findOne(id: string): Promise<PipelineEntity> {
  //   const pipeline =
  //     await this.databaseService.db.query.PipelineTable.findFirst({
  //       where: eq(PipelineTable.id, id),
  //       with: {
  //         steps: {
  //           with: {
  //             dependents: {
  //               columns: {
  //                 pipelineStepId: true
  //               }
  //             },
  //             prerequisites: {
  //               columns: {
  //                 pipelineStepId: true
  //               }
  //             },
  //             tool: true
  //           }
  //         }
  //       }
  //     })
  //   if (!pipeline) {
  //     throw new Error('Pipeline not found')
  //   }
  //   return this.toEntity(pipeline)
  // }

  // override async update(
  //   id: string,
  //   data: Partial<CreatePipelineRequest & typeof PipelineTable.$inferInsert> & {
  //     organizationId: string
  //   }
  // ) {
  //   const [pipeline] = await this.databaseService.db
  //     .update(PipelineTable)
  //     .set(data)
  //     .where(eq(PipelineTable.id, id))
  //     .returning()
  //   if (!pipeline) {
  //     throw new Error('Failed to update entity')
  //   }

  //   const steps =
  //     await this.databaseService.db.query.PipelineStepTable.findMany({
  //       where: eq(PipelineStepTable.pipelineId, id)
  //     })
  //   const stepsToDelete = steps.map((tool) => tool.id)
  //   await this.databaseService.db
  //     .delete(PipelineStepTable)
  //     .where(inArray(PipelineStepTable.id, stepsToDelete))

  //   for (const pipelineStep of data.steps || []) {
  //     await this.databaseService.db.insert(PipelineStepTable).values({
  //       name: pipelineStep.name,
  //       organizationId: data.organizationId,
  //       pipelineId: id,
  //       toolId: pipelineStep.toolId
  //     })
  //   }

  //   return this.findOne(pipeline.id)
  // }
}
