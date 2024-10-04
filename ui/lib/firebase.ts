import { getAnalytics } from "firebase/analytics";
import { initializeApp } from "firebase/app";
import { connectAuthEmulator, getAuth } from "firebase/auth";

const firebaseConfig = {
  apiKey: "AIzaSyAiWVxGSbPEX6jhwCbIL47q_z8Jm3Lucvk",
  appId: "1:153920903783:web:fbace847f681d7a832d2c2",
  authDomain: "filechat-io.firebaseapp.com",
  measurementId: "G-LH12KBG1LS",
  messagingSenderId: "153920903783",
  projectId: "filechat-io",
  storageBucket: "filechat-io.appspot.com",
};

const app = initializeApp(firebaseConfig);
const auth = getAuth(app);
let analytics = null;

if (typeof window !== "undefined" && !window["_init" as any]) {
  if (location) {
    analytics = getAnalytics();
    if (process.env.NEXT_PUBLIC_USE_FIREBASE_AUTH_EMULATOR == "true") {
      connectAuthEmulator(
        auth,
        process.env.NEXT_PUBLIC_FIREBASE_AUTH_EMULATOR_HOST as string,
        {
          disableWarnings: true,
        }
      );
    }
    window["_init" as any] = true as any;
  }
}

export { analytics, auth };
