import { Box, styled, Typography } from "@mui/material";
import { Book } from "@shared/types/book";
import placeholderCover from "@assets/placeholder-cover.png";

const BookCard = styled(Box)<{ $width: number }>(({ $width }) => ({
  display: "flex",
  flexDirection: "column",
  alignItems: "center",
  width: $width,
  gap: 8,
}));

const BookCover = styled("img")<{ $width: number }>(({ $width }) => ({
  width: $width,
  height: $width * 1.4,
  objectFit: "cover",
  borderRadius: 8,
}));

type BookPreviewProps = {
  book: Book;
  width?: number;
};

export const BookPreview = ({ book, width = 150 }: BookPreviewProps) => {
  return (
    <BookCard $width={width}>
      <BookCover $width={width} src={book.cover ?? placeholderCover} alt={book.title} />
      <Typography
        variant="body2"
        sx={{
          width: "100%",
          overflow: "hidden",
          textOverflow: "ellipsis",
          display: "-webkit-box",
          WebkitLineClamp: 2,
          WebkitBoxOrient: "vertical",
          height: "2.4em",
          lineHeight: "1.2em",
        }}
      >
        {book.title}
      </Typography>
    </BookCard>
  );
};
