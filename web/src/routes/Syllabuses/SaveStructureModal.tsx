import { forwardRef, useImperativeHandle, useState, useEffect } from 'react';
import { useFetcher } from 'react-router-dom';
import { useForm } from 'react-hook-form';
import * as v from 'valibot';
import { valibotResolver } from '@hookform/resolvers/valibot';

import Button from 'components/Button';
import Modal from 'components/Modal';
import Input from 'components/Input';

import type { Structure } from './types';

const schema = v.object({
  title: v.pipe(v.string(), v.nonEmpty('Title is required.'), v.trim()),
});

export type SaveStructureModalHandle = {
  onOpen(prev: Structure | null): void;
};

export default forwardRef<SaveStructureModalHandle>(function SaveStructureModal(_, ref) {
  const fetcher = useFetcher<{ success: boolean; error: { code: string } | null }>();

  const [prev, setPrev] = useState<Structure | null>(null);

  useImperativeHandle(ref, () => ({
    onOpen(prev) {
      setPrev(prev);
    },
  }));

  const { register, handleSubmit, formState, setError, reset } = useForm({
    resolver: valibotResolver(schema),
    defaultValues: { title: '' },
  });

  useEffect(() => {
    if (!fetcher.data) return;
    if (fetcher.data.success) {
      setPrev(null);
      reset();
    } else {
      if (fetcher.data.error?.code === 'TITLE_SHOULD_BE_UNIQUE') {
        setError('title', { message: 'Title should be unique.' });
      }
      alert(fetcher.data.error?.code);
    }
  }, [fetcher.data]);

  return (
    <Modal
      title="Add Structure"
      description={
        prev ? (
          <>
            The new structure will be added after <b>{prev.title}</b>
          </>
        ) : null
      }
      className="max-w-md"
      isOpen={!!prev}
      onClose={() => setPrev(null)}
    >
      <form
        className="mt-4"
        onSubmit={handleSubmit((values) => {
          fetcher.submit(
            {
              type: 'SAVE_STRUCTURE',
              prevId: prev ? prev.id : null,
              ...values,
            },
            { method: 'PUT', encType: 'application/json' },
          );
        })}
      >
        <Input label="Title" {...register('title')} error={formState.errors.title?.message} required />

        <Button type="submit" className="mt-4 h-12 w-full" disabled={!formState.isValid}>
          Save
        </Button>
      </form>
    </Modal>
  );
});
