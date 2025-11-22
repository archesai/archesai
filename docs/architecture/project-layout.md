# Project Layout

## Directory Structure

```text
.
├── .vscode
│   ├── extensions.json
│   └── settings.json
├── api
│   ├── components
│   │   ├── headers
│   │   │   ├── RateLimitLimit.yaml
│   │   │   ├── RateLimitRemaining.yaml
│   │   │   ├── RateLimitReset.yaml
│   │   │   ├── RetryAfter.yaml
│   │   │   └── SetCookie.yaml
│   │   ├── parameters
│   │   │   ├── APIKeysFilter.yaml
│   │   │   ├── APIKeysSort.yaml
│   │   │   ├── AccountsFilter.yaml
│   │   │   ├── AccountsSort.yaml
│   │   │   ├── ArtifactsFilter.yaml
│   │   │   ├── ArtifactsSort.yaml
│   │   │   ├── ExecutorsFilter.yaml
│   │   │   ├── ExecutorsSort.yaml
│   │   │   ├── InvitationsFilter.yaml
│   │   │   ├── InvitationsSort.yaml
│   │   │   ├── LabelsFilter.yaml
│   │   │   ├── LabelsSort.yaml
│   │   │   ├── MembersFilter.yaml
│   │   │   ├── MembersSort.yaml
│   │   │   ├── OrganizationID.yaml
│   │   │   ├── OrganizationsFilter.yaml
│   │   │   ├── OrganizationsSort.yaml
│   │   │   ├── PageQuery.yaml
│   │   │   ├── PipelinesFilter.yaml
│   │   │   ├── PipelinesSort.yaml
│   │   │   ├── ResourceID.yaml
│   │   │   ├── RunsFilter.yaml
│   │   │   ├── RunsSort.yaml
│   │   │   ├── SessionsFilter.yaml
│   │   │   ├── SessionsSort.yaml
│   │   │   ├── ToolsFilter.yaml
│   │   │   ├── ToolsSort.yaml
│   │   │   ├── UsersFilter.yaml
│   │   │   └── UsersSort.yaml
│   │   ├── responses
│   │   │   ├── APIKeyListResponse.yaml
│   │   │   ├── APIKeyResponse.yaml
│   │   │   ├── AccountListResponse.yaml
│   │   │   ├── AccountResponse.yaml
│   │   │   ├── ArtifactListResponse.yaml
│   │   │   ├── ArtifactResponse.yaml
│   │   │   ├── BadRequest.yaml
│   │   │   ├── ConfigResponse.yaml
│   │   │   ├── Conflict.yaml
│   │   │   ├── EmailVerificationResponse.yaml
│   │   │   ├── ExecutorListResponse.yaml
│   │   │   ├── ExecutorResponse.yaml
│   │   │   ├── HealthResponse.yaml
│   │   │   ├── InternalServerError.yaml
│   │   │   ├── InvitationListResponse.yaml
│   │   │   ├── InvitationResponse.yaml
│   │   │   ├── LabelListResponse.yaml
│   │   │   ├── LabelResponse.yaml
│   │   │   ├── LogoutResponse.yaml
│   │   │   ├── MemberListResponse.yaml
│   │   │   ├── MemberResponse.yaml
│   │   │   ├── NoContent.yaml
│   │   │   ├── NotFound.yaml
│   │   │   ├── OrganizationListResponse.yaml
│   │   │   ├── OrganizationResponse.yaml
│   │   │   ├── PipelineExecutionPlanResponse.yaml
│   │   │   ├── PipelineListResponse.yaml
│   │   │   ├── PipelineResponse.yaml
│   │   │   ├── PipelineStepListResponse.yaml
│   │   │   ├── PipelineStepResponse.yaml
│   │   │   ├── RunListResponse.yaml
│   │   │   ├── RunResponse.yaml
│   │   │   ├── SessionCreated.yaml
│   │   │   ├── SessionListResponse.yaml
│   │   │   ├── SessionResponse.yaml
│   │   │   ├── TooManyRequests.yaml
│   │   │   ├── ToolListResponse.yaml
│   │   │   ├── ToolResponse.yaml
│   │   │   ├── Unauthorized.yaml
│   │   │   ├── UnprocessableEntity.yaml
│   │   │   ├── UserListResponse.yaml
│   │   │   └── UserResponse.yaml
│   │   └── schemas
│   │       ├── config
│   │       │   ├── Config.yaml
│   │       │   ├── ConfigAPI.yaml
│   │       │   ├── ConfigAuth.yaml
│   │       │   ├── ConfigAuthFirebase.yaml
│   │       │   ├── ConfigAuthGithub.yaml
│   │       │   ├── ConfigAuthGoogle.yaml
│   │       │   ├── ConfigAuthLocal.yaml
│   │       │   ├── ConfigAuthMagicLink.yaml
│   │       │   ├── ConfigAuthMicrosoft.yaml
│   │       │   ├── ConfigAuthTwitter.yaml
│   │       │   ├── ConfigBilling.yaml
│   │       │   ├── ConfigDatabase.yaml
│   │       │   ├── ConfigEmail.yaml
│   │       │   ├── ConfigGrafana.yaml
│   │       │   ├── ConfigImage.yaml
│   │       │   ├── ConfigImages.yaml
│   │       │   ├── ConfigInfrastructure.yaml
│   │       │   ├── ConfigIngress.yaml
│   │       │   ├── ConfigIntelligence.yaml
│   │       │   ├── ConfigKubernetes.yaml
│   │       │   ├── ConfigLLM.yaml
│   │       │   ├── ConfigLogging.yaml
│   │       │   ├── ConfigLoki.yaml
│   │       │   ├── ConfigMigrations.yaml
│   │       │   ├── ConfigMonitoring.yaml
│   │       │   ├── ConfigPersistence.yaml
│   │       │   ├── ConfigPlatform.yaml
│   │       │   ├── ConfigRedis.yaml
│   │       │   ├── ConfigResource.yaml
│   │       │   ├── ConfigRunPod.yaml
│   │       │   ├── ConfigScraper.yaml
│   │       │   ├── ConfigServiceAccount.yaml
│   │       │   ├── ConfigSpeech.yaml
│   │       │   ├── ConfigStorage.yaml
│   │       │   ├── ConfigStripe.yaml
│   │       │   ├── ConfigTLS.yaml
│   │       │   └── ConfigUnstructured.yaml
│   │       ├── APIKey.yaml
│   │       ├── Account.yaml
│   │       ├── Artifact.yaml
│   │       ├── Base.yaml
│   │       ├── Executor.yaml
│   │       ├── FilterNode.yaml
│   │       ├── Health.yaml
│   │       ├── Invitation.yaml
│   │       ├── Label.yaml
│   │       ├── MagicLinkToken.yaml
│   │       ├── Member.yaml
│   │       ├── Organization.yaml
│   │       ├── Page.yaml
│   │       ├── PaginationMeta.yaml
│   │       ├── Pipeline.yaml
│   │       ├── PipelineStep.yaml
│   │       ├── Problem.yaml
│   │       ├── Run.yaml
│   │       ├── Session.yaml
│   │       ├── Tool.yaml
│   │       ├── UUID.yaml
│   │       └── User.yaml
│   ├── paths
│   │   ├── api-keys.yaml
│   │   ├── api-keys_id.yaml
│   │   ├── artifacts.yaml
│   │   ├── artifacts_id.yaml
│   │   ├── auth_accounts.yaml
│   │   ├── auth_accounts_id.yaml
│   │   ├── auth_change-email.yaml
│   │   ├── auth_confirm-email.yaml
│   │   ├── auth_forgot-password.yaml
│   │   ├── auth_link.yaml
│   │   ├── auth_login.yaml
│   │   ├── auth_logout-all.yaml
│   │   ├── auth_logout.yaml
│   │   ├── auth_magic-links_request.yaml
│   │   ├── auth_magic-links_verify.yaml
│   │   ├── auth_oauth_provider_authorize.yaml
│   │   ├── auth_oauth_provider_callback.yaml
│   │   ├── auth_register.yaml
│   │   ├── auth_request-verification.yaml
│   │   ├── auth_reset-password.yaml
│   │   ├── auth_sessions.yaml
│   │   ├── auth_sessions_id.yaml
│   │   ├── auth_verify-email.yaml
│   │   ├── config.yaml
│   │   ├── executors.yaml
│   │   ├── executors_id.yaml
│   │   ├── executors_id_execute.yaml
│   │   ├── health.yaml
│   │   ├── labels.yaml
│   │   ├── labels_id.yaml
│   │   ├── me.yaml
│   │   ├── organizations.yaml
│   │   ├── organizations_id.yaml
│   │   ├── organizations_organizationID_invitations.yaml
│   │   ├── organizations_organizationID_invitations_id.yaml
│   │   ├── organizations_organizationID_members.yaml
│   │   ├── organizations_organizationID_members_id.yaml
│   │   ├── pipelines.yaml
│   │   ├── pipelines_id.yaml
│   │   ├── pipelines_id_execution-plans.yaml
│   │   ├── pipelines_id_steps.yaml
│   │   ├── runs.yaml
│   │   ├── runs_id.yaml
│   │   ├── tools.yaml
│   │   ├── tools_id.yaml
│   │   ├── users.yaml
│   │   └── users_id.yaml
│   ├── openapi.bundled.yaml
│   └── openapi.yaml
├── apps
│   ├── docs
│   │   ├── apis
│   │   │   └── openapi.yaml
│   │   ├── pages
│   │   │   ├── api-reference
│   │   │   │   └── overview.md
│   │   │   ├── architecture
│   │   │   │   ├── authentication.md
│   │   │   │   ├── overview.md
│   │   │   │   ├── project-layout.md
│   │   │   │   ├── project-structure.xml
│   │   │   │   └── system-design.md
│   │   │   ├── deployment
│   │   │   │   ├── docker.md
│   │   │   │   ├── kubernetes.md
│   │   │   │   ├── overview.md
│   │   │   │   └── production.md
│   │   │   ├── documentation
│   │   │   ├── features
│   │   │   │   ├── auth.md
│   │   │   │   ├── content.md
│   │   │   │   ├── organizations.md
│   │   │   │   ├── overview.md
│   │   │   │   ├── tui.md
│   │   │   │   └── workflows.md
│   │   │   ├── guides
│   │   │   │   ├── code-generation.md
│   │   │   │   ├── custom-handlers.md
│   │   │   │   ├── development.md
│   │   │   │   ├── makefile-commands.md
│   │   │   │   ├── overview.md
│   │   │   │   ├── test-coverage-report.md
│   │   │   │   └── testing.md
│   │   │   ├── monitoring
│   │   │   │   └── overview.md
│   │   │   ├── performance
│   │   │   │   ├── optimization.md
│   │   │   │   └── overview.md
│   │   │   ├── security
│   │   │   │   ├── best-practices.md
│   │   │   │   └── overview.md
│   │   │   ├── troubleshooting
│   │   │   │   └── common-issues.md
│   │   │   ├── codebase-analysis.md
│   │   │   ├── contributing.md
│   │   │   ├── getting-started.md
│   │   │   └── studio-ui-design.md
│   │   ├── src
│   │   │   ├── landing.tsx
│   │   │   ├── landing_content.ts
│   │   │   └── sidebar.tsx
│   │   ├── package.json
│   │   ├── tsconfig.app.json
│   │   ├── tsconfig.json
│   │   ├── tsconfig.spec.json
│   │   ├── vite.config.ts
│   │   └── zudoku.config.tsx
│   └── studio
│       ├── generated
│       │   ├── adapters
│       │   │   └── http
│       │   ├── application
│       │   │   ├── commands
│       │   │   └── queries
│       │   ├── core
│       │   │   ├── events
│       │   │   ├── models
│       │   │   └── repositories
│       │   └── infrastructure
│       │       ├── bootstrap
│       │       └── persistence
│       ├── handlers
│       │   ├── auth
│       │   │   ├── confirm_email_change_handler.go
│       │   │   ├── confirm_email_verification_handler.go
│       │   │   ├── confirm_password_reset_handler.go
│       │   │   ├── delete_account_handler.go
│       │   │   ├── delete_session_handler.go
│       │   │   ├── link_account_handler.go
│       │   │   ├── login_handler.go
│       │   │   ├── logout_all_handler.go
│       │   │   ├── logout_handler.go
│       │   │   ├── oauth_authorize_handler.go
│       │   │   ├── oauth_callback_handler.go
│       │   │   ├── register_handler.go
│       │   │   ├── request_email_change_handler.go
│       │   │   ├── request_email_verification_handler.go
│       │   │   ├── request_magic_link_handler.go
│       │   │   ├── request_password_reset_handler.go
│       │   │   ├── update_account_handler.go
│       │   │   ├── update_session_handler.go
│       │   │   └── verify_magic_link_handler.go
│       │   ├── config
│       │   │   └── get_config_handler.go
│       │   ├── executor
│       │   │   └── execute_executor_handler.go
│       │   ├── health
│       │   │   └── get_health_handler.go
│       │   ├── pipeline
│       │   │   ├── create_pipeline_step_handler.go
│       │   │   ├── get_pipeline_execution_plan_handler.go
│       │   │   ├── get_pipeline_steps_handler.go
│       │   │   └── validate_pipeline_execution_plan_handler.go
│       │   └── user
│       │       ├── delete_current_user_handler.go
│       │       ├── get_current_user_handler.go
│       │       └── update_current_user_handler.go
│       └── main.go
├── assets
│   ├── android-chrome-192x192.png
│   ├── android-chrome-512x512.png
│   ├── apple-touch-icon.png
│   ├── favicon-16x16.png
│   ├── favicon-32x32.png
│   ├── favicon.ico
│   ├── github-hero.svg
│   ├── large-logo-white.svg
│   ├── large-logo.svg
│   ├── site.webmanifest
│   ├── small-logo-adaptive.svg
│   ├── small-logo-white.svg
│   └── small-logo.svg
├── cmd
│   └── archesai
│       └── main.go
├── deployments
│   ├── containers
│   │   └── runners
│   │       ├── go
│   │       │   ├── Dockerfile
│   │       │   ├── execute.example.go
│   │       │   ├── execute.go
│   │       │   ├── go.mod
│   │       │   ├── go.sum
│   │       │   └── main.go
│   │       ├── node
│   │       │   ├── src
│   │       │   ├── Dockerfile
│   │       │   ├── execute.example.ts
│   │       │   ├── package.json
│   │       │   ├── tsconfig.json
│   │       │   ├── tsconfig.lib.json
│   │       │   └── tsconfig.spec.json
│   │       └── python
│   │           ├── Dockerfile
│   │           ├── execute.example.py
│   │           ├── execute.py
│   │           ├── requirements.txt
│   │           └── runner.py
│   ├── development
│   │   └── skaffold.yaml
│   ├── docker
│   │   ├── Dockerfile
│   │   ├── Dockerfile.goreleaser
│   │   └── docker-compose.yaml
│   ├── gcp
│   │   └── clouddeploy.yaml
│   ├── helm
│   │   ├── arches
│   │   │   ├── charts
│   │   │   │   └── ingress-nginx-4.13.0.tgz
│   │   │   ├── files
│   │   │   │   └── certs
│   │   │   ├── templates
│   │   │   │   ├── components
│   │   │   │   ├── _helpers.tpl
│   │   │   │   ├── configmap.yaml
│   │   │   │   ├── namespace.yaml
│   │   │   │   ├── secrets.yaml
│   │   │   │   └── serviceaccount.yaml
│   │   │   ├── Chart.yaml
│   │   │   └── values.yaml
│   │   └── dev-overrides.yaml
│   ├── helm-minimal
│   │   ├── charts
│   │   │   └── ingress-nginx-4.13.0.tgz
│   │   ├── files
│   │   │   ├── certs
│   │   │   │   ├── .gitkeep
│   │   │   │   ├── fullchain.pem
│   │   │   │   └── privkey.pem
│   │   │   └── kustomize
│   │   │       ├── base
│   │   │       └── components
│   │   ├── templates
│   │   │   └── kustomization.yaml
│   │   ├── Chart.yaml
│   │   ├── dev-values.yaml
│   │   └── values.yaml
│   ├── k3d
│   │   └── k3d.yaml
│   └── kustomize
│       ├── base
│       │   ├── kustomization.yaml
│       │   ├── namespace.yaml
│       │   └── serviceaccount.yaml
│       └── components
│           ├── api
│           │   ├── deployment.yaml
│           │   ├── kustomization.yaml
│           │   └── service.yaml
│           ├── database
│           │   ├── kustomization.yaml
│           │   ├── service.yaml
│           │   └── statefulset.yaml
│           ├── ingress
│           │   ├── api-ingress.yaml
│           │   ├── kustomization.yaml
│           │   └── platform-ingress.yaml
│           ├── migrations
│           │   ├── job.yaml
│           │   └── kustomization.yaml
│           ├── monitoring
│           │   ├── grafana-deployment.yaml
│           │   ├── grafana-service.yaml
│           │   ├── kustomization.yaml
│           │   ├── loki-deployment.yaml
│           │   └── loki-service.yaml
│           ├── platform
│           │   ├── deployment.yaml
│           │   ├── kustomization.yaml
│           │   └── service.yaml
│           ├── redis
│           │   ├── deployment.yaml
│           │   ├── kustomization.yaml
│           │   ├── pvc.yaml
│           │   └── service.yaml
│           ├── scraper
│           │   ├── deployment.yaml
│           │   ├── kustomization.yaml
│           │   └── service.yaml
│           ├── storage
│           │   ├── deployment.yaml
│           │   ├── kustomization.yaml
│           │   ├── pvc.yaml
│           │   └── service.yaml
│           └── unstructured
│               ├── deployment.yaml
│               ├── kustomization.yaml
│               └── service.yaml
├── docs
│   ├── api-reference
│   │   └── overview.md
│   ├── architecture
│   │   ├── authentication.md
│   │   ├── overview.md
│   │   ├── project-layout.md
│   │   ├── project-structure.xml
│   │   └── system-design.md
│   ├── deployment
│   │   ├── docker.md
│   │   ├── kubernetes.md
│   │   ├── overview.md
│   │   └── production.md
│   ├── features
│   │   ├── auth.md
│   │   ├── content.md
│   │   ├── organizations.md
│   │   ├── overview.md
│   │   ├── tui.md
│   │   └── workflows.md
│   ├── guides
│   │   ├── code-generation.md
│   │   ├── custom-handlers.md
│   │   ├── development.md
│   │   ├── makefile-commands.md
│   │   ├── overview.md
│   │   ├── test-coverage-report.md
│   │   └── testing.md
│   ├── monitoring
│   │   └── overview.md
│   ├── performance
│   │   ├── optimization.md
│   │   └── overview.md
│   ├── security
│   │   ├── best-practices.md
│   │   └── overview.md
│   ├── troubleshooting
│   │   └── common-issues.md
│   ├── codebase-analysis.md
│   ├── contributing.md
│   ├── getting-started.md
│   └── studio-ui-design.md
├── internal
│   ├── cli
│   │   ├── completion.go
│   │   ├── config.go
│   │   ├── dev.go
│   │   ├── generate.go
│   │   ├── root.go
│   │   ├── tui.go
│   │   └── version.go
│   ├── codegen
│   │   ├── tmpl
│   │   │   ├── bootstrap.go.tmpl
│   │   │   ├── command_handler.go.tmpl
│   │   │   ├── controller.go.tmpl
│   │   │   ├── db.hcl.tmpl
│   │   │   ├── events.go.tmpl
│   │   │   ├── header.tmpl
│   │   │   ├── infrastructure.go.tmpl
│   │   │   ├── query_handler.go.tmpl
│   │   │   ├── repository.go.tmpl
│   │   │   ├── repository_postgres.go.tmpl
│   │   │   ├── repository_sqlite.go.tmpl
│   │   │   ├── schema.go.tmpl
│   │   │   └── sqlc.yaml.tmpl
│   │   ├── generate.go
│   │   ├── generate_bootstrap.go
│   │   ├── generate_controllers.go
│   │   ├── generate_cqrs.go
│   │   ├── generate_events.go
│   │   ├── generate_hcl.go
│   │   ├── generate_js_client.go
│   │   ├── generate_migrations.go
│   │   ├── generate_repositories.go
│   │   ├── generate_schemas.go
│   │   ├── generate_sqlc.go
│   │   ├── renderer.go
│   │   └── templates.go
│   ├── dev
│   │   ├── manager.go
│   │   ├── process.go
│   │   └── watcher.go
│   ├── parsers
│   │   ├── jsonschema.go
│   │   ├── jsonschema_test.go
│   │   ├── openapi.go
│   │   ├── openapi_orvalfix.go
│   │   ├── openapi_test.go
│   │   ├── operation.go
│   │   ├── operation_test.go
│   │   ├── response.go
│   │   ├── schema.go
│   │   ├── spec.go
│   │   ├── strings.go
│   │   ├── typeconv.go
│   │   └── xcodegenextension.go
│   ├── testutil
│   │   └── containers.go
│   └── tui
│       ├── config.go
│       └── dev.go
├── pkg
│   ├── auth
│   │   ├── oauth
│   │   │   ├── github.go
│   │   │   ├── google.go
│   │   │   ├── microsoft.go
│   │   │   └── types.go
│   │   ├── auth_tokens.go
│   │   ├── magic_link.go
│   │   ├── password.go
│   │   ├── ports.go
│   │   ├── service.go
│   │   └── token_manager.go
│   ├── cache
│   │   ├── cache.go
│   │   ├── memory.go
│   │   ├── noop.go
│   │   └── redis.go
│   ├── config
│   │   ├── config.go
│   │   ├── loader.go
│   │   └── loader_test.go
│   ├── database
│   │   ├── crud_repository.go
│   │   ├── database.go
│   │   └── migrate.go
│   ├── errors
│   │   └── errors.go
│   ├── events
│   │   ├── events.go
│   │   ├── noop.go
│   │   ├── publisher.go
│   │   └── redis.go
│   ├── executor
│   │   ├── testdata
│   │   │   └── execute.ts
│   │   ├── builder.go
│   │   ├── builder_test.go
│   │   ├── config.go
│   │   ├── container.go
│   │   ├── container_test.go
│   │   ├── local.go
│   │   ├── local_test.go
│   │   ├── ports.go
│   │   ├── schemas.go
│   │   └── service.go
│   ├── llm
│   │   ├── chat.go
│   │   ├── interfaces.go
│   │   ├── llm.go
│   │   ├── ollama.go
│   │   └── openai.go
│   ├── logger
│   │   ├── config.go
│   │   └── logger.go
│   ├── notifications
│   │   ├── console.go
│   │   ├── email.go
│   │   ├── otp.go
│   │   └── service.go
│   ├── optional
│   │   └── optional.go
│   ├── redis
│   │   ├── client.go
│   │   ├── config.go
│   │   ├── errors.go
│   │   ├── pubsub.go
│   │   ├── queue.go
│   │   └── redis.go
│   ├── server
│   │   ├── cookies.go
│   │   ├── middleware.go
│   │   ├── middleware_auth.go
│   │   ├── middleware_logger.go
│   │   ├── middleware_ratelimit.go
│   │   ├── middleware_recover.go
│   │   ├── middleware_requestid.go
│   │   ├── middleware_security.go
│   │   ├── middleware_timeout.go
│   │   ├── responses.go
│   │   ├── server.go
│   │   └── websocket.go
│   └── storage
│       ├── disk.go
│       ├── memory.go
│       └── storage.go
├── scripts
│   ├── generate-coverage-report.sh
│   ├── generate-project-structure-xml.sh
│   ├── update-makefile-docs.sh
│   └── update-project-layout-docs.sh
├── test
│   └── data
│       ├── parsers
│       │   ├── invalid
│       │   │   └── missing-type.yaml
│       │   ├── openapi
│       │   │   └── simple-api.yaml
│       │   ├── schemas
│       │   │   ├── complex.yaml
│       │   │   ├── simple.yaml
│       │   │   ├── with-inheritance.yaml
│       │   │   └── with-x-codegen.yaml
│       │   └── x-codegen
│       ├── book.pdf
│       ├── pdf.png
│       ├── text.png
│       └── website.png
├── tools
│   └── tsconfig
│       ├── src
│       │   ├── base.json
│       │   ├── lib.json
│       │   ├── react.json
│       │   └── spec.json
│       └── package.json
├── web
│   ├── client
│   │   ├── src
│   │   │   ├── generated
│   │   │   │   ├── apikey
│   │   │   │   ├── artifact
│   │   │   │   ├── auth
│   │   │   │   ├── config
│   │   │   │   ├── executor
│   │   │   │   ├── health
│   │   │   │   ├── invitation
│   │   │   │   ├── label
│   │   │   │   ├── member
│   │   │   │   ├── organization
│   │   │   │   ├── pipeline
│   │   │   │   ├── run
│   │   │   │   ├── tool
│   │   │   │   ├── user
│   │   │   │   ├── orval.schemas.ts
│   │   │   │   └── zod.ts
│   │   │   ├── fetcher.ts
│   │   │   ├── index.ts
│   │   │   └── validators.ts
│   │   ├── orval.config.ts
│   │   ├── package.json
│   │   ├── tsconfig.json
│   │   ├── tsconfig.lib.json
│   │   └── tsconfig.spec.json
│   ├── platform
│   │   ├── public
│   │   │   └── .gitkeep
│   │   ├── src
│   │   │   ├── components
│   │   │   │   ├── auth
│   │   │   │   ├── containers
│   │   │   │   ├── datatables
│   │   │   │   ├── forms
│   │   │   │   ├── navigation
│   │   │   │   ├── selectors
│   │   │   │   ├── create-pipeline.tsx
│   │   │   │   ├── default-catch-boundary.tsx
│   │   │   │   ├── example.spec.ts
│   │   │   │   ├── file-upload.tsx
│   │   │   │   ├── not-found.tsx
│   │   │   │   └── terms-indicator.tsx
│   │   │   ├── hooks
│   │   │   │   ├── use-data-table.tsx
│   │   │   │   ├── use-filter-state.tsx
│   │   │   │   ├── use-offline-indicator.tsx
│   │   │   │   ├── use-toggle-view.tsx
│   │   │   │   └── use-websockets.tsx
│   │   │   ├── lib
│   │   │   │   ├── config.ts
│   │   │   │   ├── get-session-ssr.ts
│   │   │   │   ├── site-config.ts
│   │   │   │   └── site-utils.ts
│   │   │   ├── routes
│   │   │   │   ├── _app
│   │   │   │   ├── auth
│   │   │   │   └── __root.tsx
│   │   │   ├── styles
│   │   │   │   └── globals.css
│   │   │   ├── routeTree.gen.ts
│   │   │   └── router.tsx
│   │   ├── Dockerfile
│   │   ├── package.json
│   │   ├── playwright.config.js
│   │   ├── tsconfig.app.json
│   │   ├── tsconfig.json
│   │   ├── tsconfig.spec.json
│   │   └── vite.config.ts
│   └── ui
│       ├── src
│       │   ├── components
│       │   │   ├── custom
│       │   │   ├── datatable
│       │   │   ├── primitives
│       │   │   ├── shadcn
│       │   │   ├── zudoku
│       │   │   └── index.ts
│       │   ├── hooks
│       │   │   ├── use-callback-ref.tsx
│       │   │   ├── use-debounced-callback.tsx
│       │   │   ├── use-is-top.tsx
│       │   │   └── use-mobile.tsx
│       │   ├── layouts
│       │   │   ├── app-sidebar
│       │   │   ├── page-header
│       │   │   └── index.ts
│       │   ├── lib
│       │   │   ├── base-colors.ts
│       │   │   ├── compose-refs.ts
│       │   │   ├── constants.ts
│       │   │   ├── export.ts
│       │   │   ├── format.ts
│       │   │   ├── seo.ts
│       │   │   ├── site-config.interface.ts
│       │   │   └── utils.ts
│       │   ├── providers
│       │   │   ├── index.ts
│       │   │   ├── theme-provider.tsx
│       │   │   └── vite-theme-provider.tsx
│       │   ├── styles
│       │   │   └── globals.css
│       │   ├── types
│       │   │   ├── entities.ts
│       │   │   ├── simple-data-table.ts
│       │   │   └── table-meta.d.ts
│       │   └── index.ts
│       ├── components.json
│       ├── package.json
│       ├── tsconfig.json
│       ├── tsconfig.lib.json
│       ├── tsconfig.spec.json
│       └── vite.config.ts
├── .air.toml
├── .cspell.json
├── .editorconfig
├── .env
├── .gitignore
├── .golangci.yaml
├── .goreleaser.yaml
├── .lefthook.yaml
├── .markdownlint.json
├── .mockery.yaml
├── .prettierignore
├── LICENSE
├── Makefile
├── README.md
├── arches.yaml
├── biome.json
├── go.mod
├── go.sum
├── package.json
├── pnpm-lock.yaml
├── pnpm-workspace.yaml
├── tools.mod
├── tools.sum
└── tsconfig.json

184 directories, 611 files
```
