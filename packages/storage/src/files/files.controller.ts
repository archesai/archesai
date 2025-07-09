import { Type } from '@sinclair/typebox'

import type { BaseService, Controller } from '@archesai/core'
import type { FileEntity } from '@archesai/schemas'

import { BaseController } from '@archesai/core'
import { FILE_ENTITY_KEY, FileEntitySchema } from '@archesai/schemas'

import type { FilesService } from '#files/files.service'

import { CreateSignedUrlRequestSchema } from '#storage/dto/create-signed-url.req.dto'

/**
 * Controller for files.
 */
export class FilesController
  extends BaseController<FileEntity>
  implements Controller
{
  constructor(filesService: FilesService) {
    super(
      FILE_ENTITY_KEY,
      FileEntitySchema,
      CreateSignedUrlRequestSchema,
      Type.Object({}),
      filesService as unknown as BaseService<FileEntity>
    )
  }
}
