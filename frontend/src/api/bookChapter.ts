import {
  CreateBookChapter,
  EditBookChapter,
  DeleteBookChapter,
  ListBookChapter,
  ListBookChapterResponse,
} from "@shared/types/bookChapter";
import { handleApiError } from "./utils/handleApiError";

export const bookChapterApi = {
  PREFIX: "/api/v1/book-chapter",

  listBookChapters(body: ListBookChapter): Promise<ListBookChapterResponse> {
    return fetch(this.PREFIX, {
      method: "GET",
      credentials: "include",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify(body),
    }).then((res) => {
      if (!res.ok) {
        return handleApiError(res, () => this.listBookChapters(body));
      }
      return res.json();
    });
  },

  createBookChapter(body: CreateBookChapter): Promise<unknown> {
    return fetch(this.PREFIX, {
      method: "POST",
      credentials: "include",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify(body),
    }).then((res) => {
      if (!res.ok) {
        return handleApiError(res, () => this.createBookChapter(body));
      }
      return res.json();
    });
  },

  editBookChapter(body: EditBookChapter): Promise<unknown> {
    return fetch(this.PREFIX, {
      method: "PATCH",
      credentials: "include",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify(body),
    }).then((res) => {
      if (!res.ok) {
        return handleApiError(res, () => this.editBookChapter(body));
      }
      return res.json();
    });
  },

  deleteBookChapter(body: DeleteBookChapter): Promise<unknown> {
    return fetch(this.PREFIX, {
      method: "DELETE",
      credentials: "include",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify(body),
    }).then((res) => {
      if (!res.ok) {
        return handleApiError(res, () => this.deleteBookChapter(body));
      }
      return res.json();
    });
  },
} as const;
