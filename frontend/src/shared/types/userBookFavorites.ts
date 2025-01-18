import { Book } from "@shared/types/book";

export type UserBookFavoritesType = "READING" | "PLANNED" | "DEFERRED" | "READ" | "DROPPED" | "FAVORITE";

export type UserBookFavouritesBooks = {
  type: UserBookFavoritesType;
  books: Book[];
};

export type ListBookByUserBookFavouritesResponse = {
  userBookFavouritesBooks: UserBookFavouritesBooks[];
};

export type GetFavoriteTypeByBook = {
  bookId: string;
};

export type ListBookByUserBookFavorites = {
  types: UserBookFavoritesType[];
};

export type StoreUserBookFavorites = {
  bookId: string;
  type: UserBookFavoritesType;
};

export type DeleteUserBookFavorites = {
  bookId: string;
};
