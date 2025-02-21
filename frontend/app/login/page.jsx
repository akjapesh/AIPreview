"use client";
import { Divider } from "../../components/infra/Divider";

export default function SignupForm() {
  return (
    <div className="my-36 border solid-1 border-black rounded-md p-4">
      <h1>Login Form</h1>
      <Divider />
      <form onSubmit={() => console.log("submit")} className="py-4">
        <div>
          <label htmlFor="name">Name</label>
          <input id="name" name="name" placeholder="Name" />
        </div>
        <div>
          <label htmlFor="email">Email</label>
          <input id="email" name="email" type="email" placeholder="Email" />
        </div>
        <div>
          <label htmlFor="password">Password</label>
          <input id="password" name="password" type="password" />
        </div>
        <button
          className="px-4 py-2 mt-4 border border-black rounded-lg solid-1 bg-black text-white hover:bg-gray-700"
          type="submit"
        >
          Sign Up
        </button>
      </form>
    </div>
  );
}
