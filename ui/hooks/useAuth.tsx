import { useToast } from "@/components/ui/use-toast";
// hooks/useAuth.ts
import { baseUrl } from "@/generated/archesApiFetcher";
import { TokenDto, UserEntity } from "@/generated/archesApiSchemas";
import { GoogleAuthProvider, signInWithPopup } from "firebase/auth";
import { useAtom } from "jotai";
import { jwtDecode } from "jwt-decode";
import { useRouter } from "next/navigation";
import { useCallback, useEffect, useState } from "react";

import { auth } from "../lib/firebase";
import {
  accessTokenAtom,
  authStateAtom,
  refreshTokenAtom,
} from "../state/authState";

interface DecodedToken {
  exp: number;
}

export const useAuth = () => {
  const [authState, setAuthState] = useAtom(authStateAtom);
  const [accessToken, setAccessToken] = useAtom(accessTokenAtom);
  const [refreshToken, setRefreshToken] = useAtom(refreshTokenAtom);
  const router = useRouter();
  const { toast } = useToast();
  const [authStatus, setAuthStatus] = useState<
    "Authenticated" | "Loading" | "Unauthenticated"
  >("Loading");

  // Prevent multiple refreshes
  let isRefreshing = false;
  let refreshSubscribers: ((token: string) => void)[] = [];

  const logout = useCallback(async () => {
    setAuthState({
      defaultOrgname: "",
      isLoading: false,
      memberships: [],
      user: null,
    });
    setAccessToken("");
    setRefreshToken("");
    setAuthStatus("Unauthenticated");
    router.push("/");
  }, [setAuthState, setAccessToken, setRefreshToken, router]);

  const subscribeTokenRefresh = (cb: (token: string) => void) => {
    refreshSubscribers.push(cb);
  };

  const onRefreshed = (token: string) => {
    refreshSubscribers.forEach((cb) => cb(token));
    refreshSubscribers = [];
  };

  const getNewRefreshToken = useCallback(async (): Promise<null | string> => {
    if (isRefreshing) {
      return new Promise<string>((resolve) => {
        subscribeTokenRefresh(resolve);
      });
    }

    isRefreshing = true;
    try {
      const response = await fetch(baseUrl + "/auth/refresh-token", {
        body: JSON.stringify({ refreshToken }),
        headers: { "Content-Type": "application/json" },
        method: "POST",
        mode: "cors",
      });

      const data = (await response.json()) as TokenDto;
      if (response.status !== 201) {
        await logout();
        toast({
          description: "An error occurred. Please log in again.",
          variant: "destructive",
        });
        return null;
      }

      setAccessToken(data.accessToken);
      setRefreshToken(data.refreshToken);
      onRefreshed(data.accessToken);
      return data.accessToken;
    } catch (error) {
      console.error("Error refreshing token:", error);
      await logout();
      toast({
        description: "An error occurred. Please log in again.",
        variant: "destructive",
      });
      return null;
    } finally {
      isRefreshing = false;
    }
  }, [refreshToken, setAccessToken, setRefreshToken, logout]);

  const getUserFromToken = useCallback(async () => {
    try {
      const decoded: DecodedToken = jwtDecode(accessToken);
      const now = Math.floor(Date.now() / 1000);
      let currentToken = accessToken;

      if (decoded.exp < now) {
        const newToken = await getNewRefreshToken();
        if (newToken) {
          currentToken = newToken;
        } else {
          return;
        }
      }

      const response = await fetch(baseUrl + "/user", {
        headers: {
          Authorization: "Bearer " + currentToken,
        },
        method: "GET",
        mode: "cors",
      });
      console.log("HITTING USER ENDPOINT", new Date());

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
      setAuthStatus("Authenticated");
      router.push("/playground");
    } catch (error) {
      await logout();
      toast({
        description: "An error occurred. Please log in again.",
        variant: "destructive",
      });
      console.error("Error in getUserFromToken: ", error);
    }
  }, [accessToken, getNewRefreshToken, logout, router, setAuthState]);

  const signInWithEmailAndPassword = useCallback(
    async (email: string, password: string) => {
      try {
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
        setAuthStatus("Authenticated");
        router.push("/playground");
      } catch (error) {
        toast({
          description: "Invalid credentials.",
          variant: "destructive",
        });
        console.error("Error signing in with email and password:", error);
        throw error;
      }
    },
    [router, setAccessToken, setRefreshToken]
  );

  const signInWithGoogle = useCallback(async () => {
    try {
      const provider = new GoogleAuthProvider();
      const credential = await signInWithPopup(auth, provider);
      const token = await credential.user.getIdToken();
      const response = await fetch(baseUrl + "/auth/firebase/callback", {
        body: JSON.stringify({ accessToken: token }),
        headers: {
          Authorization: "Bearer " + token,
          "Content-Type": "application/json",
        },
        method: "POST",
        mode: "cors",
      });

      const data = (await response.json()) as TokenDto;
      setAccessToken(data.accessToken);
      setRefreshToken(data.refreshToken);
      setAuthStatus("Authenticated");
      router.push("/playground");
    } catch (error) {
      await logout();
      toast({
        description: "Google sign-in failed.",
        variant: "destructive",
      });
      console.error("Error signing in with Google: ", error);
    }
  }, [router, setAccessToken, setRefreshToken, logout]);

  const registerWithEmailAndPassword = useCallback(
    async (email: string, password: string) => {
      try {
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
        setAuthStatus("Authenticated");
        router.push("/playground");
      } catch (error) {
        toast({
          description: "Registration failed.",
          variant: "destructive",
        });
        console.error("Error registering user:", error);
        throw error;
      }
    },
    [router, setAccessToken, setRefreshToken]
  );

  useEffect(() => {
    if (accessToken) {
      getUserFromToken();
    } else {
      setAuthStatus("Unauthenticated");
    }
  }, [accessToken, getUserFromToken]);

  return {
    ...authState,
    accessToken,
    authStatus,
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
