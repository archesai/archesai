import { DownloadResponse, Storage } from "@google-cloud/storage";
import {
  ConflictException,
  Injectable,
  NotFoundException,
} from "@nestjs/common";
import axios from "axios";
import * as fs from "fs";
import * as os from "os";
import * as ospath from "path";

import { archesaiSa } from "./archesai-sa";
import { StorageItemDto } from "./dto/storage-item.dto";
import { StorageService } from "./storage.service";

@Injectable()
export class GoogleCloudStorageService implements StorageService {
  private bucket: string;
  private expirationTime = 60 * 60 * 1000;
  private storage: Storage;

  constructor() {
    this.storage = new Storage({
      credentials: {
        ...archesaiSa,
      },
      projectId: "archesai",
    });
    this.bucket = "archesai";
  }

  async checkFileExists(orgname: string, path: string): Promise<boolean> {
    const [exists] = await this.storage
      .bucket(this.bucket)
      .file(ospath.join("storage", orgname, path))
      .exists();

    return exists;
  }

  async createDirectory(orgname: string, path: string): Promise<void> {
    const exists = await this.checkFileExists(orgname, path);
    if (exists) {
      throw new ConflictException(
        "Cannot create directory. File or path already exists at this location"
      );
    }
    await this.storage
      .bucket(this.bucket)
      .file(ospath.join("storage", orgname, path) + "/")
      .save("");
  }

  async delete(orgname: string, path: string) {
    const exists = await this.checkFileExists(orgname, path);
    if (!exists) {
      throw new NotFoundException(`File at ${path} does not exist`);
    }
    await this.storage
      .bucket(this.bucket)
      .file(ospath.join("storage", orgname, path))
      .delete({ ignoreNotFound: true });
  }

  async download(orgname: string, path: string, destination?: string) {
    const exists = await this.checkFileExists(orgname, path);
    if (!exists) {
      throw new NotFoundException(`File at ${path} does not exist`);
    }
    const fileResponse: DownloadResponse = await this.storage
      .bucket(this.bucket)
      .file(ospath.join("storage", orgname, path))
      .download({ destination });
    const [buffer] = fileResponse;

    return {
      buffer,
    };
  }

  async getMetaData(orgname: string, path: string) {
    const exists = await this.checkFileExists(orgname, path);
    if (!exists) {
      throw new NotFoundException(`File at ${path} does not exist`);
    }
    const [metadata] = await this.storage
      .bucket(this.bucket)
      .file(ospath.join("storage", orgname, path))
      .getMetadata();

    return {
      metadata,
    };
  }

  async getSignedUrl(orgname: string, path: string, action: "read" | "write") {
    if (action === "write") {
      let conflict = true;
      let i = 0;
      for (; i < 1000; i++) {
        conflict = await this.checkFileExists(orgname, path);
        if (!conflict) {
          break;
        }
        path = path.replace(/(\.[\w\d_-]+)$/i, `(${++i})$1`);
      }
      if (conflict) {
        throw new ConflictException("File already exists");
      }
    } else {
      const exists = await this.checkFileExists(orgname, path);
      if (!exists) {
        throw new NotFoundException(`File at ${path} does not exist`);
      }
    }

    const [url] = await this.storage
      .bucket(this.bucket)
      .file(ospath.join("storage", orgname, path))
      .getSignedUrl({
        action: action,
        expires: Date.now() + this.expirationTime,
        version: "v4",
      });

    return url;
  }

  // this function should list all of the files in the directory
  async listDirectory(orgname: string, path: string) {
    const [files, , apiResponse] = await this.storage
      .bucket(this.bucket)
      .getFiles({
        autoPaginate: false,
        delimiter: "/",
        prefix: ospath.join("storage", orgname, path),
      });

    const directories = (apiResponse as any).prefixes || [];
    const directoriesInDir = directories.map(
      (dir) =>
        new StorageItemDto({
          createdAt: null,
          id: dir,
          isDir: true,
          name: dir.split("/").at(-2) + "/",
          size: 0,
        })
    );

    const fileDetails = await Promise.all(
      files.map(async (file) => {
        const metadata = await file.getMetadata();
        return new StorageItemDto({
          createdAt: new Date(metadata[0].timeCreated),
          id: file.id,
          isDir: false,
          name: file.name.split("/").at(-1),
          size: Number(metadata[0].size),
        });
      })
    );

    return [
      ...fileDetails.filter((file) => file.size > 0),
      ...directoriesInDir,
    ];
  }

  async upload(orgname: string, path: string, file: Express.Multer.File) {
    let conflict = true;
    for (let i = 0; i < 1000; i++) {
      let i = 0;
      conflict = await this.checkFileExists(orgname, path);
      if (conflict) {
        // put the number before the extension
        // path = path + `(${++i})`;
        path = path.replace(/(\.[\w\d_-]+)$/i, `(${++i})$1`);
        continue;
      }
      break;
    }

    const ref = this.storage
      .bucket(this.bucket)
      .file(ospath.join("storage", orgname, path));
    const stream = ref.createWriteStream();
    stream.end(file.buffer);
    await new Promise((resolve) => stream.on("finish", resolve));
    const read = await this.getSignedUrl(orgname, path, "read");

    return read;
  }

  async uploadFromUrl(orgname: string, path: string, url: string) {
    const tmpPath = ospath.join(os.tmpdir(), "tmpfile");
    const response = await axios.get(url, {
      responseType: "arraybuffer",
    });
    fs.writeFileSync(tmpPath, response.data);
    try {
      const read = await this.upload(orgname, path, {
        buffer: fs.readFileSync(tmpPath),
        originalname: ospath.basename(path),
        size: fs.statSync(tmpPath).size,
      } as Express.Multer.File);
      fs.unlinkSync(tmpPath);
      return read;
    } catch (err) {
      fs.unlinkSync(tmpPath);
      throw err;
    }
  }
}
