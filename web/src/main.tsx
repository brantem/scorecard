import { createRoot } from 'react-dom/client';
import { RouterProvider, createBrowserRouter } from 'react-router-dom';

import Root from 'routes/Root';
import Users from 'routes/Users';
import Syllabuses from 'routes/Syllabuses';
import Structures from 'routes/Structures';

import './index.css';

const router = createBrowserRouter([
  {
    path: '/',
    element: <Root />,
    children: [
      {
        path: 'users',
        loader: Users.loader,
        action: Users.action,
        element: <Users />,
      },
      {
        path: 'syllabuses',
        loader: Syllabuses.loader,
        action: Syllabuses.action,
        element: <Syllabuses />,
      },
      {
        path: 'structures',
        loader: Structures.loader,
        element: <Structures />,
      },
    ],
  },
]);

createRoot(document.getElementById('root')!).render(<RouterProvider router={router} />);
