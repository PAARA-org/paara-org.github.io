// script/prepare-vendor.js
// Copies vendor assets from node_modules into assets/vendor/ for local serving and deployment.
const fs = require('fs');
const path = require('path');

const rootDir = path.resolve(__dirname, '..');
const vendorDir = path.join(rootDir, 'assets', 'vendor', 'glightbox');
const nodeDistDir = path.join(rootDir, 'node_modules', 'glightbox', 'dist');

if (!fs.existsSync(nodeDistDir)) {
  console.error('Error: glightbox dist not found in node_modules. Run "npm install" first.');
  process.exit(1);
}

fs.mkdirSync(vendorDir, { recursive: true });

const filesToCopy = [
  {
    src: path.join(nodeDistDir, 'js', 'glightbox.min.js'),
    dest: path.join(vendorDir, 'glightbox.min.js')
  },
  {
    src: path.join(nodeDistDir, 'css', 'glightbox.min.css'),
    dest: path.join(vendorDir, 'glightbox.min.css')
  }
];

for (const file of filesToCopy) {
  if (fs.existsSync(file.src)) {
    fs.copyFileSync(file.src, file.dest);
    console.log(`✓ Copied ${path.basename(file.src)} -> assets/vendor/glightbox/`);
  } else {
    console.error(`Error: Source file not found: ${file.src}`);
    process.exit(1);
  }
}
