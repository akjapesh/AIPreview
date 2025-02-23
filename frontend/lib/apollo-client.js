import { ApolloClient, InMemoryCache, HttpLink } from "@apollo/client";
import { onError } from "@apollo/client/link/error";

const httpLink = new HttpLink({
  uri: "http://localhost:8080/graphql",
  credentials: "include", // 🔥 This ensures cookies are sent with requests
});

// ✅ Error handling middleware to detect 401 Unauthorized responses
const errorLink = onError(({ graphQLErrors, networkError }) => {
  if (graphQLErrors) {
    for (const err of graphQLErrors) {
      if (err.extensions?.code === "UNAUTHENTICATED") {
        console.warn("Session expired. Redirecting to login...");
        window.location.href = "/login"; // ✅ Redirect user to login page
      }
    }
  }

  if (networkError && networkError.statusCode === 401) {
    console.warn("Unauthorized request. Redirecting to login...");
    window.location.href = "/login"; // ✅ Redirect user to login page
  }
});

// Create a client-side only Apollo Client instance
export const apolloClient = new ApolloClient({
  link: errorLink.concat(httpLink),
  cache: new InMemoryCache(),
  ssrMode: false, // running on client side
  defaultOptions: {
    watchQuery: {
      fetchPolicy: "cache-and-network",
    },
  },
});

export const apolloServerClient = new ApolloClient({
  link: new HttpLink({
    uri: "http://localhost:8080/graphql",
    // Include credentials in requests
    credentials: "include",
  }),
  cache: new InMemoryCache(),
  ssrMode: true, // running on server side only
  defaultOptions: {
    watchQuery: {
      fetchPolicy: "network-only",
    },
  },
});
