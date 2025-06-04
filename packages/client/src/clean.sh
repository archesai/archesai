find ./src/generated -type f -name "*.ts" -exec sed -i "s|'../orval.schemas'|'../orval.schemas.ts'|g" {} +
