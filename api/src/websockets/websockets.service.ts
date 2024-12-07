import { Injectable } from '@nestjs/common'
import { Server } from 'socket.io'

@Injectable()
export class WebsocketsService {
  public socket: Server = null

  constructor() {}
}
