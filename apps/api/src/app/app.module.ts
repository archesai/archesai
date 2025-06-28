import { DiscoveryModule } from '@nestjs/core'

import type { DynamicModule, ModuleMetadata } from '@archesai/core'

import {
  AccessTokensModule,
  AccountsModule,
  ApiTokensModule,
  EmailChangeModule,
  EmailVerificationModule,
  InvitationsModule,
  JwtModule,
  MembersModule,
  OAuthModule,
  OrganizationsModule,
  PassportModule,
  PasswordResetModule,
  RegistrationModule,
  SessionsModule,
  UsersModule,
  VerificationTokensModule
} from '@archesai/auth'
import {
  CallbacksModule,
  CheckoutSessionsModule,
  CustomersModule,
  PaymentMethodsModule,
  PlansModule,
  PortalModule,
  StripeModule,
  SubscriptionModule
} from '@archesai/billing'
import {
  ConfigModule,
  CorsModule,
  DatabaseModule,
  DocsModule,
  EmailModule,
  EventBusModule,
  ExceptionsModule,
  FetcherModule,
  HealthModule,
  LoggingModule,
  RedisModule,
  WebsocketsModule
} from '@archesai/core'
import { createDrizzleDatabaseService } from '@archesai/database'
import {
  AudioModule,
  LlmModule,
  RunpodModule,
  ScraperModule,
  SpeechModule
} from '@archesai/intelligence'
import {
  ArtifactsModule,
  JobsModule,
  LabelsModule,
  PipelinesModule,
  RunsModule,
  ToolsModule
} from '@archesai/orchestration'
import { FilesModule, StorageModule } from '@archesai/storage'

export const AppModuleDefinition: ModuleMetadata = {
  imports: [
    // DISCOVERY MODULE
    DiscoveryModule,

    // CORE MODULES
    ConfigModule,
    CorsModule,
    DatabaseModule.forRootAsync(createDrizzleDatabaseService),
    DocsModule,
    EmailModule,
    ExceptionsModule,
    HealthModule,
    LoggingModule,
    WebsocketsModule,
    FetcherModule,
    RedisModule.forRoot(),

    // AI MODULES
    LlmModule,
    RunpodModule,
    SpeechModule,
    AudioModule,
    ScraperModule,

    // STORAGE MODULES
    StorageModule,
    FilesModule,

    // EXTRA GLOBAL MODULES
    EventBusModule,

    // ORCHESTRATION MODULES
    PipelinesModule,
    ToolsModule,
    ArtifactsModule,
    RunsModule,
    LabelsModule,
    JobsModule,

    // AUTH MODULES
    OrganizationsModule,
    UsersModule,
    MembersModule,
    InvitationsModule,

    // BILLING MODULES
    CallbacksModule,
    CheckoutSessionsModule,
    CustomersModule,
    PaymentMethodsModule,
    PlansModule,
    StripeModule,
    PortalModule,
    SubscriptionModule,

    // AUTH MODULES
    JwtModule,
    ApiTokensModule,
    OAuthModule,
    AccountsModule,
    SessionsModule,
    AccessTokensModule,
    VerificationTokensModule,
    PasswordResetModule,
    EmailChangeModule,
    EmailVerificationModule,
    PassportModule,
    RegistrationModule
  ]
}

export class AppModule {
  public static forRoot(): DynamicModule {
    return {
      module: AppModule,
      ...AppModuleDefinition
    }
  }
}
