import { CredentialsForm } from "@/features/auth/components/CredentialsForm"

export default function SigninPage() {
  const onSubmit = (creds: any) => {
    console.log(1, creds);
    
  }

  return (
    <CredentialsForm onSubmit={onSubmit} />
  )
}
