import { Box, styled } from "@mui/material";
import { Book } from "@shared/types/book";

const BOOK_WIDTH = 135;
const BOOK_HEIGHT = 190;

const BookCard = styled(Box)({
  display: "flex",
  flexDirection: "column",
  alignItems: "center",
  width: BOOK_WIDTH,
  gap: 8,
});

const BookCover = styled("img")({
  width: BOOK_WIDTH,
  height: BOOK_HEIGHT,
  objectFit: "cover",
  borderRadius: 8,
});

type BookPreviewProps = {
  book: Book;
};

export const MenuDesktop = ({ book }: BookPreviewProps) => {};
