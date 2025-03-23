import { createBrowserRouter, RouterProvider } from 'react-router-dom';
import Login from './pages/Login';
import Register from './pages/Register';
// import Dashboard from './pages/Dashboard';
import LoginOTPVerification from './pages/LoginOtpPage';

import Inventory from './pages/Inventory';
import { RoutesPathName } from './constants';
// import PrivateRoute from './context/PrivateRoute';

const router = createBrowserRouter([
  {
    path: RoutesPathName.SIGNUP_PAGE,
    element: <Register />,
  },
  {
    path: RoutesPathName.LOGIN_PAGE,
    element: <Login />,
  },
  {
    path: RoutesPathName.LoginOTPVerification_Page,
    element: <LoginOTPVerification />,
  },

  {
    path:RoutesPathName.INVENTORY_PAGE,
    element: <Inventory />
  }

  // {
  //   path: RoutesPathName.DASHBOARD_PAGE,
  //   element: (
  //     <PrivateRoute>
  //         <Dashboard />
  //     </PrivateRoute>
  //   ),
  // },
]);

function App() {
  return <RouterProvider router={router} />;
}

export default App;