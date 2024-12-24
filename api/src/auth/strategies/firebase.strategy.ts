import { UserEntity } from '@/src/users/entities/user.entity'
import { Injectable, Logger, UnauthorizedException } from '@nestjs/common'
import { ConfigService } from '@nestjs/config'
import { PassportStrategy } from '@nestjs/passport'
import { AuthProviderType } from '@prisma/client'
import admin from 'firebase-admin'
import { ExtractJwt, Strategy } from 'passport-firebase-jwt'

import { UsersService } from '../../users/users.service'

export const firebaseConfig = {
  auth_provider_x509_cert_url: 'https://www.googleapis.com/oauth2/v1/certs',
  auth_uri: 'https://accounts.google.com/o/oauth2/auth',
  client_email: 'firebase-adminsdk-60lur@filechat-io.iam.gserviceaccount.com',
  client_id: '110829470716646406677',
  client_x509_cert_url:
    'https://www.googleapis.com/robot/v1/metadata/x509/firebase-adminsdk-60lur%40filechat-io.iam.gserviceaccount.com',
  private_key:
    '-----BEGIN PRIVATE KEY-----\nMIIEvAIBADANBgkqhkiG9w0BAQEFAASCBKYwggSiAgEAAoIBAQCkFifUh+Girwsp\n+geLoC/sIClWd4V2Ww1bRgU9t7n22sEliyIze4+cWeab7HA5aeH8y/8KdSXibo2g\nxOrWHMf5Zdy41mrb5b45IHfkyGBPnuku1HIoHIBYk+HD4UXvdPS8/9rZTNKZG48+\nAs58lnl/p3/2v3Rk96syJpibOo8ZJp3H5s1EiRXIkRHqkny9cu1P70bSaw9P6wZR\ncQB72qKlx7tuclmTI5tMrAZdlGM1fOjc/2EKxZ0WjEfS83B+7idshCDlj+sgu5W0\nV9ygjEfEVFIFZHUiVntsOsBATzvGVHEZWQlg1al/Lso0zVkdvYUj9TmOq2oP8f6p\nWCSE3WbZAgMBAAECggEAHGZ+Td2X+ve3o7QUPsFC0iiN5xqoVcA0O15O9Wv5qsWZ\npSUkDgueo3X3AXlmwjU2qXHoipeUr1CpwFAdAmK4ZQ0Rq0dndviYmFwrjER3UQik\n/RdFy/YE6+/qpWP2Hhhc7OKO70oJ8Hix7g5/zVYhIOxtrFhebcRlU/iUtNdpobVD\nmOzxY4yqwQjpA2er+0sF8xQ5dSyt9g13YOH0T+AvSAoMDrgSwvtl0c1PldMQMQaM\n1IszYpKtrQi5rneNMf9gvclxwiVSinhk1itkp721mGAQTaAvAxh1FB0zdJRTdnqC\nF7C9OIGWbJZOEIlo1e+jg6QoJmW9cyI4TdALvLYpJQKBgQDWMbi4vjlgcwPajVt4\nMGkYJHfbve2OR0HKweaHBDZReiRj5dzocB5T0SPR1udFHFrW9AkjpKS5PTtAuB8k\nrzm0P9tf+CFO/8mi/SgzQxBsz5hmTo7m3qoC4KmXIm6AX3C1Q9J8tJ0motZpBnse\nY76M+VZ8AtTN7Kgt9aQ5Rt66TQKBgQDEHMnxPS68DKVkRK825urUU6mkxHKwz80o\nWNglKqZIWje3XVbu7Tmb1svLcbN8tu0vGDbFOKEEWp/5HReXegtmoNfdRAxQcy+n\nmzC/swV8ve8nIIomuxKvbh1PMdg5BzYM+X93Cnt8F4UxClBTQxD2l3bLEyfjXHeo\nrQ9NoRtMvQKBgDLSbV/4Uqjd4WYz8CYeZnFCBeZvtDP0GFpBk68pgrHmZ0gEvFuy\nbp+4meUqNomhZrRmBt0cLbF+I9cBWPJdWTW5iRXGTDDwZCl2I9m16enHgAOWVDXX\nU0OHhvXDR7DR9G4t/31zZW5LaNBWp1PYmtfcOXcHPPL3Whg9lo+4jxRpAoGAKbU9\ntZKfh9rgqex5nyGJO9L3N1WYVsY7CaOrhGwHpUeapeKyBGprYBtUiFYMKC/3TZbG\nvzcF95kWgLKRO+P23MLEZgh83fdBYVH+EicOubLjU9z1xLrwhGLU1Ozy4V4JPsUp\nOLYAASo3Z3CcGLkguHEKELJoP1CBGlyD7qye07kCgYArEutvN6C7CnBZ1YiCQHVG\nZVvl810gHTqZ+NA4dryehXQFqtQx//HO0F/N46qC90OICIeMaJivUrQbkdOwPgvL\n1pbczgtmJ+b0ke2Ex4y/pRQ/ykl/ebOpI++rwdHZRS51PFrBn73eGZ/wu/lzOp1b\nmG/NeEjF7j3cA1KfGmtPwg==\n-----END PRIVATE KEY-----\n',
  private_key_id: '3d18328ad51f3ad682d65d249c1b4a53718ca8e5',
  project_id: 'filechat-io',
  token_uri: 'https://oauth2.googleapis.com/token',
  type: 'service_account'
}

