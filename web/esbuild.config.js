const esbuild = require("esbuild");
const path = require("path");

const isWatch = process.argv.includes("--watch");

const buildOptions = {
  entryPoints: ["ts/main.ts"],
  bundle: true,
  outfile: "public/js/dist/main.js",
  format: "iife",
  globalName: "App",
  resolveExtensions: [".ts", ".js"],
  alias: {
    "@web": path.resolve(__dirname, "ts"),
    "@components": path.resolve(__dirname, "components"),
    "@domain": path.resolve(__dirname, "../domain"),
  },
  sourcemap: true,
  minify: false,
  target: "es2020",
};

if (isWatch) {
  esbuild
    .context(buildOptions)
    .then((ctx) => {
      ctx.watch();
      console.log("Watching for changes...");
    })
    .catch(() => process.exit(1));
} else {
  esbuild.build(buildOptions).catch(() => process.exit(1));
}
