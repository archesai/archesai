import { ContentService } from "@/src/content/content.service";
import { ContentEntity } from "@/src/content/entities/content.entity";
import { SpeechService } from "@/src/speech/speech.service";
import { StorageService } from "@/src/storage/storage.service";
import { Logger } from "@nestjs/common";

export const transformTextToSpeech = async (
  runId: string,
  runInputContents: ContentEntity[],
  logger: Logger,
  contentService: ContentService,
  storageService: StorageService,
  speechService: SpeechService
): Promise<ContentEntity[]> => {
  logger.log(`Processing text to speech for run ${runId}`);
  const audioBuffer = await speechService.generateSpeech(
    runInputContents.map((x) => x.text).join(" ")
  );

  const multerFile = {
    buffer: audioBuffer,
    mimetype: "audio/mpeg",
    originalname: `${runId}.mp3`,
    size: audioBuffer.length,
  } as Express.Multer.File;
  const url = await storageService.upload(
    runInputContents[0].orgname,
    `contents/${runId}.mp3`,
    multerFile
  );

  const content = await contentService.create(runInputContents[0].orgname, {
    name:
      "Text to Speech Tool -" + runInputContents.map((x) => x.name).join(", "),
    url,
  });

  return [new ContentEntity(content)];
};
