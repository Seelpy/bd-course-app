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
import { Add, Bookmark, PhotoCamera, Edit, EditNote } from "@mui/icons-material";
import { UserBookFavoritesType } from "@shared/types/userBookFavorites";
import { userBookFavoritesApi } from "@api/userBookFavorites";
import { SelectionModal } from "./components/SelectionModal";
import { Author, CreateAuthor } from "@shared/types/author";
import { Genre } from "@shared/types/genre";
import { authorApi } from "@api/author";
import { genreApi } from "@api/genre";
import { bookAuthorApi } from "@api/bookAuthor";
import { bookGenreApi } from "@api/bookGenre";
import { UserRole } from "@shared/types/user";
import { CreateAuthorModal } from "./components/CreateAuthorModal";
import { EditGenreModal } from "./components/EditGenreModal";
import { ConfirmDialog } from "@shared/components/ConfirmDialog"; // Предполагаем, что такой компонент существует
import { EditAuthorModal } from "./components/EditAuthorModal";
import { IconButton } from "@mui/material";
import { UploadImageDialog } from "../../shared/components/UploadImageDialog";
import { imageApi } from "@api/image";
import { EditBookDialog } from "./components/EditBookDialog";

const formatRatingCount = (count: number) => {
  return count >= 1000 ? `${(count / 1000).toFixed(1)}K` : count.toString();
};

