import { Fragment, useEffect, useState } from 'react';
import { useFetcher } from 'react-router-dom';
import { ChevronRightIcon } from '@heroicons/react/16/solid';

import Button from 'components/Button';

import type { Syllabus } from 'types/syllabus';
import { cn } from 'lib/helpers';

type SyllabusListProps = {
  parentId: number | null;
  onCompleted(): void;
};

export default function SyllabusList({ parentId, onCompleted }: SyllabusListProps) {
  const fetcher = useFetcher();

  const [syllabuses, setSyllabuses] = useState<Map<Syllabus['parentId'], Syllabus[]>>(new Map());
  const [parents, setParents] = useState<Syllabus[]>([]);

  // not efficient
  useEffect(() => {
    (async () => {
      const res = await fetch(`${import.meta.env.VITE_API_URL}/v1/syllabuses`);
      const nodes = (await res.json()).nodes as Syllabus[];

      const m = new Map();
      nodes.forEach((node) => m.set(node.parentId, [...(m.get(node.parentId) || []), node]));
      setSyllabuses(m);
    })();
  }, []);

  useEffect(() => {
    if (!fetcher.data) return;
    if (fetcher.data.success) {
      onCompleted();
    } else {
      alert(fetcher.data.error?.code || 'INTERNAL_SERVER_ERROR');
    }
  }, [fetcher.data]);

  const lastParent = parents[parents.length - 1];

  return (
    <div className="mt-4">
      {parents.length ? (
        <div className="mb-2 flex items-center gap-1 text-sm">
          {parents.map((parent, i) => (
            <Fragment key={parent.id}>
              {i > 0 && <ChevronRightIcon className="size-4 [&:not(:last-child)]:text-neutral-500" />}
              <button
                className="hover:underline [&:not(:last-child)]:text-neutral-500"
                onClick={() => {
                  setParents((prev) => prev.slice(0, i));
                }}
              >
                {parent.title}
              </button>
            </Fragment>
          ))}
        </div>
      ) : null}

      <div className="flex flex-col divide-y divide-neutral-200 overflow-hidden rounded-lg border border-neutral-200">
        {(syllabuses.get(lastParent?.id || null) || []).map((syllabus) => {
          const isClickable = syllabuses.has(syllabus.id);
          return (
            <div key={syllabus.id} className="flex items-baseline justify-between gap-1.5">
              <button
                className={cn(
                  'inline-block h-full flex-1 pb-[10px] pl-3 text-left',
                  isClickable ? 'cursor-pointer hover:underline' : 'cursor-default',
                )}
                onClick={() => {
                  if (!isClickable) return;
                  setParents((prev) => [...prev, syllabus]);
                }}
              >
                {syllabus.title}
              </button>
              <Button
                className="m-1.5 bg-neutral-100 px-3 py-1.5 text-sm text-neutral-900 hover:bg-neutral-900 hover:text-white"
                onClick={() => {
                  fetcher.submit(
                    { type: 'COPY', _syllabusId: syllabus.id, parentId },
                    { method: 'PUT', encType: 'application/json' },
                  );
                }}
              >
                Copy
              </Button>
            </div>
          );
        })}
      </div>
    </div>
  );
}
