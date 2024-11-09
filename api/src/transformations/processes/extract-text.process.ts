import { StorageService } from "@/src/storage/storage.service";
import { HttpService } from "@nestjs/axios";
import { BadRequestException, Logger } from "@nestjs/common";
import { ConfigService } from "@nestjs/config";
import { AxiosError } from "axios";
import { catchError, firstValueFrom } from "rxjs";

import { ContentService } from "../../content/content.service";
import { ContentEntity } from "../../content/entities/content.entity";

export const processExtractText = async (
  runId: string,
  runInputContents: ContentEntity[],
  logger: Logger,
  contentService: ContentService,
  storageService: StorageService,
  httpService: HttpService,
  configService: ConfigService
): Promise<ContentEntity[]> => {
  logger.log(`Extracting text for run ${runId}`);

  let content = runInputContents[0];
  const { data } = await firstValueFrom(
    httpService
      .post(configService.get("LOADER_ENDPOINT"), {
        url: content.url,
      })
      .pipe(
        catchError((err: AxiosError) => {
          logger.error("Error hitting loader endpoint: " + err.message);
          throw new BadRequestException();
        })
      )
  );
  const { preview, title } = data as {
    contentType: string;
    preview: string;
    textContent: { page: number; text: string; tokens: number }[];
    title: string;
  };

  logger.log(`Extracted text for ${content.name}`);

  // const sanitizedTextContent = textContent.map((data) => ({
  //   ...data,
  //   text: data.text
  //     .replaceAll(/\0/g, "")
  //     .replaceAll(/[^ -~\u00A0-\uD7FF\uE000-\uFDCF\uFDF0-\uFFFD\n]/g, ""),
  // }));

  // const totalTokens = textContent.reduce((acc, curr) => acc + curr.tokens, 0);

  // update name
  if (title.indexOf("http") == -1) {
    content = await contentService.setTitle(content.orgname, content.id, title);
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
    await contentService.setPreviewImage(content.orgname, content.id, url);
  })();

  const uploadTextChunks = (async () => {
    // const start = Date.now();
    // await contentService.upsertTextChunks(
    //   content.orgname,
    //   content.id,
    //   textContent
    // );
    // logger.log(
    //   `Upserted embeddings for ${content.name}. Completed in ${
    //     (Date.now() - start) / 1000
    //   }s`
    // ); // FIXME
  })();

  // if any of this fail, throw an error
  await Promise.all([uploadTextChunks, uploadPreviewPromise]);
  return;
};
