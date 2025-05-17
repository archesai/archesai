// eslint-disable-next-line @typescript-eslint/no-explicit-any
export type Type<T = unknown> = new (...args: any[]) => T
