import { Injectable, Logger } from '@nestjs/common'
import { readFileSync, existsSync, statSync, readdirSync } from 'fs'
import yaml from 'js-yaml'
import { join } from 'path'
import merge from 'deepmerge'
import { archesConfigSchema } from './schemas/config.schema'

import { ArchesConfig } from './schemas/config.schema'
import { z, ZodSchema } from 'zod'
import { RunStatusEnum } from '@/src/runs/entities/run.entity'
import { HealthDto } from '@/src/health/dto/health.dto'

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
export class ConfigService {
  private readonly logger = new Logger(ConfigService.name)
  private config: ArchesConfig
  private health: HealthDto = {
    status: RunStatusEnum.QUEUED
  }

  constructor() {
    try {
      this.logger.log('initializing configuration')
      this.health.status = RunStatusEnum.PROCESSING
      this.config = this.loadConfiguration(archesConfigSchema)
      this.logger.log('loaded configuration')
      this.health.status = RunStatusEnum.COMPLETE
    } catch (err) {
      this.logger.error(err, `error loading config`)
      this.health.status = RunStatusEnum.ERROR
      this.health.error = err
    }
  }

  public get<T extends Leaves<ArchesConfig>>(
    propertyPath: T
  ): LeafTypes<ArchesConfig, T> {
    return propertyPath
      .split('.') // split on dot
      .reduce((acc, key) => acc?.[key], this.config as any) as LeafTypes<
      ArchesConfig,
      T
    >
  }

  public getHealth(): { status: RunStatusEnum; error?: any } {
    return this.health
  }

  public getConfig(): ArchesConfig {
    return this.config
  }

  private loadYamlConfig(): Record<string, any> {
    const configDir = process.env['ARCHES.CONFDIR'] || '/etc/archesai'
    const baseConfigPath = join(configDir, 'config.yaml')

    let baseConfig: Record<string, any> = {}
    if (existsSync(baseConfigPath) && statSync(baseConfigPath).isFile()) {
      baseConfig = yaml.load(readFileSync(baseConfigPath, 'utf8')) as Record<
        string,
        any
      >
    }

    const configDropInDir = join(baseConfigPath + '.d')
    if (
      existsSync(configDropInDir) &&
      statSync(configDropInDir).isDirectory()
    ) {
      const files = readdirSync(configDropInDir)
        .filter((file) => file.endsWith('.yaml'))
        .sort() // Alphabetical order
        .map((file) => join(configDropInDir, file))

      for (const file of files) {
        const fileConfig = yaml.load(readFileSync(file, 'utf8')) as Record<
          string,
          any
        >
        baseConfig = merge(baseConfig, fileConfig)
      }
    }

    return baseConfig
  }

  private loadEnvConfig(): Record<string, any> {
    const envConfig: Record<string, any> = {}
    for (const [key, value] of Object.entries(process.env)) {
      if (key.startsWith('ARCHES.')) {
        const keys = key.replace('ARCHES.', '').toLowerCase().split('.')
        let current = envConfig
        keys.forEach((k, index) => {
          if (index === keys.length - 1) {
            current[k] = value // Set the final key
          } else {
            current[k] = current[k] || {}
            current = current[k]
          }
        })
      }
    }
    return envConfig
  }

  public loadConfiguration = <T extends ZodSchema>(schema: T) => {
    const yamlConfig = this.loadYamlConfig()
    this.logger.error(yamlConfig, 'loaded yaml configuration')

    const envConfig = this.loadEnvConfig()
    this.logger.error(envConfig, 'loaded env configuration')

    const mergedConfig = merge(yamlConfig, envConfig)
    this.logger.error(mergedConfig, 'merged configuration')

    try {
      const validationResult = schema.parse(mergedConfig)
      return validationResult.data as z.infer<T>
    } catch (err: any) {
      if (mergedConfig.config.validate === 'false') {
        this.logger.error(err, 'configuration failed - ignoring')
        return mergedConfig as z.infer<T>
      }
      throw err
    }
  }
}
