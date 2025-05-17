import type {
  ArchesApiRequest,
  ArchesApiResponse
} from '#utils/get-req.transformer'

export interface CanActivate {
  canActivate(
    request: ArchesApiRequest,
    reply: ArchesApiResponse
  ): boolean | Promise<boolean>
}
