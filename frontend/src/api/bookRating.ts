import { UpdateBookRating, GetBookRatingResponse } from "@shared/types/bookRating";
import { handleApiError } from "./utils/handleApiError";

export const bookRatingApi = {
  PREFIX: "/api/v1/book",

  updateRating(bookId: string, body: UpdateBookRating): Promise<unknown> {
    return fetch(`${this.PREFIX}/${bookId}/raiting`, {
      method: "POST",
      credentials: "include",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify(body),
    }).then((res) => {
      if (!res.ok) {
        return handleApiError(res, () => this.updateRating(bookId, body));
      }
      return res.json();
    });
  },

  deleteRating(bookId: string): Promise<unknown> {
    return fetch(`${this.PREFIX}/${bookId}/raiting`, {
      method: "DELETE",
      credentials: "include",
    }).then((res) => {
      if (!res.ok) {
        return handleApiError(res, () => this.deleteRating(bookId));
      }
      return res.json();
    });
  },

  getRating(bookId: string): Promise<GetBookRatingResponse> {
    return fetch(`${this.PREFIX}/${bookId}/raiting`, {
      method: "GET",
      credentials: "include",
    }).then((res) => {
      if (!res.ok) {
        return handleApiError(res, () => this.getRating(bookId));
      }
      return res.json();
    });
  },
} as const;
