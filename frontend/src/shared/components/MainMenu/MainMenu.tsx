import { useMediaQuery, useTheme } from "@mui/material";
import { AutoStories, LibraryAddCheck } from "@mui/icons-material";
import { useNavigate } from "react-router-dom";
import { AppRoute } from "@shared/constants/routes";
import { useUserStore } from "@shared/stores/userStore";
import { authApi } from "@api/auth";
import { useSnackbar } from "notistack";
import { MenuDesktop } from "./MenuDesktop";
import { MenuMobile } from "./MenuMobile";
import { useShallow } from "zustand/shallow";

export const MainMenu = () => {
  const theme = useTheme();
  const isMobile = useMediaQuery(theme.breakpoints.down("md"));
  const navigate = useNavigate();
  const { enqueueSnackbar } = useSnackbar();

  const { userInfo, setUserInfo } = useUserStore(
    useShallow((state) => ({
      userInfo: state.userInfo,
      setUserInfo: state.setUserInfo,
    })),
  );

  const handleLogout = () => {
    authApi
      .logout()
      .then(() => {
        setUserInfo(null);
        navigate(AppRoute.Root);
      })
      .catch((error: Error) => {
        enqueueSnackbar(error.message, { variant: "error" });
      });
  };

  const menuItems = [
    {
      text: "Catalog",
      icon: <AutoStories />,
      onClick: () => {
        console.log("TODO");
      },
    },
    {
      text: "Requests",
      icon: <LibraryAddCheck />,
      onClick: () => {
        console.log("TODO");
      },
    },
  ];

  return isMobile ? (
    <MenuMobile userInfo={userInfo} handleLogout={handleLogout} menuItems={menuItems} />
  ) : (
    <MenuDesktop userInfo={userInfo} handleLogout={handleLogout} menuItems={menuItems} />
  );
};
