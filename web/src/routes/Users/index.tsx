import { useRef } from 'react';
import { ActionFunctionArgs, useLoaderData } from 'react-router-dom';

import Button from 'components/Button';
import SaveModal, { type SaveModalHandle } from './SaveModal';

type User = {
  id: number;
  name: string;
};

function Users() {
  const data = useLoaderData() as { users: User[] };

  const saveModalRef = useRef<SaveModalHandle>(null);

  return (
    <>
      <div className="mx-4 mt-4 flex items-start justify-between">
        <h1 className="font-semibold">Users</h1>
        <Button className="pl-2.5 text-sm" onClick={() => saveModalRef.current?.onOpen()}>
          <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor" className="size-5">
            <path d="M10.75 4.75a.75.75 0 0 0-1.5 0v4.5h-4.5a.75.75 0 0 0 0 1.5h4.5v4.5a.75.75 0 0 0 1.5 0v-4.5h4.5a.75.75 0 0 0 0-1.5h-4.5v-4.5Z" />
          </svg>
          <span>Add User</span>
        </Button>
      </div>

      <div className="m-4 overflow-x-auto rounded-lg border border-neutral-200">
        <table className="w-full">
          <thead>
            <tr className="bg-neutral-50">
              <th className="h-12 whitespace-nowrap py-0 font-medium">
                <div className="flex h-full w-full items-center justify-between px-3">Name</div>
              </th>
            </tr>
          </thead>

          <tbody>
            {data.users.map((user) => (
              <tr key={user.id} className="border-t border-neutral-200 text-neutral-700">
                <td className="h-12 whitespace-nowrap py-0 text-neutral-700">
                  <div className="flex h-full items-center px-3">{user.name}</div>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>

      <SaveModal ref={saveModalRef} />
    </>
  );
}

Users.loader = async () => {
  try {
    const res = await fetch(`${import.meta.env.VITE_API_URL}/v1/users`);
    return { users: (await res.json()).nodes };
  } catch {
    return { users: [] };
  }
};

Users.action = async ({ request }: ActionFunctionArgs) => {
  try {
    const { _userId, ...body } = await request.json();
    const res = await fetch(`${import.meta.env.VITE_API_URL}/v1/users/${_userId || ''}`, {
      method: 'PUT',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(body),
    });
    return await res.json();
  } catch {
    return { success: false, error: { code: 'INTERNAL_SERVER_ERROR' } };
  }
};

export default Users;
