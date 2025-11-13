table "account" {
  schema = schema.public

  column "id" {
    null    = false
    type    = sql("uuid")
    default = sql("gen_random_uuid()")
  }

  column "created_at" {
    null    = false
    type    = sql("timestamptz")
    default = sql("CURRENT_TIMESTAMP")
  }

  column "updated_at" {
    null    = false
    type    = sql("timestamptz")
    default = sql("CURRENT_TIMESTAMP")
  }

  column "access_token" {
    null = true
    type = sql("text")
  }

  column "access_token_expires_at" {
    null = true
    type = sql("timestamptz")
  }

  column "account_identifier" {
    null = false
    type = sql("text")
  }

  column "id_token" {
    null = true
    type = sql("text")
  }

  column "provider" {
    null = false
    type = sql("text")
  }

  column "refresh_token" {
    null = true
    type = sql("text")
  }

  column "refresh_token_expires_at" {
    null = true
    type = sql("timestamptz")
  }

  column "scope" {
    null = true
    type = sql("text")
  }

  column "user_id" {
    null = false
    type = sql("uuid")
  }
  primary_key {
    columns = [column.id]
  }
  foreign_key "account_user_id_fkey" {
    columns     = [column.user_id]
    ref_columns = [table.user.column.id]
    on_update   = "NO_ACTION"
    on_delete   = "CASCADE"
  }
  index "idx_account_user_id" {
    columns = [column.user_id]
  }
  check "account_provider_check" {
    expr = "(provider = ANY (ARRAY['local'::text, 'google'::text, 'github'::text, 'microsoft'::text, 'apple'::text]))"
  }
}

table "api_key" {
  schema = schema.public

  column "id" {
    null    = false
    type    = sql("uuid")
    default = sql("gen_random_uuid()")
  }

  column "created_at" {
    null    = false
    type    = sql("timestamptz")
    default = sql("CURRENT_TIMESTAMP")
  }

  column "updated_at" {
    null    = false
    type    = sql("timestamptz")
    default = sql("CURRENT_TIMESTAMP")
  }

  column "expires_at" {
    null = true
    type = sql("timestamptz")
  }

  column "key_hash" {
    null = false
    type = sql("text")
  }

  column "last_used_at" {
    null = true
    type = sql("timestamptz")
  }

  column "name" {
    null = true
    type = sql("text")
  }

  column "organization_id" {
    null = false
    type = sql("uuid")
  }

  column "prefix" {
    null = true
    type = sql("text")
  }

  column "rate_limit" {
    null    = false
    type    = sql("integer")
    default = 60
  }

  column "scopes" {
    null    = false
    type    = sql("text[]")
    default = "{}"
  }

  column "user_id" {
    null = false
    type = sql("uuid")
  }
  primary_key {
    columns = [column.id]
  }
  foreign_key "api_key_organization_id_fkey" {
    columns     = [column.organization_id]
    ref_columns = [table.organization.column.id]
    on_update   = "NO_ACTION"
    on_delete   = "CASCADE"
  }
  foreign_key "api_key_user_id_fkey" {
    columns     = [column.user_id]
    ref_columns = [table.user.column.id]
    on_update   = "NO_ACTION"
    on_delete   = "CASCADE"
  }
  index "idx_api_key_key_hash" {
    columns = [column.key_hash]
  }
  index "idx_api_key_last_used_at" {
    columns = [column.last_used_at]
  }
  index "idx_api_key_organization_id" {
    columns = [column.organization_id]
  }
  index "idx_api_key_user_id" {
    columns = [column.user_id]
  }
}

