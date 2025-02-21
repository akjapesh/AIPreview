import { Divider } from "./Divider";

export default function Fallback({ error, resetErrorBoundary }) {
  return (
    <div className="flex items-center justify-center w-full h-full">
      <div
        className="flex flex-col gap-6 w-fit my-12  items-center justify-center py-4 border solid-1 bg-gray-200 border-black rounded-xl"
        style={{ minWidth: "600px" }}
      >
        <div className="w-full py-2 px-4 overflow-scroll">
          <h1 className="bold text-2xl">ERROR</h1>
          <Divider />
          <p>Something went wrong:</p>
          <p>{error.message}</p>
        </div>
        <button
          className="mx-10 px-2 py-2 w-fit border solid-1 border-black rounded-xl bg-white hover:bg-black hover:text-white "
          onClick={resetErrorBoundary}
        >
          Try Again
        </button>
      </div>
    </div>
  );
}
