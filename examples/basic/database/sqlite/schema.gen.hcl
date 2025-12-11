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
  index "idx_invitation_organization_id" {
    columns = [column.organization_id]
  }
  index "idx_invitation_inviter_id" {
    columns = [column.inviter_id]
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
  index "idx_session_token" {
    columns = [column.token]
  }
}

table "todo" {
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

  column "completed" {
    null = false
    type = sql("INTEGER")
  }

  column "title" {
    null = false
    type = sql("TEXT")
  }
  primary_key {
    columns = [column.id]
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