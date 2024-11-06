import {
  BadRequestException,
  Inject,
  Injectable,
  Logger,
} from "@nestjs/common";
import { Content, Prisma } from "@prisma/client";
import * as mime from "mime-types";

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
    let mimeType: string;
    if (createContentDto.url) {
      mimeType = await this.detectMimeTypeFromUrl(createContentDto.url);
    } else {
      mimeType = "text/plain";
    }
    const content = await this.contentRepository.create(
      orgname,
      createContentDto,
      mimeType
    );
    this.websocketsService.socket.to(orgname).emit("update", {
      queryKey: ["organizations", orgname, "content"],
    });
    return new ContentEntity(content);
  }

  async detectMimeTypeFromUrl(url: string) {
    try {
      // Extract the file name from the URL
      const urlObj = new URL(url);
      const pathname = urlObj.pathname;
      const fileName = pathname.split("/").pop();

      if (!fileName) {
        throw new BadRequestException("Unable to extract file name from URL");
      }

      // Get MIME type based on file extension
      const mimeType = mime.lookup(fileName);
      return mimeType || null;
    } catch (error) {
      throw new BadRequestException("Failed to detect MIME type");
    }
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
    this.websocketsService.socket.to(orgname).emit("update", {
      queryKey: ["organizations", orgname, "content"],
    });
    return new ContentEntity(content);
  }
  async populateReadUrl(content: Content) {
    if (
      content.url?.startsWith(
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

  async query(
    orgname: string,
    embedding: number[],
    topK: number,
    contentIds?: string[]
  ) {
    return this.contentRepository.query(orgname, embedding, topK, contentIds);
  }

  async remove(orgname: string, contentId: string): Promise<void> {
    await this.contentRepository.remove(orgname, contentId);
    this.websocketsService.socket.to(orgname).emit("update", {
      queryKey: ["organizations", orgname, "content"],
    });
  }

  async removeMany(orgname: string, ids: string[]) {
    return this.contentRepository.removeMany(orgname, ids);
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
    this.websocketsService.socket.to(orgname).emit("update", {
      queryKey: ["organizations", orgname, "content"],
    });
    return new ContentEntity(content);
  }

  async updateRaw(orgname: string, id: string, raw: Prisma.ContentUpdateInput) {
    const content = await this.contentRepository.updateRaw(orgname, id, raw);
    this.websocketsService.socket.to(orgname).emit("update", {
      queryKey: ["organizations", orgname, "content"],
    });
    return new ContentEntity(content);
  }

  async upsertTextChunks(
    orgname: string,
    contentId: string,
    records: {
      text: string;
    }[]
  ): Promise<void> {
    return this.contentRepository.upsertTextChunks(orgname, contentId, records);
  }

  async upsertVectors(
    orgname: string,
    contentId: string,
    records: {
      embedding: number[];
      textChunkId: string;
    }[]
  ): Promise<void> {
    return this.contentRepository.upsertVectors(orgname, contentId, records);
  }
}
