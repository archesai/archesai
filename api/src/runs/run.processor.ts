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
import { ToolEntity } from "../tools/entities/tool.entity";
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
    const content = job.data.content as ContentEntity;
    const toolId = job.data.toolId as ToolEntity;
    this.logger.log(
      `Processing job ${job.id} for content ${content.id} with tool ${toolId}`
    );
    await this.runsService.updateStatus(job.id.toString(), "PROCESSING");
  }

  @OnWorkerEvent("completed")
  async onCompleted(job: Job) {
    const content = job.data.content as ContentEntity;
    this.logger.log(`Completed job ${job.id} for content ${content.id}`);
    await this.runsService.updateStatus(job.id.toString(), "COMPLETE");
  }

  @OnWorkerEvent("error")
  async onError(job: Job, error: any) {
    const content = job.data.content as ContentEntity;
    this.logger.error(
      `Failed job ${job.id} for content ${content.id}: ${error?.message}`
    );
    try {
      await this.runsService.updateStatus(job.id.toString(), "ERROR");
      await this.runsService.setRunError(job.id.toString(), error?.message);
    } catch {}
  }

  @OnWorkerEvent("failed")
  async onFailed(job: Job, error: any) {
    const content = job.data.content as ContentEntity;
    this.logger.error(
      `Failed job ${job.id} for content ${content.id}: ${error?.message}`
    );
    try {
      await this.runsService.updateStatus(job.id.toString(), "ERROR");
      await this.runsService.setRunError(job.id.toString(), error?.message);
    } catch {}
  }

  async process(job: Job) {
    const content = job.data.content as ContentEntity;
    const tool = job.data.tool as ToolEntity;
    switch (tool.name) {
      case "extract-text":
        return processExtractText(
          content,
          this.logger,
          this.loaderService,
          this.contentService,
          this.storageService
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
          this.speechService,
          this.contentService
        );
      case "summarize":
        return processSummarize(
          content,
          this.logger,
          this.loaderService,
          this.contentService,
          this.llmService
        );
      case "create-embeddings":
        return processCreateEmbeddings(
          content,
          this.logger,
          this.openAiEmbeddingsService,
          this.contentService
        );
      default:
        this.logger.error(`Unknown toolId ${job.name}`);
    }
  }
}
