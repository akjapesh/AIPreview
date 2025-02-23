import { gql } from "@apollo/client";

export const SIGNUP_MUTATION = gql`
  mutation SignUp($input: SignUpInput!) {
    signUp(input: $input)
  }
`;

export const SIGNIN_MUTATION = gql`
  mutation SignIn($input: SignInInput!) {
    signIn(input: $input) {
      token
      user{
        id
        email
        name
      }
    }
  }
`;
