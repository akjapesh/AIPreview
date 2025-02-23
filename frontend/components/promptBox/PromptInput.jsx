import React, { useEffect, useState } from "react";
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
  const [prompt, setPrompt] = useState("abcd");
  //   const { generateCode, loading } = useCodeGeneration();
  const { generateCode, loading } = {};
  const [getUserById, { loading: userLoading, data }] = useLazyQuery(
    GET_USER_BY_ID,
    {
      variables: { userId: "user_1" },
    }
  );

  const handleSubmit = getUserById;

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
