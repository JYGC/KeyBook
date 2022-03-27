import Axios from 'axios';

export function getDeviceTypes(
    successCallback: (response) => any,
    failedCallback: (any) => any
): void {
    Axios.get(`/Device/GetDeviceTypes`)
      .then(response => successCallback(response))
      .catch(e => failedCallback(e));
}

export function getPersonTypes(
    successCallback: (response) => any,
    failedCallback: (any) => any
): void {
    Axios.get(`/Person/GetPersonTypes`)
        .then(response => successCallback(response))
        .catch(e => failedCallback(e));
}