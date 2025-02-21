"use client";

import { Divider } from "../../components/infra/Divider";

export default function ErrorFallback({ error, reset }) {
  return (
    <div
      className="flex flex-col gap-6 items-center justify-center py-4 border solid-1 bg-gray-200 border-black rounded-xl"
      style={{ width: "600px", marginLeft: "30vw", marginTop: "25vh" }}
    >
      <div className="py-2 px-4 overflow-scroll">
        <h1 className="bold text-2xl">ERROR</h1>
        <Divider />
        <p>Something went wrong:</p>
        <p>{error.message}</p>
      </div>
      <button
        className="px-2 py-2 border solid-1 border-black rounded-xl bg-white hover:bg-gray-200 "
        onClick={reset}
      >
        Try Again
      </button>
    </div>
  );
}
