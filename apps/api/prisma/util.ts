import { PrismaClient } from '@prisma/client'

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
