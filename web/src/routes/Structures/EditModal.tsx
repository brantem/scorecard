import { forwardRef, useImperativeHandle, useState } from 'react';

import Modal from 'components/Modal';
import SaveForm from './SaveForm';

import type { Structure } from 'types';

export type EditModalHandle = {
  onOpen(structure: Structure): void;
};

export default forwardRef<EditModalHandle>(function EditModal(_, ref) {
  const [data, setData] = useState<{ structure: Structure | null } | null>(null);

  useImperativeHandle(ref, () => ({
    onOpen(structure) {
      setData({ structure });
    },
  }));

  return (
    <Modal title="Edit Structure" isOpen={!!data} onClose={() => setData(null)}>
      {data?.structure && (
        <SaveForm
          parentId={data.structure.parentId || null}
          structure={data.structure}
          onCompleted={() => setData(null)}
        />
      )}
    </Modal>
  );
});
