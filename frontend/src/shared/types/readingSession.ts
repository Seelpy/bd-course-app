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
