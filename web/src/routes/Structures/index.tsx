import { useRef } from 'react';
import { useLoaderData, useFetcher, type ActionFunctionArgs } from 'react-router-dom';

import Button from 'components/Button';
import Tree from 'components/Tree';
import AddStructureModal, { type AddStructureModalHandle } from './AddStructureModal';

import { cn } from 'lib/helpers';
import type { Structure } from './types';

function Structures() {
  const data = useLoaderData() as { structures: Structure[] };
  const fetcher = useFetcher();

  const addStructureModalRef = useRef<AddStructureModalHandle>(null);

  const structures = new Map();
  data.structures.forEach((structure) => {
    structures.set(structure.parentId, [...(structures.get(structure.parentId) || []), structure]);
  });

  return (
    <>
      <div className="relative flex h-full flex-col items-center overflow-y-auto p-4 pt-16">
        <div className="absolute left-4 right-4 top-4 flex items-start justify-between gap-4 font-semibold">
          <h2>Structures</h2>

          {data.structures.length ? (
            <Button
              className="-mr-1.5 -mt-1.5 bg-red-50 px-3 py-1.5 pl-2 text-sm text-red-400 hover:bg-red-100 hover:text-red-500"
              onClick={() => fetcher.submit({ type: 'RESET' }, { method: 'DELETE', encType: 'application/json' })}
            >
              <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 16 16" fill="currentColor" className="size-4">
                <path
                  fillRule="evenodd"
                  d="M5 3.25V4H2.75a.75.75 0 0 0 0 1.5h.3l.815 8.15A1.5 1.5 0 0 0 5.357 15h5.285a1.5 1.5 0 0 0 1.493-1.35l.815-8.15h.3a.75.75 0 0 0 0-1.5H11v-.75A2.25 2.25 0 0 0 8.75 1h-1.5A2.25 2.25 0 0 0 5 3.25Zm2.25-.75a.75.75 0 0 0-.75.75V4h3v-.75a.75.75 0 0 0-.75-.75h-1.5ZM6.05 6a.75.75 0 0 1 .787.713l.275 5.5a.75.75 0 0 1-1.498.075l-.275-5.5A.75.75 0 0 1 6.05 6Zm3.9 0a.75.75 0 0 1 .712.787l-.275 5.5a.75.75 0 0 1-1.498-.075l.275-5.5a.75.75 0 0 1 .786-.711Z"
                  clipRule="evenodd"
                />
              </svg>

              <span>Reset</span>
            </Button>
          ) : null}
        </div>

        {data.structures.length ? (
          <Tree
            items={structures}
            renderAdd={(parent) => (
              <Tree.Item
                className={cn(
                  'flex min-w-0 gap-2 border-neutral-800 bg-neutral-900 pl-2 text-white hover:bg-neutral-800',
                  parent && 'ml-[26px]',
                )}
                asChild
              >
                <button onClick={() => addStructureModalRef.current?.onOpen(parent)}>
                  <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor" className="size-5">
                    <path d="M10.75 4.75a.75.75 0 0 0-1.5 0v4.5h-4.5a.75.75 0 0 0 0 1.5h4.5v4.5a.75.75 0 0 0 1.5 0v-4.5h4.5a.75.75 0 0 0 0-1.5h-4.5v-4.5Z" />
                  </svg>
                  <span>Add Structure</span>
                </button>
              </Tree.Item>
            )}
          />
        ) : (
          <div className="flex h-full flex-col items-center justify-center gap-2">
            <span>Start building the structures for the Scorecard</span>

            <div className="flex items-center gap-2">
              <Button className="text-sm" onClick={() => addStructureModalRef.current?.onOpen(null)}>
                Manual
              </Button>
              <span>Or</span>
              <Button
                className="text-sm"
                onClick={() => fetcher.submit({ type: 'GENERATE' }, { method: 'POST', encType: 'application/json' })}
              >
                Generate from Syllabuses
              </Button>
            </div>
          </div>
        )}
      </div>

      <AddStructureModal ref={addStructureModalRef} />
    </>
  );
}

Structures.loader = async () => {
  try {
    const res = await fetch(`${import.meta.env.VITE_API_URL}/v1/structures`);
    return { structures: (await res.json()).nodes };
  } catch {
    return { structures: [] };
  }
};

Structures.action = async ({ request }: ActionFunctionArgs) => {
  try {
    const { type, _structureId, ...body } = await request.json();

    let res;
    switch (type) {
      case 'GENERATE':
        res = await fetch(`${import.meta.env.VITE_API_URL}/v1/structures/generate`, { method: 'POST' });
        break;
      case 'SAVE':
        res = await fetch(`${import.meta.env.VITE_API_URL}/v1/structures/${_structureId || ''}`, {
          method: 'PUT',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify(body),
        });
        break;
      case 'RESET':
        res = await fetch(`${import.meta.env.VITE_API_URL}/v1/structures/all`, { method: 'DELETE' });
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
