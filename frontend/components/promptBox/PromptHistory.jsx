import React from "react";
import { useStyletron } from "baseui";
import { useQuery } from "@apollo/client";
import { GET_HISTORY } from "../../lib/fragments";
import { Spinner } from "baseui/spinner";

export function PromptHistory() {
  const [css] = useStyletron();
  //   const { data, loading, error } = useQuery(GET_HISTORY, {
  //     variables: { limit: 10, offset: 0 },
  //   });
  const { data, loading, error } = {};

  if (loading) return <Spinner />;
  if (error) return <div>Error loading history: {error.message}</div>;

  return (
    <div>
      <h3>Previous Prompts</h3>
      {data?.getHistory.map((item) => (
        <div
          key={item.id}
          className={css({
            padding: "12px",
            marginBottom: "8px",
            cursor: "pointer",
            ":hover": {
              backgroundColor: "rgba(0, 0, 0, 0.05)",
            },
          })}
        >
          {item.prompt}
        </div>
      ))}
    </div>
  );
}
