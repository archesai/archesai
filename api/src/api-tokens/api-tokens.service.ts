import { Injectable } from "@nestjs/common";
import { ConfigService } from "@nestjs/config";
import { JwtService } from "@nestjs/jwt";
import { v4 } from "uuid";

import { BaseService } from "../common/base.service";
import { WebsocketsService } from "../websockets/websockets.service";
import { ApiTokenRepository } from "./api-token.repository";
import { CreateApiTokenDto } from "./dto/create-api-token.dto";
import { UpdateApiTokenDto } from "./dto/update-api-token.dto";
import { ApiTokenEntity, ApiTokenModel } from "./entities/api-token.entity";

@Injectable()
export class ApiTokensService extends BaseService<
  ApiTokenEntity,
  CreateApiTokenDto,
  UpdateApiTokenDto,
  ApiTokenRepository,
  ApiTokenModel
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
      username: string;
    }
  ) {
    const id = v4();
    const token = this.jwtService.sign(
      {
        domains: createTokenDto.domains,
        id,
        orgname,
        role: createTokenDto.role,
        username: additionalData.username,
      },
      {
        expiresIn: `${this.configService.get(
          "JWT_API_TOKEN_EXPIRATION_TIME"
        )}s`,
        secret: this.configService.get("JWT_API_TOKEN_SECRET"),
      }
    );
    const key = "*********" + token.slice(-5);
    const apiToken = await this.apiTokenRepository.create(
      orgname,
      createTokenDto,
      {
        id,
        key,
        username: additionalData.username,
      }
    );
    this.websocketsService.socket.to(orgname).emit("update", {
      queryKey: ["organizations", orgname, "api-tokens"],
    });

    return this.toEntity({
      ...apiToken,
      key: token,
    });
  }

  protected toEntity(model: ApiTokenModel): ApiTokenEntity {
    return new ApiTokenEntity(model);
  }
}
