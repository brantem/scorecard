import { useEffect, useRef } from 'react';
import { useLoaderData, useFetcher, type ActionFunctionArgs, Link } from 'react-router-dom';
import { ExclamationTriangleIcon, InformationCircleIcon } from '@heroicons/react/24/outline';
import { PlusIcon } from '@heroicons/react/20/solid';
import { TrashIcon } from '@heroicons/react/16/solid';

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
                className="bg-red-50 px-3 py-1.5 pl-2 text-sm text-red-500 hover:bg-red-100"
                onClick={() => {
                  resetModalRef.current?.onOpen('Structures', { type: 'RESET_STRUCTURE', _structureId: '0' });
                }}
              >
                <TrashIcon className="size-4" />
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
                            _structureId: data.structures.length === 2 ? '0' : structure.id, // Reset if there is only 2 structures
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
                <PlusIcon className="size-5" />
                {!data.structures.length ? <span>Structure</span> : null}
              </button>

              {isStructuresLocked ? (
                <Tooltip content="You need to reset the syllabuses before updating the structures." side="right">
                  <ExclamationTriangleIcon className="absolute left-[calc(100%+theme(spacing.2))] top-2 size-6 text-yellow-500" />
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
                  <InformationCircleIcon className="absolute left-full top-[7px] ml-2 size-6 text-neutral-500" />
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
                className="bg-red-50 px-3 py-1.5 pl-2 text-sm text-red-500 hover:bg-red-100"
                onClick={() => {
                  resetModalRef.current?.onOpen('Syllabuses', { type: 'RESET_SYLLABUS', _syllabusId: '0' });
                }}
              >
                <TrashIcon className="size-4" />
                <span>Reset</span>
              </Button>
            ) : null}
          </div>

          {data.structures.length ? (
            data.syllabuses.length ? (
              <div className="no-scrollbar flex size-full flex-col items-center overflow-y-auto p-4 pt-12">
                <Tree
                  items={syllabuses}
                  renderOptions={(syllabus) => (
                    <>
                      <Tree.Item className="min-w-0 hover:bg-neutral-100" asChild>
                        <button
                          onClick={() => {
                            saveSyllabusModalRef.current?.onOpen(structures.get(syllabus.structureId), null, syllabus);
                          }}
                        >
                          Edit
                        </button>
                      </Tree.Item>

                      <Tree.Item className="min-w-0 border-red-200 bg-red-50 text-red-500 hover:bg-red-100" asChild>
                        <button
                          onClick={() => {
                            deleteModalRef.current?.onOpen('Syllabus', {
                              type: 'DELETE_SYLLABUS',
                              _syllabusId: syllabus.id,
                            });
                          }}
                        >
                          Delete
                        </button>
                      </Tree.Item>

                      {syllabus.structureId === assignmentStructure.id && (
                        <Tree.Item className="ml-4 hover:bg-neutral-100" asChild>
                          <Link to={`/syllabuses/${syllabus.id}/scores`}>Score</Link>
                        </Tree.Item>
                      )}
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
                          <PlusIcon className="size-5" />
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
                  <PlusIcon className="size-5" />
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
