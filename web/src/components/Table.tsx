import { cn } from 'lib/helpers';

function Table({ children }: { children: React.ReactNode }) {
  return (
    <table
      className={cn(
        'w-full table-auto',
        '[&_tr:has(th)]:bg-mono-50 [&_tr:has(th)]:text-left [&_tr:has(th)]:text-neutral-900',
        '[&_tr:not(:has(th))]:border-mono-200 [&_tr:not(:has(th))]:border-t [&_tr:not(:has(th))]:text-neutral-700',
      )}
    >
      {children}
    </table>
  );
}

type ThProps = {
  className?: string;
  children: React.ReactNode;
};

function Th({ className, children }: ThProps) {
  return (
    <th className={cn('h-12 whitespace-nowrap py-0 font-medium', className)}>
      <div className="flex h-full w-full items-center px-3">{children}</div>
    </th>
  );
}
Table.Th = Th;

type TdProps = {
  className?: string;
  children: React.ReactNode;
};

function Td({ className, children }: TdProps) {
  return (
    <td className={cn('text-mono-700 h-12 whitespace-nowrap py-0', className)}>
      <div className="flex h-full items-center px-3">{children}</div>
    </td>
  );
}
Table.Td = Td;

export default Table;
