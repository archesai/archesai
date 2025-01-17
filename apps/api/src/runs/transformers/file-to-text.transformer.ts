import { Logger } from '@nestjs/common'

import { ContentService } from '../../content/content.service'
import { ContentEntity } from '../../content/entities/content.entity'
import { IToolRunProcess } from '../interfaces/tool-run-processor.interface'
import { UnstructuredLoader } from '@langchain/community/document_loaders/fs/unstructured'
import { ConfigService } from '@/src/config/config.service'
import { IStorageService } from '@/src/storage/interfaces/storage-provider.interface'

export const transformFileToText: IToolRunProcess = async (
  runId: string,
  inputs: ContentEntity[],
  logger: Logger,
  contentService: ContentService,
  configService: ConfigService,
  storageService: IStorageService
): Promise<ContentEntity[]> => {
  logger.log(`Extracting text for run ${runId} with url ${inputs[0].url}`)

  let content = inputs[0]
  if (!content.url) {
    throw new Error('No url provided')
  }
  const { buffer } = await storageService.download(content.orgname, content.url)
  const loader = new UnstructuredLoader(
    {
      buffer: buffer,
      fileName: content.url
    },
    {
      apiUrl: configService.get('unstructured.endpoint')
    }
  )

  const docs = await loader.load()

  const textContent = docs.map((doc) => {
    return {
      text: doc.pageContent,
      page: doc.metadata.page,
      tokens: doc.metadata.tokens
    }
  })

  const title = 'UNIMPLEMENTED TITLE'

  const sanitizedTextContent = textContent.map((data) => ({
    ...data,
    text: data.text
      .replaceAll(/\0/g, '')
      .replaceAll(/[^ -~\u00A0-\uD7FF\uE000-\uFDCF\uFDF0-\uFFFD\n]/g, '')
  }))

  // update name
  if (title?.indexOf('http') == -1) {
    content = await contentService.update(content.id, {
      name: title
    })
  }

  const chunkedContent = await Promise.all(
    sanitizedTextContent.map((data, i) =>
      contentService.create({
        name: `${content.name} - Page ${data.page} - Index ${i}`,
        text: data.text,
        labels: [],
        orgname: content.orgname,
        url: null
      })
    )
  )

  return chunkedContent
}