table "artifact" {
  schema = schema.public

  column "id" {
    null    = false
    type    = sql("uuid")
    default = sql("gen_random_uuid()")
  }

  column "created_at" {
    null    = false
    type    = sql("timestamptz")
    default = sql("CURRENT_TIMESTAMP")
  }

  column "updated_at" {
    null    = false
    type    = sql("timestamptz")
    default = sql("CURRENT_TIMESTAMP")
  }

  column "credits" {
    null    = false
    type    = sql("integer")
    default = 0
  }

  column "description" {
    null = true
    type = sql("text")
  }

  column "mime_type" {
    null    = false
    type    = sql("text")
    default = "application/octet-stream"
  }

  column "name" {
    null = true
    type = sql("text")
  }

  column "organization_id" {
    null = false
    type = sql("uuid")
  }

  column "preview_image" {
    null = true
    type = sql("text")
  }

  column "producer_id" {
    null = true
    type = sql("uuid")
  }

  column "text" {
    null = true
    type = sql("text")
  }

  column "url" {
    null = true
    type = sql("text")
  }
  primary_key {
    columns = [column.id]
  }
  foreign_key "artifact_organization_id_fkey" {
    columns     = [column.organization_id]
    ref_columns = [table.organization.column.id]
    on_update   = "CASCADE"
    on_delete   = "CASCADE"
  }
  foreign_key "artifact_producer_id_fkey" {
    columns     = [column.producer_id]
    ref_columns = [table.run.column.id]
    on_update   = "CASCADE"
    on_delete   = "SET_NULL"
  }
  index "idx_artifact_organization_id" {
    columns = [column.organization_id]
  }
  index "idx_artifact_producer_id" {
    columns = [column.producer_id]
  }
}

table "executor" {
  schema = schema.public

  column "id" {
    null    = false
    type    = sql("uuid")
    default = sql("gen_random_uuid()")
  }

  column "created_at" {
    null    = false
    type    = sql("timestamptz")
    default = sql("CURRENT_TIMESTAMP")
  }

  column "updated_at" {
    null    = false
    type    = sql("timestamptz")
    default = sql("CURRENT_TIMESTAMP")
  }

  column "cpu_shares" {
    null    = false
    type    = sql("integer")
    default = 512
  }

  column "dependencies" {
    null = true
    type = sql("text")
  }

  column "description" {
    null = false
    type = sql("text")
  }

  column "env" {
    null = true
    type = sql("text")
  }

  column "execute_code" {
    null = false
    type = sql("text")
  }

  column "extra_files" {
    null = true
    type = sql("text")
  }

  column "is_active" {
    null    = false
    type    = sql("boolean")
    default = true
  }

  column "language" {
    null = false
    type = sql("text")
  }

  column "memory_mb" {
    null    = false
    type    = sql("integer")
    default = 256
  }

  column "name" {
    null = false
    type = sql("text")
  }

  column "organization_id" {
    null = false
    type = sql("uuid")
  }

  column "schema_in" {
    null = true
    type = sql("text")
  }

  column "schema_out" {
    null = true
    type = sql("text")
  }

  column "timeout" {
    null    = false
    type    = sql("integer")
    default = 30
  }

  column "version" {
    null    = false
    type    = sql("integer")
    default = 1
  }
  primary_key {
    columns = [column.id]
  }
  foreign_key "executor_organization_id_fkey" {
    columns     = [column.organization_id]
    ref_columns = [table.organization.column.id]
    on_update   = "CASCADE"
    on_delete   = "CASCADE"
  }
  index "idx_executor_language" {
    columns = [column.language]
  }
  index "idx_executor_organization_id" {
    columns = [column.organization_id]
  }
  check "executor_language_check" {
    expr = "(language = ANY (ARRAY['nodejs'::text, 'python'::text, 'go'::text]))"
  }
}

