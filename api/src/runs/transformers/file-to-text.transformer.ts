import { HttpService } from '@nestjs/axios'
import { BadRequestException, Logger } from '@nestjs/common'
import { ConfigService } from '@nestjs/config'
import { AxiosError } from 'axios'
import { catchError, firstValueFrom } from 'rxjs'

import { ContentService } from '../../content/content.service'
import { ContentEntity } from '../../content/entities/content.entity'
import { IToolRunProcess } from '../interfaces/tool-run-processor.interface'

export const transformFileToText: IToolRunProcess = async (
  runId: string,
  runInputContents: ContentEntity[],
  logger: Logger,
  contentService: ContentService,
  httpService: HttpService,
  configService: ConfigService
): Promise<ContentEntity[]> => {
  logger.log(`Extracting text for run ${runId} with url ${runInputContents[0].url}`)

  let content = runInputContents[0]
  const { data } = await firstValueFrom(
    httpService
      .post(configService.get('LOADER_ENDPOINT') + '/indexDocument', {
        url: content.url
      })
      .pipe(
        catchError((err: AxiosError) => {
          logger.error('Error hitting loader endpoint: ' + err.message)
          throw new BadRequestException()
        })
      )
  )
  const { textContent, title } = data as {
    contentType: string
    preview: string
    textContent: { page: number; text: string; tokens: number }[]
    title: string
  }

  logger.log(`Extracted text for ${content.name}`)

  const sanitizedTextContent = textContent.map((data) => ({
    ...data,
    text: data.text.replaceAll(/\0/g, '').replaceAll(/[^ -~\u00A0-\uD7FF\uE000-\uFDCF\uFDF0-\uFFFD\n]/g, '')
  }))

  // update name
  if (title?.indexOf('http') == -1) {
    content = await contentService.setTitle(content.orgname, content.id, title)
  }

  const chunkedContent = await Promise.all(
    sanitizedTextContent.map((data, i) =>
      contentService.create(content.orgname, {
        name: `${content.name} - Page ${data.page} - Index ${i}`,
        text: data.text
      })
    )
  )

  return chunkedContent
}
