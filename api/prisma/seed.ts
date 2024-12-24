import { faker } from '@faker-js/faker'
import { NestFactory } from '@nestjs/core'
import { PrismaClient } from '@prisma/client'
import * as bcrypt from 'bcryptjs'

import { AppModule } from '../src/app.module'
import { UsersService } from '../src/users/users.service'
import { Logger } from '@nestjs/common'

const prisma = new PrismaClient()

export const resetDatabase = async () => {
  await prisma.$transaction([
    prisma.user.deleteMany(),
    prisma.organization.deleteMany(),
    prisma.apiToken.deleteMany(),
    prisma.authProvider.deleteMany(),
    prisma.member.deleteMany(),
    prisma.user.deleteMany(),
    prisma.label.deleteMany(),
    prisma.content.deleteMany(),
    prisma.aRToken.deleteMany(),
    prisma.pipeline.deleteMany(),
    prisma.pipelineStep.deleteMany(),
    prisma.run.deleteMany(),
    prisma.tool.deleteMany()
    // Add more tables as needed
  ])
}

async function main() {
  const logger = new Logger('SeedProcess')

  await resetDatabase()

  const app = await NestFactory.createApplicationContext(AppModule, {
    logger
  })
  const usersService = app.get(UsersService)
  const hashedPassword = await bcrypt.hash('password', 10)
  const user = await usersService.create({
    email: 'user@example.com',
    emailVerified: true,
    firstName: 'Jonathan',
    lastName: 'King',
    password: hashedPassword,
    photoUrl:
      'https://nsabers.com/cdn/shop/articles/bebec223da75d29d8e03027fd2882262.png?v=1708781179',
    username: 'user'
  })

  const roles = ['USER', 'ADMIN']
  const labels = ['work', 'personal', 'school stuff']

  await prisma.organization.update({
    data: {
      credits: 1000000,
      plan: 'UNLIMITED'
    },
    where: {
      orgname: user.defaultOrgname
    }
  })

  // Create labels
  for (let i = 0; i < 3; i++) {
    await prisma.label.create({
      data: {
        name: labels[i]!,
        orgname: user.defaultOrgname
      }
    })
  }

  // Create a bunch of content
  for (let i = 0; i < 100; i++) {
    const fakeDate = faker.date.past({ years: 1 })
    await prisma.content.create({
      data: {
        createdAt: fakeDate,
        credits: faker.number.int(10000),
        description: faker.lorem.paragraphs(2),
        labels: {
          connect: {
            name_orgname: {
              name: faker.helpers.arrayElement(labels),
              orgname: user.defaultOrgname
            }
          }
        },
        mimeType: 'application/pdf',
        name: faker.commerce.productName(),
        orgname: user.defaultOrgname,
        previewImage: 'https://picsum.photos/200/300',
        url: 'https://s26.q4cdn.com/900411403/files/doc_downloads/test.pdf'
      }
    })
  }

  // Create some API tokens
  for (let i = 0; i < 10; i++) {
    await prisma.apiToken.create({
      data: {
        key: '*******-2131',
        name: faker.commerce.productName(),
        organization: {
          connect: {
            orgname: user.defaultOrgname
          }
        },
        role: faker.helpers.arrayElement(roles) as any,
        user: {
          connect: {
            id: user.id
          }
        }
      }
    })
  }

  await app.close()

  logger.log('Seeding complete')
}

;(async () => {
  await main()
})()
