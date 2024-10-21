import { Slot } from '@radix-ui/react-slot';

import { cn } from 'lib/helpers';

type ButtonProps = React.ComponentPropsWithoutRef<'button'> & {
  asChild?: boolean;
};

export default function Button({ className, asChild, ...props }: ButtonProps) {
  const Comp = asChild ? Slot : 'button';
  return (
    <Comp
      className={cn(
        'flex items-center justify-center gap-2 rounded-lg bg-neutral-900 px-4 py-2 text-center font-medium text-white hover:bg-neutral-800 disabled:bg-neutral-200 disabled:text-neutral-400',
        className,
      )}
      {...props}
    />
  );
}
