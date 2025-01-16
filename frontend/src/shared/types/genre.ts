export type GetReadingSession = {
  bookId: string;
};

export type StoreReadingSession = {
  bookId: string;
  bookChaptedId: string;
};

export type GetReadingSessionResponse = {
  bookChaptedId: string;
};

export type CreateGenre = {
  name: string;
};

export type EditGenre = {
  id: string;
  name: string;
};

export type DeleteGenre = {
  id: string;
};

export type Genre = {
  id: string;
  name: string;
};

export type ListGenreResponse = {
  genres: Genre[];
};
