import { InjectQueue } from "@nestjs/bull";
import { Injectable, Logger } from "@nestjs/common";
import { Content } from "@prisma/client";
import { Queue } from "bull";

import { BaseService } from "../common/base.service";
import { ContentService } from "../content/content.service";
import { ContentEntity } from "../content/entities/content.entity";
import { WebsocketsService } from "../websockets/websockets.service";
import { CreateDocumentDto } from "./dto/create-document.dto";

@Injectable()
export class DocumentsService
  implements BaseService<Content, CreateDocumentDto, undefined, undefined>
{
  private logger = new Logger(DocumentsService.name);
  constructor(
    private contentService: ContentService,
    private websocketsService: WebsocketsService,
    @InjectQueue("document") private readonly documentQueue: Queue
  ) {}

  async create(
    orgname: string,
    createDocumentDto: CreateDocumentDto
  ): Promise<ContentEntity> {
    const content = await this.contentService.create(orgname, {
      buildArgs: {
        chunkSize: createDocumentDto.chunkSize,
        delimiter: createDocumentDto.delimiter,
      },
      name: createDocumentDto.name,
      type: "DOCUMENT",
      url: createDocumentDto.url,
    });

    const contentEntity = new ContentEntity(content);
    await this.documentQueue.add(
      "document",
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
