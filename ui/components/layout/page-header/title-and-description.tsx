export const TitleAndDescription = ({
  description,
  Icon,
  title,
}: {
  description?: string;
  Icon: any;
  title?: string;
}) => {
  if (!title) return null;
  return (
    <div className="flex items-center gap-2 px-4 pt-4">
      {Icon && <Icon className="-ml-1 h-8 w-8 text-primary" />}
      <div>
        <p className="text-xl font-semibold text-foreground">{title}</p>
        <p className="text-sm text-muted-foreground">{description}</p>
      </div>
    </div>
  );
};
