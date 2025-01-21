module.exports = {
  env: {
    AI_SDK_API_KEY: process.env.AI_SDK_API_KEY,
    POSTGRES_URL: process.env.POSTGRES_URL,
    GO_AGENT_URL: process.env.GO_AGENT_URL,
  },
  async rewrites() {
    return [
      {
        source: '/api/:path*',
        destination: 'http://localhost:8080/:path*',
      },
      {
        source: '/go-agent/:path*',
        destination: `${process.env.GO_AGENT_URL}/:path*`,
      },
    ];
  },
  serverRuntimeConfig: {
    AI_SDK_API_KEY: process.env.AI_SDK_API_KEY,
    POSTGRES_URL: process.env.POSTGRES_URL,
  },
  publicRuntimeConfig: {
    AI_SDK_API_KEY: process.env.AI_SDK_API_KEY,
    POSTGRES_URL: process.env.POSTGRES_URL,
  },
  webpack: (config, { isServer }) => {
    if (!isServer) {
      config.experiments = {
        asyncWebAssembly: true,
        layers: true,
      };
    }
    return config;
  },
};
