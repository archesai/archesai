import { atom } from "jotai";

export const viewAtom = atom<"grid" | "table">("table");
