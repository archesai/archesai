import { ContentService } from "@/src/content/content.service";
import { Logger } from "@nestjs/common";
import * as ospath from "path";

import { ContentEntity } from "../../content/entities/content.entity";
import { RunpodService } from "../../runpod/runpod.service";
import { StorageService } from "../../storage/storage.service";

export const processTextToImage = async (
  runId: string,
  runInputContents: ContentEntity[],
  logger: Logger,
  contentService: ContentService,
  runpodService: RunpodService,
  storageService: StorageService
): Promise<ContentEntity[]> => {
  logger.log(`Processing text to image for run ${runId}`);
  const input = {
    input: {
      prompt: "a man running in a circle",
    },
  };

  const orgname = runInputContents[0].orgname;

  const { image_url } = await runpodService.runPod(
    runId,
    "y55cw5fvbum8q6",
    input
  );

  const base64String = image_url.replace(/^data:image\/\w+;base64,/, "");

  // Convert the remaining base64 string to a buffer
  const buffer = Buffer.from(base64String, "base64");
  const path = `images/${runId}.png`;

  // Use the upload function
  const url = await storageService.upload(orgname, path, {
    buffer: buffer,
    originalname: ospath.basename(path),
    size: buffer.length,
  } as Express.Multer.File);
  logger.log(`Text to image completed and uploaded for run ${runId}`);

  const content = await contentService.create(runInputContents[0].orgname, {
    name:
      "Text to Speech Tool -" + runInputContents.map((x) => x.name).join(", "),
    url,
  });

  return [new ContentEntity(content)];
};
