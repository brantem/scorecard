import { useContext, createContext } from 'react';
import { Slot } from '@radix-ui/react-slot';

import { cn } from 'lib/helpers';

type TreeItem = {
  id: number;
  parentId: number | null;
  title: string;
};

const TreeContext = createContext<Map<number | null, TreeItem[]>>(new Map());

type TreeProps<T extends TreeItem> = {
  items: Map<number | null, T[]>;
  renderAdd?: (parent: T | null) => React.ReactNode;
};

function Tree<T extends TreeItem>({ items, renderAdd }: TreeProps<T>) {
  return (
    <Tree.Provider value={items}>
      <TreeItem item={null} renderAdd={renderAdd} />
    </Tree.Provider>
  );
}

type TreeItemProps<T extends TreeItem> = {
  item: T | null;
  renderAdd?: (parent: T | null) => React.ReactNode;
};

function TreeItem<T extends TreeItem>({ item, renderAdd }: TreeItemProps<T>) {
  const m = useContext(TreeContext);

  return (
    <div className={cn('relative flex flex-col gap-2 text-sm', item?.parentId && 'ml-[calc(theme(spacing.8)+2px)]')}>
      {item ? (
        <>
          <TreeInnerItem>{item.title}</TreeInnerItem>
          <div className={cn('absolute left-4 top-0 w-px bg-neutral-200', renderAdd ? '-bottom-2' : 'bottom-0')} />
        </>
      ) : null}

      {(m.get(item?.id || null) || []).map((item) => (
        <TreeItem key={item.id} item={item as T} renderAdd={renderAdd} />
      ))}

      {renderAdd?.(item || null)}
    </div>
  );
}

Tree.Provider = TreeContext.Provider;

type TreeInnerItemProps = {
  className?: string;
  asChild?: boolean;
  children: React.ReactNode;
};

function TreeInnerItem({ className, asChild, children }: TreeInnerItemProps) {
  const Comp = asChild ? Slot : 'div';
  return (
    <Comp
      className={cn(
        'relative z-10 w-fit min-w-32 rounded-lg border border-neutral-200 bg-neutral-50 px-3 py-1.5 text-center',
        className,
      )}
    >
      {children}
    </Comp>
  );
}
Tree.Item = TreeInnerItem;

export default Tree;
