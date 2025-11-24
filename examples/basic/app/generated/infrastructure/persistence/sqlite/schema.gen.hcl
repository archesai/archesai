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