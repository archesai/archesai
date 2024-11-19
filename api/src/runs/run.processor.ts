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
import { RunsService } from "./runs.service";
import { transformFileToText } from "./transformers/file-to-text.transformer";
import { transformTextToEmbeddings } from "./transformers/text-to-embeddings.transformer";
import { transformTextToImage } from "./transformers/text-to-image.transformer";
import { transformTextToSpeech } from "./transformers/text-to-speech.transformer";
import { transformTextToText } from "./transformers/text-to-text.transformer";

type RunJob = Job<ContentEntity[], ContentEntity[], string>;

@Processor("run")
export class RunProcessor extends WorkerHost {
  private readonly logger: Logger = new Logger("Run Processor");

  constructor(
    private runsService: RunsService,
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
  async onActive(job: RunJob) {
    this.logger.log(`Processing job ${job.id} with toolBase ${job.name}`);
    await this.runsService.setStatus(job.id.toString(), "PROCESSING");
  }

  @OnWorkerEvent("completed")
  async onCompleted(job: RunJob) {
    this.logger.log(`Completed job ${job.id}`);
    await this.runsService.setStatus(job.id.toString(), "COMPLETE");
  }

  @OnWorkerEvent("error")
  async onError(job: RunJob, error: any) {
    this.logger.error(`Error running job ${job.id}: ${error?.message}`);
    try {
      await this.runsService.setStatus(job.id.toString(), "ERROR");
      await this.runsService.setRunError(job.id.toString(), error?.message);
    } catch {}
  }

  @OnWorkerEvent("failed")
  async onFailed(job: RunJob, error: any) {
    this.logger.error(`Failed job ${job.id} : ${error?.message}`);
    try {
      await this.runsService.setStatus(job.id.toString(), "ERROR");
      await this.runsService.setRunError(job.id.toString(), error?.message);
    } catch {}
  }

  async process(job: RunJob) {
    const inputs = job.data as ContentEntity[];
    let outputs: ContentEntity[] = [];
    switch (job.name) {
      case "extract-text":
        outputs = await transformFileToText(
          job.id,
          inputs,
          this.logger,
          this.contentService,
          this.httpService,
          this.configService
        );
        break;
      case "text-to-image":
        outputs = await transformTextToImage(
          job.id,
          inputs,
          this.logger,
          this.contentService,
          this.runpodService,
          this.storageService
        );
        break;
      case "text-to-speech":
        outputs = await transformTextToSpeech(
          job.id,
          inputs,
          this.logger,
          this.contentService,
          this.storageService,
          this.speechService
        );
        break;
      case "summarize":
        outputs = await transformTextToText(
          job.id,
          inputs,
          this.logger,
          this.contentService,
          this.llmService
        );
        break;
      case "create-embeddings":
        outputs = await transformTextToEmbeddings();
        // job.id,
        // inputs,
        // this.logger,
        // this.contentService,
        // this.openAiEmbeddingsService
        break;
      default:
        this.logger.error(`Unknown toolId ${job.name}`);
        throw new Error(`Unknown toolId ${job.name}`);
    }

    this.logger.log(`Adding run output contents to run ${job.id}`);
    await this.runsService.setInputsOrOutputs(
      job.id.toString(),
      "outputs",
      outputs
    );
  }
}
