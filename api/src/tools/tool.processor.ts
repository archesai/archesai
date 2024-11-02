import { OnWorkerEvent, Processor, WorkerHost } from "@nestjs/bullmq";
import { Inject, Logger } from "@nestjs/common";
import { Job } from "bullmq";

import { ContentService } from "../content/content.service";
import { OpenAiEmbeddingsService } from "../embeddings/embeddings.openai.service";
import { JobsService } from "../jobs/jobs.service";
import { LLMService } from "../llm/llm.service";
import { LoaderService } from "../loader/loader.service";
import { OrganizationsService } from "../organizations/organizations.service";
import { RunpodService } from "../runpod/runpod.service";
import { SpeechService } from "../speech/speech.service";
import { STORAGE_SERVICE, StorageService } from "../storage/storage.service";
import { TextChunksService } from "../text-chunks/text-chunks.service";
import { processCreateEmbeddings } from "./processes/create-embeddings.process";
import { processExtractText } from "./processes/extract-text.process";
import { processSummarize } from "./processes/summarize.process";
import { processTextToImage } from "./processes/text-to-image.process";
import { processTextToSpeech } from "./processes/text-to-speech.process";

@Processor("tool")
export class ToolProcessor extends WorkerHost {
  private readonly logger: Logger = new Logger("Tool Processor");

  constructor(
    private runpodService: RunpodService,
    @Inject(STORAGE_SERVICE)
    private storageService: StorageService,
    private organizationsService: OrganizationsService,
    private jobsService: JobsService,
    private contentService: ContentService,
    private loaderService: LoaderService,
    private llmService: LLMService,
    private openAiEmbeddingsService: OpenAiEmbeddingsService,
    private textChunksService: TextChunksService,
    private speechService: SpeechService
  ) {
    super();
  }

  @OnWorkerEvent("active")
  async onActive(job: Job) {
    this.logger.log(
      `Processing job ${job.id} for content ${job.data.content.id} with tool ${job.data.job.toolId}`
    );
    await this.jobsService.updateStatus(job.id.toString(), "PROCESSING");
  }

  @OnWorkerEvent("completed")
  async onCompleted(job: Job) {
    this.logger.log(
      `Completed job ${job.id} for content ${job.data.content.id}`
    );
    await this.jobsService.updateStatus(job.id.toString(), "COMPLETE");
  }

  @OnWorkerEvent("failed")
  async onFailed(job: Job, error: any) {
    this.logger.error(
      `Failed job ${job.id} for content ${job.data?.content.id}: ${error?.message}`
    );
    try {
      await this.jobsService.updateStatus(job.id.toString(), "ERROR");
      await this.jobsService.setJobError(job.id.toString(), error?.message);
    } catch {}
  }

  async process(job: Job) {
    const content = job.data.content;
    switch (job.data.job.toolId) {
      case "extract-text":
        return processExtractText(
          content,
          this.logger,
          this.loaderService,
          this.contentService,
          this.storageService,
          this.textChunksService
        );
      case "text-to-image":
        return processTextToImage(
          content,
          job.data.job,
          this.runpodService,
          this.storageService,
          this.contentService
        );
      case "text-to-speech":
        return processTextToSpeech(
          content,
          this.storageService,
          this.speechService
        );
      case "summarize":
        return processSummarize(
          content,
          this.logger,
          this.loaderService,
          this.contentService,
          this.llmService,
          this.textChunksService
        );
      case "create-embeddings":
        return processCreateEmbeddings(
          content,
          this.logger,
          this.textChunksService,
          this.openAiEmbeddingsService
        );
      default:
        this.logger.error(`Unknown toolId ${job.data.job.toolId}`);
    }
  }
}
