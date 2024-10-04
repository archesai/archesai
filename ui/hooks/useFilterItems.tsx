import { limitAtom, pageAtom, queryAtom, rangeAtom } from "@/state/filterAtoms";
import { useAtom } from "jotai";

export const useFilterItems = () => {
  // Atoms state management
  const [page, setPage] = useAtom(pageAtom);
  const [limit, setLimit] = useAtom(limitAtom);
  const [query, setQuery] = useAtom(queryAtom);
  const [range, setRange] = useAtom(rangeAtom);

  return {
    limit,
    page,
    query,
    range,
    setLimit,
    setPage,
    setQuery,
    setRange,
  };
};
