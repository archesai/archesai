# Project Layout

## Directory Structure

```text
.
├── api
│   ├── components
│   │   ├── parameters
│   │   │   ├── AccountsFilter.yaml
│   │   │   ├── AccountsSort.yaml
│   │   │   ├── APIKeysFilter.yaml
│   │   │   ├── APIKeysSort.yaml
│   │   │   ├── ArtifactsFilter.yaml
│   │   │   ├── ArtifactsSort.yaml
│   │   │   ├── InvitationsFilter.yaml
│   │   │   ├── InvitationsSort.yaml
│   │   │   ├── LabelsFilter.yaml
│   │   │   ├── LabelsSort.yaml
│   │   │   ├── MembersFilter.yaml
│   │   │   ├── MembersSort.yaml
│   │   │   ├── OrganizationsFilter.yaml
│   │   │   ├── OrganizationsSort.yaml
│   │   │   ├── PageQuery.yaml
│   │   │   ├── PipelinesFilter.yaml
│   │   │   ├── PipelinesSort.yaml
│   │   │   ├── RunsFilter.yaml
│   │   │   ├── RunsSort.yaml
│   │   │   ├── SessionsFilter.yaml
│   │   │   ├── SessionsSort.yaml
│   │   │   ├── ToolsFilter.yaml
│   │   │   ├── ToolsSort.yaml
│   │   │   ├── UsersFilter.yaml
│   │   │   └── UsersSort.yaml
│   │   ├── responses
│   │   │   ├── BadRequest.yaml
│   │   │   ├── InternalServerError.yaml
│   │   │   ├── NoContent.yaml
│   │   │   ├── NotFound.yaml
│   │   │   ├── TooManyRequests.yaml
│   │   │   ├── Unauthorized.yaml
│   │   │   └── UnprocessableEntity.yaml
│   │   └── schemas
│   │       ├── config
│   │       │   ├── ConfigAPI.yaml
│   │       │   ├── ConfigAuthFirebase.yaml
│   │       │   ├── ConfigAuthGithub.yaml
│   │       │   ├── ConfigAuthGoogle.yaml
│   │       │   ├── ConfigAuthLocal.yaml
│   │       │   ├── ConfigAuthMagicLink.yaml
│   │       │   ├── ConfigAuthMicrosoft.yaml
│   │       │   ├── ConfigAuthTwitter.yaml
│   │       │   ├── ConfigAuth.yaml
│   │       │   ├── ConfigBilling.yaml
│   │       │   ├── ConfigDatabase.yaml
│   │       │   ├── ConfigEmail.yaml
│   │       │   ├── ConfigGrafana.yaml
│   │       │   ├── ConfigImages.yaml
│   │       │   ├── ConfigImage.yaml
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
│   │       │   ├── ConfigUnstructured.yaml
│   │       │   └── Config.yaml
│   │       ├── xcodegen
│   │       │   ├── CodegenExtension.yaml
│   │       │   └── JSONSchema2020Extended.yaml
│   │       ├── Account.yaml
│   │       ├── APIKey.yaml
│   │       ├── Artifact.yaml
│   │       ├── Base.yaml
│   │       ├── FilterNode.yaml
│   │       ├── Health.yaml
│   │       ├── Invitation.yaml
│   │       ├── Label.yaml
│   │       ├── MagicLinkToken.yaml
│   │       ├── Member.yaml
│   │       ├── Organization.yaml
│   │       ├── Page.yaml
│   │       ├── PipelineStep.yaml
│   │       ├── Pipeline.yaml
│   │       ├── Problem.yaml
│   │       ├── Run.yaml
│   │       ├── Session.yaml
│   │       ├── Tool.yaml
│   │       └── User.yaml
│   ├── paths
│   │   ├── api-keys_id.yaml
│   │   ├── api-keys.yaml
│   │   ├── artifacts_id.yaml
│   │   ├── artifacts.yaml
│   │   ├── auth_accounts_id.yaml
│   │   ├── auth_accounts.yaml
│   │   ├── auth_change-email.yaml
│   │   ├── auth_confirm-email.yaml
│   │   ├── auth_forgot-password.yaml
│   │   ├── auth_link.yaml
│   │   ├── auth_login.yaml
│   │   ├── auth_logout-all.yaml
│   │   ├── auth_logout.yaml
│   │   ├── auth_magic-link-request.yaml
│   │   ├── auth_magic-link-verify.yaml
│   │   ├── auth_register.yaml
│   │   ├── auth_request-verification.yaml
│   │   ├── auth_reset-password.yaml
│   │   ├── auth_sessions_id.yaml
│   │   ├── auth_sessions.yaml
│   │   ├── auth_verify-email.yaml
│   │   ├── config.yaml
│   │   ├── health.yaml
│   │   ├── invitations_id.yaml
│   │   ├── invitations.yaml
│   │   ├── labels_id.yaml
│   │   ├── labels.yaml
│   │   ├── members_id.yaml
│   │   ├── members.yaml
│   │   ├── oauth_authorize.yaml
│   │   ├── oauth_callback.yaml
│   │   ├── organizations_id.yaml
│   │   ├── organizations.yaml
│   │   ├── pipelines_id_execution-plans.yaml
│   │   ├── pipelines_id_steps.yaml
│   │   ├── pipelines_id.yaml
│   │   ├── pipelines.yaml
│   │   ├── runs_id.yaml
│   │   ├── runs.yaml
│   │   ├── tools_id.yaml
│   │   ├── tools.yaml
│   │   ├── users_id.yaml
│   │   ├── users_me.yaml
│   │   └── users.yaml
│   ├── openapi.bundled.yaml
│   └── openapi.yaml
├── assets
│   ├── android-chrome-192x192.png
│   ├── android-chrome-512x512.png
│   ├── apple-touch-icon.png
│   ├── favicon-16x16.png
│   ├── favicon-32x32.png
│   ├── favicon.ico
│   ├── github-hero.svg
│   ├── large-logo.svg
│   ├── large-logo-white.svg
│   ├── site.webmanifest
│   ├── small-logo-adaptive.svg
│   ├── small-logo.svg
│   └── small-logo-white.svg
├── cmd
│   ├── archesai
│   │   └── main.go
│   └── codegen
│       └── main.go
├── deployments
│   ├── development
│   │   └── skaffold.yaml
│   ├── docker
│   │   ├── docker-compose.yaml
│   │   ├── Dockerfile
│   │   └── Dockerfile.goreleaser
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
│   │   │   │   ├── configmap.yaml
│   │   │   │   ├── _helpers.tpl
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
│   │   │   │   ├── fullchain.pem
│   │   │   │   ├── .gitkeep
│   │   │   │   └── privkey.pem
│   │   │   └── kustomize
│   │   │       ├── base
│   │   │       └── components
│   │   ├── templates
│   │   │   └── kustomization.yaml
│   │   ├── Chart.yaml
│   │   ├── dev-values.yaml
│   │   ├── values.schema.json
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
│   ├── contributing.md
│   └── getting-started.md
├── internal
│   ├── adapters
│   │   ├── cli
│   │   │   ├── all.go
│   │   │   ├── api.go
│   │   │   ├── completion.go
│   │   │   ├── config.go
│   │   │   ├── root.go
│   │   │   ├── tui.go
│   │   │   ├── version.go
│   │   │   ├── web.go
│   │   │   └── worker.go
│   │   ├── http
│   │   │   ├── controllers
│   │   │   └── server
│   │   │       ├── assets
│   │   │       ├── cookies.go
│   │   │       ├── docs.go
│   │   │       ├── errors.go
│   │   │       ├── infra.go
│   │   │       ├── middleware.go
│   │   │       ├── responses.go
│   │   │       ├── router.go
│   │   │       ├── server.go
│   │   │       └── websocket.go
│   │   └── tui
│   │       ├── screens
│   │       ├── config_tui.go
│   │       └── tui.go
│   ├── application
│   │   ├── commands
│   │   │   ├── apikey
│   │   │   ├── artifact
│   │   │   ├── auth
│   │   │   │   ├── confirm_email_change_handler.go
│   │   │   │   ├── confirm_email_verification_handler.go
│   │   │   │   ├── confirm_password_reset_handler.go
│   │   │   │   ├── delete_account_handler.go
│   │   │   │   ├── delete_session_handler.go
│   │   │   │   ├── link_account_handler.go
│   │   │   │   ├── login_handler.go
│   │   │   │   ├── logout_all_handler.go
│   │   │   │   ├── logout_handler.go
│   │   │   │   ├── register_handler.go
│   │   │   │   ├── request_email_change_handler.go
│   │   │   │   ├── request_email_verification_handler.go
│   │   │   │   ├── request_magic_link_handler.go
│   │   │   │   ├── request_password_reset_handler.go
│   │   │   │   ├── update_account_handler.go
│   │   │   │   ├── update_session_handler.go
│   │   │   │   └── verify_magic_link_handler.go
│   │   │   ├── invitation
│   │   │   ├── label
│   │   │   ├── member
│   │   │   ├── organization
│   │   │   ├── pipeline
│   │   │   │   ├── create_pipeline_step_handler.go
│   │   │   │   └── validate_pipeline_execution_plan_handler.go
│   │   │   ├── run
│   │   │   ├── tool
│   │   │   └── user
│   │   │       ├── delete_current_user_handler.go
│   │   │       └── update_current_user_handler.go
│   │   └── queries
│   │       ├── apikey
│   │       ├── artifact
│   │       ├── auth
│   │       │   ├── oauth_authorize_handler.go
│   │       │   └── oauth_callback_handler.go
│   │       ├── config
│   │       │   └── get_config_handler.go
│   │       ├── health
│   │       │   └── get_health_handler.go
│   │       ├── invitation
│   │       ├── label
│   │       ├── member
│   │       ├── organization
│   │       ├── pipeline
│   │       │   ├── get_pipeline_execution_plan_handler.go
│   │       │   └── get_pipeline_steps_handler.go
│   │       ├── run
│   │       ├── tool
│   │       └── user
│   │           └── get_current_user_handler.go
│   ├── codegen
│   │   ├── tmpl
│   │   │   ├── bootstrap.tmpl
│   │   │   ├── command_handler.tmpl
│   │   │   ├── controller.tmpl
│   │   │   ├── events.tmpl
│   │   │   ├── header.tmpl
│   │   │   ├── infrastructure.tmpl
│   │   │   ├── query_handler.tmpl
│   │   │   ├── repository_postgres.tmpl
│   │   │   ├── repository_sqlite.tmpl
│   │   │   ├── repository.tmpl
│   │   │   └── schema.tmpl
│   │   ├── filewriter.go
│   │   ├── generate_bootstrap.go
│   │   ├── generate_controllers.go
│   │   ├── generate_cqrs.go
│   │   ├── generate_events.go
│   │   ├── generate.go
│   │   ├── generate_repositories.go
│   │   ├── generate_schemas.go
│   │   └── templates.go
│   ├── core
│   │   ├── entities
│   │   ├── errors
│   │   │   └── errors.go
│   │   ├── events
│   │   │   ├── event.go
│   │   │   └── publisher.go
│   │   ├── repositories
│   │   │   └── health.go
│   │   ├── services
│   │   │   ├── auth.go
│   │   │   └── llm.go
│   │   └── valueobjects
│   │       ├── auth_tokens.go
│   │       ├── llm.go
│   │       └── stub.go
│   ├── infrastructure
│   │   ├── auth
│   │   │   ├── oauth
│   │   │   │   ├── github.go
│   │   │   │   ├── google.go
│   │   │   │   ├── microsoft.go
│   │   │   │   └── types.go
│   │   │   ├── magic_link.go
│   │   │   ├── password.go
│   │   │   ├── service.go
│   │   │   └── token_manager.go
│   │   ├── bootstrap
│   │   ├── cache
│   │   │   ├── cache.go
│   │   │   ├── memory.go
│   │   │   ├── noop.go
│   │   │   └── redis.go
│   │   ├── config
│   │   │   ├── config.go
│   │   │   ├── loader.go
│   │   │   └── loader_test.go
│   │   ├── events
│   │   │   ├── events.go
│   │   │   ├── noop.go
│   │   │   ├── publisher.go
│   │   │   └── redis.go
│   │   ├── http
│   │   ├── llm
│   │   │   ├── chat.go
│   │   │   ├── llm.go
│   │   │   ├── ollama.go
│   │   │   └── openai.go
│   │   ├── notifications
│   │   │   ├── console.go
│   │   │   ├── email.go
│   │   │   ├── otp.go
│   │   │   └── service.go
│   │   ├── persistence
│   │   │   ├── postgres
│   │   │   │   ├── migrations
│   │   │   │   ├── queries
│   │   │   │   ├── repositories
│   │   │   │   └── sqlc.yaml
│   │   │   ├── sqlite
│   │   │   │   ├── migrations
│   │   │   │   ├── queries
│   │   │   │   └── repositories
│   │   │   ├── database.go
│   │   │   └── migrate.go
│   │   ├── redis
│   │   │   ├── client.go
│   │   │   ├── config.go
│   │   │   ├── errors.go
│   │   │   ├── pubsub.go
│   │   │   ├── queue.go
│   │   │   └── redis.go
│   │   └── storage
│   │       ├── local
│   │       ├── s3
│   │       ├── storage.go
│   │       └── storage_test.go
│   ├── parsers
│   │   ├── jsonschema.go
│   │   ├── jsonschema_test.go
│   │   ├── openapi.go
│   │   ├── openapi_test.go
│   │   ├── strings.go
│   │   ├── typeconv.go
│   │   ├── types.go
│   │   ├── types_test.go
│   │   └── xcodegenextension.go
│   └── shared
│       ├── logger
│       │   ├── config.go
│       │   └── logger.go
│       └── testutil
│           └── containers.go
├── scripts
│   ├── generate-coverage-report.sh
│   ├── generate-helm-schema.py
│   ├── generate-project-structure-xml.sh
│   ├── update-makefile-docs.sh
│   └── update-project-layout-docs.sh
├── .taskmaster
│   ├── docs
│   │   └── prd.txt
│   ├── reports
│   │   └── task-complexity-report_main.json
│   ├── tasks
│   │   └── tasks.json
│   ├── templates
│   │   └── example_prd.txt
│   ├── CLAUDE.md
│   ├── config.json
│   └── state.json
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
│   ├── pg-to-sqlite
│   │   └── main.go
│   └── tsconfig
│       ├── src
│       │   ├── base.json
│       │   ├── lib.json
│       │   ├── react.json
│       │   └── spec.json
│       └── package.json
├── .vscode
│   ├── extensions.json
│   ├── mcp.json
│   └── settings.json
├── web
│   ├── client
│   │   ├── src
│   │   │   ├── generated
│   │   │   │   ├── apikey
│   │   │   │   ├── artifact
│   │   │   │   ├── auth
│   │   │   │   ├── config
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
│   │   │   ├── contributing.md
│   │   │   └── getting-started.md
│   │   ├── public
│   │   │   └── .gitkeep
│   │   ├── src
│   │   │   ├── landing_content.ts
│   │   │   ├── landing.tsx
│   │   │   └── sidebar.tsx
│   │   ├── package.json
│   │   ├── tsconfig.app.json
│   │   ├── tsconfig.json
│   │   ├── tsconfig.spec.json
│   │   ├── vite.config.ts
│   │   └── zudoku.config.tsx
│   ├── platform
│   │   ├── public
│   │   │   └── .gitkeep
│   │   ├── src
│   │   │   ├── app
│   │   │   │   ├── _app
│   │   │   │   ├── auth
│   │   │   │   └── __root.tsx
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
│   │   │   │   ├── api
│   │   │   │   ├── config.ts
│   │   │   │   ├── get-session-ssr.ts
│   │   │   │   ├── site-config.ts
│   │   │   │   └── site-utils.ts
│   │   │   ├── styles
│   │   │   │   └── globals.css
│   │   │   ├── router.tsx
│   │   │   └── routeTree.gen.ts
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
├── arches.yaml
├── biome.json
├── .cspell.json
├── .editorconfig
├── .env
├── .gitignore
├── .golangci.yaml
├── go.mod
├── .goreleaser.yaml
├── go.sum
├── .lefthook.yaml
├── LICENSE
├── Makefile
├── .markdownlint.json
├── .mcp.json
├── .mockery.yaml
├── opencode.json
├── package.json
├── pnpm-lock.yaml
├── pnpm-workspace.yaml
├── .prettierignore
├── README.md
└── tsconfig.json

211 directories, 529 files
```
