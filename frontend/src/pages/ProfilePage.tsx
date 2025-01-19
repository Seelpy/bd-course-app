import { useEffect, useRef, useState } from "react";
import { useParams, useNavigate } from "react-router-dom";
import {
  Container,
  Paper,
  Avatar,
  Typography,
  Box,
  Tabs,
  Tab,
  Grid2,
  useMediaQuery,
  useTheme,
  IconButton,
  Button,
  Stack,
  TextareaAutosize,
  styled,
} from "@mui/material";
import EditIcon from "@mui/icons-material/Edit";
import { useUserStore } from "@shared/stores/userStore";
import { useShallow } from "zustand/shallow";
import { userApi } from "@api/user";
import { userBookFavoritesApi } from "@api/userBookFavorites";
import { User } from "@shared/types/user";
import { Book } from "@shared/types/book";
import { UserBookFavoritesType } from "@shared/types/userBookFavorites";
import { useSnackbar } from "notistack";
import { AppRoute } from "@shared/constants/routes";
import { BookPreview } from "@shared/components/BookPreview/BookPreview";
import { useBookWidth } from "@shared/hooks/useBookWidth";

type FavoriteType = UserBookFavoritesType | "ALL";

const StyledTextarea = styled(TextareaAutosize)(({ theme }) => ({
  width: "100%",
  padding: "8px 12px",
  borderRadius: "4px",
  borderColor: theme.palette.divider,
  backgroundColor: theme.palette.background.paper,
  color: theme.palette.text.primary,
  fontFamily: theme.typography.fontFamily,
  fontSize: theme.typography.body1.fontSize,
  resize: "none",
  "&:focus": {
    outline: "none",
    borderColor: theme.palette.primary.main,
    borderWidth: "2px",
  },
  transition: theme.transitions.create(["border-color", "box-shadow"]),
}));

export function ProfilePage() {
  const { id } = useParams<{ id: string }>();
  const { enqueueSnackbar } = useSnackbar();
  const navigate = useNavigate();
  const theme = useTheme();
  const isMediumUp = useMediaQuery(theme.breakpoints.up("md"));

  const [user, setUser] = useState<User | null>(null);
  const [books, setBooks] = useState<Book[]>([]);
  const [selectedType, setSelectedType] = useState<FavoriteType>("ALL");
  const [isEditing, setIsEditing] = useState(false);
  const [newAboutMe, setNewAboutMe] = useState("");
  const { userInfo } = useUserStore(useShallow((state) => ({ userInfo: state.userInfo })));

  const containerRef = useRef<HTMLDivElement>(null);
  const [containerWidth, setContainerWidth] = useState(0);

  useEffect(() => {
    const updateWidth = () => {
      if (containerRef.current) {
        setContainerWidth(containerRef.current.offsetWidth);
      }
    };

    updateWidth();
    window.addEventListener("resize", updateWidth);
    return () => {
      window.removeEventListener("resize", updateWidth);
    };
  }, []);

  const { bookWidth } = useBookWidth(containerWidth);
  const isOwnProfile = userInfo?.id === id;

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
    if (id) {
      const types: UserBookFavoritesType[] =
        selectedType === "ALL" ? ["READING", "PLANNED", "DEFERRED", "READ", "DROPPED", "FAVORITE"] : [selectedType];
      userBookFavoritesApi
        .listBooksByFavorites({ types, userId: id })
        .then((data) => {
          setBooks(data.userBookFavouritesBooks.flatMap((item) => item.books));
        })
        .catch((error: Error) => {
          enqueueSnackbar(error.message);
        });
    }
  }, [id, selectedType]);

  useEffect(() => {
    if (user?.aboutMe) {
      setNewAboutMe(user.aboutMe);
    }
  }, [user?.aboutMe]);

  const handleEditSubmit = () => {
    if (!user || !userInfo) return;
    userApi
      .editUser({
        id: userInfo.id,
        aboutMe: newAboutMe,
      })
      .then(() => {
        setUser({ ...user, aboutMe: newAboutMe });
        setIsEditing(false);
        enqueueSnackbar("Profile updated successfully", { variant: "success" });
      })
      .catch((error: Error) => {
        enqueueSnackbar((error.message, { variant: "error" }));
      });
  };

  return (
    <Box sx={{ width: "100%", py: 4 }}>
      <Container>
        <Paper elevation={3} sx={{ p: 3, mb: 4, borderRadius: 4 }}>
          <Box display="flex" alignItems="center" gap={3}>
            <Avatar src={user?.avatar} sx={{ width: 120, height: 120 }} />
            <Box sx={{ flex: 1 }}>
              <Typography variant="h4" gutterBottom>
                {user?.login}
              </Typography>
              {isEditing ? (
                <Box sx={{ display: "flex", gap: 1, alignItems: "flex-start" }}>
                  <StyledTextarea
                    value={newAboutMe}
                    onChange={(e) => {
                      setNewAboutMe(e.target.value);
                    }}
                    minRows={3}
                  />
                  <Stack direction="column" spacing={1}>
                    <Button variant="contained" onClick={handleEditSubmit}>
                      Save
                    </Button>
                    <Button
                      variant="outlined"
                      onClick={() => {
                        setIsEditing(false);
                        setNewAboutMe(user?.aboutMe ?? "");
                      }}
                    >
                      Cancel
                    </Button>
                  </Stack>
                </Box>
              ) : (
                <Box sx={{ display: "flex", alignItems: "flex-start", gap: 1 }}>
                  <Typography variant="body1" color="text.secondary">
                    {user?.aboutMe}
                  </Typography>
                  {isOwnProfile && user?.aboutMe.trim() && (
                    <IconButton
                      size="small"
                      onClick={() => {
                        setIsEditing(true);
                      }}
                    >
                      <EditIcon fontSize="small" />
                    </IconButton>
                  )}
                  {isOwnProfile && !user?.aboutMe.trim() && (
                    <Button
                      variant="text"
                      endIcon={<EditIcon />}
                      sx={{ marginLeft: -1 }}
                      onClick={() => {
                        setIsEditing(true);
                      }}
                    >
                      Write something about yourself
                    </Button>
                  )}
                </Box>
              )}
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
            <Grid2 size={{ xs: 12, sm: 9, md: 10 }} order={{ xs: 2, sm: 2 }} ref={containerRef}>
              <Grid2 container spacing={2}>
                {books.map((book) => (
                  <BookPreview key={book.bookId} book={book} width={bookWidth} />
                ))}
              </Grid2>
            </Grid2>
          </Grid2>
        </Paper>
      </Container>
    </Box>
  );
}
