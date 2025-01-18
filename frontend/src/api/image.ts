import { GetImage, Image, DeleteImage, StoreImageBook, StoreImageUser, StoreImageAuthor } from "@shared/types/image";
import { handleApiError } from "./utils/handleApiError";

export const imageApi = {
  PREFIX: "/api/v1/image",

  getImage(body: GetImage): Promise<Image> {
    return fetch(this.PREFIX, {
      method: "POST",
      credentials: "include",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify(body),
    }).then((res) => {
      if (!res.ok) {
        return handleApiError(res, () => this.getImage(body));
      }
      return res.json();
    });
  },

  deleteImage(body: DeleteImage): Promise<unknown> {
    return fetch(this.PREFIX, {
      method: "DELETE",
      credentials: "include",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify(body),
    }).then((res) => {
      if (!res.ok) {
        return handleApiError(res, () => this.deleteImage(body));
      }
      return res.json();
    });
  },

  storeBookImage(body: StoreImageBook): Promise<unknown> {
    return fetch(`${this.PREFIX}/book`, {
      method: "POST",
      credentials: "include",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify(body),
    }).then((res) => {
      if (!res.ok) {
        return handleApiError(res, () => this.storeBookImage(body));
      }
      return res.json();
    });
  },

  storeUserImage(body: StoreImageUser): Promise<unknown> {
    return fetch(`${this.PREFIX}/user`, {
      method: "POST",
      credentials: "include",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify(body),
    }).then((res) => {
      if (!res.ok) {
        return handleApiError(res, () => this.storeUserImage(body));
      }
      return res.json();
    });
  },

  storeAuthorImage(body: StoreImageAuthor): Promise<unknown> {
    return fetch(`${this.PREFIX}/author`, {
      method: "POST",
      credentials: "include",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify(body),
    }).then((res) => {
      if (!res.ok) {
        return handleApiError(res, () => this.storeAuthorImage(body));
      }
      return res.json();
    });
  },
} as const;
