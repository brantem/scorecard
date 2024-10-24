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
import { TrashIcon } from '@heroicons/react/16/solid';

import Button from 'components/Button';
import Tree from 'components/Tree';
import ResetModal, { type ResetModalHandle } from 'components/ResetModal';
import AddModal, { type AddModalHandle } from './AddModal';
import EditModal, { type EditModalHandle } from './EditModal';
import DeleteModal, { type DeleteModalHandle } from 'components/DeleteModal';

import { cn } from 'lib/helpers';
import type { Structure } from 'types/scorecard';

function Structures() {
  const params = useParams();
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
          resetModalRef.current?.close();
          break;
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
      <div className="relative h-full">
        <div className="absolute left-0 right-0 top-0 flex items-start justify-between gap-4 p-4 pb-0">
          <div className="flex flex-col gap-1">
            <h2 className="font-semibold">Structures</h2>
            <span className="inline-block text-sm text-neutral-500">
              Defines how scores are gathered into scorecards
            </span>
          </div>

          {data.structures.length ? (
            <Button
              className="bg-red-50 pl-2.5 text-sm text-red-500 hover:bg-red-100"
              onClick={() => resetModalRef.current?.open('Structures', { type: 'RESET', _structureId: '0' })}
            >
              <TrashIcon className="size-4" />
              <span>Reset</span>
            </Button>
          ) : null}
        </div>

        {data.canCreate ? (
          data.structures.length ? (
            <div className="no-scrollbar flex size-full flex-col items-center overflow-y-auto p-4 pt-16">
              <Tree
                items={structures}
                renderOptions={(structure) => (
                  <>
                    <Tree.Item className="min-w-0 hover:bg-neutral-100" asChild>
                      <button onClick={() => editModalRef.current?.open(structure)}>Edit</button>
                    </Tree.Item>

                    <Tree.Item className="min-w-0 border-red-200 bg-red-50 text-red-500 hover:bg-red-100" asChild>
                      <button
                        onClick={() => {
                          deleteModalRef.current?.open('Structure', { type: 'DELETE', _structureId: structure.id });
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
                      <button onClick={() => addModalRef.current?.open(parent)}>
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
                <Button className="text-sm" onClick={() => addModalRef.current?.open(null)}>
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
              <Link to={`/${params.programId}/syllabuses`}>Syllabuses</Link>
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

Structures.loader = async ({ params }: LoaderFunctionArgs) => {
  const [canCreate, structures] = await Promise.all([
    (async () => {
      try {
        const url = `${import.meta.env.VITE_API_URL}/v1/programs/${params.programId}/syllabuses`;
        const res = await fetch(url, { method: 'HEAD' });
        return parseInt(res.headers.get('X-Total-Count') || '0') > 0;
      } catch {
        return false;
      }
    })(),
    (async () => {
      try {
        const url = `${import.meta.env.VITE_API_URL}/v1/programs/${params.programId}/scorecards/structures`;
        const res = await fetch(url);
        return (await res.json()).nodes;
      } catch {
        return [];
      }
    })(),
  ]);
  return { canCreate, structures };
};

Structures.action = async ({ request, params }: ActionFunctionArgs) => {
  try {
    const { type, _syllabusId, _structureId, ...body } = await request.json();

    let url, res;
    switch (type) {
      case 'COPY':
        url = `${import.meta.env.VITE_API_URL}/v1/programs/${params.programId}/scorecards/structures/copy/${_syllabusId}`;
        res = await fetch(url, {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify(body),
        });
        break;
      case 'SAVE':
        url = `${import.meta.env.VITE_API_URL}/v1/programs/${params.programId}/scorecards/structures/${_structureId || ''}`;
        res = await fetch(url, {
          method: 'PUT',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify(body),
        });
        break;
      case 'RESET':
      case 'DELETE':
        url = `${import.meta.env.VITE_API_URL}/v1/programs/${params.programId}/scorecards/structures/${_structureId || ''}`;
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

export default Structures;
