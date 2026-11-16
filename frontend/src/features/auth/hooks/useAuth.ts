import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { changePassword, getCurrentUser, signIn, signOut, type User } from '../api/authApi';

export const currentUserKey = ['current-user'] as const;

/** The signed-in User: null when signed out, undefined while loading. */
export function useCurrentUser() {
  return useQuery({
    queryKey: currentUserKey,
    queryFn: ({ signal }) => getCurrentUser(signal),
  });
}

export function useSignIn() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: signIn,
    onSuccess: (user) => {
      queryClient.setQueryData(currentUserKey, user);
    },
  });
}

export function useSignOut() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: signOut,
    onSuccess: () => {
      queryClient.setQueryData(currentUserKey, null);
      // Drop every other cached response so nothing of the previous User lingers.
      queryClient.removeQueries({ predicate: (query) => query.queryKey[0] !== currentUserKey[0] });
    },
  });
}

export function useChangePassword() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: changePassword,
    onSuccess: () => {
      queryClient.setQueryData<User | null>(currentUserKey, (user) =>
        user ? { ...user, mustChangePassword: false } : user,
      );
    },
  });
}
