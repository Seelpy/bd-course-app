import { Box, Button, Card, Chip, Container, Stack, Tooltip, Typography, Menu, MenuItem } from "@mui/material";
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
import { Add } from "@mui/icons-material";
import { UserBookFavoritesType } from "@shared/types/userBookFavorites";
import { userBookFavoritesApi } from "@api/userBookFavorites";

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
  const [favoriteType, setFavoriteType] = useState<UserBookFavoritesType | null>(null);
  const [menuAnchor, setMenuAnchor] = useState<null | HTMLElement>(null);
  const [menuWidth, setMenuWidth] = useState<number>(0);

  useEffect(() => {
    if (id) {
      bookApi
        .getBook(id)
        .then((data) => {
          setBook(data.book);
        })
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

  useEffect(() => {
    if (id && userInfo) {
      userBookFavoritesApi
        .getFavoriteTypeByBook({ bookId: id })
        .then(setFavoriteType)
        .catch(() => {
          setFavoriteType(null);
        });
    }
  }, [id, userInfo]);

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

  const handleListButtonClick = (event: React.MouseEvent<HTMLElement>) => {
    setMenuAnchor(event.currentTarget);
    setMenuWidth(event.currentTarget.offsetWidth);
  };

  const handleListClose = () => {
    setMenuAnchor(null);
  };

  const handleFavoriteSelect = (type: UserBookFavoritesType | "REMOVE") => {
    if (!id) return;

    if (type === "REMOVE") {
      userBookFavoritesApi
        .deleteFavorite({ bookId: id })
        .then(() => {
          setFavoriteType(null);
          enqueueSnackbar("Book removed from list", { variant: "success" });
        })
        .catch((error: Error) => {
          enqueueSnackbar(error.message, { variant: "error" });
        });
    } else {
      userBookFavoritesApi
        .storeFavorite({ bookId: id, type })
        .then(() => {
          setFavoriteType(type);
          enqueueSnackbar("Book added to list", { variant: "success" });
        })
        .catch((error: Error) => {
          enqueueSnackbar(error.message, { variant: "error" });
        });
    }

    handleListClose();
  };

  const favoriteTypes: UserBookFavoritesType[] = ["READING", "PLANNED", "DEFERRED", "READ", "DROPPED", "FAVORITE"];

  return (
    <Container maxWidth="lg" sx={{ py: 4 }}>
      <Stack spacing={2}>
        <Card elevation={0} sx={{ p: 3, borderRadius: 2 }}>
          <Stack alignItems="center" spacing={2}>
            <Box position="relative">
              <img
                src={book.cover ?? placeholderCover}
                alt={book.title}
                style={{ width: 250, height: 350, objectFit: "cover", borderRadius: 8 }}
              />
              <Chip
                sx={{ position: "absolute", bottom: 8, left: `50%`, transform: `translate(-50%, -50%)` }}
                icon={<StarIcon color="warning" />}
                label={`${rating.toFixed(1)} (${formatRatingCount(parseInt(book.rating))})`}
                onClick={handleOpenRatingModal}
              />
            </Box>
            <Typography
              variant="h4"
              component="h1"
              sx={{
                textAlign: "center",
                width: "100%",
              }}
            >
              {book.title}
            </Typography>
            <Stack direction="row" spacing={2}>
              <Tooltip title={!userInfo ? "Authorize first" : ""}>
                <span>
                  <Button
                    variant="contained"
                    color="inherit"
                    startIcon={<Add />}
                    onClick={handleListButtonClick}
                    disabled={!userInfo}
                  >
                    {favoriteType ? `In ${favoriteType.toLowerCase()}` : "Add to List"}
                  </Button>
                </span>
              </Tooltip>
              <Menu
                anchorEl={menuAnchor}
                open={Boolean(menuAnchor)}
                onClose={handleListClose}
                anchorOrigin={{
                  vertical: "bottom",
                  horizontal: "center",
                }}
                transformOrigin={{
                  vertical: "top",
                  horizontal: "center",
                }}
                slotProps={{
                  paper: {
                    style: {
                      width: menuWidth,
                      maxWidth: "none",
                    },
                  },
                }}
              >
                {favoriteTypes.map((type) => (
                  <MenuItem
                    key={type}
                    onClick={() => {
                      handleFavoriteSelect(type);
                    }}
                    selected={type === favoriteType}
                    sx={{ justifyContent: "center" }}
                  >
                    {type.charAt(0) + type.slice(1).toLowerCase()}
                  </MenuItem>
                ))}
                {favoriteType && (
                  <MenuItem
                    onClick={() => {
                      handleFavoriteSelect("REMOVE");
                    }}
                    sx={{ color: "error.main", justifyContent: "center" }}
                  >
                    Remove
                  </MenuItem>
                )}
              </Menu>
              <Button variant="contained">Start Reading</Button>
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
