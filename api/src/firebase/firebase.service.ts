import { Injectable } from "@nestjs/common";
import { ConfigService } from "@nestjs/config";
import admin from "firebase-admin";

import { firebaseConfig } from "../firebase.config";

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
export class FirebaseService {
  public firebase: admin.app.App;

  constructor(private readonly configService: ConfigService) {
    const useLocalIdentityToolkit =
      this.configService.get("NODE_ENV") !== "production";

    try {
      this.firebase = admin.initializeApp({
        credential: admin.credential.cert(firebase_params),
        projectId: useLocalIdentityToolkit
          ? "filechat-io"
          : firebase_params.projectId,
      });
    } catch (err) {
      this.firebase = admin.app();
    }
  }
}
