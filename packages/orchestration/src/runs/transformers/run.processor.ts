// import { readFileSync } from 'node:fs'
// import type { Job } from 'bullmq'

// import { Worker } from 'bullmq'

// import type { ConfigService } from '@archesai/core'
// import type { ArtifactEntity } from '@archesai/schemas'

// import { Logger } from '@archesai/core'
// import { RUN_ENTITY_KEY } from '@archesai/schemas'

// import type { RunsService } from '#runs/runs.service'

// export type RunJob = Job<ArtifactEntity[], ArtifactEntity[]>

// /**
//  * Processor for runs.
//  */
// export  RunProcessor {
//   private readonly logger =
//   private readonly runsService: RunsService
//   private readonly worker: Worker

//   constructor(configService: ConfigService, runsService: RunsService) {
//     this.runsService = runsService
//     this.worker = new Worker(
//       RUN_ENTITY_KEY,
//       async (job: RunJob) => this.process(job),
//       {
//         autorun: false,
//         connection: {
//           host: configService.get('redis.host'),
//           password: configService.get('redis.auth')!,
//           port: configService.get('redis.port'),
//           ...(configService.get('redis.ca') ?
//             { tls: { ca: readFileSync(configService.get('redis.ca')!) } }
//           : {})
//         }
//       }
//     )

//     this.registerEvents()
//   }

//   public async process(job: RunJob) {
//     if (!job.id) {
//       this.logger.error(`job has no id`, { job })
//       return []
//     }
//     this.logger.log(`processing run`, { job })
//     return Promise.resolve([])
//     // const inputs = job.data
//     // const outputs: ArtifactEntity[] = []
//     // switch (job.name) {
//     // case 'create-embeddings':
//     //   outputs = transformTextToEmbeddings()
//     //   // job.id,
//     //   // inputs,
//     //   // this.logger,
//     //   // this.artifactsService,
//     //   // this.openAiEmbeddingsService
//     //   break
//     // case 'extract-text':
//     //   outputs = await transformFileToText(
//     //     job.id,
//     //     inputs,
//     //     this.logger,
//     //     this.artifactsService,
//     //     this.configService,
//     //     this.storageService
//     //   )
//     //   break
//     // case 'summarize':
//     //   outputs = await transformTextToText(
//     //     job.id,
//     //     inputs,
//     //     this.logger,
//     //     this.artifactsService,
//     //     this.llmService
//     //   )
//     //   break
//     // case 'text-to-image':
//     //   outputs = await transformTextToImage(
//     //     job.id,
//     //     inputs,
//     //     this.logger,
//     //     this.artifactsService,
//     //     this.runpodService,
//     //     this.storageService
//     //   )
//     //   break
//     // case 'text-to-speech':
//     //   outputs = await transformTextToSpeech(
//     //     job.id,
//     //     inputs,
//     //     this.logger,
//     //     this.artifactsService,
//     //     this.storageService,
//     //     this.speechService
//     //   )
//     //   break
//     //   default:
//     //     throw new Error(`Unknown toolId ${job.name}`)
//     // }

//     // this.logger.log(job, `adding run output`)
//     // await this.runsService.updateRelationship(
//     //   job.id.toString(),
//     //   'outputs',
//     //   outputs.map((output) => output.id),
//     //   'add'
//     // )
//   }

//   private onActive(job: RunJob): void {
//     this.logger.log(`processing`, { job })
//     if (!job.id) {
//       this.logger.error(`job has no id`, { job })
//       return
//     }
//     this.runsService
//       .setStatus(job.id.toString(), 'PROCESSING')
//       .catch((error: unknown) => {
//         this.logger.error('Failed to update run status', { error })
//       })
//   }

//   private onCompleted(job: RunJob): void {
//     this.logger.log(`completed`, { job })
//     if (!job.id) {
//       this.logger.error(`job has no id`, { job })
//       return
//     }
//     this.runsService
//       .setStatus(job.id.toString(), 'COMPLETED')
//       .catch((error: unknown) => {
//         this.logger.error('Failed to update run status', { error })
//       })
//   }

//   private onError(error: Error): void {
//     this.logger.error(`error`, {
//       error
//     })
//   }

//   private onFailed(job: RunJob, error: Error): void {
//     if (!job.id) {
//       this.logger.error(`job has no id`, { job })
//       return
//     }
//     this.logger.error(`failed job ${job.id}`, {
//       error,
//       job
//     })
//     try {
//       const jobId = job.id.toString()
//       this.runsService
//         .setStatus(jobId, 'FAILED')
//         .then(() => this.runsService.setRunError(jobId, error.message))
//         .catch(() => {
//           this.logger.error(`Failed to update run status for job ${jobId}`)
//         })
//     } catch {
//       this.logger.error(`Failed to update run status for job ${job.id}`)
//     }
//   }

//   private registerEvents() {
//     this.logger.debug(`registering events`)
//     this.worker.on('active', (job: RunJob) => {
//       this.onActive(job)
//     })
//     this.worker.on('completed', (job: RunJob) => {
//       this.onCompleted(job)
//     })

//     this.worker.on('error', (error) => {
//       this.onError(error)
//     })
//     this.worker.on('failed', (job, error) => {
//       if (job) {
//         this.onFailed(job as RunJob, error)
//       }
//     })
//     this.logger.debug(`events registered`)
//   }
// }
