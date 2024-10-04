import {
  CanActivate,
  ExecutionContext,
  HttpException,
  Injectable,
  NotFoundException,
  UnauthorizedException,
} from "@nestjs/common";
import { Logger } from "@nestjs/common";
import { Reflector } from "@nestjs/core";

import { ApiTokensService } from "../../api-tokens/api-tokens.service";
@Injectable()
export class RestrictedAPIKeyGuard implements CanActivate {
  private readonly logger: Logger = new Logger("RestrictedAPIKeyGuard");

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

      const { domains, id, orgname } = JSON.parse(
        Buffer.from(bearerToken.split(".")[1], "base64").toString()
      );
      if (!orgname) {
        return true;
      }
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

      // Verify token is valid for this chatbot if we are looking at an chatbot
      const chatbotId = request.params.chatbotId;
      if (chatbotId) {
        // Get current token
        const tokens = await this.apiTokensService.findAll(orgname, {});
        if (!tokens || !tokens.results.find((t) => t.id == id)) {
          throw new UnauthorizedException("Token not found");
        }
        const token = tokens.results.find((t) => t.id == id);
        // admin tokens are valid for all chatbots
        if (token.role == "ADMIN") {
          return true;
        }

        const isTokenValidForChatbot = Boolean(
          token.chatbots.find((val) => val.id == chatbotId)
        );
        if (!isTokenValidForChatbot) {
          throw new NotFoundException(
            "You tried to access an chatbot that doesn't exist or that you don't have access to."
          );
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
