import {
  OnQueueActive,
  OnQueueCompleted,
  OnQueueFailed,
  OnQueueRemoved,
  Process,
  Processor,
} from "@nestjs/bull";
import { Inject, Logger } from "@nestjs/common";
import { Job } from "bull";

import { AudioService } from "../audio/audio.service";
import { retry } from "../common/retry";
import { OpenAiEmbeddingsService } from "../embeddings/embeddings.openai.service";
import { JobsService } from "../jobs/jobs.service";
import { LLMService } from "../llm/llm.service";
import { LoaderService } from "../loader/loader.service";
import { OrganizationsService } from "../organizations/organizations.service";
import { RunpodService } from "../runpod/runpod.service";
import { STORAGE_SERVICE, StorageService } from "../storage/storage.service";
import {
  VECTOR_DB_SERVICE,
  VectorDBService,
} from "../vector-db/vector-db.service";
import { ContentService } from "./content.service";
import { ContentEntity } from "./entities/content.entity";

@Processor("content")
export class ContentProcessor {
  private readonly logger: Logger = new Logger("Content Processor");

  chunkArray = <T>(array: T[], chunkSize: number): T[][] =>
    Array.from({ length: Math.ceil(array.length / chunkSize) }, (v, i) =>
      array.slice(i * chunkSize, i * chunkSize + chunkSize)
    );

  constructor(
    private runpodService: RunpodService,
    private audioService: AudioService,
    private organizationsService: OrganizationsService,
    private jobsService: JobsService,
    private contentService: ContentService,
    private loaderService: LoaderService,
    @Inject(STORAGE_SERVICE)
    private storageService: StorageService,
    private llmService: LLMService,
    private openAiEmbeddingsService: OpenAiEmbeddingsService,
    @Inject(VECTOR_DB_SERVICE)
    private vectorDBService: VectorDBService
  ) {}

  private async createEmbeddings(textContent: { text: string }[]) {
    let embeddings = [] as {
      embedding: number[];
      tokens: number;
    }[];
    const textContentChunks = this.chunkArray(textContent, 100);
    for (const textContentChunk of textContentChunks) {
      const embeddingsChunk = await retry(
        this.logger,
        async () =>
          await this.openAiEmbeddingsService.createEmbeddings(
            textContentChunk.map((x) => x.text)
          ),
        3
      );
      embeddings = embeddings.concat(embeddingsChunk);
    }
    return embeddings;
  }

  private mergeAndFilterEmbeddings(
    embeddings: { embedding: number[] }[],
    textContent: { page: number; text: string; tokens: number }[]
  ) {
    return textContent.map((textContent, index) => {
      return { ...textContent, ...embeddings[index] };
    });
  }

  private async uploadMappings(
    data: {
      embedding: number[];
      page: number;
      text: string;
      tokens: number;
    }[],
    content: ContentEntity
  ) {
    const start = Date.now();
    const dataCopy = JSON.parse(JSON.stringify(data)) as {
      embedding: number[];
      page: number;
      text: string;
      tokens: number;
    }[];

    this.logger.log(`Took ${(Date.now() - start) / 1000}s to copy data`);
    const uploadStart = Date.now();
    while (dataCopy.length) {
      // const currentIndex = data.length - dataCopy.length;
      const docs = dataCopy.splice(0, 1000);
      this.logger.log(`Uploading ${docs.length} contents`, content);
      // for (const doc of docs) {
      //   await this.contentService.create(
      //     content.orgname,
      //     {
      //       buildArgs: {},
      //       name: doc.text,
      //       type: "DOCUMENT",
      //       url: "",
      //     },
      //     {
      //       page: doc.page,
      //       tokens: doc.tokens,
      //       vectorDbId: content.id + "__" + currentIndex.toString(),
      //     }
      //   );
      // }
    }
    this.logger.log(
      `Took ${(Date.now() - uploadStart) / 1000}s to upload data`
    );
    this.logger.log(
      `Took ${(Date.now() - start) / 1000}s to complete whole process`
    );
  }

  @OnQueueActive()
  async onActive(job: Job) {
    const content = job.data.content as ContentEntity;
    await this.jobsService.updateStatus(content.job.id, "PROCESSING");
    // await this.organizationsService.removeCredits(
    //   content.orgname,
    //   (content.maxFrames / 12) * 1000
    // );
    this.logger.log(`Processing job ${job.id} for content ${content.id}...`);
  }

