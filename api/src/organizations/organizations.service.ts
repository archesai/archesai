import { Injectable } from "@nestjs/common";
import { ForbiddenException } from "@nestjs/common";
import { ConfigService } from "@nestjs/config";
import { Organization, PlanType } from "@prisma/client";

import { CurrentUserDto } from "../auth/decorators/current-user.decorator";
import { BillingService } from "../billing/billing.service";
import { PrismaService } from "../prisma/prisma.service";
import { CreateOrganizationDto } from "./dto/create-organization.dto";
import { UpdateOrganizationDto } from "./dto/update-organization.dto";

@Injectable()
export class OrganizationsService {
  constructor(
    private billingService: BillingService,
    private prisma: PrismaService,
    private configService: ConfigService
  ) {}

  async addCredits(orgname: string, numCredits: number) {
    return this.prisma.organization.update({
      data: { credits: { increment: Math.ceil(numCredits) } },
      where: { orgname },
    });
  }

  async checkCredits(orgname: string, numCredits: number) {
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
    const freeUser = this.configService.get("FEATURE_BILLING") === true;
    // If billing is enabled, create a stripe user, otherwsie set it to orgname
    let stripeCustomerId = createOrganizationDto.orgname;
    if (freeUser) {
      const stripeCustomer = await this.billingService.createCustomer(
        createOrganizationDto.orgname,
        createOrganizationDto.billingEmail
      );
      stripeCustomerId = stripeCustomer.id;
    }
    const organization = await this.prisma.organization.create({
      data: {
        ...createOrganizationDto,
        chatbots: {
          create: {
            description:
              "You are an AI called Arches AI and you are designed to help people understand documents. You always answer the user's question as briefly and concisely as possible, unless the user asks you to elaborate. If the answer was found in the text, you include a direct citation. If the answer was not found in the text, you tell the user. If the user asks something like 'what is this?', they are usually talking about the included documents as a whole. You always respond in the same language that the most recent question was asked. In your citation, do not include any file names",
            llmBase: "gpt-4",
            name: "Default",
          },
        },
        credits:
          // If this is their first org and their e-mail is verified, give them free credits
          // Otherwise, if billing is disabled, give them free credits
          freeUser
            ? user.memberships?.length == 0 && user.emailVerified
              ? 0
              : 0
            : 100000000, // if this is their first org and their e-mail is verified, give them free credits
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
        plan: freeUser ? PlanType.FREE : PlanType.API,
        stripeCustomerId: stripeCustomerId,
      },
    });

    await this.prisma.user.update({
      data: { defaultOrg: organization.orgname },
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
