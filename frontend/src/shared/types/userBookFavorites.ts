export type UserBookFavoritesType = "READING" | "PLANNED" | "DEFERRED" | "READ" | "DROPPED" | "FAVORITE";

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