table "invitation" {
  schema = schema.public

  column "id" {
    null    = false
    type    = sql("uuid")
    default = sql("gen_random_uuid()")
  }

  column "created_at" {
    null    = false
    type    = sql("timestamptz")
    default = sql("CURRENT_TIMESTAMP")
  }

  column "updated_at" {
    null    = false
    type    = sql("timestamptz")
    default = sql("CURRENT_TIMESTAMP")
  }

  column "email" {
    null = false
    type = sql("text")
  }

  column "expires_at" {
    null = false
    type = sql("timestamptz")
  }

  column "inviter_id" {
    null = false
    type = sql("uuid")
  }

  column "organization_id" {
    null = false
    type = sql("uuid")
  }

  column "role" {
    null    = false
    type    = sql("text")
    default = "basic"
  }

  column "status" {
    null    = false
    type    = sql("text")
    default = "pending"
  }
  primary_key {
    columns = [column.id]
  }
  foreign_key "invitation_inviter_id_fkey" {
    columns     = [column.inviter_id]
    ref_columns = [table.user.column.id]
    on_update   = "NO_ACTION"
    on_delete   = "CASCADE"
  }
  foreign_key "invitation_organization_id_fkey" {
    columns     = [column.organization_id]
    ref_columns = [table.organization.column.id]
    on_update   = "NO_ACTION"
    on_delete   = "CASCADE"
  }
  index "idx_invitation_inviter_id" {
    columns = [column.inviter_id]
  }
  index "idx_invitation_organization_id" {
    columns = [column.organization_id]
  }
  check "invitation_role_check" {
    expr = "(role = ANY (ARRAY['admin'::text, 'owner'::text, 'basic'::text]))"
  }
  check "invitation_status_check" {
    expr = "(status = ANY (ARRAY['pending'::text, 'accepted'::text, 'declined'::text, 'expired'::text]))"
  }
}

table "label" {
  schema = schema.public

  column "id" {
    null    = false
    type    = sql("uuid")
    default = sql("gen_random_uuid()")
  }

  column "created_at" {
    null    = false
    type    = sql("timestamptz")
    default = sql("CURRENT_TIMESTAMP")
  }

  column "updated_at" {
    null    = false
    type    = sql("timestamptz")
    default = sql("CURRENT_TIMESTAMP")
  }

  column "name" {
    null = false
    type = sql("text")
  }

  column "organization_id" {
    null = false
    type = sql("uuid")
  }
  primary_key {
    columns = [column.id]
  }
  foreign_key "label_organization_id_fkey" {
    columns     = [column.organization_id]
    ref_columns = [table.organization.column.id]
    on_update   = "CASCADE"
    on_delete   = "CASCADE"
  }
  index "idx_label_name" {
    columns = [column.name]
  }
  index "idx_label_organization_id" {
    columns = [column.organization_id]
  }
}

table "member" {
  schema = schema.public

  column "id" {
    null    = false
    type    = sql("uuid")
    default = sql("gen_random_uuid()")
  }

  column "created_at" {
    null    = false
    type    = sql("timestamptz")
    default = sql("CURRENT_TIMESTAMP")
  }

  column "updated_at" {
    null    = false
    type    = sql("timestamptz")
    default = sql("CURRENT_TIMESTAMP")
  }

  column "organization_id" {
    null = false
    type = sql("uuid")
  }

  column "role" {
    null    = false
    type    = sql("text")
    default = "basic"
  }

  column "user_id" {
    null = false
    type = sql("uuid")
  }
  primary_key {
    columns = [column.id]
  }
  foreign_key "member_organization_id_fkey" {
    columns     = [column.organization_id]
    ref_columns = [table.organization.column.id]
    on_update   = "NO_ACTION"
    on_delete   = "CASCADE"
  }
  foreign_key "member_user_id_fkey" {
    columns     = [column.user_id]
    ref_columns = [table.user.column.id]
    on_update   = "NO_ACTION"
    on_delete   = "CASCADE"
  }
  index "idx_member_organization_id" {
    columns = [column.organization_id]
  }
  index "idx_member_user_id" {
    columns = [column.user_id]
  }
  check "member_role_check" {
    expr = "(role = ANY (ARRAY['admin'::text, 'owner'::text, 'basic'::text]))"
  }
}

