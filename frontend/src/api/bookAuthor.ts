import { DeleteBookAuthor, StoreBookAuthor } from "@shared/types/bookAuthor";
import { handleApiError } from "./utils/handleApiError";

export const bookAuthorApi = {
  PREFIX: "/api/v1/book-author",

  storeBookAuthor(body: StoreBookAuthor): Promise<unknown> {
    return fetch(this.PREFIX, {
      method: "POST",
      credentials: "include",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify(body),
    }).then((res) => {
      if (!res.ok) {
        return handleApiError(res, () => this.storeBookAuthor(body));
      } else {
        return res.json();
      }
    });
  },

  deleteBookAuthor(body: DeleteBookAuthor): Promise<unknown> {
    return fetch(this.PREFIX, {
      method: "DELETE",
      credentials: "include",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify(body),
    }).then((res) => {
      if (!res.ok) {
        return handleApiError(res, () => this.deleteBookAuthor(body));
      } else {
        return res.json();
      }
    });
  },
} as const;
