# Project Layout

## Directory Structure

```text
.
├── .claude
│   ├── CLAUDE.md
│   └── settings.json
├── .github
│   ├── DISCUSSION_TEMPLATE
│   │   └── ideas.yml
│   ├── ISSUE_TEMPLATE
│   │   ├── bug_report.yaml
│   │   └── config.yaml
│   ├── actions
│   │   ├── build-project
│   │   │   └── action.yml
│   │   ├── create-github-release
│   │   │   └── action.yml
│   │   ├── create-release-tag
│   │   │   └── action.yml
│   │   ├── docker-retag
│   │   │   └── action.yml
│   │   ├── docker-setup
│   │   │   └── action.yml
│   │   ├── goreleaser-run
│   │   │   └── action.yml
│   │   └── setup-build-env
│   │       └── action.yml
│   ├── workflows
│   │   ├── claude-code-review.yaml
│   │   ├── claude.yaml
│   │   ├── deploy-docs.yaml
│   │   ├── docker-build-and-push.yaml
│   │   ├── lint-go.yaml
│   │   ├── lint-typescript.yaml
│   │   ├── release-edge.yaml
│   │   ├── release-nightly.yaml
│   │   ├── release.yaml
│   │   ├── test-coverage.yaml
│   │   ├── test-go.yaml
│   │   └── update-docs.yaml
│   └── dependabot.yaml
├── .vscode
│   ├── extensions.json
│   └── settings.json
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
│   │   │   ├── NoContent.yaml
│   │   │   ├── NotFound.yaml
│   │   │   └── Unauthorized.yaml
│   │   └── schemas
│   │       ├── APIConfig.yaml
│   │       ├── Account.yaml
│   │       ├── ApiKey.yaml
│   │       ├── ApiKeyResponse.yaml
│   │       ├── ArchesConfig.yaml
│   │       ├── Artifact.yaml
│   │       ├── AuthConfig.yaml
│   │       ├── Base.yaml
│   │       ├── BillingConfig.yaml
│   │       ├── CORSConfig.yaml
│   │       ├── CodegenConfig.yaml
│   │       ├── DatabaseAuth.yaml
│   │       ├── DatabaseConfig.yaml
│   │       ├── DevServiceConfig.yaml
│   │       ├── DevelopmentConfig.yaml
│   │       ├── Email.yaml
│   │       ├── EmailConfig.yaml
│   │       ├── EmbeddingConfig.yaml
│   │       ├── FilterNode.yaml
│   │       ├── FirebaseAuth.yaml
│   │       ├── GrafanaConfig.yaml
│   │       ├── HealthResponse.yaml
│   │       ├── ImageConfig.yaml
│   │       ├── ImagesConfig.yaml
│   │       ├── InfrastructureConfig.yaml
│   │       ├── IngressConfig.yaml
│   │       ├── IntelligenceConfig.yaml
│   │       ├── Invitation.yaml
│   │       ├── LLMConfig.yaml
│   │       ├── Label.yaml
│   │       ├── LocalAuth.yaml
│   │       ├── LoggingConfig.yaml
│   │       ├── LokiConfig.yaml
│   │       ├── Member.yaml
│   │       ├── MigrationsConfig.yaml
│   │       ├── MonitoringConfig.yaml
│   │       ├── Organization.yaml
│   │       ├── OrganizationReference.yaml
│   │       ├── Page.yaml
│   │       ├── PersistenceConfig.yaml
│   │       ├── Pipeline.yaml
│   │       ├── PipelineStep.yaml
│   │       ├── PlatformConfig.yaml
│   │       ├── Problem.yaml
│   │       ├── RedisConfig.yaml
│   │       ├── ResourceConfig.yaml
│   │       ├── ResourceLimits.yaml
│   │       ├── ResourceRequests.yaml
│   │       ├── Run.yaml
│   │       ├── RunPodConfig.yaml
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
│   │       ├── UUID.yaml
│   │       ├── UnstructuredConfig.yaml
│   │       ├── User.yaml
│   │       ├── ValidationError.yaml
│   │       ├── XCodegen.yaml
│   │       └── XCodegenWrapper.yaml
│   ├── paths
│   │   ├── accounts.yaml
│   │   ├── accounts_email-change_request.yaml
│   │   ├── accounts_email-change_verify.yaml
│   │   ├── accounts_email-verification_request.yaml
│   │   ├── accounts_email-verification_verify.yaml
│   │   ├── accounts_password-reset_request.yaml
│   │   ├── accounts_password-reset_verify.yaml
│   │   ├── accounts_register.yaml
│   │   ├── accounts_{id}.yaml
│   │   ├── artifacts.yaml
│   │   ├── artifacts_{id}.yaml
│   │   ├── config.yaml
│   │   ├── health.yaml
│   │   ├── invitations.yaml
│   │   ├── invitations_{id}.yaml
│   │   ├── labels.yaml
│   │   ├── labels_{id}.yaml
│   │   ├── members.yaml
│   │   ├── members_{id}.yaml
│   │   ├── oauth_authorize.yaml
│   │   ├── oauth_callback.yaml
│   │   ├── organizations.yaml
│   │   ├── organizations_{id}.yaml
│   │   ├── pipelines.yaml
│   │   ├── pipelines_{id}.yaml
│   │   ├── pipelines_{id}_execution-plans.yaml
│   │   ├── pipelines_{id}_steps.yaml
│   │   ├── runs.yaml
│   │   ├── runs_{id}.yaml
│   │   ├── sessions.yaml
│   │   ├── sessions_create.yaml
│   │   ├── sessions_delete.yaml
│   │   ├── sessions_{id}.yaml
│   │   ├── tokens.yaml
│   │   ├── tokens_{id}.yaml
│   │   ├── tools.yaml
│   │   ├── tools_{id}.yaml
│   │   ├── users.yaml
│   │   └── users_{id}.yaml
│   └── openapi.yaml
├── assets
│   ├── github-hero.png
│   ├── large-logo.svg
│   └── small-logo.svg
├── cmd
│   └── archesai
│       └── main.go
├── deployments
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
│   │   │   ├── files
│   │   │   │   └── certs
│   │   │   │       └── .gitkeep
│   │   │   ├── templates
│   │   │   │   ├── components
│   │   │   │   │   ├── api
│   │   │   │   │   │   ├── deployment.yaml
│   │   │   │   │   │   └── service.yaml
│   │   │   │   │   ├── database
│   │   │   │   │   │   ├── deployment.yaml
│   │   │   │   │   │   ├── pvc.yaml
│   │   │   │   │   │   └── service.yaml
│   │   │   │   │   ├── ingress
│   │   │   │   │   │   ├── cert-issuer.yaml
│   │   │   │   │   │   ├── ingress.yaml
│   │   │   │   │   │   └── tls-secret.yaml
│   │   │   │   │   ├── migrations
│   │   │   │   │   │   └── migrations.yaml
│   │   │   │   │   ├── monitoring
│   │   │   │   │   │   ├── grafana-deployment.yaml
│   │   │   │   │   │   ├── grafana-service.yaml
│   │   │   │   │   │   ├── loki-deployment.yaml
│   │   │   │   │   │   └── loki-service.yaml
│   │   │   │   │   ├── platform
│   │   │   │   │   │   ├── deployment.yaml
│   │   │   │   │   │   └── service.yaml
│   │   │   │   │   ├── redis
│   │   │   │   │   │   ├── deployment.yaml
│   │   │   │   │   │   ├── pvc.yaml
│   │   │   │   │   │   └── service.yaml
│   │   │   │   │   ├── scraper
│   │   │   │   │   │   ├── deployment.yaml
│   │   │   │   │   │   └── service.yaml
│   │   │   │   │   ├── storage
│   │   │   │   │   │   ├── deployment.yaml
│   │   │   │   │   │   ├── pvc.yaml
│   │   │   │   │   │   └── service.yaml
│   │   │   │   │   └── unstructured
│   │   │   │   │       ├── deployment.yaml
│   │   │   │   │       └── service.yaml
│   │   │   │   ├── _helpers.tpl
│   │   │   │   ├── configmap.yaml
│   │   │   │   ├── namespace.yaml
│   │   │   │   ├── secrets.yaml
│   │   │   │   └── serviceaccount.yaml
│   │   │   ├── Chart.yaml
│   │   │   ├── values.schema.json
│   │   │   └── values.yaml
│   │   └── dev-overrides.yaml
│   ├── helm-minimal
│   │   ├── files
│   │   │   └── certs
│   │   │       └── .gitkeep
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
├── docs
│   ├── api-reference
│   │   └── overview.md
│   ├── architecture
│   │   ├── authentication.md
│   │   ├── overview.md
│   │   ├── project-layout.md
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
│   │   ├── password.go
│   │   └── service_test.go
│   ├── app
│   │   ├── app.go
│   │   ├── infrastructure.go
│   │   └── routes.go
│   ├── artifacts
│   │   ├── artifacts.go
│   │   ├── errors.go
│   │   └── service_test.go
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
│   │   │   ├── repository.go.tmpl
│   │   │   ├── repository_postgres.go.tmpl
│   │   │   ├── repository_sqlite.go.tmpl
│   │   │   ├── server.gen.go.tmpl
│   │   │   └── service.go.tmpl
│   │   ├── codegen.go
│   │   ├── defaults.go
│   │   ├── filewriter.go
│   │   ├── helpers.go
│   │   ├── parser.go
│   │   ├── sql_generator.go
│   │   └── template_funcs.go
│   ├── config
│   │   ├── config.go
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
│   │   │   ├── pipeline-step-dependencies.sql.go
│   │   │   ├── pipeline-steps.sql.go
│   │   │   ├── pipelines.sql.go
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
│   │   │   ├── pipeline-step-dependencies.sql
│   │   │   ├── pipeline-steps.sql
│   │   │   ├── pipelines.sql
│   │   │   ├── runs.sql
│   │   │   ├── sessions.sql
│   │   │   ├── tools.sql
│   │   │   ├── users.sql
│   │   │   └── verification-tokens.sql
│   │   ├── sqlite
│   │   │   └── stub.go
│   │   ├── database.go
│   │   ├── factory.go
│   │   ├── postgres.go
│   │   ├── sqlc.yaml
│   │   └── sqlite.go
│   ├── events
│   │   ├── events.go
│   │   ├── noop.go
│   │   ├── publisher.go
│   │   └── redis.go
│   ├── health
│   │   ├── errors.go
│   │   ├── handler.go
│   │   ├── health.go
│   │   ├── service.go
│   │   └── service_test.go
│   ├── invitations
│   │   ├── errors.go
│   │   └── invitations.go
│   ├── labels
│   │   ├── errors.go
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
│   │   └── members.go
│   ├── middleware
│   │   ├── auth.go
│   │   └── ratelimit.go
│   ├── migrations
│   │   ├── postgresql
│   │   │   ├── 20250908082709_init.sql
│   │   │   ├── 20250908124238_enable_pgvector.sql
│   │   │   └── 20250915075921_update_api_token_structure.sql
│   │   ├── sqlite
│   │   │   ├── 20250908082709_init.sql
│   │   │   ├── 20250908124238_enable_pgvector.sql
│   │   │   └── 20250915075921_update_api_token_structure.sql
│   │   └── migrate.go
│   ├── oauth
│   │   ├── github.go
│   │   ├── github_test.go
│   │   ├── google.go
│   │   ├── google_test.go
│   │   ├── microsoft.go
│   │   ├── microsoft_test.go
│   │   ├── oauth.go
│   │   └── provider.go
│   ├── organizations
│   │   ├── errors.go
│   │   ├── organizations.go
│   │   └── service_test.go
│   ├── pipelines
│   │   ├── dag.go
│   │   ├── dag_test.go
│   │   ├── errors.go
│   │   ├── executor.go
│   │   ├── executor_test.go
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
│   │   └── runs.go
│   ├── server
│   │   ├── assets
│   │   │   └── docs.html
│   │   ├── docs.go
│   │   ├── infra.go
│   │   ├── middleware.go
│   │   ├── server.go
│   │   └── websocket.go
│   ├── sessions
│   │   ├── claims.go
│   │   ├── claims_test.go
│   │   ├── context.go
│   │   ├── errors.go
│   │   ├── permissions.go
│   │   ├── permissions_test.go
│   │   ├── service.go
│   │   ├── service_test.go
│   │   ├── sessionmanager.go
│   │   ├── sessionmanager_test.go
│   │   ├── sessions.go
│   │   └── tokens.go
│   ├── storage
│   │   ├── storage.go
│   │   └── storage_test.go
│   ├── testutil
│   │   └── containers.go
│   ├── tokens
│   │   ├── errors.go
│   │   └── tokens.go
│   ├── tools
│   │   ├── errors.go
│   │   └── tools.go
│   ├── tui
│   │   ├── config_tui.go
│   │   └── tui.go
│   └── users
│       ├── errors.go
│       ├── service_test.go
│       └── users.go
├── scripts
│   ├── add-mapstructure-tags.sh
│   ├── generate-coverage-report.sh
│   └── generate-helm-schema.py
├── test
│   └── data
│       ├── book.pdf
│       ├── pdf.png
│       ├── text.png
│       └── website.png
├── tools
│   ├── codegen
│   │   └── main.go
│   └── pg-to-sqlite
│       └── main.go
├── web
│   ├── client
│   │   ├── src
│   │   │   ├── generated
│   │   │   │   ├── accounts
│   │   │   │   │   └── accounts.ts
│   │   │   │   ├── artifacts
│   │   │   │   │   └── artifacts.ts
│   │   │   │   ├── auth
│   │   │   │   │   └── auth.ts
│   │   │   │   ├── config
│   │   │   │   │   └── config.ts
│   │   │   │   ├── health
│   │   │   │   │   └── health.ts
│   │   │   │   ├── invitations
│   │   │   │   │   └── invitations.ts
│   │   │   │   ├── labels
│   │   │   │   │   └── labels.ts
│   │   │   │   ├── members
│   │   │   │   │   └── members.ts
│   │   │   │   ├── oauth
│   │   │   │   │   └── oauth.ts
│   │   │   │   ├── organizations
│   │   │   │   │   └── organizations.ts
│   │   │   │   ├── pipelines
│   │   │   │   │   └── pipelines.ts
│   │   │   │   ├── runs
│   │   │   │   │   └── runs.ts
│   │   │   │   ├── sessions
│   │   │   │   │   └── sessions.ts
│   │   │   │   ├── tokens
│   │   │   │   │   └── tokens.ts
│   │   │   │   ├── tools
│   │   │   │   │   └── tools.ts
│   │   │   │   ├── users
│   │   │   │   │   └── users.ts
│   │   │   │   └── orval.schemas.ts
│   │   │   ├── fetcher.ts
│   │   │   └── index.ts
│   │   ├── eslint.config.js
│   │   ├── orval.config.ts
│   │   ├── package.json
│   │   ├── tsconfig.json
│   │   ├── tsconfig.lib.json
│   │   └── tsconfig.spec.json
│   ├── docs
│   │   ├── public
│   │   │   ├── pagefind
│   │   │   │   └── pagefind.js
│   │   │   ├── preview
│   │   │   │   ├── dark-api.svg
│   │   │   │   └── dark-portal.svg
│   │   │   ├── search
│   │   │   │   └── search.svg
│   │   │   ├── background.png
│   │   │   ├── background.svg
│   │   │   ├── banner-dark.svg
│   │   │   └── banner.svg
│   │   ├── src
│   │   │   └── sidebar.tsx
│   │   ├── eslint.config.js
│   │   ├── package.json
│   │   ├── tsconfig.app.json
│   │   ├── tsconfig.json
│   │   ├── tsconfig.spec.json
│   │   ├── vite.config.ts
│   │   └── zudoku.config.tsx
│   ├── eslint
│   │   ├── src
│   │   │   ├── base.js
│   │   │   └── react.js
│   │   └── package.json
│   ├── platform
│   │   ├── public
│   │   │   ├── android-chrome-192x192.png
│   │   │   ├── android-chrome-512x512.png
│   │   │   ├── apple-touch-icon.png
│   │   │   ├── favicon-16x16.png
│   │   │   ├── favicon-32x32.png
│   │   │   ├── favicon.ico
│   │   │   └── site.webmanifest
│   │   ├── src
│   │   │   ├── app
│   │   │   │   ├── _app
│   │   │   │   │   ├── artifacts
│   │   │   │   │   │   ├── $artifactID
│   │   │   │   │   │   │   ├── -details.tsx
│   │   │   │   │   │   │   └── index.tsx
│   │   │   │   │   │   └── index.tsx
│   │   │   │   │   ├── labels
│   │   │   │   │   │   └── index.tsx
│   │   │   │   │   ├── organization
│   │   │   │   │   │   ├── members
│   │   │   │   │   │   │   └── index.tsx
│   │   │   │   │   │   └── index.tsx
│   │   │   │   │   ├── pipelines
│   │   │   │   │   │   ├── $pipelineID
│   │   │   │   │   │   │   └── index.tsx
│   │   │   │   │   │   ├── create
│   │   │   │   │   │   │   └── index.lazy.tsx
│   │   │   │   │   │   └── index.tsx
│   │   │   │   │   ├── profile
│   │   │   │   │   │   └── index.tsx
│   │   │   │   │   ├── runs
│   │   │   │   │   │   ├── $runID
│   │   │   │   │   │   │   └── index.tsx
│   │   │   │   │   │   └── index.tsx
│   │   │   │   │   ├── tools
│   │   │   │   │   │   └── index.tsx
│   │   │   │   │   ├── index.tsx
│   │   │   │   │   └── route.tsx
│   │   │   │   ├── auth
│   │   │   │   │   ├── forgot-password
│   │   │   │   │   │   └── index.tsx
│   │   │   │   │   ├── login
│   │   │   │   │   │   └── index.tsx
│   │   │   │   │   ├── register
│   │   │   │   │   │   └── index.tsx
│   │   │   │   │   └── route.tsx
│   │   │   │   ├── landing
│   │   │   │   │   ├── -content.ts
│   │   │   │   │   └── index.tsx
│   │   │   │   └── __root.tsx
│   │   │   ├── components
│   │   │   │   ├── containers
│   │   │   │   │   ├── app-sidebar-container.tsx
│   │   │   │   │   └── page-header-container.tsx
│   │   │   │   ├── datatables
│   │   │   │   │   ├── artifact-datatable.tsx
│   │   │   │   │   ├── data-table-container.tsx
│   │   │   │   │   ├── label-datatable.tsx
│   │   │   │   │   ├── member-datatable.tsx
│   │   │   │   │   ├── pipeline-datatable.tsx
│   │   │   │   │   ├── run-datatable.tsx
│   │   │   │   │   └── tool-datatable.tsx
│   │   │   │   ├── forms
│   │   │   │   │   ├── artifact-form.tsx
│   │   │   │   │   ├── label-form.tsx
│   │   │   │   │   ├── member-form.tsx
│   │   │   │   │   ├── organization-form.tsx
│   │   │   │   │   └── user-form.tsx
│   │   │   │   ├── navigation
│   │   │   │   │   ├── link.tsx
│   │   │   │   │   └── user-button-container.tsx
│   │   │   │   ├── selectors
│   │   │   │   │   └── data-selector-container.tsx
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
│   │   │   │   ├── get-headers.ts
│   │   │   │   ├── site-config.ts
│   │   │   │   └── site-utils.ts
│   │   │   ├── styles
│   │   │   │   └── globals.css
│   │   │   ├── routeTree.gen.ts
│   │   │   └── router.tsx
│   │   ├── Dockerfile
│   │   ├── eslint.config.js
│   │   ├── package.json
│   │   ├── playwright.config.js
│   │   ├── tsconfig.app.json
│   │   ├── tsconfig.json
│   │   ├── tsconfig.spec.json
│   │   └── vite.config.js
│   ├── typescript
│   │   ├── src
│   │   │   ├── base.json
│   │   │   ├── lib.json
│   │   │   ├── react.json
│   │   │   └── spec.json
│   │   └── package.json
│   └── ui
│       ├── src
│       │   ├── components
│       │   │   ├── custom
│       │   │   │   ├── Callout.tsx
│       │   │   │   ├── Introduction.tsx
│       │   │   │   ├── LandingPage.tsx
│       │   │   │   ├── ThemeEditor.tsx
│       │   │   │   ├── arches-logo.tsx
│       │   │   │   ├── artifact-viewer.tsx
│       │   │   │   ├── delete-items.tsx
│       │   │   │   ├── discord-icon.tsx
│       │   │   │   ├── faceted.tsx
│       │   │   │   ├── floating-label.tsx
│       │   │   │   ├── generic-form.tsx
│       │   │   │   ├── github-icon.tsx
│       │   │   │   ├── icons.tsx
│       │   │   │   ├── import-card.tsx
│       │   │   │   ├── pure-data-selector.tsx
│       │   │   │   ├── pure-user-button.tsx
│       │   │   │   ├── run-status-button.tsx
│       │   │   │   ├── schema-builder.tsx
│       │   │   │   ├── scroll-button.tsx
│       │   │   │   ├── sortable.tsx
│       │   │   │   └── timestamp.tsx
│       │   │   ├── datatable
│       │   │   │   ├── components
│       │   │   │   │   ├── filters
│       │   │   │   │   │   ├── data-table-date-filter.tsx
│       │   │   │   │   │   ├── data-table-faceted-filter.tsx
│       │   │   │   │   │   ├── data-table-range-filter.tsx
│       │   │   │   │   │   └── data-table-slider-filter.tsx
│       │   │   │   │   ├── toolbar
│       │   │   │   │   │   ├── data-table-filter-menu.tsx
│       │   │   │   │   │   ├── data-table-sort-list.tsx
│       │   │   │   │   │   └── data-table-toolbar.tsx
│       │   │   │   │   ├── views
│       │   │   │   │   │   ├── grid-view.tsx
│       │   │   │   │   │   └── table-view.tsx
│       │   │   │   │   ├── data-table-action-bar.tsx
│       │   │   │   │   ├── data-table-column-header.tsx
│       │   │   │   │   ├── data-table-pagination.tsx
│       │   │   │   │   ├── data-table-row-actions.tsx
│       │   │   │   │   ├── data-table-view-options.tsx
│       │   │   │   │   ├── tasks-table-action-bar.tsx
│       │   │   │   │   └── view-toggle.tsx
│       │   │   │   ├── data-table.spec.tsx
│       │   │   │   └── pure-data-table.tsx
│       │   │   ├── primitives
│       │   │   │   └── link.tsx
│       │   │   ├── shadcn
│       │   │   │   ├── accordion.tsx
│       │   │   │   ├── alert-dialog.tsx
│       │   │   │   ├── alert.tsx
│       │   │   │   ├── aspect-ratio.tsx
│       │   │   │   ├── avatar.tsx
│       │   │   │   ├── badge.tsx
│       │   │   │   ├── breadcrumb.tsx
│       │   │   │   ├── button.tsx
│       │   │   │   ├── calendar.tsx
│       │   │   │   ├── card.tsx
│       │   │   │   ├── carousel.tsx
│       │   │   │   ├── checkbox.tsx
│       │   │   │   ├── collapsible.tsx
│       │   │   │   ├── command.tsx
│       │   │   │   ├── context-menu.tsx
│       │   │   │   ├── dialog.tsx
│       │   │   │   ├── drawer.tsx
│       │   │   │   ├── dropdown-menu.tsx
│       │   │   │   ├── form.tsx
│       │   │   │   ├── hover-card.tsx
│       │   │   │   ├── input-otp.tsx
│       │   │   │   ├── input.tsx
│       │   │   │   ├── label.tsx
│       │   │   │   ├── menubar.tsx
│       │   │   │   ├── navigation-menu.tsx
│       │   │   │   ├── pagination.tsx
│       │   │   │   ├── popover.tsx
│       │   │   │   ├── progress.tsx
│       │   │   │   ├── radio-group.tsx
│       │   │   │   ├── resizable.tsx
│       │   │   │   ├── scroll-area.tsx
│       │   │   │   ├── select.tsx
│       │   │   │   ├── separator.tsx
│       │   │   │   ├── sheet.tsx
│       │   │   │   ├── sidebar.tsx
│       │   │   │   ├── skeleton.tsx
│       │   │   │   ├── slider.tsx
│       │   │   │   ├── sonner.tsx
│       │   │   │   ├── switch.tsx
│       │   │   │   ├── table.tsx
│       │   │   │   ├── tabs.tsx
│       │   │   │   ├── textarea.tsx
│       │   │   │   ├── toggle-group.tsx
│       │   │   │   ├── toggle.tsx
│       │   │   │   └── tooltip.tsx
│       │   │   └── index.ts
│       │   ├── hooks
│       │   │   ├── use-callback-ref.tsx
│       │   │   ├── use-debounced-callback.tsx
│       │   │   ├── use-is-top.tsx
│       │   │   └── use-mobile.tsx
│       │   ├── layouts
│       │   │   ├── app-sidebar
│       │   │   │   ├── app-sidebar.tsx
│       │   │   │   ├── organization-button.tsx
│       │   │   │   └── sidebar-links.tsx
│       │   │   ├── page-header
│       │   │   │   ├── components
│       │   │   │   │   ├── breadcrumbs.tsx
│       │   │   │   │   ├── command-menu.tsx
│       │   │   │   │   ├── email-verify.tsx
│       │   │   │   │   ├── theme-toggle.tsx
│       │   │   │   │   └── title-and-description.tsx
│       │   │   │   └── page-header.tsx
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
│       ├── eslint.config.js
│       ├── package.json
│       ├── tsconfig.json
│       ├── tsconfig.lib.json
│       ├── tsconfig.spec.json
│       └── vite.config.ts
├── .air.toml
├── .codegen.archesai.yaml
├── .codegen.server.yaml
├── .codegen.types.yaml
├── .cspell.json
├── .editorconfig
├── .gitignore
├── .golangci.yaml
├── .goreleaser.yaml
├── .lefthook.yaml
├── .markdownlint.json
├── .mockery.yaml
├── .prettierignore
├── .redocly.yaml
├── LICENSE
├── Makefile
├── README.md
├── biome.json
├── go.mod
├── go.sum
├── package.json
├── pnpm-lock.yaml
├── pnpm-workspace.yaml
└── tsconfig.json

198 directories, 725 files
```

