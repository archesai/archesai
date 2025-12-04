# Project Layout

## Directory Structure

```text
.
├── api
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
│   │   │   │   ├── code-generation.md
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
│   │   │   │   ├── quickstart.md
│   │   │   │   ├── test-coverage-report.md
│   │   │   │   └── testing.md
│   │   │   ├── security
│   │   │   │   ├── best-practices.md
│   │   │   │   └── overview.md
│   │   │   ├── troubleshooting
│   │   │   │   └── common-issues.md
│   │   │   ├── cli-reference.md
│   │   │   ├── contributing.md
│   │   │   ├── getting-started.md
│   │   │   └── ROADMAP.md
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
│   └── studio
│       ├── bootstrap
│       ├── infrastructure
│       │   ├── postgres
│       │   │   ├── migrations
│       │   │   ├── queries
│       │   │   ├── repositories
│       │   │   ├── schema.gen.hcl
│       │   │   └── sqlc.gen.yaml
│       │   └── sqlite
│       │       ├── migrations
│       │       ├── queries
│       │       ├── repositories
│       │       ├── schema.gen.hcl
│       │       └── sqlc.gen.yaml
│       ├── models
│       ├── go.mod
│       └── go.sum
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
│   │   ├── code-generation.md
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
│   │   ├── quickstart.md
│   │   ├── test-coverage-report.md
│   │   └── testing.md
│   ├── security
│   │   ├── best-practices.md
│   │   └── overview.md
│   ├── troubleshooting
│   │   └── common-issues.md
│   ├── cli-reference.md
│   ├── contributing.md
│   ├── getting-started.md
│   └── ROADMAP.md
├── examples
│   ├── auth
│   │   ├── infrastructure
│   │   │   ├── postgres
│   │   │   │   ├── migrations
│   │   │   │   ├── queries
│   │   │   │   ├── repositories
│   │   │   │   ├── schema.gen.hcl
│   │   │   │   └── sqlc.gen.yaml
│   │   │   └── sqlite
│   │   │       ├── migrations
│   │   │       ├── queries
│   │   │       ├── repositories
│   │   │       └── schema.gen.hcl
│   │   ├── models
│   │   ├── spec
│   │   │   ├── openapi.bundled.yaml
│   │   │   └── openapi.yaml
│   │   ├── go.mod
│   │   └── go.sum
│   └── basic
│       ├── application
│       ├── bootstrap
│       ├── controllers
│       ├── infrastructure
│       │   ├── postgres
│       │   │   ├── migrations
│       │   │   ├── queries
│       │   │   ├── schema.gen.hcl
│       │   │   └── sqlc.gen.yaml
│       │   └── sqlite
│       │       ├── migrations
│       │       ├── queries
│       │       ├── repositories
│       │       └── schema.gen.hcl
│       ├── models
│       ├── spec
│       │   ├── openapi.bundled.yaml
│       │   └── openapi.yaml
│       └── go.mod
├── internal
│   ├── cli
│   │   ├── completion.go
│   │   ├── config.go
│   │   ├── dev.go
│   │   ├── generate.go
│   │   ├── root.go
│   │   └── version.go
│   ├── codegen
│   │   ├── tmpl
│   │   │   ├── app.go.tmpl
│   │   │   ├── container.go.tmpl
│   │   │   ├── go.mod.tmpl
│   │   │   ├── handler_controller.go.tmpl
│   │   │   ├── handler.gen.go.tmpl
│   │   │   ├── handler_stub.go.tmpl
│   │   │   ├── hcl.tmpl
│   │   │   ├── header.tmpl
│   │   │   ├── main.go.tmpl
│   │   │   ├── repository.go.tmpl
│   │   │   ├── repository_postgres.go.tmpl
│   │   │   ├── repository_sqlite.go.tmpl
│   │   │   ├── routes.go.tmpl
│   │   │   ├── schema.go.tmpl
│   │   │   ├── sqlc_postgres.yaml.tmpl
│   │   │   ├── sql_queries.sql.tmpl
│   │   │   └── wire.go.tmpl
│   │   ├── generate_app.go
│   │   ├── generate_container.go
│   │   ├── generate.go
│   │   ├── generate_gomod.go
│   │   ├── generate_handler_controllers.go
│   │   ├── generate_handler_stubs.go
│   │   ├── generate_hcl.go
│   │   ├── generate_js_client.go
│   │   ├── generate_migrations.go
│   │   ├── generate_repositories.go
│   │   ├── generate_routes.go
│   │   ├── generate_schemas.go
│   │   ├── generate_sqlc.go
│   │   ├── generate_wire.go
│   │   ├── renderer.go
│   │   └── templates.go
│   ├── dev
│   │   ├── manager.go
│   │   ├── process.go
│   │   └── watcher.go
│   ├── parsers
│   │   ├── handlers.go
│   │   ├── handlers_test.go
│   │   ├── jsonschema.go
│   │   ├── jsonschema_test.go
│   │   ├── linter.go
│   │   ├── linter_test.go
│   │   ├── openapi.go
│   │   ├── openapi_includes.go
│   │   ├── openapi_includes_registry.go
│   │   ├── openapi_includes_test.go
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
│   │   ├── application
│   │   │   ├── confirm_email_change.impl.go
│   │   │   ├── confirm_email_verification.impl.go
│   │   │   ├── confirm_password_reset.impl.go
│   │   │   ├── delete_account.impl.go
│   │   │   ├── delete_current_user.impl.go
│   │   │   ├── delete_session.impl.go
│   │   │   ├── get_current_user.impl.go
│   │   │   ├── link_account.impl.go
│   │   │   ├── login.impl.go
│   │   │   ├── logout_all.impl.go
│   │   │   ├── logout.impl.go
│   │   │   ├── oauth_authorize.impl.go
│   │   │   ├── oauth_callback.impl.go
│   │   │   ├── register.impl.go
│   │   │   ├── request_email_change.impl.go
│   │   │   ├── request_email_verification.impl.go
│   │   │   ├── request_magic_link.impl.go
│   │   │   ├── request_password_reset.impl.go
│   │   │   ├── update_account.impl.go
│   │   │   ├── update_current_user.impl.go
│   │   │   ├── update_session.impl.go
│   │   │   └── verify_magic_link.impl.go
│   │   ├── bootstrap
│   │   ├── controllers
│   │   ├── models
│   │   ├── oauth
│   │   │   ├── github.go
│   │   │   ├── google.go
│   │   │   ├── microsoft.go
│   │   │   └── oauth.go
│   │   ├── repositories
│   │   ├── spec
│   │   │   ├── components
│   │   │   │   ├── headers
│   │   │   │   ├── parameters
│   │   │   │   ├── responses
│   │   │   │   └── schemas
│   │   │   ├── paths
│   │   │   │   ├── api-keys_id.yaml
│   │   │   │   ├── api-keys.yaml
│   │   │   │   ├── auth_accounts_id.yaml
│   │   │   │   ├── auth_accounts.yaml
│   │   │   │   ├── auth_change-email.yaml
│   │   │   │   ├── auth_confirm-email.yaml
│   │   │   │   ├── auth_forgot-password.yaml
│   │   │   │   ├── auth_link.yaml
│   │   │   │   ├── auth_login.yaml
│   │   │   │   ├── auth_logout-all.yaml
│   │   │   │   ├── auth_logout.yaml
│   │   │   │   ├── auth_magic-links_request.yaml
│   │   │   │   ├── auth_magic-links_verify.yaml
│   │   │   │   ├── auth_me.yaml
│   │   │   │   ├── auth_oauth_provider_authorize.yaml
│   │   │   │   ├── auth_oauth_provider_callback.yaml
│   │   │   │   ├── auth_register.yaml
│   │   │   │   ├── auth_request-verification.yaml
│   │   │   │   ├── auth_reset-password.yaml
│   │   │   │   ├── auth_sessions_id.yaml
│   │   │   │   ├── auth_sessions.yaml
│   │   │   │   ├── auth_verify-email.yaml
│   │   │   │   ├── organizations_id.yaml
│   │   │   │   ├── organizations_organizationID_invitations_id.yaml
│   │   │   │   ├── organizations_organizationID_invitations.yaml
│   │   │   │   ├── organizations_organizationID_members_id.yaml
│   │   │   │   ├── organizations_organizationID_members.yaml
│   │   │   │   ├── organizations.yaml
│   │   │   │   ├── users_id.yaml
│   │   │   │   └── users.yaml
│   │   │   ├── openapi.bundled.yaml
│   │   │   └── openapi.yaml
│   │   ├── auth.go
│   │   ├── auth_tokens.go
│   │   ├── magic_link.go
│   │   ├── password.go
│   │   ├── service.go
│   │   ├── spec.go
│   │   └── token_manager.go
│   ├── cache
│   │   ├── cache.go
│   │   ├── memory.go
│   │   ├── noop.go
│   │   └── redis.go
│   ├── config
│   │   ├── application
│   │   │   └── get_config.impl.go
│   │   ├── bootstrap
│   │   ├── controllers
│   │   ├── models
│   │   ├── spec
│   │   │   ├── components
│   │   │   │   ├── responses
│   │   │   │   └── schemas
│   │   │   ├── paths
│   │   │   │   └── config.yaml
│   │   │   ├── openapi.bundled.yaml
│   │   │   └── openapi.yaml
│   │   ├── config.go
│   │   ├── loader.go
│   │   ├── loader_test.go
│   │   └── spec.go
│   ├── database
│   │   ├── database.go
│   │   ├── migrate.go
│   │   └── repository.go
│   ├── events
│   │   ├── events.go
│   │   ├── noop.go
│   │   ├── publisher.go
│   │   └── redis.go
│   ├── executor
│   │   ├── application
│   │   │   └── execute_executor.impl.go
│   │   ├── bootstrap
│   │   ├── controllers
│   │   ├── models
│   │   ├── repositories
│   │   ├── spec
│   │   │   ├── components
│   │   │   │   ├── parameters
│   │   │   │   ├── responses
│   │   │   │   └── schemas
│   │   │   ├── paths
│   │   │   │   ├── executors_id_execute.yaml
│   │   │   │   ├── executors_id.yaml
│   │   │   │   └── executors.yaml
│   │   │   ├── openapi.bundled.yaml
│   │   │   └── openapi.yaml
│   │   ├── testdata
│   │   │   └── execute.ts
│   │   ├── builder.go
│   │   ├── builder_test.go
│   │   ├── config.go
│   │   ├── container.go
│   │   ├── container_test.go
│   │   ├── executor.go
│   │   ├── local.go
│   │   ├── local_test.go
│   │   ├── ports.go
│   │   ├── schemas.go
│   │   ├── service.go
│   │   └── spec.go
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
│   │   ├── deliverer.go
│   │   ├── email.go
│   │   └── otp.go
│   ├── pipelines
│   │   ├── application
│   │   │   ├── create_pipeline_step.impl.go
│   │   │   ├── get_pipeline_execution_plan.impl.go
│   │   │   ├── get_pipeline_steps.impl.go
│   │   │   └── validate_pipeline_execution_plan.impl.go
│   │   ├── bootstrap
│   │   ├── controllers
│   │   ├── models
│   │   ├── repositories
│   │   ├── spec
│   │   │   ├── components
│   │   │   │   ├── parameters
│   │   │   │   ├── responses
│   │   │   │   └── schemas
│   │   │   ├── paths
│   │   │   │   ├── pipelines_id_execution-plans.yaml
│   │   │   │   ├── pipelines_id_steps.yaml
│   │   │   │   ├── pipelines_id.yaml
│   │   │   │   ├── pipelines.yaml
│   │   │   │   ├── runs_id.yaml
│   │   │   │   ├── runs.yaml
│   │   │   │   ├── tools_id.yaml
│   │   │   │   └── tools.yaml
│   │   │   ├── openapi.bundled.yaml
│   │   │   └── openapi.yaml
│   │   └── spec.go
│   ├── redis
│   │   ├── client.go
│   │   ├── config.go
│   │   ├── errors.go
│   │   ├── pubsub.go
│   │   ├── queue.go
│   │   └── redis.go
│   ├── server
│   │   ├── application
│   │   │   └── get_health.impl.go
│   │   ├── bootstrap
│   │   ├── controllers
│   │   ├── models
│   │   ├── spec
│   │   │   ├── components
│   │   │   │   ├── headers
│   │   │   │   ├── parameters
│   │   │   │   ├── responses
│   │   │   │   └── schemas
│   │   │   ├── paths
│   │   │   │   └── health.yaml
│   │   │   ├── openapi.bundled.yaml
│   │   │   └── openapi.yaml
│   │   ├── cookies.go
│   │   ├── middleware_auth.go
│   │   ├── middleware_cors.go
│   │   ├── middleware.go
│   │   ├── middleware_logger.go
│   │   ├── middleware_ratelimit.go
│   │   ├── middleware_recover.go
│   │   ├── middleware_requestid.go
│   │   ├── middleware_security.go
│   │   ├── middleware_timeout.go
│   │   ├── responses.go
│   │   ├── server.go
│   │   ├── spec.go
│   │   └── websocket.go
│   └── storage
│       ├── application
│       ├── bootstrap
│       ├── controllers
│       ├── models
│       ├── repositories
│       ├── spec
│       │   ├── components
│       │   │   ├── parameters
│       │   │   ├── responses
│       │   │   └── schemas
│       │   ├── paths
│       │   │   ├── artifacts_id.yaml
│       │   │   ├── artifacts.yaml
│       │   │   ├── labels_id.yaml
│       │   │   └── labels.yaml
│       │   ├── openapi.bundled.yaml
│       │   └── openapi.yaml
│       ├── disk.go
│       ├── memory.go
│       ├── spec.go
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
├── .vscode
│   ├── extensions.json
│   └── settings.json
├── web
│   ├── client
│   │   ├── src
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
├── coverage.txt
├── .cspell.json
├── .editorconfig
├── .env
├── .gitignore
├── .golangci.yaml
├── go.mod
├── .goreleaser.yaml
├── go.sum
├── go.work
├── go.work.sum
├── .lefthook.yaml
├── LICENSE
├── Makefile
├── .markdownlint.json
├── .mockery.yaml
├── package.json
├── pnpm-lock.yaml
├── pnpm-workspace.yaml
├── .prettierignore
├── README.md
├── tools.mod
├── tools.sum
└── tsconfig.json

240 directories, 524 files
```
