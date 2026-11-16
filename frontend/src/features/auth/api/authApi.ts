import { ApiError, apiRequest } from '../../../shared/api/httpClient';

/** The signed-in User. */
export interface User {
  id: string;
  username: string;
  name: string;
  isAdministrator: boolean;
  /** A Temporary Password is pending and must be changed before anything else. */
  mustChangePassword: boolean;
}

interface UserDto {
  id: string;
  username: string;
  name: string;
  is_administrator: boolean;
  must_change_password: boolean;
}

export interface Credentials {
  username: string;
  password: string;
}

export interface PasswordChange {
  currentPassword: string;
  newPassword: string;
}

function toUser(dto: UserDto): User {
  return {
    id: dto.id,
    username: dto.username,
    name: dto.name,
    isAdministrator: dto.is_administrator,
    mustChangePassword: dto.must_change_password,
  };
}

export async function signIn(credentials: Credentials): Promise<User> {
  const dto = await apiRequest<UserDto>('/api/v1/sessions', { method: 'POST', body: credentials });
  return toUser(dto);
}

export async function signOut(): Promise<void> {
  await apiRequest<undefined>('/api/v1/sessions/current', { method: 'DELETE' });
}

/** Returns the signed-in User, or null when there is no valid session. */
export async function getCurrentUser(signal?: AbortSignal): Promise<User | null> {
  try {
    return toUser(await apiRequest<UserDto>('/api/v1/users/me', { signal }));
  } catch (error) {
    if (error instanceof ApiError && error.status === 401) {
      return null;
    }
    throw error;
  }
}

export async function changePassword(change: PasswordChange): Promise<void> {
  await apiRequest<undefined>('/api/v1/users/me/password', {
    method: 'PUT',
    body: { current_password: change.currentPassword, new_password: change.newPassword },
  });
}
