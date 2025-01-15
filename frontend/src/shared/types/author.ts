export type CreateAuthor = {
  firstName: string;
  secondName: string;
  middleName?: string;
  nickName?: string;
};

export type EditAuthor = {
  id: string;
  firstName: string;
  secondName: string;
  middleName?: string;
  nickName?: string;
};

export type DeleteAuthor = {
  id: string;
};

export type Author = {
  id: string;
  firstName: string;
  secondName: string;
  middleName?: string;
  nickname?: string;
};

export type ListAuthorResponse = {
  authors: Author[];
};
