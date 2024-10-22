import { forwardRef, useImperativeHandle, useEffect, useState } from 'react';
import { useFetcher } from 'react-router-dom';
import { useForm } from 'react-hook-form';
import * as v from 'valibot';
import { valibotResolver } from '@hookform/resolvers/valibot';

import Button from 'components/Button';
import Modal from 'components/Modal';
import Input from 'components/Input';

import type { Score } from 'types/syllabus';

const schema = v.object({
  score: v.pipe(v.number(), v.minValue(0, 'Score must be 0 or higher.'), v.maxValue(100, 'Score must not exceed 100')),
});

export type SaveModalHandle = {
  onOpen(node: Score): void;
};

type SaveModalProps = {
  syllabusId: number;
};

export default forwardRef<SaveModalHandle, SaveModalProps>(function SaveModal({ syllabusId }, ref) {
  const fetcher = useFetcher<{ success: boolean; error: { code: string } | null }>();

  const [data, setData] = useState<{ node: Score } | null>(null);

  const { register, handleSubmit, formState, setValue, reset } = useForm({
    resolver: valibotResolver(schema),
    defaultValues: { score: 0 },
  });

  useImperativeHandle(ref, () => ({
    onOpen(node) {
      setData({ node });
      if (node) setValue('score', node.score || 0);
    },
  }));

  useEffect(() => {
    if (!fetcher.data) return;
    if (fetcher.data.success) {
      setData(null);
      reset();
    } else {
      alert(fetcher.data.error?.code || 'INTERNAL_SERVER_ERROR');
    }
  }, [fetcher.data]);

  return (
    <Modal
      title={`${data?.node ? 'Edit' : 'Add'} Score`}
      description={
        <>
          You are {data?.node ? 'editing the score' : 'adding a score'} for <b>{data?.node.user.name}</b>
        </>
      }
      isOpen={!!data}
      onClose={() => {
        setData(null);
        reset();
      }}
    >
      <form
        className="mt-4"
        onSubmit={handleSubmit((values) => {
          if (!data) return;
          fetcher.submit(
            { type: 'SAVE', _syllabusId: syllabusId, _userId: data.node.user.id, ...values },
            { method: 'PUT', encType: 'application/json' },
          );
        })}
      >
        <Input
          label="Score"
          type="number"
          step=".01"
          {...register('score', { valueAsNumber: true })}
          error={formState.errors.score?.message}
          required
          autoFocus
        />

        <Button type="submit" className="mt-4 h-12 w-full" disabled={!formState.isValid}>
          Save
        </Button>
      </form>
    </Modal>
  );
});
