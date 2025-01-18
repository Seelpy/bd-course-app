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
import { useEffect, useState } from "react";
import { imageApi } from "@api/image";

export const MainMenu = () => {
  const theme = useTheme();
  const isMobile = useMediaQuery(theme.breakpoints.down("md"));
  const navigate = useNavigate();
  const { enqueueSnackbar } = useSnackbar();
  const [avatar, setAvatar] = useState<string>("");

  const { userInfo, setUserInfo } = useUserStore(
    useShallow((state) => ({
      userInfo: state.userInfo,
      setUserInfo: state.setUserInfo,
    })),
  );

  useEffect(() => {
    if (userInfo?.avatarId) {
      imageApi
        .getImage({ imageId: userInfo.avatarId })
        .then((data) => {
          setAvatar(data.imageData);
        })
        // eslint-disable-next-line @typescript-eslint/no-empty-function
        .catch(() => {});
    }
  }, [userInfo?.avatarId]);

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
        navigate(AppRoute.Requests);
      },
    },
  ];

  return isMobile ? (
    <MenuMobile userInfo={userInfo} avatar={avatar} handleLogout={handleLogout} menuItems={menuItems} />
  ) : (
    <MenuDesktop userInfo={userInfo} avatar={avatar} handleLogout={handleLogout} menuItems={menuItems} />
  );
};
