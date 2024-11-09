import { HttpService } from "@nestjs/axios";
import { OnWorkerEvent, Processor, WorkerHost } from "@nestjs/bullmq";
import { Inject, Logger } from "@nestjs/common";
import { ConfigService } from "@nestjs/config";
import { Job } from "bullmq";

import { ContentService } from "../content/content.service";
import { ContentEntity } from "../content/entities/content.entity";
import { LLMService } from "../llm/llm.service";
import { RunpodService } from "../runpod/runpod.service";
import { SpeechService } from "../speech/speech.service";
import { STORAGE_SERVICE, StorageService } from "../storage/storage.service";
import { TransformationsService } from "./transformations.service";
import { transformFileToText } from "./transformers/file-to-text.transformer";
import { transformTextToEmbeddings } from "./transformers/text-to-embeddings.transformer";
import { transformTextToImage } from "./transformers/text-to-image.transformer";
import { transformTextToSpeech } from "./transformers/text-to-speech.transformer";
import { transformTextToText } from "./transformers/text-to-text.transformer";

@Processor("run")
export class TransformationProcessor extends WorkerHost {
  private readonly logger: Logger = new Logger("Transformation Processor");

  constructor(
    private transformationsService: TransformationsService,
    @Inject(STORAGE_SERVICE)
    private storageService: StorageService,
    private contentService: ContentService,
    private llmService: LLMService,
    private speechService: SpeechService,
    private httpService: HttpService,
    private configService: ConfigService,
    private runpodService: RunpodService
  ) {
    super();
  }

  @OnWorkerEvent("active")
  async onActive(job: Job) {
    this.logger.log(`Processing job ${job.id} with toolBase ${job.name}`);
    await this.transformationsService.setStatus(
      job.id.toString(),
      "PROCESSING"
    );
  }

  @OnWorkerEvent("completed")
  async onCompleted(job: Job) {
    this.logger.log(`Completed job ${job.id}`);
    await this.transformationsService.setStatus(job.id.toString(), "COMPLETE");
  }

  @OnWorkerEvent("error")
  async onError(job: Job, error: any) {
    this.logger.error(`Error running job ${job.id}: ${error?.message}`);
    try {
      await this.transformationsService.setStatus(job.id.toString(), "ERROR");
      await this.transformationsService.setRunError(
        job.id.toString(),
        error?.message
      );
    } catch {}
  }

  @OnWorkerEvent("failed")
  async onFailed(job: Job, error: any) {
    this.logger.error(`Failed job ${job.id} : ${error?.message}`);
    try {
      await this.transformationsService.setStatus(job.id.toString(), "ERROR");
      await this.transformationsService.setRunError(
        job.id.toString(),
        error?.message
      );
    } catch {}
  }

  async process(job: Job) {
    const runInputContents = job.data.runInputContents as ContentEntity[];
    let runOutputContents: ContentEntity[] = [];
    switch (job.name) {
      case "extract-text":
        runOutputContents = await transformFileToText(
          job.id,
          runInputContents,
          this.logger,
          this.contentService,
          this.storageService,
          this.httpService,
          this.configService
        );
        break;
      case "text-to-image":
        runOutputContents = await transformTextToImage(
          job.id,
          runInputContents,
          this.logger,
          this.contentService,
          this.runpodService,
          this.storageService
        );
        break;
      case "text-to-speech":
        runOutputContents = await transformTextToSpeech(
          job.id,
          runInputContents,
          this.logger,
          this.contentService,
          this.storageService,
          this.speechService
        );
        break;
      case "summarize":
        runOutputContents = await transformTextToText(
          job.id,
          runInputContents,
          this.logger,
          this.contentService,

          this.llmService
        );
        break;
      case "create-embeddings":
        runOutputContents = await transformTextToEmbeddings();
        // job.id,
        // runInputContents,
        // this.logger,
        // this.contentService,
        // this.openAiEmbeddingsService
        break;
      default:
        this.logger.error(`Unknown toolId ${job.name}`);
        throw new Error(`Unknown toolId ${job.name}`);
    }

    this.logger.log(`Adding run output contents to run ${job.id}`);
    await this.transformationsService.setOutputContent(
      job.id.toString(),
      runOutputContents
    );
  }
}
