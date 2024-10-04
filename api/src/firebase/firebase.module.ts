import { Module } from "@nestjs/common";
import { ConfigModule } from "@nestjs/config";

import { FirebaseService } from "./firebase.service";

@Module({
  exports: [FirebaseService],
  imports: [ConfigModule],
  providers: [FirebaseService],
})
export class FirebaseModule {}
