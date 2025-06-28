/* eslint-disable @typescript-eslint/no-unsafe-function-type */
/* eslint-disable @typescript-eslint/no-explicit-any */

import type { Type } from '#types/type'

import 'reflect-metadata'

export interface DynamicModule extends ModuleMetadata {
  global?: boolean
  module: Type
}

export interface ModuleMetadata {
  controllers?: Type[]
  exports?: (
    | Abstract<any>
    | DynamicModule
    | Function
    | Provider
    | string
    | symbol
  )[]
  imports?: (DynamicModule | Promise<DynamicModule> | Type)[]
  providers?: Provider[]
}

export type Provider<T = any> =
  | ClassProvider<T>
  | ExistingProvider<T>
  | FactoryProvider<T>
  | Type
  | ValueProvider<T>

interface Abstract<T> extends Function {
  prototype: T
}

interface ClassProvider<T = any> {
  durable?: boolean
  inject?: never
  provide: InjectionToken
  useClass: Type<T>
}

interface ExistingProvider<T = any> {
  provide: InjectionToken
  useExisting: T
}

interface FactoryProvider<T = any> {
  durable?: boolean
  inject?: InjectionToken[]
  provide: InjectionToken
  useFactory: (...args: any[]) => Promise<T> | T
}

type InjectionToken<T = any> = Abstract<T> | string | symbol | Type<T>

interface ValueProvider<T = any> {
  inject?: never
  provide: InjectionToken
  useValue: T
}

export function createModule<
  T extends DynamicModule | Promise<DynamicModule> | Type
>(target: T, metadata: ModuleMetadata, global = true): T {
  for (const property in metadata) {
    if (Object.hasOwnProperty.call(metadata, property)) {
      Reflect.defineMetadata(
        property,
        metadata[property as keyof typeof metadata],
        target
      )
    }
  }
  if (global) {
    Reflect.defineMetadata('__module:global__', true, target)
  }
  return target
}

export function Global(): ClassDecorator {
  return (target: Function) => {
    Reflect.defineMetadata('__module:global__', true, target)
  }
}

export function Module(metadata: ModuleMetadata): ClassDecorator {
  return (target: Function) => {
    for (const property in metadata) {
      if (Object.hasOwnProperty.call(metadata, property)) {
        Reflect.defineMetadata(
          property,
          metadata[property as keyof typeof metadata],
          target
        )
      }
    }
  }
}
