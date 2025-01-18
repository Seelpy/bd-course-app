import { Book } from "./book";

export type VerifyBookRequest = {
  verifyBookRequestId: string;
  translatorId: string;
  book: Book;
  isVerified?: boolean;
  sendDateMilli: number;
};

export type ListVerifyBookRequestResponse = {
  verifyBookRequests: VerifyBookRequest[];
};

export type CreateVerifyBook = {
  bookId: string;
};

export type DeleteVerifyBook = {
  verifyBookRequestId: string;
};

export type AcceptBookChapter = {
  verifyBookRequestId: string;
  accept: boolean;
};
