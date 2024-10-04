export type CreateUserInput = {
  email: string;
  emailVerified: boolean;
  firstName?: string;
  lastName?: string;
  password?: string;
  photoUrl: string;
  username: string;
};
