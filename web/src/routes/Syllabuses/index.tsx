import { useRef } from 'react';
import { useLoaderData, useFetcher, type ActionFunctionArgs } from 'react-router-dom';

import Button from 'components/Button';
import Tooltip from 'components/Tooltip';
import Tree from 'components/Tree';
import SaveStructureModal, { type SaveStructureModalHandle } from './SaveStructureModal';
import SaveSyllabusModal, { type SaveSyllabusModalHandle } from './SaveSyllabusModal';

import { cn } from 'lib/helpers';
import type { SyllabusStructure, Syllabus } from 'types';

function Syllabuses() {
  const data = useLoaderData() as { structures: SyllabusStructure[]; syllabuses: Syllabus[] };
  const fetcher = useFetcher();

  const saveStructureModalRef = useRef<SaveStructureModalHandle>(null);
  const saveSyllabusModalRef = useRef<SaveSyllabusModalHandle>(null);

  const structures = new Map();
  data.structures.forEach((structure) => structures.set(structure.prevId, structure));
  const lastStructure = data.structures[data.structures.length - 1];
  structures.set(lastStructure.id, structures.get(-1));

  const assignmentStructureId = structures.get(-1).id;

  const syllabuses = new Map<number | null, Syllabus[]>();
  data.syllabuses.forEach((structure) => {
    syllabuses.set(structure.parentId, [...(syllabuses.get(structure.parentId) || []), structure]);
  });

  return (
    <>
      <div className="grid h-full grid-cols-4 items-center justify-center divide-x overflow-hidden">
        <div className="relative flex h-full items-center justify-center p-4">
          <div className="absolute left-4 right-4 top-4 flex items-start justify-between gap-4 font-semibold">
            <h2>Structures</h2>

            {!data.syllabuses.length ? (
              <Button
                className="-mr-1.5 -mt-1.5 bg-red-50 px-3 py-1.5 pl-2 text-sm text-red-400 hover:bg-red-100 hover:text-red-500"
                onClick={() => {
                  fetcher.submit({ type: 'RESET_SYLLABUS' }, { method: 'DELETE', encType: 'application/json' });
                }}
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

          <div className="relative flex flex-col items-center gap-6 font-medium">
            {data.structures.map((structure) => {
              if (structure.prevId === -1) return;
              return <Structure key={structure.id}>{structure.title}</Structure>;
            })}
            <div className="relative">
              <button
                className={cn(
                  'relative z-10 flex items-center justify-center rounded-full bg-neutral-900 text-white hover:bg-neutral-800 disabled:bg-neutral-200 disabled:text-neutral-400',
                  data.structures.length ? 'size-10' : 'h-10 gap-2 px-3',
                )}
                onClick={() => saveStructureModalRef.current?.onOpen(lastStructure)}
                disabled={!!data.syllabuses.length}
              >
                <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor" className="size-5">
                  <path d="M10.75 4.75a.75.75 0 0 0-1.5 0v4.5h-4.5a.75.75 0 0 0 0 1.5h4.5v4.5a.75.75 0 0 0 1.5 0v-4.5h4.5a.75.75 0 0 0 0-1.5h-4.5v-4.5Z" />
                </svg>
                {!data.structures.length ? <span>Structure</span> : null}
              </button>

              {data.syllabuses.length ? (
                <Tooltip content="You need to reset the syllabuses before updating this structure." side="right">
                  <svg
                    xmlns="http://www.w3.org/2000/svg"
                    fill="none"
                    viewBox="0 0 24 24"
                    strokeWidth={1.5}
                    stroke="currentColor"
                    className="absolute left-[calc(100%+theme(spacing.2))] top-2 size-6 text-yellow-500"
                  >
                    <path
                      strokeLinecap="round"
                      strokeLinejoin="round"
                      d="M12 9v3.75m-9.303 3.376c-.866 1.5.217 3.374 1.948 3.374h14.71c1.73 0 2.813-1.874 1.948-3.374L13.949 3.378c-.866-1.5-3.032-1.5-3.898 0L2.697 16.126ZM12 15.75h.007v.008H12v-.008Z"
                    />
                  </svg>
                </Tooltip>
              ) : null}
            </div>
            {data.structures.length ? <Structure>Assignment</Structure> : null}
            <div className="absolute bottom-0 left-1/2 top-0 h-full w-px -translate-x-1/2 bg-neutral-300" />
          </div>
        </div>

        <div className="relative col-span-3 h-full overflow-hidden">
          <div className="absolute left-4 right-4 top-4 flex items-start justify-between gap-4 font-semibold">
            <h2>Syllabuses</h2>

            <Button
              className="-mr-1.5 -mt-1.5 bg-red-50 px-3 py-1.5 pl-2 text-sm text-red-400 hover:bg-red-100 hover:text-red-500"
              onClick={() => {
                fetcher.submit({ type: 'RESET_SYLLABUS' }, { method: 'DELETE', encType: 'application/json' });
              }}
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
          </div>

          {data.syllabuses.length ? (
            <div className="flex h-full w-full flex-col items-center overflow-y-auto p-4 pt-16">
              <Tree
                items={syllabuses}
                renderAdd={(parent) => {
                  const structureId = parent?.structureId;
                  if (structureId === assignmentStructureId) return null;

                  const structure = parent ? structures.get(structureId) : structures.get(null);

                  return (
                    <Tree.Item
                      className={cn(
                        'flex min-w-0 gap-2 border-neutral-800 bg-neutral-900 pl-2 text-white hover:bg-neutral-800',
                        parent && 'ml-[calc(theme(spacing.8)+2px)]',
                      )}
                      asChild
                    >
                      <button onClick={() => saveSyllabusModalRef.current?.onOpen(structure, parent)}>
                        <svg
                          xmlns="http://www.w3.org/2000/svg"
                          viewBox="0 0 20 20"
                          fill="currentColor"
                          className="size-5"
                        >
                          <path d="M10.75 4.75a.75.75 0 0 0-1.5 0v4.5h-4.5a.75.75 0 0 0 0 1.5h4.5v4.5a.75.75 0 0 0 1.5 0v-4.5h4.5a.75.75 0 0 0 0-1.5h-4.5v-4.5Z" />
                        </svg>
                        <span>Add {structure.title}</span>
                      </button>
                    </Tree.Item>
                  );
                }}
              />
            </div>
          ) : (
            <div className="flex h-full flex-col items-center justify-center gap-2">
              <Button
                className="pl-2.5 text-sm"
                onClick={() => saveSyllabusModalRef.current?.onOpen(structures.get(null), null)}
              >
                <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor" className="size-5">
                  <path d="M10.75 4.75a.75.75 0 0 0-1.5 0v4.5h-4.5a.75.75 0 0 0 0 1.5h4.5v4.5a.75.75 0 0 0 1.5 0v-4.5h4.5a.75.75 0 0 0 0-1.5h-4.5v-4.5Z" />
                </svg>
                <span>Add {structures.get(null).title}</span>
              </Button>
            </div>
          )}
        </div>
      </div>

      <SaveStructureModal ref={saveStructureModalRef} />
      <SaveSyllabusModal ref={saveSyllabusModalRef} />
    </>
  );
}

Syllabuses.loader = async () => {
  const [structures, syllabuses] = await Promise.all([
    (async () => {
      try {
        const res = await fetch(`${import.meta.env.VITE_API_URL}/v1/syllabuses/structures`);
        return (await res.json()).nodes;
      } catch {
        return [];
      }
    })(),
    (async () => {
      try {
        const res = await fetch(`${import.meta.env.VITE_API_URL}/v1/syllabuses`);
        return (await res.json()).nodes;
      } catch {
        return [];
      }
    })(),
  ]);
  return { structures, syllabuses };
};

Syllabuses.action = async ({ request }: ActionFunctionArgs) => {
  try {
    const { type, _structureId, _syllabusId, ...body } = await request.json();

    let res;
    switch (type) {
      case 'SAVE_STRUCTURE':
        res = await fetch(`${import.meta.env.VITE_API_URL}/v1/syllabuses/structures/${_structureId || ''}`, {
          method: 'PUT',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify(body),
        });
        break;
      case 'RESET_STRUCTURE':
        res = await fetch(`${import.meta.env.VITE_API_URL}/v1/syllabuses/structures/all`, { method: 'DELETE' });
        break;
      case 'SAVE_SYLLABUS':
        res = await fetch(`${import.meta.env.VITE_API_URL}/v1/syllabuses/${_syllabusId || ''}`, {
          method: 'PUT',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify(body),
        });
        break;
      case 'RESET_SYLLABUS':
        res = await fetch(`${import.meta.env.VITE_API_URL}/v1/syllabuses/all`, { method: 'DELETE' });
        break;
      default:
        return { success: false, error: null };
    }
    return await res.json();
  } catch {
    return { success: false, error: { code: 'INTERNAL_SERVER_ERROR' } };
  }
};

export default Syllabuses;

function Structure({ children }: { children: string }) {
  return (
    <div className="relative z-10 w-fit min-w-32 rounded-full border border-neutral-200 bg-neutral-50 px-3 py-1.5 text-center">
      {children}
    </div>
  );
}
