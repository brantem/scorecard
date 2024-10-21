import { Link, useLoaderData } from 'react-router-dom';
import { ArrowPathIcon } from '@heroicons/react/20/solid';
import { InformationCircleIcon } from '@heroicons/react/16/solid';
import dayjs from 'dayjs';

import Button from 'components/Button';
import Tooltip from 'components/Tooltip';

import type { Scorecard as _Scorecard } from 'types/scorecard';

function Scorecards() {
  const data = useLoaderData() as { scorecards: Omit<_Scorecard, 'items'>[] };

  return (
    <>
      <div className="flex items-start justify-between p-4 pb-0">
        <h1 className="font-semibold">Scorecards</h1>

        <Button className="pl-2.5 text-sm" onClick={() => alert('TODO')}>
          <ArrowPathIcon className="size-5" />
          <span>Generate</span>
        </Button>
      </div>

      <div className="m-4 grid grid-cols-1 gap-4 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4 xl:grid-cols-5 2xl:grid-cols-6">
        {data.scorecards.map((scorecard) => (
          <Link
            key={scorecard.id}
            to={`/${scorecard.id}`}
            className="group flex flex-col overflow-hidden rounded-lg border border-neutral-200"
          >
            <div className="flex flex-1 items-center justify-center py-16 text-5xl font-black tabular-nums">
              {scorecard.score}
            </div>
            <div className="border-t border-neutral-200 bg-neutral-50 p-3 group-hover:bg-neutral-100">
              <h4 className="font-semibold">{scorecard.user.name}</h4>
              <div className="mt-1 flex items-center gap-2">
                <span className="text-sm text-neutral-500">
                  {dayjs(scorecard.generatedAt).format('D MMM YYYY HH:mm')}
                </span>
                {scorecard.isOutdated && (
                  <Tooltip content="This scorecard needs to be regenerated." side="right">
                    <InformationCircleIcon className="size-4 text-yellow-500" />
                  </Tooltip>
                )}
              </div>
            </div>
          </Link>
        ))}
      </div>
    </>
  );
}

Scorecards.loader = async () => {
  try {
    const res = await fetch(`${import.meta.env.VITE_API_URL}/v1/scorecards`);
    return { scorecards: (await res.json()).nodes };
  } catch {
    return { scorecards: [] };
  }
};

export default Scorecards;