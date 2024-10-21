import { forwardRef, useImperativeHandle, useState, useEffect } from 'react';
import { useFetcher } from 'react-router-dom';
import { useForm } from 'react-hook-form';
import * as v from 'valibot';
import { valibotResolver } from '@hookform/resolvers/valibot';

import Button from 'components/Button';
import Modal from 'components/Modal';
import Input from 'components/Input';

import type { SyllabusStructure } from 'types';

const schema = v.object({
  title: v.pipe(v.string(), v.nonEmpty('Title is required.'), v.trim()),
});

export type SaveStructureModalHandle = {
  onOpen(prev: SyllabusStructure | null, structure: SyllabusStructure | null): void;
};

type Data = {
  prev: SyllabusStructure | null;
  structure: SyllabusStructure | null;
};

export default forwardRef<SaveStructureModalHandle>(function SaveStructureModal(_, ref) {
  const fetcher = useFetcher<{ success: boolean; error: { code: string } | null }>();

  const [data, setData] = useState<Data | null>(null);

  const { register, handleSubmit, formState, setError, setValue, reset } = useForm({
    resolver: valibotResolver(schema),
    defaultValues: { title: '' },
  });

  useImperativeHandle(ref, () => ({
    onOpen(prev, structure) {
      setData({ prev, structure });
      if (structure) setValue('title', structure.title);
    },
  }));

  useEffect(() => {
    if (!fetcher.data) return;
    if (fetcher.data.success) {
      setData(null);
      reset();
    } else {
      if (fetcher.data.error?.code === 'TITLE_SHOULD_BE_UNIQUE') {
        setError('title', { message: 'Title should be unique.' });
      }
      alert(fetcher.data.error?.code || 'INTERNAL_SERVER_ERROR');
    }
  }, [fetcher.data]);

  return (
    <Modal
      title={`${data?.structure ? 'Edit' : 'Add'} Structure`}
      description={
        data?.prev && !data.structure ? (
          <>
            The new structure will be added after <b>{data.prev.title}</b>.
          </>
        ) : null
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
            {
              type: 'SAVE_STRUCTURE',
              _structureId: data.structure?.id || null,
              prevId: data.prev?.id || null,
              ...values,
            },
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
