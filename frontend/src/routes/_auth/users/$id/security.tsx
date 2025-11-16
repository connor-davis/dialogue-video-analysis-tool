import { createFileRoute } from '@tanstack/react-router';

export const Route = createFileRoute('/_auth/users/$id/security')({
  component: RouteComponent,
});

function RouteComponent() {
  return <div>Hello "/_auth/users/$id/security"!</div>;
}
