import { atom } from "jotai";

interface DateRange {
  from: Date;
  to: Date;
}

// Define atoms
export const pageAtom = atom<number>(0);
export const limitAtom = atom<number>(10);
export const queryAtom = atom<string>("");
export const rangeAtom = atom<DateRange>();
export const sortByAtom = atom<string>("createdAt");
export const sortDirectionAtom = atom<"asc" | "desc">("desc");
