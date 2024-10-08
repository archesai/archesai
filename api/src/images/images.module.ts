import { BullModule } from "@nestjs/bull";
import { Module } from "@nestjs/common";

import { ContentModule } from "../content/content.module";
import { JobsModule } from "../jobs/jobs.module";
import { RunpodModule } from "../runpod/runpod.module";
import { StorageModule } from "../storage/storage.module";
import { WebsocketsModule } from "../websockets/websockets.module";
import { ImageProcessor } from "./image.processor";
import { ImagesController } from "./images.controller";
import { ImagesService } from "./images.service";

@Module({
  controllers: [ImagesController],
  imports: [
    BullModule.registerQueue({
      defaultJobOptions: {
        attempts: 1,
      },
      name: "image",
    }),
    JobsModule,
    RunpodModule,
    ContentModule,
    WebsocketsModule,
    StorageModule.forRoot(),
  ],
  providers: [ImageProcessor, ImagesService],
})
export class ImagesModule {}
