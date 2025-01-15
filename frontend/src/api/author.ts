import { CreateAuthor, EditAuthor, DeleteAuthor, Author, ListAuthorResponse } from "@shared/types/author";
import { handleApiError } from "./utils/handleApiError";

export const authorApi = {
  PREFIX: "/api/v1/author",

  getAuthor(id: string): Promise<Author> {
    return fetch(`${this.PREFIX}/${id}`, {
      method: "GET",
      credentials: "include",
    }).then((res) => {
      if (!res.ok) {
        return handleApiError(res, () => this.getAuthor(id));
      }
      return res.json();
    });
  },

  listAuthors(): Promise<ListAuthorResponse> {
    return fetch(this.PREFIX, {
      method: "GET",
      credentials: "include",
    }).then((res) => {
      if (!res.ok) {
        return handleApiError(res, () => this.listAuthors());
      }
      return res.json();
    });
  },

  createAuthor(body: CreateAuthor): Promise<unknown> {
    return fetch(this.PREFIX, {
      method: "POST",
      credentials: "include",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify(body),
    }).then((res) => {
      if (!res.ok) {
        return handleApiError(res, () => this.createAuthor(body));
      }
      return res.json();
    });
  },

  editAuthor(body: EditAuthor): Promise<unknown> {
    return fetch(this.PREFIX, {
      method: "PATCH",
      credentials: "include",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify(body),
    }).then((res) => {
      if (!res.ok) {
        return handleApiError(res, () => this.editAuthor(body));
      }
      return res.json();
    });
  },

  deleteAuthor(body: DeleteAuthor): Promise<unknown> {
    return fetch(this.PREFIX, {
      method: "DELETE",
      credentials: "include",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify(body),
    }).then((res) => {
      if (!res.ok) {
        return handleApiError(res, () => this.deleteAuthor(body));
      }
      return res.json();
    });
  },
} as const;
