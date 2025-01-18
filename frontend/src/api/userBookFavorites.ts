import {
  GetFavoriteTypeByBook,
  ListBookByUserBookFavorites,
  StoreUserBookFavorites,
  DeleteUserBookFavorites,
  UserBookFavoritesType,
  ListBookByUserBookFavouritesResponse,
} from "@shared/types/userBookFavorites";
import { handleApiError } from "./utils/handleApiError";

export const userBookFavoritesApi = {
  PREFIX: "/api/v1/user-book-favourites",

  getFavoriteTypeByBook(body: GetFavoriteTypeByBook): Promise<UserBookFavoritesType> {
    return fetch(this.PREFIX, {
      method: "POST",
      credentials: "include",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify(body),
    }).then((res) => {
      if (!res.ok) {
        return handleApiError(res, () => this.getFavoriteTypeByBook(body));
      }
      return res.json();
    });
  },

  listBooksByFavorites(body: ListBookByUserBookFavorites): Promise<ListBookByUserBookFavouritesResponse> {
    return fetch(`${this.PREFIX}/book`, {
      method: "POST",
      credentials: "include",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify(body),
    }).then((res) => {
      if (!res.ok) {
        return handleApiError(res, () => this.listBooksByFavorites(body));
      }
      return res.json();
    });
  },

  storeFavorite(body: StoreUserBookFavorites): Promise<unknown> {
    return fetch(this.PREFIX, {
      method: "POST",
      credentials: "include",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify(body),
    }).then((res) => {
      if (!res.ok) {
        return handleApiError(res, () => this.storeFavorite(body));
      }
      return res.json();
    });
  },

  deleteFavorite(body: DeleteUserBookFavorites): Promise<unknown> {
    return fetch(this.PREFIX, {
      method: "DELETE",
      credentials: "include",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify(body),
    }).then((res) => {
      if (!res.ok) {
        return handleApiError(res, () => this.deleteFavorite(body));
      }
      return res.json();
    });
  },
} as const;
