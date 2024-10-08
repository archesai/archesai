import { InjectQueue } from "@nestjs/bull";
import { Injectable, Logger } from "@nestjs/common";
import { Content } from "@prisma/client";
import { Queue } from "bull";

import { BaseService } from "../common/base.service";
import { ContentService } from "../content/content.service";
import { ContentEntity } from "../content/entities/content.entity";
import { WebsocketsService } from "../websockets/websockets.service";
import { CreateImageDto } from "./dto/create-image.dto";

@Injectable()
export class ImagesService
  implements BaseService<Content, CreateImageDto, undefined, undefined>
{
  private logger = new Logger(ImagesService.name);
  constructor(
    private contentService: ContentService,
    private websocketsService: WebsocketsService,
    @InjectQueue("image") private readonly imageQueue: Queue
  ) {}

  async create(
    orgname: string,
    createImageDto: CreateImageDto
  ): Promise<ContentEntity> {
    const content = await this.contentService.create(orgname, {
      buildArgs: {
        height: createImageDto.height,
        prompt: createImageDto.prompt,
        width: createImageDto.width,
      },
      name: createImageDto.name,
      type: "IMAGE",
      url: "",
    });

    const contentEntity = new ContentEntity(
      await this.contentService.updateRaw(orgname, content.id, {
        mimeType: "image/png",
      })
    );

    await this.imageQueue.add(
      "image",
      {
        content: contentEntity,
      },
      {
        jobId: contentEntity.id,
      }
    );
    this.websocketsService.socket.to(orgname).emit("update");
    return contentEntity;
  }
}
