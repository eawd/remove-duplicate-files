import * as fs from 'fs/promises';
import * as crypto from 'crypto';
import { move } from 'fs-extra';

const hashes = new Map<string, string>();

// Fill in with a folder to keep duplicate files.
const temporaryFolder = '';
// Preferred location to keep.
const desiredLocations = [
];

const folderToScan = '';

// loop over files.
// if hash doesn't exists.
//   store hash -> file path
// if hash exists
//   if the new file is in the desired location: update the hash.
//   move the undesired file to a temporary location keeping the subdirectory structure.

let counter = 0;

function isDesiredFile(filePath: string): boolean {
    for (const location of desiredLocations) {
        if (location.test(filePath)) {
            return true;
        }
    }

    return false;
}

async function removeFile(filePath: string): Promise<void> {
    console.log('Removing:', filePath);
    const newPath = `${temporaryFolder}/${filePath.replace('F:\\', '')}`;
    try {
        await move(filePath, newPath);
    } catch (err) {
        console.log(`Couldn't move ${filePath}`);
        throw err;
    }
}

async function checkFile(filePath: string): Promise<void> {
    const contents = await fs.readFile(filePath);
    const hash = crypto.createHash('md5').update(contents).digest('hex');
    if (hashes.has(hash)) {
        const previousFile = hashes.get(hash) as string;
        let fileToRemove = filePath;
        console.log('Found duplicate: ', filePath);
        if (isDesiredFile(filePath)) {
            fileToRemove = previousFile;
            hashes.set(hash, filePath);
        }

        await removeFile(fileToRemove);
    } else {
        hashes.set(hash, filePath);
    }

    if (counter++ % 200 === 0) {
        console.log(`Scanned ${counter} files, currently on ${filePath}`);
    }
}

async function scanFolder(directory: string): Promise<void> {
    const dir = await fs.opendir(directory);
    for await (const entry of dir) {
        const fullPath = await fs.realpath(`${dir.path}/${entry.name}`);
        if (entry.isDirectory()) {
            scanFolder(fullPath);
        } else {
            checkFile(fullPath);
        }
    }
}

scanFolder(folderToScan).catch((err) => {
    console.log('Error happened:', err);
}).finally(() => console.log('Done'));
