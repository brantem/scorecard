import { cn } from 'lib/helpers';

function Table({ children }: { children: React.ReactNode }) {
  return (
    <table
      className={cn(
        'w-full table-auto',
        '[&_tr:has(th)]:bg-neutral-50 [&_tr:has(th)]:text-left [&_tr:has(th)]:text-neutral-900',
        '[&_tr:not(:has(th))]:border-t [&_tr:not(:has(th))]:border-neutral-200 [&_tr:not(:has(th))]:text-neutral-700',
      )}
    >
      {children}
    </table>
  );
}

type ThProps = React.ComponentPropsWithoutRef<'th'>;

function Th({ className, children, ...props }: ThProps) {
  return (
    <th className={cn('h-12 whitespace-nowrap py-0 font-medium', className)} {...props}>
      <div className="flex size-full items-center px-3">{children}</div>
    </th>
  );
}
Table.Th = Th;

type TdProps = React.ComponentPropsWithoutRef<'td'>;

function Td({ className, children, ...props }: TdProps) {
  return (
    <td className={cn('h-12 whitespace-nowrap py-0 text-neutral-700', className)} {...props}>
      <div className="flex h-full items-center px-3">{children}</div>
    </td>
  );
}
Table.Td = Td;

export default Table;
