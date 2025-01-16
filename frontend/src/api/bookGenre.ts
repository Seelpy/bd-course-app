import { StoreBookGenre, DeleteBookGenre } from "@shared/types/bookGenre";
import { handleApiError } from "./utils/handleApiError";

export const bookGenreApi = {
  PREFIX: "/api/v1/book-genre",

  storeBookGenre(body: StoreBookGenre): Promise<unknown> {
    return fetch(this.PREFIX, {
      method: "POST",
      credentials: "include",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify(body),
    }).then((res) => {
      if (!res.ok) {
        return handleApiError(res, () => this.storeBookGenre(body));
      } else {
        return res.json();
      }
    });
  },

  deleteBookGenre(body: DeleteBookGenre): Promise<unknown> {
    return fetch(this.PREFIX, {
      method: "DELETE",
      credentials: "include",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify(body),
    }).then((res) => {
      if (!res.ok) {
        return handleApiError(res, () => this.deleteBookGenre(body));
      } else {
        return res.json();
      }
    });
  },
} as const;
