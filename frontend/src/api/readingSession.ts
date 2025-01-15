import { GetReadingSession, GetReadingSessionResponse, StoreReadingSession } from "@shared/types/readingSession";
import { handleApiError } from "./utils/handleApiError";

export const readingSessionApi = {
  PREFIX: "/api/v1/reading-session",

  getLastReadingSession(body: GetReadingSession): Promise<GetReadingSessionResponse> {
    return fetch(this.PREFIX, {
      method: "GET",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify(body),
    }).then((res) => {
      if (!res.ok) {
        return handleApiError(res, () => this.getLastReadingSession(body));
      } else {
        return res.json();
      }
    });
  },
  storeReadingSession(body: StoreReadingSession): Promise<unknown> {
    return fetch(this.PREFIX, {
      method: "POST",
      credentials: "include",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify(body),
    }).then((res) => {
      if (!res.ok) {
        return handleApiError(res, () => this.storeReadingSession(body));
      } else {
        return res.json();
      }
    });
  },
} as const;
