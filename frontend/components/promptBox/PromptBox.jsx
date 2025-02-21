import React from "react";
import { useStyletron } from "baseui";
import { PromptInput } from "./PromptInput";
import { PromptHistory } from "./PromptHistory";

export function PromptBox() {
  const [css] = useStyletron();

  return (
    <div
      className={css({
        height: "100%",
        display: "flex",
        flexDirection: "column",
        padding: "16px",
        boxSizing: "border-box",
      })}
    >
      <div
        style={{
          height: "100%",
          display: "flex",
          flexDirection: "column",
        }}
      >
        <div
          className={css({ marginTop: "16px", flexGrow: 1, overflow: "auto" })}
        >
          <PromptHistory />
        </div>
        <PromptInput />
      </div>
    </div>
  );
}
