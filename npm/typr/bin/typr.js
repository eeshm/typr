#!/usr/bin/env node

const path = require('node:path');
const fs = require('node:fs');
const { spawn } = require('node:child_process');

const exeName = process.platform === 'win32' ? 'typr.exe' : 'typr';
const exePath = path.join(__dirname, exeName);

if (!fs.existsSync(exePath)) {
  console.error('[typr] Binary not found. Reinstall the package to download platform binary.');
  process.exit(1);
}

const child = spawn(exePath, process.argv.slice(2), {
  stdio: 'inherit',
});

child.on('error', (err) => {
  console.error('[typr] Failed to launch binary:', err.message);
  process.exit(1);
});

child.on('exit', (code, signal) => {
  if (signal) {
    process.kill(process.pid, signal);
    return;
  }
  process.exit(code ?? 0);
});
