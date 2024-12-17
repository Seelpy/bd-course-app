import { create } from "zustand";
import { persist, createJSONStorage } from "zustand/middleware";
import { immer } from "zustand/middleware/immer";
import { UserInfo } from "@shared/types/user";

type UserState = {
  userInfo: UserInfo | null;
  setUserInfo: (userInfo: UserInfo | null) => void;
};

export const useUserStore = create<UserState>()(
  persist(
    immer((set) => ({
      userInfo: null,
      setUserInfo: (userInfo: UserInfo | null) => {
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
