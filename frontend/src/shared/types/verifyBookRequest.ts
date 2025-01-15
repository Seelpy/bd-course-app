export type VerifyBookRequest = {
  verifyBookRequestId: string;
  translatorId: string;
  bookId: string;
  isVerified: boolean;
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
