import { cn } from 'lib/helpers';
import { forwardRef } from 'react';

type InputProps = React.ComponentPropsWithoutRef<'input'> & {
  label?: string;
  error?: string | undefined;
};

export default forwardRef<HTMLInputElement, InputProps>(function Input({ label, className, ...props }, ref) {
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
          className={cn('h-12 flex-1 rounded-lg border border-neutral-200 px-3', className)}
        />
      </div>
      {props.error ? <span className="text-sm text-red-500">{props.error}</span> : null}
    </label>
  );
});
