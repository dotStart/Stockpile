/*
 * Copyright 2018 Johannes Donath <johannesd@torchmind.com>
 * and other copyright owners as documented in the project's IP log.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
const path = require('path');
const HtmlWebpackPlugin = require("html-webpack-plugin");
const UglifyJsPlugin = require("uglifyjs-webpack-plugin");

const SRC_DIR = path.join(__dirname, "src");
const DIST_DIR = path.join(__dirname, "dist/ui"); // assetfs requires a root dir

console.log(' UI Input Directory: ' + SRC_DIR);
console.log('UI Output Directory: ' + DIST_DIR);

module.exports = {
  entry: path.join(SRC_DIR, "script/app.js"),
  output: {
    path: DIST_DIR,
    filename: "[name].js"
  },
  module: {
    rules: [
      {
        test: /\.js$/,
        exclude: /node_modules/,
        loader: "babel-loader"
      },
      {
        test: /\.(png|woff|woff2|eot|ttf|svg)$/,
        loader: 'url-loader?limit=100000'
      },
      {
        test: /\.css$/,
        use: [
          "css-loader"
        ]
      }
    ]
  },
  resolve: {
    alias: {
      vue: 'vue/dist/vue.js'
    }
  },
  optimization: {
    minimizer: [
      new UglifyJsPlugin({
        sourceMap: true
      })
    ]
  },
  plugins: [
    new HtmlWebpackPlugin({
      filename: 'index.html',
      template: path.join(SRC_DIR, 'index.html'),
      minify: true,
      hash: true,
      xhtml: true
    })
  ]
};
