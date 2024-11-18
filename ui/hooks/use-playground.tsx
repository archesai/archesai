import { ContentEntity, ToolEntity } from "@/generated/archesApiSchemas";
import {
  parseAsArrayOf,
  parseAsJson,
  parseAsString,
  useQueryState,
} from "nuqs";
import { useCallback } from "react";

export const usePlayground = () => {
  const [selectedTool, setSelectedTool] = useQueryState(
    "selectedTool",
    parseAsJson<ToolEntity>((tool) => tool as ToolEntity)
  );
  const [selectedContent, setSelectedContent] = useQueryState(
    "selectedContent",
    parseAsArrayOf(parseAsJson<ContentEntity>((tool) => tool as ContentEntity))
  );
  const [selectedRunId, setSelectedRunId] = useQueryState(
    "selectedRunId",
    parseAsString
  );

  const clearParams = useCallback(() => {
    setSelectedContent(null);
    setSelectedRunId(null);
    setSelectedTool(null);
  }, []);

  return {
    clearParams,
    selectedContent,
    selectedRunId,
    selectedTool,
    setSelectedContent,
    setSelectedRunId,
    setSelectedTool,
  };
};
