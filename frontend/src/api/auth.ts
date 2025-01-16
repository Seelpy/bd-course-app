import { AuthForm } from "@shared/types/auth";
import { handleApiError } from "./utils/handleApiError";

export const authApi = {
  PREFIX: "/api/v1/auth",

  login(form: AuthForm): Promise<unknown> {
    return fetch(this.PREFIX, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify(form),
    }).then((res) => {
      if (!res.ok) {
        return handleApiError(res);
      } else {
        return res.json();
      }
    });
  },
  refreshToken(): Promise<unknown> {
    return fetch(this.PREFIX, {
      method: "PATCH",
      credentials: "include",
      headers: {
        "Content-Type": "application/json",
      },
    }).then((res) => {
      if (!res.ok) {
        return handleApiError(res);
      } else {
        return res.json();
      }
    });
  },
  logout(): Promise<unknown> {
    return fetch(this.PREFIX, {
      method: "DELETE",
      credentials: "include",
      headers: {
        "Content-Type": "application/json",
      },
    }).then((res) => {
      if (!res.ok) {
        return handleApiError(res);
      } else {
        return res.json();
      }
    });
  },
} as const;
