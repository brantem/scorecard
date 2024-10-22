import {
  useLoaderData,
  useFetcher,
  redirect,
  type LoaderFunctionArgs,
  type ActionFunctionArgs,
} from 'react-router-dom';
import { ArrowPathIcon } from '@heroicons/react/20/solid';
import { InformationCircleIcon } from '@heroicons/react/16/solid';
import dayjs from 'dayjs';

import Button from 'components/Button';
import Tooltip from 'components/Tooltip';

import type * as types from 'types/scorecard';

function Scorecard() {
  const data = useLoaderData() as { scorecard: types.Scorecard; structures: Omit<types.Structure, 'syllabus'>[] };
  const fetcher = useFetcher();

  const structures = new Map<types.Structure['parentId'], typeof data.structures>();
  data.structures.forEach((structure) => {
    structures.set(structure.parentId, [...(structures.get(structure.parentId) || []), structure]);
  });

  const scores = new Map<types.ScorecardItem['structureId'], number>();
  data.scorecard.items.forEach((item) => scores.set(item.structureId, item.score));

  const roots = structures.get(null) || [];
  roots.forEach((root) => {
    const children = structures.get(root.id) || [];
    const total = children.reduce((total, structure) => total + (scores.get(structure.id) || 0), 0);
    scores.set(root.id, total / children.length);
  });

  return (
    <>
      <div className="flex items-start justify-between p-4 pb-0">
        <div className="flex flex-col gap-1">
          <h1 className="font-semibold">{data.scorecard.user.name}</h1>
          <div className="flex items-center gap-2">
            <span className="text-sm text-neutral-500">
              {dayjs(data.scorecard.generatedAt).format('D MMM YYYY HH:mm')}
            </span>
            {data.scorecard.isOutdated && (
              <Tooltip content="This scorecard needs to be regenerated." side="right">
                <InformationCircleIcon className="size-4 text-yellow-500" />
              </Tooltip>
            )}
          </div>
        </div>

        <Button
          className="pl-2.5 text-sm"
          onClick={() => fetcher.submit({ type: 'GENERATE' }, { method: 'POST', encType: 'application/json' })}
        >
          <ArrowPathIcon className="size-5" />
          <span>Generate</span>
        </Button>
      </div>

      <div className="flex flex-col gap-4 p-4">
        {roots.map((structure, i) => {
          const children = structures.get(structure.id) || [];
          return (
            <div
              key={structure.id}
              className="divide-y divide-neutral-200 overflow-hidden rounded-lg border border-neutral-200"
            >
              <div className="grid h-10 grid-cols-[30px_1fr] items-center bg-neutral-50">
                <span className="text-center tabular-nums">{i + 1}</span>
                <span className="pl-1">{structure.title}</span>
              </div>

              <div className="grid h-10 grid-cols-[repeat(6,1fr)_90px_90px_110px] items-center divide-x divide-neutral-200 bg-neutral-50">
                <div className="col-span-7 pl-2.5">Item</div>
                <div className="flex h-full items-center pl-2.5">Score</div>
                <div className="flex h-full items-center pl-2.5">Total Score</div>
              </div>

              <div className="grid grid-cols-[1fr_110px] items-center divide-y divide-neutral-200">
                {children.map((child) => (
                  <div className="grid h-10 grid-cols-[repeat(6,1fr)_90px_90px] items-center">
                    <span className="col-span-7 pl-2.5">{child.title}</span>
                    <div className="flex h-full items-center border-l px-2.5 tabular-nums">
                      {scores.get(child.id) || 0}
                    </div>
                  </div>
                ))}
                <div
                  className="col-start-2 flex h-full items-center justify-center !border-t-0 border-l px-2.5 tabular-nums"
                  style={{ gridRow: `1 / span ${children.length}` }}
                >
                  {scores.get(structure.id)}
                </div>
              </div>
            </div>
          );
        })}

        <div>
          <div className="grid h-10 grid-cols-[repeat(6,1fr)_90px_90px_111px] items-center divide-x divide-neutral-700 rounded-t-lg bg-neutral-900 text-white">
            <div className="col-span-7 pl-[11px]">Item</div>
            <div className="flex h-full items-center pl-2.5">Score</div>
            <div className="flex h-full items-center pl-2.5">Total Score</div>
          </div>

          <div className="grid grid-cols-[1fr_110px] items-center divide-y divide-neutral-200 overflow-hidden rounded-b-lg border border-neutral-200">
            {roots.map((structure) => (
              <div className="grid h-10 grid-cols-[repeat(6,1fr)_90px_90px] items-center">
                <span className="col-span-7 pl-2.5">{structure.title}</span>
                <div className="flex h-full items-center border-l px-2.5 tabular-nums">{scores.get(structure.id)}</div>
              </div>
            ))}
            <div
              className="col-start-2 flex h-full items-center justify-center !border-t-0 border-l px-2.5 tabular-nums"
              style={{ gridRow: `1 / span ${roots.length}` }}
            >
              {data.scorecard.score}
            </div>
          </div>
        </div>
      </div>
    </>
  );
}

Scorecard.loader = async ({ params }: LoaderFunctionArgs) => {
  const [structures, scorecard] = await Promise.all([
    (async () => {
      try {
        const res = await fetch(`${import.meta.env.VITE_API_URL}/v1/scorecards/structures`);
        return (await res.json()).nodes;
      } catch {
        return [];
      }
    })(),
    (async () => {
      try {
        const res = await fetch(`${import.meta.env.VITE_API_URL}/v1/scorecards/${params.scorecardId}`);
        return (await res.json()).scorecard || null;
      } catch {
        return null;
      }
    })(),
  ]);
  if (!scorecard) return redirect('/');
  return { structures, scorecard };
};

Scorecard.action = async ({ request, params }: ActionFunctionArgs) => {
  try {
    const { type } = await request.json();

    let res;
    switch (type) {
      case 'GENERATE':
        res = await fetch(`${import.meta.env.VITE_API_URL}/v1/scorecards/generate/${params.scorecardId}`, {
          method: 'POST',
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

export default Scorecard;
