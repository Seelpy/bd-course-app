export type BookChapter = {
  bookChapterId: string;
  index: number;
  title: string;
};

export type CreateBookChapter = {
  bookId: string;
  title: string;
};

export type EditBookChapter = {
  bookChapterId: string;
  title: string;
};

export type DeleteBookChapter = {
  bookChapterId: string;
};

export type ListBookChapter = {
  bookId: string;
};

export type ListBookChapterResponse = {
  bookChapters: BookChapter[];
};
