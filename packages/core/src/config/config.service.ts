import type { LeafTypes, Leaves } from '#config/interfaces/leaves.interface'
import type { ArchesConfig } from '#config/schemas/config.schema'
import type { HealthCheck } from '#health/interfaces/health-check.interface'
import type { HealthStatus } from '#health/interfaces/health-status.interface'

/**
 * Service for managing the application configuration.
 */
export class ConfigService implements HealthCheck {
  private config: ArchesConfig
  private readonly health: HealthStatus = {
    status: 'COMPLETED'
  }

  constructor(config: ArchesConfig) {
    this.config = config
  }

  public get<Path extends Leaves<ArchesConfig>>(
    propertyPath: Path
  ): LeafTypes<ArchesConfig, Path> {
    return propertyPath
      .split('.')
      .reduce<unknown>(
        (acc, key) =>
          acc && typeof acc === 'object' ?
            (acc as Record<string, unknown>)[key]
          : undefined,
        this.config
      ) as LeafTypes<ArchesConfig, Path>
  }

  public getConfig(): ArchesConfig {
    return this.config
  }

  public getHealth(): HealthStatus {
    return this.health
  }

  public setConfig(config: ArchesConfig): void {
    this.config = config
  }
}
