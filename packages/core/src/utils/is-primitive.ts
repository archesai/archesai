export const isUndefined = (obj: unknown): obj is undefined =>
  typeof obj === 'undefined'

export const isObject = (fn: unknown): fn is object =>
  !isNil(fn) && typeof fn === 'object'

export const isNil = (val: unknown): val is null | undefined =>
  isUndefined(val) || val === null

export const isString = (val: unknown): val is string => typeof val === 'string'

export const isEmpty = (array: undefined | unknown[]): boolean =>
  !(array && array.length > 0)
