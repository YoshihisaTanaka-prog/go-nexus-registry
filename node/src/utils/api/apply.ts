import { axios, convertError, alertError } from './_base'

export type ApplyProps = {
  name: string;
  v1?: number;
  v2?: number;
  v3?: number;
}

const onSuccessApply = () => {}
const onFailedApply = (statusCode: number, data: unknown) => {
  alertError(statusCode, data);
}

const applyUnit = async (libKind: string, index: number, props: ApplyProps) => {
  console.log(props)
  try {
    const result = await axios.post('apply', {...props, index, kind: libKind});
    if (result.status < 400) {
      onSuccessApply();
    } else {
      onFailedApply(result.status, result.data);
    }
  } catch (error) {
    const { status, data } = convertError(error);
    onFailedApply(status, data)
  }
}

export const apply = (libKind: string|undefined, ...propsArray: ApplyProps[]) => {
  if (libKind === undefined) {
    alert('ライブラリの種類を選択してください。')
  } else {
    console.log(propsArray)
    propsArray.map((props, index)=>{
      void applyUnit(libKind, index, props);
    });
  }
}
