import { CreateUser, EditUser, DeleteUser, User } from "@shared/types/user";
import { handleApiError } from "./utils/handleApiError";

export const userApi = {
  PREFIX: "/api/v1/user",

  getUser(id: string): Promise<User> {
    return fetch(`${this.PREFIX}/${id}`, {
      method: "GET",
      credentials: "include",
    }).then((res) => {
      if (!res.ok) {
        return handleApiError(res, () => this.getUser(id));
      }
      return res.json();
    });
  },

  listUsers(): Promise<User[]> {
    return fetch(this.PREFIX, {
      method: "GET",
      credentials: "include",
    }).then((res) => {
      if (!res.ok) {
        return handleApiError(res, () => this.listUsers());
      }
      return res.json();
    });
  },

  createUser(body: CreateUser): Promise<unknown> {
    return fetch(this.PREFIX, {
      method: "POST",
      credentials: "include",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify(body),
    }).then((res) => {
      if (!res.ok) {
        return handleApiError(res, () => this.createUser(body));
      }
      return res.json();
    });
  },

  editUser(body: EditUser): Promise<unknown> {
    return fetch(this.PREFIX, {
      method: "PATCH",
      credentials: "include",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify(body),
    }).then((res) => {
      if (!res.ok) {
        return handleApiError(res, () => this.editUser(body));
      }
      return res.json();
    });
  },

  deleteUser(body: DeleteUser): Promise<unknown> {
    return fetch(this.PREFIX, {
      method: "DELETE",
      credentials: "include",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify(body),
    }).then((res) => {
      if (!res.ok) {
        return handleApiError(res, () => this.deleteUser(body));
      }
      return res.json();
    });
  },
} as const;
