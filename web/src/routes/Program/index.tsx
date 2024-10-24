import {
  Link,
  Outlet,
  useLocation,
  useParams,
  useLoaderData,
  redirect,
  type LoaderFunctionArgs,
} from 'react-router-dom';
import { ArrowLeftIcon } from '@heroicons/react/20/solid';

import { cn } from 'lib/helpers';

import type * as types from 'types/program';

const getRoutes = (programId: string) => [
  { to: `/${programId}`, text: 'Scorecards', exact: true },
  { to: `/${programId}/users`, text: 'Users' },
  { to: `/${programId}/syllabuses`, text: 'Syllabuses' },
  { to: `/${programId}/structures`, text: 'Structures' },
];

function Program() {
  const params = useParams();
  const location = useLocation();
  const data = useLoaderData() as { program: types.Program };

  const routes = getRoutes(params.programId!);

  return (
    <>
      <div className="grid grid-cols-3 items-center border-b border-neutral-200 bg-white text-sm font-medium">
        <div className="flex h-full items-center gap-2 p-2">
          <Link
            to="/"
            className="flex aspect-square h-full items-center justify-center rounded-lg bg-neutral-50 hover:bg-neutral-100"
          >
            <ArrowLeftIcon className="size-5" />
          </Link>
          <h1>{data.program.title}</h1>
        </div>

        <ul className="flex justify-center gap-2 p-3">
          {routes.map((route) => {
            const pathname = location.pathname.replace(/\/$/, ''); // remove trailing slash
            const isActive = route.exact ? pathname === route.to : pathname.includes(route.to);
            return (
              <li key={route.to}>
                <Link
                  to={route.to}
                  className={cn(
                    'inline-block w-full rounded-full px-3 py-1',
                    isActive ? 'bg-neutral-900 text-white' : 'hover:bg-neutral-200',
                  )}
                >
                  {route.text}
                </Link>
              </li>
            );
          })}
        </ul>
      </div>

      <div className="relative mx-auto size-full max-w-[1920px] overflow-hidden border-neutral-200 bg-white min-[1920px]:border-x">
        <Outlet />
      </div>
    </>
  );
}

Program.loader = async ({ params }: LoaderFunctionArgs) => {
  try {
    const res = await fetch(`${import.meta.env.VITE_API_URL}/v1/programs/${params.programId}`);
    const program = (await res.json()).program;
    if (!program) return redirect('/');
    return { program };
  } catch {
    return redirect('/');
  }
};

export default Program;
