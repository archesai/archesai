/** @type {import('next').NextConfig} */
const nextConfig = {
  experimental: {
    missingSuspenseWithCSRBailout: false,
  },
  // add https://picsum.photos/200/300 to the list of domains
  images: {
    domains: ["picsum.photos", "storage.googleapis.com", "arches-minio"],
  },
  images: {
    unoptimized: true,
  },
  output: "export",
  async redirects() {
    return [
      {
        destination: "/settings/organization/general",
        permanent: true,
        source: "/settings/organization",
      },
      {
        destination: "/settings/organization/general",
        permanent: true,
        source: "/settings",
      },
      {
        destination: "/import/file",
        permanent: true,
        source: "/import",
      },
    ];
  },
};

export default nextConfig;
