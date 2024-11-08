import {
  BadRequestException,
  Inject,
  Injectable,
  Logger,
} from "@nestjs/common";
import { Content } from "@prisma/client";
import * as mime from "mime-types";

import { BaseService } from "../common/base.service";
import { STORAGE_SERVICE, StorageService } from "../storage/storage.service";
import { WebsocketsService } from "../websockets/websockets.service";
import { ContentRepository } from "./content.repository";
import { CreateContentDto } from "./dto/create-content.dto";
import { UpdateContentDto } from "./dto/update-content.dto";
import { ContentEntity } from "./entities/content.entity";

@Injectable()
export class ContentService extends BaseService<
  ContentEntity,
  CreateContentDto,
  UpdateContentDto,
  ContentRepository,
  Content
> {
  private logger = new Logger(ContentService.name);
  constructor(
    @Inject(STORAGE_SERVICE)
    private storageService: StorageService,
    private contentRepository: ContentRepository,
    private websocketsService: WebsocketsService
  ) {
    super(contentRepository);
  }

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
      { mimeType }
    );
    this.websocketsService.socket.to(orgname).emit("update", {
      queryKey: ["organizations", orgname, "content"],
    });
    return this.toEntity(content);
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

  async incrementCredits(orgname: string, id: string, credits: number) {
    const content = await this.contentRepository.incrementCredits(id, credits);
    this.websocketsService.socket.to(orgname).emit("update", {
      queryKey: ["organizations", orgname, "content"],
    });
    return this.toEntity(content);
  }

  async populateReadUrl(content: Content) {
    const url = `https://storage.googleapis.com/archesai/storage/${content.orgname}/`;
    if (content.url?.startsWith(url)) {
      const path = content.url.replace(url, "").split("?")[0];
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
    return this.toEntity(content);
  }

  async query(
    orgname: string,
    embedding: number[],
    topK: number,
    contentIds?: string[]
  ) {
    return this.contentRepository.query(orgname, embedding, topK, contentIds);
  }

  async removeMany(orgname: string, ids: string[]) {
    return this.contentRepository.removeMany(orgname, ids);
  }

  protected toEntity(model: Content): ContentEntity {
    return new ContentEntity(model);
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
