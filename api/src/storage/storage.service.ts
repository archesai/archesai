import { StorageItemDto } from "./dto/storage-item.dto";

export const STORAGE_SERVICE = "STORAGE_SERVICE";

export interface StorageService {
  createDirectory(orgname: string, path: string): Promise<void>;
  delete(orgname: string, path: string): Promise<void>;
  download(
    orgname: string,
    path: string,
    destination?: string,
  ): Promise<{ buffer: Buffer }>;
  getMetaData(orgname: string, path: string): Promise<{ metadata: any }>;
  getSignedUrl(
    orgname: string,
    path: string,
    action: "read" | "write",
  ): Promise<string>;
  listDirectory(orgname: string, path: string): Promise<StorageItemDto[]>;
  upload(
    orgname: string,
    path: string,
    file: Express.Multer.File,
  ): Promise<string>;
  uploadFromUrl(orgname: string, path: string, url: string): Promise<string>;
}
