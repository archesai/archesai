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