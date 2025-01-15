import { create } from "zustand";
import { persist, createJSONStorage } from "zustand/middleware";
import { immer } from "zustand/middleware/immer";
import { User } from "@shared/types/user";

type UserState = {
  userInfo: User | null;
  setUserInfo: (userInfo: User | null) => void;
};

export const useUserStore = create<UserState>()(
  persist(
    immer((set) => ({
      userInfo: null,
      setUserInfo: (userInfo: User | null) => {
        set((state) => {
          state.userInfo = userInfo;
        });
      },
    })),
    {
      name: "bdapp",
      storage: createJSONStorage(() => localStorage),
    },
  ),
);
