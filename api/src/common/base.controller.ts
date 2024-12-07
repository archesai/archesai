import { Body, Delete, Get, Param, Patch, Post, Query } from '@nestjs/common'

import { BaseService } from './base.service'
import { PaginatedDto } from './dto/paginated.dto'
import { SearchQueryDto } from './dto/search-query.dto'

export class BaseController<
  Entity,
  CreateDto,
  UpdateDto,
  Service extends BaseService<Entity, CreateDto, UpdateDto, any, any>
> {
  constructor(readonly service: Service) {}

  /**
   * Create a new entity
   * @throws {400} BadRequestException
   * @throws {401} UnauthorizedException
   * @throws {404} NotFoundException
   */
  @Post()
  async create(
    @Param('orgname') orgname: string,
    @Body() createDto: CreateDto,
    ...additionalData: any[]
  ): Promise<Entity> {
    return this.service.create(orgname, createDto, additionalData)
  }

  /**
   * Get all entities
   * @throws {401} UnauthorizedException
   * @throws {404} NotFoundException
   */
  @Get()
  async findAll(
    @Param('orgname') orgname: string,
    @Query() searchQueryDto: SearchQueryDto
  ): Promise<PaginatedDto<Entity>> {
    return this.service.findAll(orgname, searchQueryDto)
  }

  /**
   * Get a single entity
   * @throws {401} UnauthorizedException
   * @throws {404} NotFoundException
   */
  @Get(':id')
  async findOne(@Param('orgname') orgname: string, @Param('id') id: string): Promise<Entity> {
    return this.service.findOne(orgname, id)
  }

  /**
   * Delete an entity
   * @throws {401} UnauthorizedException
   * @throws {404} NotFoundException
   */
  @Delete(':id')
  async remove(@Param('orgname') orgname: string, @Param('id') id: string): Promise<void> {
    return this.service.remove(orgname, id)
  }

  /**
   * Update an entity
   * @throws {401} UnauthorizedException
   * @throws {404} NotFoundException
   */
  @Patch(':id')
  async update(
    @Param('orgname') orgname: string,
    @Param('id') id: string,
    @Body() updateDto: UpdateDto
  ): Promise<Entity> {
    return this.service.update(orgname, id, updateDto)
  }
}
