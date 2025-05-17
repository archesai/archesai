import type { Type } from '#types/type'

/* eslint-disable @typescript-eslint/no-unnecessary-type-parameters */
export interface ArgumentsHost {
  getArgByIndex<T = unknown>(index: number): T
  getArgs<T extends unknown[] = unknown[]>(): T
  getType<TContext extends string = ContextType>(): TContext
  switchToHttp(): HttpArgumentsHost
  switchToRpc(): RpcArgumentsHost
  switchToWs(): WsArgumentsHost
}

export type ContextType = 'http' | 'rpc' | 'ws'

export interface ExecutionContext extends ArgumentsHost {
  getClass<T>(): Type<T>
  // eslint-disable-next-line @typescript-eslint/no-unsafe-function-type
  getHandler(): Function
}

export interface HttpArgumentsHost {
  getNext<T = unknown>(): T
  getRequest<T = unknown>(): T
  getResponse<T = unknown>(): T
}

export interface RpcArgumentsHost {
  getContext<T = unknown>(): T
  getData<T = unknown>(): T
}

export interface WsArgumentsHost {
  getClient<T = unknown>(): T
  getData<T = unknown>(): T
  getPattern(): string
}
