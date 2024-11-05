import {
  CanActivate,
  ExecutionContext,
  HttpException,
  Injectable,
  UnauthorizedException,
} from "@nestjs/common";
import { Logger } from "@nestjs/common";
import { Reflector } from "@nestjs/core";

import { ApiTokensService } from "../../api-tokens/api-tokens.service";
@Injectable()
export class RestrictedAPIKeyGuard implements CanActivate {
  private readonly logger: Logger = new Logger("Restricted API Key Guard");

  constructor(
    private reflector: Reflector,
    private apiTokensService: ApiTokensService
  ) {}
  async canActivate(context: ExecutionContext) {
    const isPublic = this.reflector.getAllAndOverride<boolean>("public", [
      context.getHandler(),
      context.getClass(),
    ]);
    if (isPublic) {
      return true;
    }

    const request = context.switchToHttp().getRequest();
    try {
      const bearerToken = request.headers.authorization.split(" ")[1];

      const { domains, orgname } = JSON.parse(
        Buffer.from(bearerToken.split(".")[1], "base64").toString()
      );
      if (!orgname) {
        return true;
      }
      this.logger.log(`Validating API Key is valid for this resource`);
      // Make sure if there was a domain in their header that this request is authorized for it
      if (domains && domains != "*") {
        const url = new URL(request.headers.origin);
        const origin = url.hostname
          .replace("www.", "")
          .replace("http://", "")
          .replace("https://", "")
          .split(":")[0];
        if (!domains.split(",").includes(origin)) {
          throw new UnauthorizedException("Unauthorized domain");
        }
      }
    } catch (error) {
      if (error instanceof HttpException) {
        throw error;
      }
      const message = `ERROR PARSING REQUEST FOR DOMAINS: ${error}`;
      this.logger.error(message);
      throw new UnauthorizedException(message);
    }

    return true;
  }
}
