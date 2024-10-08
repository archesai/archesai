// import { faker } from "@faker-js/faker";
import { NestFactory } from "@nestjs/core";
import { AuthProviderType } from "@prisma/client";
import * as bcrypt from "bcryptjs";

import { AppModule } from "../src/app.module";
import { CurrentUserDto } from "../src/auth/decorators/current-user.decorator";
import { OrganizationsService } from "../src/organizations/organizations.service";
// import { PrismaService } from "../src/prisma/prisma.service";
import { UsersService } from "../src/users/users.service";

// const roles = ["USER", "ADMIN"];
// const contentTypes = ["application/pdf"];

async function main() {
  const app = await NestFactory.createApplicationContext(AppModule);
  const usersService = app.get<UsersService>(UsersService);
  // const prismaService = app.get<PrismaService>(PrismaService);

  const organizationsService =
    app.get<OrganizationsService>(OrganizationsService);

  // Create init user
  let user = null as CurrentUserDto;
  const email = "user@example.com";

  const hashedPassword = await bcrypt.hash("password", 10);

  try {
    user = await usersService.create({
      email: email,
      emailVerified: true,
      firstName: "Jonathan",
      lastName: "King",
      password: hashedPassword,
      photoUrl:
        "https://nsabers.com/cdn/shop/articles/bebec223da75d29d8e03027fd2882262.png?v=1708781179",
      username:
        email.split("@")[0] + "-" + Math.random().toString(36).substring(2, 6),
    });
    user = await usersService.syncAuthProvider(
      email,
      AuthProviderType.LOCAL,
      email
    );

    // Create init organization
    await organizationsService.setPlan(user.defaultOrgname, "UNLIMITED");
    await organizationsService.addCredits(user.defaultOrgname, 1000000000);
  } catch (e) {
    console.log("User already exists");
    user = await usersService.findOneByEmail(email);
    console.log(user);
  }

  // const chatbot = await prismaService.chatbot.findFirst({
  //   where: {
  //     orgname: user.defaultOrgname,
  //   },
  // });

  try {
    // for (let i = 0; i < 100; i++) {
    //   const fakeDate = faker.date.past({ years: 1 });
    //   await prismaService.content.create({
    //     data: {
    //       buildArgs: {},
    //       createdAt: fakeDate,
    //       credits: faker.number.int(10000),
    //       description: faker.lorem.paragraphs(2),
    //       job: {
    //         create: {
    //           createdAt: fakeDate,
    //           jobType: "DOCUMENT",
    //           organization: {
    //             connect: {
    //               orgname: user.defaultOrgname,
    //             },
    //           },
    //           status: "COMPLETE",
    //         },
    //       },
    //       mimeType: "application/pdf",
    //       name: faker.commerce.productName(),
    //       organization: {
    //         connect: {
    //           orgname: user.defaultOrgname,
    //         },
    //       },
    //       previewImage: "https://picsum.photos/200/300",
    //       type: "DOCUMENT",
    //       url: "https://s26.q4cdn.com/900411403/files/doc_downloads/test.pdf",
    //     },
    //   });
    //   await prismaService.thread.create({
    //     data: {
    //       chatbotId: chatbot.id,
    //       credits: faker.number.int(10000),
    //       messages: {
    //         createMany: {
    //           data: new Array(100).fill({
    //             answer: faker.lorem.sentence(),
    //             answerLength: faker.number.int(100),
    //             contextLength: faker.number.int(100),
    //             question: faker.lorem.sentence(),
    //             topK: faker.number.int(10),
    //           }),
    //         },
    //       },
    //       name: faker.commerce.productName(),
    //       orgname: user.defaultOrgname,
    //     },
    //   });
    //   await prismaService.apiToken.create({
    //     data: {
    //       chatbots: {
    //         connect: {
    //           id: chatbot.id,
    //         },
    //       },
    //       key: "*******-2131",
    //       name: faker.commerce.productName(),
    //       organization: {
    //         connect: {
    //           orgname: user.defaultOrgname,
    //         },
    //       },
    //       role: faker.helpers.arrayElement(roles) as any,
    //       user: {
    //         connect: {
    //           id: user.id,
    //         },
    //       },
    //     },
    //   });
    // }
  } catch (e) {
    console.error("Error during data seeding: ", e);
  }

  await app.close();

  console.log("Successfully seeded database");
}

(async () => {
  await main();
})();
