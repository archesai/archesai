import { ContentService } from '@/src/content/content.service'
import { ContentEntity } from '@/src/content/entities/content.entity'
import { SpeechService } from '@/src/speech/speech.service'
import { Logger } from '@nestjs/common'

import { IToolRunProcess } from '../interfaces/tool-run-processor.interface'
import { IStorageService } from '@/src/storage/interfaces/storage-provider.interface'

export const transformTextToSpeech: IToolRunProcess = async (
  runId: string,
  inputs: ContentEntity[],
  logger: Logger,
  contentService: ContentService,
  storageService: IStorageService,
  speechService: SpeechService
): Promise<ContentEntity[]> => {
  logger.log(`Processing text to speech for run ${runId}`)
  const audioBuffer = await speechService.generateSpeech(
    inputs.map((x) => x.text).join(' ')
  )

  const multerFile = {
    buffer: audioBuffer,
    mimetype: 'audio/mpeg',
    originalname: `${runId}.mp3`,
    size: audioBuffer.length
  } as Express.Multer.File
  const url = await storageService.upload(
    inputs[0].orgname,
    `contents/${runId}.mp3`,
    multerFile
  )

  const content = await contentService.create({
    name: 'Text to Speech Tool -' + inputs.map((x) => x.name).join(', '),
    url,
    labels: [],
    orgname: inputs[0].orgname,
    text: null
  })

  return [new ContentEntity(content)]
}
