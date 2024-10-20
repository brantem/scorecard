import { useLoaderData } from 'react-router-dom';

type User = {
  id: number;
  name: string;
};

function Users() {
  const data = useLoaderData() as { users: User[] };

  return (
    <>
      <div className="mx-4 mt-4 flex items-center justify-between">
        <h1 className="font-semibold">Users</h1>
        <button className="rounded-lg bg-neutral-900 px-4 py-2 text-sm font-medium text-white hover:bg-neutral-800">
          Add User
        </button>
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
    </>
  );
}

Users.loader = async () => {
  const res = await fetch(`${import.meta.env.VITE_API_URL}/v1/users`);
  return { users: (await res.json()).nodes };
};

export default Users;
