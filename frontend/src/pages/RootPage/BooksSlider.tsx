import { Book } from "@shared/types/book";
import { Box, IconButton, Paper, Typography, styled } from "@mui/material";
import { ChevronLeft, ChevronRight } from "@mui/icons-material";
import { useSwipeable } from "react-swipeable";
import { useState } from "react";
import placeholderCover from "@assets/placeholder-cover.png";

type BookSlider = Book & { chapterChip?: string };

type BooksSliderProps = {
  sliderName: string;
  books: BookSlider[];
};

const SLIDER_HEIGHT = 220;
const BOOK_WIDTH = 135;
const BOOK_HEIGHT = 190;
const BUTTON_SIZE = SLIDER_HEIGHT * 0.2;

const SliderContainer = styled(Paper)(({ theme }) => ({
  height: SLIDER_HEIGHT,
  position: "relative",
  borderRadius: theme.spacing(2),
  padding: theme.spacing(2),
  marginBottom: theme.spacing(2),
  overflow: "hidden",
}));

const BooksContainer = styled(Box)({
  display: "flex",
  gap: 16,
  height: "100%",
  transition: "transform 0.3s ease-out",
  paddingTop: 24,
});

const NavigationButton = styled(IconButton)(({ theme }) => ({
  position: "absolute",
  top: "50%",
  transform: "translateY(-50%)",
  width: BUTTON_SIZE,
  height: BUTTON_SIZE,
  backgroundColor: theme.palette.background.paper,
  "&:hover": {
    backgroundColor: theme.palette.action.hover,
  },
  zIndex: 2,
}));

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

export const BooksSlider = ({ sliderName, books }: BooksSliderProps) => {
  const [offset, setOffset] = useState(0);
  const maxOffset = -Math.max(0, books.length * (BOOK_WIDTH + 16) - window.innerWidth + 64);

  const handlers = useSwipeable({
    onSwipedLeft: () => {
      setOffset((curr) => Math.max(maxOffset, curr - BOOK_WIDTH - 16));
    },
    onSwipedRight: () => {
      setOffset((curr) => Math.min(0, curr + BOOK_WIDTH + 16));
    },
  });

  return (
    <SliderContainer elevation={2}>
      <Typography variant="h6" sx={{ position: "absolute", top: 16, left: 16 }}>
        {sliderName}
      </Typography>

      <NavigationButton
        sx={{ left: 8 }}
        onClick={() => {
          setOffset((curr) => Math.min(0, curr + BOOK_WIDTH + 16));
        }}
        disabled={offset === 0}
      >
        <ChevronLeft />
      </NavigationButton>

      <NavigationButton
        sx={{ right: 8 }}
        onClick={() => {
          setOffset((curr) => Math.max(maxOffset, curr - BOOK_WIDTH - 16));
        }}
        disabled={offset <= maxOffset}
      >
        <ChevronRight />
      </NavigationButton>

      <Box overflow="hidden">
        <BooksContainer {...handlers} sx={{ transform: `translateX(${offset.toString()}px)` }}>
          {books.map((book) => (
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
          ))}
        </BooksContainer>
      </Box>
    </SliderContainer>
  );
};
