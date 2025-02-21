import { useMutation, gql } from "@apollo/client";

export const GET_HISTORY = gql`
  mutation GetHistory {
    getHistory {
      id
      history
    }
  }
`;

export const GENERATE_COMPONENT = gql`
  mutation GenerateComponent($prompt: String!) {
    generateComponent(input: { prompt: $prompt }) {
      id
      code
      previewUrl
    }
  }
`;