  @OnQueueCompleted()
  async onCompleted(job: Job) {
    const content = job.data.content as ContentEntity;
    await this.jobsService.updateStatus(content.job.id, "COMPLETE");
    this.logger.log(
      `Completed job ${job.id} for content ${job.data.content.id}`
    );
  }

  @OnQueueFailed()
  async onFailed(job: Job, error: any) {
    const content = job.data.content as ContentEntity;
    try {
      await this.jobsService.updateStatus(content.job.id, "ERROR");
      // await this.organizationsService.addCredits(
      //   content.orgname,
      //   (content.maxFrames / 12) * 1000
      // );
    } catch {}

    this.logger.error(
      `Failed job ${job.id} for content ${job.data?.content.id}: ${error?.message}`
    );
  }

  @OnQueueRemoved()
  async onRemoved(job: Job) {
    this.logger.error(
      `Cancelled job ${job.id} for content ${job.data?.content.id}`
    );
  }

  @Process({ concurrency: 8, name: "ANIMATION" })
  async processAnimation(job: Job) {
    const content = job.data.content as ContentEntity;

    // If content uses audio, get the strength
    let strengthSchedule = "0:(0.6)";
    let translationX = "0:(0)";
    let translationZ = "0:(0)";
    const mode = "3D";
    const border = "replicate";
    const zoom = "0: (0.1)";
    if (content.buildArgs.useAudio) {
      this.logger.log("Trimming audio " + content.buildArgs.audioUrl);
      const trimmedUrl = await this.audioService.trimAudio(
        content.orgname,
        content.buildArgs.audioUrl,
        content.buildArgs.audioStart,
        content.buildArgs.length
      );
      this.logger.log("Trimmed audio to " + trimmedUrl);
      content.buildArgs.audioUrl = trimmedUrl;

      const { bassSrc, drumsSrc } =
        await this.audioService.splitAudio(trimmedUrl);

      strengthSchedule = await this.audioService.getKeyframes(
        bassSrc,
        content.buildArgs.fps,
        "0.70 - x^1.5",
        false
      );

      // Determine which to set baesd on the mode
      if (mode === "3D") {
        translationZ = await this.audioService.getKeyframes(
          drumsSrc,
          content.buildArgs.fps,
          "1 + x^4",
          true
        );
      } else {
        translationX = await this.audioService.getKeyframes(
          drumsSrc,
          content.buildArgs.fps,
          "1 + x^4",
          true
        );
      }

      this.logger.log("Strength schedule: " + strengthSchedule);
      this.logger.log("Translation X: " + translationX);
      this.logger.log("Translation Z: " + translationZ);
    }

    const input = {
      input: {
        audio_url: content.buildArgs.audioUrl,
        border: border,
        content_mode: mode,
        content_prompts: content.buildArgs.contentPrompts,
        fps: content.buildArgs.fps,
        height: content.buildArgs.height,
        max_frames: content.buildArgs.maxFrames,
        name: content.name,
        orgname: content.orgname,
        strength_schedule: strengthSchedule,
        translation_x: translationX,
        translation_z: translationZ,
        use_audio: content.buildArgs.useAudio,
        width: content.buildArgs.width,
        zoom: zoom,
      },
    };

    await this.runpodService.runPod(
      content.orgname,
      content.id,
      job.id.toString(),
      "x4mve8pg5zn7bt",
      input
    );
  }
  @Process({ concurrency: 32, name: "DOCUMENT" })
  async processDocument(job: Job) {
    let progress = 0;
    let content = job.data.content as ContentEntity;

    // hit loader endpoint
    const { mimeType, preview, textContent, title, totalTokens } =
      await this.loaderService.extractUrl(
        content.url,
        200,
        content.buildArgs.delimiter
      );
    progress = 0.25;
    await this.jobsService.setProgress(content.job.id, progress);
    this.logger.log(`Extracted text from ${content.name} with ${mimeType}`);

    // update content type
    await this.contentService.updateRaw(content.orgname, content.id, {
      mimeType,
    });

    // update name
    if (title.indexOf("http") == -1) {
      content = new ContentEntity(
        await this.contentService.updateRaw(content.orgname, content.id, {
          name: title,
        })
      );
    }

    // update organization credits
    const organization = await this.organizationsService.findOneByName(
      content.orgname
    );
    const plan = organization.plan;
    if (plan != "PREMIUM" && organization.credits < totalTokens / 50) {
      await this.jobsService.updateStatus(content.job.id, "ERROR");
      await this.contentService.updateRaw(content.orgname, content.id, {
        description: "YOU DO NOT HAVE ENOUGH CREDITS TO PROCESS THIS DOCUMENT",
      });
      return;
    }

    // Create embeddings
    const t1 = Date.now();
    const embeddings = await this.createEmbeddings(textContent);
    progress = 0.5;
    await this.jobsService.setProgress(content.job.id, progress);
    this.logger.log(
      `Created embeddings for ${content.name}.  Completed in ${
        (Date.now() - t1) / 1000
      }s`
    );

    const tokensUsed = embeddings.reduce(
      (prev, embedding) => prev + embedding.tokens,
      0
    );

    // show credits for content
    await this.contentService.incrementCredits(
      content.orgname,
      content.id,
      Math.ceil(tokensUsed / 50)
    );

    // remove credits from organization
    await this.organizationsService.removeCredits(
      content.orgname,
      Math.ceil(tokensUsed / 50)
    );

    // Merge and filter embeddings
    const t2 = Date.now();
    const populatedTextContent = this.mergeAndFilterEmbeddings(
      embeddings,
      textContent
    );
    this.logger.log(
      `Merged and filtered embeddings for ${content.name}. Completed in ${
        (Date.now() - t2) / 1000
      }s`
    );
    progress = 0.6;
    await this.jobsService.setProgress(content.job.id, progress);

    const uploadPreview = (async () => {
      const previewFilename = `${content.name}-preview.png`;
      const decodedImage = Buffer.from(preview, "base64");
      const multerFile = {
        buffer: decodedImage,
        mimetype: "image/png",
        originalname: previewFilename,
        size: decodedImage.length,
      } as Express.Multer.File;
      const url = await this.storageService.upload(
        content.orgname,
        `contents/${content.name}-preview.png`,
        multerFile
      );
      await this.contentService.updateRaw(content.orgname, content.id, {
        previewImage: url,
      });
      progress += 0.1;
      await this.jobsService.setProgress(content.job.id, progress);
    })();

    const upsertVectorsPromise = (async () => {
      const start = Date.now();
      await this.vectorDBService.upsert(
        organization.orgname,
        content.id,
        populatedTextContent.map((e) => e.embedding),
        populatedTextContent.map((e) => e.text)
      );
      this.logger.log(
        `Upserted embeddings for ${content.name}. Completed in ${
          (Date.now() - start) / 1000
        }s`
      );
      progress += 0.1;
      await this.jobsService.setProgress(content.job.id, progress);
    })();

    const upsertMappingsPromise = (async () => {
      const start = Date.now();
      await this.uploadMappings(populatedTextContent, content);
      this.logger.log(
        `Saved mappings for ${content.name}. Completed in ${
          (Date.now() - start) / 1000
        }s`
      );
      progress += 0.1;
      await this.jobsService.setProgress(content.job.id, progress);
    })();

    const summaryPromise = (async () => {
      const start = Date.now();
      const c = this.loaderService.getFirstTokens(
        populatedTextContent.map((x) => x.text),
        3000
      );
      const { summary, tokens } = await retry(
        this.logger,
        async () => await this.llmService.createSummary(c),
        3
      );

      this.logger.log(`Got summary for content for ${content.name}`);

      await this.contentService.updateRaw(content.orgname, content.id, {
        description: summary,
      });

      // show credits for content
      await this.contentService.incrementCredits(
        content.orgname,
        content.id,
        Math.ceil(tokens / 50)
      );

      // remove credits from organization
      await this.organizationsService.removeCredits(
        content.orgname,
        Math.ceil(tokens / 50)
      );

      this.logger.log(
        "Summary saved. Completed in " + (Date.now() - start) / 1000 + "s"
      );
      progress += 0.1;
      await this.jobsService.setProgress(content.job.id, progress);
    })();

    // if any of this fail, throw an error
    await Promise.all([
      upsertVectorsPromise,
      upsertMappingsPromise,
      summaryPromise,
      uploadPreview,
      // clustersPromise,
    ]);
  }

  @Process({ concurrency: 8, name: "IMAGE" })
  async processImage(job: Job) {
    const content = job.data.content as ContentEntity;

    const input = {
      input: {
        height: content.buildArgs.height,
        name: content.name,
        orgname: content.orgname,
        prompt: content.buildArgs.prompt,
        width: content.buildArgs.width,
      },
    };

    await this.runpodService.runPod(
      content.orgname,
      content.id,
      job.id.toString(),
      "7f6wc3v2b1vl1s",
      input
    );
  }
}
