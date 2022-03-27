import Axios from 'axios';
import { _ApiCall } from './_api-call';

export class GetDeviceTypesAPI extends _ApiCall {
    _apiPath: string = `/Device/GetDeviceTypes`;
}