const firebase_params = {
  authProviderX509CertUrl: firebaseConfig.auth_provider_x509_cert_url,
  authUri: firebaseConfig.auth_uri,
  clientC509CertUrl: firebaseConfig.client_x509_cert_url,
  clientEmail: firebaseConfig.client_email,
  clientId: firebaseConfig.client_id,
  privateKey: firebaseConfig.private_key,
  privateKeyId: firebaseConfig.private_key_id,
  projectId: firebaseConfig.project_id,
  tokenUri: firebaseConfig.token_uri,
  type: firebaseConfig.type
}

@Injectable()
export class FirebaseStrategy extends PassportStrategy(
  Strategy,
  'firebase-auth'
) {
  private readonly logger = new Logger(FirebaseStrategy.name)

  constructor(
    private configService: ConfigService,
    private usersService: UsersService
  ) {
    super({
      jwtFromRequest: ExtractJwt.fromAuthHeaderAsBearerToken()
    })

    const useLocalIdentityToolkit =
      this.configService.get('NODE_ENV') !== 'production'

    if (!admin.apps.length) {
      this.logger.log('Initializing Firebase Admin SDK')
      admin.initializeApp({
        credential: admin.credential.cert(firebase_params),
        projectId: useLocalIdentityToolkit
          ? 'filechat-io'
          : firebase_params.projectId
      })
    }
  }

  async validate(token: string): Promise<UserEntity> {
    this.logger.debug(`Validating Firebase Token: ${token}`)
    const decodedToken = await admin.auth().verifyIdToken(token)
    if (!decodedToken.email) {
      throw new UnauthorizedException('Token does not contain email')
    }

    try {
      await this.usersService.findOneByEmail(decodedToken.email)
      return this.usersService.syncAuthProvider(
        decodedToken.email,
        AuthProviderType.FIREBASE,
        decodedToken.uid
      )
    } catch (e) {
      this.logger.debug(`User not found: ${decodedToken.email}: ${e}`)
      const username =
        decodedToken.email.split('@')[0] +
        '-' +
        Math.random().toString(36).substring(2, 6)
      await this.usersService.create({
        email: decodedToken.email,
        emailVerified: true,
        photoUrl: decodedToken.picture || '',
        username
      })
      return this.usersService.syncAuthProvider(
        decodedToken.email,
        AuthProviderType.FIREBASE,
        decodedToken.uid
      )
    }
  }
}
