import React, { Suspense, useState } from "react";
import { useStyletron } from "baseui";
import { Tabs, Tab } from "baseui/tabs-motion";
import dynamic from "next/dynamic";
import { PreviewTombstone } from "../tombstones/PreviewTombstone";
import { EditorTombstone } from "../tombstones/EditorTombstone";

// Dynamically import Monaco Editor to avoid SSR issues
const CodeEditor = dynamic(
  () => import("./CodeEditor").then((mod) => mod.CodeEditor),
  { ssr: false }
);

const Preview = dynamic(() => import("./Preview").then((mod) => mod.Preview), {
  ssr: false,
});

export function CodePreviewBox() {
  const [css] = useStyletron();
  const [activeTab, setActiveTab] = useState("0");
  //   const { currentCode } = useCodeGeneration();
  const { currentCode } = {};

  return (
    <div
      className={css({
        height: "100%",
        padding: "16px",
        boxSizing: "border-box",
      })}
    >
      <div
        style={{
          height: "100%",
        }}
      >
        <Tabs
          activeKey={activeTab}
          onChange={({ activeKey }) => setActiveTab(String(activeKey))}
          overrides={{
            Root: {
              style: {
                height: "100%",
                display: "flex",
                flexDirection: "column",
              },
            },
            TabContent: {
              style: {
                flexGrow: 1,
                padding: "16px",
              },
            },
          }}
        >
          <Tab title="Code">
            <Suspense fallback={<EditorTombstone />}>
              <CodeEditor code={currentCode} />
            </Suspense>
          </Tab>
          <Tab title="Preview">
            <Suspense fallback={<PreviewTombstone />}>
              <Preview code={currentCode} />
            </Suspense>
          </Tab>
        </Tabs>
      </div>
    </div>
  );
}
