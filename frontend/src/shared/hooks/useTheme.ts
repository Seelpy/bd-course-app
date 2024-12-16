import { useMemo } from "react";
import { usePreferencesStore } from "@shared/stores/preferencesStore";
import { useMediaQuery } from "@mui/material";
import { createTheme } from "@mui/material/styles";

export function useTheme() {
  const preferredTheme = usePreferencesStore((state) => state.theme);
  const prefersDarkMode = useMediaQuery("(prefers-color-scheme: dark)");

  const { mode, theme } = useMemo(() => {
    const updatePreferredDark = prefersDarkMode || preferredTheme !== "auto";
    const resolvedMode = preferredTheme === "auto" ? (updatePreferredDark ? "dark" : "light") : preferredTheme;

    const createdTheme = createTheme({
      palette: {
        mode: resolvedMode,
      },
    });

    return { mode: resolvedMode, theme: createdTheme };
  }, [preferredTheme, prefersDarkMode]);

  return { mode, theme };
}
