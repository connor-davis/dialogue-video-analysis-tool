export default function PageLoader({ label }: { label?: string }) {
  return (
    <div className="flex w-full h-full items-center justify-center">
      <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-primary"></div>
      <span className="mt-4 text-sm text-gray-600">
        {label ?? 'Loading...'}
      </span>
    </div>
  );
}
