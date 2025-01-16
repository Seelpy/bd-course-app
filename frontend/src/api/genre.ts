import { CreateGenre, EditGenre, DeleteGenre, ListGenreResponse } from "@shared/types/genre";
import { handleApiError } from "./utils/handleApiError";

export const genreApi = {
  PREFIX: "/api/v1/genre",

  listGenres(): Promise<ListGenreResponse> {
    return fetch(this.PREFIX, {
      method: "GET",
      credentials: "include",
    }).then((res) => {
      if (!res.ok) {
        return handleApiError(res, () => this.listGenres());
      } else {
        return res.json();
      }
    });
  },

  createGenre(body: CreateGenre): Promise<unknown> {
    return fetch(this.PREFIX, {
      method: "POST",
      credentials: "include",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify(body),
    }).then((res) => {
      if (!res.ok) {
        return handleApiError(res, () => this.createGenre(body));
      } else {
        return res.json();
      }
    });
  },

  editGenre(body: EditGenre): Promise<unknown> {
    return fetch(this.PREFIX, {
      method: "PATCH",
      credentials: "include",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify(body),
    }).then((res) => {
      if (!res.ok) {
        return handleApiError(res, () => this.editGenre(body));
      } else {
        return res.json();
      }
    });
  },

  deleteGenre(body: DeleteGenre): Promise<unknown> {
    return fetch(this.PREFIX, {
      method: "DELETE",
      credentials: "include",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify(body),
    }).then((res) => {
      if (!res.ok) {
        return handleApiError(res, () => this.deleteGenre(body));
      } else {
        return res.json();
      }
    });
  },
} as const;
