import { axios, convertError, alertError } from './_base'

export type ApplyProps = {
  type: "npm";
  name: string;
  version: string;
  v1: number;
  v2: number;
  v3: number;
}

const onSuccessApply = () => {}
const onFailedApply = (statusCode: number, data: unknown) => {
  alertError(statusCode, data);
}

const applyUnit = async (props: ApplyProps) => {
  try {
    const result = await axios.post('apply', props);
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

export const apply = (...propsArray: ApplyProps[]) => {
  propsArray.map((props)=>{
    void applyUnit(props);
  })
}
