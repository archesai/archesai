export const TitleAndDescription = ({
  description,
  title,
}: {
  description: string;
  title: string;
}) => {
  return (
    <div className="px-4 pt-4">
      <p className="text-xl font-semibold text-foreground">{title}</p>
      <p className="text-sm text-muted-foreground">{description}</p>
    </div>
  );
};
