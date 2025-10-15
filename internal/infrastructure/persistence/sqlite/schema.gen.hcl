table "account" {
  schema = schema.main
  column "id" {
    null    = false
    type    = sql("TEXT")
    default = sql("lower(hex(randomblob(16)))")
  }
  column "created_at" {
    null    = false
    type    = sql("TEXT")
    default = sql("CURRENT_TIMESTAMP")
  }
  column "updated_at" {
    null    = false
    type    = sql("TEXT")
    default = sql("CURRENT_TIMESTAMP")
  }
  column "access_token" {
    null = true
    type = sql("TEXT")
  }
  column "access_token_expires_at" {
    null = true
    type = sql("TEXT")
  }
  column "account_identifier" {
    null = false
    type = sql("TEXT")
  }
  column "id_token" {
    null = true
    type = sql("TEXT")
  }
  column "provider" {
    null = false
    type = sql("TEXT")
  }
  column "refresh_token" {
    null = true
    type = sql("TEXT")
  }
  column "refresh_token_expires_at" {
    null = true
    type = sql("TEXT")
  }
  column "scope" {
    null = true
    type = sql("TEXT")
  }
  column "user_id" {
    null = false
    type = sql("TEXT")
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
}

table "api_key" {
  schema = schema.main
  column "id" {
    null    = false
    type    = sql("TEXT")
    default = sql("lower(hex(randomblob(16)))")
  }
  column "created_at" {
    null    = false
    type    = sql("TEXT")
    default = sql("CURRENT_TIMESTAMP")
  }
  column "updated_at" {
    null    = false
    type    = sql("TEXT")
    default = sql("CURRENT_TIMESTAMP")
  }
  column "expires_at" {
    null = true
    type = sql("TEXT")
  }
  column "key_hash" {
    null = false
    type = sql("TEXT")
  }
  column "last_used_at" {
    null = true
    type = sql("TEXT")
  }
  column "name" {
    null = true
    type = sql("TEXT")
  }
  column "organization_id" {
    null = false
    type = sql("TEXT")
  }
  column "prefix" {
    null = true
    type = sql("TEXT")
  }
  column "rate_limit" {
    null    = false
    type    = sql("INTEGER")
    default = 60
  }
  column "scopes" {
    null    = false
    type    = sql("TEXT")
    default = "[]"
  }
  column "user_id" {
    null = false
    type = sql("TEXT")
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
  schema = schema.main
  column "id" {
    null    = false
    type    = sql("TEXT")
    default = sql("lower(hex(randomblob(16)))")
  }
  column "created_at" {
    null    = false
    type    = sql("TEXT")
    default = sql("CURRENT_TIMESTAMP")
  }
  column "updated_at" {
    null    = false
    type    = sql("TEXT")
    default = sql("CURRENT_TIMESTAMP")
  }
  column "credits" {
    null    = false
    type    = sql("INTEGER")
    default = 0
  }
  column "description" {
    null = true
    type = sql("TEXT")
  }
  column "mime_type" {
    null    = false
    type    = sql("TEXT")
    default = "application/octet-stream"
  }
  column "name" {
    null = true
    type = sql("TEXT")
  }
  column "organization_id" {
    null = false
    type = sql("TEXT")
  }
  column "preview_image" {
    null = true
    type = sql("TEXT")
  }
  column "producer_id" {
    null = true
    type = sql("TEXT")
  }
  column "text" {
    null = true
    type = sql("TEXT")
  }
  column "url" {
    null = true
    type = sql("TEXT")
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

table "invitation" {
  schema = schema.main
  column "id" {
    null    = false
    type    = sql("TEXT")
    default = sql("lower(hex(randomblob(16)))")
  }
  column "created_at" {
    null    = false
    type    = sql("TEXT")
    default = sql("CURRENT_TIMESTAMP")
  }
  column "updated_at" {
    null    = false
    type    = sql("TEXT")
    default = sql("CURRENT_TIMESTAMP")
  }
  column "email" {
    null = false
    type = sql("TEXT")
  }
  column "expires_at" {
    null = false
    type = sql("TEXT")
  }
  column "inviter_id" {
    null = false
    type = sql("TEXT")
  }
  column "organization_id" {
    null = false
    type = sql("TEXT")
  }
  column "role" {
    null    = false
    type    = sql("TEXT")
    default = "basic"
  }
  column "status" {
    null    = false
    type    = sql("TEXT")
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
}

table "label" {
  schema = schema.main
  column "id" {
    null    = false
    type    = sql("TEXT")
    default = sql("lower(hex(randomblob(16)))")
  }
  column "created_at" {
    null    = false
    type    = sql("TEXT")
    default = sql("CURRENT_TIMESTAMP")
  }
  column "updated_at" {
    null    = false
    type    = sql("TEXT")
    default = sql("CURRENT_TIMESTAMP")
  }
  column "name" {
    null = false
    type = sql("TEXT")
  }
  column "organization_id" {
    null = false
    type = sql("TEXT")
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
  schema = schema.main
  column "id" {
    null    = false
    type    = sql("TEXT")
    default = sql("lower(hex(randomblob(16)))")
  }
  column "created_at" {
    null    = false
    type    = sql("TEXT")
    default = sql("CURRENT_TIMESTAMP")
  }
  column "updated_at" {
    null    = false
    type    = sql("TEXT")
    default = sql("CURRENT_TIMESTAMP")
  }
  column "organization_id" {
    null = false
    type = sql("TEXT")
  }
  column "role" {
    null    = false
    type    = sql("TEXT")
    default = "basic"
  }
  column "user_id" {
    null = false
    type = sql("TEXT")
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
}

table "organization" {
  schema = schema.main
  column "id" {
    null    = false
    type    = sql("TEXT")
    default = sql("lower(hex(randomblob(16)))")
  }
  column "created_at" {
    null    = false
    type    = sql("TEXT")
    default = sql("CURRENT_TIMESTAMP")
  }
  column "updated_at" {
    null    = false
    type    = sql("TEXT")
    default = sql("CURRENT_TIMESTAMP")
  }
  column "billing_email" {
    null = true
    type = sql("TEXT")
  }
  column "credits" {
    null    = false
    type    = sql("INTEGER")
    default = 0
  }
  column "logo" {
    null = true
    type = sql("TEXT")
  }
  column "name" {
    null = false
    type = sql("TEXT")
  }
  column "plan" {
    null    = false
    type    = sql("TEXT")
    default = "FREE"
  }
  column "slug" {
    null = false
    type = sql("TEXT")
  }
  column "stripe_customer_identifier" {
    null = false
    type = sql("TEXT")
  }
  primary_key {
    columns = [column.id]
  }
  index "idx_organization_slug" {
    unique  = true
    columns = [column.slug]
  }
}

table "pipeline" {
  schema = schema.main
  column "id" {
    null    = false
    type    = sql("TEXT")
    default = sql("lower(hex(randomblob(16)))")
  }
  column "created_at" {
    null    = false
    type    = sql("TEXT")
    default = sql("CURRENT_TIMESTAMP")
  }
  column "updated_at" {
    null    = false
    type    = sql("TEXT")
    default = sql("CURRENT_TIMESTAMP")
  }
  column "description" {
    null = true
    type = sql("TEXT")
  }
  column "name" {
    null = true
    type = sql("TEXT")
  }
  column "organization_id" {
    null = false
    type = sql("TEXT")
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
  schema = schema.main
  column "id" {
    null    = false
    type    = sql("TEXT")
    default = sql("lower(hex(randomblob(16)))")
  }
  column "created_at" {
    null    = false
    type    = sql("TEXT")
    default = sql("CURRENT_TIMESTAMP")
  }
  column "updated_at" {
    null    = false
    type    = sql("TEXT")
    default = sql("CURRENT_TIMESTAMP")
  }
  column "pipeline_id" {
    null = false
    type = sql("TEXT")
  }
  column "tool_id" {
    null = false
    type = sql("TEXT")
  }
  primary_key {
    columns = [column.id]
  }
}

table "run" {
  schema = schema.main
  column "id" {
    null    = false
    type    = sql("TEXT")
    default = sql("lower(hex(randomblob(16)))")
  }
  column "created_at" {
    null    = false
    type    = sql("TEXT")
    default = sql("CURRENT_TIMESTAMP")
  }
  column "updated_at" {
    null    = false
    type    = sql("TEXT")
    default = sql("CURRENT_TIMESTAMP")
  }
  column "completed_at" {
    null = true
    type = sql("TEXT")
  }
  column "error" {
    null = true
    type = sql("TEXT")
  }
  column "organization_id" {
    null = false
    type = sql("TEXT")
  }
  column "pipeline_id" {
    null = false
    type = sql("TEXT")
  }
  column "progress" {
    null    = false
    type    = sql("INTEGER")
    default = 0
  }
  column "started_at" {
    null = true
    type = sql("TEXT")
  }
  column "status" {
    null    = false
    type    = sql("TEXT")
    default = "QUEUED"
  }
  column "tool_id" {
    null = false
    type = sql("TEXT")
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
}

table "session" {
  schema = schema.main
  column "id" {
    null    = false
    type    = sql("TEXT")
    default = sql("lower(hex(randomblob(16)))")
  }
  column "created_at" {
    null    = false
    type    = sql("TEXT")
    default = sql("CURRENT_TIMESTAMP")
  }
  column "updated_at" {
    null    = false
    type    = sql("TEXT")
    default = sql("CURRENT_TIMESTAMP")
  }
  column "auth_method" {
    null = true
    type = sql("TEXT")
  }
  column "auth_provider" {
    null = true
    type = sql("TEXT")
  }
  column "expires_at" {
    null = false
    type = sql("TEXT")
  }
  column "ip_address" {
    null = true
    type = sql("TEXT")
  }
  column "organization_id" {
    null = true
    type = sql("TEXT")
  }
  column "token" {
    null = false
    type = sql("TEXT")
  }
  column "user_agent" {
    null = true
    type = sql("TEXT")
  }
  column "user_id" {
    null = false
    type = sql("TEXT")
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
  index "idx_session_user_id" {
    columns = [column.user_id]
  }
}

table "tool" {
  schema = schema.main
  column "id" {
    null    = false
    type    = sql("TEXT")
    default = sql("lower(hex(randomblob(16)))")
  }
  column "created_at" {
    null    = false
    type    = sql("TEXT")
    default = sql("CURRENT_TIMESTAMP")
  }
  column "updated_at" {
    null    = false
    type    = sql("TEXT")
    default = sql("CURRENT_TIMESTAMP")
  }
  column "description" {
    null = false
    type = sql("TEXT")
  }
  column "input_mime_type" {
    null    = false
    type    = sql("TEXT")
    default = "application/octet-stream"
  }
  column "name" {
    null = false
    type = sql("TEXT")
  }
  column "organization_id" {
    null = false
    type = sql("TEXT")
  }
  column "output_mime_type" {
    null    = false
    type    = sql("TEXT")
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
  schema = schema.main
  column "id" {
    null    = false
    type    = sql("TEXT")
    default = sql("lower(hex(randomblob(16)))")
  }
  column "created_at" {
    null    = false
    type    = sql("TEXT")
    default = sql("CURRENT_TIMESTAMP")
  }
  column "updated_at" {
    null    = false
    type    = sql("TEXT")
    default = sql("CURRENT_TIMESTAMP")
  }
  column "email" {
    null = false
    type = sql("TEXT")
  }
  column "email_verified" {
    null    = false
    type    = sql("INTEGER")
    default = 0
  }
  column "image" {
    null = true
    type = sql("TEXT")
  }
  column "name" {
    null = false
    type = sql("TEXT")
  }
  primary_key {
    columns = [column.id]
  }
}


schema "main" {
}