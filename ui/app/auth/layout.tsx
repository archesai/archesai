"use client";
import { LogoSVG } from "@/components/logo-svg";
export default function AuthLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <div className="flex flex-col h-screen items-center justify-center bg-gray-100 ">
      <div className="pb-6 -mt-14">
        <LogoSVG scale={0.25} />
      </div>
      <div>{children}</div>
    </div>
  );
}
