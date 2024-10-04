import { Module } from "@nestjs/common";
import { ConfigModule, ConfigService } from "@nestjs/config";

import { PrismaModule } from "../prisma/prisma.module";
import { PrismaService } from "../prisma/prisma.service";
import { PgVectorDBService } from "./pgvector.vector-db.service";
import { PineconeVectorDBService } from "./pinecone.vector-db.service";
import { VectorDBService } from "./vector-db.service";
import { VECTOR_DB_SERVICE } from "./vector-db.service";

@Module({
  exports: [VECTOR_DB_SERVICE],
  imports: [ConfigModule, PrismaModule],
  providers: [
    {
      inject: [ConfigService, PrismaService],
      provide: VECTOR_DB_SERVICE,
      useFactory: (
        configService: ConfigService,
        prismaService: PrismaService
      ): VectorDBService => {
        const vectorDBType = configService.get<string>("VECTOR_DB_TYPE");
        switch (vectorDBType) {
          case "pgvector":
            return new PgVectorDBService(prismaService);
          case "pinecone":
            return new PineconeVectorDBService(configService);
          default:
            throw new Error("Invalid vector db type");
        }
      },
    },
  ],
})
export class VectorDBModule {}
