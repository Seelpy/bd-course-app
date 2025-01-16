export type StoreBookChapterTranslation = {
  bookChapterId: string;
  text: string;
};

export type GetBookChapterTranslation = {
  bookChapterId: string;
  translatorId: string;
};

export type BookChapterTranslation = {
  text: string;
};

export type ListTranslatorsByBookChapterId = {
  bookChapterId: string;
};

export type ListTranslatorsByBookChapterIdResponse = {
  translatorsId: string[];
};
