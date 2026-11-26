import { keepPreviousData, useQuery } from "@tanstack/react-query";

import { fetchImages } from "../api/fetchImages";

interface Params {
  offset: number;
  limit: number;
}

export const useImages = ({ offset, limit }: Params) => {
  const query = useQuery({
    queryKey: ["images", offset, limit],
    queryFn: ({ signal }) => {
      return fetchImages(offset, limit, signal);
    },
    placeholderData: keepPreviousData,
    staleTime: 0,
    refetchOnMount: "always"
  });
  return query;
};
