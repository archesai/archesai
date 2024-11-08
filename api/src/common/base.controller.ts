import { PaginatedDto } from "./dto/paginated.dto";

export interface BaseController<T, CreateDto, QueryDto, UpdateDto> {
  create?(
    orgname: string,
    createDto: CreateDto,
    ...additionalParams: any[]
  ): Promise<T>;
  findAll?(
    orgname: string,
    queryDto: QueryDto,
    ...additionalParams: any[]
  ): Promise<PaginatedDto<T>>;
  findOne?(orgname: string, id: string): Promise<null | T>;
  remove?(orgname: string, id: string): Promise<void>;
  update?(orgname: string, id: string, updateDto: UpdateDto): Promise<T>;
}
