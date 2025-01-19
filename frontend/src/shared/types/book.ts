import { Author } from "./author";
import { Genre } from "./genre";

export type SortBy = "TITLE" | "RATING" | "RATING_COUNT" | "CHAPTERS_COUNT";
export type SortType = "ASC" | "DESC";

export type CreateBook = {
  title: string;
  description: string;
};

export type EditBook = {
  id: string;
  title: string;
  description: string;
};

export type DeleteBook = {
  id: string;
};

export type Book = {
  bookId: string;
  cover?: string;
  title: string;
  description: string;
  authors: Author[];
  isLoggedUserTranslator: string;
  rating: string;
  genres: Genre[];
};

export type GetBookResponse = {
  book: Book;
};

export type ListBookResponse = {
  books: Book[];
  countPages: number;
};
