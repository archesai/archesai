// TypeSystem utilities for Zod

import { z } from 'zod'

// eslint-disable-next-line @typescript-eslint/unbound-method
const originalAdd = z.globalRegistry.add

// No global configuration needed for Zod like TypeBox had
z.globalRegistry.add = (
  schema: Parameters<typeof originalAdd>[0],
  meta: Parameters<typeof originalAdd>[1]
) => {
  if (!meta.id) {
    return originalAdd.call(z.globalRegistry, schema, meta)
  }
  const existingSchema = z.globalRegistry._idmap.get(meta.id)
  if (existingSchema) {
    z.globalRegistry.remove(existingSchema)
    z.globalRegistry._idmap.delete(meta.id)
  }
  return originalAdd.call(z.globalRegistry, schema, meta)
}
