import { createFileRoute } from '@tanstack/react-router';

export const Route = createFileRoute('/_auth/roles/$id/permissions')({
  component: RouteComponent,
});

function RouteComponent() {
  return <div>Hello "/_auth/roles/$id/permissions"!</div>;
}
