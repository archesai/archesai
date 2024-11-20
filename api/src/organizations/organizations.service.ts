import { forwardRef, Inject, Injectable, Logger } from "@nestjs/common";
import { ForbiddenException } from "@nestjs/common";
import { ConfigService } from "@nestjs/config";
import { PlanType } from "@prisma/client";

import { BillingService } from "../billing/billing.service";
import { BaseService } from "../common/base.service";
import { PipelinesService } from "../pipelines/pipelines.service";
import { ToolsService } from "../tools/tools.service";
import { UserEntity } from "../users/entities/user.entity";
import { WebsocketsService } from "../websockets/websockets.service";
import { CreateOrganizationDto } from "./dto/create-organization.dto";
import { UpdateOrganizationDto } from "./dto/update-organization.dto";
import {
  OrganizationEntity,
  OrganizationModel,
} from "./entities/organization.entity";
import { OrganizationRepository } from "./organization.repository";

@Injectable()
export class OrganizationsService extends BaseService<
  OrganizationEntity,
  CreateOrganizationDto,
  UpdateOrganizationDto,
  OrganizationRepository,
  OrganizationModel
> {
  private readonly logger = new Logger(OrganizationsService.name);
  constructor(
    @Inject(forwardRef(() => BillingService))
    private billingService: BillingService,
    private configService: ConfigService,
    private organizationRepository: OrganizationRepository,
    private toolsService: ToolsService,
    private pipelinesService: PipelinesService,
    private websocketsService: WebsocketsService
  ) {
    super(organizationRepository);
  }

  async addOrRemoveCredits(orgname: string, numCredits: number) {
    this.logger.log(`Adding ${numCredits} credits to ${orgname}`);
    const organization = await this.findByOrgname(orgname);
    const organizationEntity = await this.organizationRepository.updateRaw(
      orgname,
      organization.id,
      {
        credits:
          numCredits < 0
            ? { decrement: -numCredits }
            : { increment: numCredits },
      }
    );
    this.emitMutationEvent(orgname);
    return this.toEntity(organizationEntity);
  }

  async checkCredits(orgname: string, numCredits: number) {
    this.logger.log(`Checking ${numCredits} credits for ${orgname}`);
    const organization = await this.findByOrgname(orgname);
    if (organization.plan != "PREMIUM" && organization.credits <= numCredits) {
      throw new ForbiddenException(
        "Sorry, you do not have enough credits. Please purchase more credits to continue" +
          (organization.credits < numCredits
            ? ` (estimated cost: ${numCredits})`
            : "")
      );
    }
  }

  async create(
    orgname: string,
    createOrganizationDto: CreateOrganizationDto,
    user: UserEntity
  ) {
    this.logger.log(
      `Creating organization for user ${user.username}: ${JSON.stringify(
        createOrganizationDto,
        null,
        2
      )}`
    );
    // If billing is enabled, create a stripe user, otherwsie set it to orgname
    const billingEnabled = this.configService.get("FEATURE_BILLING") === true;
    let stripeCustomerId = createOrganizationDto.orgname;
    if (billingEnabled) {
      this.logger.log("BILLING ENABLED - Creating stripe customer");
      const stripeCustomer = await this.billingService.createCustomer(
        createOrganizationDto.orgname,
        createOrganizationDto.billingEmail
      );
      stripeCustomerId = stripeCustomer.id;
    }

    // Create organization and tools
    const organization = await this.organizationRepository.create(
      null,
      createOrganizationDto,
      {
        billingEnabled,
        stripeCustomerId,
        user,
      }
    );

    await this.toolsService.createDefaultTools(organization.orgname);
    await this.pipelinesService.createDefaultPipeline(organization.orgname);

    return this.toEntity(organization);
  }

  async findByOrgname(orgname: string) {
    return this.toEntity(
      await this.organizationRepository.findByOrgname(orgname)
    );
  }

  async findByStripeCustomerId(stripeCustomerId: string) {
    return this.toEntity(
      await this.organizationRepository.findByStripeCustomerId(stripeCustomerId)
    );
  }

  async setPlan(orgname: string, plan: PlanType) {
    const organization = await this.findByOrgname(orgname);
    const organizationEntity = this.toEntity(
      await this.organizationRepository.updateRaw(orgname, organization.id, {
        plan,
      })
    );
    this.emitMutationEvent(orgname);
    return organizationEntity;
  }

  protected emitMutationEvent(orgname: string): void {
    this.websocketsService.socket.to(orgname).emit("update", {
      queryKey: ["organizations", orgname],
    });
  }

  protected toEntity(model: OrganizationModel): OrganizationEntity {
    return new OrganizationEntity(model);
  }
}
