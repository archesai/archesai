import {
  ConflictException,
  Injectable,
  NotFoundException,
} from "@nestjs/common";
import axios from "axios";
import * as fs from "fs";
import * as os from "os";
import * as path from "path";

import { StorageItemDto } from "./dto/storage-item.dto";
import { StorageService } from "./storage.service";

@Injectable()
export class LocalStorageService implements StorageService {
  private baseDir: string;
  private expirationTime = 60 * 60 * 1000; // Not used in mock but kept for consistency

  constructor() {
    // Set the base directory for local storage (e.g., project_root/localstorage)
    this.baseDir = path.join(process.cwd(), "localstorage");
  }

  async checkFileExists(orgname: string, filePath: string): Promise<boolean> {
    const fullPath = this.getFullPath(orgname, filePath);
    try {
      await fs.promises.access(fullPath);
      return true;
    } catch {
      return false;
    }
  }

  async createDirectory(orgname: string, dirPath: string): Promise<void> {
    const fullPath = this.getFullPath(orgname, dirPath);
    const exists = await this.checkFileExists(orgname, dirPath);
    if (exists) {
      throw new ConflictException(
        "Cannot create directory. File or path already exists at this location"
      );
    }
    await fs.promises.mkdir(fullPath, { recursive: true });
  }

  async delete(orgname: string, filePath: string): Promise<void> {
    const fullPath = this.getFullPath(orgname, filePath);
    const exists = await this.checkFileExists(orgname, filePath);
    if (!exists) {
      throw new NotFoundException(`File at ${filePath} does not exist`);
    }
    const stat = await fs.promises.lstat(fullPath);
    if (stat.isDirectory()) {
      await fs.promises.rmdir(fullPath, { recursive: true });
    } else {
      await fs.promises.unlink(fullPath);
    }
  }

  async download(
    orgname: string,
    filePath: string,
    destination?: string
  ): Promise<{ buffer: Buffer }> {
    const fullPath = this.getFullPath(orgname, filePath);
    const exists = await this.checkFileExists(orgname, filePath);
    if (!exists) {
      throw new NotFoundException(`File at ${filePath} does not exist`);
    }
    const stat = await fs.promises.lstat(fullPath);
    if (stat.isDirectory()) {
      throw new Error("Cannot download a directory");
    }
    const buffer = await fs.promises.readFile(fullPath);
    if (destination) {
      await fs.promises.writeFile(destination, buffer);
    }
    return { buffer };
  }

  async getMetaData(
    orgname: string,
    filePath: string
  ): Promise<{ metadata: any }> {
    const fullPath = this.getFullPath(orgname, filePath);
    const exists = await this.checkFileExists(orgname, filePath);
    if (!exists) {
      throw new NotFoundException(`File at ${filePath} does not exist`);
    }
    const stats = await fs.promises.stat(fullPath);
    return { metadata: stats };
  }

  async getSignedUrl(
    orgname: string,
    filePath: string,
    action: "read" | "write"
  ): Promise<string> {
    const fullPath = this.getFullPath(orgname, filePath);

    if (action === "write") {
      let conflict = true;
      let i = 0;
      const originalFilePath = filePath;
      while (conflict && i < 1000) {
        conflict = await this.checkFileExists(orgname, filePath);
        if (conflict) {
          const ext = path.extname(originalFilePath);
          const baseName = path.basename(originalFilePath, ext);
          const dirName = path.dirname(originalFilePath);
          filePath = path.join(dirName, `${baseName}(${++i})${ext}`);
        }
      }
      if (conflict) {
        throw new ConflictException("File already exists");
      }
      return `file://${this.getFullPath(orgname, filePath)}`;
    } else {
      const exists = await this.checkFileExists(orgname, filePath);
      if (!exists) {
        throw new NotFoundException(`File at ${filePath} does not exist`);
      }
      return `file://${fullPath}`;
    }
  }

  async listDirectory(
    orgname: string,
    dirPath: string
  ): Promise<StorageItemDto[]> {
    const fullPath = this.getFullPath(orgname, dirPath);
    const exists = await this.checkFileExists(orgname, dirPath);
    if (!exists) {
      throw new NotFoundException(`Directory at ${dirPath} does not exist`);
    }
    const items = await fs.promises.readdir(fullPath);

    const result: StorageItemDto[] = [];
    for (const item of items) {
      const itemFullPath = path.join(fullPath, item);
      const stats = await fs.promises.stat(itemFullPath);
      const isDir = stats.isDirectory();
      result.push(
        new StorageItemDto({
          createdAt: stats.birthtime,
          id: itemFullPath,
          isDir: isDir,
          name: item,
          size: stats.size,
        })
      );
    }
    return result;
  }

  async upload(
    orgname: string,
    filePath: string,
    file: Express.Multer.File
  ): Promise<string> {
    let conflict = true;
    let i = 0;
    const originalFilePath = filePath;
    while (conflict && i < 1000) {
      conflict = await this.checkFileExists(orgname, filePath);
      if (conflict) {
        const ext = path.extname(originalFilePath);
        const baseName = path.basename(originalFilePath, ext);
        const dirName = path.dirname(originalFilePath);
        filePath = path.join(dirName, `${baseName}(${++i})${ext}`);
      }
    }
    if (conflict) {
      throw new ConflictException("File already exists");
    }

    const fullPath = this.getFullPath(orgname, filePath);
    const dirName = path.dirname(fullPath);

    await fs.promises.mkdir(dirName, { recursive: true });
    await fs.promises.writeFile(fullPath, file.buffer);
    return `file://${fullPath}`;
  }

  async uploadFromUrl(
    orgname: string,
    filePath: string,
    url: string
  ): Promise<string> {
    const response = await axios.get(url, { responseType: "arraybuffer" });
    const fileBuffer = Buffer.from(response.data);

    const tempFileName = path.basename(url);
    const tempFilePath = path.join(os.tmpdir(), tempFileName);

    await fs.promises.writeFile(tempFilePath, fileBuffer);

    const file: Express.Multer.File = {
      buffer: fileBuffer,
      destination: "",
      encoding: "",
      fieldname: "",
      filename: "",
      mimetype: "",
      originalname: tempFileName,
      path: "",
      size: fileBuffer.length,
      stream: null,
    };

    try {
      const readUrl = await this.upload(orgname, filePath, file);
      await fs.promises.unlink(tempFilePath);
      return readUrl;
    } catch (err) {
      await fs.promises.unlink(tempFilePath);
      throw err;
    }
  }

  // Helper method to get the full path on the local filesystem
  private getFullPath(orgname: string, filePath: string): string {
    return path.join(this.baseDir, "storage", orgname, filePath);
  }
}
