import { forwardRef, useImperativeHandle, useEffect, useState, useRef } from 'react';
import { useFetcher } from 'react-router-dom';
import { useForm } from 'react-hook-form';
import * as v from 'valibot';
import { valibotResolver } from '@hookform/resolvers/valibot';

import Button from 'components/Button';
import Modal from 'components/Modal';
import Input from 'components/Input';

import type { User } from 'types/user';
import { withMergedRefs } from 'lib/helpers';

const schema = v.object({
  name: v.pipe(v.string(), v.nonEmpty('Name is required.'), v.trim()),
});

export type SaveModalHandle = {
  open(user: User | null): void;
};

export default forwardRef<SaveModalHandle>(function SaveModal(_, ref) {
  const fetcher = useFetcher<{ success: boolean; error: { code: string } | null }>();

  const inputRef = useRef<HTMLInputElement>(null);

  const [data, setData] = useState<{ user: User | null } | null>(null);

  const { register, handleSubmit, formState, setError, setValue, reset } = useForm({
    resolver: valibotResolver(schema),
    defaultValues: { name: '' },
  });

  useImperativeHandle(ref, () => ({
    open(user) {
      setData({ user });
      if (user) setValue('name', user.name);
      setTimeout(() => inputRef.current?.focus(), 0);
    },
  }));

  useEffect(() => {
    if (!fetcher.data) return;
    if (fetcher.data.success) {
      setData(null);
      reset();
    } else {
      if (fetcher.data.error?.code === 'NAME_SHOULD_BE_UNIQUE') setError('name', { message: 'Name should be unique.' });
      alert(fetcher.data.error?.code || 'INTERNAL_SERVER_ERROR');
    }
  }, [fetcher.data]);

  return (
    <Modal
      title={`${data?.user ? 'Edit' : 'Add'} User`}
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
            { type: 'SAVE', _userId: data?.user?.id || null, ...values },
            { method: 'PUT', encType: 'application/json' },
          );
        })}
      >
        <Input
          label="Name"
          {...withMergedRefs(register('name'), inputRef)}
          error={formState.errors.name?.message}
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
