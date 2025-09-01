import z from 'zod'

export const StringSchema: z.ZodString = z.string().describe('A string value')
