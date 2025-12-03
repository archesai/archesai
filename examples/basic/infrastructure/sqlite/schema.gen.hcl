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


schema "main" {
}