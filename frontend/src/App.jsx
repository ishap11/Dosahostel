import { createBrowserRouter, RouterProvider } from 'react-router-dom';
import Login from './pages/Login';
import Register from './pages/Register';
// import Dashboard from './pages/Dashboard';
import LoginOTPVerification from './pages/LoginOtpPage';
import PrivateRoute from './context/PrivateRoute';
import MainLayout from './pages/layout/MainLayout';
import LandingPage from './pages/landingPage';

import Inventory from './pages/Inventory';
import { RoutesPathName } from './constants';
// import PrivateRoute from './context/PrivateRoute';


const router = createBrowserRouter([
  {
    path: '/',
    element: <LandingPage />,  // Landing Page should be first
  },
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


  // {
  //   path: RoutesPathName.DASHBOARD_PAGE,
  //   element: (
  //     <PrivateRoute>
  //         <Dashboard />
  //     </PrivateRoute>
  //   ),
  // },
  {
   
    path: '/',
    element: (
      <PrivateRoute>
        {/* <Dashboard > */}
          <MainLayout />
        {/* </Dashboard> */}
      </PrivateRoute>
    ),
    children: [
      {
        path: RoutesPathName.Inventory_page,
        element: <Inventory />,
      },
    ],
  },

]);

function App() {
  return <RouterProvider router={router} />;
}

export default App;