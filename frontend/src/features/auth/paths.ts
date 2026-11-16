export const SIGN_IN_PATH = '/login';
export const CHANGE_PASSWORD_PATH = '/change-password';

/** Router state carried to the sign-in page: where to return afterwards. */
export interface SignInState {
  from?: string;
}
