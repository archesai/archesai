import { MemberEntity, UserEntity } from "@/generated/archesApiSchemas";
import { atom } from "jotai";

// Define the initial authentication state
export interface AuthState {
  defaultOrgname: string;
  memberships: MemberEntity[];
  status: "Authenticated" | "Loading" | "Refreshing" | "Unauthenticated";
  user: null | UserEntity;
}

// Atom to hold authentication state
export const authStateAtom = atom<AuthState>({
  defaultOrgname: "",
  memberships: [],
  status: "Loading",
  user: null,
});
