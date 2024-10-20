import { useLoaderData } from 'react-router-dom';

import Tree from 'components/Tree';
import Button from 'components/Button';

import { cn } from 'lib/helpers';

type Structure = {
  id: number;
  parentId: number | null;
  title: string;
};

function Structures() {
  const data = useLoaderData() as { structures: Structure[] };
  const structures = new Map();
  data.structures.forEach((structure) => {
    structures.set(structure.parentId, [...(structures.get(structure.parentId) || []), structure]);
  });

  return (
    <div className="relative flex h-full flex-col items-center overflow-y-auto p-4 pt-16">
      {data.structures.length ? (
        <>
          <div className="absolute right-4 top-4 flex items-start">
            <Button className="-mr-1.5 -mt-1.5 bg-red-50 px-3 py-1.5 text-sm text-red-400 hover:bg-red-100 hover:text-red-500">
              Reset
            </Button>
          </div>

          <Tree
            items={structures}
            renderAdd={(parent) => (
              <Tree.Item
                className={cn(
                  'min-w-0 border-neutral-800 bg-neutral-900 text-white hover:bg-neutral-700',
                  parent && 'ml-[26px]',
                )}
                asChild
              >
                <button>Add Structure</button>
              </Tree.Item>
            )}
          />
        </>
      ) : (
        <div className="flex h-full flex-col items-center justify-center gap-2">
          <span>Start building the structures for the Scorecard</span>

          <div className="flex items-center gap-2">
            <Button className="text-sm" onClick={() => alert('TODO')}>
              Copy from Syllabuses
            </Button>
            <span>Or</span>
            <Button className="text-sm" onClick={() => alert('TODO')}>
              Manual
            </Button>
          </div>
        </div>
      )}
    </div>
  );
}

Structures.loader = async () => {
  const res = await fetch(`${import.meta.env.VITE_API_URL}/v1/structures`);
  return { structures: (await res.json()).nodes };
};

export default Structures;
