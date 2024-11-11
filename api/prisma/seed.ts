import { faker } from "@faker-js/faker";
import { AuthProviderType, PrismaClient } from "@prisma/client";
import * as bcrypt from "bcryptjs";

const prisma = new PrismaClient();

export const resetDatabase = async () => {
  await prisma.$transaction([
    prisma.user.deleteMany(),
    prisma.organization.deleteMany(),
    prisma.apiToken.deleteMany(),
    prisma.authProvider.deleteMany(),
    prisma.member.deleteMany(),
    prisma.user.deleteMany(),
    prisma.tool.deleteMany(),
    prisma.label.deleteMany(),
    prisma.content.deleteMany(),
    prisma.aRToken.deleteMany(),
    prisma.pipeline.deleteMany(),
    prisma.pipelineStep.deleteMany(),
    prisma.pipelineRun.deleteMany(),
    prisma.transformation.deleteMany(),
    // Add more tables as needed
  ]);
};

async function main() {
  await resetDatabase();

  const roles = ["USER", "ADMIN"];
  const labels = ["work", "personal", "school stuff"];

  // Create init user
  const email = "user@example.com";
  const hashedPassword = await bcrypt.hash("password", 10);
  const orgname =
    email.split("@")[0] + "-" + Math.random().toString(36).substring(2, 6);
  const user = await prisma.user.create({
    data: {
      authProviders: {
        create: {
          provider: AuthProviderType.LOCAL,
          providerId: email,
        },
      },
      defaultOrgname: orgname,
      email: email,
      emailVerified: true,
      firstName: "Jonathan",
      lastName: "King",
      memberships: {
        create: {
          inviteEmail: email,
          organization: {
            create: {
              billingEmail: email,
              orgname: orgname,
              plan: "UNLIMITED",
              stripeCustomerId: "cus_123",
            },
          },
          role: "ADMIN",
        },
      },
      password: hashedPassword,
      photoUrl:
        "https://nsabers.com/cdn/shop/articles/bebec223da75d29d8e03027fd2882262.png?v=1708781179",
      username: orgname,
    },
  });

  // Create a bunch of content
  for (let i = 0; i < 100; i++) {
    const fakeDate = faker.date.past({ years: 1 });
    await prisma.content.create({
      data: {
        createdAt: fakeDate,
        credits: faker.number.int(10000),
        description: faker.lorem.paragraphs(2),
        mimeType: "application/pdf",
        name: faker.commerce.productName(),
        orgname: user.defaultOrgname,
        previewImage: "https://picsum.photos/200/300",
        url: "https://s26.q4cdn.com/900411403/files/doc_downloads/test.pdf",
      },
    });
  }

  // Create labels
  for (let i = 0; i < 3; i++) {
    await prisma.label.create({
      data: {
        name: labels[i],
        orgname: user.defaultOrgname,
      },
    });
  }

  // Create some API tokens
  for (let i = 0; i < 10; i++) {
    await prisma.apiToken.create({
      data: {
        key: "*******-2131",
        name: faker.commerce.productName(),
        organization: {
          connect: {
            orgname: user.defaultOrgname,
          },
        },
        role: faker.helpers.arrayElement(roles) as any,
        user: {
          connect: {
            id: user.id,
          },
        },
      },
    });
  }

  console.log("Successfully seeded database");
}

(async () => {
  await main();
})();
