import { forwardRef, useImperativeHandle, useEffect, useState } from 'react';
import { useFetcher } from 'react-router-dom';
import { useForm } from 'react-hook-form';
import * as v from 'valibot';
import { valibotResolver } from '@hookform/resolvers/valibot';

import Button from 'components/Button';
import Modal from 'components/Modal';
import Input from 'components/Input';

import type { Program } from 'types/program';

const schema = v.object({
  title: v.pipe(v.string(), v.nonEmpty('Title is required.'), v.trim()),
});

export type SaveModalHandle = {
  open(program: Program | null): void;
};

export default forwardRef<SaveModalHandle>(function SaveModal(_, ref) {
  const fetcher = useFetcher<{ success: boolean; error: { code: string } | null }>();

  const [data, setData] = useState<{ program: Program | null } | null>(null);

  const { register, handleSubmit, formState, setError, setValue, reset } = useForm({
    resolver: valibotResolver(schema),
    defaultValues: { title: '' },
  });

  useImperativeHandle(ref, () => ({
    open(program) {
      setData({ program });
      if (program) setValue('title', program.title);
    },
  }));

  useEffect(() => {
    if (!fetcher.data) return;
    if (fetcher.data.success) {
      setData(null);
      reset();
    } else {
      if (fetcher.data.error?.code === 'TITLE_SHOULD_BE_UNIQUE')
        setError('title', { message: 'Title should be unique.' });
      alert(fetcher.data.error?.code || 'INTERNAL_SERVER_ERROR');
    }
  }, [fetcher.data]);

  return (
    <Modal
      title={`${data?.program ? 'Edit' : 'Add'} Program`}
      isOpen={!!data}
      onClose={() => {
        setData(null);
        reset();
      }}
    >
      <form
        className="mt-4"
        onSubmit={handleSubmit((values) => {
          fetcher.submit(
            { type: 'SAVE', _programId: data?.program?.id || null, ...values },
            { method: 'PUT', encType: 'application/json' },
          );
        })}
      >
        <Input label="Title" {...register('title')} error={formState.errors.title?.message} required autoFocus />

        <Button type="submit" className="mt-4 h-12 w-full" disabled={!formState.isValid}>
          Save
        </Button>
      </form>
    </Modal>
  );
});
