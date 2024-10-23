import { createRoot } from 'react-dom/client';
import { RouterProvider, createBrowserRouter } from 'react-router-dom';

import Programs from 'routes/Programs';
import Program from 'routes/Program';
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
    children: [
      {
        index: true,
        loader: Programs.loader,
        action: Programs.action,
        element: <Programs />,
      },
      {
        path: '/:programId',
        loader: Program.loader,
        element: <Program />,
        children: [
          {
            index: true,
            loader: Scorecards.loader,
            action: Scorecards.action,
            element: <Scorecards />,
          },
          {
            path: '/:programId/:scorecardId',
            loader: Scorecard.loader,
            action: Scorecard.action,
            element: <Scorecard />,
          },
          {
            path: '/:programId/users',
            loader: Users.loader,
            action: Users.action,
            element: <Users />,
          },
          {
            path: '/:programId/users/:userId/scores',
            loader: UserScores.loader,
            action: UserScores.action,
            element: <UserScores />,
          },
          {
            path: '/:programId/syllabuses',
            loader: Syllabuses.loader,
            action: Syllabuses.action,
            element: <Syllabuses />,
          },
          {
            path: '/:programId/syllabuses/:syllabusId/scores',
            loader: SyllabusScores.loader,
            action: SyllabusScores.action,
            element: <SyllabusScores />,
          },
          {
            path: '/:programId/structures',
            loader: Structures.loader,
            action: Structures.action,
            element: <Structures />,
          },
        ],
      },
    ],
  },
]);

createRoot(document.getElementById('root')!).render(<RouterProvider router={router} />);
