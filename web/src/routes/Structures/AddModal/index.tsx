import { forwardRef, useImperativeHandle, useState } from 'react';
import { PencilSquareIcon, DocumentDuplicateIcon } from '@heroicons/react/20/solid';

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
          <SaveForm
            parentId={data?.parent?.id || null}
            structure={null}
            onCompleted={() => {
              setIsManual(null);
              setData(null);
            }}
          />
        ) : (
          <SyllabusList
            parentId={data?.parent?.id || null}
            onCompleted={() => {
              setIsManual(null);
              setData(null);
            }}
          />
        )
      ) : (
        <div className="mt-4 flex w-full items-center gap-4">
          <Button
            className="flex flex-1 items-center justify-center gap-2 rounded-lg border border-neutral-900 text-sm"
            onClick={() => setIsManual(true)}
          >
            <PencilSquareIcon className="size-5" />
            <span>Manual</span>
          </Button>
          <span>Or</span>
          <Button
            className="flex flex-1 items-center justify-center gap-2 rounded-lg border border-neutral-900 text-sm"
            onClick={() => setIsManual(false)}
          >
            <DocumentDuplicateIcon className="size-5" />
            <span>From Syllabus</span>
          </Button>
        </div>
      )}
    </Modal>
  );
});
