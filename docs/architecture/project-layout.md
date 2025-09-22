# Project Layout

## Directory Structure

```text
.
├── api
│   ├── components
│   │   ├── parameters
│   │   │   ├── AccountsFilter.yaml
│   │   │   ├── AccountsSort.yaml
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
│   │   │   └── Unauthorized.yaml
│   │   └── schemas
│   │       ├── Account.yaml
│   │       ├── APIConfig.yaml
│   │       ├── ApiKeyResponse.yaml
│   │       ├── ApiKey.yaml
│   │       ├── ArchesConfig.yaml
│   │       ├── Artifact.yaml
│   │       ├── AuthConfig.yaml
│   │       ├── Base.yaml
│   │       ├── BillingConfig.yaml
│   │       ├── CodegenConfig.yaml
│   │       ├── DatabaseAuth.yaml
│   │       ├── DatabaseConfig.yaml
│   │       ├── EmailConfig.yaml
│   │       ├── Email.yaml
│   │       ├── EmbeddingConfig.yaml
│   │       ├── FilterNode.yaml
│   │       ├── FirebaseAuth.yaml
│   │       ├── GitHubAuth.yaml
│   │       ├── GoogleAuth.yaml
│   │       ├── GrafanaConfig.yaml
│   │       ├── HealthResponse.yaml
│   │       ├── ImageConfig.yaml
│   │       ├── ImagesConfig.yaml
│   │       ├── InfrastructureConfig.yaml
│   │       ├── IngressConfig.yaml
│   │       ├── IntelligenceConfig.yaml
│   │       ├── Invitation.yaml
│   │       ├── Label.yaml
│   │       ├── LLMConfig.yaml
│   │       ├── LocalAuth.yaml
│   │       ├── LoggingConfig.yaml
│   │       ├── LokiConfig.yaml
│   │       ├── MagicLinkAuth.yaml
│   │       ├── MagicLinkToken.yaml
│   │       ├── Member.yaml
│   │       ├── MicrosoftAuth.yaml
│   │       ├── MigrationsConfig.yaml
│   │       ├── MonitoringConfig.yaml
│   │       ├── OrganizationReference.yaml
│   │       ├── Organization.yaml
│   │       ├── Page.yaml
│   │       ├── PersistenceConfig.yaml
│   │       ├── PipelineStep.yaml
│   │       ├── Pipeline.yaml
│   │       ├── PlatformConfig.yaml
│   │       ├── Problem.yaml
│   │       ├── RedisConfig.yaml
│   │       ├── ResourceConfig.yaml
│   │       ├── ResourceLimits.yaml
│   │       ├── ResourceRequests.yaml
│   │       ├── RunPodConfig.yaml
│   │       ├── Run.yaml
│   │       ├── ScraperConfig.yaml
│   │       ├── ServiceAccountConfig.yaml
│   │       ├── Session.yaml
│   │       ├── SpeechConfig.yaml
│   │       ├── StorageConfig.yaml
│   │       ├── StripeConfig.yaml
│   │       ├── TLSConfig.yaml
│   │       ├── TokenResponse.yaml
│   │       ├── Tool.yaml
│   │       ├── TwitterAuth.yaml
│   │       ├── UnstructuredConfig.yaml
│   │       ├── User.yaml
│   │       ├── UUID.yaml
│   │       ├── ValidationError.yaml
│   │       ├── XCodegenWrapper.yaml
│   │       └── XCodegen.yaml
│   ├── paths
│   │   ├── accounts_email-change_request.yaml
│   │   ├── accounts_email-change_verify.yaml
│   │   ├── accounts_email-verification_request.yaml
│   │   ├── accounts_email-verification_verify.yaml
│   │   ├── accounts_{id}.yaml
│   │   ├── accounts_password-reset_request.yaml
│   │   ├── accounts_password-reset_verify.yaml
│   │   ├── accounts.yaml
│   │   ├── artifacts_{id}.yaml
│   │   ├── artifacts.yaml
│   │   ├── auth_magic-link.yaml
│   │   ├── auth_register.yaml
│   │   ├── config.yaml
│   │   ├── health.yaml
│   │   ├── invitations_{id}.yaml
│   │   ├── invitations.yaml
│   │   ├── labels_{id}.yaml
│   │   ├── labels.yaml
│   │   ├── members_{id}.yaml
│   │   ├── members.yaml
│   │   ├── oauth_authorize.yaml
│   │   ├── oauth_callback.yaml
│   │   ├── organizations_{id}.yaml
│   │   ├── organizations.yaml
│   │   ├── pipelines_{id}_execution-plans.yaml
│   │   ├── pipelines_{id}_steps.yaml
│   │   ├── pipelines_{id}.yaml
│   │   ├── pipelines.yaml
│   │   ├── runs_{id}.yaml
│   │   ├── runs.yaml
│   │   ├── sessions_create.yaml
│   │   ├── sessions_delete.yaml
│   │   ├── sessions_{id}.yaml
│   │   ├── sessions.yaml
│   │   ├── tokens_{id}.yaml
│   │   ├── tokens.yaml
│   │   ├── tools_{id}.yaml
│   │   ├── tools.yaml
│   │   ├── users_{id}.yaml
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
│   └── archesai
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
│   ├── kustomize
│   │   ├── base
│   │   │   ├── kustomization.yaml
│   │   │   ├── namespace.yaml
│   │   │   └── serviceaccount.yaml
│   │   └── components
│   │       ├── api
│   │       │   ├── deployment.yaml
│   │       │   ├── kustomization.yaml
│   │       │   └── service.yaml
│   │       ├── database
│   │       │   ├── kustomization.yaml
│   │       │   ├── service.yaml
│   │       │   └── statefulset.yaml
│   │       ├── ingress
│   │       │   ├── api-ingress.yaml
│   │       │   ├── kustomization.yaml
│   │       │   └── platform-ingress.yaml
│   │       ├── migrations
│   │       │   ├── job.yaml
│   │       │   └── kustomization.yaml
│   │       ├── monitoring
│   │       │   ├── grafana-deployment.yaml
│   │       │   ├── grafana-service.yaml
│   │       │   ├── kustomization.yaml
│   │       │   ├── loki-deployment.yaml
│   │       │   └── loki-service.yaml
│   │       ├── platform
│   │       │   ├── deployment.yaml
│   │       │   ├── kustomization.yaml
│   │       │   └── service.yaml
│   │       ├── redis
│   │       │   ├── deployment.yaml
│   │       │   ├── kustomization.yaml
│   │       │   ├── pvc.yaml
│   │       │   └── service.yaml
│   │       ├── scraper
│   │       │   ├── deployment.yaml
│   │       │   ├── kustomization.yaml
│   │       │   └── service.yaml
│   │       ├── storage
│   │       │   ├── deployment.yaml
│   │       │   ├── kustomization.yaml
│   │       │   ├── pvc.yaml
│   │       │   └── service.yaml
│   │       └── unstructured
│   │           ├── deployment.yaml
│   │           ├── kustomization.yaml
│   │           └── service.yaml
│   └── scripts
│       └── deploy.sh
├── .devcontainer
│   ├── devcontainer.json
│   ├── Dockerfile
│   └── init-firewall.sh
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
│   ├── accounts
│   │   ├── accounts.go
│   │   ├── errors.go
│   │   ├── generate.go
│   │   └── service_test.go
│   ├── app
│   │   ├── app.go
│   │   ├── infrastructure.go
│   │   └── routes.go
│   ├── artifacts
│   │   ├── artifacts.go
│   │   ├── errors.go
│   │   ├── generate.go
│   │   └── service_test.go
│   ├── auth
│   │   ├── deliverers
│   │   │   ├── console.go
│   │   │   ├── email.go
│   │   │   └── otp.go
│   │   ├── providers
│   │   │   ├── github.go
│   │   │   ├── google.go
│   │   │   └── microsoft.go
│   │   ├── stores
│   │   │   ├── api_token_repository.go
│   │   │   ├── api_token_store.go
│   │   │   ├── magic_link_repository.go
│   │   │   ├── magic_link_store.go
│   │   │   └── session_store.go
│   │   ├── strategies
│   │   ├── tokens
│   │   │   ├── api_validator.go
│   │   │   ├── manager.go
│   │   │   └── utils.go
│   │   ├── claims.go
│   │   ├── claims_test.go
│   │   ├── context.go
│   │   ├── generate.go
│   │   ├── handler.go
│   │   ├── middleware.go
│   │   ├── permissions.go
│   │   ├── permissions_test.go
│   │   ├── routes.go
│   │   ├── service.go
│   │   ├── sessionmanager_test.go
│   │   ├── sessions_repository.go
│   │   ├── sessions_types.go
│   │   └── tokens_types.go
│   ├── cache
│   │   ├── cache.go
│   │   ├── memory.go
│   │   ├── noop.go
│   │   └── redis.go
│   ├── cli
│   │   ├── all.go
│   │   ├── api.go
│   │   ├── completion.go
│   │   ├── config.go
│   │   ├── root.go
│   │   ├── tui.go
│   │   ├── version.go
│   │   ├── web.go
│   │   └── worker.go
│   ├── codegen
│   │   ├── templates
│   │   │   ├── config.go.tmpl
│   │   │   ├── events.go.tmpl
│   │   │   ├── events_nats.go.tmpl
│   │   │   ├── events_redis.go.tmpl
│   │   │   ├── handler.gen.go.tmpl
│   │   │   ├── repository.go.tmpl
│   │   │   ├── repository_postgres.go.tmpl
│   │   │   ├── repository_sqlite.go.tmpl
│   │   │   └── service.go.tmpl
│   │   ├── codegen.go
│   │   ├── defaults.go
│   │   ├── filewriter.go
│   │   ├── generate.go
│   │   ├── helpers.go
│   │   ├── parser.go
│   │   ├── sql_generator.go
│   │   └── template_funcs.go
│   ├── config
│   │   ├── config.go
│   │   ├── generate.go
│   │   ├── handler.go
│   │   ├── loader.go
│   │   └── loader_test.go
│   ├── database
│   │   ├── postgresql
│   │   │   ├── accounts.sql.go
│   │   │   ├── api-tokens.sql.go
│   │   │   ├── artifacts.sql.go
│   │   │   ├── db.go
│   │   │   ├── invitations.sql.go
│   │   │   ├── labels.sql.go
│   │   │   ├── members.sql.go
│   │   │   ├── models.go
│   │   │   ├── organizations.sql.go
│   │   │   ├── pipelines.sql.go
│   │   │   ├── pipeline-step-dependencies.sql.go
│   │   │   ├── pipeline-steps.sql.go
│   │   │   ├── querier.go
│   │   │   ├── runs.sql.go
│   │   │   ├── sessions.sql.go
│   │   │   ├── tools.sql.go
│   │   │   ├── users.sql.go
│   │   │   └── verification-tokens.sql.go
│   │   ├── queries
│   │   │   ├── accounts.sql
│   │   │   ├── api-tokens.sql
│   │   │   ├── artifacts.sql
│   │   │   ├── invitations.sql
│   │   │   ├── labels.sql
│   │   │   ├── members.sql
│   │   │   ├── organizations.sql
│   │   │   ├── pipelines.sql
│   │   │   ├── pipeline-step-dependencies.sql
│   │   │   ├── pipeline-steps.sql
│   │   │   ├── runs.sql
│   │   │   ├── sessions.sql
│   │   │   ├── tools.sql
│   │   │   ├── users.sql
│   │   │   └── verification-tokens.sql
│   │   ├── sqlite
│   │   │   └── stub.go
│   │   ├── database.go
│   │   ├── generate.go
│   │   └── sqlc.yaml
│   ├── events
│   │   ├── events.go
│   │   ├── noop.go
│   │   ├── publisher.go
│   │   └── redis.go
│   ├── health
│   │   ├── errors.go
│   │   ├── generate.go
│   │   ├── handler.go
│   │   ├── health.go
│   │   ├── postgres.go
│   │   ├── repository.go
│   │   ├── service.go
│   │   ├── service_test.go
│   │   └── sqlite.go
│   ├── invitations
│   │   ├── errors.go
│   │   ├── generate.go
│   │   └── invitations.go
│   ├── labels
│   │   ├── errors.go
│   │   ├── generate.go
│   │   ├── labels.go
│   │   └── service_test.go
│   ├── llm
│   │   ├── chat.go
│   │   ├── llm.go
│   │   ├── ollama.go
│   │   └── openai.go
│   ├── logger
│   │   ├── config.go
│   │   └── logger.go
│   ├── members
│   │   ├── errors.go
│   │   ├── generate.go
│   │   └── members.go
│   ├── middleware
│   │   ├── auth.go
│   │   └── ratelimit.go
│   ├── migrations
│   │   ├── postgresql
│   │   │   ├── 20250908082709_init.sql
│   │   │   ├── 20250908124238_enable_pgvector.sql
│   │   │   ├── 20250915075921_update_api_token_structure.sql
│   │   │   └── 20250919061512_add_magic_link_auth.sql
│   │   ├── sqlite
│   │   │   ├── 20250908082709_init.sql
│   │   │   ├── 20250908124238_enable_pgvector.sql
│   │   │   ├── 20250915075921_update_api_token_structure.sql
│   │   │   └── 20250919061512_add_magic_link_auth.sql
│   │   └── migrate.go
│   ├── organizations
│   │   ├── errors.go
│   │   ├── generate.go
│   │   ├── organizations.go
│   │   └── service_test.go
│   ├── pipelines
│   │   ├── dag.go
│   │   ├── dag_test.go
│   │   ├── errors.go
│   │   ├── executor.go
│   │   ├── executor_test.go
│   │   ├── generate.go
│   │   ├── handler.go
│   │   ├── manager.go
│   │   ├── pipelines.go
│   │   ├── queue_redis.go
│   │   ├── service_test.go
│   │   └── step_repository.go
│   ├── redis
│   │   ├── client.go
│   │   ├── config.go
│   │   ├── errors.go
│   │   ├── pubsub.go
│   │   ├── queue.go
│   │   └── redis.go
│   ├── runs
│   │   ├── errors.go
│   │   ├── generate.go
│   │   └── runs.go
│   ├── server
│   │   ├── assets
│   │   │   └── docs.html
│   │   ├── docs.go
│   │   ├── errors.go
│   │   ├── infra.go
│   │   ├── middleware.go
│   │   ├── server.go
│   │   └── websocket.go
│   ├── storage
│   │   ├── storage.go
│   │   └── storage_test.go
│   ├── testutil
│   │   └── containers.go
│   ├── tools
│   │   ├── errors.go
│   │   ├── generate.go
│   │   └── tools.go
│   ├── tui
│   │   ├── config_tui.go
│   │   └── tui.go
│   └── users
│       ├── errors.go
│       ├── generate.go
│       ├── oauth.go
│       ├── service_extra.go
│       ├── service_test.go
│       └── users.go
├── scripts
│   ├── add-mapstructure-tags.sh
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
│       ├── book.pdf
│       ├── pdf.png
│       ├── text.png
│       └── website.png
├── tmp
│   └── archesai
├── tools
│   ├── codegen
│   │   └── main.go
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
│   │   │   │   ├── accounts
│   │   │   │   ├── artifacts
│   │   │   │   ├── auth
│   │   │   │   ├── config
│   │   │   │   ├── health
│   │   │   │   ├── invitations
│   │   │   │   ├── labels
│   │   │   │   ├── members
│   │   │   │   ├── oauth
│   │   │   │   ├── organizations
│   │   │   │   ├── pipelines
│   │   │   │   ├── runs
│   │   │   │   ├── sessions
│   │   │   │   ├── tokens
│   │   │   │   ├── tools
│   │   │   │   ├── users
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
├── .codegen.archesai.yaml
├── .codegen.server.yaml
├── .codegen.types.yaml
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
├── .redocly.yaml
└── tsconfig.json

179 directories, 605 files
```
