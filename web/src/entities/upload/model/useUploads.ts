import { keepPreviousData, useQuery } from "@tanstack/react-query";

import { fetchUploads } from "../api/fetchUploads";

interface Params {
  offset: number;
  limit: number;
}

export const useUploads = ({ offset, limit }: Params) => {
  return useQuery({
    queryKey: ["uploads", offset, limit],
    queryFn: ({ signal }) => {
      return fetchUploads(offset, limit, signal);
    },
    placeholderData: keepPreviousData,
    staleTime: 0,
    refetchInterval: 5_000
  });
};
