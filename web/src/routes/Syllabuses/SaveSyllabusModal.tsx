import { forwardRef, useImperativeHandle, useEffect, useState, useRef } from 'react';
import { useFetcher } from 'react-router-dom';
import { useForm } from 'react-hook-form';
import * as v from 'valibot';
import { valibotResolver } from '@hookform/resolvers/valibot';

import Button from 'components/Button';
import Modal from 'components/Modal';
import Input from 'components/Input';

import type { Structure, Syllabus } from 'types/syllabus';
import { withMergedRefs } from 'lib/helpers';

const schema = v.object({
  title: v.pipe(v.string(), v.nonEmpty('Title is required.'), v.trim()),
});

type Data = {
  structure: Structure;
  parent: Syllabus | null;
  syllabus: Syllabus | null;
};

export type SaveSyllabusModalHandle = {
  open(structure: Structure, parent: Syllabus | null, syllabus: Syllabus | null): void;
};

export default forwardRef<SaveSyllabusModalHandle>(function SaveSyllabusModal(_, ref) {
  const fetcher = useFetcher<{ success: boolean; error: { code: string } | null }>();

  const inputRef = useRef<HTMLInputElement>(null);

  const [data, setData] = useState<Data | null>(null);

  const { register, handleSubmit, formState, setError, setValue, reset } = useForm({
    resolver: valibotResolver(schema),
    defaultValues: { title: '' },
  });

  useImperativeHandle(ref, () => ({
    open(structure, parent, syllabus) {
      setData({ structure, parent, syllabus });
      if (syllabus) setValue('title', syllabus.title);
      setTimeout(() => inputRef.current?.focus(), 0);
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
      title={`${data?.syllabus ? 'Edit' : 'Add'} ${data?.structure.title}`}
      description={
        data?.parent && !data.syllabus ? (
          <>
            This will be added under <b>{data.parent.title}</b>.
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
              type: 'SAVE_SYLLABUS',
              _syllabusId: data.syllabus?.id || null,
              parentId: data.syllabus ? data.syllabus.parentId : data.parent?.id || null,
              structureId: data.syllabus ? null : data.structure.id,
              ...values,
            },
            { method: 'PUT', encType: 'application/json' },
          );
        })}
      >
        <Input
          label="Title"
          {...withMergedRefs(register('title'), inputRef)}
          error={formState.errors.title?.message}
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
