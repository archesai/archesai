import {
  limitAtom,
  pageAtom,
  queryAtom,
  rangeAtom,
  sortByAtom,
  sortDirectionAtom,
} from "@/state/filterAtoms";
import { useAtom } from "jotai";

export const useFilterItems = () => {
  // Atoms state management
  const [page, setPage] = useAtom(pageAtom);
  const [limit, setLimit] = useAtom(limitAtom);
  const [query, setQuery] = useAtom(queryAtom);
  const [range, setRange] = useAtom(rangeAtom);
  const [sortBy, setSortBy] = useAtom(sortByAtom);
  const [sortDirection, setSortDirection] = useAtom(sortDirectionAtom);

  return {
    limit,
    page,
    query,
    range,
    setLimit,
    setPage,
    setQuery,
    setRange,
    setSortBy,
    setSortDirection,
    sortBy,
    sortDirection,
  };
};
