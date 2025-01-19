import { useEffect, useState, useCallback, useRef } from "react";
import {
  Container,
  Paper,
  Typography,
  TextField,
  Box,
  Radio,
  RadioGroup,
  FormControlLabel,
  FormControl,
  FormLabel,
  Drawer,
  useTheme,
  useMediaQuery,
  Chip,
  Grid2,
  styled,
  Menu,
  Button,
} from "@mui/material";
import FilterListIcon from "@mui/icons-material/FilterList";
import SortIcon from "@mui/icons-material/Sort";
import { bookApi } from "@api/book";
import { authorApi } from "@api/author";
import { genreApi } from "@api/genre";
import { Book, SortBy, SortType } from "@shared/types/book";
import { Author } from "@shared/types/author";
import { Genre } from "@shared/types/genre";
import { BookPreview } from "@shared/components/BookPreview/BookPreview";
import { BookPreviewSkeleton } from "@shared/components/BookPreview/BookPreviewSkeleton";
import { useSnackbar } from "notistack";
import { debounce } from "lodash";
import { useBookWidth } from "@shared/hooks/useBookWidth";

const HalfedTextField = styled(TextField)({
  width: "50%",
  "& input::-webkit-outer-spin-button, & input::-webkit-inner-spin-button": {
    display: "none",
  },
  "& input[type=number]": {
    MozAppearance: "textfield",
  },
});

const BOOKS_PER_PAGE = 20;

type Range = {
  from?: number;
  to?: number;
};

