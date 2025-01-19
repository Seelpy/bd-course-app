import { Box, styled, Typography } from "@mui/material";
import { Book } from "@shared/types/book";
import placeholderCover from "@assets/placeholder-cover.png";
import { useNavigate } from "react-router-dom";

const BookCard = styled(Box)<{ width: number }>(({ width }) => ({
  display: "flex",
  flexDirection: "column",
  alignItems: "center",
  width: width,
  gap: 8,
  cursor: "pointer",
}));

const BookCover = styled("img")<{ width: number }>(({ width }) => ({
  width: width,
  height: width * 1.4,
  objectFit: "cover",
  borderRadius: 8,
}));

type BookPreviewProps = {
  book: Book;
  width?: number;
};

export const BookPreview = ({ book, width = 150 }: BookPreviewProps) => {
  const navigate = useNavigate();

  return (
    <BookCard
      width={width}
      onClick={() => {
        navigate(`/book/${book.bookId}`);
      }}
    >
      <BookCover width={width} src={book.cover ?? placeholderCover} alt={book.title} />
      <Typography
        variant="body2"
        sx={{
          width: "100%",
          paddingX: 0.5,
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
