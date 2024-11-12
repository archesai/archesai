import { useToast } from "@/components/ui/use-toast";
import { baseUrl } from "@/generated/archesApiFetcher";
import { TokenDto, UserEntity } from "@/generated/archesApiSchemas";
import { GoogleAuthProvider, signInWithPopup } from "firebase/auth";
import { useAtom } from "jotai";
import { useRouter } from "next/navigation";
import { useCallback } from "react";

import { auth } from "../lib/firebase";
import { AuthState, authStateAtom } from "../state/authState";

export const useAuth = () => {
  const [authState, setAuthState] = useAtom(authStateAtom);
  const router = useRouter();
  const { toast } = useToast();

  const logout = useCallback(async () => {
    const response = await fetch(baseUrl + "/auth/logout", {
      credentials: "include",
      method: "POST",
      mode: "cors",
    });
    if (response.status !== 201) {
      console.error("Failed to logout");
    }
    setAuthState({
      defaultOrgname: "",
      memberships: [],
      status: "Unauthenticated",
      user: null,
    });
    router.push("/");
  }, [setAuthState, router]);

  const getNewRefreshToken = async (
    authState: AuthState,
    setAuthState: any
  ) => {
    console.log("Getting new refresh token");
    if (authState.status === "Refreshing") {
      console.log("Already refreshing token, skipping");
      return;
    }

    setAuthState((prev: AuthState) => ({ ...prev, status: "Refreshing" }));

    try {
      const response = await fetch(baseUrl + "/auth/refresh-token", {
        credentials: "include",
        method: "POST",
        mode: "cors",
      });

      const data = (await response.json()) as TokenDto;
      if (response.status !== 201) {
        throw new Error("Failed to refresh token");
      }

      setAuthState((prev: AuthState) => ({
        ...prev,
        status: "Authenticated",
      }));
      console.log("Got new refresh token");
      return data.accessToken;
    } catch (error) {
      console.error("Error refreshing token:", error);
      setAuthState((prev: AuthState) => ({
        ...prev,
        status: "Unauthenticated",
      }));
      console.log("Logging out due to error refreshing token");
      await logout();
      return null;
    }
  };

  const getUserFromToken = useCallback(async () => {
    console.log("Getting user from token");
    try {
      const response = await fetch(baseUrl + "/user", {
        credentials: "include", // Include cookies
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
      toast({
        description: "An error occurred. Please log in again.",
        variant: "destructive",
      });
      console.error("Error in getUserFromToken: ", error);
      await logout();
    }
  }, [logout, router, setAuthState]);

  const signInWithEmailAndPassword = useCallback(
    async (email: string, password: string) => {
      try {
        const result = await fetch(baseUrl + "/auth/login", {
          body: JSON.stringify({ email, password }),
          credentials: "include", // Include cookies
          headers: { "Content-Type": "application/json" },
          method: "POST",
          mode: "cors",
        });
        if (result.status !== 201) {
          throw new Error("Invalid credentials");
        }
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
    [router]
  );

  const signInWithGoogle = useCallback(async () => {
    try {
      const provider = new GoogleAuthProvider();
      const credential = await signInWithPopup(auth, provider);
      const token = await credential.user.getIdToken();
      await fetch(baseUrl + "/auth/firebase/callback", {
        body: JSON.stringify({ accessToken: token }),
        credentials: "include", // Include cookies
        headers: {
          "Content-Type": "application/json",
        },
        method: "POST",
        mode: "cors",
      });

      router.push("/playground");
    } catch (error) {
      await logout();
      toast({
        description: "Google sign-in failed.",
        variant: "destructive",
      });
      console.error("Error signing in with Google: ", error);
    }
  }, [router, logout]);

  const registerWithEmailAndPassword = useCallback(
    async (email: string, password: string) => {
      try {
        const result = await fetch(baseUrl + "/auth/register", {
          body: JSON.stringify({ email, password }),
          credentials: "include", // Include cookies
          headers: { "Content-Type": "application/json" },
          method: "POST",
          mode: "cors",
        });
        if (result.status !== 201) {
          throw new Error("Could not register user");
        }
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
    [router]
  );

  return {
    ...authState,
    getNewRefreshToken: () => getNewRefreshToken(authState, setAuthState),
    getUserFromToken,
    logout,
    registerWithEmailAndPassword,
    signInWithEmailAndPassword,
    signInWithGoogle,
  };
};
