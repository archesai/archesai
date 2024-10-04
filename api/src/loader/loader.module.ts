import { HttpModule } from "@nestjs/axios";
import { Module } from "@nestjs/common";
import { ConfigModule } from "@nestjs/config";

import { KMeansService } from "./k-means.service";
import { LoaderService } from "./loader.service";

@Module({
  exports: [LoaderService, KMeansService],
  imports: [ConfigModule, HttpModule],
  providers: [LoaderService, KMeansService],
})
export class LoaderModule {}
