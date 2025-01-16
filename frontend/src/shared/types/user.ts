export type User = {
  id: string;
  name: string;
  email: string;
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
