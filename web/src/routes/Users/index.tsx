import { useRef, useEffect } from 'react';
import {
  Link,
  useParams,
  useLoaderData,
  useFetcher,
  type LoaderFunctionArgs,
  type ActionFunctionArgs,
} from 'react-router-dom';
import { PlusIcon } from '@heroicons/react/20/solid';

import Button from 'components/Button';
import Table from 'components/Table';
import SaveModal, { type SaveModalHandle } from './SaveModal';
import DeleteModal, { type DeleteModalHandle } from 'components/DeleteModal';

import type { User } from 'types/user';

function Users() {
  const params = useParams();
  const data = useLoaderData() as { users: User[] };
  const fetcher = useFetcher<{ success: boolean; error: { code: string | null } | null }>();

  const saveModalRef = useRef<SaveModalHandle>(null);
  const deleteModalRef = useRef<DeleteModalHandle>(null);

  useEffect(() => {
    if (!fetcher.data) return;
    if (fetcher.data.success) {
      switch ((fetcher.json as { type?: string }).type) {
        case 'DELETE':
          deleteModalRef.current?.close();
          break;
      }
    } else {
      alert(fetcher.data.error?.code || 'INTERNAL_SERVER_ERROR');
    }
  }, [fetcher.data]);

  return (
    <>
      <div className="flex items-start justify-between p-4 pb-0">
        <div className="flex flex-col gap-1">
          <h2 className="font-semibold">Users</h2>
          <span className="inline-block text-sm text-neutral-500">{data.users.length} Users</span>
        </div>

        <Button className="pl-2.5 text-sm" onClick={() => saveModalRef.current?.open(null)}>
          <PlusIcon className="size-5" />
          <span>Add User</span>
        </Button>
      </div>

      <div className="m-4 overflow-x-auto rounded-lg border border-neutral-200">
        <Table>
          <thead>
            <tr>
              <Table.Th>Name</Table.Th>

              <Table.Th className="w-32 [&>div]:justify-end">Actions</Table.Th>
            </tr>
          </thead>

          <tbody>
            {data.users.length ? (
              <>
                {data.users.map((user) => (
                  <tr key={user.id}>
                    <Table.Td>{user.name}</Table.Td>

                    <Table.Td className="text-sm [&>div]:justify-end [&>div]:gap-1.5 [&>div]:pr-1.5">
                      <Link
                        to={`/${params.programId}/users/${user.id}/scores`}
                        className="mr-2 flex h-8 items-center rounded-lg border border-neutral-200 bg-neutral-50 px-3 hover:bg-neutral-100"
                      >
                        Scores
                      </Link>

                      <button
                        className="flex h-8 items-center rounded-lg border border-neutral-200 bg-neutral-50 px-3 hover:bg-neutral-100"
                        onClick={() => saveModalRef.current?.open(user)}
                      >
                        Edit
                      </button>

                      <button
                        className="flex h-8 items-center rounded-lg border border-red-200 bg-red-50 px-3 text-red-500 hover:bg-red-100"
                        onClick={() => deleteModalRef.current?.open('User', { type: 'DELETE', _userId: user.id })}
                      >
                        Delete
                      </button>
                    </Table.Td>
                  </tr>
                ))}

                {[...new Array(10 - data.users.length)].map((_, i) => (
                  <tr key={i}>
                    <td className="h-12" colSpan={2} />
                  </tr>
                ))}
              </>
            ) : (
              <tr>
                <Table.Td className="text-neutral-400 [&>div]:justify-center" colSpan={2}>
                  No users available
                </Table.Td>
              </tr>
            )}
          </tbody>
        </Table>
      </div>

      <SaveModal ref={saveModalRef} />
      <DeleteModal
        ref={deleteModalRef}
        onAccept={(body) => fetcher.submit(body, { method: 'DELETE', encType: 'application/json' })}
      />
    </>
  );
}

Users.loader = async ({ params }: LoaderFunctionArgs) => {
  try {
    const res = await fetch(`${import.meta.env.VITE_API_URL}/v1/programs/${params.programId}/users`);
    return { users: (await res.json()).nodes };
  } catch {
    return { users: [] };
  }
};

Users.action = async ({ request, params }: ActionFunctionArgs) => {
  try {
    const { type, _userId, ...body } = await request.json();

    let res;
    switch (type) {
      case 'SAVE':
        res = await fetch(`${import.meta.env.VITE_API_URL}/v1/programs/${params.programId}/users/${_userId || ''}`, {
          method: 'PUT',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify(body),
        });
        break;
      case 'DELETE':
        res = await fetch(`${import.meta.env.VITE_API_URL}/v1/programs/${params.programId}/users/${_userId || ''}`, {
          method: 'DELETE',
        });
        break;
      default:
        return { success: false, error: null };
    }
    return await res.json();
  } catch {
    return { success: false, error: { code: 'INTERNAL_SERVER_ERROR' } };
  }
};

export default Users;
