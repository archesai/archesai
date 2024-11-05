import { baseUrl } from "@/generated/archesApiFetcher";
import { TokenDto, UserEntity } from "@/generated/archesApiSchemas";
import { GoogleAuthProvider, signInWithPopup } from "firebase/auth";
import { useAtom } from "jotai"; // The auth state atom
import { useRouter } from "next/navigation";

import { auth } from "../lib/firebase";
import {
  accessTokenAtom,
  authStateAtom,
  refreshTokenAtom,
} from "../state/authState";

export const useAuth = () => {
  const [authState, setAuthState] = useAtom(authStateAtom);
  const [accessToken, setAccessToken] = useAtom(accessTokenAtom);
  const [refreshToken, setRefreshToken] = useAtom(refreshTokenAtom);
  const router = useRouter();

  const getUserFromToken = async () => {
    console.log("Getting user from token");
    try {
      // get expiration
      const [, payload] = accessToken.split(".");
      const data = JSON.parse(atob(payload));
      const exp = data.exp;
      const now = Math.floor(Date.now() / 1000);
      if (exp < now) {
        await getNewRefreshToken();
      }

      const response = await fetch(baseUrl + "/user", {
        headers: {
          Authorization: "Bearer " + accessToken,
        },
        method: "GET",
        mode: "cors",
      });

      if (response.status !== 200) {
        console.error("Error loading user");
        await logout();
        router.push("/");
        return;
      }

      const user = (await response.json()) as UserEntity;
      setAuthState((prevState) => ({
        ...prevState,
        defaultOrgname: user.defaultOrgname,
        isLoading: false,
        memberships: user.memberships,
        user,
      }));
    } catch (error) {
      await logout();
      console.error("Error in getUserFromToken: ", error);
    }
  };

  const signInWithEmailAndPassword = async (
    email: string,
    password: string
  ) => {
    const result = await fetch(baseUrl + "/auth/login", {
      body: JSON.stringify({ email, password }),
      headers: { "Content-Type": "application/json" },
      method: "POST",
      mode: "cors",
    });
    if (result.status !== 201) {
      throw new Error("Invalid credentials");
    }
    const data = (await result.json()) as TokenDto;
    setAccessToken(data.accessToken);
    setRefreshToken(data.refreshToken);
    router.push("/playground");
  };

  const getNewRefreshToken = async () => {
    console.log("Getting new refresh token");
    const response = await fetch(baseUrl + "/auth/refresh-token", {
      body: JSON.stringify({ refreshToken }),
      headers: { "Content-Type": "application/json" },
      method: "POST",
      mode: "cors",
    });

    const data = (await response.json()) as TokenDto;
    if (response.status !== 201) {
      await logout();
      console.error("Error refreshing token");
      return;
    }

    setAccessToken(data.accessToken);
    setRefreshToken(data.refreshToken);
  };

  const signInWithGoogle = async () => {
    try {
      const provider = new GoogleAuthProvider();
      const credential = await signInWithPopup(auth, provider);
      const accessToken = await credential.user.getIdToken();
      const response = await fetch(baseUrl + "/auth/firebase/callback", {
        body: JSON.stringify({ accessToken }),
        headers: {
          Authorization: "Bearer " + accessToken,
          "Content-Type": "application/json",
        },
        method: "POST",
        mode: "cors",
      });

      const data = (await response.json()) as TokenDto;
      setAccessToken(data.accessToken);
      setRefreshToken(data.refreshToken);
      router.push("/playground");
    } catch (error) {
      await logout();
      console.error("Error signing in with Google: ", error);
    }
  };

  const registerWithEmailAndPassword = async (
    email: string,
    password: string
  ) => {
    const result = await fetch(baseUrl + "/auth/register", {
      body: JSON.stringify({
        email,
        password,
      }),
      headers: { "Content-Type": "application/json" },
      method: "POST",
      mode: "cors",
    });
    if (result.status !== 201) {
      throw new Error("Could not register user");
    }
    const data = (await result.json()) as TokenDto;
    setAccessToken(data.accessToken);
    setRefreshToken(data.refreshToken);
    router.push("/playground");
  };

  const logout = async () => {
    setAuthState((prevState) => ({
      ...prevState,
      defaultOrgname: "",
      isLoading: false,
      memberships: [],
      user: null,
    }));
    setAccessToken("");
    setRefreshToken("");
    router.push("/");
  };

  return {
    ...authState,
    accessToken,
    getNewRefreshToken,
    getUserFromToken,
    logout,
    refreshToken,
    registerWithEmailAndPassword,
    setAccessToken,
    setRefreshToken,
    signInWithEmailAndPassword,
    signInWithGoogle,
  };
};
