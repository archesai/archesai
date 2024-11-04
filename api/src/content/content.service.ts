import { Inject, Injectable, Logger } from "@nestjs/common";
import { Content, Prisma } from "@prisma/client";

import { BaseService } from "../common/base.service";
import { PaginatedDto } from "../common/paginated.dto";
import { PipelinesService } from "../pipelines/pipelines.service";
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
    BaseService<
      ContentEntity,
      CreateContentDto,
      ContentQueryDto,
      UpdateContentDto
    >
{
  private logger = new Logger(ContentService.name);
  constructor(
    @Inject(STORAGE_SERVICE)
    private storageService: StorageService,
    private contentRepository: ContentRepository,
    private websocketsService: WebsocketsService,
    private pipelinesService: PipelinesService
  ) {}

  async create(orgname: string, createContentDto: CreateContentDto) {
    const content = await this.contentRepository.create(
      orgname,
      createContentDto
    );
    this.websocketsService.socket.to(orgname).emit("update");
    const contentEntity = new ContentEntity(content);
    await this.pipelinesService.runPipeline(contentEntity);
    return contentEntity;
  }

  async findAll(orgname: string, contentQueryDto: ContentQueryDto) {
    const { count, results } = await this.contentRepository.findAll(
      orgname,
      contentQueryDto
    );
    const contentEntities = results.map(
      (content) => new ContentEntity(content)
    );
    return new PaginatedDto<ContentEntity>({
      metadata: {
        limit: contentQueryDto.limit,
        offset: contentQueryDto.offset,
        totalResults: count,
      },
      results: contentEntities,
    });
  }

  async findOne(id: string) {
    const content = await this.contentRepository.findOne(id);
    const populated = await this.populateReadUrl(content);
    return new ContentEntity(populated);
  }

  async incrementCredits(orgname: string, id: string, credits: number) {
    const content = await this.contentRepository.incrementCredits(id, credits);
    this.websocketsService.socket.to(orgname).emit("update");
    return new ContentEntity(content);
  }

  async populateReadUrl(content: Content) {
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
    return new ContentEntity(content);
  }

  async updateRaw(orgname: string, id: string, raw: Prisma.ContentUpdateInput) {
    const content = await this.contentRepository.updateRaw(orgname, id, raw);
    this.websocketsService.socket.to(orgname).emit("update");
    return new ContentEntity(content);
  }
}
