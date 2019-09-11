module.exports = {
  //set up reverse proxy on the dev server
  devServer: {
    proxy: {
      "^/api": {
        target: "http://localhost:8000",
        ws: true,
        changeOrigin: true
      }
    }
  },
  //for debugging (see https://vuejs.org/v2/cookbook/debugging-in-vscode.html)
  configureWebpack: {
    devtool: "source-map"
  }
};
