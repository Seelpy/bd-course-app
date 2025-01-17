import { userBookFavoritesApi } from "@api/userBookFavorites";
import { useUserStore } from "@shared/stores/userStore";
import { useSnackbar } from "notistack";
import { useShallow } from "zustand/shallow";

export const RootPage = () => {
  const { enqueueSnackbar } = useSnackbar();

  const { userInfo } = useUserStore(
    useShallow((state) => ({
      userInfo: state.userInfo,
    })),
  );

  let continueReadingBooks = null;
  let plannedBooks = null;

  if (userInfo) {
    userBookFavoritesApi
      .listBooksByFavorites({ types: ["READING"] })
      .then((books) => {
        continueReadingBooks = books;
      })
      .catch((error: Error) => {
        enqueueSnackbar(error.message);
      });

    userBookFavoritesApi
      .listBooksByFavorites({ types: ["PLANNED"] })
      .then((books) => {
        plannedBooks = books;
      })
      .catch((error: Error) => {
        enqueueSnackbar(error.message);
      });
  }

  return <></>;
};
