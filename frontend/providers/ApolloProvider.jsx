// app/providers.js
"use client";

import { ApolloProvider as BaseApolloProvider } from "@apollo/client";
import { apolloClient } from "../lib/apollo-client";

export function ApolloProvider({ children }) {
  return (
    <BaseApolloProvider client={apolloClient}>{children}</BaseApolloProvider>
  );
}
