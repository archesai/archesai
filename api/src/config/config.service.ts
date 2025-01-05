import { Injectable } from '@nestjs/common'
import { ConfigService } from '@nestjs/config'
import { ArchesConfig } from './schema'

export type Leaves<T> = T extends object
  ? {
      [K in keyof T]: `${Exclude<K, symbol>}${Leaves<T[K]> extends never
        ? ''
        : `.${Leaves<T[K]>}`}`
    }[keyof T]
  : never

export type LeafTypes<T, S extends string> = S extends `${infer T1}.${infer T2}`
  ? T1 extends keyof T
    ? LeafTypes<T[T1], T2>
    : never
  : S extends keyof T
    ? T[S]
    : never

@Injectable()
export class ArchesConfigService {
  constructor(private configService: ConfigService<ArchesConfig, true>) {}

  get<T extends Leaves<ArchesConfig>>(
    propertyPath: T
  ): LeafTypes<ArchesConfig, T> {
    return this.configService.get(propertyPath as any)
  }
}
