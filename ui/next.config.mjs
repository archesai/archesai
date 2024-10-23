/** @type {import('next').NextConfig} */
const nextConfig = {
  experimental: {
    missingSuspenseWithCSRBailout: false,
  },
  // add https://picsum.photos/200/300 to the list of domains
  images: {
    domains: ["picsum.photos", "storage.googleapis.com", "arches-minio"],
  },
  async redirects() {
    return [
      {
        destination: "/organization/general",
        permanent: true,
        source: "/organization",
      },
      {
        destination: "/profile/general",
        permanent: true,
        source: "/profile",
      },
    ];
  },
};

export default nextConfig;
