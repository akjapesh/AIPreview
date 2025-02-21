import dynamic from "next/dynamic";
import { useState } from "react";
import { useStyletron } from "baseui";
import { CodePreviewBox } from "./codePreviewBox";
import { PromptBox } from "./promptBox";
import { ErrorBoundary } from "react-error-boundary";
import Fallback from "./infra/ErrorBoundaryComponent";

// Use dynamic imports for client components that depend on window
const ResizableDivider = dynamic(
  () => import("./ResizableDivider").then((mod) => mod.ResizableDivider),
  {
    ssr: false,
  }
);

export const MainLayout = () => {
  const [css] = useStyletron();
  const [leftPaneWidth, setLeftPaneWidth] = useState("40%");

  return (
    <div className="flex h-hull w-full overflow-hidden">
      <div
        className={css({
          width: leftPaneWidth,
          height: "100%",
          transition: "width 0.1s ease",
        })}
      >
        <ErrorBoundary FallbackComponent={Fallback}>
          <PromptBox />
        </ErrorBoundary>
      </div>

      <ResizableDivider onResize={setLeftPaneWidth} />

      <div
        className={css({
          width: `calc(100% - ${leftPaneWidth})`,
          height: "100%",
          transition: "width 0.1s ease",
        })}
      >
        <ErrorBoundary FallbackComponent={Fallback}>
          <CodePreviewBox />
        </ErrorBoundary>
      </div>
    </div>
  );
};
