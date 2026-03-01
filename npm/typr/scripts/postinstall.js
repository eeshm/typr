const fs = require('node:fs');
const path = require('node:path');
const https = require('node:https');
const AdmZip = require('adm-zip');

const pkg = require('../package.json');

const PLATFORM_MAP = {
  win32: 'windows',
  darwin: 'darwin',
  linux: 'linux',
};

const ARCH_MAP = {
  x64: 'amd64',
  arm64: 'arm64',
};

function requestBuffer(url, redirectCount = 0) {
  return new Promise((resolve, reject) => {
    const req = https.get(url, {
      headers: {
        'User-Agent': 'typr-npm-installer',
      },
    }, (res) => {
      if (res.statusCode >= 300 && res.statusCode < 400 && res.headers.location) {
        if (redirectCount > 5) {
          reject(new Error('Too many redirects while downloading release asset'));
          return;
        }
        resolve(requestBuffer(res.headers.location, redirectCount + 1));
        return;
      }

      if (res.statusCode !== 200) {
        reject(new Error(`Download failed (${res.statusCode}) for ${url}`));
        return;
      }

      const chunks = [];
      res.on('data', (chunk) => chunks.push(chunk));
      res.on('end', () => resolve(Buffer.concat(chunks)));
    });

    req.on('error', reject);
  });
}

function failWithHint(message) {
  console.error(`[typr] ${message}`);
  console.error('[typr] You can still download binaries manually from GitHub Releases.');
  process.exit(1);
}

async function main() {
  const platform = PLATFORM_MAP[process.platform];
  const arch = ARCH_MAP[process.arch];

  if (!platform || !arch) {
    failWithHint(`Unsupported platform/arch: ${process.platform}/${process.arch}`);
  }

  const repo = process.env.TYPR_REPO || pkg.typr?.repo;
  if (!repo) {
    failWithHint('Missing repository slug (expected pkg.typr.repo)');
  }

  const version = process.env.TYPR_VERSION || `v${pkg.version}`;
  const assetName = `typr_${version}_${platform}_${arch}.zip`;
  const url = `https://github.com/${repo}/releases/download/${version}/${assetName}`;

  console.log(`[typr] Downloading ${assetName}...`);
  const zipBuffer = await requestBuffer(url);

  const binDir = path.resolve(__dirname, '../bin');
  fs.mkdirSync(binDir, { recursive: true });

  const zip = new AdmZip(zipBuffer);
  zip.extractAllTo(binDir, true);

  const exeName = process.platform === 'win32' ? 'typr.exe' : 'typr';
  const exePath = path.join(binDir, exeName);

  if (!fs.existsSync(exePath)) {
    failWithHint(`Installed archive did not contain ${exeName}`);
  }

  if (process.platform !== 'win32') {
    fs.chmodSync(exePath, 0o755);
  }

  console.log('[typr] Installed successfully.');
}

main().catch((err) => {
  failWithHint(err.message || String(err));
});
