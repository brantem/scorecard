import { forwardRef, useImperativeHandle, useEffect, useState, useRef } from 'react';
import { useFetcher } from 'react-router-dom';
import { useForm } from 'react-hook-form';
import * as v from 'valibot';
import { valibotResolver } from '@hookform/resolvers/valibot';

import Button from 'components/Button';
import Modal from 'components/Modal';
import Input from 'components/Input';

import type { User } from 'types/user';
import type { BaseSyllabus } from 'types/syllabus';
import { withMergedRefs } from 'lib/helpers';

const schema = v.object({
  score: v.pipe(v.number(), v.minValue(0, 'Score must be 0 or higher.'), v.maxValue(100, 'Score must not exceed 100')),
});

export type SaveModalHandle = {
  open(syllabus: BaseSyllabus, user: User, score: number | null): void;
};

type Data = {
  syllabus: BaseSyllabus;
  user: User;
  isEditing: boolean;
};

type SaveModalProps = {
  description(data: Data): React.ReactNode;
};

export default forwardRef<SaveModalHandle, SaveModalProps>(function SaveModal({ description }, ref) {
  const fetcher = useFetcher<{ success: boolean; error: { code: string } | null }>();

  const inputRef = useRef<HTMLInputElement>(null);

  const [data, setData] = useState<Data | null>(null);

  const { register, handleSubmit, formState, setValue, reset } = useForm({
    resolver: valibotResolver(schema),
    defaultValues: { score: 0 },
  });

  useImperativeHandle(ref, () => ({
    open(syllabus, user, score) {
      const isEditing = typeof score === 'number';
      setData({ syllabus, user, isEditing });
      if (isEditing) setValue('score', score || 0);
      setTimeout(() => inputRef.current?.focus(), 0);
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
      title={`${data?.isEditing ? 'Edit' : 'Add'} Score`}
      description={data ? description(data) : null}
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
            { type: 'SAVE', _syllabusId: data.syllabus.id, _userId: data.user.id, ...values },
            { method: 'PUT', encType: 'application/json' },
          );
        })}
      >
        <Input
          label="Score"
          type="number"
          step=".01"
          {...withMergedRefs(register('score', { valueAsNumber: true }), inputRef)}
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
