import { Injectable } from "@nestjs/common";
import { PassportStrategy } from "@nestjs/passport";
import { AuthProviderType } from "@prisma/client";
import { ExtractJwt, Strategy } from "passport-firebase-jwt";

import { FirebaseService } from "../../firebase/firebase.service";
import { UsersService } from "../../users/users.service";
import { CurrentUserDto } from "../decorators/current-user.decorator";

@Injectable()
export class FirebaseStrategy extends PassportStrategy(
  Strategy,
  "firebase-auth"
) {
  constructor(
    private firebaseService: FirebaseService,
    private usersService: UsersService
  ) {
    super({
      jwtFromRequest: ExtractJwt.fromAuthHeaderAsBearerToken(),
    });
  }

  async validate(token: string): Promise<CurrentUserDto> {
    const decodedToken = await this.firebaseService.firebase
      .auth()
      .verifyIdToken(token);
    let user: CurrentUserDto;
    try {
      const user = await this.usersService.findOneByEmail(decodedToken.email);
      return user;
    } catch (e) {
      user = await this.usersService.create({
        email: decodedToken.email,
        emailVerified: true,
        password: null,
        photoUrl: decodedToken.picture,
        username: decodedToken.email,
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
