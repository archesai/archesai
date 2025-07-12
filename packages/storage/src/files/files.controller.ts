import type { BaseService, Controller } from '@archesai/core'
import type { FileEntity } from '@archesai/schemas'

import { BaseController } from '@archesai/core'
import {
  CreateSignedUrlDtoSchema,
  FILE_ENTITY_KEY,
  FileEntitySchema
} from '@archesai/schemas'

import type { FilesService } from '#files/files.service'

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
      CreateSignedUrlDtoSchema,
      CreateSignedUrlDtoSchema,
      filesService as unknown as BaseService<FileEntity>
    )
  }
}
