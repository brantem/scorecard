import { forwardRef, useImperativeHandle, useState, useEffect } from 'react';
import { useFetcher } from 'react-router-dom';
import { useForm } from 'react-hook-form';
import * as v from 'valibot';
import { valibotResolver } from '@hookform/resolvers/valibot';

import Button from 'components/Button';
import Modal from 'components/Modal';
import Input from 'components/Input';

import type { Structure, Syllabus } from './types';

const schema = v.object({
  title: v.pipe(v.string(), v.nonEmpty('Title is required.'), v.trim()),
});

export type SaveSyllabusModalHandle = {
  onOpen(structure: Structure, parent: Syllabus | null): void;
};

export default forwardRef<SaveSyllabusModalHandle>(function SaveSyllabusModal(_, ref) {
  const fetcher = useFetcher<{ success: boolean; error: { code: string } | null }>();

  const [data, setData] = useState<{ structure: Structure; parent: Syllabus | null } | null>(null);

  useImperativeHandle(ref, () => ({
    onOpen(structure, parent) {
      setData({ structure, parent });
    },
  }));

  const { register, handleSubmit, formState, setError } = useForm({
    resolver: valibotResolver(schema),
    defaultValues: { title: '' },
  });

  useEffect(() => {
    if (!fetcher.data) return;
    if (fetcher.data.success) {
      setData(null);
    } else {
      if (fetcher.data.error?.code === 'TITLE_SHOULD_BE_UNIQUE')
        setError('title', { message: 'Title should be unique.' });
      alert(fetcher.data.error?.code);
    }
  }, [fetcher.data]);

  return (
    <Modal
      title={`Add ${data?.structure.title}`}
      description={
        data?.structure.prevId && data.parent ? (
          <>
            This will be added under <b>{data.parent.title}</b>
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
              type: 'SAVE_SYLLABUS',
              parentId: data.parent?.id || null,
              structureId: data.structure.id || null,
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
