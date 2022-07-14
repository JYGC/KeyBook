import Axios from 'axios';
import { _ApiCall } from './_api-call';

export class GetPersonNamesAPI extends _ApiCall {
    _apiPath = `/Person/GetPersonNamesAPI`;
}

export class GetPersonTypesAPI extends _ApiCall {
    _apiPath: string = `/Person/GetPersonTypesAPI`;
}