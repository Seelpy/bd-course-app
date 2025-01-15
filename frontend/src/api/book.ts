import { CreateBook, EditBook, DeleteBook, Book, ListBookResponse } from "@shared/types/book";
import { handleApiError } from "./utils/handleApiError";

export const bookApi = {
  PREFIX: "/api/v1/book",

  getBook(id: string): Promise<Book> {
    return fetch(`${this.PREFIX}/${id}`, {
      method: "GET",
      credentials: "include",
    }).then((res) => {
      if (!res.ok) {
        return handleApiError(res, () => this.getBook(id));
      }
      return res.json();
    });
  },

  listBooks(page: number, size: number): Promise<ListBookResponse> {
    return fetch(`${this.PREFIX}/${page.toString()}/${size.toString()}`, {
      method: "GET",
      credentials: "include",
    }).then((res) => {
      if (!res.ok) {
        return handleApiError(res, () => this.listBooks(page, size));
      }
      return res.json();
    });
  },

  createBook(body: CreateBook): Promise<unknown> {
    return fetch(this.PREFIX, {
      method: "POST",
      credentials: "include",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify(body),
    }).then((res) => {
      if (!res.ok) {
        return handleApiError(res, () => this.createBook(body));
      }
      return res.json();
    });
  },

  editBook(body: EditBook): Promise<unknown> {
    return fetch(this.PREFIX, {
      method: "PATCH",
      credentials: "include",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify(body),
    }).then((res) => {
      if (!res.ok) {
        return handleApiError(res, () => this.editBook(body));
      }
      return res.json();
    });
  },

  deleteBook(body: DeleteBook): Promise<unknown> {
    return fetch(this.PREFIX, {
      method: "DELETE",
      credentials: "include",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify(body),
    }).then((res) => {
      if (!res.ok) {
        return handleApiError(res, () => this.deleteBook(body));
      }
      return res.json();
    });
  },
} as const;
