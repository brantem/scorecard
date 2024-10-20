import { useEffect } from 'react';
import { useFetcher } from 'react-router-dom';
import { useForm } from 'react-hook-form';
import * as v from 'valibot';
import { valibotResolver } from '@hookform/resolvers/valibot';

import Input from 'components/Input';
import Button from 'components/Button';

const schema = v.object({
  title: v.pipe(v.string(), v.nonEmpty('Title is required.'), v.trim()),
});

type SaveStructureFormProps = {
  parentId: number | null;
  onCompleted(): void;
};

export default function SaveStructureForm({ parentId, onCompleted }: SaveStructureFormProps) {
  const fetcher = useFetcher<{ success: boolean; error: { code: string } | null }>();

  const { register, handleSubmit, formState, setError, reset } = useForm({
    resolver: valibotResolver(schema),
    defaultValues: { title: '' },
  });

  useEffect(() => {
    if (!fetcher.data) return;
    if (fetcher.data.success) {
      onCompleted();
      reset();
    } else {
      if (fetcher.data.error?.code === 'TITLE_SHOULD_BE_UNIQUE') {
        setError('title', { message: 'Title should be unique.' });
      }
      alert(fetcher.data.error?.code);
    }
  }, [fetcher.data]);

  return (
    <form
      className="mt-4"
      onSubmit={handleSubmit((values) => {
        fetcher.submit({ type: 'SAVE', parentId, ...values }, { method: 'PUT', encType: 'application/json' });
      })}
    >
      <Input label="Title" {...register('title')} error={formState.errors.title?.message} required />

      <Button type="submit" className="mt-4 h-12 w-full" disabled={!formState.isValid}>
        Save
      </Button>
    </form>
  );
}
