import { createRoot } from 'react-dom/client';
import { RouterProvider, createBrowserRouter } from 'react-router-dom';

import Root from 'routes/Root';
import Users from 'routes/Users';
import Syllabuses from 'routes/Syllabuses';
import Scores from 'routes/Scores';
import Scorecards from 'routes/Scorecards';

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
        path: 'syllabuses/:syllabusId/scores',
        loader: Scores.loader,
        action: Scores.action,
        element: <Scores />,
      },
      {
        path: 'scorecards',
        loader: Scorecards.loader,
        action: Scorecards.action,
        element: <Scorecards />,
      },
    ],
  },
]);

createRoot(document.getElementById('root')!).render(<RouterProvider router={router} />);
