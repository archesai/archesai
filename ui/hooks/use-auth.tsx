import { useToast } from "@/components/ui/use-toast";
import { baseUrl } from "@/generated/archesApiFetcher";
import { TokenDto, UserEntity } from "@/generated/archesApiSchemas";
import { GoogleAuthProvider, signInWithPopup } from "firebase/auth";
import { useAtom } from "jotai";
import { useRouter } from "next/navigation";
import { useCallback } from "react";

import { auth } from "../lib/firebase";
import { authStatusAtom, defaultOrgnameAtom } from "../state/authState";

export const useAuth = () => {
  const [defaultOrgname, setDefaultOrgname] = useAtom(defaultOrgnameAtom);
  const [status, setStatus] = useAtom(authStatusAtom);

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
      return;
    }
    router.push("/");
    setStatus("Unauthenticated");
  }, [setStatus]);

  const getNewRefreshToken = async () => {
    if (status === "Refreshing") {
      console.log("Already refreshing token, skipping");
      return;
    }
    console.log("Getting new refresh token");
    setStatus("Refreshing");

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

      setStatus("Authenticated");
      console.log("Got new refresh token");
      return data.accessToken;
    } catch (error) {
      console.log("Error refreshing token, logging out:", error);
      await logout();
      return null;
    }
  };

  const authenticate = useCallback(async () => {
    console.log("Attempting to authenticate");
    setStatus("Loading");
    try {
      let response = await fetch(baseUrl + "/user", {
        credentials: "include", // Include cookies
        method: "GET",
        mode: "cors",
      });

      if (response.status === 401) {
        await getNewRefreshToken();
      }

      response = await fetch(baseUrl + "/user", {
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
      setDefaultOrgname(user.defaultOrgname);
      setStatus("Authenticated");
    } catch (error) {
      toast({
        description: "An error occurred. Please log in again.",
        variant: "destructive",
      });
      console.error("Error in getUserFromToken: ", error);
      await logout();
    }
  }, [logout, router, setStatus, setDefaultOrgname]);

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
    authenticate,
    defaultOrgname,
    getNewRefreshToken,
    logout,
    registerWithEmailAndPassword,
    setStatus,
    signInWithEmailAndPassword,
    signInWithGoogle,
    status,
  };
};
