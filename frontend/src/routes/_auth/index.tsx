import { createFileRoute } from '@tanstack/react-router';

export const Route = createFileRoute('/_auth/')({
  component: App,
});

function App() {
  return (
    <div className="relative flex flex-col w-full h-full overflow-hidden bg-background rounded-2xl border"></div>
  );
}
