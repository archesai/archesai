import { MemberEntity, UserEntity } from "@/generated/archesApiSchemas";
import { atom } from "jotai";
import { atomWithStorage } from "jotai/utils";

// Define the initial authentication state
interface AuthState {
  defaultOrgname: string;
  isLoading: boolean;
  memberships: MemberEntity[];
  user: null | UserEntity;
}

// Atom to hold authentication state
export const authStateAtom = atom<AuthState>({
  defaultOrgname: "",
  isLoading: true,
  memberships: [],
  user: null,
});

export const accessTokenAtom = atomWithStorage<string>("accessToken", ""); // "accessToken" is the localStorage key
export const refreshTokenAtom = atomWithStorage<string>("refreshToken", ""); // "refreshToken" is the localStorage key
export const readStoragAtom = atomWithStorage<boolean>("readStorage", false); // "refreshToken" is the localStorage key
