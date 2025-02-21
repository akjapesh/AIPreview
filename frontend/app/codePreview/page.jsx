"use client";

//Components
import React, { Suspense } from "react";
import { MainLayout } from "../../components/MainLayout";
import { MainLayoutTombstone } from "../../components/tombstones/MainLayoutTombstone";

export default function CodePreviewPage() {
  return (
    <div className="flex h-full w-full">
      <Suspense fallback={<MainLayoutTombstone />}>
        <MainLayout />
      </Suspense>
    </div>
  );
}
