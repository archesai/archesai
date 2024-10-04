import { PaginatedDto } from "./paginated.dto";

export interface BaseController<T, CreateDto, QueryDto, UpdateDto> {
  create?(orgname: string, createDto: CreateDto): Promise<T>;
  findAll?(orgname: string, queryDto: QueryDto): Promise<PaginatedDto<T>>;
  findOne?(orgname: string, id: string): Promise<null | T>;
  remove?(orgname: string, id: string): Promise<void>;
  update?(orgname: string, id: string, updateDto: UpdateDto): Promise<T>;
}
