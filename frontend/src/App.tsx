import "./App.css";
import Layout from "./layout/Layout";
import CssBaseline from "@mui/material/CssBaseline";
import { useTheme } from "@shared/hooks/useTheme.ts";
import { ThemeProvider } from "@mui/material";
import { RouterProvider, createBrowserRouter } from "react-router-dom";
import { AppRoute } from "@shared/constants/routes";
import { SnackbarProvider } from "notistack";
import ErrorBoundary from "./pages/ErrorBoundary.tsx";
import { useEffect } from "react";
import { userApi } from "@api/user.ts";
import { useUserStore } from "@shared/stores/userStore.ts";
import { useShallow } from "zustand/shallow";

function App() {
  const { theme } = useTheme();

  const { userInfo, setUserInfo } = useUserStore(
    useShallow((state) => ({
      userInfo: state.userInfo,
      setUserInfo: state.setUserInfo,
    })),
  );

  useEffect(() => {
    if (userInfo) {
      userApi
        .getAuthorizedUser()
        .then((data) => {
          setUserInfo(data);
        })
        .catch(() => {
          setUserInfo(null);
        });
    }
  }, []);

  const routes = createBrowserRouter([
    {
      element: <Layout />,
      errorElement: <ErrorBoundary />,
      children: [
        {
          path: AppRoute.Root,
          lazy: () => import("./pages/RootPage/RootPage.tsx").then((m) => ({ Component: m.RootPage })),
        },
        {
          path: AppRoute.NotFound,
          lazy: () => import("./pages/PageNotFound.tsx").then((m) => ({ Component: m.PageNotFound })),
        },
        {
          path: AppRoute.Login,
          lazy: () => import("./pages/LoginPage.tsx").then((m) => ({ Component: m.LoginPage })),
        },
        {
          path: AppRoute.Register,
          lazy: () => import("./pages/RegisterPage.tsx").then((m) => ({ Component: m.RegisterPage })),
        },
        {
          path: AppRoute.Profile,
          lazy: () => import("./pages/ProfilePage.tsx").then((m) => ({ Component: m.ProfilePage })),
        },
        {
          path: AppRoute.Requests,
          lazy: () => import("./pages/RequestsPage.tsx").then((m) => ({ Component: m.RequestsPage })),
        },
      ],
    },
  ]);

  console.log(routes);

  return (
    <SnackbarProvider maxSnack={3} anchorOrigin={{ vertical: "bottom", horizontal: "right" }}>
      <ThemeProvider theme={theme}>
        <CssBaseline />
        <RouterProvider router={routes} />
      </ThemeProvider>
    </SnackbarProvider>
  );
}

export default App;
