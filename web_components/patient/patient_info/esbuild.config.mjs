import { build } from 'esbuild'
import { resolve } from 'node:path'
import { fileURLToPath } from 'node:url'

const __dirname = resolve(fileURLToPath(import.meta.url), '..')

const buildOptions = {
  entryPoints: [resolve(__dirname, 'src/index.ts')],
  bundle: true,
  minify: false,
  sourcemap: true,
  target: 'es2020',
  format: 'esm',
  outfile: resolve(__dirname, 'dist/index.js'),
  platform: 'browser'
}

async function run() {
  await build(buildOptions)

  await build({
    ...buildOptions,
    format: 'cjs',
    outfile: resolve(__dirname, 'dist/index.cjs')
  })
}

run().catch(error => {
  console.error(error)
  process.exit(1)
})

