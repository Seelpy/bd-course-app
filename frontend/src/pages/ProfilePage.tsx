import { useEffect, useState } from "react";
import { useParams, Navigate, useNavigate } from "react-router-dom";
import { Container, Paper, Avatar, Typography, Box, Tabs, Tab, Grid2, useMediaQuery, useTheme } from "@mui/material";
import { userApi } from "@api/user";
import { imageApi } from "@api/image";
import { userBookFavoritesApi } from "@api/userBookFavorites";
import { User } from "@shared/types/user";
import { Book } from "@shared/types/book";
import { UserBookFavoritesType } from "@shared/types/userBookFavorites";
import { useSnackbar } from "notistack";
import { AppRoute } from "@shared/constants/routes";

type FavoriteType = UserBookFavoritesType | "ALL";

export function ProfilePage() {
  const { id } = useParams<{ id: string }>();
  const { enqueueSnackbar } = useSnackbar();
  const navigate = useNavigate();
  const theme = useTheme();
  const isMediumUp = useMediaQuery(theme.breakpoints.up("md"));

  const [user, setUser] = useState<User | null>(null);
  const [avatar, setAvatar] = useState<string>("");
  const [books, setBooks] = useState<Book[]>([]);
  const [selectedType, setSelectedType] = useState<FavoriteType>("ALL");

  useEffect(() => {
    if (id) {
      userApi
        .getUser(id)
        .then(setUser)
        .catch((error: Error) => {
          enqueueSnackbar(error.message);
          navigate(AppRoute.NotFound, { replace: true });
        });
    } else {
      navigate(AppRoute.NotFound, { replace: true });
    }
  }, [id]);

  useEffect(() => {
    if (user?.avatarId) {
      imageApi
        .getImage({ imageId: user.avatarId })
        .then((data) => {
          setAvatar(data.imageData);
        })
        .catch((error: Error) => {
          enqueueSnackbar(error.message);
        });
    }
  }, [user?.avatarId]);

  useEffect(() => {
    if (id) {
      const types: UserBookFavoritesType[] =
        selectedType === "ALL" ? ["READING", "PLANNED", "DEFERRED", "READ", "DROPPED", "FAVORITE"] : [selectedType];
      userBookFavoritesApi
        .listBooksByFavorites({ types })
        .then(setBooks)
        .catch((error: Error) => {
          enqueueSnackbar(error.message);
        });
    }
  }, [id, selectedType]);

  return (
    <Box sx={{ width: "100%", py: 4 }}>
      <Container>
        <Paper elevation={3} sx={{ p: 3, mb: 4, borderRadius: 4 }}>
          <Box display="flex" alignItems="center" gap={3}>
            <Avatar src={avatar} sx={{ width: 120, height: 120 }} />
            <Box>
              <Typography variant="h4" gutterBottom>
                {user?.login}
              </Typography>
              <Typography variant="body1" color="text.secondary">
                {user?.aboutMe}
              </Typography>
            </Box>
          </Box>
        </Paper>

        <Paper elevation={3} sx={{ p: 3, borderRadius: 4 }}>
          <Typography variant="h5" gutterBottom>
            Book Lists
          </Typography>
          <Grid2 container spacing={2}>
            <Grid2 order={{ xs: 1, sm: 1 }}>
              <Tabs
                value={selectedType}
                onChange={(_, value: FavoriteType) => {
                  setSelectedType(value);
                }}
                sx={{
                  mt: { xs: 2, md: 0 },
                  mb: { xs: 2, md: 0 },
                  pr: { xs: 2, md: 0 },
                  borderRight: { xs: 1, md: 0 },
                  borderColor: { xs: "divider" },
                  "& .MuiTabs-indicator": {
                    left: { xs: 0, md: "auto" },
                  },
                }}
                orientation={isMediumUp ? "horizontal" : "vertical"}
                scrollButtons="auto"
              >
                <Tab label="All" value="ALL" />
                <Tab label="Reading" value="READING" />
                <Tab label="Planned" value="PLANNED" />
                <Tab label="Deferred" value="DEFERRED" />
                <Tab label="Read" value="READ" />
                <Tab label="Dropped" value="DROPPED" />
                <Tab label="Favorite" value="FAVORITE" />
              </Tabs>
            </Grid2>
            <Grid2 size={{ xs: 12, sm: 9, md: 10 }} order={{ xs: 2, sm: 2 }}>
              <Grid2 container spacing={2}>
                {books.map((book) => (
                  <Grid2 size={{ xs: 12, sm: 6, md: 4, lg: 3 }} key={book.bookId}>
                    <Paper elevation={2} sx={{ p: 2, borderRadius: 2, height: "100%" }}>
                      <Typography variant="h6">{book.title}</Typography>
                      <Typography variant="body2" color="text.secondary">
                        {book.description}
                      </Typography>
                    </Paper>
                  </Grid2>
                ))}
              </Grid2>
            </Grid2>
          </Grid2>
        </Paper>
      </Container>
    </Box>
  );
}
