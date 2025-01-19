export enum UserRole {
  Admin = 0,
  Client = 1,
}

export type User = {
  id: string;
  login: string;
  role: UserRole;
  aboutMe: string;
  avatar: string;
};

export type CreateUser = {
  login: string;
  password: string;
  aboutMe: string;
};

export type EditUser = {
  id: string;
  login: string;
  password: string;
  aboutMe: string;
};

export type DeleteUser = {
  id: string;
};
