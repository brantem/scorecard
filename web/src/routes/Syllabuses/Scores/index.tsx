import { Fragment, useRef } from 'react';
import { Link, useLoaderData, redirect, type LoaderFunctionArgs, type ActionFunctionArgs } from 'react-router-dom';
import { ChevronRightIcon } from '@heroicons/react/16/solid';

import Table from 'components/Table';
import SaveModal, { type SaveModalHandle } from 'components/scores/SaveModal';

import type { BaseSyllabus, Score } from 'types/syllabus';

type Syllabus = BaseSyllabus & {
  parents: BaseSyllabus[];
};

function SyllabusScores() {
  const data = useLoaderData() as { syllabus: Syllabus; scores: Score[] };

  const saveModalRef = useRef<SaveModalHandle>(null);

  return (
    <>
      <div className="flex flex-col gap-1 p-4 pb-0">
        <h1 className="font-semibold">{data.syllabus.title}</h1>
        <div className="flex items-center gap-1 text-sm text-neutral-500">
          <Link to="/syllabuses" className="hover:underline">
            Syllabuses
          </Link>
          {data.syllabus.parents.map((parent) => (
            <Fragment key={parent.id}>
              <ChevronRightIcon className="size-4" />
              <span>{parent.title}</span>
            </Fragment>
          ))}
        </div>
      </div>

      <div className="m-4 overflow-x-auto rounded-lg border border-neutral-200">
        <Table>
          <thead>
            <tr>
              <Table.Th>Name</Table.Th>
              <Table.Th>Score</Table.Th>
              <Table.Th className="w-32 [&>div]:justify-end">Actions</Table.Th>
            </tr>
          </thead>

          <tbody>
            {data.scores.map((node) => (
              <tr key={node.user.id}>
                <Table.Td>{node.user.name}</Table.Td>
                <Table.Td className="tabular-nums">
                  {typeof node.score === 'number' ? node.score : <span className="text-neutral-300">-</span>}
                </Table.Td>

                <Table.Td className="text-sm [&>div]:justify-end [&>div]:pr-1.5">
                  <button
                    className="flex h-8 items-center rounded-lg border border-neutral-200 bg-neutral-50 px-3 hover:bg-neutral-100"
                    onClick={() => saveModalRef.current?.onOpen(data.syllabus, node.user, node.score)}
                  >
                    {typeof node.score === 'number' ? 'Edit' : 'Add'}
                  </button>
                </Table.Td>
              </tr>
            ))}
          </tbody>
        </Table>
      </div>

      <SaveModal
        ref={saveModalRef}
        description={({ user, isEditing }) => (
          <>
            You are {isEditing ? 'editing the score' : 'adding a score'} for <b>{user.name}</b>
          </>
        )}
      />
    </>
  );
}

SyllabusScores.loader = async ({ params }: LoaderFunctionArgs) => {
  const [syllabus, scores] = await Promise.all([
    (async () => {
      try {
        const res = await fetch(`${import.meta.env.VITE_API_URL}/v1/syllabuses/${params.syllabusId}`);
        return (await res.json()).syllabus;
      } catch {
        return null;
      }
    })(),
    (async () => {
      try {
        const res = await fetch(`${import.meta.env.VITE_API_URL}/v1/syllabuses/${params.syllabusId}/scores`);
        return (await res.json()).nodes;
      } catch {
        return [];
      }
    })(),
  ]);
  if (!syllabus?.isAssignment) return redirect('/syllabuses');
  return { syllabus, scores };
};

SyllabusScores.action = async ({ request }: ActionFunctionArgs) => {
  try {
    const { type, _syllabusId, _userId, ...body } = await request.json();

    let res;
    switch (type) {
      case 'SAVE':
        res = await fetch(`${import.meta.env.VITE_API_URL}/v1/syllabuses/${_syllabusId}/scores/${_userId}`, {
          method: 'PUT',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify(body),
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

export default SyllabusScores;
