import { create } from "zustand";
import { persist, createJSONStorage } from "zustand/middleware";
import { immer } from "zustand/middleware/immer";
import { Theme } from "@shared/types/Theme";

type PreferencesState = {
  theme: Theme;
  changeTheme: (theme: Theme) => void;
};

export const usePreferencesStore = create<PreferencesState>()(
  persist(
    immer((set) => ({
      theme: "auto",
      changeTheme: (theme: Theme) => {
        set((state) => {
          state.theme = theme;
        });
      },
    })),
    {
      name: "bdapp",
      storage: createJSONStorage(() => localStorage),
    },
  ),
);
