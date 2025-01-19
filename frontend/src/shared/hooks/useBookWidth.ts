const MIN_BOOKS_PER_ROW = 1;
const MAX_BOOKS_PER_ROW = 5;
const MIN_BOOK_WIDTH = 150;
const GRID_SPACING = 16;

export const useBookWidth = (containerWidth: number) => {
  let booksPerRow = MAX_BOOKS_PER_ROW;
  let bookWidth = 0;

  while (booksPerRow >= MIN_BOOKS_PER_ROW) {
    const availableWidth = containerWidth - GRID_SPACING * (booksPerRow - 1);
    const calculatedWidth = Math.floor(availableWidth / booksPerRow);

    if (calculatedWidth >= MIN_BOOK_WIDTH) {
      bookWidth = calculatedWidth;
      break;
    }

    booksPerRow--;
  }

  if (bookWidth < MIN_BOOK_WIDTH) {
    bookWidth = MIN_BOOK_WIDTH;
  }

  return { bookWidth, booksPerRow };
};
