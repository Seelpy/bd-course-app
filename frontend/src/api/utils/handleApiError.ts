import { AppRoute } from "@shared/constants/routes";
import { authApi } from "../auth";

class SelfThrownError extends Error {}

type MessageError = {
  message: string;
};

type CustomError = {
  error: string;
};

type SystemError = {
  errors: Record<string, string[]>;
};

type TitledError = {
  title: string;
};

const parseError = (res: Response) => {
  return res
    .json()
    .then((json: MessageError | CustomError | SystemError | TitledError) => {
      if ("message" in json) {
        return Promise.reject(new SelfThrownError(json.message));
      } else if ("error" in json) {
        return Promise.reject(new SelfThrownError(json.error));
      } else if ("errors" in json) {
        const parsedErrors = Object.entries(json.errors)
          .map(([key, messages]) => `${key}: ${messages.join(", ")}`)
          .join("\n");
        return Promise.reject(new SelfThrownError(parsedErrors));
      } else if ("title" in json) {
        return Promise.reject(new SelfThrownError(json.title));
      }
      return Promise.reject(new SelfThrownError(JSON.stringify(json)));
    })
    .catch((error) => {
      if (error instanceof SelfThrownError) {
        throw error;
      }
      throw new Error(res.statusText);
    });
};

export const handleApiError = (res: Response, retryRequest?: () => Promise<unknown>): Promise<unknown> => {
  if (res.status === 401 && retryRequest) {
    return authApi
      .refreshToken()
      .then(() => {
        return retryRequest();
      })
      .catch(() => {
        if (res.status === 401 && window.location.pathname !== AppRoute.Root.toString()) {
          window.location.href = AppRoute.Root.toString();
        }
        return parseError(res);
      });
  }

  return parseError(res);
};
