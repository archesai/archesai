// src/auth/twitter.strategy.ts

import { Injectable } from "@nestjs/common";
import { PassportStrategy } from "@nestjs/passport";
import { AuthProviderType } from "@prisma/client";
import { Profile, Strategy } from "passport-twitter";

import { UsersService } from "../../users/users.service";
import { CurrentUserDto } from "../decorators/current-user.decorator";

@Injectable()
export class TwitterStrategy extends PassportStrategy(Strategy, "twitter") {
  constructor(private readonly usersService: UsersService) {
    super({
      callbackURL:
        "https://grizzly-content-worm.ngrok-free.app/auth/twitter/callback", // Your callback URL
      consumerKey: "IM7VJZOzMCrcT6QNqy1Phpjpu", // Twitter API Key
      consumerSecret: "W3pPP5tZI6XZzFmyv9NA3DldgUneLDzQ4F8D8yFKEajFVD5lD5", // Twitter API Secret Key
      includeEmail: true, // Request email from Twitter
    });
  }

  async validate(
    token: string,
    tokenSecret: string,
    profile: Profile,
    cb: any
  ) {
    try {
      const twitterId = profile.id;
      const email = profile.emails?.[0]?.value;
      const username = profile.username;

      let user: CurrentUserDto;
      try {
        user = await this.usersService.findOneByEmail(email);
      } catch (e) {
        user = await this.usersService.create({
          email: email,
          emailVerified: true,
          password: null,
          photoUrl: profile.photos?.[0]?.value,
          username,
        });
      } finally {
        user = await this.usersService.syncAuthProvider(
          email,
          AuthProviderType.TWITTER,
          twitterId
        );
        return cb(null, user);
      }
    } catch (err) {
      return cb(err, false);
    }
  }
}
