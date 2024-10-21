import { useRef, useEffect } from 'react';
import { Link, useLoaderData, useFetcher, type ActionFunctionArgs } from 'react-router-dom';
import { PlusIcon } from '@heroicons/react/20/solid';
import { TrashIcon } from '@heroicons/react/16/solid';

import Button from 'components/Button';
import Tree from 'components/Tree';
import ResetModal, { type ResetModalHandle } from 'components/ResetModal';
import AddModal, { type AddModalHandle } from './AddModal';
import EditModal, { type EditModalHandle } from './EditModal';
import DeleteModal, { type DeleteModalHandle } from 'components/DeleteModal';

import { cn } from 'lib/helpers';
import type { Structure } from 'types/scorecard';

function Scorecards() {
  const data = useLoaderData() as { canCreate: boolean; structures: Structure[] };
  const fetcher = useFetcher<{ success: boolean; error: { code: string | null } | null }>();

  const resetModalRef = useRef<ResetModalHandle>(null);
  const addModalRef = useRef<AddModalHandle>(null);
  const editModalRef = useRef<EditModalHandle>(null);
  const deleteModalRef = useRef<DeleteModalHandle>(null);

  const structures = new Map<Structure['parentId'], Structure[]>();
  data.structures.forEach((structure) => {
    structures.set(structure.parentId, [...(structures.get(structure.parentId) || []), structure]);
  });

  useEffect(() => {
    if (!fetcher.data) return;
    if (fetcher.data.success) {
      switch ((fetcher.json as { type?: string }).type) {
        case 'RESET':
          resetModalRef.current?.onClose();
          break;
        case 'DELETE':
          deleteModalRef.current?.onClose();
          break;
      }
    } else {
      alert(fetcher.data.error?.code || 'INTERNAL_SERVER_ERROR');
    }
  }, [fetcher.data]);

  return (
    <>
      <div className="relative h-full">
        <div className="absolute left-4 right-4 top-4 flex items-start justify-between gap-4 font-semibold">
          <h2>Structures</h2>

          {data.structures.length ? (
            <Button
              className="bg-red-50 px-3 py-1.5 pl-2 text-sm text-red-500 hover:bg-red-100"
              onClick={() => resetModalRef.current?.onOpen('Scorecards', { type: 'RESET', _structureId: 'all' })}
            >
              <TrashIcon className="size-4" />
              <span>Reset</span>
            </Button>
          ) : null}
        </div>

        {data.canCreate ? (
          data.structures.length ? (
            <div className="no-scrollbar flex size-full flex-col items-center overflow-y-auto p-4 pt-12">
              <Tree
                items={structures}
                renderOptions={(structure) => (
                  <>
                    <Tree.Item className="min-w-0 hover:bg-neutral-100" asChild>
                      <button onClick={() => editModalRef.current?.onOpen(structure)}>Edit</button>
                    </Tree.Item>

                    <Tree.Item className="min-w-0 border-red-200 bg-red-50 text-red-500 hover:bg-red-100" asChild>
                      <button
                        onClick={() => {
                          deleteModalRef.current?.onOpen('Structure', { type: 'DELETE', _structureId: structure.id });
                        }}
                      >
                        Delete
                      </button>
                    </Tree.Item>
                  </>
                )}
                renderAdd={(parent) => {
                  if (parent?.syllabus?.isAssignment) return null;
                  return (
                    <Tree.Item
                      className={cn(
                        'flex min-w-0 gap-2 border-neutral-800 bg-neutral-900 pl-2 text-white hover:bg-neutral-800',
                        parent && 'ml-[calc(theme(spacing.8)+2px)]',
                      )}
                      asChild
                    >
                      <button onClick={() => addModalRef.current?.onOpen(parent)}>
                        <PlusIcon className="size-5" />
                        <span>Add Structure</span>
                      </button>
                    </Tree.Item>
                  );
                }}
              />
            </div>
          ) : (
            <div className="flex h-full flex-col items-center justify-center gap-2">
              <span>Start building the structures for the Scorecard</span>

              <div className="flex items-center gap-2">
                <Button className="text-sm" onClick={() => addModalRef.current?.onOpen(null)}>
                  Manual
                </Button>
                <span>Or</span>
                <Button
                  className="text-sm"
                  onClick={() => {
                    fetcher.submit({ type: 'COPY', _syllabusId: 0 }, { method: 'POST', encType: 'application/json' });
                  }}
                >
                  Generate from Syllabuses
                </Button>
              </div>
            </div>
          )
        ) : (
          <div className="flex h-full flex-col items-center justify-center gap-2 text-neutral-500">
            <span>You need to create the syllabuses before adding any structures</span>
            <Button className="rounded-full" asChild>
              <Link to="/syllabuses">Syllabuses</Link>
            </Button>
          </div>
        )}
      </div>

      <ResetModal
        ref={resetModalRef}
        onAccept={(body) => fetcher.submit(body, { method: 'DELETE', encType: 'application/json' })}
      />
      <AddModal ref={addModalRef} />
      <EditModal ref={editModalRef} />
      <DeleteModal
        ref={deleteModalRef}
        onAccept={(body) => fetcher.submit(body, { method: 'DELETE', encType: 'application/json' })}
      />
    </>
  );
}

Scorecards.loader = async () => {
  const [canCreate, structures] = await Promise.all([
    (async () => {
      try {
        const res = await fetch(`${import.meta.env.VITE_API_URL}/v1/syllabuses`, { method: 'HEAD' });
        return parseInt(res.headers.get('X-Total-Count') || '0') > 0;
      } catch {
        return [];
      }
    })(),
    (async () => {
      try {
        const res = await fetch(`${import.meta.env.VITE_API_URL}/v1/scorecards/structures`);
        return (await res.json()).nodes;
      } catch {
        return [];
      }
    })(),
  ]);
  return { canCreate, structures };
};

Scorecards.action = async ({ request }: ActionFunctionArgs) => {
  try {
    const { type, _syllabusId, _structureId, ...body } = await request.json();

    let res;
    switch (type) {
      case 'COPY':
        res = await fetch(`${import.meta.env.VITE_API_URL}/v1/scorecards/structures/copy/${_syllabusId}`, {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify(body),
        });
        break;
      case 'SAVE':
        res = await fetch(`${import.meta.env.VITE_API_URL}/v1/scorecards/structures/${_structureId || ''}`, {
          method: 'PUT',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify(body),
        });
        break;
      case 'RESET':
      case 'DELETE':
        res = await fetch(`${import.meta.env.VITE_API_URL}/v1/scorecards/structures/${_structureId || ''}`, {
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

export default Scorecards;
