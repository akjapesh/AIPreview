// components/CodePane/CodeEditor.tsx
"use client";
import {
  SandpackCodeEditor,
  SandpackFileExplorer,
} from "@codesandbox/sandpack-react";
import React from "react";
import { useStyletron } from "baseui";
import MonacoEditor from "@monaco-editor/react";

export function CodeEditor({ code }) {
  const [css] = useStyletron();

  return (
    <div className="flex h-full w-full">
      <div className="flex h-full">
        <SandpackFileExplorer />
      </div>
      <div className="w-0.5 bg-gray-200"></div>
      <div className="flex flex-1 overflow-auto h-full">
        <SandpackCodeEditor
          showLineNumbers
          showInlineErrors
          wrapContent
          showTabs
          closableTabs
        />
      </div>
    </div>
  );
}
