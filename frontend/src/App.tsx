import "./App.css";
import Layout from "./layout/Layout";
import CssBaseline from "@mui/material/CssBaseline";
import { useTheme } from "@shared/hooks/useTheme.ts";
import { ThemeProvider } from "@mui/material";
import { RouterProvider, createBrowserRouter } from "react-router-dom";
import { AppRoute } from "@shared/constants/routes";
import { SnackbarProvider } from "notistack";
import ErrorBoundary from "./pages/ErrorBoundary.tsx";

function App() {
  const { mode, theme } = useTheme();
  console.log(mode);

  const routes = createBrowserRouter([
    {
      element: <Layout />,
      errorElement: <ErrorBoundary />,
      children: [
        {
          path: AppRoute.Root,
          lazy: () => import("./pages/RootPage.tsx").then((m) => ({ Component: m.RootPage })),
        },
        {
          path: AppRoute.NotFound,
          lazy: () => import("./pages/PageNotFound.tsx").then((m) => ({ Component: m.PageNotFound })),
        },
      ],
    },
  ]);

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
