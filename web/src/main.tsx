import { createRoot } from 'react-dom/client';
import { RouterProvider, createBrowserRouter } from 'react-router-dom';

import Root from 'routes/Root';
import Users from 'routes/Users';
import UserScores from 'routes/Users/Scores';
import Scorecards from 'routes/Scorecards';
import Scorecard from 'routes/Scorecard';
import Syllabuses from 'routes/Syllabuses';
import SyllabusScores from 'routes/Syllabuses/Scores';
import Structures from 'routes/Structures';

import './index.css';

const router = createBrowserRouter([
  {
    path: '/',
    element: <Root />,
    children: [
      {
        path: '/',
        loader: Scorecards.loader,
        action: Scorecards.action,
        element: <Scorecards />,
      },
      {
        path: '/:scorecardId',
        loader: Scorecard.loader,
        action: Scorecard.action,
        element: <Scorecard />,
      },
      {
        path: '/users',
        loader: Users.loader,
        action: Users.action,
        element: <Users />,
      },
      {
        path: '/users/:userId/scores',
        loader: UserScores.loader,
        action: UserScores.action,
        element: <UserScores />,
      },
      {
        path: '/syllabuses',
        loader: Syllabuses.loader,
        action: Syllabuses.action,
        element: <Syllabuses />,
      },
      {
        path: '/syllabuses/:syllabusId/scores',
        loader: SyllabusScores.loader,
        action: SyllabusScores.action,
        element: <SyllabusScores />,
      },
      {
        path: '/structures',
        loader: Structures.loader,
        action: Structures.action,
        element: <Structures />,
      },
    ],
  },
]);

createRoot(document.getElementById('root')!).render(<RouterProvider router={router} />);
