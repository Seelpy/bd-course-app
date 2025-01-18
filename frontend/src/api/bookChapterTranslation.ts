import {
  GetBookChapterTranslation,
  BookChapterTranslation,
  ListTranslatorsByBookChapterId,
  ListTranslatorsByBookChapterIdResponse,
  StoreBookChapterTranslation,
} from "@shared/types/bookChapterTranslation";
import { handleApiError } from "./utils/handleApiError";

export const bookChapterTranslationApi = {
  PREFIX: "/api/v1/book-chapter-translation",

  getTranslation(body: GetBookChapterTranslation): Promise<BookChapterTranslation> {
    return fetch(this.PREFIX, {
      method: "POST",
      credentials: "include",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify(body),
    }).then((res) => {
      if (!res.ok) {
        return handleApiError(res, () => this.getTranslation(body));
      }
      return res.json();
    });
  },

  storeTranslation(body: StoreBookChapterTranslation): Promise<unknown> {
    return fetch(this.PREFIX, {
      method: "POST",
      credentials: "include",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify(body),
    }).then((res) => {
      if (!res.ok) {
        return handleApiError(res, () => this.storeTranslation(body));
      }
      return res.json();
    });
  },

  listTranslators(body: ListTranslatorsByBookChapterId): Promise<ListTranslatorsByBookChapterIdResponse> {
    return fetch(`${this.PREFIX}/translator`, {
      method: "POST",
      credentials: "include",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify(body),
    }).then((res) => {
      if (!res.ok) {
        return handleApiError(res, () => this.listTranslators(body));
      }
      return res.json();
    });
  },
} as const;
