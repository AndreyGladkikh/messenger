import type { Credentials } from "@/features/auth/model/credentials";
import { CredentialsForm } from "./CredentialsForm";

export function RegisterForm() {
  const onSubmit = (creds: Credentials) => {
    console.log(1, creds);

  }

  return <CredentialsForm onSubmit={onSubmit} />
}
