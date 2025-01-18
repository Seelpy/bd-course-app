import { userBookFavoritesApi } from "@api/userBookFavorites";
import { useUserStore } from "@shared/stores/userStore";
import { Book } from "@shared/types/book";
import { useSnackbar } from "notistack";
import { useEffect, useState } from "react";
import { useShallow } from "zustand/shallow";

export const RootPage = () => {
  const { enqueueSnackbar } = useSnackbar();

  const { userInfo } = useUserStore(
    useShallow((state) => ({
      userInfo: state.userInfo,
    })),
  );

  const [continueReadingBooks, setContinueReadingBooks] = useState<Book[]>([]);
  const [plannedBooks, setPlannedBooks] = useState<Book[]>([]);

  useEffect(() => {
    if (userInfo) {
      userBookFavoritesApi
        .listBooksByFavorites({ types: ["READING"], userId: userInfo.id })
        .then((data) => {
          setContinueReadingBooks(data.userBookFavouritesBooks.flatMap((item) => item.books));
        })
        .catch((error: Error) => {
          enqueueSnackbar(error.message);
        });

      userBookFavoritesApi
        .listBooksByFavorites({ types: ["PLANNED"], userId: userInfo.id })
        .then((data) => {
          setPlannedBooks(data.userBookFavouritesBooks.flatMap((item) => item.books));
        })
        .catch((error: Error) => {
          enqueueSnackbar(error.message);
        });
    }
  }, [userInfo]);

  console.log(continueReadingBooks);
  console.log(plannedBooks);

  return <></>;
};
