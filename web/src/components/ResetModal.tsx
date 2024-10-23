import { useState, useImperativeHandle, forwardRef } from 'react';

import Modal from './Modal';
import Button from './Button';

export type ResetModalHandle = {
  open(title: string, values: Record<string, any>): void;
  close(): void;
};

type ResetModalProps = {
  onAccept(values: Record<string, any>): void;
};

export default forwardRef<ResetModalHandle, ResetModalProps>(function ResetModal({ onAccept }, ref) {
  const [data, setData] = useState<{ title: string; values: Record<string, any> } | null>(null);

  useImperativeHandle(ref, () => ({
    open(title, values) {
      setData({ title, values });
    },
    close() {
      setData(null);
    },
  }));

  return (
    <Modal
      title="Are you sure you want to reset?"
      description="This action is irreversible. Deleted data cannot be recovered."
      className="max-w-md"
      isOpen={!!data}
      onClose={() => setData(null)}
    >
      <div className="mt-4 flex items-center justify-end gap-2">
        <Button onClick={() => setData(null)}>Cancel</Button>
        <Button className="bg-red-50 text-red-500 hover:bg-red-100" onClick={() => onAccept(data?.values || {})}>
          Reset {data?.title}
        </Button>
      </div>
    </Modal>
  );
});
