import { Body, Controller, Delete, Get, Inject, Query } from '@nestjs/common'
import { Param, Post } from '@nestjs/common'
import {
  ApiBearerAuth,
  ApiOperation,
  ApiResponse,
  ApiTags
} from '@nestjs/swagger'

import { Roles } from '../auth/decorators/roles.decorator'
import { PathDto } from './dto/path.dto'
import { ReadUrlDto } from './dto/read-url.dto'
import { StorageItemDto } from './dto/storage-item.dto'
import { WriteUrlDto } from './dto/write-url.dto'
import { STORAGE_SERVICE, StorageService } from './storage.service'

@ApiBearerAuth()
@ApiTags('Storage')
@Controller('/organizations/:orgname/storage')
@Roles('ADMIN')
export class StorageController {
  constructor(
    @Inject(STORAGE_SERVICE) private readonly storageService: StorageService
  ) {}

  @ApiOperation({
    description:
      "This endpoint will delete a file or directory in the organization's secure storage at the specified path. ADMIN ONLY.",
    summary: 'Delete file or directory'
  })
  @ApiResponse({ description: 'Not Found', status: 404 })
  @ApiResponse({ description: 'Unauthorized', status: 401 })
  @ApiResponse({
    description: 'Path was successfully deleted',
    status: 200
  })
  @ApiResponse({ description: 'Forbidden', status: 403 })
  @Delete('delete')
  async delete(@Param('orgname') orgname: string, @Query('path') path: string) {
    await this.storageService.delete(orgname, path)
  }

  @ApiOperation({
    description:
      "This endpoint will return a url for reading a file in the organization's secure storage. It will be valid for 15 minutes. ADMIN ONLY.",
    summary: 'Read file'
  })
  @ApiResponse({ description: 'Not Found', status: 404 })
  @ApiResponse({ description: 'Unauthorized', status: 401 })
  @ApiResponse({
    description: 'Read  url was successfully created',
    status: 201,
    type: ReadUrlDto
  })
  @ApiResponse({ description: 'Forbidden', status: 403 })
  @Post('read')
  async getReadUrl(
    @Param('orgname') orgname: string,
    @Body() pathDto: PathDto
  ) {
    const read = await this.storageService.getSignedUrl(
      orgname,
      pathDto.path,
      'read'
    )

    return { read } as ReadUrlDto
  }

  @ApiOperation({
    description:
      "This endpoint will return a url for writing to a file location in the organization's secure storage. You must write your file to the url returned by this endpoint. If you use is isDir param, it will create a directory instead of a file and you do not need to write to the url. ADMIN ONLY.",
    summary: 'Write file'
  })
  @ApiResponse({ description: 'Not Found', status: 404 })
  @ApiResponse({ description: 'Unauthorized', status: 401 })
  @ApiResponse({
    description: 'Write urls was successfully created',
    status: 201,
    type: WriteUrlDto
  })
  @ApiResponse({ description: 'Forbidden', status: 403 })
  @Post('write')
  async getWriteUrl(
    @Param('orgname') orgname: string,
    @Body() pathDto: PathDto
  ) {
    if (pathDto.isDir) {
      await this.storageService.createDirectory(orgname, pathDto.path)
      return { write: '' } as WriteUrlDto
    }
    const write = await this.storageService.getSignedUrl(
      orgname,
      pathDto.path,
      'write'
    )
    return { write } as WriteUrlDto
  }

  @ApiOperation({
    description:
      "This endpoint will return a list of files and directories in the organization's secure storage at the specified path. ADMIN ONLY.",
    summary: 'Show directory'
  })
  @ApiResponse({ description: 'Not Found', status: 404 })
  @ApiResponse({ description: 'Unauthorized', status: 401 })
  @ApiResponse({
    description: 'Path was successfully retrieved',
    status: 200,
    type: [StorageItemDto]
  })
  @ApiResponse({ description: 'Forbidden', status: 403 })
  @Get('items')
  async listDirectory(
    @Param('orgname') orgname: string,
    @Query('path') path: string
  ) {
    const files = await this.storageService.listDirectory(orgname, path)
    return files
  }
}
