import {
  CreateBucketCommand,
  DeleteObjectCommand,
  GetObjectCommand,
  HeadObjectCommand,
  ListObjectsV2Command,
  PutObjectCommand,
  S3Client,
} from "@aws-sdk/client-s3";
import { getSignedUrl } from "@aws-sdk/s3-request-presigner";
import {
  ConflictException,
  Injectable,
  NotFoundException,
} from "@nestjs/common";
import axios from "axios";
import * as fs from "fs";
import * as path from "path";
import { Readable } from "stream";

import { StorageItemDto } from "./dto/storage-item.dto";
import { StorageService } from "./storage.service";

@Injectable()
export class S3StorageProvider implements StorageService {
  private bucketName: string;
  private expirationTime = 60 * 60 * 1000; // 1 hour in milliseconds
  private s3Client: S3Client;

  constructor() {
    this.bucketName = process.env.MINIO_BUCKET || "my-bucket"; // Replace with your bucket name
    this.s3Client = new S3Client({
      credentials: {
        accessKeyId: process.env.MINIO_ACCESS_KEY || "minioadmin",
        secretAccessKey: process.env.MINIO_SECRET_KEY || "minioadmin",
      },
      endpoint: process.env.MINIO_ENDPOINT || "http://localhost:9000",
      forcePathStyle: true, // Required for MinIO
      region: "us-east-1",
    });

    // Ensure the bucket exists
    this.createBucketIfNotExists();
  }

  async checkFileExists(orgname: string, filePath: string): Promise<boolean> {
    try {
      await this.s3Client.send(
        new HeadObjectCommand({
          Bucket: this.bucketName,
          Key: this.getKey(orgname, filePath),
        })
      );
      return true;
    } catch (error) {
      if (
        error.name === "NotFound" ||
        error.$metadata?.httpStatusCode === 404
      ) {
        return false;
      }
      throw error;
    }
  }

  async createDirectory(orgname: string, dirPath: string): Promise<void> {
    const exists = await this.checkFileExists(orgname, dirPath);
    if (exists) {
      throw new ConflictException(
        "Cannot create directory. File or path already exists at this location"
      );
    }
    const key = this.getKey(orgname, dirPath).replace(/\/?$/, "/") + "/";
    await this.s3Client.send(
      new PutObjectCommand({
        Body: "",
        Bucket: this.bucketName,
        Key: key,
      })
    );
  }

  async delete(orgname: string, filePath: string): Promise<void> {
    const exists = await this.checkFileExists(orgname, filePath);
    if (!exists) {
      throw new NotFoundException(`File at ${filePath} does not exist`);
    }
    await this.s3Client.send(
      new DeleteObjectCommand({
        Bucket: this.bucketName,
        Key: this.getKey(orgname, filePath),
      })
    );
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
    const result = await this.s3Client.send(
      new GetObjectCommand({
        Bucket: this.bucketName,
        Key: this.getKey(orgname, filePath),
      })
    );
    const stream = result.Body as Readable;
    const chunks = [];
    for await (const chunk of stream) {
      chunks.push(chunk);
    }
    const buffer = Buffer.concat(chunks);

    // Optionally save to destination
    if (destination) {
      await fs.promises.writeFile(destination, buffer);
    }

    return { buffer };
  }

  async getMetaData(
    orgname: string,
    filePath: string
  ): Promise<{ metadata: any }> {
    const exists = await this.checkFileExists(orgname, filePath);
    if (!exists) {
      throw new NotFoundException(`File at ${filePath} does not exist`);
    }
    const result = await this.s3Client.send(
      new HeadObjectCommand({
        Bucket: this.bucketName,
        Key: this.getKey(orgname, filePath),
      })
    );
    return { metadata: result.Metadata };
  }

