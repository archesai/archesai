import fs from 'node:fs'
import path from 'node:path'
import type { StaticDecode } from '@sinclair/typebox'

import { deepmerge } from '@fastify/deepmerge'
import { parse } from 'yaml'

import type { TObject } from '@archesai/schemas'

import { Type, Value } from '@archesai/schemas'

import type { ArchesConfig } from '#config/schemas/config.schema'

import { ArchesConfigSchema } from '#config/schemas/config.schema'

export type LeafTypes<T, S extends string> =
  T extends (
    unknown // distribute across unions explicitly
  ) ?
    S extends `${infer Head}.${infer Tail}` ?
      Head extends keyof T ?
        LeafTypes<NonNullable<T[Head]>, Tail>
      : never
    : S extends keyof T ? T[S]
    : never
  : never

export type Leaves<T> =
  T extends object ?
    {
      [K in keyof T & string]: NonNullable<T[K]> extends object ?
        // If it's an object, include both "K" and deeper paths "K.xxx"
        `${K}.${Leaves<NonNullable<T[K]>>}` | K
      : // Otherwise it's a leaf, just "K"
        K
    }[keyof T & string]
  : never

export const createConfigService = (): {
  get<Path extends Leaves<ArchesConfig>>(
    propertyPath: Path
  ): LeafTypes<ArchesConfig, Path>
  getConfig(): ArchesConfig
  load: <T extends TObject>(schema: T) => StaticDecode<T>
} => {
  const loadEnvConfig = (): Record<string, unknown> => {
    const envConfig: Record<string, unknown> = {}
    for (const [key, value] of Object.entries(process.env)) {
      if (key.startsWith('ARCHES_')) {
        const keys = key.replace('ARCHES_', '').toLowerCase().split('_')
        let current = envConfig
        keys.forEach((k, index) => {
          if (index === keys.length - 1) {
            current[k] = value
          } else {
            current[k] = current[k] ?? {}
            current = current[k] as Record<string, unknown>
          }
        })
      }
    }
    return envConfig
  }

  const loadYamlConfig = (): Record<string, unknown> => {
    const configDir = process.env.ARCHES_CONFIG_ROOT ?? process.cwd()
    const baseConfigPath = path.join(configDir, 'config.yaml')

    let baseConfig: Record<string, unknown> = {}
    if (fs.existsSync(baseConfigPath) && fs.statSync(baseConfigPath).isFile()) {
      const yamlValidation = Type.Record(Type.String(), Type.Unknown())
      baseConfig = Value.Parse(
        yamlValidation,
        parse(fs.readFileSync(baseConfigPath, 'utf8'))
      )
    }

    const configDropInDir = baseConfigPath + '.d'
    if (
      fs.existsSync(configDropInDir) &&
      fs.statSync(configDropInDir).isDirectory()
    ) {
      const files = fs
        .readdirSync(configDropInDir)
        .filter((file) => file.endsWith('.yaml'))
        .sort()
        .map((file) => path.join(configDropInDir, file))

      for (const file of files) {
        const yamlValidation = Type.Record(Type.String(), Type.Unknown())
        const fileConfig = Value.Parse(
          yamlValidation,
          parse(fs.readFileSync(file, 'utf8'))
        )
        baseConfig = deepmerge({ all: true })({
          baseConfig,
          fileConfig
        })
      }
    }
    return baseConfig
  }

  const load = <T extends TObject>(schema: T): StaticDecode<T> => {
    const yamlConfig = loadYamlConfig()
    const envConfig = loadEnvConfig()
    const mergedConfig = deepmerge({ all: true })({
      ...envConfig,
      ...yamlConfig
    })
    return Value.Parse(schema, mergedConfig)
  }

  const config = load(ArchesConfigSchema)

  return {
    get<Path extends Leaves<ArchesConfig>>(
      propertyPath: Path
    ): LeafTypes<ArchesConfig, Path> {
      return propertyPath
        .split('.')
        .reduce<unknown>(
          (acc, key) =>
            acc && typeof acc === 'object' ?
              (acc as Record<string, unknown>)[key]
            : undefined,
          config
        ) as LeafTypes<ArchesConfig, Path>
    },
    getConfig(): ArchesConfig {
      return config
    },
    load
  }
}

export type ConfigService = ReturnType<typeof createConfigService>
