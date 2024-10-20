import { cn } from 'lib/helpers';

export default function Button({ className, ...props }: React.ComponentPropsWithoutRef<'button'>) {
  return (
    <button
      className={cn(
        'flex items-center gap-2 rounded-lg bg-neutral-900 px-4 py-2 font-medium text-white hover:bg-neutral-800 disabled:bg-neutral-200 disabled:text-neutral-400',
        className,
      )}
      {...props}
    />
  );
}
