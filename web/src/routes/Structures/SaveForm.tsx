import { useEffect } from 'react';
import { useFetcher } from 'react-router-dom';
import { useForm } from 'react-hook-form';
import * as v from 'valibot';
import { valibotResolver } from '@hookform/resolvers/valibot';

import Input from 'components/Input';
import Button from 'components/Button';

import type { Structure } from 'types/scorecard';
import { withMergedRefs } from 'lib/helpers';

const schema = v.object({
  title: v.pipe(v.string(), v.nonEmpty('Title is required.'), v.trim()),
});

type SaveFormProps = {
  initialFocusRef: React.RefObject<HTMLElement>;
  parentId: number | null;
  structure: Structure | null;
  onCompleted(): void;
};

export default function SaveForm({ initialFocusRef, parentId, structure, onCompleted }: SaveFormProps) {
  const fetcher = useFetcher<{ success: boolean; error: { code: string } | null }>();

  const { register, handleSubmit, formState } = useForm({
    resolver: valibotResolver(schema),
    defaultValues: { title: structure?.title || '' },
  });

  useEffect(() => {
    if (!fetcher.data) return;
    if (fetcher.data.success) {
      onCompleted();
    } else {
      alert(fetcher.data.error?.code || 'INTERNAL_SERVER_ERROR');
    }
  }, [fetcher.data]);

  return (
    <form
      className="mt-4"
      onSubmit={handleSubmit((values) => {
        fetcher.submit(
          { type: 'SAVE', _structureId: structure?.id || null, parentId, ...values },
          { method: 'PUT', encType: 'application/json' },
        );
      })}
    >
      <Input
        label="Title"
        {...withMergedRefs(register('title'), initialFocusRef)}
        error={formState.errors.title?.message}
        required
        autoFocus
      />

      <Button type="submit" className="mt-4 h-12 w-full" disabled={!formState.isValid}>
        Save
      </Button>
    </form>
  );
}
