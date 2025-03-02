import { useSubscription } from "@apollo/client";
import { useState } from "react";
import { GENERATE_RESPONSE } from "../lib/fragments/generateComponent";

export function useCodeGeneration({ prompt }) {
  const [currentCode, setCurrentCode] = useState("");
  const [shouldSubscribe, setShouldSubscribe] = useState(false);
  const [isGenerating, setIsGenerating] = useState(false);

  const { data, loading, error } = useSubscription(GENERATE_RESPONSE, {
    variables: { userId: "user_1", prompt, discussionId: "problem_1" },
    skip: !shouldSubscribe,
    onSubscriptionData: ({ subscriptionData }) => {
      console.log(subscriptionData?.type);
      if (subscriptionData.data?.generateResponse?.isComplete) {
        setIsGenerating(false);
        setShouldSubscribe(false);
      } else {
        setCurrentCode(
          (prev) => prev + subscriptionData.data?.generateComponent?.chunk
        );
      }
    },
    onError: (error) => {
      console.error("-----Subscription error:", error);
      setShouldSubscribe(false);
      setIsGenerating(false);
    },
  });

  return {
    data,
    currentCode,
    setCurrentCode,
    shouldSubscribe,
    setShouldSubscribe,
    isGenerating,
    setIsGenerating,
    loading,
    error,
  };
}
