import { Box, styled, Typography } from "@mui/material";
import { Book } from "@shared/types/book";
import placeholderCover from "@assets/placeholder-cover.png";

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

export const BookPreview = ({ book }: BookPreviewProps) => {
  return (
    <BookCard key={book.bookId}>
      <BookCover src={book.cover ?? placeholderCover} alt={book.title} />
      <Typography
        variant="body2"
        sx={{
          width: "100%",
          textAlign: "center",
          overflow: "hidden",
          textOverflow: "ellipsis",
          whiteSpace: "nowrap",
        }}
      >
        {book.title}
      </Typography>
    </BookCard>
  );
};
