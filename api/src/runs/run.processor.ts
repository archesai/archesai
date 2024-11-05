import { OnWorkerEvent, Processor, WorkerHost } from "@nestjs/bullmq";
import { Inject, Logger } from "@nestjs/common";
import { Job } from "bullmq";

import { ContentService } from "../content/content.service";
import { ContentEntity } from "../content/entities/content.entity";
import { OpenAiEmbeddingsService } from "../embeddings/embeddings.openai.service";
import { LLMService } from "../llm/llm.service";
import { LoaderService } from "../loader/loader.service";
import { OrganizationsService } from "../organizations/organizations.service";
import { RunpodService } from "../runpod/runpod.service";
import { SpeechService } from "../speech/speech.service";
import { STORAGE_SERVICE, StorageService } from "../storage/storage.service";
import { processCreateEmbeddings } from "./processes/create-embeddings.process";
import { processExtractText } from "./processes/extract-text.process";
import { processSummarize } from "./processes/summarize.process";
import { processTextToImage } from "./processes/text-to-image.process";
import { processTextToSpeech } from "./processes/text-to-speech.process";
import { RunsService } from "./runs.service";

@Processor("run")
export class RunProcessor extends WorkerHost {
  private readonly logger: Logger = new Logger("Tool Processor");

  constructor(
    private runpodService: RunpodService,
    @Inject(STORAGE_SERVICE)
    private storageService: StorageService,
    private organizationsService: OrganizationsService,
    private runsService: RunsService,
    private contentService: ContentService,
    private loaderService: LoaderService,
    private llmService: LLMService,
    private openAiEmbeddingsService: OpenAiEmbeddingsService,
    private speechService: SpeechService
  ) {
    super();
  }

  @OnWorkerEvent("active")
  async onActive(job: Job) {
    this.logger.log(`Processing job ${job.id} with toolBase ${job.name}`);
    await this.runsService.updateStatus(job.id.toString(), "PROCESSING");
  }

  @OnWorkerEvent("completed")
  async onCompleted(job: Job) {
    this.logger.log(`Completed job ${job.id}`);
    await this.runsService.updateStatus(job.id.toString(), "COMPLETE");
  }

  @OnWorkerEvent("error")
  async onError(job: Job, error: any) {
    this.logger.error(`Error running job ${job.id}: ${error?.message}`);
    try {
      await this.runsService.updateStatus(job.id.toString(), "ERROR");
      await this.runsService.setRunError(job.id.toString(), error?.message);
    } catch {}
  }

  @OnWorkerEvent("failed")
  async onFailed(job: Job, error: any) {
    this.logger.error(`Failed job ${job.id} : ${error?.message}`);
    try {
      await this.runsService.updateStatus(job.id.toString(), "ERROR");
      await this.runsService.setRunError(job.id.toString(), error?.message);
    } catch {}
  }

  async process(job: Job) {
    const runInputContents = job.data.runInputContents as ContentEntity[];
    switch (job.name) {
      case "extract-text":
        return processExtractText(
          job.id,
          runInputContents,
          this.logger,
          this.contentService,
          this.loaderService,
          this.storageService
        );
      case "text-to-image":
        return processTextToImage(
          job.id,
          runInputContents,
          this.logger,
          this.contentService,
          this.runpodService,
          this.storageService
        );
      case "text-to-speech":
        return processTextToSpeech(
          job.id,
          runInputContents,
          this.logger,
          this.contentService,
          this.storageService,
          this.speechService
        );
      case "summarize":
        return processSummarize(
          job.id,
          runInputContents,
          this.logger,
          this.contentService,
          this.loaderService,
          this.llmService
        );
      case "create-embeddings":
        return processCreateEmbeddings(
          job.id,
          runInputContents,
          this.logger,
          this.contentService,
          this.openAiEmbeddingsService
        );
      default:
        this.logger.error(`Unknown toolId ${job.name}`);
    }
  }
}
