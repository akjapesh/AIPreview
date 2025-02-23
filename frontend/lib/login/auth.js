"use server";

import { SignupFormSchema, SigninFormSchema } from "./definitions";

//helpers
import { apolloServerClient } from "../apollo-client";
import { cookies } from "next/headers";

//fragments
import { SIGNUP_MUTATION, SIGNIN_MUTATION } from "../fragments/auth";

export async function signup(state, formData) {
  // Validate form fields
  const validatedFields = SignupFormSchema.safeParse({
    name: formData.get("name"),
    email: formData.get("email"),
    password: formData.get("password"),
  });

  // If any form fields are invalid, return early
  if (!validatedFields.success) {
    return {
      errors: validatedFields.error.flatten().fieldErrors,
    };
  }

  // Call the provider or db to create a user...
  // send to go-backend with apollo client
  try {
    const { data } = await apolloServerClient.mutate({
      mutation: SIGNUP_MUTATION,
      variables: { input: validatedFields.data },
      context: {
        headers: {
          "X-Auth-Optional": "true",
        },
      },
    });

    return { success: true, user: data.signUp };
  } catch (error) {
    console.error("Signup failed:", error);
    return { errors: { form: error.message } };
  }
}

export async function signin(state, formData) {
  // Validate form fields
  const validatedFields = SigninFormSchema.safeParse({
    email: formData.get("email"),
    password: formData.get("password"),
  });

  // If any form fields are invalid, return early
  if (!validatedFields.success) {
    return {
      errors: validatedFields.error.flatten().fieldErrors,
    };
  }

  // Call the provider or db to create a user...
  // send to go-backend with apollo client
  try {
    const { data } = await apolloServerClient.mutate({
      mutation: SIGNIN_MUTATION,
      variables: { input: validatedFields.data },
      context: {
        headers: {
          "X-Auth-Optional": "true",
        },
      },
    });

    // Store the token in a cookie
    cookies().set("auth_token", data.signIn.token, {
      httpOnly: true,
      secure: process.env.NODE_ENV === "production",
      sameSite: "Strict",
      maxAge: 60, // 1 day
    });

    return { success: true, user: data.signIn };
  } catch (error) {
    console.error("Signin failed:", error);
    return { errors: { form: error.message } };
  }
}
