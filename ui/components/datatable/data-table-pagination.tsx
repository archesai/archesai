import { Button } from "@/components/ui/button";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { useFilterItems } from "@/hooks/useFilterItems";
import { useSelectItems } from "@/hooks/useSelectItems";
import {
  ChevronLeftIcon,
  ChevronRightIcon,
  DoubleArrowLeftIcon,
  DoubleArrowRightIcon,
} from "@radix-ui/react-icons";

interface DataTablePaginationProps<TData> {
  data: {
    metadata: {
      limit: number;
      offset: number;
      totalResults: number;
    };
    results: TData[];
  };
}

export function DataTablePagination<TData>({
  data,
}: DataTablePaginationProps<TData>) {
  const { limit, page, setLimit, setPage } = useFilterItems();
  const { selectedItems } = useSelectItems({ items: data?.results || [] });
  return (
    <div className="flex items-center justify-between px-2 backdrop-blur-md mb-3">
      <div className="flex-1 text-sm text-muted-foreground sm:block hidden">
        {data?.metadata?.totalResults} found - {selectedItems.length} of{" "}
        {Math.min(limit, data?.results.length)} item(s) selected.
      </div>
      <div className="flex items-center space-x-6 lg:space-x-8">
        <div className="flex items-center space-x-0 sm:space-x-2">
          <p className="text-sm font-medium sm:block hidden">Rows per page</p>
          <Select
            onValueChange={(value) => {
              setLimit(Number(value));
            }}
            value={`${limit}`}
          >
            <SelectTrigger className="h-8 w-[70px]">
              <SelectValue placeholder={limit} />
            </SelectTrigger>
            <SelectContent side="top">
              {[10, 20, 30, 40, 50].map((limit) => (
                <SelectItem key={limit} value={`${limit}`}>
                  {limit}
                </SelectItem>
              ))}
            </SelectContent>
          </Select>
        </div>
        <div className="flex w-[100px] items-center justify-center text-sm font-medium">
          Page {page + 1} of{" "}
          {isNaN(Math.max(Math.ceil(data?.metadata?.totalResults / limit), 1))
            ? 1
            : Math.max(Math.ceil(data?.metadata?.totalResults / limit), 1)}
        </div>
        <div className="flex items-center space-x-2">
          <Button
            className="hidden h-8 w-8 p-0 lg:flex"
            disabled={page === 0}
            onClick={() => setPage(0)}
            variant="outline"
          >
            <span className="sr-only">Go to first page</span>
            <DoubleArrowLeftIcon className="h-4 w-4" />
          </Button>
          <Button
            className="h-8 w-8 p-0"
            disabled={page === 0}
            onClick={() => setPage(page - 1)}
            variant="outline"
          >
            <span className="sr-only">Go to previous page</span>
            <ChevronLeftIcon className="h-4 w-4" />
          </Button>
          <Button
            className="h-8 w-8 p-0"
            disabled={
              page >= Math.ceil(data?.metadata?.totalResults / limit) - 1
            }
            onClick={() => setPage(page + 1)}
            variant="outline"
          >
            <span className="sr-only">Go to next page</span>
            <ChevronRightIcon className="h-4 w-4" />
          </Button>
          <Button
            className="hidden h-8 w-8 p-0 lg:flex"
            disabled={
              page >= Math.ceil(data?.metadata?.totalResults / limit) - 1
            }
            onClick={() =>
              setPage(Math.ceil(data?.metadata?.totalResults / limit) - 1)
            }
            variant="outline"
          >
            <span className="sr-only">Go to last page</span>
            <DoubleArrowRightIcon className="h-4 w-4" />
          </Button>
        </div>
      </div>
    </div>
  );
}