table "organization" {
  schema = schema.public

  column "id" {
    null    = false
    type    = sql("uuid")
    default = sql("gen_random_uuid()")
  }

  column "created_at" {
    null    = false
    type    = sql("timestamptz")
    default = sql("CURRENT_TIMESTAMP")
  }

  column "updated_at" {
    null    = false
    type    = sql("timestamptz")
    default = sql("CURRENT_TIMESTAMP")
  }

  column "billing_email" {
    null = true
    type = sql("text")
  }

  column "credits" {
    null    = false
    type    = sql("integer")
    default = 0
  }

  column "logo" {
    null = true
    type = sql("text")
  }

  column "name" {
    null = false
    type = sql("text")
  }

  column "plan" {
    null    = false
    type    = sql("text")
    default = "FREE"
  }

  column "slug" {
    null = false
    type = sql("text")
  }

  column "stripe_customer_identifier" {
    null = false
    type = sql("text")
  }
  primary_key {
    columns = [column.id]
  }
  check "organization_plan_check" {
    expr = "(plan = ANY (ARRAY['FREE'::text, 'BASIC'::text, 'STANDARD'::text, 'PREMIUM'::text, 'UNLIMITED'::text]))"
  }
}

table "pipeline" {
  schema = schema.public

  column "id" {
    null    = false
    type    = sql("uuid")
    default = sql("gen_random_uuid()")
  }

  column "created_at" {
    null    = false
    type    = sql("timestamptz")
    default = sql("CURRENT_TIMESTAMP")
  }

  column "updated_at" {
    null    = false
    type    = sql("timestamptz")
    default = sql("CURRENT_TIMESTAMP")
  }

  column "description" {
    null = true
    type = sql("text")
  }

  column "name" {
    null = true
    type = sql("text")
  }

  column "organization_id" {
    null = false
    type = sql("uuid")
  }
  primary_key {
    columns = [column.id]
  }
  foreign_key "pipeline_organization_id_fkey" {
    columns     = [column.organization_id]
    ref_columns = [table.organization.column.id]
    on_update   = "CASCADE"
    on_delete   = "CASCADE"
  }
  index "idx_pipeline_organization_id" {
    columns = [column.organization_id]
  }
}

table "pipeline_step" {
  schema = schema.public

  column "id" {
    null    = false
    type    = sql("uuid")
    default = sql("gen_random_uuid()")
  }

  column "created_at" {
    null    = false
    type    = sql("timestamptz")
    default = sql("CURRENT_TIMESTAMP")
  }

  column "updated_at" {
    null    = false
    type    = sql("timestamptz")
    default = sql("CURRENT_TIMESTAMP")
  }

  column "pipeline_id" {
    null = false
    type = sql("uuid")
  }

  column "tool_id" {
    null = false
    type = sql("uuid")
  }
  primary_key {
    columns = [column.id]
  }
}

table "run" {
  schema = schema.public

  column "id" {
    null    = false
    type    = sql("uuid")
    default = sql("gen_random_uuid()")
  }

  column "created_at" {
    null    = false
    type    = sql("timestamptz")
    default = sql("CURRENT_TIMESTAMP")
  }

  column "updated_at" {
    null    = false
    type    = sql("timestamptz")
    default = sql("CURRENT_TIMESTAMP")
  }

  column "completed_at" {
    null = true
    type = sql("timestamptz")
  }

  column "error" {
    null = true
    type = sql("text")
  }

  column "organization_id" {
    null = false
    type = sql("uuid")
  }

  column "pipeline_id" {
    null = false
    type = sql("uuid")
  }

  column "progress" {
    null    = false
    type    = sql("integer")
    default = 0
  }

  column "started_at" {
    null = true
    type = sql("timestamptz")
  }

  column "status" {
    null    = false
    type    = sql("text")
    default = "QUEUED"
  }

  column "tool_id" {
    null = false
    type = sql("uuid")
  }
  primary_key {
    columns = [column.id]
  }
  foreign_key "run_organization_id_fkey" {
    columns     = [column.organization_id]
    ref_columns = [table.organization.column.id]
    on_update   = "CASCADE"
    on_delete   = "CASCADE"
  }
  foreign_key "run_pipeline_id_fkey" {
    columns     = [column.pipeline_id]
    ref_columns = [table.pipeline.column.id]
    on_update   = "CASCADE"
    on_delete   = "SET_NULL"
  }
  foreign_key "run_tool_id_fkey" {
    columns     = [column.tool_id]
    ref_columns = [table.tool.column.id]
    on_update   = "CASCADE"
    on_delete   = "SET_NULL"
  }
  index "idx_run_organization_id" {
    columns = [column.organization_id]
  }
  index "idx_run_pipeline_id" {
    columns = [column.pipeline_id]
  }
  index "idx_run_tool_id" {
    columns = [column.tool_id]
  }
  check "run_status_check" {
    expr = "(status = ANY (ARRAY['COMPLETED'::text, 'FAILED'::text, 'PROCESSING'::text, 'QUEUED'::text]))"
  }
}

