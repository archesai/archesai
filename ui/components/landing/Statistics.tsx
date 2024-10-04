export const Statistics = () => {
  interface statsProps {
    description: string;
    quantity: string;
  }

  const stats: statsProps[] = [
    {
      description: "Users",
      quantity: "2.7K+",
    },
    {
      description: "Subscribers",
      quantity: "1.8K+",
    },
    {
      description: "Downloads",
      quantity: "112",
    },
    {
      description: "Products",
      quantity: "4",
    },
  ];

  return (
    <section id="statistics">
      <div className="grid grid-cols-2 lg:grid-cols-4 gap-8">
        {stats.map(({ description, quantity }: statsProps) => (
          <div
            className="space-y-2 text-center"
            key={description}
          >
            <h2 className="text-3xl sm:text-4xl font-bold ">{quantity}</h2>
            <p className="text-xl text-muted-foreground">{description}</p>
          </div>
        ))}
      </div>
    </section>
  );
};
