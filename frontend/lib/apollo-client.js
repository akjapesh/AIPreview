import { ApolloClient, InMemoryCache, HttpLink } from "@apollo/client";

// Create a client-side only Apollo Client instance
export const apolloClient = new ApolloClient({
  link: new HttpLink({
    uri: "/api/graphql",
    // Include credentials in requests
    credentials: "same-origin",
  }),
  cache: new InMemoryCache(),
  ssrMode: typeof window === "undefined", // running on client side
  defaultOptions: {
    watchQuery: {
      fetchPolicy: "cache-and-network",
    },
  },
});
