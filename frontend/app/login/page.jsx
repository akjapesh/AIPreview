"use client";
import { Divider } from "../../components/infra/Divider";
import { signup } from "../../lib/login/auth";

export default function SignupForm() {
  const [state, action, pending] = useActionState(signup, undefined);
  return (
    <div className="my-36 border solid-1 border-black rounded-md p-4">
      <h1>Login Form</h1>
      <Divider />
      <form action={action} className="py-4">
        <div>
          <label htmlFor="name">Name</label>
          <input id="name" name="name" placeholder="Name" />
        </div>
        {state?.errors?.name && <p>{state.errors.name}</p>}
        <div>
          <label htmlFor="email">Email</label>
          <input id="email" name="email" type="email" placeholder="Email" />
        </div>
        {state?.errors?.email && <p>{state.errors.email}</p>}
        <div>
          <label htmlFor="password">Password</label>
          <input id="password" name="password" type="password" />
        </div>
        {state?.errors?.password && (
          <div>
            <p>Password must:</p>
            <ul>
              {state.errors.password.map((error) => (
                <li key={error}>- {error}</li>
              ))}
            </ul>
          </div>
        )}
        <button
          disabled={pending}
          className="px-4 py-2 mt-4 border border-black rounded-lg solid-1 bg-black text-white hover:bg-gray-700"
          type="submit"
        >
          Sign Up
        </button>
      </form>
    </div>
  );
}