## Domain Package Structure

Each domain in `/internal` follows this pattern:

```text
domain/
├── domain.go          # Package documentation, constants, errors
├── service.go         # Business logic implementation
├── handler.go         # HTTP request/response handling
├── middleware.go      # Domain-specific middleware (optional)
├── repository.gen.go  # Generated repository interface
├── postgres.gen.go    # Generated PostgreSQL implementation
├── sqlite.gen.go      # Generated SQLite implementation
├── service.gen.go     # Generated service interface
├── server.gen.go      # Generated HTTP server implementation
├── types.gen.go       # Generated types from OpenAPI
├── api.gen.go         # Generated API client interface
├── service_test.go    # Unit tests with mocked dependencies
├── handler_test.go    # HTTP handler tests (optional)
├── mocks_test.gen.go  # Generated test mocks
└── postgres_test.go   # Integration tests (optional)
```

## Generated Files

### Go Generated Files

- `*.gen.go` - Do not edit manually
- `types.gen.go` - OpenAPI struct definitions
- `api.gen.go` - API client interfaces
- `repository.gen.go` - Repository interface from x-codegen
- `postgres.gen.go` - PostgreSQL repository implementation
- `sqlite.gen.go` - SQLite repository implementation
- `service.gen.go` - Service interface from x-codegen
- `server.gen.go` - HTTP server implementation
- `mocks_test.gen.go` - Test mocks from mockery

### TypeScript Generated Files

- `web/client/src/generated/` - Complete API client
- Generated from `api/openapi.bundled.yaml`

### SQL Generated Files

- Database queries in `internal/database/queries/*.sql`
- Generate Go code with `sqlc generate`

## File Naming Conventions

### Go Files

- `domain.go` - Package documentation, constants, errors
- `service.go` - Business logic implementation
- `handler.go` - HTTP handlers
- `middleware.go` - Middleware functions (optional)
- `*_test.go` - Test files
- `*.gen.go` - Generated files (don't edit manually)
- `mocks_test.gen.go` - Generated test mocks

### Config Files

- `.*.yaml` - YAML configs (golangci, mockery, sqlc)
- `.*.toml` - TOML configs (air)
- `.*.json` - JSON configs (markdownlint, tsconfig)
- `.*rc` - RC files (prettierrc)

### Documentation

- `*.md` - Markdown docs
- `*.mdx` - MDX with components (web/docs)
- `README.md` - Package/directory docs
