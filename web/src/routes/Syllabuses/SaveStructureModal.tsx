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

  const [data, setData] = useState<{ prev: Structure | null } | null>(null);

  useImperativeHandle(ref, () => ({
    onOpen(prev) {
      setData({ prev });
    },
  }));

  const { register, handleSubmit, formState, setError, reset } = useForm({
    resolver: valibotResolver(schema),
    defaultValues: { title: '' },
  });

  useEffect(() => {
    if (!fetcher.data) return;
    if (fetcher.data.success) {
      setData(null);
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
        data?.prev ? (
          <>
            The new structure will be added after <b>{data.prev.title}</b>
          </>
        ) : null
      }
      className="max-w-md"
      isOpen={!!data}
      onClose={() => setData(null)}
    >
      <form
        className="mt-4"
        onSubmit={handleSubmit((values) => {
          if (!data) return;
          fetcher.submit(
            {
              type: 'SAVE_STRUCTURE',
              prevId: data.prev?.id || null,
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
