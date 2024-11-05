import { StorageService } from "@/src/storage/storage.service";
import { Logger } from "@nestjs/common";

import { ContentService } from "../../content/content.service";
import { ContentEntity } from "../../content/entities/content.entity";
import { LoaderService } from "../../loader/loader.service";

export const processExtractText = async (
  runId: string,
  runInputContentIds: ContentEntity[],
  logger: Logger,
  contentService: ContentService,
  loaderService: LoaderService,
  storageService: StorageService
) => {
  // hit loader endpoint
  const { mimeType, preview, textContent, title } =
    await loaderService.extractUrl(content.url, 200);
  logger.log(`Extracted text from ${content.name} with ${mimeType}`);

  // update content type
  await contentService.updateRaw(content.orgname, content.id, {
    mimeType,
  });

  // update name
  if (title.indexOf("http") == -1) {
    content = await contentService.updateRaw(content.orgname, content.id, {
      name: title,
    });
  }

  const uploadPreviewPromise = (async () => {
    const previewFilename = `${content.name}-preview.png`;
    const decodedImage = Buffer.from(preview, "base64");
    const multerFile = {
      buffer: decodedImage,
      mimetype: "image/png",
      originalname: previewFilename,
      size: decodedImage.length,
    } as Express.Multer.File;
    const url = await storageService.upload(
      content.orgname,
      `contents/${content.name}-preview.png`,
      multerFile
    );
    await contentService.updateRaw(content.orgname, content.id, {
      previewImage: url,
    });
  })();

  const uploadTextChunks = (async () => {
    const start = Date.now();
    await contentService.upsertTextChunks(
      content.orgname,
      content.id,
      textContent
    );
    logger.log(
      `Upserted embeddings for ${content.name}. Completed in ${
        (Date.now() - start) / 1000
      }s`
    );
  })();

  // if any of this fail, throw an error
  await Promise.all([uploadTextChunks, uploadPreviewPromise]);
};
