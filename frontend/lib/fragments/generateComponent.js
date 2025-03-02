import { gql } from "@apollo/client";

export const GET_HISTORY = gql`
  mutation GetHistory($userId: String!, $discussionId: String!) {
    getHistory(input: { userId: $userId, discussionId: $discussionId }) {
      id
      history
    }
  }
`;

export const GENERATE_RESPONSE = gql`
  subscription GenerateResponse(
    $userId: String!
    $discussionId: String!
    $prompt: String!
  ) {
    generateResponse(
      userId: $userId
      discussionId: $discussionId
      prompt: $prompt
    ) {
      files {
        name
        content
      }
      chunk
      isComplete
    }
  }
`;
