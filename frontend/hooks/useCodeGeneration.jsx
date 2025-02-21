import { useMutation } from "@apollo/client";
import { useState } from "react";
import { GENERATE_COMPONENT } from "../lib/fragments";

export function useCodeGeneration() {
  const [currentCode, setCurrentCode] = useState("");
  const [currentPreviewUrl, setCurrentPreviewUrl] = useState("");

  const [generateMutation, { loading, error }] = useMutation(
    GENERATE_COMPONENT,
    {
      onCompleted: (data) => {
        setCurrentCode(data.generateComponent.code);
        setCurrentPreviewUrl(data.generateComponent.previewUrl);
      },
    }
  );

  const generateCode = async (prompt) => {
    try {
      await generateMutation({
        variables: { prompt },
      });
    } catch (err) {
      console.error("Error generating code:", err);
    }
  };

  return {
    generateCode,
    currentCode,
    currentPreviewUrl,
    loading,
    error,
  };
}
