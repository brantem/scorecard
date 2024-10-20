import { Dialog, DialogBackdrop, DialogPanel, DialogTitle, Description } from '@headlessui/react';
import { cn } from 'lib/helpers';

type ModalProps = {
  isOpen: boolean;
  onClose(): void;
  className?: string;
  title: React.ReactNode;
  description?: React.ReactNode;
  children: React.ReactNode;
};

export default function Modal({ isOpen, onClose, className, title, description, children }: ModalProps) {
  return (
    <Dialog open={isOpen} onClose={onClose} className="relative z-10 focus:outline-none">
      <DialogBackdrop className="fixed inset-0 bg-black/30" />

      <div className="fixed inset-0 flex w-screen items-start justify-center p-4 pt-32">
        <DialogPanel
          className={cn(
            'w-full max-w-lg rounded-lg border border-neutral-200 bg-white p-6 data-[closed]:opacity-0',
            className,
          )}
        >
          <DialogTitle className="text-lg font-semibold">{title}</DialogTitle>
          {description && <Description className="mt-1">{description}</Description>}
          {children}
        </DialogPanel>
      </div>
    </Dialog>
  );
}
