import { Book } from "@shared/types/book";
import { Box, IconButton, Paper, Typography, styled } from "@mui/material";
import { ChevronLeft, ChevronRight } from "@mui/icons-material";
import { useSwipeable } from "react-swipeable";
import { useState } from "react";
import { BookPreview } from "@shared/components/BookPreview/BookPreview";

type BookSlider = Book & { chapterChip?: string };

type BooksSliderProps = {
  sliderName: string;
  books: BookSlider[];
};

const SLIDER_HEIGHT = 200;
const BOOK_WIDTH = 90;
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
    <>
      <Typography variant="h6">{sliderName}</Typography>
      <SliderContainer elevation={2}>
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
              <BookPreview key={book.bookId} book={book} width={BOOK_WIDTH} />
            ))}
          </BooksContainer>
        </Box>
      </SliderContainer>
    </>
  );
};
