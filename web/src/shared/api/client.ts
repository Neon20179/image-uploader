import axios from "axios";

const apiOrigin = (
  process.env.NEXT_PUBLIC_API_URL || "http://localhost:8000"
).replace(/\/+$/, "");

export const api = axios.create({
  baseURL: `${apiOrigin}/api/v1`,
  timeout: 20_000,
  headers: { Accept: "application/json" }
});

export const getErrorMessage = (error: unknown): string => {
  if (axios.isAxiosError(error)) {
    if (!error.response) {
      return "Could not connect to the server. Check your connection and try again.";
    }
    if (error.response.status >= 500) {
      return "The server is temporarily unavailable. Please try again later.";
    }
    if (error.response.status === 429) {
      return "Too many requests. Please try again in a moment.";
    }
    if (error.response.status === 403) {
      return "Storage rejected the upload. Check the file type and size, or select the file again.";
    }
    return "The request failed. Please try again.";
  }
  return "Something went wrong. Please try again.";
};
