import "@tanstack/react-table";

declare module "@tanstack/react-table" {
  interface TableMeta<_TData> {
    entityKey?: string;
  }
}
