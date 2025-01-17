import { Injectable } from '@nestjs/common'

import { BaseService } from '../common/base.service'
import { WebsocketsService } from '../websockets/websockets.service'
import { MemberEntity, MemberModel } from './entities/member.entity'
import { MemberRepository } from './member.repository'
import { UsersService } from '../users/users.service'

@Injectable()
export class MembersService extends BaseService<
  MemberEntity,
  MemberModel,
  MemberRepository
> {
  constructor(
    private usersService: UsersService,
    private memberRepository: MemberRepository,
    private websocketsService: WebsocketsService
  ) {
    super(memberRepository)
  }

  async create(data: Pick<MemberEntity, 'inviteEmail' | 'orgname' | 'role'>) {
    let membership = await this.memberRepository.create({
      ...data
    })
    const existingUser = await this.usersService.findOneByEmail(
      data.inviteEmail
    )
    if (existingUser) {
      membership = await this.memberRepository.findByInviteEmailAndOrgname(
        membership.inviteEmail,
        membership.orgname
      )
      membership = await this.memberRepository.update(membership.id, {
        user: {
          connect: {
            id: existingUser.id
          }
        }
      })
    }
    return this.toEntity(membership)
  }

  async join(orgname: string, inviteEmail: string, username: string) {
    this.logger.log(
      {
        orgname,
        inviteEmail,
        username
      },
      `attempting join`
    )
    return this.toEntity(
      await this.memberRepository.join(orgname, inviteEmail, username)
    )
  }

  protected emitMutationEvent(entity: MemberEntity): void {
    this.websocketsService.socket?.to(entity.orgname).emit('update', {
      queryKey: ['organizations', entity.orgname, 'members']
    })
  }

  protected toEntity(model: MemberModel): MemberEntity {
    return new MemberEntity(model)
  }
}
