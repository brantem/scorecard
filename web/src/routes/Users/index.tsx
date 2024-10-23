import { useRef, useEffect } from 'react';
import {
  Link,
  useParams,
  useSearchParams,
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

const LIMIT = 10;

function Users() {
  const params = useParams();
  const [searchParams, setSearchParams] = useSearchParams();
  const data = useLoaderData() as { users: { totalCount: number; nodes: User[] } };
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
          <span className="inline-block text-sm text-neutral-500">{data.users.totalCount} Users</span>
        </div>

        <Button className="pl-2.5 text-sm" onClick={() => saveModalRef.current?.open(null)}>
          <PlusIcon className="size-5" />
          <span>Add User</span>
        </Button>
      </div>

      <div className="p-4">
        <Table.Provider
          totalItems={data.users.totalCount}
          onStateChange={(state) => {
            setSearchParams((prev) => {
              prev.set('page', state.page.toString());
              return prev;
            });
          }}
          initialState={{ itemsPerPage: LIMIT, page: parseInt(searchParams.get('page') || '1') }}
        >
          <div className="overflow-x-auto rounded-lg border border-neutral-200">
            <Table>
              <thead>
                <tr>
                  <Table.Th>Name</Table.Th>

                  <Table.Th className="w-32 [&>div]:justify-end">Actions</Table.Th>
                </tr>
              </thead>

              <tbody>
                {data.users.nodes.length ? (
                  <>
                    {data.users.nodes.map((node) => (
                      <tr key={node.id}>
                        <Table.Td>{node.name}</Table.Td>

                        <Table.Td className="text-sm [&>div]:justify-end [&>div]:gap-1.5 [&>div]:pr-1.5">
                          <Link
                            to={`/${params.programId}/users/${node.id}/scores`}
                            className="mr-2 flex h-8 items-center rounded-lg border border-neutral-200 bg-neutral-50 px-3 hover:bg-neutral-100"
                          >
                            Scores
                          </Link>

                          <button
                            className="flex h-8 items-center rounded-lg border border-neutral-200 bg-neutral-50 px-3 hover:bg-neutral-100"
                            onClick={() => saveModalRef.current?.open(node)}
                          >
                            Edit
                          </button>

                          <button
                            className="flex h-8 items-center rounded-lg border border-red-200 bg-red-50 px-3 text-red-500 hover:bg-red-100"
                            onClick={() => deleteModalRef.current?.open('User', { type: 'DELETE', _userId: node.id })}
                          >
                            Delete
                          </button>
                        </Table.Td>
                      </tr>
                    ))}

                    {[...new Array(LIMIT - data.users.nodes.length)].map((_, i) => (
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

          <div className="mt-4 flex items-center justify-end text-sm">
            <Table.Pagination />
          </div>
        </Table.Provider>
      </div>

      <SaveModal ref={saveModalRef} />
      <DeleteModal
        ref={deleteModalRef}
        onAccept={(body) => fetcher.submit(body, { method: 'DELETE', encType: 'application/json' })}
      />
    </>
  );
}

Users.loader = async ({ request, params }: LoaderFunctionArgs) => {
  try {
    const page = parseInt(new URL(request.url).searchParams.get('page') || '1') - 1;
    const url = `${import.meta.env.VITE_API_URL}/v1/programs/${params.programId}/users?limit=${LIMIT}&offset=${page * LIMIT}`;
    const res = await fetch(url);
    return {
      users: {
        totalCount: parseInt(res.headers.get('X-Total-Count') || '0'),
        nodes: (await res.json()).nodes,
      },
    };
  } catch {
    return {
      users: {
        totalCount: 0,
        nodes: [],
      },
    };
  }
};

Users.action = async ({ request, params }: ActionFunctionArgs) => {
  try {
    const { type, _userId, ...body } = await request.json();

    let url, res;
    switch (type) {
      case 'SAVE':
        url = `${import.meta.env.VITE_API_URL}/v1/programs/${params.programId}/users/${_userId || ''}`;
        res = await fetch(url, {
          method: 'PUT',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify(body),
        });
        break;
      case 'DELETE':
        url = `${import.meta.env.VITE_API_URL}/v1/programs/${params.programId}/users/${_userId || ''}`;
        res = await fetch(url, { method: 'DELETE' });
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
