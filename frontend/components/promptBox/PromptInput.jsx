import React, { useState } from "react";
import { useStyletron } from "baseui";
import { Textarea } from "baseui/textarea";
import { Button } from "baseui/button";
import { useCodeGeneration } from "../../hooks/useCodeGeneration";

export function PromptInput() {
  const [css] = useStyletron();
  const [prompt, setPrompt] = useState("");
  //   const { generateCode, loading } = useCodeGeneration();
  const { generateCode, loading } = {};

  const handleSubmit = async () => {
    if (prompt.trim()) {
      await generateCode(prompt);
      setPrompt("");
    }
  };

  return (
    <div>
      <Textarea
        value={prompt}
        onChange={(e) => setPrompt(e.target.value)}
        placeholder="Describe the component you want to generate..."
        rows={4}
      />
      <div className={css({ marginTop: "16px" })}>
        <Button
          onClick={handleSubmit}
          isLoading={loading}
          disabled={!prompt.trim() || loading}
        >
          Generate Component
        </Button>
      </div>
    </div>
  );
}
