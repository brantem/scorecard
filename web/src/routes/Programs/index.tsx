import { useRef, useEffect } from 'react';
import { Link, useLoaderData, useFetcher, type ActionFunctionArgs } from 'react-router-dom';
import { PlusIcon } from '@heroicons/react/20/solid';

import Table from 'components/Table';
import Button from 'components/Button';
import SaveModal, { SaveModalHandle } from './SaveModal';
import DeleteModal, { type DeleteModalHandle } from 'components/DeleteModal';

import { Program } from 'types/program';

function Programs() {
  const data = useLoaderData() as { programs: Program[] };
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
      <div className="flex w-full items-center justify-center border-b border-neutral-200 bg-white p-2 py-3 text-sm font-medium">
        <span className="inline-block rounded-full bg-neutral-900 px-3 py-1 text-white">Programs</span>
      </div>

      <div className="mx-auto size-full max-w-[1920px] overflow-hidden border-neutral-200 bg-white min-[1920px]:border-x">
        <div className="flex items-start justify-between p-4 pb-0">
          <div className="flex flex-col gap-1">
            <h2 className="font-semibold">Programs</h2>
            <span className="inline-block text-sm text-neutral-500">{data.programs.length} Programs</span>
          </div>

          <Button className="pl-2.5 text-sm" onClick={() => saveModalRef.current?.open(null)}>
            <PlusIcon className="size-5" />
            <span>Add Program</span>
          </Button>
        </div>

        <div className="m-4 overflow-x-auto rounded-lg border border-neutral-200">
          <Table>
            <thead>
              <tr>
                <Table.Th>Title</Table.Th>
                <Table.Th className="w-24 [&>div]:justify-end">Actions</Table.Th>
              </tr>
            </thead>

            <tbody>
              {data.programs.length ? (
                data.programs.map((node) => (
                  <tr key={node.id}>
                    <Table.Td>
                      <Link to={`/${node.id}`} className="w-full hover:underline">
                        {node.title}
                      </Link>
                    </Table.Td>

                    <Table.Td className="text-sm [&>div]:justify-end [&>div]:gap-1.5 [&>div]:pr-1.5">
                      <button
                        className="flex h-8 items-center rounded-lg border border-neutral-200 bg-neutral-50 px-3 hover:bg-neutral-100"
                        onClick={() => saveModalRef.current?.open(node)}
                      >
                        Edit
                      </button>

                      <button
                        className="flex h-8 items-center rounded-lg border border-red-200 bg-red-50 px-3 text-red-500 hover:bg-red-100"
                        onClick={() => deleteModalRef.current?.open('Program', { type: 'DELETE', _programId: node.id })}
                      >
                        Delete
                      </button>
                    </Table.Td>
                  </tr>
                ))
              ) : (
                <tr>
                  <Table.Td className="text-neutral-400 [&>div]:justify-center" colSpan={2}>
                    No programs available
                  </Table.Td>
                </tr>
              )}
            </tbody>
          </Table>
        </div>
      </div>

      <SaveModal ref={saveModalRef} />
      <DeleteModal
        ref={deleteModalRef}
        onAccept={(body) => fetcher.submit(body, { method: 'DELETE', encType: 'application/json' })}
      />
    </>
  );
}

Programs.loader = async () => {
  try {
    const res = await fetch(`${import.meta.env.VITE_API_URL}/v1/programs`);
    return { programs: (await res.json()).nodes };
  } catch {
    return { programs: [] };
  }
};

Programs.action = async ({ request }: ActionFunctionArgs) => {
  try {
    const { type, _programId, ...body } = await request.json();

    let res;
    switch (type) {
      case 'SAVE':
        res = await fetch(`${import.meta.env.VITE_API_URL}/v1/programs/${_programId || ''}`, {
          method: 'PUT',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify(body),
        });
        break;
      case 'DELETE':
        res = await fetch(`${import.meta.env.VITE_API_URL}/v1/programs/${_programId}`, {
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

export default Programs;
