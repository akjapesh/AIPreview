import { useActionState, useEffect } from "react";
import { signin } from "../../lib/login/auth";
import { useRouter } from "next/navigation";

export const SignInForm = () => {
  const router = useRouter();
  const [state, action, pending] = useActionState(signin, undefined);

  useEffect(() => {
    if (state?.success) {
      console.log("Redirecting to /codePreview");
      router.push("/codePreview");
    }
  }, [state, router]);

  return (
    <form action={action} className="py-4 flex flex-col gap-2">
      <div className="flex justify-between">
        <label htmlFor="email" className="mr-2">
          Email
        </label>
        <input
          id="email"
          name="email"
          type="email"
          placeholder="enter Email"
          className="px-2 border"
        />
      </div>
      {state?.errors?.email && (
        <div className="text-red-500 text-xs">{state.errors.email}</div>
      )}
      <div className="flex justify-between">
        <label htmlFor="password" className="mr-2">
          Password
        </label>
        <input
          id="password"
          name="password"
          type="password"
          className="px-2 border"
          placeholder="enter password"
        />
      </div>
      {state?.errors?.password && (
        <div className="text-red-500 text-xs">
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
        Log In
      </button>
    </form>
  );
};
