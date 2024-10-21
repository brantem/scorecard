import { Link, Outlet, useLocation } from 'react-router-dom';

import { cn } from 'lib/helpers';

const routes = [
  { to: '/users', text: 'Users' },
  { to: '/syllabuses', text: 'Syllabuses' },
  { to: '/structures', text: 'Structures' },
];

function Root() {
  const location = useLocation();

  return (
    <>
      <ul className="flex w-full items-center justify-center gap-2 border-b border-neutral-200 bg-white p-2 py-3 text-sm font-medium">
        {routes.map((route) => (
          <li key={route.to}>
            <Link
              to={route.to}
              className={cn(
                'inline-block w-full rounded-full px-3 py-1',
                location.pathname.includes(route.to) ? 'bg-neutral-900 text-white' : 'hover:bg-neutral-200',
              )}
            >
              {route.text}
            </Link>
          </li>
        ))}
      </ul>

      <div className="mx-auto h-full w-full max-w-[1920px] overflow-hidden border-neutral-200 bg-white min-[1920px]:border-x">
        <Outlet />
      </div>
    </>
  );
}

export default Root;
