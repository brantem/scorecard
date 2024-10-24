import { createContext, useContext, useState } from 'react';
import {
  ChevronDoubleLeftIcon,
  ChevronDoubleRightIcon,
  ChevronLeftIcon,
  ChevronRightIcon,
} from '@heroicons/react/20/solid';

import { cn } from 'lib/helpers';

type State = {
  isLoading: boolean;

  totalItems: number;
  setTotalItems(totalItems: number): void;

  itemsPerPage: number;
  setItemsPerPage(itemsPerPage: number): void;

  page: number;
  setPage(page: number): void;
};

export type TableState = Pick<State, 'itemsPerPage' | 'page'>;

const TableContext = createContext<State>({
  isLoading: false,

  totalItems: 0,
  setTotalItems() {},

  itemsPerPage: 10,
  setItemsPerPage() {},

  page: 1,
  setPage() {},
});

export const useTable = () => useContext(TableContext);

function Table({ children }: { children: React.ReactNode }) {
  return (
    <table
      className={cn(
        'w-full table-auto',
        '[&_tr:has(th)]:bg-neutral-50 [&_tr:has(th)]:text-left',
        '[&_tr:not(:has(th))]:border-t [&_tr:not(:has(th))]:border-neutral-200',
      )}
    >
      {children}
    </table>
  );
}

type RootProps = {
  totalItems?: number;
  stickyId?: string;
  initialState?: Partial<TableState>;
  onStateChange(data: TableState): void;
  children: React.ReactNode;
  isLoading?: boolean;
};

const Provider = ({
  totalItems: _totalItems = 0,
  initialState,
  onStateChange,
  children,
  isLoading = false,
}: RootProps) => {
  const [totalItems, setTotalItems] = useState(() => _totalItems);
  const _itemsPerPage = initialState?.itemsPerPage || 10;
  const [itemsPerPage, setItemsPerPage] = useState(() => _itemsPerPage);
  const [page, setPage] = useState(() => initialState?.page || 1);

  return (
    <TableContext.Provider
      value={{
        isLoading,

        totalItems,
        setTotalItems(v) {
          setTotalItems(v);
          if (v <= (page - 1) * itemsPerPage && page !== 1) setPage(1);
        },

        itemsPerPage,
        setItemsPerPage(v) {
          setItemsPerPage(v);
          setPage(1);
          onStateChange({ itemsPerPage: v, page });
        },

        page,
        setPage(v) {
          setPage(v);
          onStateChange({ itemsPerPage, page: v });
        },
      }}
    >
      {children}
    </TableContext.Provider>
  );
};
Table.Provider = Provider;

const Cell = (props: React.ComponentPropsWithoutRef<'div'>) => {
  return <div {...props} className="flex h-full items-center px-3" />;
};

type ThProps = React.ComponentPropsWithoutRef<'th'>;

function Th({ className, children, ...props }: ThProps) {
  return (
    <th className={cn('h-12 font-medium', className)} {...props}>
      <Cell>{children}</Cell>
    </th>
  );
}
Table.Th = Th;

type TdProps = React.ComponentPropsWithoutRef<'td'>;

function Td({ className, children, ...props }: TdProps) {
  return (
    <td className={cn('h-12', className)} {...props}>
      <Cell>{children}</Cell>
    </td>
  );
}
Table.Td = Td;

const Stats = () => {
  const { isLoading, totalItems, itemsPerPage, page } = useTable();

  return isLoading ? (
    <span className="select-none">Loading...</span>
  ) : totalItems ? (
    <span className="select-none tabular-nums text-neutral-500">
      <span className="font-semibold text-black">{(page - 1) * itemsPerPage + 1}</span> to{' '}
      <span className="font-semibold text-black">
        {((v) => (totalItems < v ? totalItems : v))(page * itemsPerPage)}
      </span>{' '}
      of <span className="font-semibold text-black">{totalItems}</span> items
    </span>
  ) : (
    <span className="select-none tabular-nums text-neutral-500">
      <span className="font-semibold text-black">0</span> items
    </span>
  );
};
Table.Stats = Stats;

type PaginationProps = {
  className?: string;
};

const Pagination = ({ className }: PaginationProps) => {
  const { totalItems, itemsPerPage, page, setPage } = useTable();

  const totalPages = Math.ceil(totalItems / itemsPerPage);
  const canPrev = page > 1;
  const canNext = page < totalPages;

  return (
    <div className={cn('flex select-none font-bold', className)}>
      <button
        disabled={!canPrev}
        className="px-2 disabled:cursor-not-allowed disabled:text-neutral-300"
        onClick={() => setPage(1)}
      >
        <ChevronDoubleLeftIcon className="size-5" />
      </button>

      <button
        disabled={!canPrev}
        className="px-2 disabled:cursor-not-allowed disabled:text-neutral-300"
        onClick={() => setPage(page - 1)}
      >
        <ChevronLeftIcon className="size-5" />
      </button>

      {[...new Array(Math.min(2, page - 1))].map((_, i) => {
        const v = Math.max(1, page - 1 + i - 1);
        return (
          <button key={i} className="w-9 font-medium text-neutral-400" onClick={() => setPage(v)}>
            {v}
          </button>
        );
      })}

      <button className="w-9 font-medium" disabled>
        {page}
      </button>

      {[...new Array(Math.min(2, Math.abs(totalPages - (page - 1) - 1)))].map((_, i) => {
        const v = page + i + 1;
        return (
          <button key={i} className="w-9 font-medium text-neutral-400" onClick={() => setPage(v)}>
            {v}
          </button>
        );
      })}

      <button
        disabled={!canNext}
        className="px-2 disabled:cursor-not-allowed disabled:text-neutral-300"
        onClick={() => setPage(page + 1)}
      >
        <ChevronRightIcon className="size-5" />
      </button>

      <button
        disabled={!canNext}
        className="px-2 disabled:cursor-not-allowed disabled:text-neutral-300"
        onClick={() => setPage(totalPages)}
      >
        <ChevronDoubleRightIcon className="size-5" />
      </button>
    </div>
  );
};
Table.Pagination = Pagination;

export default Table;
