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

const applyUnit = async (libType: string, props: ApplyProps) => {
  console.log(props)
  try {
    const result = await axios.post('apply', {...props, type: libType});
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

export const apply = (libType: string|undefined, ...propsArray: ApplyProps[]) => {
  if (libType === undefined) {
    alert('ライブラリの種類を選択してください。')
  } else {
    console.log(propsArray)
    propsArray.map((props)=>{
      void applyUnit(libType, props);
    });
  }
}
