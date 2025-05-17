import fs from 'node:fs'
import path from 'node:path'
import type { StaticDecode, TObject } from '@sinclair/typebox'

import { deepmerge } from '@fastify/deepmerge'
import { Type } from '@sinclair/typebox'
import { Value } from '@sinclair/typebox/value'
import { parse } from 'yaml'

/**
 * ConfigLoader is a utility class that loads configuration from the following sources:
 * 1. Environment variables
 * 2. YAML files
 * 3. YAML files in a drop-in directory
 */
export class ConfigLoader {
  /**
   * Loads and merges configuration data from YAML and environment variables,
   * then validates the merged configuration against the provided schema.
   * @template T - The type of the schema object.
   * @param schema - The schema object used to validate the configuration.
   * @returns The validated and parsed configuration object.
   */
  public static load<T extends TObject>(schema: T): StaticDecode<T> {
    const yamlConfig = this.loadYamlConfig()
    const envConfig = this.loadEnvConfig()
    const mergedConfig = deepmerge({ all: true })({
      ...envConfig,
      ...yamlConfig
    })
    return Value.Parse(schema, mergedConfig)
  }

  /**
   * Loads environment variables prefixed with "ARCHES_" and transforms them into a nested configuration object.
   *
   * The method iterates over all environment variables, filtering those that start with "ARCHES_".
   * It removes the prefix, converts the remaining key to lowercase, and splits it by underscores ("_").
   * The resulting keys are used to construct a nested object structure, where each segment of the key
   * represents a level in the hierarchy.
   * @returns A nested configuration object derived from the environment variables.
   */
  private static loadEnvConfig(): Record<string, unknown> {
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

  /**
   * Loads and parses YAML configuration files from a specified directory.
   *
   * This method first attempts to load a base configuration file named `config.yaml`
   * from the directory specified by the `ARCHES_CONFIG_ROOT` environment variable.
   * If the environment variable is not set, it defaults to the current working directory.
   *
   * After loading the base configuration, it checks for an optional drop-in directory
   * named `config.yaml.d` located in the same directory as the base configuration file.
   * If the drop-in directory exists, it loads and merges all `.yaml` files within it
   * in lexicographical order. The merging process ensures that configurations from
   * later files override those from earlier ones.
   *
   * Validation is performed on each YAML file to ensure it conforms to a structure
   * where keys are strings and values are unknown types.
   * @returns The merged configuration object.
   * @throws {Error} If any YAML file is invalid or cannot be parsed.
   */
  private static loadYamlConfig(): Record<string, unknown> {
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
}
