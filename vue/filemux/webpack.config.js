const path = require('path')
const webpack = require('webpack')
const { VueLoaderPlugin } = require('vue-loader')
const TerserPlugin = require('terser-webpack-plugin')

const VUE_APP_ROOT_API = process.env.VUE_APP_ROOT_API || ''

module.exports = {
  entry: './src/main.js',
  output: {
    path: path.resolve(__dirname, './dist'),
    publicPath: '/dist/',
    filename: 'build.js'
  },
  module: {
    rules: [
      {
        test: /\.vue$/,
        loader: 'vue-loader',
        options: { loaders: {} }
      },
      {
        test: /\.js$/,
        loader: 'babel-loader',
        exclude: /node_modules/
      },
      {
        test: /\.css$/,
        use: [
          'vue-style-loader',
          'css-loader'
        ],
      },
      {
        test: /\.(png|jpg|gif|svg)$/,
        loader: 'file-loader',
        options: { name: '[name].[ext]?[hash]' }
      }
    ]
  },
  resolve: {
    alias: {
      'vue$': 'vue/dist/vue.esm.js',
      '@': path.resolve(__dirname, 'src')
    }
  },
  devServer: {
    hot: true,
    proxy: process.env.NODE_DEV_PROXY_API === undefined ? {} : {
      [`${VUE_APP_ROOT_API}/ws/events`]: {
        target: `${process.env.NODE_DEV_PROXY_API}/ws/events`,
        ws: true
      }
    }
  },
  devtool: '#eval-source-map',
  plugins: [ new VueLoaderPlugin() ]
}

module.exports.plugins = (module.exports.plugins || []).concat([
  new webpack.EnvironmentPlugin({ VUE_APP_ROOT_API })
])

if (process.env.NODE_ENV === 'production') {
  module.exports.devtool = '#source-map'

  module.exports.optimization = {
    minimizer: [
      new TerserPlugin({
        cache: true,
        parallel: true,
        terserOptions: {
          compress: false,
          ecma: 6,
          mangle: true
        },
        sourceMap: true
      })
    ]
  }
}
