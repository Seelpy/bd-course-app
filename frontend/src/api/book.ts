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

  searchBooks(
    page: number,
    size: number,
    params?: {
      bookTitle?: string;
      authorIds?: string[];
      rating?: "MIN_RATING" | "MAX_RATING";
      genreIds?: string[];
      numberBookChapter?: "MIN_BOOK_CHAPTERS" | "MAX_BOOK_CHAPTERS";
    },
  ): Promise<ListBookResponse> {
    const searchParams = new URLSearchParams({
      page: page.toString(),
      size: size.toString(),
    });

    if (params?.bookTitle) {
      searchParams.append("bookTitle", params.bookTitle);
    }
    if (params?.authorIds) {
      params.authorIds.forEach((id) => {
        searchParams.append("authorIds[]", id);
      });
    }
    if (params?.rating) {
      searchParams.append("rating", params.rating);
    }
    if (params?.genreIds) {
      params.genreIds.forEach((id) => {
        searchParams.append("genreIds[]", id);
      });
    }
    if (params?.numberBookChapter) {
      searchParams.append("numberBookChapter", params.numberBookChapter);
    }

    return fetch(`${this.PREFIX}/search?${searchParams}`, {
      method: "GET",
      credentials: "include",
    }).then((res) => {
      if (!res.ok) {
        return handleApiError(res, () => this.searchBooks(page, size, params));
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
