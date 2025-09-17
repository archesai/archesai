export const seo = ({
  description,
  image,
  keywords,
  title,
}: {
  description?: string;
  image?: string;
  keywords?: string;
  title: string;
}): {
  content?: string;
  name?: string;
  title?: string;
}[] => {
  const tags = [
    {
      title,
    },
    ...(description
      ? [
          {
            content: description,
            name: "description",
          },
        ]
      : []),
    ...(keywords
      ? [
          {
            content: keywords,
            name: "keywords",
          },
        ]
      : []),

    {
      content: title,
      name: "twitter:title",
    },
    ...(description
      ? [
          {
            content: description,
            name: "twitter:description",
          },
        ]
      : []),
    {
      content: "https://www.archesai.com/",
      name: "twitter:url",
    },
    {
      content: "@archesai",
      name: "twitter:creator",
    },
    {
      content: "@archesai",
      name: "twitter:site",
    },

    {
      content: title,
      name: "og:title",
    },
    ...(description
      ? [
          {
            content: description,
            name: "og:description",
          },
        ]
      : []),
    {
      content: "https://www.archesai.com/",
      name: "og:url",
    },
    {
      content: "website",
      name: "og:type",
    },
    ...(image
      ? [
          {
            content: image,
            name: "twitter:image",
          },
          {
            content: "600",
            name: "twitter:image:height",
          },
          {
            content: "800",
            name: "twitter:image:width",
          },
          {
            content: "Arches AI",
            name: "twitter:image:alt",
          },
          {
            content: "summary_large_image",
            name: "twitter:card",
          },

          {
            content: image,
            name: "og:image",
          },
          {
            content: "600",
            name: "og:image:height",
          },
          {
            content: "800",
            name: "og:image:width",
          },
          {
            content: "Arches AI",
            name: "og:image:alt",
          },
        ]
      : []),
  ];

  return tags;
};
