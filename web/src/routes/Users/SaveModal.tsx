import { forwardRef, useImperativeHandle, useEffect, useState } from 'react';
import { useFetcher } from 'react-router-dom';
import { useForm } from 'react-hook-form';
import * as v from 'valibot';
import { valibotResolver } from '@hookform/resolvers/valibot';

import Button from 'components/Button';
import Modal from 'components/Modal';
import Input from 'components/Input';

const schema = v.object({
  name: v.pipe(v.string(), v.nonEmpty('Name is required.'), v.trim()),
});

export type SaveModalHandle = {
  onOpen(): void;
};

export default forwardRef<SaveModalHandle>(function SaveModal(_, ref) {
  const fetcher = useFetcher<{ success: boolean; error: { code: string } | null }>();

  const [isOpen, setIsOpen] = useState(false);

  useImperativeHandle(ref, () => ({
    onOpen() {
      setIsOpen(true);
    },
  }));

  const { register, handleSubmit, formState, setError } = useForm({
    resolver: valibotResolver(schema),
    defaultValues: { name: '' },
  });

  useEffect(() => {
    if (!fetcher.data) return;
    if (fetcher.data.success) {
      setIsOpen(false);
    } else {
      if (fetcher.data.error?.code === 'NAME_SHOULD_BE_UNIQUE') setError('name', { message: 'Name should be unique.' });
      alert(fetcher.data.error?.code);
    }
  }, [fetcher.data]);

  return (
    <Modal title="Add User" className="max-w-md" isOpen={isOpen} onClose={() => setIsOpen(false)}>
      <form
        className="mt-4"
        onSubmit={handleSubmit((values) => {
          fetcher.submit(values, { method: 'PUT', encType: 'application/json' });
        })}
      >
        <Input label="Name" {...register('name')} error={formState.errors.name?.message} required />

        <Button type="submit" className="mt-4 h-12 w-full" disabled={!formState.isValid}>
          Save
        </Button>
      </form>
    </Modal>
  );
});
