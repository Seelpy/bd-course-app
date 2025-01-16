export type GetImage = {
  imageId: string;
};

export type Image = {
  imageData: string;
};

export type DeleteImage = {
  imageId: string;
};

export type StoreImageAuthor = {
  authorId: string;
  imageData: string;
};

export type StoreImageBook = {
  bookId: string;
  imageData: string;
};

export type StoreImageUser = {
  imageData: string;
};
