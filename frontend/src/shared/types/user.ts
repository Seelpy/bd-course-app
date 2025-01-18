export enum UserRole {
  Admin = 0,
  Client = 1,
}

export type User = {
  id: string;
  login: string;
  role: UserRole;
  aboutMe: string;
  avatarId: string;
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
