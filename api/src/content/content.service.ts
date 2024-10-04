import { InjectQueue } from "@nestjs/bull";
import { Inject, Injectable, Logger } from "@nestjs/common";
import { Content, Job, Prisma } from "@prisma/client";
import { Queue } from "bull";

import { BaseService } from "../common/base.service";
import { STORAGE_SERVICE, StorageService } from "../storage/storage.service";
import { WebsocketsService } from "../websockets/websockets.service";
import { ContentRepository } from "./content.repository";
import { ContentQueryDto } from "./dto/content-query.dto";
import { CreateContentDto } from "./dto/create-content.dto";
import { UpdateContentDto } from "./dto/update-content.dto";
import { ContentEntity } from "./entities/content.entity";

@Injectable()
export class ContentService
  implements
    BaseService<Content, CreateContentDto, ContentQueryDto, UpdateContentDto>
{
  private logger = new Logger(ContentService.name);
  constructor(
    @Inject(STORAGE_SERVICE)
    private storageService: StorageService,
    private contentRepository: ContentRepository,
    private websocketsService: WebsocketsService,
    @InjectQueue("content") private readonly contentQueue: Queue
  ) {}

  async create(
    orgname: string,
    createContentDto: CreateContentDto,
    annotations?: object
  ): Promise<ContentEntity> {
    // await this.organizationsService.checkCredits(orgname, estimatedCredits); FIXME

    const content = await this.contentRepository.create(
      orgname,
      createContentDto,
      annotations
    );
    const contentEntity = new ContentEntity(content);

    await this.contentQueue.add(
      createContentDto.type,
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

  async findAll(orgname: string, contentQueryDto: ContentQueryDto) {
    return this.contentRepository.findAll(orgname, contentQueryDto);
  }

  async findOne(id: string) {
    const content = await this.contentRepository.findOne(id);
    const populated = await this.populateReadUrl(content);
    return populated;
  }

  async incrementCredits(orgname: string, id: string, credits: number) {
    const content = await this.contentRepository.incrementCredits(id, credits);
    this.websocketsService.socket.to(orgname).emit("update");
    return content;
  }

  async populateReadUrl(
    content: { job: Job } & Content
  ): Promise<{ job: Job } & Content> {
    if (
      content.url.startsWith(
        `https://storage.googleapis.com/archesai/storage/${content.orgname}/`
      )
    ) {
      const path = content.url
        .replace(
          `https://storage.googleapis.com/archesai/storage/${content.orgname}/`,
          ""
        )
        .split("?")[0];

      try {
        const read = await this.storageService.getSignedUrl(
          content.orgname,
          decodeURIComponent(path),
          "read"
        );
        content.url = read;
      } catch (e) {
        this.logger.warn(e);
        content.url = "";
      }
    }

    return content;
  }

  async remove(orgname: string, contentId: string): Promise<void> {
    await this.contentRepository.remove(orgname, contentId);
    this.websocketsService.socket.to(orgname).emit("update");
  }

  async update(
    orgname: string,
    id: string,
    updateContentDto: UpdateContentDto
  ) {
    const content = await this.contentRepository.update(
      orgname,
      id,
      updateContentDto
    );
    this.websocketsService.socket.to(orgname).emit("update");
    return content;
  }

  async updateRaw(orgname: string, id: string, raw: Prisma.ContentUpdateInput) {
    const content = await this.contentRepository.updateRaw(orgname, id, raw);
    this.websocketsService.socket.to(orgname).emit("update");
    return content;
  }
}
