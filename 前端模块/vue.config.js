module.exports = {
  devServer: {
    port: 8081, // 修改前端开发服务器端口为8081，避免与后端的8080端口冲突
    proxy: {
      '/api': {
        target: 'http://localhost:8080', // 代理到后端服务器
        changeOrigin: true
      }
    }
  },
  // 禁用生产环境的source map以减小构建体积
  productionSourceMap: false
}
