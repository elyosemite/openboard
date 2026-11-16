import { Navigate, type RouteObject } from 'react-router';
import {
  CHANGE_PASSWORD_PATH,
  ChangePasswordPage,
  LoginPage,
  RequireAuth,
  SIGN_IN_PATH,
} from '../features/auth';
import { AppLayout } from './AppLayout';
import { HomePage } from './HomePage';

export const routes: RouteObject[] = [
  { path: SIGN_IN_PATH, element: <LoginPage /> },
  {
    element: <RequireAuth />,
    children: [
      {
        element: <AppLayout />,
        children: [
          { path: '/', element: <HomePage /> },
          { path: CHANGE_PASSWORD_PATH, element: <ChangePasswordPage /> },
        ],
      },
    ],
  },
  { path: '*', element: <Navigate to="/" replace /> },
];
