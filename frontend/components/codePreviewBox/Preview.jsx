import React, { useState } from "react";
import { useStyletron } from "baseui";
import {
  SandpackProvider,
  SandpackLayout,
  SandpackCodeEditor,
  SandpackPreview,
  SandpackFileExplorer,
  SandpackConsole,
} from "@codesandbox/sandpack-react";

export function Preview({ code }) {
  const [css] = useStyletron();
  const [previewUrl, setPreviewUrl] = useState("");

  return (
    <div className="flex flex-col w-full h-full">
      <div className="h-full flex">
        <SandpackPreview />
      </div>
      
      {/* to be added later if needed the console
      <div className="bg-gray-200 h-1 w-full"></div>
      <div className="bg-white sticky bottom-0" style={{ height: "200px" }}>
        <div className="p-2 text-gray-500 sticky top-0">Console:</div>
        <div className="bg-black w-full" style={{ height: "1px" }}></div>
        <div className="flex overflow-scroll">
          <SandpackConsole />
        </div>
      </div> */}
    </div>
  );
}
