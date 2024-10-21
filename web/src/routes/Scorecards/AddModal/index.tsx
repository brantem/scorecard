import { forwardRef, useImperativeHandle, useState } from 'react';

import Modal from 'components/Modal';
import SaveForm from '../SaveForm';
import SyllabusList from './SyllabusList';
import Button from 'components/Button';

import type { Structure } from 'types/scorecard';

export type AddModalHandle = {
  onOpen(parent: Structure | null): void;
};

export default forwardRef<AddModalHandle>(function AddModal(_, ref) {
  const [isManual, setIsManual] = useState<boolean | null>(null);
  const [data, setData] = useState<{ parent: Structure | null } | null>(null);

  useImperativeHandle(ref, () => ({
    onOpen(parent) {
      setData({ parent });
    },
  }));

  return (
    <Modal
      title="Add Structure"
      description={
        typeof isManual === 'boolean' ? (
          data?.parent ? (
            <>
              The new structure will be added under <b>{data.parent.title}</b>.
            </>
          ) : null
        ) : (
          'You can create a structure manually or by copying data from a syllabus.'
        )
      }
      className={isManual === false ? 'max-w-4xl' : ''}
      isOpen={!!data}
      onClose={() => {
        setIsManual(null);
        setData(null);
      }}
    >
      {typeof isManual === 'boolean' ? (
        isManual ? (
          <SaveForm parentId={data?.parent?.id || null} structure={null} onCompleted={() => setData(null)} />
        ) : (
          <SyllabusList parentId={data?.parent?.id || null} onCompleted={() => setData(null)} />
        )
      ) : (
        <div className="mt-4 flex w-full items-center gap-4">
          <Button
            className="flex flex-1 items-center justify-center gap-2 rounded-lg border border-neutral-900 text-sm"
            onClick={() => setIsManual(true)}
          >
            <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor" className="size-5">
              <path d="m5.433 13.917 1.262-3.155A4 4 0 0 1 7.58 9.42l6.92-6.918a2.121 2.121 0 0 1 3 3l-6.92 6.918c-.383.383-.84.685-1.343.886l-3.154 1.262a.5.5 0 0 1-.65-.65Z" />
              <path d="M3.5 5.75c0-.69.56-1.25 1.25-1.25H10A.75.75 0 0 0 10 3H4.75A2.75 2.75 0 0 0 2 5.75v9.5A2.75 2.75 0 0 0 4.75 18h9.5A2.75 2.75 0 0 0 17 15.25V10a.75.75 0 0 0-1.5 0v5.25c0 .69-.56 1.25-1.25 1.25h-9.5c-.69 0-1.25-.56-1.25-1.25v-9.5Z" />
            </svg>

            <span>Manual</span>
          </Button>
          <span>Or</span>
          <Button
            className="flex flex-1 items-center justify-center gap-2 rounded-lg border border-neutral-900 text-sm"
            onClick={() => setIsManual(false)}
          >
            <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor" className="size-5">
              <path d="M7 3.5A1.5 1.5 0 0 1 8.5 2h3.879a1.5 1.5 0 0 1 1.06.44l3.122 3.12A1.5 1.5 0 0 1 17 6.622V12.5a1.5 1.5 0 0 1-1.5 1.5h-1v-3.379a3 3 0 0 0-.879-2.121L10.5 5.379A3 3 0 0 0 8.379 4.5H7v-1Z" />
              <path d="M4.5 6A1.5 1.5 0 0 0 3 7.5v9A1.5 1.5 0 0 0 4.5 18h7a1.5 1.5 0 0 0 1.5-1.5v-5.879a1.5 1.5 0 0 0-.44-1.06L9.44 6.439A1.5 1.5 0 0 0 8.378 6H4.5Z" />
            </svg>

            <span>From Syllabus</span>
          </Button>
        </div>
      )}
    </Modal>
  );
});
