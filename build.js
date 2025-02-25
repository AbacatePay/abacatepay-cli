const { execSync } = require('child_process');
const { copyFileSync, rmSync, renameSync, mkdir, mkdirSync} = require('fs');
const os = process.platform;

mkdirSync('build', {recursive: true});

execSync('npx esbuild index.ts --bundle --platform=node --target=es2020 --outfile=dist/index.js', { stdio: 'inherit' });
execSync('node --experimental-sea-config sea-config.json', { stdio: 'inherit' });

if(os === 'win32') {
    copyFileSync(process.execPath, 'index.exe');
    execSync('npx postject index.exe NODE_SEA_BLOB sea-prep.blob --sentinel-fuse NODE_SEA_FUSE_fce680ab2cc467b6e072b8b5df1996b2', { stdio: 'inherit' });
    renameSync('index.exe', 'build/abacate.exe');
}
else {
    execSync('cp $(command -v node) index', { stdio: 'inherit' })
    execSync('npx postject index NODE_SEA_BLOB sea-prep.blob --sentinel-fuse NODE_SEA_FUSE_fce680ab2cc467b6e072b8b5df1996b2')
    renameSync('index', 'build/abacate');
};

rmSync('sea-prep.blob');
rmSync('dist', { recursive: true });