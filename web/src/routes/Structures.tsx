import { useLoaderData } from 'react-router-dom';

import Tree from 'components/Tree';

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
    </div>
  );
}

Structures.loader = async () => {
  return {
    structures: [
      {
        id: 1,
        parentId: null,
        title: 'History',
      },
      {
        id: 2,
        parentId: null,
        title: 'Science',
      },
      {
        id: 3,
        parentId: 1,
        title: 'World War 1',
      },
      {
        id: 4,
        parentId: 1,
        title: 'World War 2',
      },
      {
        id: 5,
        parentId: 3,
        title: 'Battle of Gallipoli',
      },
      {
        id: 6,
        parentId: 3,
        title: 'Battle of Jutland',
      },
      {
        id: 7,
        parentId: 4,
        title: 'Dunkirk',
      },
      {
        id: 8,
        parentId: null,
        title: 'Math',
      },
      {
        id: 9,
        parentId: 8,
        title: '3D',
      },
      {
        id: 10,
        parentId: 9,
        title: 'Cube',
      },
      {
        id: 11,
        parentId: 10,
        title: 'Area',
      },
      {
        id: 12,
        parentId: 10,
        title: 'Volume',
      },
      {
        id: 13,
        parentId: 9,
        title: 'Ball',
      },
      {
        id: 14,
        parentId: 13,
        title: 'Area',
      },
      {
        id: 15,
        parentId: 13,
        title: 'Volume',
      },
    ],
  };
  const res = await fetch(`${import.meta.env.VITE_API_URL}/v1/structures`);
  return { structures: (await res.json()).nodes };
};

export default Structures;
