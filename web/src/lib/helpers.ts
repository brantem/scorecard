import { twMerge } from 'tailwind-merge';
import clsx, { type ClassValue } from 'clsx';

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs));
}

export function mergeRefs<T>(...refs: (React.MutableRefObject<T> | React.LegacyRef<T>)[]): React.RefCallback<T> {
  return (v) => {
    refs.forEach((ref) => {
      if (typeof ref === 'function') {
        ref(v);
      } else if (ref !== null) {
        (ref as React.MutableRefObject<T | null>).current = v;
      }
    });
  };
}

export function withMergedRefs<T>(
  { ref, ...props }: React.HTMLProps<T>,
  ...refs: Parameters<typeof mergeRefs>
): React.HTMLProps<T> & { ref: React.Ref<T> } {
  return { ...props, ref: ref ? mergeRefs(ref, ...refs) : mergeRefs(...refs) };
}

export function formatNumber(v: number) {
  return Intl.NumberFormat(undefined, { maximumFractionDigits: 2 }).format(v);
}
