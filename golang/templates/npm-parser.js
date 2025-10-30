const { packages } = require('./package-lock.json')

const mainLibrary = {name: '', version: '', numOfLibraries: 0};
const subLibraries = [];

for (const [key, val] of Object.entries(packages)) {
  if (key === '') {
    const { dependencies } = val;
    for (const [k, v] of Object.entries(dependencies)) {
      mainLibrary.name = k;
      mainLibrary.version = v;
      if (mainLibrary.version.startsWith('^')) {
        mainLibrary.version = mainLibrary.version.slice(1);
      }
    }
  } else {
    const [libName] = key.split('node_modules/').slice(-1);
    const { version, resolved } = val;
    subLibraries.push({name: libName, version, resolved});
  }
}

mainLibrary.numOfLibraries = subLibraries.length;

console.log(
  JSON.stringify({ mainLibrary, subLibraries })
);