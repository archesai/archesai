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
import * as ospath from "path";

import { ContentService } from "../content/content.service";
import { ContentEntity } from "../content/entities/content.entity";
import { JobsService } from "../jobs/jobs.service";
import { RunpodService } from "../runpod/runpod.service";
import { STORAGE_SERVICE, StorageService } from "../storage/storage.service";

@Processor("image")
export class ImageProcessor {
  private readonly logger: Logger = new Logger("Image Processor");

  constructor(
    private runpodService: RunpodService,
    private jobsService: JobsService,
    private contentService: ContentService,
    @Inject(STORAGE_SERVICE)
    private storageService: StorageService
  ) {}

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
    console.log(content);
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

  @Process({ concurrency: 8, name: "image" })
  async processImage(job: Job) {
    const content = job.data.content as ContentEntity;

    const input = {
      input: {
        height: content.buildArgs.height,
        prompt: content.buildArgs.prompt,
        width: content.buildArgs.width,
      },
    };

    const { image_url } = await this.runpodService.runPod(
      content.orgname,
      content.id,
      content.job.id,
      "y55cw5fvbum8q6",
      input
    );

    const base64String = image_url.replace(/^data:image\/\w+;base64,/, "");

    // Convert the remaining base64 string to a buffer
    const buffer = Buffer.from(base64String, "base64");
    const path = `images/${content.id}.png`;

    // Use the upload function
    const url = await this.storageService.upload(content.orgname, path, {
      buffer: buffer,
      originalname: ospath.basename(path),
      size: buffer.length,
    } as Express.Multer.File);

    console.log(url);

    await this.contentService.updateRaw(content.orgname, content.id, {
      previewImage: url,
      url,
    });
  }
}
