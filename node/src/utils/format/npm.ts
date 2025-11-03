type PackageData = {
  dependencies?: {
    [key: string]: string;
  };
  devDependencies?: {
    [key: string]: string;
  };
}

function formatVersion(version: string) {
  return version.split('').filter(c => ['0' ,'1' ,'2' ,'3' ,'4' ,'5' ,'6', '7', '8', '9', '.'].includes(c)).join('')
}

type PackageList = {name: string, version: string}[]

function formatUnit(packageData: PackageData) {
  const packageList: PackageList = [];
  if (packageData.dependencies) {
    for(const [key, value] of Object.entries(packageData.dependencies)) {
      packageList.push({name: key, version: formatVersion(value)});
    }
  }
  if (packageData.devDependencies) {
    for(const [key, value] of Object.entries(packageData.devDependencies)) {
      packageList.push({name: key, version: formatVersion(value)});
    }
  }
  return packageList;
}

function readPackageFile(file: File): Promise<PackageList> {
  return new Promise<PackageList>((resolve, reject) => {
    const reader = new FileReader();

    reader.onload = (e) => {
      const content = e.target?.result as string;
      try {
        const parsedData = JSON.parse(content);
        if (parsedData.dependencies || parsedData.devDependencies) {
          resolve(formatUnit(parsedData))
        } else if (parsedData.packages) {
          if (parsedData.packages['']) {
            resolve(formatUnit(parsedData.packages['']))
          }
        } else {
          throw new Error('Invalid Format')
        }
      } catch (error) {
        reject(error)
      }
    };

    reader.onerror = () => {
      console.error("ファイルの読み込みに失敗しました。");
    };

    reader.readAsText(file, "utf-8");
  })
}

export async function formatNpm(files: FileList): Promise<PackageList>{
  const packages: PackageList = [];

  try {
    for(const file of files) {
      const loadPackages = await readPackageFile(file)
      packages.push(...loadPackages)
    }
  } catch (error) {
    alert(JSON.stringify(error))
  }
  return packages;
}