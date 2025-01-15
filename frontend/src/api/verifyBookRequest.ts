import {
  CreateVerifyBook,
  DeleteVerifyBook,
  AcceptBookChapter,
  ListVerifyBookRequestResponse,
} from "@shared/types/verifyBookRequest";
import { handleApiError } from "./utils/handleApiError";

export const verifyBookRequestApi = {
  PREFIX: "/api/v1/verify-book-request",

  listVerifyBookRequests(): Promise<ListVerifyBookRequestResponse> {
    return fetch(this.PREFIX, {
      method: "GET",
      credentials: "include",
    }).then((res) => {
      if (!res.ok) {
        return handleApiError(res, () => this.listVerifyBookRequests());
      }
      return res.json();
    });
  },

  createVerifyBookRequest(body: CreateVerifyBook): Promise<unknown> {
    return fetch(this.PREFIX, {
      method: "POST",
      credentials: "include",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify(body),
    }).then((res) => {
      if (!res.ok) {
        return handleApiError(res, () => this.createVerifyBookRequest(body));
      }
      return res.json();
    });
  },

  deleteVerifyBookRequest(body: DeleteVerifyBook): Promise<unknown> {
    return fetch(this.PREFIX, {
      method: "DELETE",
      credentials: "include",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify(body),
    }).then((res) => {
      if (!res.ok) {
        return handleApiError(res, () => this.deleteVerifyBookRequest(body));
      }
      return res.json();
    });
  },

  acceptVerifyBookRequest(body: AcceptBookChapter): Promise<unknown> {
    return fetch(`${this.PREFIX}/accept`, {
      method: "POST",
      credentials: "include",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify(body),
    }).then((res) => {
      if (!res.ok) {
        return handleApiError(res, () => this.acceptVerifyBookRequest(body));
      }
      return res.json();
    });
  },
} as const;
