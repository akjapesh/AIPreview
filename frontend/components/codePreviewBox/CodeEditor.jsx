// components/CodePane/CodeEditor.tsx
"use client";

import React from "react";
import { useStyletron } from "baseui";
import MonacoEditor from "@monaco-editor/react";

export function CodeEditor({ code }) {
  const [css] = useStyletron();

  return (
    <div className={css({ height: "100%" })}>
      <MonacoEditor
        height="100%"
        language="typescript"
        theme="vs"
        value={code}
        options={{
          readOnly: true,
          minimap: { enabled: false },
        }}
      />
    </div>
  );
}
