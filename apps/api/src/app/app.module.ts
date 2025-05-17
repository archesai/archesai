import { DiscoveryModule } from '@nestjs/core'

import type { DynamicModule, ModuleMetadata } from '@archesai/core'

import {
  AudioModule,
  LlmModule,
  RunpodModule,
  ScraperModule,
  SpeechModule
} from '@archesai/ai'
import {
  AccessTokensModule,
  AccountsModule,
  ApiTokensModule,
  AuthenticationModule,
  EmailChangeModule,
  EmailVerificationModule,
  JwtModule,
  OAuthModule,
  PassportModule,
  PasswordResetModule,
  SessionsModule,
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
  WebsocketsModule
} from '@archesai/core'
import { createDrizzleDatabaseService } from '@archesai/drizzle'
import {
  InvitationsModule,
  MembersModule,
  OrganizationsModule,
  UsersModule
} from '@archesai/organizations'
import { FilesModule, StorageModule } from '@archesai/storage'
import {
  ContentModule,
  JobsModule,
  LabelsModule,
  PipelinesModule,
  RunsModule,
  ToolsModule
} from '@archesai/workflows'

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

    // WORKFLOWS MODULES
    PipelinesModule,
    ToolsModule,
    ContentModule,
    RunsModule,
    LabelsModule,
    JobsModule,

    // ORGANIZATIONS MODULES
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
    AuthenticationModule,
    OAuthModule,
    AccountsModule,
    SessionsModule,
    AccessTokensModule,
    VerificationTokensModule,
    PasswordResetModule,
    EmailChangeModule,
    EmailVerificationModule,
    PassportModule
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
