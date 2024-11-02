import * as ospath from "path";

import { ContentService } from "../../content/content.service";
import { ContentEntity } from "../../content/entities/content.entity";
import { JobEntity } from "../../jobs/entities/job.entity";
import { RunpodService } from "../../runpod/runpod.service";
import { StorageService } from "../../storage/storage.service";

export const processTextToImage = async (
  content: ContentEntity,
  toolRun: JobEntity,
  runpodService: RunpodService,
  storageService: StorageService,
  contentService: ContentService
) => {
  const input = {
    input: {
      prompt: toolRun.input,
    },
  };

  const { image_url } = await runpodService.runPod(
    content.orgname,
    content.id,
    toolRun.id,
    "y55cw5fvbum8q6",
    input
  );

  const base64String = image_url.replace(/^data:image\/\w+;base64,/, "");

  // Convert the remaining base64 string to a buffer
  const buffer = Buffer.from(base64String, "base64");
  const path = `images/${content.id}.png`;

  // Use the upload function
  const url = await storageService.upload(content.orgname, path, {
    buffer: buffer,
    originalname: ospath.basename(path),
    size: buffer.length,
  } as Express.Multer.File);

  await contentService.updateRaw(content.orgname, content.id, {
    previewImage: url,
    url,
  });
};
