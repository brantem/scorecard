import { useLoaderData } from 'react-router-dom';

import Tree from 'components/Tree';

import { cn } from 'lib/helpers';

type Structure = {
  id: number;
  prevId: number | null;
  title: string;
};

type Syllabus = {
  id: number;
  parentId: number | null;
  structureId: number;
  title: string;
};

function Syllabuses() {
  const data = useLoaderData() as { structures: Structure[]; syllabuses: Syllabus[] };

  const structures = new Map();
  let lastStructureId;
  data.structures.forEach((structure) => {
    structures.set(structure.prevId, structure);
    if (structure.prevId !== -1) lastStructureId = structure.id;
  });
  structures.set(lastStructureId, structures.get(-1));

  const assignmentStructureId = structures.get(-1).id;

  const syllabuses = new Map<number | null, Syllabus[]>();
  data.syllabuses.forEach((structure) => {
    syllabuses.set(structure.parentId, [...(syllabuses.get(structure.parentId) || []), structure]);
  });

  return (
    <div className="grid h-full grid-cols-4 items-center justify-center divide-x overflow-hidden">
      <div className="relative flex h-full items-center justify-center p-4">
        <h2 className="absolute left-4 top-4 font-semibold">Structures</h2>

        <div className="relative flex flex-col items-center gap-6 font-medium">
          {data.structures.map((structure) => {
            if (structure.prevId === -1) return;
            return <Structure key={structure.id}>{structure.title}</Structure>;
          })}
          <button
            className={cn(
              'relative z-10 flex items-center justify-center rounded-full bg-neutral-900 text-white hover:bg-neutral-700',
              data.structures.length ? 'size-10' : 'h-10 gap-2 px-3',
            )}
          >
            <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor" className="size-5">
              <path d="M10.75 4.75a.75.75 0 0 0-1.5 0v4.5h-4.5a.75.75 0 0 0 0 1.5h4.5v4.5a.75.75 0 0 0 1.5 0v-4.5h4.5a.75.75 0 0 0 0-1.5h-4.5v-4.5Z" />
            </svg>
            {!data.structures.length ? <span>Structure</span> : null}
          </button>
          {data.structures.length ? <Structure>Assignment</Structure> : null}
          <div className="absolute bottom-0 left-1/2 top-0 h-full w-px -translate-x-1/2 bg-neutral-300" />
        </div>
      </div>

      <div className="relative col-span-3 h-full overflow-hidden">
        <h2 className="absolute left-4 top-4 font-semibold">Syllabuses</h2>
        <div className="flex h-full w-full flex-col items-center overflow-y-auto p-4 pt-16">
          <Tree
            items={syllabuses}
            renderAdd={(parent) => {
              const structureId = parent?.structureId;
              if (structureId === assignmentStructureId) return null;
              return (
                <Tree.Item
                  className={cn(
                    'min-w-0 border-neutral-800 bg-neutral-900 text-white hover:bg-neutral-700',
                    parent && 'ml-[26px]',
                  )}
                  asChild
                >
                  <button>Add {parent ? structures.get(structureId).title : structures.get(null).title}</button>
                </Tree.Item>
              );
            }}
          />
        </div>
      </div>
    </div>
  );
}

Syllabuses.loader = async () => {
  const [structures, syllabuses] = await Promise.all([
    (async () => {
      const res = await fetch(`${import.meta.env.VITE_API_URL}/v1/syllabuses/structures`);
      return (await res.json()).nodes;
    })(),
    (async () => {
      const res = await fetch(`${import.meta.env.VITE_API_URL}/v1/syllabuses`);
      return (await res.json()).nodes;
    })(),
  ]);
  return { structures, syllabuses };
};

export default Syllabuses;

function Structure({ children }: { children: string }) {
  return (
    <div className="relative z-10 w-fit min-w-32 rounded-full border border-neutral-200 bg-neutral-50 px-3 py-1.5 text-center">
      {children}
    </div>
  );
}
