import React, { useCallback, useState } from "react";
import { useStyletron } from "baseui";
import { Textarea } from "baseui/textarea";
import { Button } from "baseui/button";
import { useCodeGeneration } from "../../hooks/useCodeGeneration";
import { gql, useLazyQuery } from "@apollo/client";

const GET_USER_BY_ID = gql`
  query GetUserById($userId: String!) {
    getUserById(userId: $userId) {
      name
      email
    }
  }
`;
export function PromptInput() {
  const [css] = useStyletron();
  const [prompt, setPrompt] = useState("");

  const {
    currentCode,
    setCurrentCode,
    shouldSubscribe,
    setShouldSubscribe,
    isGenerating,
    setIsGenerating,
    loading,
    error,
  } = useCodeGeneration({ prompt });

  const handleSubmit = useCallback(() => {
    //need to check this
    setCurrentCode("");
    setShouldSubscribe(true);
    setIsGenerating(true);
  }, []);

  const stopGeneration = useCallback(() => {
    setShouldSubscribe(false);
    setIsGenerating(false);
  }, []);

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
          disabled={!prompt.trim() || loading || isGenerating}
        >
          Generate Component
        </Button>
        {isGenerating ?? (
          <Button onClick={stopGeneration}>Stop Generation</Button>
        )}
      </div>
    </div>
  );
}
