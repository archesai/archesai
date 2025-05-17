export type LeafTypes<T, S extends string> = T extends unknown // distribute across unions explicitly
  ? S extends `${infer Head}.${infer Tail}`
    ? Head extends keyof T
      ? LeafTypes<NonNullable<T[Head]>, Tail>
      : never
    : S extends keyof T
      ? T[S]
      : never
  : never

export type Leaves<T> = T extends object
  ? {
      [K in keyof T & string]: NonNullable<T[K]> extends object
        ? // If it's an object, include both "K" and deeper paths "K.xxx"
          `${K}.${Leaves<NonNullable<T[K]>>}` | K
        : // Otherwise it's a leaf, just "K"
          K
    }[keyof T & string]
  : never
