import { Provider, Root, Portal, Trigger, Content, TooltipContentProps, Arrow } from '@radix-ui/react-tooltip';

export type TooltipProps = {
  isOpen?: boolean;
  content: React.ReactNode;
  children: React.ReactNode;
} & Partial<Pick<TooltipContentProps, 'side' | 'align'>>;

export default function Tooltip({ children, isOpen, content, side = 'top', align = 'center' }: TooltipProps) {
  return (
    <Provider delayDuration={300}>
      <Root {...(isOpen != undefined && { open: isOpen })}>
        <Trigger asChild>{children}</Trigger>

        <Portal>
          <Content
            className="relative z-[1000] max-w-xs rounded-lg bg-neutral-900 px-3 py-2 text-sm font-normal text-white"
            side={side}
            align={align}
            sideOffset={2}
          >
            {content}
            <Arrow />
          </Content>
        </Portal>
      </Root>
    </Provider>
  );
}
