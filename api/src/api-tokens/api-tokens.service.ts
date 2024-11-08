import { Injectable } from "@nestjs/common";
import { ConfigService } from "@nestjs/config";
import { JwtService } from "@nestjs/jwt";
import { ApiToken } from "@prisma/client";
import { v4 } from "uuid";

import { BaseService } from "../common/base.service";
import { WebsocketsService } from "../websockets/websockets.service";
import { ApiTokenRepository } from "./api-token.repository";
import { CreateApiTokenDto } from "./dto/create-api-token.dto";
import { UpdateApiTokenDto } from "./dto/update-api-token.dto";
import { ApiTokenEntity } from "./entities/api-token.entity";

@Injectable()
export class ApiTokensService extends BaseService<
  ApiTokenEntity,
  CreateApiTokenDto,
  UpdateApiTokenDto,
  ApiTokenRepository,
  ApiToken
> {
  constructor(
    private apiTokenRepository: ApiTokenRepository,
    private configService: ConfigService,
    private jwtService: JwtService,
    private websocketsService: WebsocketsService
  ) {
    super(apiTokenRepository);
  }

  async create(
    orgname: string,
    createTokenDto: CreateApiTokenDto,
    additionalData: {
      uid: string;
    }
  ) {
    const id = v4();
    const token = this.jwtService.sign(
      {
        domains: createTokenDto.domains,
        id,
        orgname,
        role: createTokenDto.role,
        uid: additionalData.uid,
      },
      {
        expiresIn: `${this.configService.get(
          "JWT_API_TOKEN_EXPIRATION_TIME"
        )}s`,
        secret: this.configService.get("JWT_API_TOKEN_SECRET"),
      }
    );
    const snippet = "*********" + token.slice(-5);
    const apiToken = await this.apiTokenRepository.create(
      orgname,
      createTokenDto,
      {
        id,
        snippet,
        uid: additionalData.uid,
      }
    );
    this.websocketsService.socket.to(orgname).emit("update", {
      queryKey: ["organizations", orgname, "api-tokens"],
    });

    apiToken.key = token;

    return this.toEntity(apiToken);
  }

  protected toEntity(model: ApiToken): ApiTokenEntity {
    return new ApiTokenEntity(model);
  }
}
