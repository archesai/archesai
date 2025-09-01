export const seo = ({
  description,
  image,
  keywords,
  title
}: {
  description?: string
  image?: string
  keywords?: string
  title: string
}) => {
  const tags = [
    { title },
    { content: description, name: 'description' },
    { content: keywords, name: 'keywords' },

    { content: title, name: 'twitter:title' },
    { content: description, name: 'twitter:description' },
    { content: 'https://www.archesai.com/', name: 'twitter:url' },
    { content: '@archesai', name: 'twitter:creator' },
    { content: '@archesai', name: 'twitter:site' },

    { content: title, name: 'og:title' },
    { content: description, name: 'og:description' },
    { content: 'https://www.archesai.com/', name: 'og:url' },
    { content: 'website', name: 'og:type' },
    ...(image ?
      [
        { content: image, name: 'twitter:image' },
        { content: '600', name: 'twitter:image:height' },
        { content: '800', name: 'twitter:image:width' },
        { content: 'Arches AI', name: 'twitter:image:alt' },
        { content: 'summary_large_image', name: 'twitter:card' },

        { content: image, name: 'og:image' },
        { content: '600', name: 'og:image:height' },
        { content: '800', name: 'og:image:width' },
        { content: 'Arches AI', name: 'og:image:alt' }
      ]
    : [])
  ]

  return tags
}
