import { AlertCircleIcon } from 'lucide-react';

import type { ErrorResponse } from '@/api-client';
import { Alert, AlertDescription, AlertTitle } from '@/components/ui/alert';

export default function ErrorMessage({
  error,
}: {
  error: Error | ErrorResponse;
}) {
  if (error instanceof Error) {
    return <ErrorMessageItem error={error} />;
  }

  return (
    <ErrorMessageItem error={{ name: error.error, message: error.message }} />
  );
}

function ErrorMessageItem({
  error,
}: {
  error: { name: string; message: string };
}) {
  return (
    <div className="flex flex-col w-full h-full items-center justify-center">
      <Alert variant="destructive" className="max-w-lg">
        <AlertCircleIcon />
        <AlertTitle>{error.name}</AlertTitle>
        <AlertDescription>{error.message}</AlertDescription>
      </Alert>
    </div>
  );
}
