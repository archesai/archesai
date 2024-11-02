import { ContentEntity } from "@/src/content/entities/content.entity";
import { SpeechService } from "@/src/speech/speech.service";
import { StorageService } from "@/src/storage/storage.service";

export const processTextToSpeech = async (
  content: ContentEntity,
  storageService: StorageService,
  speechService: SpeechService
) => {
  const text = content.text;
  const audioBuffer = await speechService.generateSpeech(text);
  const multerFile = {
    buffer: audioBuffer,
    mimetype: "audio/mpeg",
    originalname: `${content.name}.mp3`,
    size: audioBuffer.length,
  } as Express.Multer.File;
  const url = await storageService.upload(
    content.orgname,
    `contents/${content.name}.mp3`,
    multerFile
  );
  // await contentService.updateRaw(content.orgname, content.id, {
  //   audio: url,
  // });
  console.log(url);
};
