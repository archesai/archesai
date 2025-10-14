table "account" {
  schema = schema.public
  column "id" {
    null    = false
    type    = uuid
    default = sql("gen_random_uuid()")
  }
  column "created_at" {
    null    = false
    type    = timestamptz
    default = sql("CURRENT_TIMESTAMP")
  }
  column "updated_at" {
    null    = false
    type    = timestamptz
    default = sql("CURRENT_TIMESTAMP")
  }
  column "access_token" {
    null = true
    type = text
  }
  column "access_token_expires_at" {
    null = true
    type = timestamptz
  }
  column "account_identifier" {
    null = false
    type = text
  }
  column "id_token" {
    null = true
    type = text
  }
  column "password" {
    null = true
    type = text
  }
  column "provider" {
    null = false
    type = text
  }
  column "refresh_token" {
    null = true
    type = text
  }
  column "refresh_token_expires_at" {
    null = true
    type = timestamptz
  }
  column "scope" {
    null = true
    type = text
  }
  column "user_id" {
    null = false
    type = uuid
  }
  primary_key {
    columns = [column.id]
  }
  foreign_key "account_user_id_fkey" {
    columns     = [column.user_id]
    ref_columns = [table.user.column.id]
    on_update   = NO_ACTION
    on_delete   = CASCADE
  }
  index "idx_account_user_id" {
    columns = [column.user_id]
  }
}
table "api_keys" {
  schema = schema.public
  column "id" {
    null    = false
    type    = uuid
    default = sql("gen_random_uuid()")
  }
  column "created_at" {
    null    = false
    type    = timestamptz
    default = sql("CURRENT_TIMESTAMP")
  }
  column "updated_at" {
    null    = false
    type    = timestamptz
    default = sql("CURRENT_TIMESTAMP")
  }
  column "expires_at" {
    null = true
    type = timestamptz
  }
  column "key_hash" {
    null = false
    type = text
  }
  column "last_used_at" {
    null = true
    type = timestamptz
  }
  column "name" {
    null = true
    type = text
  }
  column "organization_id" {
    null = false
    type = uuid
  }
  column "prefix" {
    null = true
    type = text
  }
  column "rate_limit" {
    null    = false
    type    = integer
    default = 60
  }
  column "scopes" {
    null    = false
    type    = sql("text[]")
    default = "{}"
  }
  column "user_id" {
    null = false
    type = uuid
  }
  primary_key {
    columns = [column.id]
  }
  foreign_key "api_keys_organization_id_fkey" {
    columns     = [column.organization_id]
    ref_columns = [table.organization.column.id]
    on_update   = NO_ACTION
    on_delete   = CASCADE
  }
  foreign_key "api_keys_user_id_fkey" {
    columns     = [column.user_id]
    ref_columns = [table.user.column.id]
    on_update   = NO_ACTION
    on_delete   = CASCADE
  }
  index "idx_api_keys_key_hash" {
    columns = [column.key_hash]
  }
  index "idx_api_keys_last_used_at" {
    columns = [column.last_used_at]
  }
  index "idx_api_keys_organization_id" {
    columns = [column.organization_id]
  }
  index "idx_api_keys_user_id" {
    columns = [column.user_id]
  }
}
table "artifact" {
  schema = schema.public
  column "id" {
    null    = false
    type    = uuid
    default = sql("gen_random_uuid()")
  }
  column "created_at" {
    null    = false
    type    = timestamptz
    default = sql("CURRENT_TIMESTAMP")
  }
  column "updated_at" {
    null    = false
    type    = timestamptz
    default = sql("CURRENT_TIMESTAMP")
  }
  column "credits" {
    null    = false
    type    = integer
    default = 0
  }
  column "description" {
    null = true
    type = text
  }
  column "embedding" {
    null = true
    type = sql("vector(1536)")
  }
  column "mime_type" {
    null    = false
    type    = text
    default = "application/octet-stream"
  }
  column "name" {
    null = true
    type = text
  }
  column "organization_id" {
    null = false
    type = uuid
  }
  column "preview_image" {
    null = true
    type = text
  }
  column "producer_id" {
    null = true
    type = uuid
  }
  column "text" {
    null = true
    type = text
  }
  column "url" {
    null = true
    type = text
  }
  primary_key {
    columns = [column.id]
  }
  foreign_key "artifact_organization_id_fkey" {
    columns     = [column.organization_id]
    ref_columns = [table.organization.column.id]
    on_update   = CASCADE
    on_delete   = CASCADE
  }
  foreign_key "artifact_producer_id_fkey" {
    columns     = [column.producer_id]
    ref_columns = [table.run.column.id]
    on_update   = CASCADE
    on_delete   = SET_NULL
  }
  index "artifact_embedding_idx" {
    type = "ivfflat"
    on {
      column = column.embedding
      ops    = "vector_cosine_ops"
    }
  }
  index "idx_artifact_organization_id" {
    columns = [column.organization_id]
  }
  index "idx_artifact_producer_id" {
    columns = [column.producer_id]
  }
}
table "goose_db_version" {
  schema = schema.public
  column "id" {
    null = false
    type = integer
    identity {
      generated = BY_DEFAULT
    }
  }
  column "version_id" {
    null = false
    type = bigint
  }
  column "is_applied" {
    null = false
    type = boolean
  }
  column "tstamp" {
    null    = false
    type    = timestamp
    default = sql("now()")
  }
  primary_key {
    columns = [column.id]
  }
}
table "invitation" {
  schema = schema.public
  column "id" {
    null    = false
    type    = uuid
    default = sql("gen_random_uuid()")
  }
  column "created_at" {
    null    = false
    type    = timestamptz
    default = sql("CURRENT_TIMESTAMP")
  }
  column "updated_at" {
    null    = false
    type    = timestamptz
    default = sql("CURRENT_TIMESTAMP")
  }
  column "email" {
    null = false
    type = text
  }
  column "expires_at" {
    null = false
    type = timestamptz
  }
  column "inviter_id" {
    null = false
    type = uuid
  }
  column "organization_id" {
    null = false
    type = uuid
  }
  column "role" {
    null    = false
    type    = text
    default = "member"
  }
  column "status" {
    null    = false
    type    = text
    default = "pending"
  }
  primary_key {
    columns = [column.id]
  }
  foreign_key "invitation_inviter_id_fkey" {
    columns     = [column.inviter_id]
    ref_columns = [table.user.column.id]
    on_update   = NO_ACTION
    on_delete   = CASCADE
  }
  foreign_key "invitation_organization_id_fkey" {
    columns     = [column.organization_id]
    ref_columns = [table.organization.column.id]
    on_update   = NO_ACTION
    on_delete   = CASCADE
  }
  index "idx_invitation_inviter_id" {
    columns = [column.inviter_id]
  }
  index "idx_invitation_organization_id" {
    columns = [column.organization_id]
  }
  check "invitation_role_check" {
    expr = "(role = ANY (ARRAY['admin'::text, 'owner'::text, 'member'::text]))"
  }
}
table "label" {
  schema = schema.public
  column "id" {
    null    = false
    type    = uuid
    default = sql("gen_random_uuid()")
  }
  column "created_at" {
    null    = false
    type    = timestamptz
    default = sql("CURRENT_TIMESTAMP")
  }
  column "updated_at" {
    null    = false
    type    = timestamptz
    default = sql("CURRENT_TIMESTAMP")
  }
  column "name" {
    null = false
    type = text
  }
  column "organization_id" {
    null = false
    type = uuid
  }
  primary_key {
    columns = [column.id]
  }
  foreign_key "label_organization_id_fkey" {
    columns     = [column.organization_id]
    ref_columns = [table.organization.column.id]
    on_update   = CASCADE
    on_delete   = CASCADE
  }
  index "idx_label_name_organization" {
    unique  = true
    columns = [column.name, column.organization_id]
  }
}
table "label_to_artifact" {
  schema = schema.public
  column "label_id" {
    null = false
    type = uuid
  }
  column "artifact_id" {
    null = false
    type = uuid
  }
  primary_key {
    columns = [column.label_id, column.artifact_id]
  }
  foreign_key "label_to_artifact_artifact_id_fkey" {
    columns     = [column.artifact_id]
    ref_columns = [table.artifact.column.id]
    on_update   = NO_ACTION
    on_delete   = CASCADE
  }
  foreign_key "label_to_artifact_label_id_fkey" {
    columns     = [column.label_id]
    ref_columns = [table.label.column.id]
    on_update   = NO_ACTION
    on_delete   = CASCADE
  }
}
table "magic_link_tokens" {
  schema = schema.public
  column "id" {
    null    = false
    type    = uuid
    default = sql("gen_random_uuid()")
  }
  column "code" {
    null = true
    type = character_varying(6)
  }
  column "created_at" {
    null    = true
    type    = timestamp
    default = sql("now()")
  }
  column "delivery_method" {
    null = true
    type = character_varying(50)
  }
  column "expires_at" {
    null = false
    type = timestamp
  }
  column "identifier" {
    null = false
    type = character_varying(255)
  }
  column "ip_address" {
    null = true
    type = character_varying(45)
  }
  column "token_hash" {
    null = false
    type = character_varying(255)
  }
  column "used_at" {
    null = true
    type = timestamp
  }
  column "user_agent" {
    null = true
    type = text
  }
  column "user_id" {
    null = true
    type = uuid
  }
  primary_key {
    columns = [column.id]
  }
  foreign_key "magic_link_tokens_user_id_fkey" {
    columns     = [column.user_id]
    ref_columns = [table.user.column.id]
    on_update   = NO_ACTION
    on_delete   = CASCADE
  }
  index "idx_magic_link_tokens_code" {
    columns = [column.code]
    where   = "(code IS NOT NULL)"
  }
  index "idx_magic_link_tokens_expires_at" {
    columns = [column.expires_at]
    where   = "(used_at IS NULL)"
  }
  index "idx_magic_link_tokens_identifier" {
    columns = [column.identifier]
  }
  index "idx_magic_link_tokens_token_hash" {
    columns = [column.token_hash]
  }
  unique "magic_link_tokens_token_hash_key" {
    columns = [column.token_hash]
  }
}
table "member" {
  schema = schema.public
  column "id" {
    null    = false
    type    = uuid
    default = sql("gen_random_uuid()")
  }
  column "created_at" {
    null    = false
    type    = timestamptz
    default = sql("CURRENT_TIMESTAMP")
  }
  column "updated_at" {
    null    = false
    type    = timestamptz
    default = sql("CURRENT_TIMESTAMP")
  }
  column "organization_id" {
    null = false
    type = uuid
  }
  column "role" {
    null    = false
    type    = text
    default = "member"
  }
  column "user_id" {
    null = false
    type = uuid
  }
  primary_key {
    columns = [column.id]
  }
  foreign_key "member_organization_id_fkey" {
    columns     = [column.organization_id]
    ref_columns = [table.organization.column.id]
    on_update   = NO_ACTION
    on_delete   = CASCADE
  }
  foreign_key "member_user_id_fkey" {
    columns     = [column.user_id]
    ref_columns = [table.user.column.id]
    on_update   = NO_ACTION
    on_delete   = CASCADE
  }
  index "idx_member_organization_id" {
    columns = [column.organization_id]
  }
  index "idx_member_user_id" {
    columns = [column.user_id]
  }
  check "member_role_check" {
    expr = "(role = ANY (ARRAY['admin'::text, 'owner'::text, 'member'::text]))"
  }
}
table "organization" {
  schema = schema.public
  column "id" {
    null    = false
    type    = uuid
    default = sql("gen_random_uuid()")
  }
  column "created_at" {
    null    = false
    type    = timestamptz
    default = sql("CURRENT_TIMESTAMP")
  }
  column "updated_at" {
    null    = false
    type    = timestamptz
    default = sql("CURRENT_TIMESTAMP")
  }
  column "billing_email" {
    null = true
    type = text
  }
  column "credits" {
    null    = false
    type    = integer
    default = 0
  }
  column "logo" {
    null = true
    type = text
  }
  column "name" {
    null = false
    type = text
  }
  column "plan" {
    null    = false
    type    = text
    default = "FREE"
  }
  column "slug" {
    null = false
    type = character_varying(50)
  }
  column "stripe_customer_identifier" {
    null = false
    type = text
  }
  primary_key {
    columns = [column.id]
  }
  index "idx_organization_slug" {
    unique  = true
    columns = [column.slug]
  }
  check "organization_plan_check" {
    expr = "(plan = ANY (ARRAY['BASIC'::text, 'FREE'::text, 'PREMIUM'::text, 'STANDARD'::text, 'UNLIMITED'::text]))"
  }
  check "organization_slug_format" {
    expr = "((slug)::text ~ '^[a-z0-9]+(?:-[a-z0-9]+)*$'::text)"
  }
  unique "organization_stripe_customer_identifier_key" {
    columns = [column.stripe_customer_identifier]
  }
}
table "pipeline" {
  schema = schema.public
  column "id" {
    null    = false
    type    = uuid
    default = sql("gen_random_uuid()")
  }
  column "created_at" {
    null    = false
    type    = timestamptz
    default = sql("CURRENT_TIMESTAMP")
  }
  column "updated_at" {
    null    = false
    type    = timestamptz
    default = sql("CURRENT_TIMESTAMP")
  }
  column "description" {
    null = true
    type = text
  }
  column "name" {
    null = true
    type = text
  }
  column "organization_id" {
    null = false
    type = uuid
  }
  primary_key {
    columns = [column.id]
  }
  foreign_key "pipeline_organization_id_fkey" {
    columns     = [column.organization_id]
    ref_columns = [table.organization.column.id]
    on_update   = CASCADE
    on_delete   = CASCADE
  }
  index "idx_pipeline_organization_id" {
    columns = [column.organization_id]
  }
}
table "pipeline_step" {
  schema = schema.public
  column "id" {
    null    = false
    type    = uuid
    default = sql("gen_random_uuid()")
  }
  column "created_at" {
    null    = false
    type    = timestamptz
    default = sql("CURRENT_TIMESTAMP")
  }
  column "updated_at" {
    null    = false
    type    = timestamptz
    default = sql("CURRENT_TIMESTAMP")
  }
  column "pipeline_id" {
    null = false
    type = uuid
  }
  column "tool_id" {
    null = false
    type = uuid
  }
  primary_key {
    columns = [column.id]
  }
  foreign_key "pipeline_step_pipeline_id_fkey" {
    columns     = [column.pipeline_id]
    ref_columns = [table.pipeline.column.id]
    on_update   = CASCADE
    on_delete   = CASCADE
  }
  foreign_key "pipeline_step_tool_id_fkey" {
    columns     = [column.tool_id]
    ref_columns = [table.tool.column.id]
    on_update   = CASCADE
    on_delete   = CASCADE
  }
  index "idx_pipeline_step_pipeline_id" {
    columns = [column.pipeline_id]
  }
  index "idx_pipeline_step_tool_id" {
    columns = [column.tool_id]
  }
}
table "pipeline_step_to_dependency" {
  schema = schema.public
  column "pipeline_step_id" {
    null = false
    type = uuid
  }
  column "prerequisite_id" {
    null = false
    type = uuid
  }
  primary_key {
    columns = [column.pipeline_step_id, column.prerequisite_id]
  }
  foreign_key "pipeline_step_to_dependency_pipeline_step_id_fkey" {
    columns     = [column.pipeline_step_id]
    ref_columns = [table.pipeline_step.column.id]
    on_update   = NO_ACTION
    on_delete   = CASCADE
  }
  foreign_key "pipeline_step_to_dependency_prerequisite_id_fkey" {
    columns     = [column.prerequisite_id]
    ref_columns = [table.pipeline_step.column.id]
    on_update   = NO_ACTION
    on_delete   = CASCADE
  }
}
table "run" {
  schema = schema.public
  column "id" {
    null    = false
    type    = uuid
    default = sql("gen_random_uuid()")
  }
  column "created_at" {
    null    = false
    type    = timestamptz
    default = sql("CURRENT_TIMESTAMP")
  }
  column "updated_at" {
    null    = false
    type    = timestamptz
    default = sql("CURRENT_TIMESTAMP")
  }
  column "completed_at" {
    null = true
    type = timestamptz
  }
  column "error" {
    null = true
    type = text
  }
  column "organization_id" {
    null = false
    type = uuid
  }
  column "pipeline_id" {
    null = false
    type = uuid
  }
  column "progress" {
    null    = false
    type    = double_precision
    default = 0
  }
  column "started_at" {
    null = true
    type = timestamptz
  }
  column "status" {
    null    = false
    type    = text
    default = "QUEUED"
  }
  column "tool_id" {
    null = false
    type = uuid
  }
  primary_key {
    columns = [column.id]
  }
  foreign_key "run_organization_id_fkey" {
    columns     = [column.organization_id]
    ref_columns = [table.organization.column.id]
    on_update   = CASCADE
    on_delete   = CASCADE
  }
  foreign_key "run_pipeline_id_fkey" {
    columns     = [column.pipeline_id]
    ref_columns = [table.pipeline.column.id]
    on_update   = CASCADE
    on_delete   = SET_NULL
  }
  foreign_key "run_tool_id_fkey" {
    columns     = [column.tool_id]
    ref_columns = [table.tool.column.id]
    on_update   = CASCADE
    on_delete   = SET_NULL
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
table "run_to_artifact" {
  schema = schema.public
  column "run_id" {
    null = false
    type = uuid
  }
  column "artifact_id" {
    null = false
    type = uuid
  }
  primary_key {
    columns = [column.run_id, column.artifact_id]
  }
  foreign_key "run_to_artifact_artifact_id_fkey" {
    columns     = [column.artifact_id]
    ref_columns = [table.artifact.column.id]
    on_update   = NO_ACTION
    on_delete   = CASCADE
  }
  foreign_key "run_to_artifact_run_id_fkey" {
    columns     = [column.run_id]
    ref_columns = [table.run.column.id]
    on_update   = NO_ACTION
    on_delete   = CASCADE
  }
}
table "session" {
  schema = schema.public
  column "id" {
    null    = false
    type    = uuid
    default = sql("gen_random_uuid()")
  }
  column "created_at" {
    null    = false
    type    = timestamptz
    default = sql("CURRENT_TIMESTAMP")
  }
  column "updated_at" {
    null    = false
    type    = timestamptz
    default = sql("CURRENT_TIMESTAMP")
  }
  column "auth_method" {
    null = true
    type = character_varying(50)
  }
  column "auth_provider" {
    null = true
    type = character_varying(50)
  }
  column "expires_at" {
    null = false
    type = timestamptz
  }
  column "ip_address" {
    null = true
    type = text
  }
  column "metadata" {
    null = true
    type = jsonb
  }
  column "organization_id" {
    null = true
    type = uuid
  }
  column "token" {
    null = false
    type = text
  }
  column "user_agent" {
    null = true
    type = text
  }
  column "user_id" {
    null = false
    type = uuid
  }
  primary_key {
    columns = [column.id]
  }
  foreign_key "session_user_id_fkey" {
    columns     = [column.user_id]
    ref_columns = [table.user.column.id]
    on_update   = NO_ACTION
    on_delete   = CASCADE
  }
  index "idx_session_user_id" {
    columns = [column.user_id]
  }
  unique "session_token_key" {
    columns = [column.token]
  }
}
table "tool" {
  schema = schema.public
  column "id" {
    null    = false
    type    = uuid
    default = sql("gen_random_uuid()")
  }
  column "created_at" {
    null    = false
    type    = timestamptz
    default = sql("CURRENT_TIMESTAMP")
  }
  column "updated_at" {
    null    = false
    type    = timestamptz
    default = sql("CURRENT_TIMESTAMP")
  }
  column "description" {
    null = false
    type = text
  }
  column "input_mime_type" {
    null    = false
    type    = text
    default = "application/octet-stream"
  }
  column "name" {
    null = false
    type = text
  }
  column "organization_id" {
    null = false
    type = uuid
  }
  column "output_mime_type" {
    null    = false
    type    = text
    default = "application/octet-stream"
  }
  primary_key {
    columns = [column.id]
  }
  foreign_key "tool_organization_id_fkey" {
    columns     = [column.organization_id]
    ref_columns = [table.organization.column.id]
    on_update   = CASCADE
    on_delete   = CASCADE
  }
  index "idx_tool_organization_id" {
    columns = [column.organization_id]
  }
}
table "user" {
  schema = schema.public
  column "id" {
    null    = false
    type    = uuid
    default = sql("gen_random_uuid()")
  }
  column "created_at" {
    null    = false
    type    = timestamptz
    default = sql("CURRENT_TIMESTAMP")
  }
  column "updated_at" {
    null    = false
    type    = timestamptz
    default = sql("CURRENT_TIMESTAMP")
  }
  column "email" {
    null = false
    type = text
  }
  column "email_verified" {
    null    = false
    type    = boolean
    default = false
  }
  column "image" {
    null = true
    type = text
  }
  column "name" {
    null = false
    type = text
  }
  primary_key {
    columns = [column.id]
  }
  unique "user_email_key" {
    columns = [column.email]
  }
}
table "verification_token" {
  schema = schema.public
  column "id" {
    null    = false
    type    = uuid
    default = sql("gen_random_uuid()")
  }
  column "created_at" {
    null    = false
    type    = timestamptz
    default = sql("CURRENT_TIMESTAMP")
  }
  column "updated_at" {
    null    = false
    type    = timestamptz
    default = sql("CURRENT_TIMESTAMP")
  }
  column "expires_at" {
    null = false
    type = timestamptz
  }
  column "identifier" {
    null = false
    type = text
  }
  column "value" {
    null = false
    type = text
  }
  primary_key {
    columns = [column.id]
  }
}
schema "public" {
  comment = "standard public schema"
}
