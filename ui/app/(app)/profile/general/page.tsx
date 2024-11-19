import UserForm from "@/components/forms/user-form";
import { getMetadata } from "@/config/site";
import { Metadata } from "next";

export const metadata: Metadata = getMetadata("/profile/general");

export default function ProfilePage() {
  return <UserForm />;
}
