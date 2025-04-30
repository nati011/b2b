module.exports = {
  async rewrites() {
    return {
      beforeFiles: [
        {
          source: '/random/:path',
          destination: 'http://localhost:9000/api/v1/:path*',
        },
      ],

    }
  },
}