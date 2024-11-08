import { Injectable, Logger } from "@nestjs/common";
import { ConfigService } from "@nestjs/config";
import { PassportStrategy } from "@nestjs/passport";
import { AuthProviderType } from "@prisma/client";
import admin from "firebase-admin";
import { ExtractJwt, Strategy } from "passport-firebase-jwt";

import { firebaseConfig } from "../../firebase.config";
import { UsersService } from "../../users/users.service";
import { CurrentUserDto } from "../decorators/current-user.decorator";

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
  type: firebaseConfig.type,
};

@Injectable()
export class FirebaseStrategy extends PassportStrategy(
  Strategy,
  "firebase-auth"
) {
  private readonly logger = new Logger("Firebase Strategy");

  constructor(
    private configService: ConfigService,
    private usersService: UsersService
  ) {
    super({
      jwtFromRequest: ExtractJwt.fromAuthHeaderAsBearerToken(),
    });

    const useLocalIdentityToolkit =
      this.configService.get("NODE_ENV") !== "production";

    if (!admin.apps.length) {
      admin.initializeApp({
        credential: admin.credential.cert(firebase_params),
        projectId: useLocalIdentityToolkit
          ? "filechat-io"
          : firebase_params.projectId,
      });
    }
  }

  async validate(token: string): Promise<CurrentUserDto> {
    this.logger.log(`Validating Firebase Token: ${token}`);
    const decodedToken = await admin.auth().verifyIdToken(token);
    let user: CurrentUserDto;
    try {
      const user = await this.usersService.findOneByEmail(decodedToken.email);
      return user;
    } catch (e) {
      const username =
        user.email.split("@")[0] +
        "-" +
        Math.random().toString(36).substring(2, 6);
      user = await this.usersService.create(null, {
        email: decodedToken.email,
        emailVerified: true,
        password: null,
        photoUrl: decodedToken.picture,
        // plus - and a random string of 4 letters
        username,
      });
    } finally {
      user = await this.usersService.syncAuthProvider(
        decodedToken.email,
        AuthProviderType.FIREBASE,
        decodedToken.uid
      );
      return user;
    }
  }
}