table "session" {
  schema = schema.public

  column "id" {
    null    = false
    type    = sql("uuid")
    default = sql("gen_random_uuid()")
  }

  column "created_at" {
    null    = false
    type    = sql("timestamptz")
    default = sql("CURRENT_TIMESTAMP")
  }

  column "updated_at" {
    null    = false
    type    = sql("timestamptz")
    default = sql("CURRENT_TIMESTAMP")
  }

  column "auth_method" {
    null = true
    type = sql("text")
  }

  column "auth_provider" {
    null = true
    type = sql("text")
  }

  column "expires_at" {
    null = false
    type = sql("timestamptz")
  }

  column "ip_address" {
    null = true
    type = sql("text")
  }

  column "organization_id" {
    null = true
    type = sql("uuid")
  }

  column "token" {
    null = false
    type = sql("text")
  }

  column "user_agent" {
    null = true
    type = sql("text")
  }

  column "user_id" {
    null = false
    type = sql("uuid")
  }
  primary_key {
    columns = [column.id]
  }
  foreign_key "session_user_id_fkey" {
    columns     = [column.user_id]
    ref_columns = [table.user.column.id]
    on_update   = "NO_ACTION"
    on_delete   = "CASCADE"
  }
  index "idx_session_token" {
    columns = [column.token]
  }
  index "idx_session_user_id" {
    columns = [column.user_id]
  }
  check "session_auth_provider_check" {
    expr = "(auth_provider = ANY (ARRAY['local'::text, 'google'::text, 'github'::text, 'microsoft'::text, 'apple'::text]))"
  }
}

table "tool" {
  schema = schema.public

  column "id" {
    null    = false
    type    = sql("uuid")
    default = sql("gen_random_uuid()")
  }

  column "created_at" {
    null    = false
    type    = sql("timestamptz")
    default = sql("CURRENT_TIMESTAMP")
  }

  column "updated_at" {
    null    = false
    type    = sql("timestamptz")
    default = sql("CURRENT_TIMESTAMP")
  }

  column "description" {
    null = false
    type = sql("text")
  }

  column "input_mime_type" {
    null    = false
    type    = sql("text")
    default = "application/octet-stream"
  }

  column "name" {
    null = false
    type = sql("text")
  }

  column "organization_id" {
    null = false
    type = sql("uuid")
  }

  column "output_mime_type" {
    null    = false
    type    = sql("text")
    default = "application/octet-stream"
  }
  primary_key {
    columns = [column.id]
  }
  foreign_key "tool_organization_id_fkey" {
    columns     = [column.organization_id]
    ref_columns = [table.organization.column.id]
    on_update   = "CASCADE"
    on_delete   = "CASCADE"
  }
  index "idx_tool_organization_id" {
    columns = [column.organization_id]
  }
}

table "user" {
  schema = schema.public

  column "id" {
    null    = false
    type    = sql("uuid")
    default = sql("gen_random_uuid()")
  }

  column "created_at" {
    null    = false
    type    = sql("timestamptz")
    default = sql("CURRENT_TIMESTAMP")
  }

  column "updated_at" {
    null    = false
    type    = sql("timestamptz")
    default = sql("CURRENT_TIMESTAMP")
  }

  column "email" {
    null = false
    type = sql("text")
  }

  column "email_verified" {
    null    = false
    type    = sql("boolean")
    default = false
  }

  column "image" {
    null = true
    type = sql("text")
  }

  column "name" {
    null = false
    type = sql("text")
  }
  primary_key {
    columns = [column.id]
  }
}


schema "public" {
  comment = "standard public schema"
}