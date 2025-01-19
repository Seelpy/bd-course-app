import { Box, Button, Card, Chip, Container, Stack, Tooltip, Typography } from "@mui/material";
import { useParams, useNavigate } from "react-router-dom";
import { useEffect, useState } from "react";
import { bookApi } from "@api/book";
import { Book } from "@shared/types/book";
import { bookRatingApi } from "@api/bookRating";
import { RatingModal } from "./RatingModal";
import { useUserStore } from "@shared/stores/userStore";
import { AppRoute } from "@shared/constants/routes";
import StarIcon from "@mui/icons-material/Star";
import placeholderCover from "@assets/placeholder-cover.png";
import { enqueueSnackbar } from "notistack";

const formatRatingCount = (count: number) => {
  return count >= 1000 ? `${(count / 1000).toFixed(1)}K` : count.toString();
};

export const BookPage = () => {
  const { id } = useParams();
  const navigate = useNavigate();
  const userInfo = useUserStore((state) => state.userInfo);
  const [book, setBook] = useState<Book>();
  const [rating, setRating] = useState<number>(0);
  const [isRatingModalOpen, setIsRatingModalOpen] = useState(false);

  useEffect(() => {
    if (id) {
      bookApi
        .getBook(id)
        .then(setBook)
        .catch(() => {
          navigate(AppRoute.NotFound, { replace: true });
        });
      bookRatingApi
        .getRating(id)
        .then((data) => {
          setRating(data.average);
        })
        .catch(() => {
          enqueueSnackbar("Failed to load rating", { variant: "error" });
        });
    } else {
      navigate(AppRoute.NotFound, { replace: true });
    }
  }, [id]);

  if (!book) return null;

  const handleOpenRatingModal = () => {
    if (!userInfo) {
      navigate(AppRoute.Login);
      return;
    }
    setIsRatingModalOpen(true);
  };

  const handleRate = async (value: number) => {
    if (!id) return;
    await bookRatingApi.updateRating(id, { value });
    const newRating = await bookRatingApi.getRating(id);
    setRating(newRating.average);
  };

  const handleRemoveRating = async () => {
    if (!id) return;
    await bookRatingApi.deleteRating(id);
    const newRating = await bookRatingApi.getRating(id);
    setRating(newRating.average);
  };

  return (
    <Container maxWidth="lg" sx={{ py: 4 }}>
      <Stack spacing={2}>
        <Card elevation={3} sx={{ p: 3, borderRadius: 2 }}>
          <Stack alignItems="center" spacing={2}>
            <Box position="relative">
              <img
                src={book.cover ?? placeholderCover}
                alt={book.title}
                style={{ width: 250, height: 350, objectFit: "cover" }}
              />
              <Chip
                sx={{ position: "absolute", bottom: 8, left: 8 }}
                icon={<StarIcon sx={{ color: "gold" }} />}
                label={`${rating.toFixed(1)} (${formatRatingCount(parseInt(book.rating))})`}
                onClick={handleOpenRatingModal}
              />
            </Box>
            <Typography variant="h4" component="h1">
              {book.title}
            </Typography>
            <Stack direction="row" spacing={2}>
              <Tooltip title={!userInfo ? "Authorize first" : ""}>
                <span>
                  <Button variant="contained" disabled={!userInfo}>
                    Add to List
                  </Button>
                </span>
              </Tooltip>
              <Button variant="contained" color="secondary">
                Start Reading
              </Button>
            </Stack>
          </Stack>
        </Card>

        <Card elevation={3} sx={{ p: 3, borderRadius: 2 }}>
          <Stack spacing={2}>
            <Box sx={{ display: "flex", flexWrap: "wrap", gap: 1 }}>
              {book.authors.map((author) => (
                <Chip
                  key={author.id}
                  label={
                    <Stack>
                      <Typography variant="body2">{`${author.firstName} ${author.secondName}`}</Typography>
                      {author.nickname && <Typography variant="caption">{author.nickname}</Typography>}
                    </Stack>
                  }
                />
              ))}
            </Box>
            <Typography variant="body1">{book.description}</Typography>
            {book.genres && (
              <Box sx={{ display: "flex", flexWrap: "wrap", gap: 1 }}>
                {book.genres.map((genre) => (
                  <Chip key={genre.id} label={genre.name} />
                ))}
              </Box>
            )}
          </Stack>
        </Card>
      </Stack>

      <RatingModal
        open={isRatingModalOpen}
        onClose={() => {
          setIsRatingModalOpen(false);
        }}
        currentRating={rating}
        onRate={handleRate}
        onRemove={handleRemoveRating}
      />
    </Container>
  );
};
