"use client";
import { useState, useActionState } from "react";
import { Divider } from "../../components/infra/Divider";

import { SignInForm, SignUpForm } from "../../components/login";

const FORM_TYPES = {
  SIGNUP: "SIGNUP",
  SIGNIN: "SIGNIN",
};

export default function SignupForm() {
  const [formType, setFormType] = useState(FORM_TYPES.SIGNIN);

  const handleFormTypeChange = (type) => setFormType((t) => type);
  return (
    <div className="my-36 border solid-1 border-black rounded-md p-4 flex items-center flex-col">
      <div className="text-3xl">
        {formType === FORM_TYPES.SIGNIN ? "Login Form" : "Create User"}
      </div>
      <Divider />
      {formType === FORM_TYPES.SIGNIN ? <SignInForm /> : <SignUpForm />}
      {formType === FORM_TYPES.SIGNIN ? (
        <div>
          <button
            className="text-blue-500 hover:underline"
            onClick={() => handleFormTypeChange(FORM_TYPES.SIGNUP)}
          >
            Create an account
          </button>
        </div>
      ) : (
        <div>
          <button
            className="text-blue-500 hover:underline"
            onClick={() => handleFormTypeChange(FORM_TYPES.SIGNIN)}
          >
            Already have an account?
          </button>
        </div>
      )}
    </div>
  );
}
