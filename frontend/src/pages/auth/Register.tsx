import { CredentialsForm } from "@/features/auth/ui/CredentialsForm"

export default function RegisterPage() {
  const onSubmit = (creds: any) => {
    console.log(1, creds);
  }

  return (
    <CredentialsForm onSubmit={onSubmit} />
  )
}
