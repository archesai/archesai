/* eslint-disable */
export default async () => {
  const t = {
    ["./common/dto/search-query.dto"]: await import(
      "./common/dto/search-query.dto"
    ),
    ["./users/entities/auth-provider.entity"]: await import(
      "./users/entities/auth-provider.entity"
    ),
    ["./members/entities/member.entity"]: await import(
      "./members/entities/member.entity"
    ),
    ["./common/dto/paginated.dto"]: await import("./common/dto/paginated.dto"),
    ["./common/entities/base-sub-item.entity"]: await import(
      "./common/entities/base-sub-item.entity"
    ),
    ["./tools/entities/tool.entity"]: await import(
      "./tools/entities/tool.entity"
    ),
    ["./pipelines/entities/pipeline-step.entity"]: await import(
      "./pipelines/entities/pipeline-step.entity"
    ),
    ["./pipelines/dto/create-pipeline.dto"]: await import(
      "./pipelines/dto/create-pipeline.dto"
    ),
    ["./billing/entities/payment-method.entity"]: await import(
      "./billing/entities/payment-method.entity"
    ),
    ["./llm/dto/create-chat-completion.dto"]: await import(
      "./llm/dto/create-chat-completion.dto"
    ),
    ["./api-tokens/entities/api-token.entity"]: await import(
      "./api-tokens/entities/api-token.entity"
    ),
    ["./storage/dto/read-url.dto"]: await import("./storage/dto/read-url.dto"),
    ["./storage/dto/write-url.dto"]: await import(
      "./storage/dto/write-url.dto"
    ),
    ["./storage/dto/storage-item.dto"]: await import(
      "./storage/dto/storage-item.dto"
    ),
    ["./billing/entities/billing-url.entity"]: await import(
      "./billing/entities/billing-url.entity"
    ),
    ["./billing/entities/plan.entity"]: await import(
      "./billing/entities/plan.entity"
    ),
    ["./organizations/entities/organization.entity"]: await import(
      "./organizations/entities/organization.entity"
    ),
    ["./users/entities/user.entity"]: await import(
      "./users/entities/user.entity"
    ),
    ["./auth/dto/token.dto"]: await import("./auth/dto/token.dto"),
  };
  return {
    "@nestjs/swagger": {
      models: [
        [
          import("./common/dto/search-query.dto"),
          {
            AggregateFieldQuery: {
              field: { required: true, type: () => String },
              granularity: {
                required: false,
                description: "The granularity to use for ranged aggregates",
                example: "day",
                enum: t["./common/dto/search-query.dto"].Granularity,
              },
              type: {
                required: true,
                type: () => Object,
                description: "The type of aggregate to perform",
                example: "count",
              },
            },
            FieldFieldQuery: {
              field: { required: true, type: () => String },
              operator: {
                required: false,
                enum: t["./common/dto/search-query.dto"].Operator,
              },
              value: { required: true, type: () => Object },
            },
            SearchQueryDto: {
              aggregates: {
                required: false,
                type: () => [
                  t["./common/dto/search-query.dto"].AggregateFieldQuery,
                ],
              },
              endDate: {
                required: false,
                type: () => String,
                description: "The end date to search to",
                example: "2022-01-01",
              },
              filters: {
                required: false,
                type: () => [
                  t["./common/dto/search-query.dto"].FieldFieldQuery,
                ],
              },
              limit: {
                required: false,
                type: () => Number,
                default: 10,
                minimum: 1,
              },
              offset: {
                required: false,
                type: () => Number,
                description: "The offset of the returned results",
                example: 10,
                default: 0,
              },
              sortBy: {
                required: false,
                type: () => String,
                description: "The field to sort the results by",
                example: "createdAt",
                default: "createdAt",
              },
              sortDirection: {
                required: false,
                description: "The direction to sort the results by",
                example: "desc",
                enum: t["./common/dto/search-query.dto"].SortDirection,
              },
              startDate: {
                required: false,
                type: () => Date,
                description: "The start date to search from",
                example: "2021-01-01",
              },
            },
          },
        ],
        [
          import("./common/entities/base.entity"),
          {
            BaseEntity: {
              createdAt: {
                required: true,
                type: () => Date,
                description: "The date that this item was created",
                example: "2023-07-11T21:09:20.895Z",
              },
              id: {
                required: true,
                type: () => String,
                description: "The ID of the item",
                example: "item-id",
              },
            },
          },
        ],
        [
          import("./api-tokens/entities/api-token.entity"),
          {
            ApiTokenEntity: {
              domains: {
                required: true,
                type: () => String,
                description: "The domains that can access this API token",
                example: "archesai.com,localhost:3000",
                default: "*",
              },
              key: {
                required: true,
                type: () => String,
                description: "The API token key. This will only be shown once",
                example: "********1234567890",
              },
              name: {
                required: true,
                type: () => String,
                description: "The name of the API token",
              },
              orgname: {
                required: true,
                type: () => String,
                description: "The organization name",
                example: "my-organization",
              },
              role: {
                required: true,
                type: () => Object,
                description: "The role of the API token",
                example: "ADMIN",
              },
              username: {
                required: true,
                type: () => String,
                description: "The username of the user who owns this API token",
                example: "jonathan",
              },
            },
          },
        ],
        [
          import("./api-tokens/dto/create-api-token.dto"),
          { CreateApiTokenDto: {} },
        ],
        [
          import("./api-tokens/dto/update-api-token.dto"),
          { UpdateApiTokenDto: {} },
        ],
        [
          import("./members/entities/member.entity"),
          {
            MemberEntity: {
              inviteAccepted: { required: true, type: () => Boolean },
              inviteEmail: { required: true, type: () => String },
              orgname: { required: true, type: () => String },
              role: { required: true, type: () => Object },
              username: { required: true, type: () => String, nullable: true },
            },
          },
        ],
        [
          import("./users/entities/auth-provider.entity"),
          {
            AuthProviderEntity: {
              provider: {
                required: true,
                type: () => Object,
                description: "The auth provider's provider",
                example: "LOCAL",
              },
              providerId: { required: true, type: () => String },
              userId: { required: true, type: () => String },
            },
          },
        ],
        [
          import("./users/entities/user.entity"),
          {
            UserEntity: {
              authProviders: {
                required: true,
                type: () => [
                  t["./users/entities/auth-provider.entity"].AuthProviderEntity,
                ],
              },
              deactivated: { required: true, type: () => Boolean },
              defaultOrgname: { required: true, type: () => String },
              displayName: { required: true, type: () => String },
              email: { required: true, type: () => String, format: "email" },
              emailVerified: { required: true, type: () => Boolean },
              firstName: { required: true, type: () => String },
              lastName: { required: true, type: () => String },
              memberships: {
                required: true,
                type: () => [
                  t["./members/entities/member.entity"].MemberEntity,
                ],
              },
              photoUrl: { required: true, type: () => String },
              username: { required: true, type: () => String, minLength: 5 },
            },
          },
        ],
        [
          import("./common/dto/paginated.dto"),
          {
            AggregateFieldResult: {
              value: { required: true, type: () => Number },
            },
            Metadata: {
              limit: { required: true, type: () => Number },
              offset: { required: true, type: () => Number },
              totalResults: { required: true, type: () => Number },
            },
            PaginatedDto: {
              aggregates: {
                required: true,
                type: () => [
                  t["./common/dto/paginated.dto"].AggregateFieldResult,
                ],
              },
              metadata: {
                required: true,
                type: () => t["./common/dto/paginated.dto"].Metadata,
              },
              results: { required: true },
            },
          },
        ],
        [
          import("./storage/dto/path.dto"),
          {
            PathDto: {
              isDir: { required: false, type: () => Boolean, default: false },
              path: { required: true, type: () => String },
            },
          },
        ],
        [
          import("./storage/dto/read-url.dto"),
          { ReadUrlDto: { read: { required: true, type: () => String } } },
        ],
        [
          import("./storage/dto/storage-item.dto"),
          {
            StorageItemDto: {
              createdAt: { required: true, type: () => Date },
              id: { required: true, type: () => String },
              isDir: { required: true, type: () => Boolean },
              name: { required: true, type: () => String },
              size: { required: true, type: () => Number },
            },
          },
        ],
        [
          import("./storage/dto/write-url.dto"),
          { WriteUrlDto: { write: { required: true, type: () => String } } },
        ],
        [
          import("./common/entities/base-sub-item.entity"),
          {
            SubItemEntity: {
              id: { required: true, type: () => String },
              name: { required: true, type: () => String },
            },
          },
        ],
        [
          import("./tools/entities/tool.entity"),
          {
            ToolEntity: {
              description: {
                required: true,
                type: () => String,
                description: "The tool description",
              },
              inputType: { required: true, type: () => Object },
              name: { required: true, type: () => String },
              orgname: { required: true, type: () => String },
              outputType: { required: true, type: () => Object },
              toolBase: { required: true, type: () => String },
            },
          },
        ],
        [
          import("./pipelines/entities/pipeline-step.entity"),
          {
            PipelineStepEntity: {
              dependents: {
                required: true,
                type: () => [
                  t["./common/entities/base-sub-item.entity"].SubItemEntity,
                ],
                description: "The order of the step in the pipeline",
              },
              dependsOn: {
                required: true,
                type: () => [
                  t["./common/entities/base-sub-item.entity"].SubItemEntity,
                ],
                description: "These are the steps that this step depends on.",
              },
              name: {
                required: true,
                type: () => String,
                description:
                  "The name of the step in the pipeline. It must be unique within the pipeline.",
              },
              pipelineId: {
                required: true,
                type: () => String,
                description: "The ID of the pipelin that this step belongs to",
                example: "pipeline-id",
              },
              tool: {
                required: true,
                type: () => t["./tools/entities/tool.entity"].ToolEntity,
                description: "The name of the tool that this step uses.",
              },
              toolId: {
                required: true,
                type: () => String,
                description: "This is the ID of the tool that this step uses.",
                example: "tool-id",
              },
            },
          },
        ],
        [
          import("./pipelines/entities/pipeline.entity"),
          {
            PipelineEntity: {
              description: {
                required: true,
                type: () => String,
                nullable: true,
                description: "The description of the pipeline",
                example: "This pipeline does something",
              },
              name: {
                required: true,
                type: () => String,
                description: "The name of the pipeline",
                example: "my-pipeline",
              },
              orgname: {
                required: true,
                type: () => String,
                description:
                  "The name of the organization that this pipeline belongs to",
                example: "my-org",
              },
              pipelineSteps: {
                required: true,
                type: () => [
                  t["./pipelines/entities/pipeline-step.entity"]
                    .PipelineStepEntity,
                ],
                description: "The steps in the pipeline",
              },
            },
          },
        ],
        [
          import("./pipelines/dto/create-pipeline.dto"),
          {
            CreatePipelineDto: {
              pipelineSteps: {
                required: true,
                type: () => [
                  t["./pipelines/dto/create-pipeline.dto"]
                    .CreatePipelineStepDto,
                ],
                description:
                  "An array of pipeline tools to be added to the pipeline",
              },
            },
            CreatePipelineStepDto: {
              dependsOn: {
                required: true,
                type: () => [String],
                description: "An array of steps that this step depends on",
                example: ["step-id", "step-id-2"],
              },
            },
          },
        ],
        [
          import("./pipelines/dto/update-pipeline.dto"),
          { UpdatePipelineDto: {} },
        ],
        [import("./tools/dto/create-tool.dto"), { CreateToolDto: {} }],
        [import("./tools/dto/update-tool.dto"), { UpdateToolDto: {} }],
        [
          import("./organizations/entities/organization.entity"),
          {
            OrganizationEntity: {
              billingEmail: {
                required: true,
                type: () => String,
                format: "email",
              },
              credits: { required: true, type: () => Number },
              orgname: { required: true, type: () => String },
              plan: { required: true, type: () => Object },
            },
          },
        ],
        [
          import("./organizations/dto/create-organization.dto"),
          { CreateOrganizationDto: {} },
        ],
        [
          import("./organizations/dto/update-organization.dto"),
          { UpdateOrganizationDto: {} },
        ],
        [
          import("./billing/entities/billing-url.entity"),
          { BillingUrlEntity: { url: { required: true, type: () => String } } },
        ],
        [
          import("./billing/entities/payment-method.entity"),
          {
            Address: {
              city: { required: true, type: () => String, nullable: true },
              country: { required: true, type: () => String, nullable: true },
              line1: { required: true, type: () => String, nullable: true },
              line2: { required: true, type: () => String, nullable: true },
              postal_code: {
                required: true,
                type: () => String,
                nullable: true,
              },
              state: { required: true, type: () => String, nullable: true },
            },
            BillingDetails: {
              address: {
                required: true,
                type: () =>
                  t["./billing/entities/payment-method.entity"].Address,
              },
              email: { required: true, type: () => String, nullable: true },
              name: { required: true, type: () => String, nullable: true },
              phone: { required: true, type: () => String, nullable: true },
            },
            CardDetails: {
              brand: { required: true, type: () => String },
              country: { required: true, type: () => String },
              exp_month: { required: true, type: () => Number },
              exp_year: { required: true, type: () => Number },
              fingerprint: {
                required: true,
                type: () => String,
                nullable: true,
              },
              funding: { required: true, type: () => String },
              last4: { required: true, type: () => String },
            },
            PaymentMethodEntity: {
              billing_details: {
                required: true,
                type: () =>
                  t["./billing/entities/payment-method.entity"].BillingDetails,
              },
              card: {
                required: true,
                type: () =>
                  t["./billing/entities/payment-method.entity"].CardDetails,
                nullable: true,
              },
              customer: { required: true, type: () => String, nullable: true },
              id: { required: true, type: () => String },
              type: { required: true, type: () => String },
            },
          },
        ],
        [
          import("./billing/entities/plan.entity"),
          {
            PlanEntity: {
              currency: { required: true, type: () => String },
              description: {
                required: true,
                type: () => String,
                nullable: true,
              },
              id: { required: true, type: () => String },
              metadata: { required: true, type: () => Object },
              name: { required: true, type: () => String },
              priceId: { required: true, type: () => String },
              priceMetadata: { required: true, type: () => Object },
              recurring: { required: true, type: () => Object, nullable: true },
              unitAmount: { required: true, type: () => Number },
            },
          },
        ],
        [
          import("./users/dto/create-user.dto"),
          {
            CreateUserDto: {
              email: { required: true, type: () => String },
              emailVerified: { required: true, type: () => Boolean },
              firstName: { required: false, type: () => String },
              lastName: { required: false, type: () => String },
              password: { required: false, type: () => String },
              photoUrl: { required: true, type: () => String },
              username: { required: true, type: () => String },
            },
          },
        ],
        [import("./users/dto/update-user.dto"), { UpdateUserDto: {} }],
        [
          import("./auth/dto/confirmation-token.dto"),
          {
            ConfirmationTokenDto: {
              token: { required: true, type: () => String },
            },
          },
        ],
        [
          import("./auth/dto/confirmation-token-with-new-password.dto"),
          {
            ConfirmationTokenWithNewPasswordDto: {
              newPassword: { required: true, type: () => String },
            },
          },
        ],
        [
          import("./auth/dto/email-request.dto"),
          {
            EmailRequestDto: {
              email: { required: true, type: () => String, format: "email" },
            },
          },
        ],
        [
          import("./auth/dto/register.dto"),
          {
            RegisterDto: {
              password: { required: true, type: () => String, minLength: 7 },
            },
          },
        ],
        [import("./auth/dto/login.dto"), { LoginDto: {} }],
        [
          import("./auth/dto/token.dto"),
          {
            TokenDto: {
              accessToken: { required: true, type: () => String },
              refreshToken: { required: true, type: () => String },
            },
          },
        ],
        [
          import("./content/entities/content.entity"),
          {
            ContentEntity: {
              children: {
                required: true,
                type: () => [
                  t["./common/entities/base-sub-item.entity"].SubItemEntity,
                ],
              },
              consumedBy: {
                required: true,
                type: () => [
                  t["./common/entities/base-sub-item.entity"].SubItemEntity,
                ],
              },
              credits: { required: true, type: () => Number },
              description: {
                required: true,
                type: () => String,
                nullable: true,
              },
              labels: {
                required: true,
                type: () => [
                  t["./common/entities/base-sub-item.entity"].SubItemEntity,
                ],
              },
              mimeType: { required: true, type: () => String, nullable: true },
              name: { required: true, type: () => String },
              orgname: { required: true, type: () => String },
              parent: {
                required: true,
                type: () =>
                  t["./common/entities/base-sub-item.entity"].SubItemEntity,
                nullable: true,
              },
              parentId: { required: true, type: () => String, nullable: true },
              previewImage: {
                required: true,
                type: () => String,
                nullable: true,
              },
              producedBy: {
                required: true,
                type: () =>
                  t["./common/entities/base-sub-item.entity"].SubItemEntity,
                nullable: true,
              },
              producedById: {
                required: true,
                type: () => String,
                nullable: true,
              },
              text: { required: true, type: () => String, nullable: true },
              url: { required: true, type: () => String, nullable: true },
            },
          },
        ],
        [
          import("./content/dto/create-content.dto"),
          {
            CreateContentDto: {
              labels: { required: false, type: () => [String] },
            },
          },
        ],
        [import("./content/dto/update-content.dto"), { UpdateContentDto: {} }],
        [
          import("./labels/entities/label.entity"),
          {
            LabelEntity: {
              name: { required: true, type: () => String },
              orgname: { required: true, type: () => String },
            },
          },
        ],
        [import("./labels/dto/create-label.dto"), { CreateLabelDto: {} }],
        [import("./labels/dto/update-label.dto"), { UpdateLabelDto: {} }],
        [
          import("./llm/dto/create-chat-completion.dto"),
          {
            CreateChatCompletionDto: {
              best_of: { required: false, type: () => Number },
              frequency_penalty: { required: false, type: () => Number },
              ignore_eos: { required: false, type: () => Boolean },
              logit_bias: { required: false, type: () => Object },
              max_tokens: { required: false, type: () => Number },
              messages: {
                required: true,
                type: () => [
                  t["./llm/dto/create-chat-completion.dto"].MessageDto,
                ],
              },
              model: { required: false, type: () => String },
              n: { required: false, type: () => Number, default: 1 },
              name: { required: false, type: () => String },
              presence_penalty: { required: false, type: () => Number },
              skip_special_tokens: { required: false, type: () => Boolean },
              stop: { required: false, type: () => [String] },
              stop_token_ids: { required: false, type: () => [Number] },
              stream: { required: false, type: () => Boolean },
              temperature: {
                required: false,
                type: () => Number,
                default: 0.7,
              },
              top_k: { required: false, type: () => Number, default: -1 },
              top_p: { required: false, type: () => Number, default: 1 },
              use_beam_search: { required: false, type: () => Boolean },
              user: { required: false, type: () => String },
            },
            MessageDto: {
              content: { required: true, type: () => String },
              name: { required: false, type: () => String },
              role: { required: true, type: () => Object },
            },
          },
        ],
        [import("./members/dto/create-member.dto"), { CreateMemberDto: {} }],
        [import("./members/dto/update-member.dto"), { UpdateMemberDto: {} }],
        [
          import("./runs/entities/run.entity"),
          {
            RunEntity: {
              completedAt: {
                required: true,
                type: () => Date,
                nullable: true,
                description: "The timestamp when the run completed",
                example: "2024-11-05T11:42:02.258Z",
              },
              error: { required: true, type: () => String, nullable: true },
              inputs: {
                required: true,
                type: () => [
                  t["./common/entities/base-sub-item.entity"].SubItemEntity,
                ],
              },
              name: { required: true, type: () => String, nullable: true },
              outputs: {
                required: true,
                type: () => [
                  t["./common/entities/base-sub-item.entity"].SubItemEntity,
                ],
              },
              pipelineId: {
                required: true,
                type: () => String,
                nullable: true,
              },
              progress: { required: true, type: () => Number },
              runType: {
                required: true,
                type: () => Object,
                description:
                  "The type of run, either an individual tool run or a pipeline run",
                example: "TOOL_RUN",
              },
              startedAt: {
                required: true,
                type: () => Date,
                nullable: true,
                description: "The timestamp when the run started",
                example: "2024-11-05T11:42:02.258Z",
              },
              status: {
                required: true,
                type: () => Object,
                description: "The status of the run",
                example: "QUEUED",
              },
              toolId: { required: true, type: () => String, nullable: true },
            },
          },
        ],
        [
          import("./runs/dto/create-run.dto"),
          {
            CreateRunDto: {
              contentIds: { required: false, type: () => [String] },
              text: { required: false, type: () => String },
              url: { required: false, type: () => String },
            },
          },
        ],
        [
          import("./llm/dto/create-completion.dto"),
          {
            CreateCompletionDto: {
              best_of: { required: false, type: () => Number },
              echo: { required: false, type: () => Boolean, default: false },
              frequency_penalty: { required: false, type: () => Number },
              ignore_eos: { required: false, type: () => Boolean },
              logit_bias: { required: false, type: () => Object },
              logprobs: { required: false, type: () => Number },
              max_tokens: { required: false, type: () => Number, default: 16 },
              model: { required: true, type: () => String },
              n: { required: true, type: () => Number, default: 1 },
              presence_penalty: { required: false, type: () => Number },
              prompt: { required: true, type: () => Object },
              skip_special_tokens: { required: false, type: () => Boolean },
              stop: { required: false, type: () => [String] },
              stop_token_ids: { required: false, type: () => [Number] },
              stream: { required: false, type: () => Boolean, default: false },
              suffix: { required: false, type: () => String },
              temperature: { required: false, type: () => Number, default: 1 },
              top_k: { required: true, type: () => Number },
              top_p: { required: false, type: () => Number, default: 1 },
              use_beam_search: { required: false, type: () => Boolean },
              user: { required: false, type: () => String },
            },
          },
        ],
        [
          import("./storage/dto/upload-document.dto"),
          {
            UploadDocumentDto: { file: { required: true, type: () => Object } },
          },
        ],
      ],
      controllers: [
        [
          import("./common/base.controller"),
          {
            BaseController: {
              create: { status: 404, description: "NotFoundException" },
              findAll: { status: 404, description: "NotFoundException" },
              findOne: { status: 404, description: "NotFoundException" },
              remove: { status: 404, description: "NotFoundException" },
              update: { status: 404, description: "NotFoundException" },
            },
          },
        ],
        [
          import("./api-tokens/api-tokens.controller"),
          {
            ApiTokensController: {
              create: {
                summary: "Create a new API token",
                description:
                  "This endpoint requires the user to be authenticated",
                type: t["./api-tokens/entities/api-token.entity"]
                  .ApiTokenEntity,
              },
            },
          },
        ],
        [
          import("./storage/storage.controller"),
          {
            StorageController: {
              delete: {},
              getReadUrl: { type: t["./storage/dto/read-url.dto"].ReadUrlDto },
              getWriteUrl: {
                type: t["./storage/dto/write-url.dto"].WriteUrlDto,
              },
              listDirectory: {
                type: [t["./storage/dto/storage-item.dto"].StorageItemDto],
              },
            },
          },
        ],
        [
          import("./billing/billing.controller"),
          {
            BillingController: {
              cancelSubscriptionPlan: {
                status: 403,
                description: "ForbiddenException",
              },
              changeSubscriptionPlan: {
                status: 400,
                description: "BadRequestException",
              },
              createBillingPortal: {
                status: 403,
                description: "ForbiddenException",
                type: t["./billing/entities/billing-url.entity"]
                  .BillingUrlEntity,
              },
              createCheckoutSession: {
                status: 400,
                description: "BadRequestException",
                type: t["./billing/entities/billing-url.entity"]
                  .BillingUrlEntity,
              },
              getPlans: {
                summary: "Get plans",
                description:
                  "This endpoint will return a list of available billing plans",
                type: [t["./billing/entities/plan.entity"].PlanEntity],
              },
              listPaymentMethods: {
                summary: "List payment methods",
                description:
                  "This endpoint will return a list of payment methods for an organization",
                type: [
                  t["./billing/entities/payment-method.entity"]
                    .PaymentMethodEntity,
                ],
              },
              removePaymentMethod: {
                status: 404,
                description: "NotFoundException",
              },
              stripe_handleIncomingEvents: {},
            },
          },
        ],
        [
          import("./organizations/organizations.controller"),
          {
            OrganizationsController: {
              create: {
                status: 404,
                description: "NotFoundException",
                type: t["./organizations/entities/organization.entity"]
                  .OrganizationEntity,
              },
              update: {
                status: 404,
                description: "NotFoundException",
                type: t["./organizations/entities/organization.entity"]
                  .OrganizationEntity,
              },
            },
          },
        ],
        [
          import("./users/users.controller"),
          {
            UsersController: {
              deactivate: { status: 400, description: "Bad Request." },
              findOne: {
                status: 404,
                description: "NotFoundException",
                type: t["./users/entities/user.entity"].UserEntity,
              },
              update: {
                status: 404,
                description: "NotFoundException",
                type: t["./users/entities/user.entity"].UserEntity,
              },
            },
          },
        ],
        [
          import("./auth/auth.controller"),
          {
            AuthController: {
              emailChangeConfirm: {
                status: 400,
                description: "BadRequestException",
                type: t["./auth/dto/token.dto"].TokenDto,
              },
              emailChangeRequest: {
                summary: "Request e-mail change with a token",
                description:
                  "This endpoint will request your e-mail change with a token",
              },
              emailVerificationConfirm: {
                status: 401,
                description: "UnauthorizedException",
                type: t["./auth/dto/token.dto"].TokenDto,
              },
              emailVerificationRequest: {
                status: 403,
                description: "ForbiddenException",
              },
              login: {
                status: 400,
                description: "BadRequestException",
                type: t["./auth/dto/token.dto"].TokenDto,
              },
              logout: { status: 401, description: "UnauthorizedException" },
              passwordResetConfirm: {
                status: 401,
                description: "UnauthorizedException",
                type: t["./auth/dto/token.dto"].TokenDto,
              },
              passwordResetRequest: {
                summary: "Request password reset",
                description: "This endpoint will request a password reset link",
              },
              refreshToken: {
                status: 401,
                description: "UnauthorizedException",
                type: t["./auth/dto/token.dto"].TokenDto,
              },
              register: {
                status: 409,
                description: "ConflictException",
                type: t["./auth/dto/token.dto"].TokenDto,
              },
              zfirebaseAuthCallback: {
                type: t["./auth/dto/token.dto"].TokenDto,
              },
              ztwitterAuth: {},
              ztwitterAuthCallback: {
                type: t["./auth/dto/token.dto"].TokenDto,
              },
            },
          },
        ],
        [
          import("./members/members.controller"),
          {
            MembersController: {
              join: {
                status: 404,
                description: "Not Found",
                type: t["./members/entities/member.entity"].MemberEntity,
              },
            },
          },
        ],
      ],
    },
  };
};