export const BookPage = () => {
  const { id } = useParams();
  const navigate = useNavigate();
  const userInfo = useUserStore((state) => state.userInfo);
  const [book, setBook] = useState<Book>();
  const [rating, setRating] = useState<number>(0);
  const [ratingCount, setRatingCount] = useState<number>(0);
  const [userRating, setUserRating] = useState<number>(0);
  const [isRatingModalOpen, setIsRatingModalOpen] = useState(false);
  const [favoriteType, setFavoriteType] = useState<UserBookFavoritesType | null>(null);
  const [menuAnchor, setMenuAnchor] = useState<null | HTMLElement>(null);
  const [menuWidth, setMenuWidth] = useState<number>(0);
  const [authors, setAuthors] = useState<Author[]>([]);
  const [genres, setGenres] = useState<Genre[]>([]);
  const [isAuthorModalOpen, setIsAuthorModalOpen] = useState(false);
  const [isGenreModalOpen, setIsGenreModalOpen] = useState(false);
  const [loading, setLoading] = useState({ authors: false, genres: false });
  const [isCreateAuthorModalOpen, setIsCreateAuthorModalOpen] = useState(false);
  const [editingGenre, setEditingGenre] = useState<Genre | null>(null);
  const [deletingItem, setDeletingItem] = useState<{ type: "genre" | "author"; item: Genre | Author } | null>(null);
  const [editingAuthor, setEditingAuthor] = useState<Author | null>(null);
  const [isUploadImageOpen, setIsUploadImageOpen] = useState(false);
  const [editMode, setEditMode] = useState<"title" | "description" | null>(null);

  const loadRating = () => {
    if (!id) return;
    bookRatingApi
      .getRating(id)
      .then((data) => {
        setRating(data.average);
        setRatingCount(data.count);
        setUserRating(data.userLoginRating);
      })
      .catch(() => {
        enqueueSnackbar("Failed to load rating", { variant: "error" });
      });
  };

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
      loadRating();
    } else {
      navigate(AppRoute.NotFound, { replace: true });
    }
  }, [id]);

  useEffect(() => {
    if (id && userInfo) {
      userBookFavoritesApi
        .getFavoriteTypeByBook({ bookId: id })
        .then((data) => {
          setFavoriteType(data.type);
        })
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

  const handleRate = (value: number) => {
    if (!id) return;
    bookRatingApi
      .updateRating(id, { value })
      .then()
      .catch(() => {
        enqueueSnackbar("Failed to update rating", { variant: "error" });
      });

    setTimeout(() => {
      loadRating();
    }, 500);
  };

  const handleRemoveRating = () => {
    if (!id) return;
    bookRatingApi
      .deleteRating(id)
      .then()
      .catch(() => {
        enqueueSnackbar("Failed to update rating", { variant: "error" });
      });

    setTimeout(() => {
      loadRating();
    }, 500);
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

  const loadAuthors = () => {
    setLoading((prev) => ({ ...prev, authors: true }));

    authorApi
      .listAuthors()
      .then((data) => {
        setAuthors(data.authors);
      })
      .catch((error: Error) => {
        enqueueSnackbar(error.message, { variant: "error" });
      });

    setLoading((prev) => ({ ...prev, authors: false }));
  };

  const loadGenres = () => {
    setLoading((prev) => ({ ...prev, genres: true }));

    genreApi
      .listGenres()
      .then((data) => {
        setGenres(data.genres);
      })
      .catch((error: Error) => {
        enqueueSnackbar(error.message, { variant: "error" });
      });

    setLoading((prev) => ({ ...prev, genres: false }));
  };

  const handleAuthorSelect = (authorId: string) => {
    const isSelected = book.authors.some((a) => a.id === authorId);

    if (isSelected) {
      bookAuthorApi
        .deleteBookAuthor({ bookId: book.bookId, authorId })
        .then(() => {
          setBook((prev) =>
            prev
              ? {
                  ...prev,
                  authors: prev.authors.filter((a) => a.id !== authorId),
                }
              : prev,
          );
        })
        .catch((error: Error) => {
          enqueueSnackbar(error.message, { variant: "error" });
        });
    } else {
      bookAuthorApi
        .storeBookAuthor({ bookId: book.bookId, authorId })
        .then(() => {
          const selectedAuthor = authors.find((a) => a.id === authorId);
          if (selectedAuthor) {
            setBook((prev) =>
              prev
                ? {
                    ...prev,
                    authors: [...prev.authors, selectedAuthor],
                  }
                : prev,
            );
          }
        })
        .catch((error: Error) => {
          enqueueSnackbar(error.message, { variant: "error" });
        });
    }
  };

  const handleGenreSelect = (genreId: string) => {
    const isSelected = book.genres.some((g) => g.id === genreId);

    if (isSelected) {
      bookGenreApi
        .deleteBookGenre({ bookId: book.bookId, genreId })
        .then(() => {
          setBook((prev) =>
            prev
              ? {
                  ...prev,
                  genres: prev.genres.filter((g) => g.id !== genreId),
                }
              : prev,
          );
        })
        .catch((error: Error) => {
          enqueueSnackbar(error.message, { variant: "error" });
        });
    } else {
      bookGenreApi
        .storeBookGenre({ bookId: book.bookId, genreId })
        .then(() => {
          const selectedGenre = genres.find((g) => g.id === genreId);
          if (selectedGenre) {
            setBook((prev) =>
              prev
                ? {
                    ...prev,
                    genres: [...prev.genres, selectedGenre],
                  }
                : prev,
            );
          }
        })
        .catch((error: Error) => {
          enqueueSnackbar(error.message, { variant: "error" });
        });
    }
  };

  const handleCreateAuthor = (author: CreateAuthor) => {
    authorApi
      .createAuthor(author)
      .then(() => {
        loadAuthors();
        enqueueSnackbar("Author created successfully", { variant: "success" });
      })
      .catch((error: Error) => {
        enqueueSnackbar(error.message, { variant: "error" });
      });
  };

  const handleCreateGenre = (name: string) => {
    genreApi
      .createGenre({ name })
      .then(() => {
        loadGenres();
        enqueueSnackbar("Genre created successfully", { variant: "success" });
      })
      .catch((error: Error) => {
        enqueueSnackbar(error.message, { variant: "error" });
      });
  };

  const handleEditGenre = (genreId: string, name: string) => {
    genreApi
      .editGenre({ id: genreId, name })
      .then(() => {
        loadGenres();
        enqueueSnackbar("Genre updated successfully", { variant: "success" });
        setEditingGenre(null);
      })
      .catch((error: Error) => {
        enqueueSnackbar(error.message, { variant: "error" });
      });
  };

  const handleEditAuthor = (
    authorId: string,
    data: { firstName: string; secondName: string; middleName?: string; nickName?: string },
  ) => {
    authorApi
      .editAuthor({ id: authorId, ...data })
      .then(() => {
        loadAuthors();
        enqueueSnackbar("Author updated successfully", { variant: "success" });
        setEditingAuthor(null);
      })
      .catch((error: Error) => {
        enqueueSnackbar(error.message, { variant: "error" });
      });
  };

  const handleDeleteConfirm = () => {
    if (!deletingItem) return;

    if (deletingItem.type === "genre") {
      genreApi
        .deleteGenre({ id: deletingItem.item.id })
        .then(() => {
          loadGenres();
          enqueueSnackbar("Genre deleted successfully", { variant: "success" });
        })
        .catch((error: Error) => {
          enqueueSnackbar(error.message, { variant: "error" });
        });
    } else {
      authorApi
        .deleteAuthor({ id: deletingItem.item.id })
        .then(() => {
          loadAuthors();
          enqueueSnackbar("Author deleted successfully", { variant: "success" });
        })
        .catch((error: Error) => {
          enqueueSnackbar(error.message, { variant: "error" });
        });
    }
    setDeletingItem(null);
  };

  const handleUploadCover = (base64: string) => {
    imageApi
      .storeBookImage({ bookId: book.bookId, imageData: base64 })
      .then(() => {
        setBook((prev) => (prev ? { ...prev, cover: base64 } : prev));
        enqueueSnackbar("Cover updated successfully", { variant: "success" });
      })
      .catch((error: Error) => {
        enqueueSnackbar(error.message, { variant: "error" });
      });
    setIsUploadImageOpen(false);
  };

  const handleEditBook = (data: { title?: string; description?: string }) => {
    bookApi
      .editBook({
        id: book.bookId,
        title: data.title ?? book.title,
        description: data.description ?? book.description,
      })
      .then(() => {
        setBook((prev) => (prev ? { ...prev, ...data } : prev));
        enqueueSnackbar("Book updated successfully", { variant: "success" });
      })
      .catch((error: Error) => {
        enqueueSnackbar(error.message, { variant: "error" });
      });

    setEditMode(null);
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
              {userInfo?.role === UserRole.Admin && (
                <IconButton
                  sx={{
                    position: "absolute",
                    top: 8,
                    right: 8,
                  }}
                  onClick={() => {
                    setIsUploadImageOpen(true);
                  }}
                >
                  <PhotoCamera />
                </IconButton>
              )}
              <Chip
                sx={{
                  position: "absolute",
                  bottom: 8,
                  left: `50%`,
                  transform: `translate(-50%, -50%)`,
                }}
                label={
                  <div style={{ display: "flex", justifyContent: "center", alignItems: "center", gap: 2 }}>
                    {!!userRating && (
                      <>
                        <StarIcon fontSize="small" sx={{ alignContent: "center" }} /> {userRating.toFixed(0)}
                      </>
                    )}
                    <StarIcon color="warning" sx={{ alignContent: "center" }} /> {rating.toFixed(1)} (
                    {formatRatingCount(ratingCount)})
                  </div>
                }
                onClick={handleOpenRatingModal}
              />
            </Box>
            <Box position="relative" width="100%">
              <Typography
                variant="h4"
                component="h1"
                sx={{
                  textAlign: "center",
                  width: "100%",
                }}
              >
                {book.title}
                {userInfo?.role === UserRole.Admin && (
                  <IconButton
                    size="small"
                    sx={{ ml: 1, verticalAlign: "middle" }}
                    onClick={() => {
                      setEditMode("title");
                    }}
                  >
                    <Edit />
                  </IconButton>
                )}
              </Typography>
            </Box>
            <Stack direction="row" spacing={2}>
              <Tooltip title={!userInfo ? "Authorize first" : ""}>
                <span>
                  <Button
                    variant="contained"
                    color="inherit"
                    startIcon={favoriteType ? <Bookmark /> : <Add />}
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
          <Stack spacing={3}>
            <Box display="flex" flexDirection="column" gap={1}>
              <Typography variant="h5" color="textSecondary">
                Authors
              </Typography>
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
                {userInfo?.role === UserRole.Admin && (
                  <Chip
                    icon={<Add />}
                    label="Add Author"
                    onClick={() => {
                      loadAuthors();
                      setIsAuthorModalOpen(true);
                    }}
                  />
                )}
              </Box>
            </Box>

            <Box display="flex" flexDirection="column" gap={1}>
              <Typography variant="h5" color="textSecondary">
                Description
                {userInfo?.role === UserRole.Admin && (
                  <IconButton
                    size="small"
                    sx={{ ml: 1, verticalAlign: "middle" }}
                    onClick={() => {
                      setEditMode("description");
                    }}
                  >
                    <EditNote />
                  </IconButton>
                )}
              </Typography>
              <Typography variant="body1">{book.description}</Typography>
            </Box>
            <Box display="flex" flexDirection="column" gap={1}>
              <Typography variant="h5" color="textSecondary">
                Genres
              </Typography>
              <Box sx={{ display: "flex", flexWrap: "wrap", gap: 1 }}>
                {book.genres.map((genre) => (
                  <Chip key={genre.id} label={genre.name} />
                ))}
                {userInfo?.role === UserRole.Admin && (
                  <Chip
                    icon={<Add />}
                    label="Add Genre"
                    onClick={() => {
                      loadGenres();
                      setIsGenreModalOpen(true);
                    }}
                  />
                )}
              </Box>
            </Box>
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

      <SelectionModal
        open={isAuthorModalOpen}
        onClose={() => {
          setIsAuthorModalOpen(false);
        }}
        title="Select Authors"
        items={authors.map((a) => ({ id: a.id, name: `${a.firstName} ${a.secondName}` }))}
        selectedIds={book.authors.map((a) => a.id)}
        onSelect={handleAuthorSelect}
        loading={loading.authors}
        searchPlaceholder="Search authors..."
        createButtonText="Create New Author"
        forceCreateNewButton={true}
        onCreate={() => {
          setIsCreateAuthorModalOpen(true);
          setIsAuthorModalOpen(false);
        }}
        onEdit={(author) => {
          const fullAuthor = authors.find((a) => a.id === author.id);
          if (fullAuthor) {
            setEditingAuthor(fullAuthor);
          }
        }}
        onDelete={(author) => {
          const fullAuthor = authors.find((a) => a.id === author.id);
          if (fullAuthor) {
            setDeletingItem({ type: "author", item: fullAuthor });
          }
        }}
      />

      <CreateAuthorModal
        open={isCreateAuthorModalOpen}
        onClose={() => {
          setIsCreateAuthorModalOpen(false);
          setIsAuthorModalOpen(true);
        }}
        onCreate={handleCreateAuthor}
      />

      <SelectionModal
        open={isGenreModalOpen}
        onClose={() => {
          setIsGenreModalOpen(false);
        }}
        title="Select Genres"
        items={genres.map((g) => ({ id: g.id, name: g.name }))}
        selectedIds={book.genres.map((g) => g.id)}
        onSelect={handleGenreSelect}
        onCreate={handleCreateGenre}
        loading={loading.genres}
        searchPlaceholder="Search or enter new genre name..."
        createButtonText="Create Genre"
        onEdit={(genre) => {
          setEditingGenre(genre);
        }}
        onDelete={(genre) => {
          setDeletingItem({ type: "genre", item: genre });
        }}
      />

      {editingGenre && (
        <EditGenreModal
          open={!!editingGenre}
          onClose={() => {
            setEditingGenre(null);
          }}
          genre={editingGenre}
          onEdit={handleEditGenre}
        />
      )}

      {editingAuthor && (
        <EditAuthorModal
          open={!!editingAuthor}
          onClose={() => {
            setEditingAuthor(null);
          }}
          author={editingAuthor}
          onEdit={handleEditAuthor}
        />
      )}

      <ConfirmDialog
        open={!!deletingItem}
        onClose={() => {
          setDeletingItem(null);
        }}
        onConfirm={handleDeleteConfirm}
        title={`Delete ${deletingItem?.type ?? ""}`}
        content={`Are you sure you want to delete? This action cannot be undone.`}
      />

      <UploadImageDialog
        open={isUploadImageOpen}
        onClose={() => {
          setIsUploadImageOpen(false);
        }}
        onUpload={handleUploadCover}
        title="Upload Book Cover"
      />

      {editMode && (
        <EditBookDialog
          open={!!editMode}
          onClose={() => {
            setEditMode(null);
          }}
          onSave={handleEditBook}
          title={book.title}
          description={book.description}
          mode={editMode}
        />
      )}
    </Container>
  );
};