  async getSignedUrl(
    orgname: string,
    filePath: string,
    action: "read" | "write"
  ): Promise<string> {
    if (action === "write") {
      let conflict = true;
      let i = 0;
      const originalFilePath = filePath;
      for (; i < 1000; i++) {
        conflict = await this.checkFileExists(orgname, filePath);
        if (!conflict) {
          break;
        }
        const ext = path.extname(originalFilePath);
        const baseName = path.basename(originalFilePath, ext);
        const dirName = path.dirname(originalFilePath);
        filePath = path.join(dirName, `${baseName}(${i + 1})${ext}`);
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

    const key = this.getKey(orgname, filePath);
    let command;
    if (action === "read") {
      command = new GetObjectCommand({
        Bucket: this.bucketName,
        Key: key,
      });
    } else {
      command = new PutObjectCommand({
        Bucket: this.bucketName,
        Key: key,
      });
    }

    const signedUrl = await getSignedUrl(this.s3Client, command, {
      expiresIn: this.expirationTime / 1000, // Convert milliseconds to seconds
    });

    return signedUrl;
  }

  async listDirectory(
    orgname: string,
    dirPath: string
  ): Promise<StorageItemDto[]> {
    const prefix = this.getKey(orgname, dirPath).replace(/\/?$/, "/") + "/";
    const result = await this.s3Client.send(
      new ListObjectsV2Command({
        Bucket: this.bucketName,
        Delimiter: "/",
        Prefix: prefix,
      })
    );

    const items: StorageItemDto[] = [];

    if (result.CommonPrefixes) {
      for (const commonPrefix of result.CommonPrefixes) {
        items.push(
          new StorageItemDto({
            createdAt: null,
            id: commonPrefix.Prefix,
            isDir: true,
            name: path.basename(commonPrefix.Prefix.replace(/\/$/, "")),
            size: 0,
          })
        );
      }
    }

    if (result.Contents) {
      for (const content of result.Contents) {
        if (content.Key === prefix) {
          continue; // Skip the directory placeholder object
        }
        items.push(
          new StorageItemDto({
            createdAt: content.LastModified,
            id: content.Key,
            isDir: false,
            name: path.basename(content.Key),
            size: content.Size,
          })
        );
      }
    }

    return items;
  }

  async upload(
    orgname: string,
    filePath: string,
    file: Express.Multer.File
  ): Promise<string> {
    let conflict = true;
    let i = 0;
    const originalFilePath = filePath;
    for (; i < 1000; i++) {
      conflict = await this.checkFileExists(orgname, filePath);
      if (!conflict) {
        break;
      }
      const ext = path.extname(originalFilePath);
      const baseName = path.basename(originalFilePath, ext);
      const dirName = path.dirname(originalFilePath);
      filePath = path.join(dirName, `${baseName}(${i + 1})${ext}`);
    }
    if (conflict) {
      throw new ConflictException("File already exists");
    }

    const key = this.getKey(orgname, filePath);

    await this.s3Client.send(
      new PutObjectCommand({
        Body: file.buffer,
        Bucket: this.bucketName,
        ContentType: file.mimetype,
        Key: key,
      })
    );

    const readUrl = await this.getSignedUrl(orgname, filePath, "read");
    return readUrl;
  }

  async uploadFromUrl(
    orgname: string,
    filePath: string,
    url: string
  ): Promise<string> {
    const response = await axios.get(url, { responseType: "arraybuffer" });
    const fileBuffer = Buffer.from(response.data);

    const readUrl = await this.upload(orgname, filePath, {
      buffer: fileBuffer,
      destination: "",
      encoding: "7bit",
      fieldname: "",
      filename: "",
      mimetype: response.headers["content-type"] || "application/octet-stream",
      originalname: path.basename(filePath),
      path: "",
      size: fileBuffer.length,
      stream: null,
    } as Express.Multer.File);

    return readUrl;
  }

  private async createBucketIfNotExists() {
    try {
      await this.s3Client.send(
        new CreateBucketCommand({ Bucket: this.bucketName })
      );
    } catch (error) {
      if (error.name !== "BucketAlreadyOwnedByYou") {
        throw error;
      }
    }
  }

  private getKey(orgname: string, filePath: string): string {
    return path.posix.join("storage", orgname, filePath);
  }
}
