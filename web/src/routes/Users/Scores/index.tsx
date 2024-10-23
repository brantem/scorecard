import { Fragment, useRef } from 'react';
import {
  Link,
  useParams,
  useLoaderData,
  redirect,
  type LoaderFunctionArgs,
  type ActionFunctionArgs,
} from 'react-router-dom';
import { ArrowLeftIcon } from '@heroicons/react/20/solid';
import { ChevronRightIcon } from '@heroicons/react/16/solid';

import Table from 'components/Table';
import Button from 'components/Button';
import SaveModal, { type SaveModalHandle } from 'components/scores/SaveModal';

import type { Score, User } from 'types/user';

function UserScores() {
  const params = useParams();
  const data = useLoaderData() as { user: User; scores: Score[] };

  const saveModalRef = useRef<SaveModalHandle>(null);

  const m = new Map<string, (Pick<Score, 'score'> & { syllabus: Omit<Score['syllabus'], 'parents'> })[]>();
  data.scores.forEach(({ syllabus: { parents, ...syllabus }, score }) => {
    const key = JSON.stringify(parents); // lol
    m.set(key, [...(m.get(key) || []), { syllabus, score }]);
  });

  return (
    <>
      <div className="flex items-start gap-4 p-4 pb-0">
        <Link
          to={`/${params.programId}/users`}
          className="flex aspect-square h-12 items-center justify-center rounded-lg bg-neutral-50 hover:bg-neutral-100"
        >
          <ArrowLeftIcon className="size-5" />
        </Link>

        <div className="flex flex-col gap-1">
          <h2 className="font-semibold">{data.user.name}</h2>
          <span className="inline-block text-sm text-neutral-500">{data.scores.length} Assignments</span>
        </div>
      </div>

      {data.scores.length ? (
        <div className="flex flex-col gap-4 p-4">
          {[...m.keys()].map((key, i) => {
            const parents = JSON.parse(key) as Score['syllabus']['parents'];
            return (
              <div key={i} className="overflow-x-auto rounded-lg border border-neutral-200">
                <div className="flex items-center gap-1 border-b bg-neutral-50 px-3 py-1.5 text-sm text-neutral-500">
                  {parents.map((parent, i) => (
                    <Fragment key={parent.id}>
                      {i > 0 && <ChevronRightIcon className="size-4" />}
                      <span>{parent.title}</span>
                    </Fragment>
                  ))}
                </div>

                <Table>
                  <thead>
                    <tr>
                      <Table.Th>Title</Table.Th>
                      <Table.Th className="w-24">Score</Table.Th>
                      <Table.Th className="w-24 [&>div]:justify-end">Actions</Table.Th>
                    </tr>
                  </thead>

                  <tbody>
                    {(m.get(key) || []).map((node) => (
                      <tr key={node.syllabus.id}>
                        <Table.Td>{node.syllabus.title}</Table.Td>
                        <Table.Td className="tabular-nums">
                          {typeof node.score === 'number' ? node.score : <span className="text-neutral-400">-</span>}
                        </Table.Td>

                        <Table.Td className="text-sm [&>div]:justify-end [&>div]:pr-1.5">
                          <button
                            className="flex h-8 items-center rounded-lg border border-neutral-200 bg-neutral-50 px-3 hover:bg-neutral-100"
                            onClick={() => saveModalRef.current?.open(node.syllabus, data.user, node.score)}
                          >
                            {typeof node.score === 'number' ? 'Edit' : 'Add'}
                          </button>
                        </Table.Td>
                      </tr>
                    ))}
                  </tbody>
                </Table>
              </div>
            );
          })}
        </div>
      ) : (
        <div className="absolute inset-0 flex h-full flex-col items-center justify-center gap-2 text-neutral-500">
          <span>You need to create the syllabuses before adding any scores</span>
          <Button className="rounded-full" asChild>
            <Link to={`/${params.programId}/syllabuses`}>Syllabuses</Link>
          </Button>
        </div>
      )}

      <SaveModal
        ref={saveModalRef}
        description={({ syllabus, isEditing }) => (
          <>
            You are {isEditing ? 'editing the score' : 'adding a score'} for <b>{syllabus.title}</b>
          </>
        )}
      />
    </>
  );
}

UserScores.loader = async ({ params }: LoaderFunctionArgs) => {
  const [user, scores] = await Promise.all([
    (async () => {
      try {
        const res = await fetch(
          `${import.meta.env.VITE_API_URL}/v1/programs/${params.programId}/users/${params.userId}`,
        );
        return (await res.json()).user;
      } catch {
        return null;
      }
    })(),
    (async () => {
      try {
        const res = await fetch(
          `${import.meta.env.VITE_API_URL}/v1/programs/${params.programId}/users/${params.userId}/scores`,
        );
        return (await res.json()).nodes || [];
      } catch {
        return [];
      }
    })(),
  ]);
  if (!user) return redirect('/users');
  return { user, scores };
};

UserScores.action = async ({ request, params }: ActionFunctionArgs) => {
  try {
    const { type, _syllabusId, _userId, ...body } = await request.json();

    let res;
    switch (type) {
      case 'SAVE':
        res = await fetch(
          `${import.meta.env.VITE_API_URL}/v1/programs/${params.programId}/syllabuses/${_syllabusId}/scores/${_userId}`,
          {
            method: 'PUT',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(body),
          },
        );
        break;
      default:
        return { success: false, error: null };
    }
    return await res.json();
  } catch {
    return { success: false, error: { code: 'INTERNAL_SERVER_ERROR' } };
  }
};

export default UserScores;