export function CatalogPage() {
  const theme = useTheme();
  const isMediumUp = useMediaQuery(theme.breakpoints.up("md"));
  const { enqueueSnackbar } = useSnackbar();
  const [isFilterOpen, setIsFilterOpen] = useState(false);
  const [books, setBooks] = useState<Book[]>([]);
  const [authors, setAuthors] = useState<Author[]>([]);
  const [genres, setGenres] = useState<Genre[]>([]);
  const [loading, setLoading] = useState(false);
  const [loadingDebounced, setLoadingDebounced] = useState(false);
  const [page, setPage] = useState(0);
  const [hasMore, setHasMore] = useState(true);
  const observer = useRef<IntersectionObserver>();
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

  // Search params
  const [searchTitle, setSearchTitle] = useState("");
  const [selectedAuthorIds, setSelectedAuthorIds] = useState<string[]>([]);
  const [selectedGenreIds, setSelectedGenreIds] = useState<string[]>([]);
  const [chaptersRange, setChaptersRange] = useState<Range>({});
  const [ratingRange, setRatingRange] = useState<Range>({});
  const [ratingCountRange, setRatingCountRange] = useState<Range>({});
  const [sortBy, setSortBy] = useState<SortBy>("TITLE");
  const [sortType, setSortType] = useState<SortType>("ASC");
  const [anchorEl, setAnchorEl] = useState<null | HTMLElement>(null);
  const open = Boolean(anchorEl);

  const handleSortClick = (event: React.MouseEvent<HTMLElement>) => {
    setAnchorEl(event.currentTarget);
  };

  const handleClose = () => {
    setAnchorEl(null);
  };

  const loadBooks = useCallback(
    (newSearch = false) => {
      if (loading || (!hasMore && !newSearch)) return;

      let currentPage = 1;

      if (newSearch) {
        setPage(1);
      } else {
        currentPage = page + 1;
        setPage((prev) => prev + 1);
      }

      setLoading(true);
      setLoadingDebounced(true);

      bookApi
        .searchBooks(currentPage, BOOKS_PER_PAGE, {
          bookTitle: searchTitle,
          authorIds: selectedAuthorIds.length ? selectedAuthorIds : undefined,
          genreIds: selectedGenreIds.length ? selectedGenreIds : undefined,
          minChaptersCount: chaptersRange.from,
          maxChaptersCount: chaptersRange.to,
          minRating: ratingRange.from,
          maxRating: ratingRange.to,
          minRatingCount: ratingCountRange.from,
          maxRatingCount: ratingCountRange.to,
          sortBy,
          sortType,
        })
        .then((result) => {
          const filteredBooks = result.books.filter((book) => {
            const hasSelectedAuthors =
              selectedAuthorIds.length === 0 || book.authors.some((author) => selectedAuthorIds.includes(author.id));

            const hasSelectedGenres =
              selectedGenreIds.length === 0 || book.genres.some((genre) => selectedGenreIds.includes(genre.id));

            return hasSelectedAuthors && hasSelectedGenres;
          });

          setBooks((prev) => (!newSearch ? [...prev, ...filteredBooks] : filteredBooks));
          setHasMore(filteredBooks.length === BOOKS_PER_PAGE);
        })
        .catch((error: Error) => {
          enqueueSnackbar(error.message, { variant: "error" });
        })
        .finally(() => {
          setLoading(false);
          setTimeout(() => {
            setLoadingDebounced(false);
          }, 300);
        });
    },
    [
      searchTitle,
      selectedAuthorIds,
      selectedGenreIds,
      chaptersRange,
      ratingRange,
      ratingCountRange,
      sortBy,
      sortType,
      page,
    ],
  );

  const debouncedSearch = debounce(() => {
    loadBooks(true);
  }, 500);

  useEffect(() => {
    Promise.all([authorApi.listAuthors(), genreApi.listGenres()])
      .then(([authorsData, genresData]) => {
        setAuthors(authorsData.authors);
        setGenres(genresData.genres);
      })
      .catch((error) => {
        enqueueSnackbar((error as Error).message, { variant: "error" });
      });
  }, []);

  useEffect(() => {
    if (page === 0) loadBooks(true);
    else debouncedSearch();

    return () => {
      debouncedSearch.cancel();
    };
  }, [
    searchTitle,
    selectedAuthorIds,
    selectedGenreIds,
    chaptersRange,
    ratingRange,
    ratingCountRange,
    sortBy,
    sortType,
  ]);

  const lastBookElementRef = useCallback(
    (node: HTMLDivElement | null) => {
      if (!node || loading) return;

      if (observer.current) observer.current.disconnect();

      observer.current = new IntersectionObserver((entries) => {
        if (entries[0].isIntersecting && hasMore) {
          loadBooks();
        }
      });

      observer.current.observe(node);
    },
    [loading, hasMore],
  );

  const FiltersContent = (
    <Box
      sx={{
        width: isMediumUp ? 300 : "100%",
        height: "auto",
        display: "flex",
        flexDirection: "column",
      }}
    >
      <Typography variant="h5" sx={{ mb: 2, flex: "0 0 auto" }}>
        Filters
      </Typography>

      <Box
        sx={{
          overflow: "auto",
          flex: 1,
          "&::-webkit-scrollbar": {
            width: "8px",
          },
          "&::-webkit-scrollbar-track": {
            background: "transparent",
          },
          "&::-webkit-scrollbar-thumb": {
            background: (theme) => theme.palette.divider,
            borderRadius: "4px",
          },
        }}
      >
        {authors.length > 0 && (
          <Box sx={{ mb: 3 }}>
            <Typography gutterBottom>Authors</Typography>
            <Box sx={{ display: "flex", flexWrap: "wrap", gap: 1 }}>
              {authors.map((author) => (
                <Chip
                  key={author.id}
                  label={`${author.firstName} ${author.secondName}`}
                  onClick={() => {
                    setSelectedAuthorIds((prev) =>
                      prev.includes(author.id) ? prev.filter((id) => id !== author.id) : [...prev, author.id],
                    );
                  }}
                  color={selectedAuthorIds.includes(author.id) ? "primary" : "default"}
                />
              ))}
            </Box>
          </Box>
        )}

        {genres.length > 0 && (
          <Box sx={{ mb: 3 }}>
            <Typography gutterBottom>Genres</Typography>
            <Box sx={{ display: "flex", flexWrap: "wrap", gap: 1 }}>
              {genres.map((genre) => (
                <Chip
                  key={genre.id}
                  label={genre.name}
                  onClick={() => {
                    setSelectedGenreIds((prev) =>
                      prev.includes(genre.id) ? prev.filter((id) => id !== genre.id) : [...prev, genre.id],
                    );
                  }}
                  color={selectedGenreIds.includes(genre.id) ? "primary" : "default"}
                />
              ))}
            </Box>
          </Box>
        )}

        <Box sx={{ mb: 3 }}>
          <Typography marginBottom={2}>Chapters count</Typography>
          <Box sx={{ display: "flex", gap: 1 }}>
            <HalfedTextField
              type="number"
              size="small"
              placeholder="From"
              value={chaptersRange.from ?? ""}
              onChange={(e) => {
                const value = e.target.value === "" ? undefined : Number(e.target.value);
                setChaptersRange((prev) => ({ ...prev, from: value }));
              }}
            />
            <HalfedTextField
              type="number"
              size="small"
              placeholder="To"
              value={chaptersRange.to ?? ""}
              onChange={(e) => {
                const value = e.target.value === "" ? undefined : Number(e.target.value);
                setChaptersRange((prev) => ({ ...prev, to: value }));
              }}
            />
          </Box>
        </Box>

        <Box sx={{ mb: 3 }}>
          <Typography marginBottom={2}>Rating</Typography>
          <Box sx={{ display: "flex", gap: 1 }}>
            <HalfedTextField
              type="number"
              size="small"
              placeholder="From"
              value={ratingRange.from ?? ""}
              onChange={(e) => {
                const value = e.target.value === "" ? undefined : Number(e.target.value);
                setRatingRange((prev) => ({ ...prev, from: value }));
              }}
            />
            <HalfedTextField
              type="number"
              size="small"
              placeholder="To"
              value={ratingRange.to ?? ""}
              onChange={(e) => {
                const value = e.target.value === "" ? undefined : Number(e.target.value);
                setRatingRange((prev) => ({ ...prev, to: value }));
              }}
            />
          </Box>
        </Box>

        <Box sx={{ mb: 3 }}>
          <Typography marginBottom={2}>Rating count</Typography>
          <Box sx={{ display: "flex", gap: 1 }}>
            <HalfedTextField
              type="number"
              size="small"
              placeholder="From"
              value={ratingCountRange.from ?? ""}
              onChange={(e) => {
                const value = e.target.value === "" ? undefined : Number(e.target.value);
                setRatingCountRange((prev) => ({ ...prev, from: value }));
              }}
            />
            <HalfedTextField
              type="number"
              size="small"
              placeholder="To"
              value={ratingCountRange.to ?? ""}
              onChange={(e) => {
                const value = e.target.value === "" ? undefined : Number(e.target.value);
                setRatingCountRange((prev) => ({ ...prev, to: value }));
              }}
            />
          </Box>
        </Box>

        {!isMediumUp && (
          <Button
            variant="contained"
            fullWidth
            onClick={() => {
              setIsFilterOpen(false);
            }}
          >
            Apply
          </Button>
        )}
      </Box>
    </Box>
  );

  return (
    <Container sx={{ py: 3 }}>
      <Box sx={{ display: "flex", gap: 2 }}>
        <Paper
          elevation={3}
          sx={{
            p: 3,
            flexGrow: 1,
            borderRadius: 2,
          }}
        >
          <Box sx={{ display: "flex", alignItems: "center", mb: 2 }}>
            <Typography variant="h5" sx={{ flexGrow: 1 }}>
              Catalog
            </Typography>

            <Box sx={{ display: "flex", gap: 1 }}>
              <Button startIcon={<SortIcon />} onClick={handleSortClick} variant="outlined" size="small">
                Sorting
              </Button>

              <Menu
                anchorEl={anchorEl}
                open={open}
                onClose={handleClose}
                anchorOrigin={{
                  vertical: "bottom",
                  horizontal: "right",
                }}
                transformOrigin={{
                  vertical: "top",
                  horizontal: "right",
                }}
              >
                <Box sx={{ p: 2, minWidth: 200 }}>
                  <FormControl component="fieldset" sx={{ mb: 2 }}>
                    <FormLabel>Sort by</FormLabel>
                    <RadioGroup
                      value={sortBy}
                      onChange={(e) => {
                        setSortBy(e.target.value as SortBy);
                      }}
                    >
                      <FormControlLabel value="TITLE" control={<Radio />} label="Title" />
                      <FormControlLabel value="RATING" control={<Radio />} label="Rating" />
                      <FormControlLabel value="RATING_COUNT" control={<Radio />} label="Rating count" />
                      <FormControlLabel value="CHAPTERS_COUNT" control={<Radio />} label="Chapters" />
                    </RadioGroup>
                  </FormControl>

                  <FormControl component="fieldset">
                    <FormLabel>Order</FormLabel>
                    <RadioGroup
                      value={sortType}
                      onChange={(e) => {
                        setSortType(e.target.value as SortType);
                      }}
                    >
                      <FormControlLabel value="ASC" control={<Radio />} label="Ascending" />
                      <FormControlLabel value="DESC" control={<Radio />} label="Descending" />
                    </RadioGroup>
                  </FormControl>
                </Box>
              </Menu>

              {!isMediumUp && (
                <Button
                  startIcon={<FilterListIcon />}
                  onClick={() => {
                    setIsFilterOpen(true);
                  }}
                  variant="outlined"
                  size="small"
                >
                  Filters
                </Button>
              )}
            </Box>
          </Box>

          <TextField
            fullWidth
            label="Search by name"
            value={searchTitle}
            onChange={(e) => {
              setSearchTitle(e.target.value);
            }}
            sx={{ mb: 3 }}
          />

          <Box ref={containerRef}>
            <Grid2 container spacing={2}>
              {!loadingDebounced &&
                books.map((book, index) => (
                  <Grid2 key={book.bookId} ref={index === books.length - 1 ? lastBookElementRef : undefined}>
                    <BookPreview book={book} width={bookWidth} />
                  </Grid2>
                ))}
              {loadingDebounced &&
                Array.from(new Array(6)).map((_, index) => (
                  <Grid2 key={index}>
                    <BookPreviewSkeleton width={bookWidth} />
                  </Grid2>
                ))}
            </Grid2>
          </Box>
        </Paper>

        {isMediumUp ? (
          <Paper
            elevation={3}
            sx={{
              flexGrow: 0,
              p: 3,
              borderRadius: 2,
              position: "sticky",
              alignSelf: "flex-start",
            }}
          >
            {FiltersContent}
          </Paper>
        ) : (
          <Drawer
            anchor="right"
            open={isFilterOpen}
            onClose={() => {
              setIsFilterOpen(false);
            }}
          >
            <Box sx={{ p: 2 }}>{FiltersContent}</Box>
          </Drawer>
        )}
      </Box>
    </Container>
  );
}
