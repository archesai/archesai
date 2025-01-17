import { RunStatusEnum } from '@/src/runs/entities/run.entity'

export class HealthDto {
  error?: any
  status: RunStatusEnum
}
