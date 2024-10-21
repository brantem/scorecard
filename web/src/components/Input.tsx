import { cn } from 'lib/helpers';
import { forwardRef } from 'react';

type InputProps = React.ComponentPropsWithoutRef<'input'> & {
  label?: string;
  error?: string | undefined;
};

export default forwardRef<HTMLInputElement, InputProps>(function Input({ className, label, error, ...props }, ref) {
  return (
    <label className="flex w-full flex-col gap-1">
      <span className="font-medium text-neutral-900">
        {label}
        {props.required && <span className="ml-1 text-red-500">*</span>}
      </span>
      <div className="flex w-full items-center">
        <input
          ref={ref}
          type="text"
          {...props}
          className={cn(
            'h-12 flex-1 rounded-lg border px-3',
            error ? 'border-red-500' : 'border-neutral-200',
            className,
          )}
        />
      </div>
      {error ? <span className="text-sm text-red-500">{error}</span> : null}
    </label>
  );
});
