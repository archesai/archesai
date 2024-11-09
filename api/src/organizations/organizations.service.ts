import { forwardRef, Inject, Injectable, Logger } from "@nestjs/common";
import { ForbiddenException } from "@nestjs/common";
import { ConfigService } from "@nestjs/config";
import { Organization, PlanType } from "@prisma/client";
import { v4 } from "uuid";

import { CurrentUserDto } from "../auth/decorators/current-user.decorator";
import { BillingService } from "../billing/billing.service";
import { PrismaService } from "../prisma/prisma.service";
import { CreateOrganizationDto } from "./dto/create-organization.dto";
import { UpdateOrganizationDto } from "./dto/update-organization.dto";

@Injectable()
export class OrganizationsService {
  private readonly logger = new Logger(OrganizationsService.name);
  constructor(
    @Inject(forwardRef(() => BillingService))
    private billingService: BillingService,
    private prisma: PrismaService,
    private configService: ConfigService
  ) {}

  async addCredits(orgname: string, numCredits: number) {
    this.logger.log(`Adding ${numCredits} credits to ${orgname}`);
    return this.prisma.organization.update({
      data: { credits: { increment: Math.ceil(numCredits) } },
      where: { orgname },
    });
  }

  async checkCredits(orgname: string, numCredits: number) {
    this.logger.log(`Checking ${numCredits} credits for ${orgname}`);
    const organization = await this.findOneByName(orgname);
    if (organization.plan != "PREMIUM" && organization.credits <= numCredits) {
      throw new ForbiddenException(
        "Sorry, you do not have enough credits. Please purchase more credits to continue" +
          (organization.credits < numCredits
            ? ` (estimated cost: ${numCredits})`
            : "")
      );
    }
  }

  // this creates an organization, but also creates stripe account
  async createAndInitialize(
    user: CurrentUserDto,
    createOrganizationDto: CreateOrganizationDto
  ): Promise<Organization> {
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
    let organization = await this.prisma.organization.create({
      data: {
        ...createOrganizationDto,
        credits:
          // If this is their first org and their e-mail is verified, give them free credits
          // Otherwise, if billing is disabled, give them free credits
          billingEnabled
            ? user.memberships?.length == 0 && user.emailVerified
              ? 0
              : 0
            : 100000000, // if this is their first org and their e-mail is verified, give them free credits

        // Add them as an admin to this organization
        members: {
          create: {
            inviteAccepted: true,
            inviteEmail: user.email, // FIXME
            role: "ADMIN",
            user: {
              connect: {
                username: user.username,
              },
            },
          },
        },
        plan: billingEnabled ? PlanType.FREE : PlanType.UNLIMITED,
        stripeCustomerId: stripeCustomerId,
        tools: {
          createMany: {
            data: [
              {
                description:
                  "Extract text from a file. This tool supports all file types.",
                inputType: "TEXT",
                name: "Extract Text",
                outputType: "TEXT",
                toolBase: "extract-text",
              },
              {
                description: "Create an image from text.",
                inputType: "TEXT",
                name: "Text to Image",
                outputType: "IMAGE",
                toolBase: "text-to-image",
              },
              {
                description:
                  "Summarize text. This tool supports all languages.",
                inputType: "TEXT",
                name: "Summarize",
                outputType: "TEXT",
                toolBase: "summarize",
              },
              {
                description:
                  "Create embeddings from text. This tool supports all languages.",
                inputType: "TEXT",
                name: "Create Embeddings",
                outputType: "TEXT", // FIXME make this none
                toolBase: "create-embeddings",
              },
              {
                description:
                  "Convert text to speech. This tool supports all languages.",
                inputType: "TEXT",
                name: "Text to Speech",
                outputType: "AUDIO",
                toolBase: "text-to-speech",
              },
            ],
          },
        },
      },
    });

    const tools = await this.prisma.tool.findMany({
      where: { orgname: organization.orgname },
    });

    // Create default pipeline
    const initialId = v4();
    organization = await this.prisma.organization.update({
      data: {
        pipelines: {
          create: {
            description:
              "This is a default pipeline for indexing arbitrary documents. It extracts text from the document, creates an image from the text, summarizes the text, creates embeddings from the text, and converts the text to speech.",
            name: "Default",
            pipelineSteps: {
              createMany: {
                data: [
                  {
                    id: initialId,
                    toolId: tools.find((t) => t.name == "Extract Text").id,
                  },
                  {
                    dependsOnId: initialId,
                    toolId: tools.find((t) => t.name == "Text to Image").id,
                  },
                  {
                    dependsOnId: initialId,
                    toolId: tools.find((t) => t.name == "Summarize").id,
                  },
                  {
                    dependsOnId: initialId,
                    toolId: tools.find((t) => t.name == "Create Embeddings").id,
                  },
                  {
                    dependsOnId: initialId,
                    toolId: tools.find((t) => t.name == "Text to Speech").id,
                  },
                ],
              },
            },
          },
        },
      },
      where: { id: organization.id },
    });

    await this.prisma.user.update({
      data: { defaultOrgname: organization.orgname },
      where: { username: user.username },
    });

    return organization;
  }

  async deactivate(orgname: string) {
    return this.prisma.organization.delete({ where: { orgname } });
  }

  async findAll(): Promise<Organization[]> {
    return this.prisma.organization.findMany({});
  }

  async findOne(id: string): Promise<Organization> {
    return this.prisma.organization.findUniqueOrThrow({ where: { id } });
  }

  async findOneByCustomerId(customerId: string): Promise<Organization> {
    return this.prisma.organization.findUniqueOrThrow({
      where: { stripeCustomerId: customerId },
    });
  }

  async findOneByName(orgname: string): Promise<Organization> {
    return this.prisma.organization.findUniqueOrThrow({ where: { orgname } });
  }

  async remove(id: string): Promise<void> {
    await this.prisma.organization.delete({ where: { id } });
  }

  async removeByName(orgname: string) {
    await this.prisma.organization.delete({ where: { orgname } });
  }

  async removeCredits(orgname: string, numCredits: number) {
    const organization = await this.prisma.organization.update({
      data: { credits: { decrement: Math.ceil(numCredits) } },
      where: { orgname },
    });
    if (organization.credits < 0) {
      await this.prisma.organization.update({
        data: { credits: 0 },
        where: { orgname },
      });
    }
  }

  async setPlan(orgname: string, plan: PlanType) {
    return this.prisma.organization.update({
      data: { plan },
      where: { orgname },
    });
  }

  async update(
    id: string,
    updateOrganizationDto: UpdateOrganizationDto
  ): Promise<Organization> {
    return this.prisma.organization.update({
      data: updateOrganizationDto,
      where: { id },
    });
  }

  async updateByName(
    orgname: string,
    updateOrganizationDto: UpdateOrganizationDto
  ): Promise<Organization> {
    const organization = await this.findOneByName(orgname);
    return this.update(organization.id, updateOrganizationDto);
  }
}
