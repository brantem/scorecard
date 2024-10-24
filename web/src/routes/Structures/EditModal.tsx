import { forwardRef, useRef, useState, useImperativeHandle } from 'react';

import Modal from 'components/Modal';
import SaveForm from './SaveForm';

import type { Structure } from 'types/scorecard';

export type EditModalHandle = {
  open(structure: Structure): void;
};

export default forwardRef<EditModalHandle>(function EditModal(_, ref) {
  const initialFocusRef = useRef<HTMLElement>(null);

  const [data, setData] = useState<{ structure: Structure | null } | null>(null);

  useImperativeHandle(ref, () => ({
    open(structure) {
      setData({ structure });
      setTimeout(() => initialFocusRef.current?.focus(), 0);
    },
  }));

  return (
    <Modal title="Edit Structure" isOpen={!!data} onClose={() => setData(null)}>
      {data?.structure && (
        <SaveForm
          initialFocusRef={initialFocusRef}
          parentId={data.structure.parentId || null}
          structure={data.structure}
          onCompleted={() => setData(null)}
        />
      )}
    </Modal>
  );
});
