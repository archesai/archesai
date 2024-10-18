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

import { retry } from "../common/retry";
import { ContentService } from "../content/content.service";
import { ContentEntity } from "../content/entities/content.entity";
import { OpenAiEmbeddingsService } from "../embeddings/embeddings.openai.service";
import { JobsService } from "../jobs/jobs.service";
import { LLMService } from "../llm/llm.service";
import { LoaderService } from "../loader/loader.service";
import { OrganizationsService } from "../organizations/organizations.service";
import { SpeechService } from "../speech/speech.service";
import { STORAGE_SERVICE, StorageService } from "../storage/storage.service";
import { VectorRecordService } from "../vector-records/vector-record.service";

@Processor("document")
export class DocumentProcessor {
  private readonly logger: Logger = new Logger("Document Processor");

  chunkArray = <T>(array: T[], chunkSize: number): T[][] =>
    Array.from({ length: Math.ceil(array.length / chunkSize) }, (v, i) =>
      array.slice(i * chunkSize, i * chunkSize + chunkSize)
    );

  constructor(
    private organizationsService: OrganizationsService,
    private jobsService: JobsService,
    private contentService: ContentService,
    private loaderService: LoaderService,
    @Inject(STORAGE_SERVICE)
    private storageService: StorageService,
    private llmService: LLMService,
    private openAiEmbeddingsService: OpenAiEmbeddingsService,
    private vectorRecordService: VectorRecordService,
    private speechService: SpeechService
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

  @Process({ concurrency: 32, name: "document" })
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

    const uploadPreviewPromise = (async () => {
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

    // const uploadSpeechPromise = (async () => {
    //   const text = populatedTextContent.map((x) => x.text).join("\n");
    //   const audioBuffer = await this.speechService.generateSpeech(text);
    //   const multerFile = {
    //     buffer: audioBuffer,
    //     mimetype: "audio/mpeg",
    //     originalname: `${content.name}.mp3`,
    //     size: audioBuffer.length,
    //   } as Express.Multer.File;
    //   const url = await this.storageService.upload(
    //     content.orgname,
    //     `contents/${content.name}.mp3`,
    //     multerFile
    //   );
    //   // await this.contentService.updateRaw(content.orgname, content.id, {
    //   //   audio: url,
    //   // });
    //   console.log(url);
    //   progress += 0.1;
    //   await this.jobsService.setProgress(content.job.id, progress);
    // })();

    const upsertVectorsPromise = (async () => {
      const start = Date.now();
      await this.vectorRecordService.upsert(
        organization.orgname,
        content.id,
        populatedTextContent
      );
      this.logger.log(
        `Upserted embeddings for ${content.name}. Completed in ${
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
      summaryPromise,
      uploadPreviewPromise,
      // uploadSpeechPromise,
    ]);
  }
}
