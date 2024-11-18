export const TitleAndDescription = ({
  description,
  title,
}: {
  description?: string;
  Icon: any;
  title?: string;
}) => {
  if (!title) return null;
  return (
    <div className="flex items-center gap-3 border-b px-4 py-4">
      {/* {Icon && <Icon className="h-6 w-6 text-accent-foreground" />} */}
      <div>
        <p className="text-xl font-semibold text-foreground">{title}</p>
        <p className="text-sm text-muted-foreground">{description}</p>
      </div>
    </div>
  );
};
