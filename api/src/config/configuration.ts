import { readFileSync, existsSync, statSync, readdirSync } from 'fs'
import yaml from 'js-yaml'
import { join } from 'path'
import { z, ZodSchema } from 'zod'
import merge from 'deepmerge'

export const loadConfiguration = <T extends ZodSchema>(schema: T) => {
  const yamlConfig = loadYamlConfig()
  const envConfig = loadEnvConfig()
  const mergedConfig = merge(yamlConfig, envConfig)
  if (process.env.ARCHES_NOVALIDATE === 'false') {
    return mergedConfig as z.infer<T>
  }
  return schema.parse(mergedConfig) as z.infer<T>
}

export function loadYamlConfig(): Record<string, any> {
  const configDir = process.env.ARCHES_CONF_DIR || '/etc/archesai'
  const baseConfigPath = join(configDir, 'config.yaml')

  let baseConfig: Record<string, any> = {}
  if (existsSync(baseConfigPath) && statSync(baseConfigPath).isFile()) {
    baseConfig = yaml.load(readFileSync(baseConfigPath, 'utf8')) as Record<
      string,
      any
    >
  }

  const configDropInDir = join(baseConfigPath + '.d')
  if (existsSync(configDropInDir) && statSync(configDropInDir).isDirectory()) {
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

export function loadEnvConfig(): Record<string, any> {
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
