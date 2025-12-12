import { Link } from "@tanstack/react-router";

export const TermsIndicator: React.FC = () => {
  return (
    <div className="max-w-sm text-balance text-center text-muted-foreground text-xs *:[a]:underline *:[a]:underline-offset-4 *:[a]:hover:text-primary">
      By clicking continue, you agree to our{" "}
      <Link to="/">Terms of Service</Link> and{" "}
      <Link to="/">Privacy Policy</Link>.
    </div>
  );
};
