const fs = require('fs');
const path = require('path');

const root = path.resolve(__dirname, '..');
const vendorDir = path.join(root, 'public', 'vendor');

const files = [
  {
    src: path.join(root, 'node_modules', 'htmx.org', 'dist', 'htmx.min.js'),
    dest: path.join(vendorDir, 'htmx.min.js'),
  },
];

fs.mkdirSync(vendorDir, { recursive: true });

for (const { src, dest } of files) {
  if (!fs.existsSync(src)) {
    console.warn(`Skipping copy: ${path.relative(root, src)} not found`);
    continue;
  }

  fs.copyFileSync(src, dest);
  console.log(`Copied ${path.relative(root, src)} -> ${path.relative(root, dest)}`);
}
