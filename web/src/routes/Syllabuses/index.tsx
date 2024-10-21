import { useEffect, useRef } from 'react';
import { useLoaderData, useFetcher, type ActionFunctionArgs, Link } from 'react-router-dom';

import Button from 'components/Button';
import Tooltip from 'components/Tooltip';
import Tree from 'components/Tree';
import ResetModal, { type ResetModalHandle } from 'components/ResetModal';
import SaveStructureModal, { type SaveStructureModalHandle } from './SaveStructureModal';
import SaveSyllabusModal, { type SaveSyllabusModalHandle } from './SaveSyllabusModal';
import DeleteModal, { type DeleteModalHandle } from 'components/DeleteModal';

import { cn } from 'lib/helpers';
import type { Structure, Syllabus } from 'types/syllabus';

function Syllabuses() {
  const data = useLoaderData() as { structures: Structure[]; syllabuses: Syllabus[] };
  const fetcher = useFetcher();

  const resetModalRef = useRef<ResetModalHandle>(null);
  const saveStructureModalRef = useRef<SaveStructureModalHandle>(null);
  const saveSyllabusModalRef = useRef<SaveSyllabusModalHandle>(null);
  const deleteModalRef = useRef<DeleteModalHandle>(null);

  useEffect(() => {
    if (!fetcher.data) return;
    if (fetcher.data.success) {
      switch ((fetcher.json as { type?: string }).type) {
        case 'RESET_STRUCTURE':
        case 'RESET_SYLLABUS':
          resetModalRef.current?.onClose();
          break;
        case 'DELETE_STRUCTURE':
        case 'DELETE_SYLLABUS':
          deleteModalRef.current?.onClose();
          break;
      }
    } else {
      alert(fetcher.data.error?.code || 'INTERNAL_SERVER_ERROR');
    }
  }, [fetcher.data]);

  const structures = new Map();
  const prevStructures = new Map();
  data.structures.forEach((structure) => {
    structures.set(structure.id, structure);
    prevStructures.set(structure.prevId, structure);
  });
  const lastStructure = data.structures[data.structures.length - 1];
  if (lastStructure) prevStructures.set(lastStructure.id, prevStructures.get(-1));

  const assignmentStructure = prevStructures.get(-1);

  const syllabuses = new Map<number | null, Syllabus[]>();
  data.syllabuses.forEach((structure) => {
    syllabuses.set(structure.parentId, [...(syllabuses.get(structure.parentId) || []), structure]);
  });

  const isStructuresLocked = data.syllabuses.length > 0;

  return (
    <>
      <div className="grid size-full grid-cols-1 items-center justify-center overflow-hidden max-lg:grid-rows-3 max-lg:divide-y lg:grid-cols-4 lg:divide-x">
        <div className="relative flex size-full items-center justify-center p-4">
          <div className="absolute left-4 right-4 top-4 flex items-start justify-between gap-4 font-semibold">
            <h2>Structures</h2>

            {data.structures.length && !isStructuresLocked ? (
              <Button
                className="-mr-1.5 -mt-1.5 bg-red-50 px-3 py-1.5 pl-2 text-sm text-red-500 hover:bg-red-100"
                onClick={() => {
                  resetModalRef.current?.onOpen('Structures', { type: 'RESET_STRUCTURE', _structureId: 'all' });
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
            {data.structures.map((structure, i) => {
              if (structure.prevId === -1) return;
              return (
                <div key={structure.id} className="group relative">
                  {!isStructuresLocked && (
                    <div className="absolute right-full top-0 flex h-[38px] items-center pr-2 opacity-0 group-hover:opacity-100">
                      <button
                        className="rounded-full border border-neutral-200 bg-neutral-50 px-3 py-1 text-sm hover:bg-neutral-100"
                        onClick={() => {
                          saveStructureModalRef.current?.onOpen(i > 0 ? data.structures[i - 1] : null, structure);
                        }}
                      >
                        Edit
                      </button>
                    </div>
                  )}

                  <Structure>{structure.title}</Structure>

                  {!isStructuresLocked && (
                    <div className="absolute left-full top-0 flex h-[38px] items-center pl-2 opacity-0 group-hover:opacity-100">
                      <button
                        className="rounded-full border border-red-200 bg-red-50 px-3 py-1 text-sm text-red-500 hover:bg-red-100"
                        onClick={() => {
                          deleteModalRef.current?.onOpen('Structure', {
                            type: 'DELETE_STRUCTURE',
                            _structureId: structure.id,
                          });
                        }}
                      >
                        Delete
                      </button>
                    </div>
                  )}
                </div>
              );
            })}
            <div className="relative">
              <button
                className={cn(
                  'relative z-10 flex items-center justify-center rounded-full bg-neutral-900 text-white hover:bg-neutral-800 disabled:bg-neutral-200 disabled:text-neutral-400',
                  data.structures.length ? 'size-10' : 'h-10 gap-2 px-3',
                )}
                onClick={() => saveStructureModalRef.current?.onOpen(lastStructure, null)}
                disabled={isStructuresLocked}
              >
                <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor" className="size-5">
                  <path d="M10.75 4.75a.75.75 0 0 0-1.5 0v4.5h-4.5a.75.75 0 0 0 0 1.5h4.5v4.5a.75.75 0 0 0 1.5 0v-4.5h4.5a.75.75 0 0 0 0-1.5h-4.5v-4.5Z" />
                </svg>
                {!data.structures.length ? <span>Structure</span> : null}
              </button>

              {isStructuresLocked ? (
                <Tooltip content="You need to reset the syllabuses before updating the structures." side="right">
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
            {data.structures.length ? (
              <div className="relative">
                <Structure>Assignment</Structure>
                <Tooltip
                  content="This structure is automatically generated and can't be edited. It must always be last."
                  side="right"
                >
                  <svg
                    xmlns="http://www.w3.org/2000/svg"
                    fill="none"
                    viewBox="0 0 24 24"
                    strokeWidth={1.5}
                    stroke="currentColor"
                    className="absolute left-full top-[7px] ml-2 size-6 text-neutral-500"
                  >
                    <path
                      strokeLinecap="round"
                      strokeLinejoin="round"
                      d="m11.25 11.25.041-.02a.75.75 0 0 1 1.063.852l-.708 2.836a.75.75 0 0 0 1.063.853l.041-.021M21 12a9 9 0 1 1-18 0 9 9 0 0 1 18 0Zm-9-3.75h.008v.008H12V8.25Z"
                    />
                  </svg>
                </Tooltip>
              </div>
            ) : null}
            <div className="absolute bottom-0 left-1/2 top-0 h-full w-px -translate-x-1/2 bg-neutral-300" />
          </div>
        </div>

        <div className="relative size-full overflow-hidden max-lg:row-span-2 lg:col-span-3">
          <div className="absolute left-4 right-4 top-4 flex items-start justify-between gap-4 font-semibold">
            <h2>Syllabuses</h2>

            {data.syllabuses.length ? (
              <Button
                className="-mr-1.5 -mt-1.5 bg-red-50 px-3 py-1.5 pl-2 text-sm text-red-500 hover:bg-red-100"
                onClick={() => {
                  resetModalRef.current?.onOpen('Syllabuses', { type: 'RESET_SYLLABUS', _syllabusId: 'all' });
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

          {data.structures.length ? (
            data.syllabuses.length ? (
              <div className="flex h-full w-full flex-col items-center overflow-y-auto p-4 pt-16">
                <Tree
                  items={syllabuses}
                  renderOptions={(syllabus) => (
                    <>
                      {syllabus.structureId === assignmentStructure.id && (
                        <Link
                          to={`/syllabuses/${syllabus.id}/scores`}
                          className="flex h-[34px] items-center justify-center rounded-lg border border-neutral-200 bg-neutral-50 px-3 hover:bg-neutral-100"
                        >
                          Score
                        </Link>
                      )}

                      <button
                        className="flex h-[34px] items-center justify-center rounded-lg border border-neutral-200 bg-neutral-50 px-3 hover:bg-neutral-100"
                        onClick={() => {
                          saveSyllabusModalRef.current?.onOpen(structures.get(syllabus.structureId), null, syllabus);
                        }}
                      >
                        Edit
                      </button>

                      <button
                        className="flex h-[34px] items-center justify-center rounded-lg border border-red-200 bg-red-50 px-3 text-red-500 hover:bg-red-100"
                        onClick={() => {
                          deleteModalRef.current?.onOpen('Syllabus', {
                            type: 'DELETE_SYLLABUS',
                            _syllabusId: syllabus.id,
                          });
                        }}
                      >
                        Delete
                      </button>
                    </>
                  )}
                  renderAdd={(parent) => {
                    const structureId = parent?.structureId;
                    if (assignmentStructure && structureId === assignmentStructure.id) return null;

                    const structure = parent ? prevStructures.get(structureId) : prevStructures.get(null);

                    return (
                      <Tree.Item
                        className={cn(
                          'flex min-w-0 gap-2 border-neutral-800 bg-neutral-900 pl-2 text-white hover:bg-neutral-800',
                          parent && 'ml-[calc(theme(spacing.8)+2px)]',
                        )}
                        asChild
                      >
                        <button onClick={() => saveSyllabusModalRef.current?.onOpen(structure, parent, null)}>
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
                  onClick={() => saveSyllabusModalRef.current?.onOpen(prevStructures.get(null), null, null)}
                >
                  <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor" className="size-5">
                    <path d="M10.75 4.75a.75.75 0 0 0-1.5 0v4.5h-4.5a.75.75 0 0 0 0 1.5h4.5v4.5a.75.75 0 0 0 1.5 0v-4.5h4.5a.75.75 0 0 0 0-1.5h-4.5v-4.5Z" />
                  </svg>
                  <span>Add {prevStructures.get(null).title}</span>
                </Button>
              </div>
            )
          ) : (
            <div className="flex h-full items-center justify-center text-neutral-500">
              You need to create the structures before adding any syllabuses
            </div>
          )}
        </div>
      </div>

      <ResetModal
        ref={resetModalRef}
        onAccept={(body) => fetcher.submit(body, { method: 'DELETE', encType: 'application/json' })}
      />
      <SaveStructureModal ref={saveStructureModalRef} />
      <SaveSyllabusModal ref={saveSyllabusModalRef} />
      <DeleteModal
        ref={deleteModalRef}
        onAccept={(body) => fetcher.submit(body, { method: 'DELETE', encType: 'application/json' })}
      />
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
      case 'DELETE_STRUCTURE':
        res = await fetch(`${import.meta.env.VITE_API_URL}/v1/syllabuses/structures/${_structureId || ''}`, {
          method: 'DELETE',
        });
        break;
      case 'SAVE_SYLLABUS':
        res = await fetch(`${import.meta.env.VITE_API_URL}/v1/syllabuses/${_syllabusId || ''}`, {
          method: 'PUT',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify(body),
        });
        break;
      case 'RESET_SYLLABUS':
      case 'DELETE_SYLLABUS':
        res = await fetch(`${import.meta.env.VITE_API_URL}/v1/syllabuses/${_syllabusId || ''}`, { method: 'DELETE' });
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
