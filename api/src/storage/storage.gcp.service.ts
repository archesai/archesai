import { Storage } from "@google-cloud/storage";
import {
  ConflictException,
  Injectable,
  NotFoundException,
} from "@nestjs/common";
import axios from "axios";
import * as path from "path";

import { archesaiSa } from "./archesai-sa";
import { StorageItemDto } from "./dto/storage-item.dto";
import { StorageService } from "./storage.service";

@Injectable()
export class GoogleCloudStorageService implements StorageService {
  private readonly bucketName: string;
  private readonly expirationTime = 60 * 60 * 1000; // 1 hour in milliseconds
  private readonly storage: Storage;

  constructor() {
    this.storage = new Storage({
      credentials: {
        ...archesaiSa,
      },
      projectId: "archesai",
    });
    this.bucketName = "archesai";
  }

  private getFilePath(orgname: string, filePath: string): string {
    return path.posix.join("storage", orgname, filePath);
  }

  async checkFileExists(orgname: string, filePath: string): Promise<boolean> {
    const [exists] = await this.storage
      .bucket(this.bucketName)
      .file(this.getFilePath(orgname, filePath))
      .exists();
    return exists;
  }

  async createDirectory(orgname: string, dirPath: string): Promise<void> {
    const exists = await this.checkFileExists(orgname, dirPath);
    if (exists) {
      throw new ConflictException(
        "Cannot create directory. File or path already exists at this location"
      );
    }
    await this.storage
      .bucket(this.bucketName)
      .file(this.getFilePath(orgname, dirPath) + "/")
      .save("");
  }

  async delete(orgname: string, filePath: string): Promise<void> {
    const exists = await this.checkFileExists(orgname, filePath);
    if (!exists) {
      throw new NotFoundException(`File at ${filePath} does not exist`);
    }
    await this.storage
      .bucket(this.bucketName)
      .file(this.getFilePath(orgname, filePath))
      .delete();
  }

  async download(
    orgname: string,
    filePath: string,
    destination?: string
  ): Promise<{ buffer: Buffer }> {
    const exists = await this.checkFileExists(orgname, filePath);
    if (!exists) {
      throw new NotFoundException(`File at ${filePath} does not exist`);
    }
    const [buffer] = await this.storage
      .bucket(this.bucketName)
      .file(this.getFilePath(orgname, filePath))
      .download({ destination });
    return { buffer };
  }

  async getMetaData(orgname: string, filePath: string) {
    const exists = await this.checkFileExists(orgname, filePath);
    if (!exists) {
      throw new NotFoundException(`File at ${filePath} does not exist`);
    }
    const [metadata] = await this.storage
      .bucket(this.bucketName)
      .file(this.getFilePath(orgname, filePath))
      .getMetadata();
    return { metadata };
  }

  async getSignedUrl(
    orgname: string,
    filePath: string,
    action: "read" | "write"
  ): Promise<string> {
    let fullPath = this.getFilePath(orgname, filePath);
    if (action === "write") {
      let conflict = true;
      let i = 0;
      for (; i < 1000; i++) {
        conflict = await this.checkFileExists(orgname, filePath);
        if (!conflict) {
          break;
        }
        filePath = filePath.replace(/(\.[\w\d_-]+)$/i, `(${i})$1`);
        fullPath = this.getFilePath(orgname, filePath);
      }
      if (conflict) {
        throw new ConflictException("File already exists");
      }
    } else {
      const exists = await this.checkFileExists(orgname, filePath);
      if (!exists) {
        throw new NotFoundException(`File at ${filePath} does not exist`);
      }
    }

    const [url] = await this.storage
      .bucket(this.bucketName)
      .file(fullPath)
      .getSignedUrl({
        action: action,
        expires: Date.now() + this.expirationTime,
        version: "v4",
      });

    return url;
  }

  async listDirectory(
    orgname: string,
    dirPath: string
  ): Promise<StorageItemDto[]> {
    const fullPath =
      this.getFilePath(orgname, dirPath).replace(/\/+$/, "") + "/";

    const [files] = await this.storage.bucket(this.bucketName).getFiles({
      delimiter: "/",
      prefix: fullPath,
    });

    const directories = new Set<string>();
    const fileItems: StorageItemDto[] = [];

    files.forEach((file) => {
      const relativePath = file.name.slice(fullPath.length);
      if (relativePath.endsWith("/")) {
        const dirName = relativePath.split("/")[0];
        directories.add(dirName);
      } else if (relativePath) {
        fileItems.push(
          new StorageItemDto({
            createdAt: new Date(file.metadata.timeCreated),
            id: file.id,
            isDir: false,
            name: relativePath,
            size: Number(file.metadata.size),
          })
        );
      }
    });

    const directoryItems = Array.from(directories).map(
      (dirName) =>
        new StorageItemDto({
          createdAt: null,
          id: `${fullPath}${dirName}/`,
          isDir: true,
          name: dirName + "/",
          size: 0,
        })
    );

    return [...directoryItems, ...fileItems];
  }

  async upload(
    orgname: string,
    filePath: string,
    file: Express.Multer.File
  ): Promise<string> {
    let conflict = await this.checkFileExists(orgname, filePath);
    const originalPath = filePath;
    let i = 1;
    while (conflict && i < 1000) {
      filePath = originalPath.replace(/(\.[\w\d_-]+)$/i, `(${i})$1`);
      conflict = await this.checkFileExists(orgname, filePath);
      i++;
    }
    if (conflict) {
      throw new ConflictException("File already exists");
    }

    const ref = this.storage
      .bucket(this.bucketName)
      .file(this.getFilePath(orgname, filePath));

    await ref.save(file.buffer, {
      contentType: file.mimetype,
      metadata: {
        metadata: {
          originalName: file.originalname,
        },
      },
    });

    return this.getSignedUrl(orgname, filePath, "read");
  }

  async uploadFromUrl(
    orgname: string,
    filePath: string,
    url: string
  ): Promise<string> {
    let conflict = await this.checkFileExists(orgname, filePath);
    const originalPath = filePath;
    let i = 1;
    while (conflict && i < 1000) {
      filePath = originalPath.replace(/(\.[\w\d_-]+)$/i, `(${i})$1`);
      conflict = await this.checkFileExists(orgname, filePath);
      i++;
    }
    if (conflict) {
      throw new ConflictException("File already exists");
    }

    const ref = this.storage
      .bucket(this.bucketName)
      .file(this.getFilePath(orgname, filePath));

    const response = await axios({
      method: "get",
      responseType: "stream",
      url: url,
    });

    const writeStream = ref.createWriteStream();

    response.data.pipe(writeStream);

    await new Promise<void>((resolve, reject) => {
      writeStream.on("finish", resolve);
      writeStream.on("error", reject);
    });

    return this.getSignedUrl(orgname, filePath, "read");
  }
}
