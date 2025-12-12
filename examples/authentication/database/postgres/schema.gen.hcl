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
  index "idx_api_key_organization_id" {
    columns = [column.organization_id]
  }
  index "idx_api_key_user_id" {
    columns = [column.user_id]
  }
  index "idx_api_key_key_hash" {
    columns = [column.key_hash]
  }
  index "idx_api_key_last_used_at" {
    columns = [column.last_used_at]
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
  index "idx_invitation_organization_id" {
    columns = [column.organization_id]
  }
  index "idx_invitation_inviter_id" {
    columns = [column.inviter_id]
  }
  check "invitation_role_check" {
    expr = "(role = ANY (ARRAY['admin'::text, 'owner'::text, 'basic'::text]))"
  }
  check "invitation_status_check" {
    expr = "(status = ANY (ARRAY['pending'::text, 'accepted'::text, 'declined'::text, 'expired'::text]))"
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
  index "idx_session_user_id" {
    columns = [column.user_id]
  }
  index "idx_session_token" {
    columns = [column.token]
  }
  check "session_auth_provider_check" {
    expr = "(auth_provider = ANY (ARRAY['local'::text, 'google'::text, 'github'::text, 'microsoft'::text, 'apple'::text]))"
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