export const Statistics = () => {
  interface statsProps {
    description: string;
    quantity: string;
  }

  const stats: statsProps[] = [
    {
      description: "Users",
      quantity: "50K+",
    },
    {
      description: "Subscribers",
      quantity: "2K+",
    },
    {
      description: "Uploads",
      quantity: "200k+",
    },
    {
      description: "Messages",
      quantity: "650k+",
    },
  ];

  return (
    <section id="statistics">
      <div className="grid grid-cols-2 lg:grid-cols-4 gap-8">
        {stats.map(({ description, quantity }: statsProps) => (
          <div className="space-y-2 text-center" key={description}>
            <h2 className="text-3xl sm:text-4xl font-bold ">{quantity}</h2>
            <p className="text-xl text-muted-foreground">{description}</p>
          </div>
        ))}
      </div>
    </section>
  );
};
